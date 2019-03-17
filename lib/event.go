package checkawsec2mainte

import (
	"bytes"
	"log"
	"text/template"
	"time"
	"strings"

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
	const tplText = "Code: {{.Code}}, InstanceId: {{.InstanceId}}, Date: {{.NotBefore}} - {{.NotAfter}}, Description: {{.Description}}"
	var buf bytes.Buffer

	tpl, err := template.New("").Parse(strings.Trim(tplText, "\t"))
	if err != nil {
		log.Fatal(err)
		return ""
	}
	if err := tpl.Execute(&buf, self); err != nil {
		log.Fatal(err)
		return ""
	}
	return buf.String()
}
