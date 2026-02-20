package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var blockedExt = map[string]struct{}{
	".png": {}, ".jpg": {}, ".jpeg": {}, ".gif": {}, ".webp": {}, ".bmp": {}, ".ico": {},
	".exe": {}, ".dll": {}, ".dylib": {}, ".so": {}, ".a": {}, ".o": {},
	".zip": {}, ".gz": {}, ".bz2": {}, ".xz": {}, ".7z": {}, ".pdf": {}, ".msi": {}, ".dmg": {}, ".pkg": {}, ".deb": {}, ".rpm": {}, ".appimage": {},
}

func main() {
	out, err := exec.Command("git", "ls-files").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "git ls-files failed: %v\n", err)
		os.Exit(1)
	}
	files := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, f := range files {
		if f == "" {
			continue
		}
		ext := strings.ToLower(filepath.Ext(f))
		if _, blocked := blockedExt[ext]; blocked {
			fmt.Fprintf(os.Stderr, "blocked binary extension: %s\n", f)
			os.Exit(1)
		}
		data, err := os.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read %s failed: %v\n", f, err)
			os.Exit(1)
		}
		sample := data
		if len(sample) > 4096 {
			sample = sample[:4096]
		}
		if bytes.IndexByte(sample, 0) >= 0 {
			fmt.Fprintf(os.Stderr, "binary content detected (NUL byte): %s\n", f)
			os.Exit(1)
		}
	}
	fmt.Println("no binary files detected")
}
