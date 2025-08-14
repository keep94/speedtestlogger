package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/keep94/context"
	"github.com/keep94/speedtestlogger/cmd/stlview/day"
	"github.com/keep94/speedtestlogger/stl/stldb/for_sqlite"
	"github.com/keep94/toolbox/build"
	"github.com/keep94/toolbox/date_util"
	"github.com/keep94/toolbox/db"
	"github.com/keep94/toolbox/db/sqlite3_db"
	"github.com/keep94/toolbox/http_util"
	"github.com/keep94/toolbox/logging"
	"github.com/keep94/weblogs"
	_ "github.com/mattn/go-sqlite3"
)

var (
	kClock date_util.Clock = date_util.SystemClock{}
)

var (
	fDb   string
	fPort string
)

var (
	kDoer  db.Doer
	kStore *for_sqlite.Store
)

func main() {
	flag.Parse()
	if fDb == "" {
		fmt.Println("Need to specify at least -db flag.")
		flag.Usage()
		os.Exit(1)
	}
	setupDb(fDb)
	http.HandleFunc("/", rootRedirect)
	version, _ := build.MainVersion()
	http.Handle(
		"/day",
		&day.Handler{
			Store:    kStore,
			BuildId:  build.BuildId(version),
			Clock:    kClock,
			Location: time.Local})
	defaultHandler := context.ClearHandler(
		weblogs.HandlerWithOptions(
			http.DefaultServeMux,
			&weblogs.Options{Logger: logging.ApacheCommonLoggerWithLatency()}))
	if err := http.ListenAndServe(fPort, defaultHandler); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func rootRedirect(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http_util.Redirect(w, r, "/day")
	} else {
		http_util.Error(w, http.StatusNotFound)
	}
}

func setupDb(filepath string) {
	rawdb, err := sql.Open("sqlite3", filepath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbase := sqlite3_db.New(rawdb)
	kDoer = sqlite3_db.NewDoer(dbase)
	kStore = for_sqlite.New(dbase)
}

func init() {
	flag.StringVar(&fPort, "http", ":8080", "Port to bind")
	flag.StringVar(&fDb, "db", "", "Path to database file")
}
