package checkawsec2mainte

import (
	"sort"
	"strings"
	"time"
)

type IEC2Events interface {
	GetCloseEvent() EC2Mainte
	Len() int
}

type EC2Events []EC2Event

// Filter ... Filter EC2Events containing substr in Description
func (evs EC2Events) Filter(substr string) EC2Events {
	events := EC2Events{}

	for _, ev := range evs {
		if strings.Contains(ev.Description, substr) {
			continue
		}
		events = append(events, ev)
	}

	return events
}

// GetCloseEvent ... Get oldest EC2Event
func (evs EC2Events) GetCloseEvent() EC2Event {
	// Copy to temporary variable
	// Prevent to mutate self
	events := make(EC2Events, cap(evs))
	copy(events, evs)

	// Sort as NotBefore date
	sort.Stable(events)

	return events[0]
}

// BeforeAll ...
func (evs EC2Events) BeforeAll(d time.Time) bool {
	for _, ev := range evs {
		if ev.NotBefore.After(d) {
			return false
		}
	}
	return true
}

// Len ... Implement sort.Interface
func (evs EC2Events) Len() int {
	return len(evs)
}

// Less ... Implement sort.Interface
func (evs EC2Events) Less(i, j int) bool {
	return evs[i].NotBefore.Unix() < evs[j].NotBefore.Unix()
}

// Swap ... Implement sort.Interface
func (evs EC2Events) Swap(i, j int) {
	evs[i], evs[j] = evs[j], evs[i]
}
