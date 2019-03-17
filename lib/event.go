package checkawsec2mainte

import (
	"fmt"
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

func (self EC2Event) CreateMessage() string {
	// Load Location from $TZ or /etc/localtime
	return fmt.Sprintf(
		"Code: %v, InstanceId: %v, Date: %v - %v, Description: %v",
		self.Code,
		self.InstanceId,
		self.NotBefore.In(time.Local).Format(time.RFC3339),
		self.NotAfter.In(time.Local).Format(time.RFC3339),
		self.Description,
	)
}
