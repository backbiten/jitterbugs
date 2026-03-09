package report

import (
	"html/template"
	"os"

	"github.com/backbiten/jitterbugs/internal/core"
)

const htmlTmpl = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>QAQC Report</title>
  <style>
    body{font-family:system-ui,sans-serif;max-width:960px;margin:2rem auto;padding:0 1rem;color:#333}
    h1{margin-bottom:.25rem}
    .meta{color:#666;font-size:.9rem;margin-bottom:1.5rem}
    .badge{display:inline-block;padding:.2em .6em;border-radius:4px;font-weight:700;text-transform:uppercase;font-size:.8rem}
    .pass{background:#d4edda;color:#155724}
    .warning{background:#fff3cd;color:#856404}
    .fail{background:#f8d7da;color:#721c24}
    table{width:100%;border-collapse:collapse;margin-top:1rem}
    th,td{border:1px solid #dee2e6;padding:.5rem .75rem;text-align:left;vertical-align:top}
    th{background:#f8f9fa;font-weight:600}
    tr:nth-child(even){background:#fafafa}
    .findings-heading{margin-top:2rem}
    code{font-size:.85em;background:#f4f4f4;padding:.1em .3em;border-radius:3px}
  </style>
</head>
<body>
<h1>QAQC Report</h1>
<div class="meta">
  <strong>Repository:</strong> {{.RepoPath}}<br>
  <strong>Scanned at:</strong> {{.Timestamp.UTC.Format "2006-01-02 15:04:05"}} UTC<br>
  <strong>Overall status:</strong> <span class="badge {{.OverallStatus}}">{{.OverallStatus}}</span>
</div>

<h2>Checks</h2>
<table>
  <tr><th>Check</th><th>Status</th><th>Message</th></tr>
  {{range .Results}}
  <tr>
    <td>{{.Name}}</td>
    <td><span class="badge {{.Status}}">{{.Status}}</span></td>
    <td>{{.Message}}</td>
  </tr>
  {{end}}
</table>

{{if .Findings}}
<h2 class="findings-heading">Findings</h2>
<table>
  <tr><th>File</th><th>Line</th><th>Pattern</th><th>Match (redacted)</th></tr>
  {{range .Findings}}
  <tr>
    <td>{{.File}}</td>
    <td>{{.Line}}</td>
    <td>{{.Pattern}}</td>
    <td><code>{{.Match}}</code></td>
  </tr>
  {{end}}
</table>
{{end}}
</body>
</html>`

// reportView is the template data container.
type reportView struct {
	*core.Report
	Findings []core.Finding
}

// WriteHTML renders the report as an HTML file at the given path.
func WriteHTML(rpt *core.Report, path string) error {
	tmpl, err := template.New("report").Parse(htmlTmpl)
	if err != nil {
		return err
	}

	view := &reportView{Report: rpt}
	for _, r := range rpt.Results {
		view.Findings = append(view.Findings, r.Findings...)
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, view)
}
