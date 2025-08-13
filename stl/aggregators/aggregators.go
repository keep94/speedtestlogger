// Package aggregators aggregates values.
package aggregators

import (
	"slices"
	"time"

	"github.com/keep94/speedtestlogger/stl"
	"github.com/keep94/speedtestlogger/stl/dates"
	"github.com/keep94/toolbox/date_util"
)

// Average represents an average
type Average struct {
	N   int
	Sum float64
}

// Add adds a new value to this average.
func (a *Average) Add(value float64) {
	a.N++
	a.Sum += value
}

// Empty returns true if this average is empty.
func (a *Average) Empty() bool {
	return a.N <= 0
}

// Avg returns the average value. Avg panics if called on an empty Average.
func (a *Average) Avg() float64 {
	if a.Empty() {
		panic("Avg() called on empty Average")
	}
	return a.Sum / float64(a.N)
}

// Summary represents a summary of internet speeds
type Summary struct {
	DownloadMbps Average
	UploadMbps   Average
}

// Add adds an stl.Entry to this summary.
func (s *Summary) Add(entry stl.Entry) {
	s.DownloadMbps.Add(entry.DownloadMbps)
	s.UploadMbps.Add(entry.UploadMbps)
}

// DatedSummary represents a dated summary.
type DatedSummary struct {
	Date time.Time
	Summary
}

// Recurring is the interface for recurring time periods. e.g monthly, yearly
type Recurring interface {

	// Normalize returns the beginning of a time period for a given date.
	Normalize(date time.Time) time.Time

	// Add returns the result of adding numPeriods time periods to a start
	// date.
	Add(date time.Time, numPeriods int) time.Time
}

func Daily() Recurring {
	return daily{}
}

func Monthly() Recurring {
	return monthly{}
}

func Yearly() Recurring {
	return yearly{}
}

type daily struct{}

func (d daily) Normalize(date time.Time) time.Time {
	return date_util.YMD(date.Year(), int(date.Month()), date.Day())
}

func (d daily) Add(date time.Time, numPeriods int) time.Time {
	return date.AddDate(0, 0, numPeriods)
}

type monthly struct{}

func (m monthly) Normalize(date time.Time) time.Time {
	return date_util.YMD(date.Year(), int(date.Month()), 1)
}

func (m monthly) Add(date time.Time, numPeriods int) time.Time {
	return date.AddDate(0, numPeriods, 0)
}

type yearly struct{}

func (y yearly) Normalize(date time.Time) time.Time {
	return date_util.YMD(date.Year(), 1, 1)
}

func (y yearly) Add(date time.Time, numPeriods int) time.Time {
	return date.AddDate(numPeriods, 0, 0)
}

// ByPeriodTotaler aggregates stl.Entry instances by period.
type ByPeriodTotaler struct {
	summaries []DatedSummary
	smap      map[time.Time]*DatedSummary
	recurring Recurring
	loc       *time.Location
}

// NewByPeriodTotaler creates a new ByPeriodTotaler that summarizes
// stl.Entry objects with timestamps between start inclusive and end
// exclusive. The recurring perameter indicates the recurring period
// such as daily, monthly or yearly. loc is the time zone used to convert
// timestamps to dates. This function converts start and end so that they
// fall on the beginning of the given recurring period.
func NewByPeriodTotaler(
	start,
	end time.Time,
	recurring Recurring,
	loc *time.Location) *ByPeriodTotaler {
	start = recurring.Normalize(start)
	end = recurring.Normalize(end)
	summaries := initializeDatedSummaries(start, end, recurring)
	smap := initializeSummaryMap(summaries)
	return &ByPeriodTotaler{
		summaries: summaries,
		smap:      smap,
		recurring: recurring,
		loc:       loc}
}

// Add adds a new entry to this instance.
func (b *ByPeriodTotaler) Add(entry stl.Entry) {
	cdate := dates.DatePart(entry.Ts, b.loc)
	datedSummaryPtr := b.smap[b.recurring.Normalize(cdate)]
	if datedSummaryPtr != nil {
		datedSummaryPtr.Add(entry)
	}
}

// DatedSummaries returns a copy of the summaries collected so far.
// Each DatedSummary falls on the beginning of a day, month, or year
// depending on the recurring parameter passed to NewByPeriodTotaler().
func (b *ByPeriodTotaler) DatedSummaries() []DatedSummary {
	return slices.Clone(b.summaries)
}

func initializeDatedSummaries(
	start, end time.Time, recurring Recurring) []DatedSummary {
	current := recurring.Add(end, -1)
	var result []DatedSummary
	for !current.Before(start) {
		result = append(result, DatedSummary{Date: current})
		current = recurring.Add(current, -1)
	}
	return result
}

func initializeSummaryMap(
	summaries []DatedSummary) map[time.Time]*DatedSummary {
	result := make(map[time.Time]*DatedSummary)
	for i := range summaries {
		result[summaries[i].Date] = &summaries[i]
	}
	return result
}
