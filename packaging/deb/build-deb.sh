#!/usr/bin/env bash
# Build a Debian .deb package from the PyInstaller-generated binary.
#
# Usage:  bash packaging/deb/build-deb.sh [VERSION]
# Input:  dist/jitterbugs   (built by PyInstaller --onefile)
# Output: packaging/deb/jitterbugs_<VERSION>_amd64.deb

set -euo pipefail

VERSION="${1:-0.0.0}"
ARCH="amd64"
PKG_NAME="jitterbugs_${VERSION}_${ARCH}"
STAGING="packaging/deb/${PKG_NAME}"

echo "Building ${PKG_NAME}.deb …"

# ── Directory layout ─────────────────────────────────────────────────────────
mkdir -p "${STAGING}/DEBIAN"
mkdir -p "${STAGING}/opt/jitterbugs"
mkdir -p "${STAGING}/usr/bin"

# ── Binary ───────────────────────────────────────────────────────────────────
cp dist/jitterbugs "${STAGING}/opt/jitterbugs/jitterbugs"
chmod 755 "${STAGING}/opt/jitterbugs/jitterbugs"

# ── /usr/bin wrapper ─────────────────────────────────────────────────────────
cat > "${STAGING}/usr/bin/jitterbugs" << 'WRAPPER'
#!/bin/sh
exec /opt/jitterbugs/jitterbugs "$@"
WRAPPER
chmod 755 "${STAGING}/usr/bin/jitterbugs"

# ── DEBIAN/control ───────────────────────────────────────────────────────────
INSTALLED_SIZE=$(du -sk "${STAGING}/opt" | awk '{print $1}')
cat > "${STAGING}/DEBIAN/control" << CONTROL
Package: jitterbugs
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: ${ARCH}
Installed-Size: ${INSTALLED_SIZE}
Depends: libc6 (>= 2.17), libgl1, libglib2.0-0
Maintainer: backbiten
Description: Consent-based real-time audio/video overlay demo
 Jitterbugs is a desktop application for real-time audio/video feature
 extraction and overlay visualization. Explicit user consent is required
 to operate.
CONTROL

# ── Build .deb ───────────────────────────────────────────────────────────────
DEB_PATH="packaging/deb/${PKG_NAME}.deb"
dpkg-deb --build --root-owner-group "${STAGING}" "${DEB_PATH}"
echo "✓ Package ready: ${DEB_PATH}"
