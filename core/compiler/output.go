package compiler

import (
	"fmt"
	"image"
	"os"
	"path/filepath"

	"pixelc/internal/imageutil"
)

func WriteOutputs(outDir string, atlasImg *image.RGBA, presetJSON []byte) error {
	if atlasImg == nil {
		return fmt.Errorf("nil atlas image")
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	if err := imageutil.SavePNG(filepath.Join(outDir, "atlas.png"), atlasImg); err != nil {
		return fmt.Errorf("write atlas.png: %w", err)
	}
	if err := os.WriteFile(filepath.Join(outDir, "atlas.json"), presetJSON, 0o644); err != nil {
		return fmt.Errorf("write atlas.json: %w", err)
	}
	return nil
}
