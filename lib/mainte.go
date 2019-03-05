package checkawsec2mainte

import "time"

type EC2Mainte struct {
	NotAfter    time.Time
	NotBefore   time.Time
	Description string
}

type EC2Maintes []EC2Mainte

func (self EC2Maintes) Len() int {
	return len(self)
}

func (self EC2Maintes) Less(i, j int) bool {
	return self[i].NotAfter.Unix() < self[j].NotAfter.Unix()
}

func (self EC2Maintes) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self EC2Mainte) IsTimeOver(d time.Duration) bool {
	return time.Now().Add(d).After(self.NotBefore)
}
