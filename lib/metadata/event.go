package metadata

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// Event ...
type Event struct {
	Code        ec2.EventCode  `json:"Code"`
	NotBefore   MetaMainteTime `json:"NotBefore"`
	NotAfter    MetaMainteTime `json:"NotAfter"`
	Description string         `json:"Description"`
	State       string         `json:"State"`
}

// Events ...
type Events []Event

// SpotEventAction ...
type SpotEventAction string

// SpotEventAction ...
const (
	SpotEventActionHibernate SpotEventAction = "hibernate"
	SpotEventActionStop      SpotEventAction = "stop"
	SpotEventActionTerminate SpotEventAction = "terminate"
)

// SpotEvent ...
type SpotEvent struct {
	Action SpotEventAction `json:"action"`
	Time   time.Time       `json:"time"`
}
