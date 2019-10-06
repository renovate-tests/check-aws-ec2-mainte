package metadata

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
)

// Mainte ...
type Mainte struct {
	Client *ec2metadata.Client
}

// GetInstanceId ... Get Instance ID from http://169.254.169.254/latest/meta-data/instance-id
func (mm *Mainte) GetInstanceId(ctx context.Context) (string, error) {
	id, err := mm.Client.GetMetadata("instance-id")
	if err != nil {
		return "", err
	}
	return id, nil
}

// GetEvents ... Get Scheduled Maintenances
func (mm *Mainte) GetEvents(ctx context.Context) (events Events, err error) {
	data, err := mm.Client.GetMetadata("events/maintenance/scheduled")
	if err != nil {
		return
	}
	if err = json.Unmarshal([]byte(data), &events); err != nil {
		return
	}

	instanceId, err := mm.GetInstanceId(ctx)
	if err != nil {
		return
	}
	for i := range events {
		events[i].InstanceId = instanceId
	}
	return
}
