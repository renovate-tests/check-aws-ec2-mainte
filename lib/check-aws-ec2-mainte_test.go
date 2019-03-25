package checkawsec2mainte_test

import (
	"encoding/json"
	"testing"

	"github.com/k0kubun/pp"
	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"

	checkawsec2mainte "github.com/ntrv/check-aws-ec2-mainte/lib"
	"github.com/ntrv/check-aws-ec2-mainte/lib/unit"
)

func TestCheckerIsCriticalViaCLI(t *testing.T) {
	events := unit.CreateEvents(t)

	c, err := checkawsec2mainte.NewChecker([]string{
		"-c", "1000h",
		"-r", "us-west-1",
	})
	if err != nil {
		t.Error(err)
	}
	c.Now = unit.CreateTime(t, "2019-03-14T12:23:12+09:00")

	ckr := c.Run(events)
	assert.Equal(t, checkers.CRITICAL, ckr.Status)
}

func TestCheckerIsWarningViaCLI(t *testing.T) {
	events := unit.CreateEvents(t)

	c, err := checkawsec2mainte.NewChecker([]string{
		"-c", "1m",
		"-r", "us-west-1",
	})
	if err != nil {
		t.Error(err)
	}
	c.Now = unit.CreateTime(t, "2019-03-14T12:23:12+09:00")

	ckr := c.Run(events)
	assert.Equal(t, checkers.WARNING, ckr.Status)
}

func TestCheckerIsOkViaCLI(t *testing.T) {
	events := checkawsec2mainte.Events{}
	assert.Len(t, events, 0)
	assert.Zero(t, events.Len())

	c, err := checkawsec2mainte.NewChecker([]string{
		"-r", "us-west-1",
	})
	if err != nil {
		t.Error(err)
	}
	c.Now = unit.CreateTime(t, "2019-03-14T12:23:12+09:00")

	ckr := c.Run(events)
	assert.Equal(t, checkers.OK, ckr.Status)
}

func TestOverCheckerIsCriticalViaCLI(t *testing.T) {
	events := unit.CreateEvents(t)

	c, err := checkawsec2mainte.NewChecker([]string{
		"-c", "1m",
		"-r", "us-west-1",
	})
	if err != nil {
		t.Error(err)
	}
	c.Now = unit.CreateTime(t, "2019-03-20T12:23:12+09:00")
	pp.Println(c)

	ckr := c.Run(events)
	pp.Println(ckr)

	assert.Equal(t, checkers.CRITICAL, ckr.Status)
}

func TestEventsFromMetadataViaCLI(t *testing.T) {
	expectedId := "i-09e032cce9ef71d84"
	expected := unit.CreateEventsMetadata(t, expectedId)

	data, _ := json.Marshal(expected)

	server := unit.StartTestServer(map[string]string{
		"/latest/meta-data/instance-id":                  expectedId,
		"/latest/meta-data/events/maintenance/scheduled": string(data),
	})
	defer server.Close()

	c, err := checkawsec2mainte.NewChecker([]string{
		"-c", "1m",
	})
	if err != nil {
		t.Error(err)
	}
	c.Now = unit.CreateTime(t, "2019-03-14T12:23:12+09:00")

	ckr := c.Run(expected)
	pp.Println(ckr)

	assert.Equal(t, checkers.WARNING, ckr.Status)
}
