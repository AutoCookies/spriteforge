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

func TestVersionAndDoctor(t *testing.T) {
	v := exec.Command(testBinary, "version")
	vout, err := v.CombinedOutput()
	if err != nil || !strings.Contains(string(vout), "pixelc") {
		t.Fatalf("version failed err=%v out=%s", err, vout)
	}
	d := exec.Command(testBinary, "doctor")
	dout, err := d.CombinedOutput()
	if err != nil || !strings.Contains(string(dout), "doctor ok") {
		t.Fatalf("doctor failed err=%v out=%s", err, dout)
	}
}

func TestCompileWritesOutputsAndReport(t *testing.T) {
	input := writeTempPNG(t)
	outDir := filepath.Join(t.TempDir(), "out")
	cmd := exec.Command(testBinary, "compile", input, "--out", outDir, "--report")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected success err=%v out=%s", err, out)
	}
	for _, f := range []string{"atlas.png", "atlas.json", "report.json"} {
		assertExists(t, filepath.Join(outDir, f))
	}
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

func writeTempPNG(t *testing.T) string {
	p := filepath.Join(t.TempDir(), "input.png")
	writePNGAt(t, p)
	return p
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
