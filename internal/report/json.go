// Package report provides JSON and HTML rendering for QAQC scan reports.
package report

import (
	"encoding/json"

	"github.com/backbiten/jitterbugs/internal/core"
)

// RenderJSON returns the report serialised as indented JSON.
func RenderJSON(rpt *core.Report) ([]byte, error) {
	return json.MarshalIndent(rpt, "", "  ")
}
