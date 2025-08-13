package aggregators

import (
	"testing"
	"time"

	"github.com/keep94/speedtestlogger/stl"
	"github.com/keep94/speedtestlogger/stl/dates"
	"github.com/keep94/toolbox/date_util"
	"github.com/stretchr/testify/assert"
)

func TestByPeriodTotaler(t *testing.T) {
	loc, err := time.LoadLocation("America/New_York")
	assert.NoError(t, err)
	totaler := NewByPeriodTotaler(
		date_util.YMD(2025, 9, 12),
		date_util.YMD(2026, 1, 12),
		Monthly(),
		loc)
	totaler.Add(stl.Entry{
		Ts:           dates.ToTimestamp(date_util.YMD(2026, 1, 1), loc),
		DownloadMbps: 57.25,
		UploadMbps:   31.25,
	})
	totaler.Add(stl.Entry{
		Ts:           dates.ToTimestamp(date_util.YMD(2025, 12, 13), loc),
		DownloadMbps: 130.0,
		UploadMbps:   13.0,
	})
	totaler.Add(stl.Entry{
		Ts:           dates.ToTimestamp(date_util.YMD(2025, 12, 1), loc),
		DownloadMbps: 110.0,
		UploadMbps:   11.0,
	})
	totaler.Add(stl.Entry{
		Ts:           dates.ToTimestamp(date_util.YMD(2025, 10, 1), loc),
		DownloadMbps: 100.0,
		UploadMbps:   10.0,
	})
	datedSummaries := totaler.DatedSummaries()
	assert.Len(t, datedSummaries, 4)

	assert.Equal(t, date_util.YMD(2025, 12, 1), datedSummaries[0].Date)
	assert.Equal(t, 120.0, datedSummaries[0].DownloadMbps.Avg())
	assert.Equal(t, 12.0, datedSummaries[0].UploadMbps.Avg())

	assert.Equal(t, date_util.YMD(2025, 11, 1), datedSummaries[1].Date)
	assert.True(t, datedSummaries[1].DownloadMbps.Empty())
	assert.True(t, datedSummaries[1].UploadMbps.Empty())

	assert.Equal(t, date_util.YMD(2025, 10, 1), datedSummaries[2].Date)
	assert.Equal(t, 100.0, datedSummaries[2].DownloadMbps.Avg())
	assert.Equal(t, 10.0, datedSummaries[2].UploadMbps.Avg())

	assert.Equal(t, date_util.YMD(2025, 9, 1), datedSummaries[3].Date)
	assert.True(t, datedSummaries[3].DownloadMbps.Empty())
	assert.True(t, datedSummaries[3].UploadMbps.Empty())
}

func TestDaily(t *testing.T) {
	r := Daily()
	assert.Equal(
		t, date_util.YMD(2025, 8, 12), r.Normalize(date_util.YMD(2025, 8, 12)))
	assert.Equal(
		t, date_util.YMD(2025, 8, 17), r.Add(date_util.YMD(2025, 8, 12), 5))
}

func TestYearly(t *testing.T) {
	r := Yearly()
	assert.Equal(
		t, date_util.YMD(2025, 1, 1), r.Normalize(date_util.YMD(2025, 8, 12)))
	assert.Equal(
		t, date_util.YMD(2030, 1, 1), r.Add(date_util.YMD(2025, 1, 1), 5))
}
