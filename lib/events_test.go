package checkawsec2mainte_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	checkawsec2mainte "github.com/ntrv/check-aws-ec2-mainte/lib"
	"github.com/ntrv/check-aws-ec2-mainte/lib/unit"
)

func TestFilterCompleted(t *testing.T) {
	events := unit.CreateEvents(t)

	events = events.Filter(checkawsec2mainte.StateCompleted)
	assert.Len(t, events, 2)
}

func TestGetCloseEvent(t *testing.T) {
	events := unit.CreateEvents(t) // Contains completed events
	event := events.GetCloseEvent()
	assert.Equal(t, events[1], event)
}

func TestBeforeAllisOk(t *testing.T) {
	ds := unit.CreateTimes(t, []string{
		"2019-03-14T16:04:05+09:00",
		"2019-03-16T16:04:05+09:00",
		"2019-03-16T18:04:05+09:00",
		"2019-03-17T17:34:35+09:00",
		"2019-03-17T18:04:05+07:00",
	})

	now := unit.CreateTime(t, "2019-03-20T19:00:00+09:00") // Future

	events := checkawsec2mainte.Events{
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
	ds := unit.CreateTimes(t, []string{
		"2019-03-14T16:04:05+09:00",
		"2019-03-16T16:04:05+09:00",
		"2019-03-16T18:04:05+09:00",
		"2019-03-17T17:34:35+09:00",
		"2019-03-17T18:04:05+07:00",
	})

	now := unit.CreateTime(t, "2019-03-16T19:00:00+09:00") // intermediate

	events := checkawsec2mainte.Events{
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
	now = unit.CreateTime(t, "2019-02-13T19:00:00+09:00") // Far Past
	assert.False(t, events.BeforeAll(now))
}
