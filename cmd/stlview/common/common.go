// Package common contains common routines for the stlview pages.
package common

import (
	"html/template"
	"net/url"
	"time"

	"github.com/keep94/speedtestlogger/stl/aggregators"
	"github.com/keep94/speedtestlogger/stl/format"
	"github.com/keep94/toolbox/date_util"
	"github.com/keep94/toolbox/http_util"
)

const (
	Date        = "date"
	DayPage     = "/day"
	SummaryPage = "/summary"
)

// NewTemplate returns a new template instance. name is the name
// of the template; templateStr is the template string.
func NewTemplate(name, templateStr string) *template.Template {
	return template.Must(template.New(name).Parse(templateStr))
}

// TimestampFormatter handles formatting timestamps
type TimestampFormatter struct {

	// The time zone.
	Location *time.Location
}

// FormatTimestamp formats a timestamp.
func (t *TimestampFormatter) FormatTimestamp(ts int64) string {
	return format.Time(ts, t.Location)
}

// SpeedFormatter formats internet speeds.
type SpeedFormatter struct {
}

// FormatSpeed formats an internet speed.
func (s *SpeedFormatter) FormatSpeed(mbps float64) string {
	return format.Speed(mbps)
}

// ParseDateParam parses the value of the date parameter. If the date
// parameter is of the form yyyy, then ParseDateParam returns the year
// with the year DateHandler. If the date parameter is of the form yyyyMM,
// then ParseDateParam returns the year and month with the month
// DateHandler. If the date parameter is of the form yyyyMMdd, then
// ParseDateParam returns the year, month, and day with the day DateHandler.
// If there is an error parsing the date, ParseDateParam returns the
// current time according to the now parameter normalized with the
// defaultHandler along with the defaultHandler.
func ParseDateParam(
	dateParam string,
	now time.Time,
	defaultHandler DateHandler) (time.Time, DateHandler) {
	var returnedHandler DateHandler
	if len(dateParam) == 4 {
		returnedHandler = Year()
		dateParam += "0101"
	} else if len(dateParam) == 6 {
		returnedHandler = Month()
		dateParam += "01"
	} else if len(dateParam) == 8 {
		returnedHandler = Day()
	}
	result, err := time.Parse(date_util.YMDFormat, dateParam)
	if err != nil || returnedHandler == nil {
		return defaultHandler.Normalize(now), defaultHandler
	}
	return result, returnedHandler
}

// DateHandler handles dates.
type DateHandler interface {

	// DrillDown returns the link for a date in the summary table.
	DrillDown(date time.Time) *url.URL

	// DrillDownFormat returns the date in the summary table properly formatted
	DrillDownFormat(date time.Time) string

	// DrillUp returns the link to the page one level up.
	DrillUp(current time.Time) *url.URL

	// Format formats the current date for the page
	Format(current time.Time) string

	// Prev returns a link to the previous page.
	Prev(current time.Time) *url.URL

	// Next returns a link to the next page.
	Next(current time.Time) *url.URL

	// Recurring says whether the summary table should be daily or monthly.
	Recurring() aggregators.Recurring

	// End returns the end date for the current page.
	End(current time.Time) time.Time

	// Normalize normalizes the current date
	Normalize(current time.Time) time.Time
}

func Day() DateHandler {
	return dayHandler{}
}

func Month() DateHandler {
	return monthHandler{}
}

func Year() DateHandler {
	return yearHandler{}
}

type dayHandler struct {
}

func (d dayHandler) DrillDown(date time.Time) *url.URL {
	return nil
}

func (d dayHandler) DrillDownFormat(date time.Time) string {
	return ""
}

func (d dayHandler) DrillUp(current time.Time) *url.URL {
	month := aggregators.Monthly().Normalize(current)
	return http_util.NewUrl(SummaryPage, Date, month.Format("200601"))
}

func (d dayHandler) Format(current time.Time) string {
	return current.Format("Mon 01/02/2006")
}

func (d dayHandler) Prev(current time.Time) *url.URL {
	prev := aggregators.Daily().Add(current, -1)
	return http_util.NewUrl(DayPage, Date, prev.Format("20060102"))
}

func (d dayHandler) Next(current time.Time) *url.URL {
	next := aggregators.Daily().Add(current, 1)
	return http_util.NewUrl(DayPage, Date, next.Format("20060102"))
}

func (d dayHandler) Recurring() aggregators.Recurring {
	return nil
}

func (d dayHandler) End(current time.Time) time.Time {
	return aggregators.Daily().Add(current, 1)
}

func (d dayHandler) Normalize(current time.Time) time.Time {
	return aggregators.Daily().Normalize(current)
}

type monthHandler struct {
}

func (m monthHandler) DrillDown(date time.Time) *url.URL {
	return http_util.NewUrl(DayPage, Date, date.Format("20060102"))
}

func (m monthHandler) DrillDownFormat(date time.Time) string {
	return date.Format("Mon 01/02/2006")
}

func (m monthHandler) DrillUp(current time.Time) *url.URL {
	year := aggregators.Yearly().Normalize(current)
	return http_util.NewUrl(SummaryPage, Date, year.Format("2006"))
}

func (m monthHandler) Format(current time.Time) string {
	return current.Format("01/2006")
}

func (m monthHandler) Prev(current time.Time) *url.URL {
	prev := aggregators.Monthly().Add(current, -1)
	return http_util.NewUrl(SummaryPage, Date, prev.Format("200601"))
}

func (m monthHandler) Next(current time.Time) *url.URL {
	next := aggregators.Monthly().Add(current, 1)
	return http_util.NewUrl(SummaryPage, Date, next.Format("200601"))
}

func (m monthHandler) Recurring() aggregators.Recurring {
	return aggregators.Daily()
}

func (m monthHandler) End(current time.Time) time.Time {
	return aggregators.Monthly().Add(current, 1)
}

func (m monthHandler) Normalize(current time.Time) time.Time {
	return aggregators.Monthly().Normalize(current)
}

type yearHandler struct {
}

func (y yearHandler) DrillDown(date time.Time) *url.URL {
	return http_util.NewUrl(SummaryPage, Date, date.Format("200601"))
}

func (y yearHandler) DrillDownFormat(date time.Time) string {
	return date.Format("01/2006")
}

func (y yearHandler) DrillUp(current time.Time) *url.URL {
	return nil
}

func (y yearHandler) Format(current time.Time) string {
	return current.Format("2006")
}

func (y yearHandler) Prev(current time.Time) *url.URL {
	prev := aggregators.Yearly().Add(current, -1)
	return http_util.NewUrl(SummaryPage, Date, prev.Format("2006"))
}

func (y yearHandler) Next(current time.Time) *url.URL {
	next := aggregators.Yearly().Add(current, 1)
	return http_util.NewUrl(SummaryPage, Date, next.Format("2006"))
}

func (y yearHandler) Recurring() aggregators.Recurring {
	return aggregators.Monthly()
}

func (y yearHandler) End(current time.Time) time.Time {
	return aggregators.Yearly().Add(current, 1)
}

func (y yearHandler) Normalize(current time.Time) time.Time {
	return aggregators.Yearly().Normalize(current)
}
