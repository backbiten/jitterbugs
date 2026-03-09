# QAQC – Repository Quality and Security Scanner

`qaqc` is a cross-platform CLI tool that scans a **local git repository** for
common quality and security issues. It produces a JSON report to stdout and an
optional HTML report to a file.

---

## Contents

- [Build & Install](#build--install)
- [Usage](#usage)
- [Check Descriptions](#check-descriptions)
- [Exit Codes](#exit-codes)
- [Configuration](#configuration)
- [Output Format](#output-format)

---

## Build & Install

**Prerequisites:** Go 1.21 or later.

```bash
# Clone the repository
git clone https://github.com/backbiten/jitterbugs.git
cd jitterbugs

# Build the binary (output: ./qaqc)
go build -o qaqc ./cmd/qaqc

# Run tests
go test ./...
```

Install to your `$GOPATH/bin` (optional):

```bash
go install github.com/backbiten/jitterbugs/cmd/qaqc@latest
```

---

## Usage

```
qaqc scan [flags]

Flags:
  --path  <dir>    Path to the git repository to scan (default: current directory)
  --html  <file>   Write an HTML report to this file in addition to JSON stdout
```

### Examples

Scan the current directory:

```bash
./qaqc scan
```

Scan a specific repository:

```bash
./qaqc scan --path /path/to/my-repo
```

Save both JSON (stdout) and HTML reports:

```bash
./qaqc scan --path /path/to/my-repo --html report.html > report.json
```

Use the exit code in a CI pipeline:

```bash
./qaqc scan --path . && echo "All checks passed"
```

---

## Check Descriptions

### 1. Required Files (`required_files`)

Verifies that key community health files exist in the repository root.

| File             | Missing severity |
|------------------|------------------|
| `README` (any extension) | **fail**  |
| `LICENSE` (any extension) | **fail** |
| `SECURITY.md`    | warning          |
| `CONTRIBUTING.md` | warning         |

Additional required files can be declared in [`.qaqc.json`](#configuration).

### 2. CI Configuration (`ci`)

Detects GitHub Actions workflow files under `.github/workflows/`.

- **pass** – at least one `.yml` or `.yaml` workflow file found.
- **warning** – directory absent or contains no workflow files.

### 3. Secret Scan (`secrets`)

Scans all git-tracked files (or all non-`.git` files if `git` is unavailable)
for common secret patterns.

| Pattern                    | Confidence |
|----------------------------|------------|
| AWS Access Key ID (`AKIA…`) | high       |
| Private key PEM header     | high       |
| GitHub token (`ghp_`, `ghs_`, …) | high  |
| Generic secret assignment (`SECRET_KEY = …`, `password: …`) | low |

- **fail** – one or more **high-confidence** matches found.
- **warning** – only low-confidence matches found.
- **pass** – no matches.

Binary files and files larger than 1 MiB are skipped. Secret values are
**redacted** in the report (only the first four characters are shown).

---

## Exit Codes

| Code | Meaning                             |
|------|-------------------------------------|
| `0`  | All checks passed                   |
| `1`  | One or more warnings (no failures)  |
| `2`  | One or more checks failed           |
| `3`  | Internal error (report generation)  |

---

## Configuration

Place a `.qaqc.json` file in the root of the target repository to customise
the scanner behaviour.

```json
{
  "required_files": ["CHANGELOG.md", "CODEOWNERS"],
  "checks": {
    "required_files": true,
    "ci": true,
    "secrets": false
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `required_files` | `string[]` | Extra files that must exist (treated as mandatory – fail if absent). |
| `checks.required_files` | `bool` | Enable/disable the Required Files check (default: `true`). |
| `checks.ci` | `bool` | Enable/disable the CI Configuration check (default: `true`). |
| `checks.secrets` | `bool` | Enable/disable the Secret Scan check (default: `true`). |

---

## Output Format

JSON report structure (emitted to stdout):

```json
{
  "repo_path": "/path/to/repo",
  "timestamp": "2024-01-15T12:00:00Z",
  "overall_status": "warning",
  "results": [
    {
      "name": "Required Files",
      "status": "warning",
      "message": "Missing recommended files: SECURITY.md"
    },
    {
      "name": "CI Configuration",
      "status": "pass",
      "message": "Found 1 GitHub Actions workflow file(s)"
    },
    {
      "name": "Secret Scan",
      "status": "fail",
      "message": "Found 1 potential secret(s) including high-confidence findings",
      "findings": [
        {
          "file": "config/settings.py",
          "line": 12,
          "pattern": "AWS Access Key ID",
          "match": "AKIA****************"
        }
      ]
    }
  ]
}
```
