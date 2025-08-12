// Package for_sqlite provides a sqlite implementation of interfaces in the
// stldb package.
package for_sqlite

import (
	"database/sql"

	"github.com/keep94/consume2"
	"github.com/keep94/speedtestlogger/stl"
	"github.com/keep94/toolbox/db"
	"github.com/keep94/toolbox/db/sqlite3_db"
	"github.com/keep94/toolbox/db/sqlite3_rw"
)

const (
	kSQLEntries       = "select id, ts, download_mbps, upload_mbps from entry where ts >= ? and ts < ? order by ts desc"
	kSQLAddEntry      = "insert into entry (ts, download_mbps, upload_mbps) values (?, ?, ?)"
	kSQLRemoveEntries = "delete from entry where ts >= ? and ts < ?"
)

type Store struct {
	db sqlite3_db.Doer
}

// New creates a sqlite implementation of the speedtestlogger app datastore.
func New(db *sqlite3_db.Db) *Store {
	return &Store{db}
}

func (s *Store) AddEntry(t db.Transaction, entry *stl.Entry) error {
	return sqlite3_db.ToDoer(s.db, t).Do(func(tx *sql.Tx) error {
		return sqlite3_rw.AddRow(
			tx, (&rawEntry{}).init(entry), &entry.Id, kSQLAddEntry)
	})
}

func (s *Store) Entries(
	t db.Transaction,
	startTime,
	endTime int64,
	consumer consume2.Consumer[stl.Entry]) error {
	return sqlite3_db.ToDoer(s.db, t).Do(func(tx *sql.Tx) error {
		return sqlite3_rw.ReadMultiple[stl.Entry](
			tx,
			(&rawEntry{}).init(&stl.Entry{}),
			consumer,
			kSQLEntries,
			startTime,
			endTime)
	})
}

func (s *Store) RemoveEntries(
	t db.Transaction, startTime, endTime int64) error {
	return sqlite3_db.ToDoer(s.db, t).Do(func(tx *sql.Tx) error {
		_, err := tx.Exec(kSQLRemoveEntries, startTime, endTime)
		return err
	})
}

type rawEntry struct {
	*stl.Entry
	sqlite3_rw.SimpleRow
}

func (r *rawEntry) init(bo *stl.Entry) *rawEntry {
	r.Entry = bo
	return r
}

func (r *rawEntry) Ptrs() []interface{} {
	return []interface{}{&r.Id, &r.Ts, &r.DownloadMbps, &r.UploadMbps}
}

func (r *rawEntry) Values() []interface{} {
	return []interface{}{r.Ts, r.DownloadMbps, r.UploadMbps, r.Id}
}

func (r *rawEntry) ValueRead() stl.Entry {
	return *r.Entry
}
