package packer

import (
	"image"
	"pixelc/core/compiler"
	"pixelc/pkg/model"
)

func Pack(sprites []model.Sprite, cfg model.Config) (model.Atlas, *image.RGBA, error) {
	return model.Atlas{}, nil, compiler.ErrNotImplemented
}
