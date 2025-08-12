// Package fixture provides test suites to test implementations of the
// datastore interfaces in the stldb package.
package fixture

import (
	"testing"

	"github.com/keep94/consume2"
	"github.com/keep94/speedtestlogger/stl"
	"github.com/keep94/speedtestlogger/stl/stldb"
	"github.com/stretchr/testify/assert"
)

var (
	kFirstEntry = stl.Entry{
		Ts:           123,
		DownloadMbps: 50.0,
		UploadMbps:   5.0,
	}
	kSecondEntry = stl.Entry{
		Ts:           234,
		DownloadMbps: 60.0,
		UploadMbps:   6.0,
	}
	kThirdEntry = stl.Entry{
		Ts:           345,
		DownloadMbps: 70.0,
		UploadMbps:   7.0,
	}
)

type Store interface {
	stldb.AddEntryRunner
	stldb.EntriesRunner
	stldb.RemoveEntriesRunner
}

func Entries(t *testing.T, store Store) {
	first := kFirstEntry
	assert.NoError(t, store.AddEntry(nil, &first))
	second := kSecondEntry
	assert.NoError(t, store.AddEntry(nil, &second))
	third := kThirdEntry
	assert.NoError(t, store.AddEntry(nil, &third))
	assert.Equal(t, int64(1), first.Id)
	assert.Equal(t, int64(2), second.Id)
	assert.Equal(t, int64(3), third.Id)

	var entries []stl.Entry
	assert.NoError(
		t,
		store.Entries(nil, 100, 400, consume2.AppendTo(&entries)))
	assert.Equal(t, []stl.Entry{third, second, first}, entries)

	entries = nil
	assert.NoError(
		t,
		store.Entries(nil, 100, 300, consume2.AppendTo(&entries)))
	assert.Equal(t, []stl.Entry{second, first}, entries)

	assert.NoError(t, store.RemoveEntries(nil, 100, 300))

	entries = nil
	assert.NoError(
		t,
		store.Entries(nil, 100, 400, consume2.AppendTo(&entries)))
	assert.Equal(t, []stl.Entry{third}, entries)
}
