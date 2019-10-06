package metadata

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"

	"github.com/stretchr/testify/assert"

	"github.com/ntrv/check-aws-ec2-mainte/lib/internal/test"
)

var testCases = []struct {
	fname      string
	instanceID string
}{
	{
		fname:      "./testdata/single.json",
		instanceID: "i-09e032cce9ef71d84",
	},
	{
		fname:      "./testdata/multi.json",
		instanceID: "i-98fa339d",
	},
	{
		fname:      "./testdata/empty.json",
		instanceID: "i-0342eeba4f394a064",
	},
}

func initMetaConfig(t *testing.T, endpoint string) aws.Config {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		t.Error(err.Error())
	}
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(endpoint + "/latest")
	return cfg
}

func readTestCase(t *testing.T, filename string, instanceID string) (events Events, data []byte) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error(err)
	}

	if err := json.Unmarshal(data, &events); err != nil {
		t.Error(err)
	}

	if instanceID != "" {
		// Create expected events to append instance id
		for i := range events {
			events[i].InstanceId = instanceID
		}
	}
	return
}

func TestMarshalNoError(t *testing.T) {
	for _, c := range testCases {
		t.Run(filepath.Base(c.fname), func(t *testing.T) {
			events, expected := readTestCase(t, c.fname, c.instanceID)

			actual, err := json.MarshalIndent(events, "", " ")
			if err != nil {
				t.Error(err)
			}

			t.Skip("TODO: Should match JSON from file and Marshaled struct")
			assert.JSONEq(t, string(expected), string(actual))
		})
	}
}

func TestGetEvents(t *testing.T) {
	for _, c := range testCases {
		t.Run(filepath.Base(c.fname), func(t *testing.T) {
			expected, data := readTestCase(t, c.fname, c.instanceID)

			// Create dummy metadata endpoint
			server := test.InitTestServer(map[string]string{
				"/latest/meta-data/events/maintenance/scheduled": string(data),
				"/latest/meta-data/instance-id":                  c.instanceID,
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
			server := test.InitTestServer(map[string]string{
				"/latest/meta-data/instance-id": c.instanceID,
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
			assert.Equal(t, c.instanceID, actual)
		})
	}
}
