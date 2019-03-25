package checkawsec2mainte_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"

	checkawsec2mainte "github.com/ntrv/check-aws-ec2-mainte/lib"
	"github.com/ntrv/check-aws-ec2-mainte/lib/unit"
)

func initMetaConfig(t *testing.T, endpoint string) aws.Config {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		t.Error(err.Error())
	}
	cfg.EndpointResolver = aws.ResolveWithEndpointURL(endpoint + "/latest")
	return cfg
}

func TestGetInstanceIdFromMetadata(t *testing.T) {
	expected := "i-09e032cce9ef71d84"

	server := unit.StartTestServer(map[string]string{
		"/latest/meta-data/instance-id": expected,
	})
	defer server.Close()

	mt := checkawsec2mainte.EC2MetaMainte{
		Client: ec2metadata.New(initMetaConfig(t, server.URL)),
	}

	actual, err := mt.GetInstanceId(context.Background())
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, expected, actual)
}

func TestEventsFromMetadata(t *testing.T) {
	expectedId := "i-09e032cce9ef71d84"
	expected := unit.CreateEventsMetadata(t, expectedId)

	data, _ := json.Marshal(expected)

	server := unit.StartTestServer(map[string]string{
		"/latest/meta-data/instance-id":                  expectedId,
		"/latest/meta-data/events/maintenance/scheduled": string(data),
	})
	defer server.Close()

	mt := checkawsec2mainte.EC2MetaMainte{
		Client: ec2metadata.New(initMetaConfig(t, server.URL)),
	}

	actual, err := mt.GetEvents(context.Background())
	if err != nil {
		t.Error(err.Error())
	}

	pp.Println(actual)
	assert.Equal(t, expected, actual)
}
