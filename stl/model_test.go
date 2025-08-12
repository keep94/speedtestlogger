package stl

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatTime(t *testing.T) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	assert.NoError(t, err)
	formatted := FormatTime(1754993618, loc)
	assert.Equal(t, "08/12/2025 03:13", formatted)
}

func TestFormatSpeeds(t *testing.T) {
	entry := Entry{DownloadMbps: 51.375, UploadMbps: 0.03125}
	assert.Equal(t, "51.38", entry.DownloadMbpsString())
	assert.Equal(t, "0.03", entry.UploadMbpsString())
}
