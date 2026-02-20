package slicer

import (
	"errors"
	"image"
	"testing"

	"pixelc/core/compiler"
	"pixelc/pkg/model"
)

func TestSliceSpritesheetNotImplemented(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	cfg := model.Config{Connectivity: 4, Padding: 0, PivotMode: "center", Preset: "unity"}
	_, err := SliceSpritesheet(img, cfg)
	if !errors.Is(err, compiler.ErrNotImplemented) {
		t.Fatalf("expected ErrNotImplemented, got %v", err)
	}
}
