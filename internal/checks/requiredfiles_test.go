package checks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/backbiten/jitterbugs/internal/core"
)

func TestRequiredFilesCheck_AllPresent(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"README.md", "LICENSE", "SECURITY.md", "CONTRIBUTING.md"} {
		must(t, os.WriteFile(filepath.Join(dir, name), []byte("content"), 0o644))
	}

	result := NewRequiredFilesCheck(nil).Run(dir)
	if result.Status != core.SeverityPass {
		t.Errorf("expected pass, got %s: %s", result.Status, result.Message)
	}
}

func TestRequiredFilesCheck_MissingReadmeAndLicense(t *testing.T) {
	dir := t.TempDir()
	result := NewRequiredFilesCheck(nil).Run(dir)
	if result.Status != core.SeverityFail {
		t.Errorf("expected fail, got %s: %s", result.Status, result.Message)
	}
}

func TestRequiredFilesCheck_MissingOnlyRecommended(t *testing.T) {
	dir := t.TempDir()
	must(t, os.WriteFile(filepath.Join(dir, "README.md"), []byte("x"), 0o644))
	must(t, os.WriteFile(filepath.Join(dir, "LICENSE"), []byte("x"), 0o644))

	result := NewRequiredFilesCheck(nil).Run(dir)
	if result.Status != core.SeverityWarning {
		t.Errorf("expected warning, got %s: %s", result.Status, result.Message)
	}
}

func TestRequiredFilesCheck_ReadmeAlternativeExtension(t *testing.T) {
	dir := t.TempDir()
	must(t, os.WriteFile(filepath.Join(dir, "README.rst"), []byte("x"), 0o644))
	must(t, os.WriteFile(filepath.Join(dir, "LICENSE"), []byte("x"), 0o644))
	must(t, os.WriteFile(filepath.Join(dir, "SECURITY.md"), []byte("x"), 0o644))
	must(t, os.WriteFile(filepath.Join(dir, "CONTRIBUTING.md"), []byte("x"), 0o644))

	result := NewRequiredFilesCheck(nil).Run(dir)
	if result.Status != core.SeverityPass {
		t.Errorf("expected pass with README.rst, got %s: %s", result.Status, result.Message)
	}
}

func TestRequiredFilesCheck_ExtraConfigFiles(t *testing.T) {
	dir := t.TempDir()
	must(t, os.WriteFile(filepath.Join(dir, "README.md"), []byte("x"), 0o644))
	must(t, os.WriteFile(filepath.Join(dir, "LICENSE"), []byte("x"), 0o644))
	must(t, os.WriteFile(filepath.Join(dir, "SECURITY.md"), []byte("x"), 0o644))
	must(t, os.WriteFile(filepath.Join(dir, "CONTRIBUTING.md"), []byte("x"), 0o644))
	// CHANGELOG.md required by config but absent.

	cfg := &core.Config{RequiredFiles: []string{"CHANGELOG.md"}}
	result := NewRequiredFilesCheck(cfg).Run(dir)
	if result.Status != core.SeverityFail {
		t.Errorf("expected fail for missing config-required file, got %s: %s", result.Status, result.Message)
	}
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
