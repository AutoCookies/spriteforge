package compiler

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"pixelc/core/anim"
	"pixelc/internal/imageutil"
	"pixelc/internal/testutil"
	"pixelc/pkg/model"
)

type BatchOptions struct {
	IgnorePatterns []string
	DryRun         bool
	WriteReport    bool
	OutDir         string
}

type BatchResult struct {
	Units []UnitResult `json:"units"`
}

type UnitResult struct {
	UnitName string      `json:"unit_name"`
	OutDir   string      `json:"out_dir"`
	Atlas    model.Atlas `json:"atlas"`
	JSON     []byte      `json:"-"`
	Report   []byte      `json:"report,omitempty"`
}

type reportJSON struct {
	UnitName       string `json:"unit_name"`
	SpriteCount    int    `json:"sprite_count"`
	AtlasWidth     int    `json:"atlas_width"`
	AtlasHeight    int    `json:"atlas_height"`
	Animations     int    `json:"animations_count"`
	AtlasPngSHA256 string `json:"atlas_png_sha256"`
	AtlasJSONSHA   string `json:"atlas_json_sha256"`
}

func CompileBatch(inputPath string, cfg model.Config, opts BatchOptions) (*BatchResult, error) {
	units, err := discoverUnits(inputPath, opts.IgnorePatterns)
	if err != nil {
		return nil, err
	}
	result := &BatchResult{Units: make([]UnitResult, 0, len(units))}
	for _, rel := range units {
		unitPath := filepath.Join(inputPath, rel)
		atlas, atlasImg, presetJSON, err := Compile(unitPath, cfg)
		if err != nil {
			return nil, fmt.Errorf("compile unit %s: %w", rel, err)
		}
		outDir := filepath.Join(opts.OutDir, rel)
		if !opts.DryRun {
			if err := WriteOutputs(outDir, atlasImg, presetJSON); err != nil {
				return nil, err
			}
		}
		unit := UnitResult{UnitName: rel, OutDir: outDir, Atlas: *atlas, JSON: presetJSON}
		if opts.WriteReport {
			rep, err := buildUnitReport(rel, *atlas, atlasImg, presetJSON)
			if err != nil {
				return nil, err
			}
			unit.Report = rep
			if !opts.DryRun {
				if err := os.MkdirAll(outDir, 0o755); err != nil {
					return nil, err
				}
				if err := os.WriteFile(filepath.Join(outDir, "report.json"), rep, 0o644); err != nil {
					return nil, err
				}
			}
		}
		result.Units = append(result.Units, unit)
	}
	return result, nil
}

func buildUnitReport(unitName string, atlas model.Atlas, atlasImg *image.RGBA, presetJSON []byte) ([]byte, error) {
	names := make([]string, 0, len(atlas.Sprites))
	for _, s := range atlas.Sprites {
		names = append(names, s.Sprite.Name)
	}
	anims, _, err := anim.BuildAnimations(names, 12)
	if err != nil {
		return nil, err
	}
	rep := reportJSON{
		UnitName:       unitName,
		SpriteCount:    len(atlas.Sprites),
		AtlasWidth:     atlas.Width,
		AtlasHeight:    atlas.Height,
		Animations:     len(anims),
		AtlasPngSHA256: imageutil.HashRGBA(atlasImg),
		AtlasJSONSHA:   testutil.HashBytes(presetJSON),
	}
	b, err := json.Marshal(rep)
	if err != nil {
		return nil, err
	}
	return testutil.CanonicalJSON(b)
}

func discoverUnits(root string, ignore []string) ([]string, error) {
	allIgnore := append([]string{".git", "node_modules", "build", "dist"}, ignore...)
	units := make([]string, 0)
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		if rel == "." {
			rel = ""
		}
		if shouldIgnore(rel, allIgnore) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		hasPNG := false
		for _, e := range entries {
			if !e.IsDir() && strings.EqualFold(filepath.Ext(e.Name()), ".png") {
				hasPNG = true
				break
			}
		}
		if hasPNG {
			if rel == "" {
				units = append(units, ".")
			} else {
				units = append(units, filepath.ToSlash(rel))
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(units)
	return units, nil
}

func shouldIgnore(rel string, patterns []string) bool {
	rel = filepath.ToSlash(rel)
	for _, p := range patterns {
		p = filepath.ToSlash(strings.TrimSpace(p))
		if p == "" {
			continue
		}
		if globMatch(p, rel) || strings.HasPrefix(rel, strings.TrimSuffix(p, "/")+"/") || rel == p {
			return true
		}
	}
	return false
}

func globMatch(pattern, value string) bool {
	if pattern == "" {
		return false
	}
	if pattern == "*" || pattern == "**" {
		return true
	}
	if strings.Contains(pattern, "**") {
		parts := strings.Split(pattern, "**")
		idx := 0
		for _, part := range parts {
			part = strings.Trim(part, "/")
			if part == "" {
				continue
			}
			pos := strings.Index(value[idx:], part)
			if pos < 0 {
				return false
			}
			idx += pos + len(part)
		}
		return true
	}
	ok, _ := filepath.Match(pattern, value)
	return ok
}
