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

// getEvents ... Get Scedule from metadata
func (mm *Mainte) getEvents(ctx context.Context) (mevs Events, err error) {
	data, err := mm.Client.GetMetadata("events/maintenance/scheduled")
	if err != nil {
		return
	}
	if err = json.Unmarshal([]byte(data), &mevs); err != nil {
		return
	}
	return
}

// Fetch ... Get Scheduled Maintenances
func (mm *Mainte) Fetch(ctx context.Context) (evs events.Events, err error) {
	mevs, err := mm.getEvents(ctx)
	if err != nil {
		return
	}

	instanceID, err := mm.getInstanceID(ctx)
	if err != nil {
		return
	}
	for _, mev := range mevs {
		evs = append(evs, events.Event{
			Code:        mev.Code,
			InstanceID:  instanceID,
			NotBefore:   time.Time(mev.NotBefore),
			NotAfter:    time.Time(mev.NotAfter),
			Description: mev.Description,
			State:       events.EventState(mev.State),
		})
	}
	return
}

// FetchSpotEvent ...
func (mm *Mainte) FetchSpotEvent(ctx context.Context) (sev SpotEvent, err error) {
	data, err := mm.Client.GetMetadata("spot/instance-action")
	if err != nil {
		return
	}
	if err = json.Unmarshal([]byte(data), &sev); err != nil {
		return
	}
	return
}
