package checkawsec2mainte

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/stretchr/testify/assert"
)

func createTime(t *testing.T, value string) time.Time {
	d, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func createTimes(t *testing.T, values []string) (ds []time.Time) {
	for _, v := range values {
		ds = append(ds, createTime(t, v))
	}
	return
}

func createEvents(t *testing.T) EC2Events {

	ds := createTimes(t, []string{
		"2019-03-14T16:04:05+09:00",
		"2019-03-16T16:04:05+09:00",
		"2019-03-16T18:04:05+09:00",
		"2019-03-17T17:34:35+09:00",
		"2019-03-17T18:04:05+09:00",
		"2019-03-17T18:04:05+07:00",
	})

	return EC2Events{
		{
			Code:        ec2.EventCodeSystemReboot,
			InstanceId:  "i-9263d590",
			NotBefore:   ds[2],
			Description: "scheduled reboot",
		},
		{
			Code:        ec2.EventCodeSystemMaintenance,
			InstanceId:  "i-07bfa293eacde7019",
			NotBefore:   ds[0], // Closest Event
			NotAfter:    ds[1],
			Description: "[Completed] Scheduled System Maintenance",
		},
		{
			Code:        ec2.EventCodeInstanceRetirement,
			InstanceId:  "i-05d9be9a",
			NotBefore:   ds[5],
			Description: "[Completed] Scheduled Instance Retirement Maintenance",
		},
		{
			Code:        ec2.EventCodeInstanceReboot,
			InstanceId:  "i-0f456b937f33abe9e",
			NotBefore:   ds[3],
			NotAfter:    ds[4],
			Description: "Scheduled Instance Reboot Maintenance",
		},
	}
}

func TestFilterCompleted(t *testing.T) {
	events := createEvents(t)

	events = events.Filter("Completed")
	assert.Len(t, events, 2)
}

func TestGetCloseEvent(t *testing.T) {
	events := createEvents(t) // Contains completed events
	event := events.GetCloseEvent()
	assert.Equal(t, events[1], event)
}

func TestBeforeAllisOk(t *testing.T) {
	ds := createTimes(t, []string{
		"2019-03-14T16:04:05+09:00",
		"2019-03-16T16:04:05+09:00",
		"2019-03-16T18:04:05+09:00",
		"2019-03-17T17:34:35+09:00",
		"2019-03-17T18:04:05+07:00",
	})

	now := createTime(t, "2019-03-20T19:00:00+09:00") // Future

	events := EC2Events{
		{
			NotBefore: ds[0],
		},
		{
			NotBefore: ds[1],
		},
		{
			NotBefore: ds[2],
		},
		{
			NotBefore: ds[3],
		},
		{
			NotBefore: ds[4],
		},
	}

	assert.True(t, events.BeforeAll(now))
}

func TestBeforeAllisNotOk(t *testing.T) {
	ds := createTimes(t, []string{
		"2019-03-14T16:04:05+09:00",
		"2019-03-16T16:04:05+09:00",
		"2019-03-16T18:04:05+09:00",
		"2019-03-17T17:34:35+09:00",
		"2019-03-17T18:04:05+07:00",
	})

	now := createTime(t, "2019-03-16T19:00:00+09:00") // intermediate

	events := EC2Events{
		{
			NotBefore: ds[0],
		},
		{
			NotBefore: ds[1],
		},
		{
			NotBefore: ds[2],
		},
		{
			NotBefore: ds[3],
		},
		{
			NotBefore: ds[4],
		},
	}

	assert.False(t, events.BeforeAll(now))
	now = createTime(t, "2019-02-13T19:00:00+09:00") // Far Past
	assert.False(t, events.BeforeAll(now))
}
