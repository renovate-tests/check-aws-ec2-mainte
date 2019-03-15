package checkawsec2mainte

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeIsOver(t *testing.T) {
	mainte := EC2Mainte{
		NotBefore: time.Now(),
	}

	assert.True(t, mainte.IsTimeOver(time.Now(), 1*time.Hour))
}
