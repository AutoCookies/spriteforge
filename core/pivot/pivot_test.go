package pivot

import (
	"errors"
	"testing"

	"pixelc/core/compiler"
	"pixelc/pkg/model"
)

func TestApplyPivotNotImplemented(t *testing.T) {
	cfg := model.Config{Connectivity: 4, Padding: 0, PivotMode: "center", Preset: "unity"}
	_, err := ApplyPivot(model.Sprite{Name: "a"}, cfg)
	if !errors.Is(err, compiler.ErrNotImplemented) {
		t.Fatalf("expected ErrNotImplemented, got %v", err)
	}
}
