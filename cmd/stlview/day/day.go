package day

import (
	"html/template"
	"net/http"
	"time"

	"github.com/keep94/consume2"
	"github.com/keep94/speedtestlogger/cmd/stlview/common"
	"github.com/keep94/speedtestlogger/stl"
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
  td, .normal {
    font-size: 30px;
  }
  td.today {
    font-style: italic;
  }
  </style>
</head>
<body>
  <h1>Speeds for {{.Format .Current}} &nbsp; &nbsp; Build: {{.BuildId}}</h1>
  <a href="{{.Prev .Current}}">prev</a> &nbsp; <a href="{{.Next .Current}}">next</a> &nbsp; <a href="{{.DrillUp .Current}}">up</a>
  <br><br>
  <span class="normal">
  {{with $top := .}}
  Download Average (Mbps): {{with .Summary.DownloadMbps}}{{if .Exists}}{{$top.FormatSpeed .Avg}}{{else}}--{{end}}{{end}}
  <br>
  Upload Average (Mbps): {{with .Summary.UploadMbps}}{{if .Exists}}{{$top.FormatSpeed .Avg}}{{else}}--{{end}}{{end}}
  {{end}}
  </span>
  <br><br>
  <table border=1>
    <tr>
      <th>Timestamp</th>
      <th>Download (Mbps)</th>
      <th>Upload (Mbps)</th>
    </tr>
    {{with $top := .}}
    {{range .Entries}}
    <tr>
      <td>{{$top.FormatTimestamp .Ts}}</td>
      <td align="right">{{$top.FormatSpeed .DownloadMbps}}</td>
      <td align="right">{{$top.FormatSpeed .UploadMbps}}</td>
    </tr>
    {{end}}
    {{end}}
  </table>
</body>
</html>`
)

var (
	kTemplate *template.Template
)

type Handler struct {
	Store    stldb.EntriesRunner
	BuildId  string
	Clock    date_util.Clock
	Location *time.Location
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	current, _ := common.ParseDateParam(
		r.Form.Get(common.Date),
		h.Clock.Now().Unix(),
		h.Location,
		common.Day())
	handler := common.Day()
	var entries []*stl.Entry
	var summary aggregators.Summary
	err := h.Store.Entries(
		nil,
		dates.ToTimestamp(current, h.Location),
		dates.ToTimestamp(handler.End(current), h.Location),
		consume2.Compose(
			consume2.AppendPtrsTo(&entries),
			consume2.Call(summary.Add),
		))
	if err != nil {
		http_util.ReportError(w, "Error reading database", err)
		return
	}
	http_util.WriteTemplate(
		w,
		kTemplate,
		&view{
			common.SpeedFormatter{Precision: 2},
			common.TimestampFormatter{Location: h.Location},
			handler,
			current,
			h.BuildId,
			entries,
			summary,
		},
	)
}

type view struct {
	common.SpeedFormatter
	common.TimestampFormatter
	common.DateHandler
	Current time.Time
	BuildId string
	Entries []*stl.Entry
	Summary aggregators.Summary
}

func init() {
	kTemplate = common.NewTemplate("day", kTemplateSpec)
}
