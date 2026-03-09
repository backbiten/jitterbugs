package checks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/backbiten/jitterbugs/internal/core"
)

func TestSecretsCheck_NoSecrets(t *testing.T) {
	dir := t.TempDir()
	must(t, os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Hello World\nThis is a readme.\n"), 0o644))

	result := NewSecretsCheck().Run(dir)
	if result.Status != core.SeverityPass {
		t.Errorf("expected pass, got %s: %s", result.Status, result.Message)
	}
}

func TestSecretsCheck_AWSAccessKey(t *testing.T) {
	dir := t.TempDir()
	must(t, os.WriteFile(filepath.Join(dir, "config.go"), []byte(`key = "AKIAIOSFODNN7EXAMPLE"`), 0o644))

	result := NewSecretsCheck().Run(dir)
	if result.Status != core.SeverityFail {
		t.Errorf("expected fail for AWS key, got %s: %s", result.Status, result.Message)
	}
	if len(result.Findings) == 0 {
		t.Error("expected at least one finding")
	}
}

func TestSecretsCheck_PrivateKeyHeader(t *testing.T) {
	dir := t.TempDir()
	content := "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA...\n-----END RSA PRIVATE KEY-----\n"
	must(t, os.WriteFile(filepath.Join(dir, "key.pem"), []byte(content), 0o644))

	result := NewSecretsCheck().Run(dir)
	if result.Status != core.SeverityFail {
		t.Errorf("expected fail for private key, got %s: %s", result.Status, result.Message)
	}
}

func TestSecretsCheck_GitHubToken(t *testing.T) {
	dir := t.TempDir()
	// 40-char token after prefix → total ≥ 43 chars, satisfies gh[pousr]_[A-Za-z0-9]{36,}
	token := "ghp_" + strings.Repeat("A", 36)
	must(t, os.WriteFile(filepath.Join(dir, "script.sh"), []byte("TOKEN="+token), 0o644))

	result := NewSecretsCheck().Run(dir)
	if result.Status != core.SeverityFail {
		t.Errorf("expected fail for GitHub token, got %s: %s", result.Status, result.Message)
	}
}

func TestSecretsCheck_GenericSecretWarning(t *testing.T) {
	dir := t.TempDir()
	// Matches the generic low-confidence pattern but not a high-confidence one.
	must(t, os.WriteFile(filepath.Join(dir, "settings.py"), []byte(`SECRET_KEY = "mysupersecretvalue123"`), 0o644))

	result := NewSecretsCheck().Run(dir)
	// generic pattern is low confidence → warning, not fail
	if result.Status == core.SeverityFail {
		t.Errorf("expected warning (not fail) for generic secret pattern, got %s", result.Status)
	}
	if result.Status == core.SeverityPass {
		t.Errorf("expected warning for generic secret pattern, got pass")
	}
}

func TestSecretsCheck_RedactsMatch(t *testing.T) {
	dir := t.TempDir()
	must(t, os.WriteFile(filepath.Join(dir, "env"), []byte(`KEY=AKIAIOSFODNN7EXAMPLE`), 0o644))

	result := NewSecretsCheck().Run(dir)
	for _, f := range result.Findings {
		if strings.Contains(f.Match, "AKIAIOSFODNN7EXAMPLE") {
			t.Error("finding should not expose the full secret value")
		}
	}
}

func TestSecretsCheck_SkipBinaryFile(t *testing.T) {
	dir := t.TempDir()
	// Write a file with NUL bytes (binary) that would match a pattern if treated as text.
	content := []byte("AKIAIOSFODNN7EXAMPLE\x00\x00\x00binary content")
	must(t, os.WriteFile(filepath.Join(dir, "image.png"), content, 0o644))

	result := NewSecretsCheck().Run(dir)
	if result.Status == core.SeverityFail {
		t.Errorf("binary files should be skipped, got fail: %s", result.Message)
	}
}

func TestRedact(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"AKIAIOSFODNN7EXAMPLE", "AKIA" + strings.Repeat("*", 16)},
		{"ab", "**"},   // shorter than 4 – fully masked, length preserved
		{"abc", "***"}, // shorter than 4 – fully masked, length preserved
		{"abcd", "abcd"}, // exactly 4 chars – first 4 shown, 0 stars appended
	}
	for _, tc := range cases {
		got := redact(tc.input)
		if got != tc.want {
			t.Errorf("redact(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}
