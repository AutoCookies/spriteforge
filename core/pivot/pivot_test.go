package pivot

import (
	"image"
	"image/color"
	"math"
	"testing"

	"pixelc/pkg/model"
)

func TestApplyPivot(t *testing.T) {
	t.Run("center symmetrical shape", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 4, 4))
		img.SetRGBA(1, 1, color.RGBA{A: 255})
		img.SetRGBA(2, 2, color.RGBA{A: 255})
		s := model.Sprite{Image: img, Width: 4, Height: 4}
		out, err := ApplyPivot(s, model.Config{PivotMode: "center"})
		if err != nil {
			t.Fatalf("pivot failed: %v", err)
		}
		assertFloat(t, out.PivotX, 0.5)
		assertFloat(t, out.PivotY, 0.5)
	})

	t.Run("bottom center off-center mass", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 6, 5))
		img.SetRGBA(1, 4, color.RGBA{A: 255})
		img.SetRGBA(2, 4, color.RGBA{A: 255})
		img.SetRGBA(3, 4, color.RGBA{A: 255})
		s := model.Sprite{Image: img, Width: 6, Height: 5}
		out, err := ApplyPivot(s, model.Config{PivotMode: "bottom-center"})
		if err != nil {
			t.Fatalf("pivot failed: %v", err)
		}
		assertFloat(t, out.PivotX, 0.4166666666666667)
		assertFloat(t, out.PivotY, 1.0)
	})

	t.Run("thin vertical sprite", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 1, 5))
		for y := 0; y < 5; y++ {
			img.SetRGBA(0, y, color.RGBA{A: 255})
		}
		s := model.Sprite{Image: img, Width: 1, Height: 5}
		out, err := ApplyPivot(s, model.Config{PivotMode: "bottom-center"})
		if err != nil {
			t.Fatalf("pivot failed: %v", err)
		}
		assertFloat(t, out.PivotX, 0.5)
		assertFloat(t, out.PivotY, 1.0)
	})

	t.Run("single row sprite", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 4, 1))
		img.SetRGBA(2, 0, color.RGBA{A: 255})
		s := model.Sprite{Image: img, Width: 4, Height: 1}
		out, err := ApplyPivot(s, model.Config{PivotMode: "bottom-center"})
		if err != nil {
			t.Fatalf("pivot failed: %v", err)
		}
		assertFloat(t, out.PivotX, 0.625)
		assertFloat(t, out.PivotY, 1.0)
	})
}

func assertFloat(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("got %v want %v", got, want)
	}
}
