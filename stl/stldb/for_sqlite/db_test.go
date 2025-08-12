package for_sqlite_test

import (
	"database/sql"
	"testing"

	"github.com/keep94/speedtestlogger/stl/stldb/fixture"
	"github.com/keep94/speedtestlogger/stl/stldb/for_sqlite"
	"github.com/keep94/speedtestlogger/stl/stldb/sqlite_setup"
	"github.com/keep94/toolbox/db/sqlite3_db"
	_ "github.com/mattn/go-sqlite3"
)

func TestEntries(t *testing.T) {
	db := openDb(t)
	defer closeDb(t, db)
	fixture.Entries(t, for_sqlite.New(db))
}

func openDb(t *testing.T) *sqlite3_db.Db {
	rawdb, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	db := sqlite3_db.New(rawdb)
	err = db.Do(sqlite_setup.SetUpTables)
	if err != nil {
		t.Fatalf("Error creating tables: %v", err)
	}
	return db
}

func closeDb(t *testing.T, db *sqlite3_db.Db) {
	if err := db.Close(); err != nil {
		t.Errorf("Error closing database: %v", err)
	}
}
