package checks

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/backbiten/jitterbugs/internal/core"
)

// secretPattern describes a single secret-detection heuristic.
type secretPattern struct {
	name           string
	pattern        *regexp.Regexp
	highConfidence bool
}

// patterns is the set of heuristics applied to every scanned line.
var patterns = []secretPattern{
	{
		name:           "AWS Access Key ID",
		pattern:        regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
		highConfidence: true,
	},
	{
		name:           "Private Key Header",
		pattern:        regexp.MustCompile(`-----BEGIN (RSA |EC |DSA |OPENSSH )?PRIVATE KEY-----`),
		highConfidence: true,
	},
	{
		name:           "GitHub Token",
		pattern:        regexp.MustCompile(`gh[pousr]_[A-Za-z0-9]{36,}`),
		highConfidence: true,
	},
	{
		name:           "Generic Secret Assignment",
		pattern:        regexp.MustCompile(`(?i)(password|passwd|secret[._-]?key|api[._-]?key|apikey|auth[._-]?token)\s*[=:]\s*['"]?[^\s'"]{8,100}`),
		highConfidence: false,
	},
}

// maxScanBytes is the per-file size limit to avoid scanning huge binary blobs.
const maxScanBytes = 1 << 20 // 1 MiB

type secretsCheck struct{}

// NewSecretsCheck returns a check that scans tracked files for secret patterns.
func NewSecretsCheck() core.Check {
	return &secretsCheck{}
}

func (c *secretsCheck) Name() string { return "secrets" }

func (c *secretsCheck) Run(repoPath string) core.CheckResult {
	result := core.CheckResult{Name: "Secret Scan"}

	files, err := trackedFiles(repoPath)
	if err != nil {
		// Fall back to a simple recursive walk when git is unavailable.
		files, err = allFiles(repoPath)
		if err != nil {
			result.Status = core.SeverityWarning
			result.Message = "Could not enumerate repository files for secret scan"
			return result
		}
	}

	var findings []core.Finding
	highConfidenceFound := false

	for _, rel := range files {
		ff, high := scanFile(filepath.Join(repoPath, rel), rel)
		findings = append(findings, ff...)
		if high {
			highConfidenceFound = true
		}
	}

	result.Findings = findings

	switch {
	case len(findings) == 0:
		result.Status = core.SeverityPass
		result.Message = "No secrets detected"
	case highConfidenceFound:
		result.Status = core.SeverityFail
		result.Message = fmt.Sprintf("Found %d potential secret(s) including high-confidence findings", len(findings))
	default:
		result.Status = core.SeverityWarning
		result.Message = fmt.Sprintf("Found %d potential secret pattern(s)", len(findings))
	}

	return result
}

// trackedFiles returns the list of files tracked by git (respects .gitignore).
func trackedFiles(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "-C", repoPath, "ls-files")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var files []string
	sc := bufio.NewScanner(strings.NewReader(string(out)))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

// allFiles walks repoPath recursively, skipping the .git directory.
func allFiles(repoPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip unreadable entries
		}
		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		rel, relErr := filepath.Rel(repoPath, path)
		if relErr != nil {
			return nil
		}
		files = append(files, rel)
		return nil
	})
	return files, err
}

// scanFile reads path and returns any findings together with a flag indicating
// whether any high-confidence pattern was matched. relPath is stored in findings
// so that reports show repository-relative paths.
func scanFile(path, relPath string) ([]core.Finding, bool) {
	f, err := os.Open(path)
	if err != nil {
		return nil, false
	}
	defer f.Close()

	// Read up to maxScanBytes to detect binary content and cap memory use.
	buf := make([]byte, maxScanBytes)
	n, _ := f.Read(buf)
	if n == 0 {
		return nil, false
	}
	data := buf[:n]

	// Heuristic: skip binary files (NUL bytes in first 8 KiB).
	probe := data
	if len(probe) > 8192 {
		probe = probe[:8192]
	}
	if bytes.IndexByte(probe, 0) >= 0 {
		return nil, false
	}

	var findings []core.Finding
	highConfidence := false

	sc := bufio.NewScanner(bytes.NewReader(data))
	lineNum := 0
	for sc.Scan() {
		lineNum++
		line := sc.Text()
		for _, sp := range patterns {
			if sp.pattern.MatchString(line) {
				match := sp.pattern.FindString(line)
				findings = append(findings, core.Finding{
					File:    relPath,
					Line:    lineNum,
					Pattern: sp.name,
					Match:   redact(match),
				})
				if sp.highConfidence {
					highConfidence = true
				}
			}
		}
	}

	return findings, highConfidence
}

// redact returns a partially masked version of the match to avoid leaking
// secrets in reports. The first four characters are preserved; everything
// beyond four characters is replaced with asterisks. Strings shorter than four
// characters are fully masked with the same number of asterisks as the input
// length.
func redact(s string) string {
	if len(s) < 4 {
		return strings.Repeat("*", len(s))
	}
	return s[:4] + strings.Repeat("*", len(s)-4)
}
