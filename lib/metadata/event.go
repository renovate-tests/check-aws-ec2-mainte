package metadata

import (
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
