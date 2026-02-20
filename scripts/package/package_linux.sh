#!/usr/bin/env bash
set -euo pipefail
VERSION=${1:-1.0.0}; COMMIT=${2:-dev}; BUILDDATE=${3:-1970-01-01T00:00:00Z}
mkdir -p dist
LDFLAGS="-X pixelc/internal/version.Version=${VERSION} -X pixelc/internal/version.Commit=${COMMIT} -X pixelc/internal/version.BuildDate=${BUILDDATE}"
go build -ldflags "$LDFLAGS" -o dist/pixelc ./cmd/pixelc
tar -czf "dist/pixelc-${VERSION}-linux-x64.tar.gz" -C dist pixelc
PKGROOT=$(mktemp -d)
mkdir -p "$PKGROOT/DEBIAN" "$PKGROOT/usr/bin" "$PKGROOT/usr/share/doc/pixelc"
cp dist/pixelc "$PKGROOT/usr/bin/pixelc"
cp README.md LICENSE "$PKGROOT/usr/share/doc/pixelc/"
cat > "$PKGROOT/DEBIAN/control" <<CTRL
Package: pixelc
Version: ${VERSION}
Section: utils
Priority: optional
Architecture: amd64
Maintainer: pixelc
Description: Pixel Asset Compiler CLI
CTRL
dpkg-deb --build "$PKGROOT" "dist/pixelc_${VERSION}_amd64.deb"
rm -rf "$PKGROOT"
find dist -maxdepth 1 -type f -print
