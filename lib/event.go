package checkawsec2mainte

import (
	"fmt"
	"log"
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

func (self EC2Event) CreateMessage(location string) string {

	loc, err := time.LoadLocation(location)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return fmt.Sprintf(
		"Code: %v, InstanceId: %v, Date: %v - %v, Description: %v",
		self.Code,
		self.InstanceId,
		self.NotBefore.In(loc).Format(time.RFC3339),
		self.NotAfter.In(loc).Format(time.RFC3339),
		self.Description,
	)
}
