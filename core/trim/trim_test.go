package trim

import (
	"image"
	"image/color"
	"testing"

	"pixelc/pkg/model"
)

func TestTrimSprite(t *testing.T) {
	t.Run("already trimmed", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 2, 2))
		fillOpaque(img)
		s := model.Sprite{Image: img, X: 4, Y: 5, Width: 2, Height: 2}
		out, err := TrimSprite(s)
		if err != nil {
			t.Fatalf("trim failed: %v", err)
		}
		if out.X != 4 || out.Y != 5 || out.Width != 2 || out.Height != 2 {
			t.Fatalf("unexpected sprite: %+v", out)
		}
	})

	t.Run("fully transparent returns error", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 3, 3))
		_, err := TrimSprite(model.Sprite{Image: img, Width: 3, Height: 3})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("one by one", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 1, 1))
		img.SetRGBA(0, 0, color.RGBA{R: 1, A: 255})
		out, err := TrimSprite(model.Sprite{Image: img, Width: 1, Height: 1})
		if err != nil {
			t.Fatalf("trim failed: %v", err)
		}
		if out.Width != 1 || out.Height != 1 {
			t.Fatalf("unexpected dimensions: %+v", out)
		}
	})

	t.Run("sprite touching edges", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 4, 4))
		img.SetRGBA(0, 1, color.RGBA{G: 255, A: 255})
		img.SetRGBA(3, 2, color.RGBA{G: 255, A: 255})
		s := model.Sprite{Image: img, X: 10, Y: 20, Width: 4, Height: 4}
		out, err := TrimSprite(s)
		if err != nil {
			t.Fatalf("trim failed: %v", err)
		}
		if out.X != 10 || out.Y != 21 || out.Width != 4 || out.Height != 2 {
			t.Fatalf("unexpected trim result: %+v", out)
		}
	})
}

func fillOpaque(img *image.RGBA) {
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.SetRGBA(x, y, color.RGBA{R: 1, G: 1, B: 1, A: 255})
		}
	}
}
