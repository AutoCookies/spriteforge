package compiler

import (
	"path/filepath"
	"testing"

	"pixelc/internal/testutil"
	"pixelc/pkg/model"
)

func TestCompileGoldenScaffold(t *testing.T) {
	t.Skip("golden tests enabled in Phase 2")

	input := filepath.FromSlash("fixtures/png/simple_1.png")
	cfg := model.Config{Connectivity: 4, Padding: 2, PivotMode: "center", Preset: "unity"}
	expectedAtlasPNG := filepath.FromSlash("fixtures/golden/simple/atlas.png")
	expectedAtlasJSON := filepath.FromSlash("fixtures/golden/simple/atlas.json")

	_ = input
	_ = cfg
	_ = expectedAtlasPNG
	_ = expectedAtlasJSON
	_ = testutil.HashBytes(nil)
}
