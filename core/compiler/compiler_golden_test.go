package compiler

import (
	"encoding/json"
	"image"
	"image/color"
	"path/filepath"
	"testing"

	"pixelc/internal/imageutil"
	"pixelc/internal/testutil"
	"pixelc/pkg/model"
)

type placementJSON struct {
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	W    int    `json:"w"`
	H    int    `json:"h"`
}

func TestCompileGoldenHashes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "golden_sheet.png")
	img := image.NewRGBA(image.Rect(0, 0, 14, 10))
	img.SetRGBA(1, 1, color.RGBA{R: 255, A: 255})
	img.SetRGBA(1, 2, color.RGBA{R: 255, A: 255})
	img.SetRGBA(5, 1, color.RGBA{G: 255, A: 255})
	img.SetRGBA(6, 1, color.RGBA{G: 255, A: 255})
	img.SetRGBA(10, 7, color.RGBA{B: 255, A: 255})
	img.SetRGBA(10, 8, color.RGBA{B: 255, A: 255})
	if err := imageutil.SavePNG(path, img); err != nil {
		t.Fatalf("save input: %v", err)
	}

	cfg := model.Config{Connectivity: 4, Padding: 1, PivotMode: "bottom-center", Preset: "unity", PowerOfTwo: true}
	atlas, atlasImg, _, err := Compile(path, cfg)
	if err != nil {
		t.Fatalf("compile failed: %v", err)
	}

	placements := make([]placementJSON, 0, len(atlas.Sprites))
	for _, ps := range atlas.Sprites {
		placements = append(placements, placementJSON{Name: ps.Sprite.Name, X: ps.AtlasX, Y: ps.AtlasY, W: ps.Sprite.Width, H: ps.Sprite.Height})
	}
	pj, err := json.Marshal(placements)
	if err != nil {
		t.Fatalf("marshal placements: %v", err)
	}
	canon, err := testutil.CanonicalJSON(pj)
	if err != nil {
		t.Fatalf("canonical json: %v", err)
	}

	atlasHash := imageutil.HashRGBA(atlasImg)
	placementsHash := testutil.HashBytes(canon)

	const expectedAtlasHash = "209eb34ca78c65eba1c8af9551265483c2522b6a6b3a4aecadca4617dd91669a"
	const expectedPlacementsHash = "89bfbdaee7daf3b8b1117b9c15981857d0d63b4b8c505c1f84f63c0128501abe"

	if atlasHash != expectedAtlasHash {
		t.Fatalf("atlas hash mismatch got=%s want=%s", atlasHash, expectedAtlasHash)
	}
	if placementsHash != expectedPlacementsHash {
		t.Fatalf("placements hash mismatch got=%s want=%s", placementsHash, expectedPlacementsHash)
	}
}
