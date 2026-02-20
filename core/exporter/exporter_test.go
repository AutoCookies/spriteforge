package exporter

import (
	"encoding/json"
	"testing"

	"pixelc/internal/testutil"
	"pixelc/pkg/model"
	"pixelc/pkg/schema"
)

func TestExportUnityDeterministic(t *testing.T) {
	atlas := model.Atlas{Width: 32, Height: 16, Sprites: []model.PlacedSprite{
		{Sprite: model.Sprite{Name: "player_idle_02", Width: 2, Height: 3, PivotX: 0.5, PivotY: 1}, AtlasX: 10, AtlasY: 4},
		{Sprite: model.Sprite{Name: "player_idle_01", Width: 1, Height: 1, PivotX: 0.5, PivotY: 0.5}, AtlasX: 1, AtlasY: 2},
	}}
	b1, err := ExportUnity(atlas, "atlas.png", "0.1.0", 12)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}
	b2, err := ExportUnity(atlas, "atlas.png", "0.1.0", 12)
	if err != nil {
		t.Fatalf("export failed: %v", err)
	}
	if string(b1) != string(b2) {
		t.Fatalf("non-deterministic output")
	}

	var out schema.UnityAtlasJSON
	if err := json.Unmarshal(b1, &out); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if out.Meta.Image != "atlas.png" || out.Meta.Size.W != 32 || out.Meta.Size.H != 16 {
		t.Fatalf("unexpected meta: %+v", out.Meta)
	}
	if _, ok := out.Frames["player_idle_01"]; !ok {
		t.Fatalf("missing frame")
	}
	if out.Animations["idle"].FPS != 12 || len(out.Animations["idle"].Frames) != 2 {
		t.Fatalf("missing animations")
	}

	canon, err := testutil.CanonicalJSON(b1)
	if err != nil {
		t.Fatalf("canonicalize failed: %v", err)
	}
	if len(canon) == 0 {
		t.Fatalf("empty canonical json")
	}
}

func TestExportUnityValidation(t *testing.T) {
	_, err := ExportUnity(model.Atlas{Width: -1}, "atlas.png", "0.1.0", 12)
	if err == nil {
		t.Fatalf("expected validation error")
	}
}
