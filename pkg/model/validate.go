package model

import (
	"fmt"
)

func (c Config) Validate() error {
	if c.Connectivity != 4 && c.Connectivity != 8 {
		return fmt.Errorf("connectivity must be 4 or 8")
	}
	if c.Padding < 0 {
		return fmt.Errorf("padding must be >= 0")
	}
	if c.PivotMode != "center" && c.PivotMode != "bottom-center" {
		return fmt.Errorf("pivot must be center or bottom-center")
	}
	if c.Preset != "unity" {
		return fmt.Errorf("preset must be unity")
	}
	return nil
}

func (a Atlas) Validate() error {
	if a.Width < 0 || a.Height < 0 {
		return fmt.Errorf("atlas dimensions must be non-negative")
	}
	for i, ps := range a.Sprites {
		if ps.AtlasX < 0 || ps.AtlasY < 0 {
			return fmt.Errorf("sprite %d atlas position must be non-negative", i)
		}
		if ps.Sprite.Width < 0 || ps.Sprite.Height < 0 {
			return fmt.Errorf("sprite %d dimensions must be non-negative", i)
		}
		if ps.Sprite.X < 0 || ps.Sprite.Y < 0 {
			return fmt.Errorf("sprite %d source position must be non-negative", i)
		}
	}
	return nil
}
