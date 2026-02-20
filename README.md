# Pixel Asset Compiler (pixelc)

Pixel Asset Compiler (PAC) is a deterministic sprite-atlas compilation toolchain written in Go.

## Commands

```bash
pixelc compile <input> --out <dir> [--batch] [--dry-run] [--config cfg.json] [--ignore "**/temp/**"] [--fps 12] [--report]
pixelc version
pixelc doctor
```

## Development

```bash
go test ./...
go vet ./...
./scripts/verify.sh
go run ./scripts/no_binary.go
go build ./cmd/pixelc
PIXELC_BIN=./pixelc go test -tags smoketool ./scripts -run TestSmokeHarness -v
go test ./... -bench=. -benchmem -run=^$ > bench.txt
go run -tags benchgate ./scripts/bench_gate.go bench.txt
```

## Packaging

```bash
./scripts/package/package_linux.sh 1.0.0 <commit> <build-date>
./scripts/package/package_macos.sh 1.0.0 <commit> <build-date>
pwsh ./scripts/package/package_windows.ps1 -Version 1.0.0 -Commit <commit> -BuildDate <build-date>
```

Artifacts are generated in `dist/` and must never be committed.
