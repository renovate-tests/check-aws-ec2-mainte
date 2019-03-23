package checkawsec2mainte_test

import (
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/ntrv/check-aws-ec2-mainte/lib"

	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
)

func TestCheckerIsCritical(t *testing.T) {
	events := createEvents(t)

	c, err := checkawsec2mainte.NewChecker([]string{
		"-c", "1000h",
		"-r", "us-west-1",
	})
	if err != nil {
		t.Error(err)
	}
	c.Now = createTime(t, "2019-03-14T12:23:12+09:00")

	ckr := c.Run(events)
	assert.Equal(t, checkers.CRITICAL, ckr.Status)
}

func TestCheckerIsWarning(t *testing.T) {
	events := createEvents(t)

	c, err := checkawsec2mainte.NewChecker([]string{
		"-c", "1m",
		"-r", "us-west-1",
	})
	if err != nil {
		t.Error(err)
	}
	c.Now = createTime(t, "2019-03-14T12:23:12+09:00")

	ckr := c.Run(events)
	assert.Equal(t, checkers.WARNING, ckr.Status)
}

func TestCheckerIsOk(t *testing.T) {
	events := checkawsec2mainte.EC2Events{}
	assert.Len(t, events, 0)
	assert.Zero(t, events.Len())

	c, err := checkawsec2mainte.NewChecker([]string{
		"-r", "us-west-1",
	})
	if err != nil {
		t.Error(err)
	}
	c.Now = createTime(t, "2019-03-14T12:23:12+09:00")

	ckr := c.Run(events)
	assert.Equal(t, checkers.OK, ckr.Status)
}

func TestOverCheckerIsCritical(t *testing.T) {
	events := createEvents(t)

	c, err := checkawsec2mainte.NewChecker([]string{
		"-c", "1m",
		"-r", "us-west-1",
	})
	if err != nil {
		t.Error(err)
	}
	c.Now = createTime(t, "2019-03-20T12:23:12+09:00")
	pp.Println(c)

	ckr := c.Run(events)
	pp.Println(ckr)

	assert.Equal(t, checkers.CRITICAL, ckr.Status)
}
