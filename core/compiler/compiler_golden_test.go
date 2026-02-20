package compiler

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"os"
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
	atlas, atlasImg, presetJSON, err := Compile(path, cfg)
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
	presetHash := testutil.HashBytes(presetJSON)

	const expectedAtlasHash = "209eb34ca78c65eba1c8af9551265483c2522b6a6b3a4aecadca4617dd91669a"
	const expectedPlacementsHash = "89bfbdaee7daf3b8b1117b9c15981857d0d63b4b8c505c1f84f63c0128501abe"
	const expectedPresetHash = "c6294dac2e74620ace7dbffeabab16f585fbe1481e171825d015de742cf5ebc8"
	if atlasHash != expectedAtlasHash {
		t.Fatalf("atlas hash mismatch got=%s want=%s", atlasHash, expectedAtlasHash)
	}
	if placementsHash != expectedPlacementsHash {
		t.Fatalf("placements hash mismatch got=%s want=%s", placementsHash, expectedPlacementsHash)
	}
	if presetHash != expectedPresetHash {
		t.Fatalf("preset hash mismatch got=%s want=%s", presetHash, expectedPresetHash)
	}
}

func TestBatchGoldenHashes(t *testing.T) {
	root := t.TempDir()
	mkGoldenPNG(t, filepath.Join(root, "a", "hero_idle_001.png"), color.RGBA{R: 255, A: 255})
	mkGoldenPNG(t, filepath.Join(root, "a", "hero_idle_002.png"), color.RGBA{G: 255, A: 255})
	mkGoldenPNG(t, filepath.Join(root, "b", "enemy_run_001.png"), color.RGBA{B: 255, A: 255})
	cfg := model.Config{Connectivity: 4, Padding: 1, PivotMode: "center", Preset: "unity", FPS: 12}
	res, err := CompileBatch(root, cfg, BatchOptions{OutDir: filepath.Join(t.TempDir(), "out")})
	if err != nil {
		t.Fatalf("batch failed: %v", err)
	}
	stable := make([]map[string]string, 0, len(res.Units))
	for _, u := range res.Units {
		stable = append(stable, map[string]string{
			"unit":  u.UnitName,
			"atlas": testutil.HashBytes([]byte(fmt.Sprintf("%dx%d-%d", u.Atlas.Width, u.Atlas.Height, len(u.Atlas.Sprites)))),
			"json":  testutil.HashBytes(u.JSON),
		})
	}
	b, _ := json.Marshal(stable)
	h := testutil.HashBytes(b)
	const expected = "51250aa477be310cf5cab8cfb428cae84c9003c88442ea7faf6f0bd4a9d19017"
	if h != expected {
		t.Fatalf("hash mismatch got=%s want=%s", h, expected)
	}
}

func TestBatchIgnoreGoldenDiffers(t *testing.T) {
	root := t.TempDir()
	mkGoldenPNG(t, filepath.Join(root, "a", "hero_idle_001.png"), color.RGBA{R: 255, A: 255})
	mkGoldenPNG(t, filepath.Join(root, "skip", "x_001.png"), color.RGBA{B: 255, A: 255})
	cfg := model.Config{Connectivity: 4, Padding: 1, PivotMode: "center", Preset: "unity", FPS: 12}
	all, err := CompileBatch(root, cfg, BatchOptions{OutDir: filepath.Join(t.TempDir(), "all")})
	if err != nil {
		t.Fatal(err)
	}
	flt, err := CompileBatch(root, cfg, BatchOptions{OutDir: filepath.Join(t.TempDir(), "flt"), IgnorePatterns: []string{"**/skip/**"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(all.Units) == len(flt.Units) {
		t.Fatalf("ignore pattern did not change unit count")
	}
}

func mkGoldenPNG(t *testing.T, p string, c color.RGBA) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	img := image.NewRGBA(image.Rect(0, 0, 4, 4))
	img.SetRGBA(1, 1, c)
	img.SetRGBA(1, 2, c)
	if err := imageutil.SavePNG(p, img); err != nil {
		t.Fatal(err)
	}
}
