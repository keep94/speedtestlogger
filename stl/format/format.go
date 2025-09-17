// Package format provides formatting routines for speedtestlogger app.
package format

import (
	"strconv"
	"time"
)

// Float formats a float64 with precision digits after the decimal
func Float(value float64, precision int) string {
	return strconv.FormatFloat(value, 'f', precision, 64)
}

// Time formats a time given seconds after Jan 1, 1970 GMT
// and the time zone.
func Time(ts int64, loc *time.Location) string {
	timestamp := time.Unix(ts, 0).In(loc)
	return timestamp.Format("Mon 01/02/2006 15:04")
}
