package metadata

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/ntrv/check-aws-ec2-mainte/lib/events"
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

// Fetch ..
func (mm *Mainte) Fetch(ctx context.Context) (evs events.Events, err error) {
	mevs, err := mm.GetEvents(ctx)
	if err != nil {
		return
	}

	for _, mev := range mevs {
		evs = append(evs, events.Event{
			Code:        mev.Code,
			InstanceID:  mev.InstanceID,
			NotBefore:   time.Time(mev.NotBefore),
			NotAfter:    time.Time(mev.NotAfter),
			Description: mev.Description,
			State:       events.EventState(mev.State),
		})
	}
	return
}
