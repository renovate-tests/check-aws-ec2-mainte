package checkawsec2mainte

import (
	"strings"
)

type IEC2Events interface {
	GetCloseEvent() EC2Mainte
	Len() int
}

type EC2Events []EC2Event

func (e EC2Events) Filter(substr string) EC2Events {
	events := EC2Maintes{}

	for _, event := range e {
		if strings.Contains(event.Description, substr) {
			continue
		}
		events = append(events, event)
	}

	return events
}

func (self EC2Events) GetCloseEvent() EC2Events {
	return self[len(self)-1]
}

func (self EC2Events) Len() int {
	return len(self)
}

func (self EC2Events) Less(i, j int) bool {
	return self[i].NotBefore.Unix() < self[j].NotBefore.Unix()
}

func (self EC2Events) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}
