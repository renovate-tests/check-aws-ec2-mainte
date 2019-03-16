package checkawsec2mainte

import (
	"sort"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/ec2iface"
)

type EC2Mainte struct {
	Client      ec2iface.EC2API
	instanceIds []string
}

func (e EC2Mainte) GetMainteInfo() (events EC2Events, err error) {
	options := &ec2.DescribeInstanceStatusInput{}
	if len(e.instanceIds) != 0 {
		options.InstanceIds = e.instanceIds
	}

	req := e.Client.DescribeInstancesStatusRequest(options)
	res, err := req.Send()
	if err != nil {
		return
	}

	for i, e := range res.InstanceStatuses[0].Events {
		events[i].Code = e.Code
		events[i].NotAfter = *e.NotAfter
		events[i].NotBefore = *e.NotBefore
		events[i].Description = *e.Description
	}

	sort.Stable(events)
	events = events.Filter("Completed")
	return
}
