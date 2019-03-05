package checkawsec2mainte

import "time"

type EC2Mainte struct {
	NotAfter    time.Time
	NotBefore   time.Time
	Description string
}

func (self EC2Mainte) IsTimeOver(d time.Duration) bool {
	return time.Now().Add(d).After(self.NotBefore)
}
