package unit

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	checkawsec2mainte "github.com/ntrv/check-aws-ec2-mainte/lib"
)

type testCaseMainte struct {
	Resp     ec2.DescribeInstanceStatusOutput
	Expected checkawsec2mainte.Events
}

func CreateCases(t *testing.T) []testCaseMainte {
	ds := CreateTimes(t, []string{
		"2019-03-14T16:04:05+09:00",
		"2019-03-16T16:04:05+09:00",
		"2019-03-16T18:04:05+09:00",
	})

	return []testCaseMainte{
		{
			Resp: ec2.DescribeInstanceStatusOutput{
				InstanceStatuses: []ec2.InstanceStatus{
					{
						InstanceId: aws.String("i-0472b8a82f226da14"),
						Events: []ec2.InstanceStatusEvent{
							{
								Code:        ec2.EventCodeSystemMaintenance,
								Description: aws.String("Scheduled System Maintenance"),
								NotBefore:   aws.Time(ds[0]),
								NotAfter:    aws.Time(ds[1]),
							},
						},
					},
					{
						InstanceId: aws.String("i-0dc818ea941b1ae18"),
						Events: []ec2.InstanceStatusEvent{
							{
								Code:        ec2.EventCodeInstanceRetirement,
								Description: aws.String("Scheduled Instance Retirement Maintenance"),
								NotBefore:   aws.Time(ds[2]),
								NotAfter:    aws.Time(ds[2]), //TODO: Check this field is required
							},
						},
					},
					{
						InstanceId: aws.String("i-0dc818ea941b1ae18"),
						Events: []ec2.InstanceStatusEvent{
							{
								Code:        ec2.EventCodeInstanceReboot,
								Description: aws.String("[Completed] Scheduled Instance Reboot Maintenance"),
								NotBefore:   aws.Time(ds[0]),
								NotAfter:    aws.Time(ds[1]),
							},
						},
					},
				},
			},
			Expected: checkawsec2mainte.Events{
				{
					Code:        ec2.EventCodeSystemMaintenance,
					InstanceId:  "i-0472b8a82f226da14",
					Description: "Scheduled System Maintenance",
					NotBefore:   ds[0],
					NotAfter:    ds[1],
					State:       checkawsec2mainte.StateActive,
				},
				{
					Code:        ec2.EventCodeInstanceRetirement,
					InstanceId:  "i-0dc818ea941b1ae18",
					Description: "Scheduled Instance Retirement Maintenance",
					NotBefore:   ds[2],
					NotAfter:    ds[2], //TODO: Check this field is required
					State:       checkawsec2mainte.StateActive,
				},
				{
					Code:        ec2.EventCodeInstanceReboot,
					InstanceId:  "i-0dc818ea941b1ae18",
					Description: "[Completed] Scheduled Instance Reboot Maintenance",
					NotBefore:   ds[0],
					NotAfter:    ds[1],
					State:       checkawsec2mainte.StateCompleted,
				},
			},
		},
	}
}
