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

// getInstanceID ... Get Instance ID from http://169.254.169.254/latest/meta-data/instance-id
func (mm *Mainte) getInstanceID(ctx context.Context) (string, error) {
	return mm.Client.GetMetadata("instance-id")
}

// GetEvents ... Get Scheduled Maintenances
func (mm *Mainte) GetEvents(ctx context.Context) (evs Events, err error) {
	data, err := mm.Client.GetMetadata("events/maintenance/scheduled")
	if err != nil {
		return
	}
	if err = json.Unmarshal([]byte(data), &evs); err != nil {
		return
	}

	instanceID, err := mm.getInstanceID(ctx)
	if err != nil {
		return
	}
	for i := range evs {
		evs[i].InstanceID = instanceID
	}
	return
}
