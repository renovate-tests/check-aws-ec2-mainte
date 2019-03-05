package checkawsec2mainte

import (
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type IEC2Mainte interface {
	GetCloseEvent() EC2Mainte
	Length() int
}

func NewEC2Mainte(svc ec2iface.EC2API, instanceIds ...string) (IEC2Mainte, error) {

	var maintes EC2Maintes

	events, err := maintes.GetMainteInfo(svc, instanceIds...)
	if err != nil {
		return nil, err
	}

	for i, event := range events {
		maintes[i].NotAfter = *event.NotAfter
		maintes[i].NotBefore = *event.NotBefore
		maintes[i].Description = *event.Description
	}

	sort.Sort(maintes)
	return maintes, nil
}

func (_ EC2Maintes) getInstanceIdFromMetadata() (string, error) {
	metadata := ec2metadata.New(session.New())
	d, err := metadata.GetInstanceIdentityDocument()
	if err != nil {
		return "", err
	}
	return d.InstanceID, nil
}

func (self EC2Maintes) GetMainteInfo(svc ec2iface.EC2API, instanceIds ...string) (
	[]*ec2.InstanceStatusEvent,
	error,
) {

	if len(instanceIds) == 0 {
		instanceId, err := self.getInstanceIdFromMetadata()
		if err != nil {
			return nil, err
		}
		instanceIds = append(instanceIds, instanceId)
	}

	d, err := svc.DescribeInstanceStatus(&ec2.DescribeInstanceStatusInput{
		InstanceIds: aws.StringSlice(instanceIds),
	})
	if err != nil {
		return nil, err
	}

	return d.InstanceStatuses[0].Events, nil
}

func (self EC2Maintes) GetCloseEvent() EC2Mainte {
	return self[len(self)-1]
}

func (self EC2Maintes) Length() int {
	return len(self)
}
