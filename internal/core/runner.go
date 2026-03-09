// Package core provides the check runner and configuration loader.
package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Check is the interface that all scan checks must implement.
type Check interface {
	// Name returns the stable identifier used in config (e.g. "ci").
	Name() string
	// Run executes the check against the repository at repoPath and returns the result.
	Run(repoPath string) CheckResult
}

// Runner orchestrates a set of checks over a repository path.
type Runner struct {
	repoPath string
	cfg      *Config
	checks   []Check
}

// NewRunner creates a Runner for the given repository path and config.
func NewRunner(repoPath string, cfg *Config) *Runner {
	if cfg == nil {
		cfg = &Config{}
	}
	return &Runner{repoPath: repoPath, cfg: cfg}
}

// AddCheck registers a check to be executed by Run.
func (r *Runner) AddCheck(c Check) {
	r.checks = append(r.checks, c)
}

// Run executes all registered checks and returns the aggregated Report.
func (r *Runner) Run() *Report {
	rpt := &Report{
		RepoPath:      r.repoPath,
		Timestamp:     time.Now().UTC(),
		OverallStatus: SeverityPass,
	}

	for _, check := range r.checks {
		if !r.cfg.CheckEnabled(check.Name()) {
			continue
		}
		result := check.Run(r.repoPath)
		rpt.Results = append(rpt.Results, result)

		switch result.Status {
		case SeverityFail:
			rpt.OverallStatus = SeverityFail
		case SeverityWarning:
			if rpt.OverallStatus != SeverityFail {
				rpt.OverallStatus = SeverityWarning
			}
		}
	}

	return rpt
}

// LoadConfig reads the optional .qaqc.json file from repoPath.
// If no config file is found a default (empty) Config is returned.
func LoadConfig(repoPath string) *Config {
	cfg := &Config{}
	path := filepath.Join(repoPath, ".qaqc.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}
	// Ignore unmarshal errors – an unreadable config is treated as default.
	_ = json.Unmarshal(data, cfg)
	return cfg
}
