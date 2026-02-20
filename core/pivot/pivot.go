package pivot

import (
	"fmt"

	"pixelc/pkg/model"
)

func ApplyPivot(s model.Sprite, cfg model.Config) (model.Sprite, error) {
	if s.Image == nil {
		return model.Sprite{}, fmt.Errorf("sprite image is nil")
	}
	if s.Width <= 0 || s.Height <= 0 {
		return model.Sprite{}, fmt.Errorf("sprite dimensions must be positive")
	}

	switch cfg.PivotMode {
	case "center":
		s.PivotX = float64(s.Width/2) / float64(s.Width)
		s.PivotY = float64(s.Height/2) / float64(s.Height)
		return s, nil
	case "bottom-center":
		return applyBottomCenter(s)
	default:
		return model.Sprite{}, fmt.Errorf("unsupported pivot mode: %s", cfg.PivotMode)
	}
}

func applyBottomCenter(s model.Sprite) (model.Sprite, error) {
	bounds := s.Image.Bounds()
	bottomY := -1
	for y := bounds.Max.Y - 1; y >= bounds.Min.Y; y-- {
		if hasOpaquePixelAtY(s, y) {
			bottomY = y
			break
		}
	}
	if bottomY < 0 {
		return model.Sprite{}, fmt.Errorf("sprite is fully transparent")
	}

	sumX := 0
	count := 0
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		if s.Image.RGBAAt(x, bottomY).A > 0 {
			sumX += x
			count++
		}
	}
	if count == 0 {
		return model.Sprite{}, fmt.Errorf("bottom row has no opaque pixels")
	}

	avgX := float64(sumX) / float64(count)
	s.PivotX = (avgX + 0.5) / float64(s.Width)
	s.PivotY = float64(bottomY+1) / float64(s.Height)
	return s, nil
}

func hasOpaquePixelAtY(s model.Sprite, y int) bool {
	bounds := s.Image.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		if s.Image.RGBAAt(x, y).A > 0 {
			return true
		}
	}
	return false
}
