package checkawsec2mainte

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/ec2iface"
)

type EC2Mainte struct {
	Client      ec2iface.EC2API
	InstanceIds []string
}

// GetMainteInfo ... Call API and get specified events
func (mt EC2Mainte) GetMainteInfo(ctx context.Context) (events EC2Events, err error) {
	options := &ec2.DescribeInstanceStatusInput{}

	// If InstanceIds is empty, get all EC2 Events
	if len(mt.InstanceIds) != 0 {
		options.InstanceIds = mt.InstanceIds
	}

	req := mt.Client.DescribeInstanceStatusRequest(options)
	req.SetContext(ctx)
	res, err := req.Send()
	if err != nil {
		return
	}

	// Create EC2 Events from InstanceStatusResponse
	for _, instance := range res.InstanceStatuses {
		if len(instance.Events) != 0 {
			for _, ev := range instance.Events {
				events = append(events, EC2Event{
					Code:        ev.Code,
					InstanceId:  *instance.InstanceId,
					NotAfter:    *ev.NotAfter,
					NotBefore:   *ev.NotBefore,
					Description: *ev.Description,
				})
			}
		}
	}

	// Remove already completed events
	// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_InstanceStatusEvent.html
	events = events.Filter("[Completed]")
	events = events.Filter("[Canceled]")
	return
}
