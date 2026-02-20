package packer

import (
	"errors"
	"testing"

	"pixelc/core/compiler"
	"pixelc/pkg/model"
)

func TestPackNotImplemented(t *testing.T) {
	cfg := model.Config{Connectivity: 4, Padding: 0, PivotMode: "center", Preset: "unity"}
	_, _, err := Pack(nil, cfg)
	if !errors.Is(err, compiler.ErrNotImplemented) {
		t.Fatalf("expected ErrNotImplemented, got %v", err)
	}
}
