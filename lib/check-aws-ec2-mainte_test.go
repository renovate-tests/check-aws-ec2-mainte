package checkawsec2mainte

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"

	"github.com/k0kubun/pp"
)

func TestCheckerCritical(t *testing.T) {

	ds1 := createTimes(t, []string{
		"2019-03-14T16:04:05+09:00",
		"2019-03-16T18:04:05+09:00",
	})

	ds2 := createTimes(t, []string{
		"2019-03-16T16:04:05+09:00",
		"2019-03-17T03:34:51+09:00",
	})

	events := EC2Events{
		{
			Code:        ec2.EventCodeSystemMaintenance,
			InstanceId:  "i-0dc818ea941b1ae18",
			NotBefore:   ds1[0],
			NotAfter:    ds1[1],
			Description: "Scheduled System Maintenance",
		},
		{
			Code:        ec2.EventCodeInstanceRetirement,
			InstanceId:  "i-0f456b937f33abe9e",
			NotBefore:   ds2[0],
			NotAfter:    ds2[1],
			Description: "Scheduled Instance Retirement Maintenance",
		},
	}

	c, err := NewChecker([]string{
		"-c", "1h",
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
