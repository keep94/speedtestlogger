package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/keep94/speedtestlogger/stl"
	"github.com/keep94/speedtestlogger/stl/stldb"
	"github.com/keep94/speedtestlogger/stl/stldb/for_sqlite"
	"github.com/keep94/toolbox/db/sqlite3_db"
	_ "github.com/mattn/go-sqlite3"
)

const (
	kMillionFloat = 1000000.0
	kEightFloat   = 8.0
)

var (
	fDb  string
	fCsv string
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
	store := for_sqlite.New(db)
	entry := stl.Entry{Ts: time.Now().Unix()}
	var csvrow []string
	if fCsv != "" {
		csvrow = readcsv(fCsv)
	}
	if len(csvrow) < 7 {
		log.Println("Not enough columns in csv:", csvrow)
	} else {
		download, _ := strconv.ParseFloat(csvrow[5], 64)
		upload, _ := strconv.ParseFloat(csvrow[6], 64)
		entry.DownloadMbps = download * kEightFloat / kMillionFloat
		entry.UploadMbps = upload * kEightFloat / kMillionFloat
	}
	addEntry(store, &entry)
}

func readcsv(csvPath string) []string {
	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatal("Unable to open csv file: ", csvPath)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	result, err := reader.Read()
	if err != nil {
		log.Fatal("Unable to read csv file: ", err)
	}
	return result
}

func addEntry(store stldb.AddEntryRunner, entry *stl.Entry) {
	if err := store.AddEntry(nil, entry); err != nil {
		log.Fatal("Error writing to db: ", err)
	}
}

func openDb(dbPath string) *sqlite3_db.Db {
	rawdb, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Unable to open database: ", dbPath)
	}
	return sqlite3_db.New(rawdb)
}

func init() {
	flag.StringVar(&fDb, "db", "", "Path to database file")
	flag.StringVar(&fCsv, "csv", "", "path to csv file")
}
