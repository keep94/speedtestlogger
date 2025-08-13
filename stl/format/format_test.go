package format

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTime(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	assert.NoError(t, err)
	formatted := Time(1754993618, loc)
	assert.Equal(t, "Tue 08/12/2025 06:13", formatted)
}

func TestSpeed(t *testing.T) {
	assert.Equal(t, "51.38", Speed(51.375))
	assert.Equal(t, "0.03", Speed(0.03125))
}
