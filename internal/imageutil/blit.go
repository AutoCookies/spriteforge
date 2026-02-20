package imageutil

import (
	"fmt"
	"image"
)

func Blit(dst *image.RGBA, src *image.RGBA, atX, atY int) error {
	if dst == nil || src == nil {
		return fmt.Errorf("dst and src must be non-nil")
	}
	sb := src.Bounds()
	db := dst.Bounds()
	if atX < db.Min.X || atY < db.Min.Y || atX+sb.Dx() > db.Max.X || atY+sb.Dy() > db.Max.Y {
		return fmt.Errorf("blit out of bounds")
	}
	for y := 0; y < sb.Dy(); y++ {
		for x := 0; x < sb.Dx(); x++ {
			dst.SetRGBA(atX+x, atY+y, src.RGBAAt(sb.Min.X+x, sb.Min.Y+y))
		}
	}
	return nil
}
