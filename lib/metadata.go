package checkawsec2mainte

import (
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

func getInstanceIdFromMetadata() (string, error) {
	metadata := ec2metadata.New(session.New())
	d, err := metadata.GetInstanceIdentityDocument()
	if err != nil {
		return "", err
	}
	return d.InstanceID, nil
}
