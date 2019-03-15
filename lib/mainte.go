package checkawsec2mainte

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2Mainte struct {
	Code        ec2.EventCode
	NotBefore   time.Time
	NotAfter    time.Time
	Description string
}

func (self EC2Mainte) IsTimeOver(now time.Time, d time.Duration) bool {
	return now.Add(d).After(self.NotBefore)
}
