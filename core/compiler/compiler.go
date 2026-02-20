package compiler

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"pixelc/core/packer"
	"pixelc/core/pivot"
	"pixelc/core/slicer"
	"pixelc/core/trim"
	"pixelc/internal/imageutil"
	"pixelc/pkg/model"
)

func Compile(inputPath string, cfg model.Config) (*model.Atlas, *image.RGBA, []byte, error) {
	if err := cfg.Validate(); err != nil {
		return nil, nil, nil, err
	}

	info, err := os.Stat(inputPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("stat input: %w", err)
	}

	sprites, err := loadSprites(inputPath, info.IsDir(), cfg)
	if err != nil {
		return nil, nil, nil, err
	}

	atlas, atlasImg, err := packer.Pack(sprites, cfg)
	if err != nil {
		return nil, nil, nil, err
	}
	return &atlas, atlasImg, nil, nil
}

func loadSprites(inputPath string, isDir bool, cfg model.Config) ([]model.Sprite, error) {
	if isDir {
		return loadFolderSprites(inputPath, cfg)
	}
	if strings.EqualFold(filepath.Ext(inputPath), ".png") {
		img, err := imageutil.LoadPNG(inputPath)
		if err != nil {
			return nil, fmt.Errorf("load spritesheet: %w", err)
		}
		sprites, err := slicer.SliceSpritesheet(img, cfg)
		if err != nil {
			return nil, err
		}
		return processSprites(sprites, cfg)
	}
	return nil, fmt.Errorf("unsupported input: expected .png file or directory")
}

func loadFolderSprites(dir string, cfg model.Config) ([]model.Sprite, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read folder: %w", err)
	}
	files := make([]string, 0)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.EqualFold(filepath.Ext(e.Name()), ".png") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	sprites := make([]model.Sprite, 0, len(files))
	for _, f := range files {
		path := filepath.Join(dir, f)
		img, err := imageutil.LoadPNG(path)
		if err != nil {
			return nil, fmt.Errorf("load frame %s: %w", f, err)
		}
		sprites = append(sprites, model.Sprite{Name: strings.TrimSuffix(f, filepath.Ext(f)), Image: img, Width: img.Bounds().Dx(), Height: img.Bounds().Dy()})
	}
	return processSprites(sprites, cfg)
}

func processSprites(sprites []model.Sprite, cfg model.Config) ([]model.Sprite, error) {
	processed := make([]model.Sprite, 0, len(sprites))
	for _, s := range sprites {
		if s.Name == "" {
			s.Name = fmt.Sprintf("sprite_%d_%d", s.X, s.Y)
		}
		trimmed, err := trim.TrimSprite(s)
		if err != nil {
			return nil, err
		}
		pivoted, err := pivot.ApplyPivot(trimmed, cfg)
		if err != nil {
			return nil, err
		}
		processed = append(processed, pivoted)
	}
	return processed, nil
}
