package compiler

import (
	"image"
	"image/color"
	"path/filepath"
	"testing"

	"pixelc/internal/imageutil"
	"pixelc/pkg/model"
)

func TestCompileReturnsValidationError(t *testing.T) {
	_, _, err := Compile("irrelevant.png", model.Config{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestCompiler_BasicSpritesheet(t *testing.T) {
	tmp := t.TempDir()
	inputPath := filepath.Join(tmp, "input.png")
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	img.SetRGBA(1, 1, color.RGBA{R: 255, A: 255})
	img.SetRGBA(1, 2, color.RGBA{R: 255, A: 255})
	img.SetRGBA(7, 5, color.RGBA{G: 255, A: 255})
	img.SetRGBA(8, 5, color.RGBA{G: 255, A: 255})
	if err := imageutil.SavePNG(inputPath, img); err != nil {
		t.Fatalf("save png: %v", err)
	}

	cfg := model.Config{Connectivity: 4, Padding: 0, PivotMode: "bottom-center", Preset: "unity"}
	atlas, preset, err := Compile(inputPath, cfg)
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}
	if preset != nil {
		t.Fatalf("expected nil preset json in phase 1")
	}
	if atlas.Width != 0 || atlas.Height != 0 {
		t.Fatalf("expected zero atlas dimensions in phase 1")
	}
	if len(atlas.Sprites) != 2 {
		t.Fatalf("expected 2 sprites, got %d", len(atlas.Sprites))
	}
	if atlas.Sprites[0].Sprite.Y > atlas.Sprites[1].Sprite.Y {
		t.Fatalf("sprites out of deterministic order")
	}
	for i, ps := range atlas.Sprites {
		if ps.Sprite.PivotY <= 0 || ps.Sprite.PivotY > 1 || ps.Sprite.PivotX <= 0 || ps.Sprite.PivotX > 1 {
			t.Fatalf("sprite %d has invalid normalized pivot: %+v", i, ps.Sprite)
		}
	}
}
