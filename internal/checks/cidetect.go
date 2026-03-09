package checks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/backbiten/jitterbugs/internal/core"
)

type ciDetectCheck struct{}

// NewCIDetectCheck returns a check that verifies GitHub Actions workflow files exist.
func NewCIDetectCheck() core.Check {
	return &ciDetectCheck{}
}

func (c *ciDetectCheck) Name() string { return "ci" }

func (c *ciDetectCheck) Run(repoPath string) core.CheckResult {
	result := core.CheckResult{Name: "CI Configuration"}

	workflowsDir := filepath.Join(repoPath, ".github", "workflows")
	entries, err := os.ReadDir(workflowsDir)
	if err != nil {
		result.Status = core.SeverityWarning
		result.Message = "No GitHub Actions workflows found (.github/workflows/ missing or empty)"
		return result
	}

	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
			count++
		}
	}

	if count == 0 {
		result.Status = core.SeverityWarning
		result.Message = "No GitHub Actions workflow files (.yml/.yaml) found in .github/workflows/"
	} else {
		result.Status = core.SeverityPass
		result.Message = fmt.Sprintf("Found %d GitHub Actions workflow file(s)", count)
	}

	return result
}
