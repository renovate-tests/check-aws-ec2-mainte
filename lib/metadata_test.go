package checkawsec2mainte_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/ntrv/check-aws-ec2-mainte/lib"
	"github.com/stretchr/testify/assert"
)

func initTestServer(path string, resp string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != path {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Write([]byte(resp))
	}))
}

func TestGetInstanceIdFromMetadata(t *testing.T) {
	ast := assert.New(t)
	expected := "i-09e032cce9ef71d84"

	server := initTestServer(
		"/latest/meta-data/instance-id",
		expected,
	)
	defer server.Close()

	cfg, err := external.LoadDefaultAWSConfig()
	ast.NoError(err)

	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL + "/latest")

	actual, err := checkawsec2mainte.GetInstanceIdFromMetadata(cfg)
	ast.NoError(err)

	ast.Equal(expected, actual)
}
