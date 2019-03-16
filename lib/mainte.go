package checkawsec2mainte

import (
	"sort"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/ec2iface"
)

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
