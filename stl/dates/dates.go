// Package dates contains routines for converting between dates and
// timestamps.
package dates

import (
	"time"

	"github.com/keep94/toolbox/date_util"
)

// DatePart returns the date part of ts (seconds since Jan 1, 1970 GMT)
// using the given time zone.
func DatePart(ts int64, loc *time.Location) time.Time {
	timestamp := time.Unix(ts, 0).In(loc)
	return date_util.YMD(
		timestamp.Year(), int(timestamp.Month()), timestamp.Day())
}

// ToTimestamp converts midnight of the given date to seconds since
// Jan 1, 1970 GMT in the given time zone.
func ToTimestamp(date time.Time, loc *time.Location) int64 {
	timestamp := time.Date(
		date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
	return timestamp.Unix()
}
