// Package stl provides the data structures for the speedtestlogger app
package stl

import (
	"strconv"
	"time"
)

// Entry represents a speed test data point
type Entry struct {

	// Id of Entry
	Id int64

	// Seconds since Jan 1 1970 GMT
	Ts int64

	// Download speed in megabits per second
	DownloadMbps float64

	// Upload speed in megabits per second
	UploadMbps float64
}

// DownloadMbpsString returns the download speed as "176.00" for 176 Mbps.
func (e *Entry) DownloadMbpsString() string {
	return strconv.FormatFloat(e.DownloadMbps, 'f', 2, 64)
}

// UploadMbpsString returns the upload speed as "17.00" for 17 Mbps.
func (e *Entry) UploadMbpsString() string {
	return strconv.FormatFloat(e.UploadMbps, 'f', 2, 64)
}

// FormatTime returns formatted time given seconds after Jan 1, 1970 GMT
// and time zone.
func FormatTime(Ts int64, loc *time.Location) string {
	timestamp := time.Unix(Ts, 0).In(loc)
	return timestamp.Format("01/02/2006 15:04")
}
