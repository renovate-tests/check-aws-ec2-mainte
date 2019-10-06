package ec2api

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/ec2iface"
	"github.com/ntrv/check-aws-ec2-mainte/lib/events"
)

// Mainte ...
type Mainte struct {
	Client      ec2iface.ClientAPI
	InstanceIds []string
}

func parseState(desc string) events.EventState {
	switch {
	case strings.Contains(desc, "[Completed]"):
		return events.StateCompleted
	case strings.Contains(desc, "[Canceled]"):
		return events.StateCanceled
	default:
		return events.StateActive
	}
}

// Fetch ... Call API and get specified events
func (mt Mainte) Fetch(ctx context.Context) (evs events.Events, err error) {
	options := &ec2.DescribeInstanceStatusInput{}

	// If InstanceIds is empty, get all EC2 Events
	if len(mt.InstanceIds) != 0 {
		options.InstanceIds = mt.InstanceIds
	}

	req := mt.Client.DescribeInstanceStatusRequest(options)
	res, err := req.Send(ctx)
	if err != nil {
		return
	}

	// Create EC2 Events from InstanceStatusResponse
	for _, instance := range res.InstanceStatuses {
		if len(instance.Events) != 0 {
			for _, ev := range instance.Events {
				evs = append(evs, events.Event{
					Code:        ev.Code,
					InstanceID:  aws.StringValue(instance.InstanceId),
					NotAfter:    aws.TimeValue(ev.NotAfter),
					NotBefore:   aws.TimeValue(ev.NotBefore),
					Description: aws.StringValue(ev.Description),
					State:       parseState(aws.StringValue(ev.Description)),
				})
			}
		}
	}
	return
}
