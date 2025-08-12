// Package sqlite_setup sets up a sqlite database for the speedtestlogger app
package sqlite_setup

import (
	"database/sql"
)

// SetUpTables creates all needed tables in database for speedtestlogger app.
func SetUpTables(tx *sql.Tx) error {
	_, err := tx.Exec("create table if not exists entry (id INTEGER PRIMARY KEY AUTOINCREMENT, ts INTEGER, download_mbps REAL, upload_mbps REAL)")
	if err != nil {
		return err
	}
	_, err = tx.Exec("create index if not exists entry_ts_idx on entry (ts)")
	return err
}
