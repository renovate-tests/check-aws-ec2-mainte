package checkawsec2mainte

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2Event struct {
	Code        ec2.EventCode
	InstanceId  string
	NotBefore   time.Time
	NotAfter    time.Time
	Description string
}

func (self EC2Event) IsTimeOver(now time.Time, d time.Duration) bool {
	return now.Add(d).After(self.NotBefore)
}
