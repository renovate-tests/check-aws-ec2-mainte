package ec2api

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type Event struct {
	Code        ec2.EventCode `json:"Code"`
	InstanceId  string        `json:"-"`
	NotBefore   time.Time     `json:"NotBefore"`
	NotAfter    time.Time     `json:"NotAfter"`
	Description string        `json:"Description"`
}

type Events []Event
