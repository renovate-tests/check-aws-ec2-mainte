package checkawsec2mainte

import (
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"
)

func getInstanceIdFromMetadata() string {

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return ""
	}

	cfg.HTTPClient = &http.Client{
		Timeout: 100 * time.Millisecond,
	}

	m := ec2metadata.New(cfg)

	id, err := m.GetMetadata("instance-id")
	if err != nil {
		return ""
	}

	return id
}
