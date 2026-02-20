package compiler

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"pixelc/internal/imageutil"
	"pixelc/pkg/model"
)

func TestCompileBatchDeterministicAndIgnore(t *testing.T) {
	root := t.TempDir()
	mkpng(t, filepath.Join(root, "a", "hero_idle_001.png"), color.RGBA{R: 255, A: 255})
	mkpng(t, filepath.Join(root, "a", "hero_idle_002.png"), color.RGBA{G: 255, A: 255})
	mkpng(t, filepath.Join(root, "b", "enemy_run_001.png"), color.RGBA{B: 255, A: 255})
	mkpng(t, filepath.Join(root, "temp", "skip_001.png"), color.RGBA{A: 255})

	cfg := model.Config{Connectivity: 4, Padding: 1, PivotMode: "center", Preset: "unity", FPS: 12}
	out := filepath.Join(t.TempDir(), "out")
	r1, err := CompileBatch(root, cfg, BatchOptions{OutDir: out, IgnorePatterns: []string{"**/temp/**"}, WriteReport: true})
	if err != nil {
		t.Fatalf("batch compile failed: %v", err)
	}
	if len(r1.Units) != 2 {
		t.Fatalf("expected 2 units got %d", len(r1.Units))
	}
	if r1.Units[0].UnitName > r1.Units[1].UnitName {
		t.Fatalf("units not sorted")
	}
	for _, u := range r1.Units {
		if _, err := os.Stat(filepath.Join(u.OutDir, "atlas.json")); err != nil {
			t.Fatalf("missing output for %s", u.UnitName)
		}
		if _, err := os.Stat(filepath.Join(u.OutDir, "report.json")); err != nil {
			t.Fatalf("missing report for %s", u.UnitName)
		}
	}

	r2, err := CompileBatch(root, cfg, BatchOptions{OutDir: filepath.Join(t.TempDir(), "out2"), IgnorePatterns: []string{"**/temp/**"}})
	if err != nil {
		t.Fatalf("batch compile failed: %v", err)
	}
	if len(r1.Units) != len(r2.Units) {
		t.Fatalf("non-deterministic unit count")
	}
}

func TestCompileBatchDryRun(t *testing.T) {
	root := t.TempDir()
	mkpng(t, filepath.Join(root, "u", "x_001.png"), color.RGBA{R: 255, A: 255})
	cfg := model.Config{Connectivity: 4, Padding: 1, PivotMode: "center", Preset: "unity", FPS: 12}
	out := filepath.Join(t.TempDir(), "out")
	res, err := CompileBatch(root, cfg, BatchOptions{OutDir: out, DryRun: true})
	if err != nil {
		t.Fatalf("dry run failed: %v", err)
	}
	if len(res.Units) != 1 {
		t.Fatalf("expected 1 unit")
	}
	if _, err := os.Stat(filepath.Join(out, "u", "atlas.json")); err == nil {
		t.Fatalf("dry run should not write files")
	}
}

func mkpng(t *testing.T, path string, c color.RGBA) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	img.SetRGBA(1, 1, c)
	img.SetRGBA(2, 1, c)
	if err := imageutil.SavePNG(path, img); err != nil {
		t.Fatal(err)
	}
}
