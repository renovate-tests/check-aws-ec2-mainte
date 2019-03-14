package checkawsec2mainte

import (
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/ec2iface"
)

type IEC2Mainte interface {
	GetCloseEvent() EC2Mainte
	Len() int
}

type EC2Maintes []EC2Mainte

func NewEC2Mainte(svc ec2iface.EC2API, instanceIds ...string) (IEC2Mainte, error) {

	var maintes EC2Maintes

	// if len(instanceIds) == 0 {
	// 	instanceId, err := getInstanceIdFromMetadata()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	instanceIds = append(instanceIds, instanceId)
	// }

	events, err := maintes.GetMainteInfo(svc, instanceIds...)
	if err != nil {
		return nil, err
	}

	for i, event := range events {
		if strings.Contains(maintes[i].Description, "Completed") {
			continue
		}

		maintes[i].NotAfter = *event.NotAfter
		maintes[i].NotBefore = *event.NotBefore
		maintes[i].Description = *event.Description
	}

	sort.Stable(maintes)
	return maintes, nil
}

func (self EC2Maintes) GetMainteInfo(svc ec2iface.EC2API, instanceIds ...string) (
	[]ec2.InstanceStatusEvent,
	error,
) {

	req := svc.DescribeInstanceStatusRequest(&ec2.DescribeInstanceStatusInput{
		InstanceIds: instanceIds,
	})

	res, err := req.Send()
	if err != nil {
		return nil, err
	}

	return res.InstanceStatuses[0].Events, nil
}

func (self EC2Maintes) GetCloseEvent() EC2Mainte {
	return self[len(self)-1]
}

func (self EC2Maintes) Len() int {
	return len(self)
}

func (self EC2Maintes) Less(i, j int) bool {
	return self[i].NotBefore.Unix() < self[j].NotBefore.Unix()
}

func (self EC2Maintes) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}
