//go:build smoketool

package main

import (
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestSmokeHarness(t *testing.T) {
	bin := os.Getenv("PIXELC_BIN")
	if bin == "" {
		t.Fatal("PIXELC_BIN is required")
	}

	if !filepath.IsAbs(bin) {
		if abs, err := filepath.Abs(bin); err == nil {
			bin = abs
		}
	}
	if _, err := os.Stat(bin); err != nil {
		alt := filepath.Clean(filepath.Join("..", filepath.Base(bin)))
		if _, e2 := os.Stat(alt); e2 == nil {
			if abs, e3 := filepath.Abs(alt); e3 == nil {
				bin = abs
			}
		}
	}
	work := t.TempDir()
	input := filepath.Join(work, "input.png")
	outDir := filepath.Join(work, "out")
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	img.SetRGBA(1, 1, color.RGBA{R: 255, A: 255})
	img.SetRGBA(1, 2, color.RGBA{R: 255, A: 255})
	f, _ := os.Create(input)
	_ = png.Encode(f, img)
	_ = f.Close()

	cmd := exec.Command(bin, "compile", input, "--out", outDir, "--preset", "unity")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("smoke compile failed: %v\n%s", err, out)
	}
	for _, fn := range []string{"atlas.png", "atlas.json"} {
		if _, err := os.Stat(filepath.Join(outDir, fn)); err != nil {
			t.Fatalf("missing output: %s", fn)
		}
	}
	data, err := os.ReadFile(filepath.Join(outDir, "atlas.json"))
	if err != nil {
		t.Fatal(err)
	}
	var v map[string]any
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("invalid atlas.json: %v", err)
	}
}
