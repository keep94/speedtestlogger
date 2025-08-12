package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/keep94/speedtestlogger/stl/stldb/sqlite_setup"
	"github.com/keep94/toolbox/db/sqlite3_db"
	_ "github.com/mattn/go-sqlite3"
)

var (
	fDb string
)

func main() {
	flag.Parse()
	if fDb == "" {
		fmt.Println("Need to specify at least -db flag.")
		flag.Usage()
		os.Exit(2)
	}
	db := openDb(fDb)
	defer db.Close()
	initDb(db)
}

func openDb(dbPath string) *sqlite3_db.Db {
	rawdb, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Printf("Unable to open database - %s\n", dbPath)
		os.Exit(1)
	}
	return sqlite3_db.New(rawdb)
}

func initDb(dbase *sqlite3_db.Db) {
	err := dbase.Do(sqlite_setup.SetUpTables)
	if err != nil {
		fmt.Printf("Unable to create tables - %v\n", err)
		os.Exit(1)
	}
}

func init() {
	flag.StringVar(&fDb, "db", "", "Path to database file")
}
