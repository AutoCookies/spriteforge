package compiler

import (
	"image"
	"image/color"
	"path/filepath"
	"testing"

	"pixelc/core/slicer"
	"pixelc/internal/imageutil"
	"pixelc/pkg/model"
)

func TestSlicerHashGoldenScaffold(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "in.png")
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	img.SetRGBA(1, 1, color.RGBA{R: 255, A: 255})
	img.SetRGBA(1, 2, color.RGBA{R: 255, A: 255})
	img.SetRGBA(5, 4, color.RGBA{B: 255, A: 255})
	img.SetRGBA(6, 4, color.RGBA{B: 255, A: 255})
	if err := imageutil.SavePNG(path, img); err != nil {
		t.Fatalf("save input: %v", err)
	}

	loaded, err := imageutil.LoadPNG(path)
	if err != nil {
		t.Fatalf("load input: %v", err)
	}
	cfg := model.Config{Connectivity: 4, Padding: 0, PivotMode: "center", Preset: "unity"}
	sprites, err := slicer.SliceSpritesheet(loaded, cfg)
	if err != nil {
		t.Fatalf("slice failed: %v", err)
	}

	if len(sprites) != 2 {
		t.Fatalf("expected 2 sprites, got %d", len(sprites))
	}
	hashes := []string{imageutil.HashRGBA(sprites[0].Image), imageutil.HashRGBA(sprites[1].Image)}
	expected := []string{
		"42af801193fad22a2c6b98d9fe22f22d9c00f8f927539f65d98a41178fa31142",
		"b978f97a664cc8fae883c8185c0e72e4c16c3be1fb8901b5306b20fdb24a34b4",
	}
	for i := range hashes {
		if hashes[i] != expected[i] {
			t.Fatalf("hash mismatch at %d: got %s want %s", i, hashes[i], expected[i])
		}
	}
}
