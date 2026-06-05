package daterange

import (
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/keep94/consume2"
	"github.com/keep94/speedtestlogger/cmd/stlview/common"
	"github.com/keep94/speedtestlogger/stl/aggregators"
	"github.com/keep94/speedtestlogger/stl/dates"
	"github.com/keep94/speedtestlogger/stl/stldb"
	"github.com/keep94/toolbox/date_util"
	"github.com/keep94/toolbox/http_util"
)

var (
	kTemplateSpec = `
<html>
<head>
  <title>Internet Speeds</title>
  <style>
  h1 {
    font-size: 40px;
  }
  th {
    font-size: 30px;
  }
  .my-submit-btn {
    font-size: 30px;
  }
  td, .normal {
    font-size: 30px;
  }
  td.today {
    font-style: italic;
  }
  .error {
    font-size: 30px;
    color:#FF0000;
  }
  input[type="text"], textarea {
    font-family: inherit;
    font-size: inherit;
    font-weight: inherit;
  }
  </style>
</head>
<body>
  <h1>Average Speeds Over Date Range &nbsp; &nbsp; Build: {{.BuildId}}</h1>
  <a href="{{.BackLink}}">Back</a>
  <br><br>
  {{with .ErrorMessage}}
    <span class="error">{{.}}</span>
  {{end}}
  <form>
    <input type="hidden" name="prev" value="{{.Get "prev"}}">
    <table>
      <tr>
        <td>Start Date (yyyyMMdd): </td>
        <td><input type="text" name="sd" value="{{.Get "sd"}}"></td>
      </tr>
      <tr>
        <td>End Date (yyyyMMdd): </td>
        <td><input type="text" name="ed" value="{{.Get "ed"}}"></td>
      </tr>
    </table>
  <input type="submit" class="my-submit-btn" value="Get">
  </form>
  <hr>
  {{with $top := .}}
  {{with .Summary}}
  <span class="normal">
  Download Average (Mbps): {{with .DownloadMbps}}{{if .Exists}}{{$top.FormatSpeed .Avg}}{{else}}--{{end}}{{end}}
  <br>
  Upload Average (Mbps): {{with .UploadMbps}}{{if .Exists}}{{$top.FormatSpeed .Avg}}{{else}}--{{end}}{{end}}
  <br>
  Percent Uptime: {{with .PercentUptime}}{{if .Exists}}{{$top.FormatPercent .Avg}}{{else}}--{{end}}{{end}}
  </span>
  {{end}}
  {{end}}
</body>
</html>`
)

var (
	kTemplate *template.Template
)

type Handler struct {
	Store    stldb.EntriesRunner
	BuildId  string
	Location *time.Location
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	startDateStr := strings.TrimSpace(r.Form.Get("sd"))
	endDateStr := strings.TrimSpace(r.Form.Get("ed"))
	prevUrl := http_util.Sanitize(r.Form.Get("prev"), common.DayPage)
	if startDateStr == "" && endDateStr == "" {
		http_util.WriteTemplate(
			w,
			kTemplate,
			&view{
				Values:   r.Form,
				BuildId:  h.BuildId,
				BackLink: prevUrl,
			},
		)
		return
	}
	startDate, endDate, err := parseDates(startDateStr, endDateStr)
	if err != nil {
		http_util.WriteTemplate(
			w,
			kTemplate,
			&view{
				Values:       r.Form,
				BuildId:      h.BuildId,
				BackLink:     prevUrl,
				ErrorMessage: "Start and End Date must be in yyyyMMdd format",
			},
		)
		return
	}
	var summary aggregators.Summary
	err = h.Store.Entries(
		nil,
		dates.ToTimestamp(startDate, h.Location),
		dates.ToTimestamp(endDate, h.Location),
		consume2.Call(summary.Add))
	if err != nil {
		http_util.ReportError(w, "Error reading database", err)
		return
	}
	http_util.WriteTemplate(
		w,
		kTemplate,
		&view{
			Values:   r.Form,
			BuildId:  h.BuildId,
			BackLink: prevUrl,
			Summary:  &summary,
		},
	)
}

type view struct {
	common.SpeedFormatter
	common.PercentFormatter
	url.Values
	BuildId      string
	BackLink     string
	Summary      *aggregators.Summary
	ErrorMessage string
}

func parseDates(startDateStr, endDateStr string) (
	startDate time.Time, endDate time.Time, err error) {
	startDate, err = parseDate(startDateStr)
	if err != nil {
		return
	}
	endDate, err = parseDate(endDateStr)
	return
}

func parseDate(dateStr string) (date time.Time, err error) {
	if len(dateStr) == 4 {
		dateStr += "0101"
	} else if len(dateStr) == 6 {
		dateStr += "01"
	}
	return time.Parse(date_util.YMDFormat, dateStr)
}

func init() {
	kTemplate = common.NewTemplate("daterange", kTemplateSpec)
}
