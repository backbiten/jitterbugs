# Building Jitterbugs

This document describes how to build the Jitterbugs desktop application locally
and how the CI pipeline produces distributable artifacts.

## Prerequisites

| Platform | Requirements |
|----------|--------------|
| Linux (Debian/Ubuntu/Mint) | Python ≥ 3.10, `portaudio19-dev`, `libgl1`, `libglib2.0-0`, `dpkg-deb` |
| Windows | Python ≥ 3.10 |

### Install system libraries (Linux)

```bash
sudo apt-get update
sudo apt-get install -y portaudio19-dev python3-dev libgl1 libglib2.0-0
```

## Setup

```bash
pip install -r requirements.txt
pip install pyinstaller
```

## Building

### Standalone binary – Linux

```bash
pyinstaller --noconfirm --clean --onefile \
    --name jitterbugs \
    --collect-all mediapipe \
    --collect-all cv2 \
    scripts/run_sniffer.py
# Output: dist/jitterbugs
```

### Debian package (.deb) – Linux

Build the binary first (see above), then:

```bash
bash packaging/deb/build-deb.sh 0.1.0
# Output: packaging/deb/jitterbugs_0.1.0_amd64.deb
```

Install the package:

```bash
sudo dpkg -i packaging/deb/jitterbugs_0.1.0_amd64.deb
# Then run:
jitterbugs
```

### Executable (.exe) – Windows

```powershell
pyinstaller --noconfirm --clean --onefile `
    --name jitterbugs `
    --collect-all mediapipe `
    --collect-all cv2 `
    scripts/run_sniffer.py
# Output: dist\jitterbugs.exe
```

## CI / GitHub Actions

The workflow at `.github/workflows/build-desktop.yml` runs automatically on
every push to `main`, on version tags (`v*`), and on pull requests targeting
`main`.

| Job | Runner | Artifacts uploaded |
|-----|--------|--------------------|
| Linux build | `ubuntu-latest` | `jitterbugs-linux-<version>` (binary), `jitterbugs-deb-<version>` (.deb) |
| Windows build | `windows-latest` | `jitterbugs-windows-<version>` (.exe) |

### Versioning

If the commit carries an exact tag (e.g. `v1.2.3`) the version is taken from
that tag (strip the leading `v`).  Otherwise the version falls back to
`0.0.0+YYYYMMDD.<short-sha>`.

To produce a tagged release:

```bash
git tag v1.2.3
git push origin v1.2.3
```
