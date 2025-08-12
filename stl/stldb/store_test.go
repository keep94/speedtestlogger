package stldb

import (
	"testing"
	"time"

	"github.com/keep94/consume2"
	"github.com/keep94/speedtestlogger/stl"
	"github.com/keep94/toolbox/db"
	"github.com/stretchr/testify/assert"
)

func TestEntriesForDay(t *testing.T) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	assert.NoError(t, err)
	var store fakeEntriesRunner
	var result []stl.Entry
	EntriesForDay(
		&store,
		nil,
		2025,
		8,
		12,
		loc,
		consume2.AppendTo(&result))
	assert.Equal(t, int64(1754982000), result[1].Ts)
	assert.Equal(t, int64(1755068400), result[0].Ts)
}

type fakeEntriesRunner struct {
}

func (f fakeEntriesRunner) Entries(
	t db.Transaction,
	start,
	end int64,
	consumer consume2.Consumer[stl.Entry]) error {
	result := []stl.Entry{{Id: 2, Ts: end}, {Id: 1, Ts: start}}
	consume2.FromSlice(result, consumer)
	return nil
}
