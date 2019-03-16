package checkawsec2mainte

import (
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/ec2iface"
)

type IEC2Events interface {
	GetCloseEvent() EC2Mainte
	Len() int
}

type EC2Events []EC2Event

func (e EC2Events) Filter(substr string) EC2Events {
	events := EC2Maintes{}

	for _, event := range e {
		if strings.Contains(event.Description, substr) {
			continue
		}
		events = append(events, event)
	}

	return events
}

func GetMainteInfo(svc ec2iface.EC2API, instanceIds ...string) (EC2Events, error) {

	maintes := EC2Events{}

	options := &ec2.DescribeInstanceStatusInput{}
	if len(instanceIds) != 0 {
		options.InstanceIds = instanceIds
	}

	req := svc.DescribeInstanceStatusRequest(options)

	res, err := req.Send()
	if err != nil {
		return nil, err
	}

	for idx, event := range res.InstanceStatuses[0].Events {
		maintes[idx].Code = event.Code
		maintes[idx].NotAfter = *event.NotAfter
		maintes[idx].NotBefore = *event.NotBefore
		maintes[idx].Description = *event.Description
	}

	sort.Stable(maintes)
	return maintes.Filter("Completed"), nil
}

func (self EC2Events) GetCloseEvent() EC2Events {
	return self[len(self)-1]
}

func (self EC2Events) Len() int {
	return len(self)
}

func (self EC2Events) Less(i, j int) bool {
	return self[i].NotBefore.Unix() < self[j].NotBefore.Unix()
}

func (self EC2Events) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}
