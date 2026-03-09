package checks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/backbiten/jitterbugs/internal/core"
)

func TestCIDetectCheck_NoWorkflowsDir(t *testing.T) {
	dir := t.TempDir()
	result := NewCIDetectCheck().Run(dir)
	if result.Status != core.SeverityWarning {
		t.Errorf("expected warning when .github/workflows absent, got %s", result.Status)
	}
}

func TestCIDetectCheck_EmptyWorkflowsDir(t *testing.T) {
	dir := t.TempDir()
	must(t, os.MkdirAll(filepath.Join(dir, ".github", "workflows"), 0o755))

	result := NewCIDetectCheck().Run(dir)
	if result.Status != core.SeverityWarning {
		t.Errorf("expected warning for empty workflows dir, got %s", result.Status)
	}
}

func TestCIDetectCheck_WithYMLWorkflow(t *testing.T) {
	dir := t.TempDir()
	wfDir := filepath.Join(dir, ".github", "workflows")
	must(t, os.MkdirAll(wfDir, 0o755))
	must(t, os.WriteFile(filepath.Join(wfDir, "ci.yml"), []byte("name: CI"), 0o644))

	result := NewCIDetectCheck().Run(dir)
	if result.Status != core.SeverityPass {
		t.Errorf("expected pass, got %s: %s", result.Status, result.Message)
	}
}

func TestCIDetectCheck_WithYAMLWorkflow(t *testing.T) {
	dir := t.TempDir()
	wfDir := filepath.Join(dir, ".github", "workflows")
	must(t, os.MkdirAll(wfDir, 0o755))
	must(t, os.WriteFile(filepath.Join(wfDir, "release.yaml"), []byte("name: Release"), 0o644))

	result := NewCIDetectCheck().Run(dir)
	if result.Status != core.SeverityPass {
		t.Errorf("expected pass, got %s: %s", result.Status, result.Message)
	}
}

func TestCIDetectCheck_OnlyNonWorkflowFiles(t *testing.T) {
	dir := t.TempDir()
	wfDir := filepath.Join(dir, ".github", "workflows")
	must(t, os.MkdirAll(wfDir, 0o755))
	must(t, os.WriteFile(filepath.Join(wfDir, "README.md"), []byte("docs"), 0o644))

	result := NewCIDetectCheck().Run(dir)
	if result.Status != core.SeverityWarning {
		t.Errorf("expected warning when only non-yml files present, got %s", result.Status)
	}
}
