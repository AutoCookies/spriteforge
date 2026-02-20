package trim

import (
	"fmt"
	"image"

	"pixelc/pkg/model"
)

func TrimSprite(s model.Sprite) (model.Sprite, error) {
	if s.Image == nil {
		return model.Sprite{}, fmt.Errorf("sprite image is nil")
	}
	bounds := s.Image.Bounds()
	minX, minY := bounds.Max.X, bounds.Max.Y
	maxX, maxY := bounds.Min.X-1, bounds.Min.Y-1

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if s.Image.RGBAAt(x, y).A == 0 {
				continue
			}
			if x < minX {
				minX = x
			}
			if y < minY {
				minY = y
			}
			if x > maxX {
				maxX = x
			}
			if y > maxY {
				maxY = y
			}
		}
	}

	if maxX < minX || maxY < minY {
		return model.Sprite{}, fmt.Errorf("sprite is fully transparent")
	}

	w := maxX - minX + 1
	h := maxY - minY + 1
	trimmed := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			trimmed.SetRGBA(x, y, s.Image.RGBAAt(minX+x, minY+y))
		}
	}

	s.Image = trimmed
	s.X += minX
	s.Y += minY
	s.Width = w
	s.Height = h
	return s, nil
}
