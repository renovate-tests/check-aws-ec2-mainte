package checkawsec2mainte

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2Mainte ec2.InstanceStatusEvent

func (self EC2Mainte) IsTimeOver(d time.Duration) bool {
	return time.Now().Add(d).After(*self.NotBefore)
}
