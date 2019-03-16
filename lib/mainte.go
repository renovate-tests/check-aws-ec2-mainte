package checkawsec2mainte

import (
	"sort"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/ec2iface"
)

type EC2Mainte struct {
	Client      ec2iface.EC2API
	InstanceIds []string
}

func (e EC2Mainte) GetMainteInfo() (events EC2Events, err error) {
	options := &ec2.DescribeInstanceStatusInput{}

	// If InstanceIds is empty, get all EC2 Events
	if len(e.InstanceIds) != 0 {
		options.InstanceIds = e.InstanceIds
	}

	req := e.Client.DescribeInstanceStatusRequest(options)
	res, err := req.Send()
	if err != nil {
		return
	}

	// Create EC2 Events from InstanceStatusResponse
	for _, i := range res.InstanceStatuses {
		if len(i.Events) != 0 {
			for _, e := range i.Events {
				events = append(events, EC2Event{
					Code:        e.Code,
					InstanceId:  *i.InstanceId,
					NotAfter:    *e.NotAfter,
					NotBefore:   *e.NotBefore,
					Description: *e.Description,
				})
			}
		}
	}

	// Sort as NotBefore date
	sort.Stable(events)

	// Remove already completed events
	// https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_InstanceStatusEvent.html
	events = events.Filter("Completed")
	return
}
