package checkawsec2mainte

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/defaults"
)

type jsonErrorResponse struct {
	Code    string `json:"__type"`
	Message string `json:"message"`
}

func newMockConfig() aws.Config {
	config := defaults.Config()
	config.Region = "mock-region"
	config.EndpointResolver = aws.ResolveWithEndpointURL("https://endpoint")
	config.Credentials = aws.StaticCredentialsProvider{
		Value: aws.Credentials{
			AccessKeyID: "AKID",
			SecretAccessKey: "SECRET",
			SessionToken: "SESSION",
			Source: "unit test credentials",
		},
	}
	return config.Copy()
}

func newMockClient(cfg aws.Config) *aws.Client {
	if cfg.Retryer == nil {
		cfg.Retryer = aws.DefaultRetryer{NumMaxRetries: 3}
	}
	return aws.NewClient(
		cfg,
		aws.Metadata{
			ServiceName: "mockService",
		},
	)
}

func body(str string) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader([]byte(str)))
}

func NewAwsMockRequest(data interface{}) *aws.Request {
	cfg := newMockConfig()
	s := newMockClient(cfg)

	s.Handlers.Validate.Clear()
	s.Handlers.Unmarshal.PushBack(func(req *aws.Request) {
		defer req.HTTPResponse.Body.Close()
		if req.Data != nil {
			json.NewDecoder(req.HTTPResponse.Body).Decode(req.Data)
		}
	})
	s.Handlers.UnmarshalError.PushBack(unmarshalError)
	s.Handlers.Send.Clear() // mock sending
	s.Handlers.Send.PushBack(func(r *aws.Request) {
		r.HTTPResponse = &http.Response{
			StatusCode: http.StatusOK,
			Body:       body(`{"data":"valid"}`),
		}
	})

	return s.NewRequest(
		&aws.Operation{Name: "Operation"},
		nil,
		data,
	)
}

func unmarshalError(req *aws.Request) {
	bodyBytes, err := ioutil.ReadAll(req.HTTPResponse.Body)
	if err != nil {
		req.Error = awserr.New("UnmarshaleError", req.HTTPResponse.Status, err)
		return
	}
	if len(bodyBytes) == 0 {
		req.Error = awserr.NewRequestFailure(
			awserr.New(
				"UnmarshaleError",
				req.HTTPResponse.Status,
				errors.New("empty body"),
			),
			req.HTTPResponse.StatusCode,
			"",
		)
		return
	}
	var jsonErr jsonErrorResponse
	if err := json.Unmarshal(bodyBytes, &jsonErr); err != nil {
		req.Error = awserr.New("UnmarshaleError", "JSON unmarshal", err)
		return
	}
	req.Error = awserr.NewRequestFailure(
		awserr.New(jsonErr.Code, jsonErr.Message, nil),
		req.HTTPResponse.StatusCode,
		"",
	)
}
