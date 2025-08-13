// Package stl provides the data structures for the speedtestlogger app
package stl

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
