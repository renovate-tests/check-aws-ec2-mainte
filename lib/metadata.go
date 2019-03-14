package checkawsec2mainte

import (
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

func getInstanceIdFromMetadata() string {
	metadata := ec2metadata.New(session.New(), &aws.Config{
		HTTPClient: &http.Client{Timeout: 100 * time.Millisecond},
	})

	if !metadata.Available() {
		return ""
	}

	d, err := metadata.GetInstanceIdentityDocument()
	if err != nil {
		return ""
	}
	return d.InstanceID
}

func getRegionFromMetadata() string {
	metadata := ec2metadata.New(session.New(), &aws.Config{
		HTTPClient: &http.Client{Timeout: 100 * time.Millisecond},
	})

	if !metadata.Available() {
		return ""
	}

	region, err := metadata.Region()
	if err != nil {
		return ""
	}
	return region
}
