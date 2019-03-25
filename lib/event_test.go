package checkawsec2mainte_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	checkawsec2mainte "github.com/ntrv/check-aws-ec2-mainte/lib"
)

func TestTimeIsOver(t *testing.T) {
	past, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")
	future, _ := time.Parse(time.RFC3339, "2018-03-15T19:14:05+09:00")

	mainte := checkawsec2mainte.Event{
		NotBefore: past,
	}

	assert.True(t, mainte.IsTimeOver(future, 1*time.Hour))
}

func TestTimeIsNotOver(t *testing.T) {
	past, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05+09:00")
	future, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05+07:00")

	mainte := checkawsec2mainte.Event{
		NotBefore: future,
	}

	assert.False(t, mainte.IsTimeOver(past, 1*time.Hour))
}
