package checkawsec2mainte_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/ec2iface"
	"github.com/k0kubun/pp"
	checkawsec2mainte "github.com/ntrv/check-aws-ec2-mainte/lib"
	"github.com/ntrv/check-aws-ec2-mainte/lib/unit"
	"github.com/stretchr/testify/assert"
)

type mockEc2Svc struct {
	ec2iface.EC2API
	Resp ec2.DescribeInstanceStatusOutput
}

func (m mockEc2Svc) DescribeInstanceStatusRequest(input *ec2.DescribeInstanceStatusInput) ec2.DescribeInstanceStatusRequest {
	r := unit.NewAwsMockRequest(&m.Resp)

	return ec2.DescribeInstanceStatusRequest{
		Request: r,
	}
}

func TestGetMaintesFromAPI(t *testing.T) {
	cases := unit.CreateCases(t)

	for _, c := range cases {
		mt := checkawsec2mainte.EC2Mainte{
			Client:      mockEc2Svc{Resp: c.Resp},
			InstanceIds: []string{},
		}
		events, err := mt.GetMainteInfo(context.Background())
		if err != nil {
			t.Error(err)
		}

		pp.Println(events)

		assert.Equal(t, c.Expected, events)
	}
}
