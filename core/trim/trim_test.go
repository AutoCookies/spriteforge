package trim

import (
	"errors"
	"testing"

	"pixelc/core/compiler"
	"pixelc/pkg/model"
)

func TestTrimSpriteNotImplemented(t *testing.T) {
	_, err := TrimSprite(model.Sprite{Name: "a"})
	if !errors.Is(err, compiler.ErrNotImplemented) {
		t.Fatalf("expected ErrNotImplemented, got %v", err)
	}
}
