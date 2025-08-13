package dates

import (
	"testing"
	"time"

	"github.com/keep94/toolbox/date_util"
	"github.com/stretchr/testify/assert"
)

func TestDatePart(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	assert.NoError(t, err)
	date := DatePart(1755057599, loc)
	assert.Equal(t, date_util.YMD(2025, 8, 12), date)
	date = DatePart(1755057600, loc)
	assert.Equal(t, date_util.YMD(2025, 8, 13), date)
}

func TestToTimestamp(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	assert.NoError(t, err)
	ts := ToTimestamp(date_util.YMD(2025, 8, 12), loc)
	assert.Equal(t, int64(1754971200), ts)
}
