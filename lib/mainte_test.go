package checkawsec2mainte

import (
	"time"
	"testing"

)


func (self EC2Mainte) IsTimeOver(now time.Time, d time.Duration) bool {
	return now.Add(d).After(*self.NotBefore)
}

func TestTimeIsOver(t *testing.T) {
	mainte := EC2Mainte{}

	mainte.IsTimeOver(time.Now(), 1*time.Hour)
}
