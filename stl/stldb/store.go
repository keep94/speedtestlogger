// Package stldb handles storing instances in the speedtestlogger app to
// persistent storage.
package stldb

import (
	"time"

	"github.com/keep94/consume2"
	"github.com/keep94/speedtestlogger/stl"
	"github.com/keep94/toolbox/db"
)

type AddEntryRunner interface {

	// AddEntry adds a new entry to persistent storage.
	AddEntry(t db.Transaction, entry *stl.Entry) error
}

type EntriesRunner interface {

	// Entries returns all entries (most recent to least recent) within a
	// given time range. startTime and endTime are seconds since Jan 1, 1970.
	Entries(
		t db.Transaction,
		startTime,
		endTime int64,
		consumer consume2.Consumer[stl.Entry]) error
}

type RemoveEntriesRunner interface {

	// RemoveEntries removes entries within a given time range.
	// startTime and endTime are seconds since Jan 1, 1970.
	RemoveEntries(t db.Transaction, startTime, endTime int64) error
}

// EntriesForDay returns the entries (most recent to least recent) for a
// particular day in the given time zone.
func EntriesForDay(
	store EntriesRunner,
	t db.Transaction,
	year,
	month,
	day int,
	loc *time.Location,
	consumer consume2.Consumer[stl.Entry]) error {
	start := time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 0, 1)
	return store.Entries(t, start.Unix(), end.Unix(), consumer)
}
