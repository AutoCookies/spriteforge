# Pixel Asset Compiler (pixelc)

<div align="center">

<img src="./assets/spriteforge-logo.png" width="200" alt="pixelc logo">

[![Discord](https://img.shields.io/discord/1341453502095478916?label=Discord&logo=discord&logoColor=white&color=5865F2)](https://discord.gg/nnkfW83n)
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Shell](https://img.shields.io/badge/Shell-Bash-4EAA25?logo=gnubash&logoColor=white)](https://www.gnu.org/software/bash/)
[![PowerShell](https://img.shields.io/badge/PowerShell-5.1+-5391FE?logo=powershell&logoColor=white)](https://github.com/PowerShell/PowerShell)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.8-3178C6?logo=typescript&logoColor=white)](https://www.typescriptlang.org/)
[![Electron](https://img.shields.io/badge/Electron-35-47848F?logo=electron&logoColor=white)](https://www.electronjs.org/)
[![License](https://img.shields.io/badge/License-NC--OSL%20v1.0-orange)](./LICENSE)

</div>

**pixelc** is a deterministic, cross-platform sprite-atlas compilation toolchain. It takes PNG sprite frames (or spritesheets), slices, trims, pivots, and packs them into an optimized atlas — then exports metadata ready to drop into your game engine.


It ships as both a **CLI** for pipeline automation and a **desktop GUI** for visual, drag-and-drop workflows.

---

## Table of Contents

- [Features](#features)
- [How It Works](#how-it-works)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [CLI Reference](#cli-reference)
- [Configuration File](#configuration-file)
- [Output Format](#output-format)
- [Desktop GUI](#desktop-gui)
- [Development](#development)
- [Packaging & Release](#packaging--release)
- [Project Structure](#project-structure)

---

## Features

- 🎯 **Deterministic output** — same input always produces the same atlas, bit-for-bit
- ✂️ **Auto-slicing** — accepts a whole spritesheet PNG or a folder of individual frames
- 🔲 **Transparent trim** — strips empty alpha border from each sprite to save atlas space
- 📌 **Pivot points** — configurable pivot per sprite (`center`, `bottom-center`)
- 📦 **Smart packing** — bin-packs sprites with configurable padding
- 🔢 **Power-of-two atlas** — optional constraint for GPU compatibility
- 🎬 **Animation metadata** — infers animation states and FPS from frame filename conventions
- 📤 **Unity export preset** — outputs `atlas.png` + `atlas.json` compatible with Unity's sprite atlas system
- 🗂️ **Batch mode** — recursively compile entire asset directories in one command
- 🧪 **Dry-run mode** — preview output dimensions without writing any files
- 🖥️ **Desktop GUI** — Electron + React app for visual compilation without touching the terminal

---

## How It Works

```
Input PNG / Folder of PNGs
        │
        ▼
   [1] Slice           ← split spritesheet rows/cols, or load individual frames
        │
        ▼
   [2] Trim            ← remove transparent padding from each sprite
        │
        ▼
   [3] Pivot           ← compute pivot point (center / bottom-center)
        │
        ▼
   [4] Pack            ← bin-pack all sprites onto an atlas with padding
        │
        ▼
   [5] Export          ← write atlas.png + atlas.json (Unity preset)
```

---

## Installation

### From Source (requires Go 1.22+)

```bash
git clone https://github.com/your-org/spriteforge
cd spriteforge

# Build the CLI (output to ./pixelc_bin to avoid conflict with source dir)
go build -o pixelc_bin ./cmd/pixelc

# Optionally add to PATH
sudo mv pixelc_bin /usr/local/bin/pixelc
```

### Pre-built Binaries

Download the latest release for your platform from the [Releases page](https://github.com/your-org/spriteforge/releases).

| Platform | File |
|---|---|
| Linux (x86_64) | `pixelc-linux-amd64.tar.gz` |
| macOS | `pixelc-macos.tar.gz` |
| Windows | `pixelc-windows.zip` |

---

## Quick Start

### Compile a single spritesheet

```bash
pixelc compile hero_walk.png --out ./out
```

### Compile a folder of frames

```bash
pixelc compile ./frames/hero_walk/ --out ./out
```

### Batch compile all asset folders recursively

```bash
pixelc compile ./assets/ --out ./out --batch
```

### Preview without writing files

```bash
pixelc compile hero_walk.png --out ./out --dry-run
```

---

## CLI Reference

### `pixelc compile <input> --out <dir> [flags]`

| Flag | Default | Description |
|---|---|---|
| `--out <dir>` | *(required)* | Output directory for `atlas.png` and `atlas.json` |
| `--preset <name>` | `unity` | Export preset. Currently supported: `unity` |
| `--connectivity <4\|8>` | `4` | Pixel connectivity for sprite boundary detection |
| `--padding <n>` | `0` | Padding in pixels between sprites on the atlas |
| `--pivot <mode>` | `center` | Pivot point mode: `center` or `bottom-center` |
| `--power2` | `false` | Force atlas dimensions to be powers of two |
| `--fps <n>` | `12` | Frames per second written into animation metadata |
| `--batch` | `false` | Recursively compile subdirectories as separate atlases |
| `--dry-run` | `false` | Plan and print output without writing any files |
| `--report` | `false` | Write a `report.json` alongside the atlas outputs |
| `--ignore <glob>` | — | Glob pattern to exclude from batch mode (repeatable) |
| `--config <file>` | — | Path to a JSON config file (see below) |

### `pixelc version`

Prints the current version string.

### `pixelc doctor`

Checks the environment and confirms pixelc is correctly installed.

```
doctor ok  os=linux  arch=amd64  temp=/tmp/pixelc-doctor-...
```

---

## Configuration File

Instead of passing flags every time, you can save your settings in a JSON config file:

```json
{
  "connectivity": 4,
  "padding": 2,
  "pivotMode": "center",
  "powerOfTwo": false,
  "preset": "unity",
  "fps": 12,
  "ignore": ["**/temp/**", "**/unused/**"]
}
```

Then reference it with:

```bash
pixelc compile ./assets/ --out ./out --batch --config cfg.json
```

Command-line flags take priority over config file values.

---

## Output Format

pixelc writes two files to the output directory:

### `atlas.png`
A tightly packed PNG sprite atlas containing all input sprites.

### `atlas.json` (Unity preset)
A JSON file describing each sprite's position, dimensions, pivot, and any detected animation states:

```json
{
  "frames": {
    "hero_walk_0": {
      "frame": { "x": 0, "y": 0, "w": 48, "h": 64 },
      "pivot": { "x": 0.5, "y": 0.5 }
    },
    "hero_walk_1": {
      "frame": { "x": 48, "y": 0, "w": 48, "h": 64 },
      "pivot": { "x": 0.5, "y": 0.5 }
    }
  },
  "animations": {
    "hero_walk": {
      "fps": 12,
      "frames": ["hero_walk_0", "hero_walk_1"]
    }
  },
  "meta": {
    "app": "pixelc",
    "version": "1.0.0",
    "image": "atlas.png",
    "size": { "w": 256, "h": 256 }
  }
}
```

> **Animation detection**: pixelc infers animations from frame filenames. Frames named `hero_walk_0.png`, `hero_walk_1.png` are grouped into an animation called `hero_walk` automatically.

---

## Desktop GUI

For a visual workflow, pixelc ships with a desktop application built on **Electron + React**.

### Requirements

- Node.js 20+
- npm

### Run in development mode

```bash
cd pixelc/apps/gui
npm install
npm run dev
```

The app window will open on your desktop. It lets you:
- Browse for your input file or folder via a native file picker
- Configure all compile options visually (preset, padding, pivot, FPS, etc.)
- Watch real-time compile logs as pixelc runs
- Open the output folder when done

### Build a distributable package

```bash
npm run package           # builds for current platform
npm run package:dir       # unpacked build (no installer)
```

Packaged apps are written to `dist/`. The pixelc CLI binary is automatically bundled inside the app.

---

## Development

### Running tests

```bash
go test ./...
```

### Linting & vetting

```bash
go vet ./...
./scripts/verify.sh
```

### Ensuring no binary files were committed

```bash
go run ./scripts/no_binary.go
```

### Smoke tests (requires a built binary)

```bash
PIXELC_BIN=./pixelc_bin go test -tags smoketool ./scripts -run TestSmokeHarness -v
```

### Benchmarks

```bash
go test ./... -bench=. -benchmem -run=^$ > bench.txt
go run -tags benchgate ./scripts/bench_gate.go bench.txt
```

---

## Packaging & Release

Cross-platform release packages are produced by platform-specific scripts in `scripts/package/`.

```bash
# Linux
./scripts/package/package_linux.sh 1.0.0 <commit> <build-date>

# macOS
./scripts/package/package_macos.sh 1.0.0 <commit> <build-date>

# Windows (PowerShell)
pwsh ./scripts/package/package_windows.ps1 -Version 1.0.0 -Commit <commit> -BuildDate <build-date>
```

Artifacts are written to `dist/` and **must never be committed** to the repository.
---

## License

[NC-OSL v1.0 — Non-Commercial, Share-Alike](./LICENSE) © Pomaieco
