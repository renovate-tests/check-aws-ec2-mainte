package checkawsec2mainte

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
)

// Get Instance ID from http://169.254.169.254/latest/meta-data/instance-id
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

// GetMaintesFromMetadata ... Get Scheduled Maintenances
func GetMaintesFromMetadata(cfg aws.Config) (events EC2Events, err error) {
	cfg.HTTPClient = &http.Client{
		Timeout: 100 * time.Millisecond,
	}

	m := ec2metadata.New(cfg)

	data, err := m.GetMetadata("events/maintenance/scheduled")
	if err != nil {
		return
	}

	if err = json.Unmarshal([]byte(data), &events); err != nil {
		return
	}

	instanceId, err := GetInstanceIdFromMetadata(cfg)
	if err != nil {
		return
	}

	for _, event := range events {
		event.InstanceId = instanceId
	}
	return
}
