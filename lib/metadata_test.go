package checkawsec2mainte_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	checkawsec2mainte "github.com/ntrv/check-aws-ec2-mainte/lib"
	"github.com/ntrv/check-aws-ec2-mainte/lib/unit"
	"github.com/stretchr/testify/assert"
)

func TestGetInstanceIdFromMetadata(t *testing.T) {
	expected := "i-09e032cce9ef71d84"

	server := unit.StartTestServer(map[string]string{
		"/latest/meta-data/instance-id": expected,
	})
	defer server.Close()

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		t.Error(err.Error())
	}

	cfg.EndpointResolver = aws.ResolveWithEndpointURL(server.URL + "/latest")

	actual, err := checkawsec2mainte.GetInstanceIdFromMetadata(cfg)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, expected, actual)
}
