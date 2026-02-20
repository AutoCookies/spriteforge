package compiler

import (
	"errors"
	"testing"

	"pixelc/pkg/model"
)

func TestCompileReturnsValidationError(t *testing.T) {
	_, _, err := Compile("fixtures/png/simple_1.png", model.Config{})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if errors.Is(err, ErrNotImplemented) {
		t.Fatal("expected validation error before not implemented")
	}
}

func TestCompileReturnsNotImplemented(t *testing.T) {
	cfg := model.Config{Connectivity: 4, Padding: 2, PivotMode: "bottom-center", Preset: "unity"}
	_, _, err := Compile("fixtures/png/simple_1.png", cfg)
	if !errors.Is(err, ErrNotImplemented) {
		t.Fatalf("expected ErrNotImplemented, got %v", err)
	}
}
