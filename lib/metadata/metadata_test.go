package metadata

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"

	"github.com/stretchr/testify/assert"
)

var testCases = []struct {
	fname      string
	instanceId string
}{
	{
		fname:      "./testdata/case1.json",
		instanceId: "i-09e032cce9ef71d84",
	},
	{
		fname:      "./testdata/case2.json",
		instanceId: "i-98fa339d",
	},
	{
		fname:      "./testdata/case3.json",
		instanceId: "i-0342eeba4f394a064",
	},
}

func initTestServer(patterns map[string]string) *httptest.Server {
	return httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if resp, ok := patterns[r.RequestURI]; ok {
				w.Write([]byte(resp))
				return
			}
			http.Error(w, "not found", http.StatusNotFound)
		}),
	)
}

func initMetaConfig(t *testing.T, endpoint string) aws.Config {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		t.Error(err.Error())
	}
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(endpoint + "/latest")
	return cfg
}

func readTestCase(t *testing.T, filename string, instanceId string) (events Events, data []byte) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error(err)
	}

	if err := json.Unmarshal(data, &events); err != nil {
		t.Error(err)
	}

	if instanceId != "" {
		// Create expected events to append instance id
		for i, _ := range events {
			events[i].InstanceId = instanceId
		}
	}
	return
}

func TestMarshalNoError(t *testing.T) {
	for _, c := range testCases {
		t.Run(filepath.Base(c.fname), func(t *testing.T) {
			events, _ := readTestCase(t, c.fname, c.instanceId)

			_, err := json.MarshalIndent(events, "", " ")
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGetEvents(t *testing.T) {
	for _, c := range testCases {
		t.Run(filepath.Base(c.fname), func(t *testing.T) {
			expected, data := readTestCase(t, c.fname, c.instanceId)

			// Create dummy metadata endpoint
			server := initTestServer(map[string]string{
				"/latest/meta-data/events/maintenance/scheduled": string(data),
				"/latest/meta-data/instance-id":                  c.instanceId,
			})
			defer server.Close()

			cfg := initMetaConfig(t, server.URL)
			mt := Mainte{
				Client: ec2metadata.New(cfg),
			}

			actual, err := mt.GetEvents(context.Background())
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, expected, actual)
		})
	}
}

func TestGetInstanceId(t *testing.T) {
	for _, c := range testCases {
		t.Run(filepath.Base(c.fname), func(t *testing.T) {
			// Create dummy metadata endpoint
			server := initTestServer(map[string]string{
				"/latest/meta-data/instance-id": c.instanceId,
			})
			defer server.Close()

			cfg := initMetaConfig(t, server.URL)
			mt := Mainte{
				Client: ec2metadata.New(cfg),
			}

			actual, err := mt.GetInstanceId(context.Background())
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, c.instanceId, actual)
		})
	}
}
