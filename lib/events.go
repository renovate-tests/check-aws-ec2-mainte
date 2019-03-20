package checkawsec2mainte

import (
	"time"

	"sort"
	"strings"
)

type IEC2Events interface {
	GetCloseEvent() EC2Mainte
	Len() int
}

type EC2Events []EC2Event

// Filter ... Filter EC2Events containing substr in Description
func (e EC2Events) Filter(substr string) EC2Events {
	events := EC2Events{}

	for _, event := range e {
		if strings.Contains(event.Description, substr) {
			continue
		}
		events = append(events, event)
	}

	return events
}

// GetCloseEvent ... Get oldest EC2Event
func (self EC2Events) GetCloseEvent() EC2Event {
	// Copy to temporary variable
	// Prevent to mutate self
	events := make(EC2Events, cap(self))
	copy(events, self)

	// Sort as NotBefore date
	sort.Stable(events)

	return events[0]
}

// BeforeAll ...
func (self EC2Events) BeforeAll(d time.Time) bool {
	for _, a := range self {
		if a.NotBefore.After(d) {
			return false
		}
	}
	return true
}

// Len ... Implement sort.Interface
func (self EC2Events) Len() int {
	return len(self)
}

// Less ... Implement sort.Interface
func (self EC2Events) Less(i, j int) bool {
	return self[i].NotBefore.Unix() < self[j].NotBefore.Unix()
}

// Swap ... Implement sort.Interface
func (self EC2Events) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}
