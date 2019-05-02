package checkawsec2mainte

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// EventState ...
type EventState string

// State ...
const (
	StateActive    EventState = "active"
	StateCompleted EventState = "completed"
	StateCanceled  EventState = "canceled"
)

// Event ... Almost same as ec2.InstanceStatusEvent
type Event struct {
	Code        ec2.EventCode `json:"Code"`
	InstanceID  string        `json:"-"`
	NotBefore   time.Time     `json:"NotBefore"`
	NotAfter    time.Time     `json:"NotAfter"`
	Description string        `json:"Description"`
	State       EventState    `json:"State"`
}

// IsTimeOver ... EC2Eventが引数より新しいかどうか
func (ev Event) IsTimeOver(now time.Time, d time.Duration) bool {
	return now.Add(d).After(ev.NotBefore)
}

// CreateMessage ... Information for displaying to Mackerel
func (ev Event) CreateMessage() string {
	// Load Location from $TZ or /etc/localtime
	return fmt.Sprintf(
		"Code: %v, InstanceId: %v, TimeRange: %v to %v, Description: %v",
		ev.Code,
		ev.InstanceID,
		ev.NotBefore.In(time.Local).Format(time.RFC3339),
		ev.NotAfter.In(time.Local).Format(time.RFC3339),
		ev.Description,
	)
}
