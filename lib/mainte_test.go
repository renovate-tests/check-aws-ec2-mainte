package checkawsec2mainte

import (
	"testing"
	"time"
)

func TestTimeIsOver(t *testing.T) {
	mainte := EC2Mainte{
		NotBefore: time.Now(),
	}

	mainte.IsTimeOver(time.Now(), 1*time.Hour)
}
