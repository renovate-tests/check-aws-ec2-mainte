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
func (evs EC2Events) Filter(states ...EC2EventState) EC2Events {
	events := EC2Events{}

EVENTS:
	for _, ev := range evs {
		for _, state := range states {
			if ev.State == state {
				continue EVENTS // Skip and Next Event
			}
		}
		events = append(events, ev)
	}
	return events
}

// UpdateStates ... Descriptionに含まれている文字列からStateを設定
func (evs EC2Events) UpdateStates() {
	for i, ev := range evs {
		switch {
		case strings.Contains(ev.Description, "[Completed]"):
			evs[i].State = StateCompleted
		case strings.Contains(ev.Description, "[Canceled]"):
			evs[i].State = StateCanceled
		default:
			evs[i].State = StateActive
		}
	}
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
