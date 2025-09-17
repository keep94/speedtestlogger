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

func TestFloat(t *testing.T) {
	assert.Equal(t, "51.38", Float(51.375, 2))
	assert.Equal(t, "0.0312", Float(0.03125, 4))
}
