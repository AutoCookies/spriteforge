# Pixel Asset Compiler (pixelc)

Pixel Asset Compiler (PAC) is a deterministic sprite-atlas compilation toolchain written in Go. This repository currently provides a production-ready Phase 0 foundation with strict contracts, stable schemas, boundary guardrails, CI checks, and compile-time-safe stubs for future implementation phases.

## Input support status

- PNG spritesheet input: parser/pipeline implementation coming in later phases.
- Folder-of-PNG-frames input: parser/pipeline implementation coming in later phases.

## CLI usage

```bash
pixelc compile <input> --out <dir> --preset unity --padding 2 --connectivity 4 --pivot bottom-center --power2
```

In Phase 0, `compile` validates inputs/config and returns a clear not-implemented message.

## Development commands

```bash
go test ./...
go vet ./...
go build ./cmd/pixelc
./scripts/verify.sh
```

## Fixtures

Fixture directories are scaffolded in-repo for Phase 0, while binary PNG fixture assets are intentionally omitted from source control. Later phases can generate deterministic fixture images during tests or reintroduce approved assets.
