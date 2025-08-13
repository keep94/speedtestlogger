// Package format provides formatting routines for speedtestlogger app.
package format

import (
	"strconv"
	"time"
)

// Speed formats internet speeds.
func Speed(mbps float64) string {
	return strconv.FormatFloat(mbps, 'f', 2, 64)
}

// Time formats a time given seconds after Jan 1, 1970 GMT
// and the time zone.
func Time(ts int64, loc *time.Location) string {
	timestamp := time.Unix(ts, 0).In(loc)
	return timestamp.Format("Mon 01/02/2006 15:04")
}
