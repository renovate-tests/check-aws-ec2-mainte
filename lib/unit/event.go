package unit

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	checkawsec2mainte "github.com/ntrv/check-aws-ec2-mainte/lib"
)

func CreateTime(t *testing.T, value string) time.Time {
	d, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func CreateTimes(t *testing.T, values []string) (ds []time.Time) {
	for _, v := range values {
		ds = append(ds, CreateTime(t, v))
	}
	return
}

func CreateEvents(t *testing.T) checkawsec2mainte.Events {
	ds := CreateTimes(t, []string{
		"2019-03-14T16:04:05+09:00",
		"2019-03-16T16:04:05+09:00",
		"2019-03-16T18:04:05+09:00",
		"2019-03-17T17:34:35+09:00",
		"2019-03-17T18:04:05+09:00",
		"2019-03-17T18:04:05+07:00",
	})

	return checkawsec2mainte.Events{
		{
			Code:        ec2.EventCodeSystemReboot,
			InstanceId:  "i-9263d590",
			NotBefore:   ds[2],
			Description: "scheduled reboot",
			State:       checkawsec2mainte.StateActive,
		},
		{
			Code:        ec2.EventCodeSystemMaintenance,
			InstanceId:  "i-07bfa293eacde7019",
			NotBefore:   ds[0], // Closest Event
			NotAfter:    ds[1],
			Description: "[Completed] Scheduled System Maintenance",
			State:       checkawsec2mainte.StateCompleted,
		},
		{
			Code:        ec2.EventCodeInstanceRetirement,
			InstanceId:  "i-05d9be9a",
			NotBefore:   ds[5],
			Description: "[Completed] Scheduled Instance Retirement Maintenance",
			State:       checkawsec2mainte.StateCompleted,
		},
		{
			Code:        ec2.EventCodeInstanceReboot,
			InstanceId:  "i-0f456b937f33abe9e",
			NotBefore:   ds[3],
			NotAfter:    ds[4],
			Description: "Scheduled Instance Reboot Maintenance",
			State:       checkawsec2mainte.StateActive,
		},
	}
}
