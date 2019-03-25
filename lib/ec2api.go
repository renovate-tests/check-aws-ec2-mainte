package checkawsec2mainte

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/ec2iface"
)

type EC2Mainte struct {
	Client      ec2iface.EC2API
	InstanceIds []string
}

// Call API and get specified events
func (mt EC2Mainte) GetEvents(ctx context.Context) (events Events, err error) {
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
				events = append(events, Event{
					Code:        ev.Code,
					InstanceId:  aws.StringValue(instance.InstanceId),
					NotAfter:    aws.TimeValue(ev.NotAfter),
					NotBefore:   aws.TimeValue(ev.NotBefore),
					Description: aws.StringValue(ev.Description),
				})
			}
		}
	}

	// Parse Descriptions and Set States
	events.UpdateStates()
	return
}
