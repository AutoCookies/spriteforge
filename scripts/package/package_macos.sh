#!/usr/bin/env bash
set -euo pipefail
VERSION=${1:-1.0.0}; COMMIT=${2:-dev}; BUILDDATE=${3:-1970-01-01T00:00:00Z}
mkdir -p dist
LDFLAGS="-X pixelc/internal/version.Version=${VERSION} -X pixelc/internal/version.Commit=${COMMIT} -X pixelc/internal/version.BuildDate=${BUILDDATE}"
go build -ldflags "$LDFLAGS" -o dist/pixelc ./cmd/pixelc
mkdir -p dist/pixelc.app/Contents/MacOS dist/pixelc.app/Contents
cp dist/pixelc dist/pixelc.app/Contents/MacOS/pixelc
cp scripts/package/macos.Info.plist dist/pixelc.app/Contents/Info.plist
tar -czf "dist/pixelc-${VERSION}-macos-universal.tar.gz" -C dist pixelc
if command -v hdiutil >/dev/null 2>&1; then
  hdiutil create -volname pixelc -srcfolder dist/pixelc.app -ov -format UDZO "dist/pixelc-${VERSION}-macos.dmg"
fi
find dist -maxdepth 1 -type f -print
