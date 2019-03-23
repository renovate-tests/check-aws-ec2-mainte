package checkawsec2mainte

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2EventState string

const (
	StateActive    EC2EventState = "active"
	StateCompleted EC2EventState = "completed"
	StateCanceled  EC2EventState = "canceled"
)

// EC2Event ... Almost same as ec2.InstanceStatusEvent
type EC2Event struct {
	Code        ec2.EventCode `json:"Code"`
	InstanceId  string        `json:"-"`
	NotBefore   time.Time     `json:"NotBefore"`
	NotAfter    time.Time     `json:"NotAfter"`
	Description string        `json:"Description"`
	State       EC2EventState `json:"State"`
}

// IsTimeOver ... EC2Eventが引数より新しいかどうか
func (ev EC2Event) IsTimeOver(now time.Time, d time.Duration) bool {
	return now.Add(d).After(ev.NotBefore)
}

// CreateMessage ... Information for displaying to Mackerel
func (ev EC2Event) CreateMessage() string {
	// Load Location from $TZ or /etc/localtime
	return fmt.Sprintf(
		"Code: %v, InstanceId: %v, Date: %v - %v, Description: %v",
		ev.Code,
		ev.InstanceId,
		ev.NotBefore.In(time.Local).Format(time.RFC3339),
		ev.NotAfter.In(time.Local).Format(time.RFC3339),
		ev.Description,
	)
}
