package exporter

import (
	"errors"
	"testing"

	"pixelc/core/compiler"
	"pixelc/pkg/model"
)

func TestExportUnityNotImplemented(t *testing.T) {
	_, err := ExportUnity(model.Atlas{}, "atlas.png", "0.0.0-dev")
	if !errors.Is(err, compiler.ErrNotImplemented) {
		t.Fatalf("expected ErrNotImplemented, got %v", err)
	}
}
