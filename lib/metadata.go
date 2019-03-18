package checkawsec2mainte

import (
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
)

func GetInstanceIdFromMetadata(cfg aws.Config) (string, error) {
	cfg.HTTPClient = &http.Client{
		Timeout: 100 * time.Millisecond,
	}

	m := ec2metadata.New(cfg)

	id, err := m.GetMetadata("instance-id")
	if err != nil {
		return "", err
	}

	return id, nil
}
