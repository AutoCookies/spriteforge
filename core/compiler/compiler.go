package compiler

import (
	"fmt"

	"pixelc/core/pivot"
	"pixelc/core/slicer"
	"pixelc/core/trim"
	"pixelc/internal/imageutil"
	"pixelc/pkg/model"
)

func Compile(inputPath string, cfg model.Config) (*model.Atlas, []byte, error) {
	if err := cfg.Validate(); err != nil {
		return nil, nil, err
	}
	img, err := imageutil.LoadPNG(inputPath)
	if err != nil {
		return nil, nil, fmt.Errorf("load input image: %w", err)
	}
	sprites, err := slicer.SliceSpritesheet(img, cfg)
	if err != nil {
		return nil, nil, err
	}

	processed := make([]model.PlacedSprite, 0, len(sprites))
	for _, s := range sprites {
		trimmed, err := trim.TrimSprite(s)
		if err != nil {
			return nil, nil, err
		}
		pivoted, err := pivot.ApplyPivot(trimmed, cfg)
		if err != nil {
			return nil, nil, err
		}
		processed = append(processed, model.PlacedSprite{Sprite: pivoted})
	}

	atlas := &model.Atlas{Width: 0, Height: 0, Sprites: processed}
	return atlas, nil, nil
}
