# Pixel Asset Compiler (pixelc)

Pixel Asset Compiler (PAC) is a deterministic sprite-atlas compilation toolchain written in Go.

## Input support status

- PNG spritesheet input: supported.
- Folder-of-PNG-frames input: supported.
- Recursive batch folder compilation: supported with `--batch`.

## CLI usage

```bash
pixelc compile <input> --out <dir> --preset unity --padding 2 --connectivity 4 --pivot bottom-center --power2
pixelc compile <input_dir> --batch --out <dir> --ignore "**/temp/**" --fps 12 --report
pixelc compile <input> --out <dir> --config config.json --dry-run
```

Phase 4 adds batch compile, animation grouping (filename-based), config loading, dry-run mode, ignore patterns, and optional report output.

## Development commands

```bash
go test ./...
go vet ./...
go build ./cmd/pixelc
./scripts/verify.sh
go run ./scripts/no_binary.go
```

## Fixtures

Fixture directories are scaffolded in-repo. Binary fixture assets are intentionally not tracked.
