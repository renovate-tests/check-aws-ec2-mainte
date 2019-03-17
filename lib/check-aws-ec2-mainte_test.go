package checkawsec2mainte

import (
	"testing"

	"github.com/mackerelio/checkers"

	"github.com/stretchr/testify/assert"
	"github.com/k0kubun/pp"
)

func TestCheckerIsCritical(t *testing.T) {
	events := createEvents(t)

	c, err := NewChecker([]string{
		"-c", "1000h",
		"-r", "us-west-1",
	})
	if err != nil {
		t.Error(err)
	}
	c.Now = createTime(t, "2019-03-18T12:23:12+09:00")
	pp.Println(c)

	ckr := c.run(events)
	pp.Println(ckr)

	assert.Equal(t, checkers.CRITICAL, ckr.Status)
}

func TestCheckerIsWarning(t *testing.T) {
	events := createEvents(t)

	c, err := NewChecker([]string{
		"-c", "1m",
		"-r", "us-west-1",
	})
	if err != nil {
		t.Error(err)
	}
	c.Now = createTime(t, "2019-03-18T12:23:12+09:00")
	pp.Println(c)

	ckr := c.run(events)
	pp.Println(ckr)

	assert.Equal(t, checkers.WARNING, ckr.Status)
}
