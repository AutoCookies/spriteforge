package slicer

import (
	"image"
	"image/color"
	"testing"

	"pixelc/internal/imageutil"
	"pixelc/pkg/model"
)

func TestSliceSpritesheetCases(t *testing.T) {
	cfg4 := model.Config{Connectivity: 4, Padding: 0, PivotMode: "center", Preset: "unity"}
	cfg8 := model.Config{Connectivity: 8, Padding: 0, PivotMode: "center", Preset: "unity"}

	t.Run("single blob", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 8, 8))
		setOpaque(img, 2, 2)
		setOpaque(img, 2, 3)
		setOpaque(img, 3, 3)
		sprites, err := SliceSpritesheet(img, cfg4)
		if err != nil {
			t.Fatalf("slice failed: %v", err)
		}
		if len(sprites) != 1 {
			t.Fatalf("expected 1 sprite, got %d", len(sprites))
		}
		if sprites[0].X != 2 || sprites[0].Y != 2 || sprites[0].Width != 2 || sprites[0].Height != 2 {
			t.Fatalf("unexpected bounds: %+v", sprites[0])
		}
	})

	t.Run("two separated blobs sorted", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 12, 12))
		setOpaque(img, 8, 1)
		setOpaque(img, 8, 2)
		setOpaque(img, 1, 8)
		setOpaque(img, 2, 8)
		sprites, err := SliceSpritesheet(img, cfg4)
		if err != nil {
			t.Fatalf("slice failed: %v", err)
		}
		if len(sprites) != 2 {
			t.Fatalf("expected 2 sprites, got %d", len(sprites))
		}
		if sprites[0].Y > sprites[1].Y || (sprites[0].Y == sprites[1].Y && sprites[0].X > sprites[1].X) {
			t.Fatalf("sprites are not sorted by minY/minX")
		}
	})

	t.Run("diagonal adjacency differs by connectivity", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 4, 4))
		setOpaque(img, 1, 1)
		setOpaque(img, 2, 2)
		s4, _ := SliceSpritesheet(img, cfg4)
		s8, _ := SliceSpritesheet(img, cfg8)
		if len(s4) != 0 {
			t.Fatalf("expected zero sprites for 4-connectivity due to noise threshold, got %d", len(s4))
		}
		if len(s8) != 1 {
			t.Fatalf("expected one sprite for 8-connectivity, got %d", len(s8))
		}
	})

	t.Run("noise pixel ignored", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 5, 5))
		setOpaque(img, 1, 1)
		sprites, err := SliceSpritesheet(img, cfg4)
		if err != nil {
			t.Fatalf("slice failed: %v", err)
		}
		if len(sprites) != 0 {
			t.Fatalf("expected 0 sprites, got %d", len(sprites))
		}
	})

	t.Run("edge touching boundary", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 4, 4))
		setOpaque(img, 0, 0)
		setOpaque(img, 0, 1)
		sprites, err := SliceSpritesheet(img, cfg4)
		if err != nil {
			t.Fatalf("slice failed: %v", err)
		}
		if len(sprites) != 1 || sprites[0].X != 0 || sprites[0].Y != 0 {
			t.Fatalf("unexpected edge sprite: %+v", sprites)
		}
	})

	t.Run("fully transparent returns empty", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 3, 3))
		sprites, err := SliceSpritesheet(img, cfg4)
		if err != nil {
			t.Fatalf("slice failed: %v", err)
		}
		if len(sprites) != 0 {
			t.Fatalf("expected 0 sprites, got %d", len(sprites))
		}
	})
}

func TestDeterministicOutput(t *testing.T) {
	cfg := model.Config{Connectivity: 4, Padding: 0, PivotMode: "center", Preset: "unity"}
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	for i := 0; i < 5; i++ {
		setOpaque(img, i*3, i)
		setOpaque(img, i*3+1, i)
	}

	a, err := SliceSpritesheet(img, cfg)
	if err != nil {
		t.Fatalf("slice failed: %v", err)
	}
	b, err := SliceSpritesheet(img, cfg)
	if err != nil {
		t.Fatalf("slice failed: %v", err)
	}
	if len(a) != len(b) {
		t.Fatalf("length mismatch %d != %d", len(a), len(b))
	}
	for i := range a {
		ha := imageutil.HashRGBA(a[i].Image)
		hb := imageutil.HashRGBA(b[i].Image)
		if ha != hb {
			t.Fatalf("hash mismatch at %d: %s vs %s", i, ha, hb)
		}
	}
}

func BenchmarkSlicer_100Sprites(b *testing.B) {
	cfg := model.Config{Connectivity: 4, Padding: 0, PivotMode: "center", Preset: "unity"}
	img := image.NewRGBA(image.Rect(0, 0, 220, 220))
	for i := 0; i < 100; i++ {
		x := (i % 10) * 20
		y := (i / 10) * 20
		setOpaque(img, x, y)
		setOpaque(img, x+1, y)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := SliceSpritesheet(img, cfg)
		if err != nil {
			b.Fatalf("slice failed: %v", err)
		}
	}
}

func setOpaque(img *image.RGBA, x, y int) {
	img.SetRGBA(x, y, color.RGBA{R: 255, A: 255})
}
