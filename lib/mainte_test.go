package checkawsec2mainte

import (
	"time"
	"testing"

)

func TestTimeIsOver(t *testing.T) {
	mainte := EC2Mainte{}

	mainte.IsTimeOver(time.Now(), 1*time.Hour)
}
