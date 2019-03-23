package checkawsec2mainte

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// EC2Event ... Almost same as ec2.InstanceStatusEvent
type EC2Event struct {
	Code        ec2.EventCode `json:"Code"`
	InstanceId  string        `json:"-"`
	NotBefore   time.Time     `json:"NotBefore"`
	NotAfter    time.Time     `json:"NotAfter"`
	Description string        `json:"Description"`
}

// IsTimeOver ... EC2Eventが引数より新しいかどうか
func (self EC2Event) IsTimeOver(now time.Time, d time.Duration) bool {
	return now.Add(d).After(self.NotBefore)
}

// CreateMessage ... Information for displaying to Mackerel
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
