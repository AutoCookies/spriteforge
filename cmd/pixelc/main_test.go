package main

import (
	"encoding/json"
	"image"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"pixelc/internal/imageutil"
)

var testBinary string

func TestMain(m *testing.M) {
	wd, _ := os.Getwd()
	root := filepath.Clean(filepath.Join(wd, "../.."))
	testBinary = filepath.Join(root, "pixelc-test-bin")
	build := exec.Command("go", "build", "-o", testBinary, "./cmd/pixelc")
	build.Dir = root
	if out, err := build.CombinedOutput(); err != nil {
		panic("build test binary failed: " + err.Error() + " output=" + string(out))
	}
	code := m.Run()
	_ = os.Remove(testBinary)
	os.Exit(code)
}

func TestHelpExitsZero(t *testing.T) {
	cmd := exec.Command(testBinary, "--help")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("expected help to exit 0 err=%v out=%s", err, out)
	}
}

func TestCompileWritesOutputsAndReport(t *testing.T) {
	input := writeTempPNG(t)
	outDir := filepath.Join(t.TempDir(), "out")
	cmd := exec.Command(testBinary, "compile", input, "--out", outDir, "--preset", "unity", "--padding", "2", "--connectivity", "4", "--pivot", "bottom-center", "--power2", "--report")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected success err=%v out=%s", err, out)
	}
	if !strings.Contains(string(out), "compiled sprites=") {
		t.Fatalf("missing summary output: %s", out)
	}
	assertExists(t, filepath.Join(outDir, "atlas.png"))
	assertExists(t, filepath.Join(outDir, "atlas.json"))
	assertExists(t, filepath.Join(outDir, "report.json"))
	data, _ := os.ReadFile(filepath.Join(outDir, "atlas.json"))
	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("invalid atlas.json: %v", err)
	}
}

func TestCompileDryRunNoFiles(t *testing.T) {
	input := writeTempPNG(t)
	outDir := filepath.Join(t.TempDir(), "out")
	cmd := exec.Command(testBinary, "compile", input, "--out", outDir, "--dry-run")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("dry run failed err=%v out=%s", err, out)
	}
	if _, err := os.Stat(filepath.Join(outDir, "atlas.json")); err == nil {
		t.Fatalf("dry-run wrote files")
	}
}

func TestBatchConfigAndIgnore(t *testing.T) {
	root := t.TempDir()
	writePNGAt(t, filepath.Join(root, "a", "player_idle_001.png"))
	writePNGAt(t, filepath.Join(root, "a", "player_idle_002.png"))
	writePNGAt(t, filepath.Join(root, "skip", "foo_001.png"))
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	_ = os.WriteFile(cfgPath, []byte(`{"padding":1,"connectivity":4,"pivotMode":"center","preset":"unity","fps":15,"ignore":["**/skip/**"]}`), 0o644)
	outDir := filepath.Join(t.TempDir(), "out")
	cmd := exec.Command(testBinary, "compile", root, "--batch", "--out", outDir, "--config", cfgPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("batch failed err=%v out=%s", err, out)
	}
	assertExists(t, filepath.Join(outDir, "a", "atlas.json"))
	if _, err := os.Stat(filepath.Join(outDir, "skip", "atlas.json")); err == nil {
		t.Fatalf("ignore not applied")
	}
}

func TestCompileInvalidConfig(t *testing.T) {
	input := writeTempPNG(t)
	cmd := exec.Command(testBinary, "compile", input, "--out", "out", "--preset", "unity", "--padding", "2", "--connectivity", "6", "--pivot", "bottom-center")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected failure output=%s", out)
	}
	if !strings.Contains(string(out), "config validation error") {
		t.Fatalf("expected validation error output=%s", out)
	}
}

func writeTempPNG(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "input.png")
	writePNGAt(t, path)
	return path
}

func writePNGAt(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	img.SetRGBA(1, 1, color.RGBA{A: 255})
	img.SetRGBA(1, 2, color.RGBA{A: 255})
	if err := imageutil.SavePNG(path, img); err != nil {
		t.Fatalf("save png: %v", err)
	}
}

func assertExists(t *testing.T, p string) {
	t.Helper()
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("missing %s: %v", p, err)
	}
}
