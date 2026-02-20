package compiler

import (
	"fmt"
	"image"
	"image/color"
	"path/filepath"
	"testing"

	"pixelc/internal/imageutil"
	"pixelc/pkg/model"
)

func TestCompileReturnsValidationError(t *testing.T) {
	_, _, _, err := Compile("irrelevant.png", model.Config{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestCompiler_BasicSpritesheet(t *testing.T) {
	path := makeSpritesheet(t)
	cfg := model.Config{Connectivity: 4, Padding: 1, PivotMode: "bottom-center", Preset: "unity"}
	a1, img1, preset, err := Compile(path, cfg)
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}
	if len(preset) == 0 {
		t.Fatalf("expected preset json")
	}
	if len(a1.Sprites) != 2 || img1 == nil {
		t.Fatalf("unexpected output")
	}
	assertWithinBounds(t, a1)

	a2, img2, preset2, err := Compile(path, cfg)
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}
	if imageutil.HashRGBA(img1) != imageutil.HashRGBA(img2) {
		t.Fatalf("atlas hash mismatch")
	}
	if string(preset) != string(preset2) {
		t.Fatalf("preset json mismatch")
	}
	if !samePlacements(a1, a2) {
		t.Fatalf("placement mismatch")
	}
}

func TestCompiler_FolderInput(t *testing.T) {
	dir := makeFolderFrames(t, 3)
	cfg := model.Config{Connectivity: 4, Padding: 2, PivotMode: "center", Preset: "unity", PowerOfTwo: true}
	a1, img1, preset1, err := Compile(dir, cfg)
	if err != nil {
		t.Fatalf("compile folder failed: %v", err)
	}
	if len(a1.Sprites) != 3 || img1 == nil || len(preset1) == 0 {
		t.Fatalf("unexpected output")
	}
	assertWithinBounds(t, a1)

	a2, img2, preset2, err := Compile(dir, cfg)
	if err != nil {
		t.Fatalf("compile folder failed: %v", err)
	}
	if imageutil.HashRGBA(img1) != imageutil.HashRGBA(img2) {
		t.Fatalf("folder atlas hash mismatch")
	}
	if string(preset1) != string(preset2) {
		t.Fatalf("folder preset mismatch")
	}
	if !samePlacements(a1, a2) {
		t.Fatalf("folder placement mismatch")
	}
}

func BenchmarkCompiler_Folder_200Frames(b *testing.B) {
	dir := makeFolderFramesBench(b, 200)
	cfg := model.Config{Connectivity: 4, Padding: 1, PivotMode: "center", Preset: "unity", PowerOfTwo: true}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, _, _, err := Compile(dir, cfg); err != nil {
			b.Fatalf("compile failed: %v", err)
		}
	}
}

func makeSpritesheet(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "sheet.png")
	img := image.NewRGBA(image.Rect(0, 0, 12, 8))
	img.SetRGBA(1, 1, color.RGBA{R: 255, A: 255})
	img.SetRGBA(1, 2, color.RGBA{R: 255, A: 255})
	img.SetRGBA(7, 5, color.RGBA{G: 255, A: 255})
	img.SetRGBA(8, 5, color.RGBA{G: 255, A: 255})
	if err := imageutil.SavePNG(path, img); err != nil {
		t.Fatalf("save png: %v", err)
	}
	return path
}

func makeFolderFrames(t *testing.T, n int) string {
	t.Helper()
	d := t.TempDir()
	for i := 0; i < n; i++ {
		img := image.NewRGBA(image.Rect(0, 0, 6, 6))
		img.SetRGBA(1, 1, color.RGBA{R: uint8(50 + i), A: 255})
		img.SetRGBA(2, 1, color.RGBA{R: uint8(50 + i), A: 255})
		img.SetRGBA(1, 2, color.RGBA{R: uint8(50 + i), A: 255})
		name := filepath.Join(d, frameName(i))
		if err := imageutil.SavePNG(name, img); err != nil {
			t.Fatalf("save frame: %v", err)
		}
	}
	return d
}

func makeFolderFramesBench(b *testing.B, n int) string {
	b.Helper()
	d := b.TempDir()
	for i := 0; i < n; i++ {
		img := image.NewRGBA(image.Rect(0, 0, 8, 8))
		img.SetRGBA(1, 1, color.RGBA{R: uint8(i), A: 255})
		img.SetRGBA(2, 1, color.RGBA{R: uint8(i), A: 255})
		img.SetRGBA(1, 2, color.RGBA{R: uint8(i), A: 255})
		if err := imageutil.SavePNG(filepath.Join(d, frameName(i)), img); err != nil {
			b.Fatalf("save frame: %v", err)
		}
	}
	return d
}

func frameName(i int) string {
	return fmt.Sprintf("frame_%03d.png", i)
}

func samePlacements(a, b *model.Atlas) bool {
	if len(a.Sprites) != len(b.Sprites) {
		return false
	}
	for i := range a.Sprites {
		x := a.Sprites[i]
		y := b.Sprites[i]
		if x.AtlasX != y.AtlasX || x.AtlasY != y.AtlasY || x.Sprite.Name != y.Sprite.Name {
			return false
		}
	}
	return true
}

func assertWithinBounds(t *testing.T, atlas *model.Atlas) {
	t.Helper()
	for _, ps := range atlas.Sprites {
		if ps.AtlasX < 0 || ps.AtlasY < 0 || ps.AtlasX+ps.Sprite.Width > atlas.Width || ps.AtlasY+ps.Sprite.Height > atlas.Height {
			t.Fatalf("placement out of bounds: %+v in %dx%d", ps, atlas.Width, atlas.Height)
		}
	}
}
