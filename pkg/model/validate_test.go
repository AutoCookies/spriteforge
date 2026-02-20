package model

import "testing"

func TestConfigValidate(t *testing.T) {
	valid := Config{Connectivity: 4, Padding: 0, PivotMode: "center", Preset: "unity", FPS: 12}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid config, got %v", err)
	}

	cases := []Config{
		{Connectivity: 5, Padding: 0, PivotMode: "center", Preset: "unity"},
		{Connectivity: 4, Padding: -1, PivotMode: "center", Preset: "unity"},
		{Connectivity: 4, Padding: 0, PivotMode: "top", Preset: "unity"},
		{Connectivity: 4, Padding: 0, PivotMode: "center", Preset: "invalid"},
		{Connectivity: 4, Padding: 0, PivotMode: "center", Preset: "unity", FPS: -1},
	}

	for _, cfg := range cases {
		if err := cfg.Validate(); err == nil {
			t.Fatalf("expected config validation failure for %+v", cfg)
		}
	}
}

func TestAtlasValidate(t *testing.T) {
	valid := Atlas{Width: 32, Height: 32, Sprites: []PlacedSprite{{Sprite: Sprite{Width: 16, Height: 16}, AtlasX: 0, AtlasY: 0}}}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid atlas, got %v", err)
	}

	cases := []Atlas{
		{Width: -1, Height: 10},
		{Width: 10, Height: -1},
		{Width: 10, Height: 10, Sprites: []PlacedSprite{{Sprite: Sprite{Width: -1, Height: 1}, AtlasX: 0, AtlasY: 0}}},
		{Width: 10, Height: 10, Sprites: []PlacedSprite{{Sprite: Sprite{Width: 1, Height: 1, X: -1}, AtlasX: 0, AtlasY: 0}}},
		{Width: 10, Height: 10, Sprites: []PlacedSprite{{Sprite: Sprite{Width: 1, Height: 1}, AtlasX: -1, AtlasY: 0}}},
	}

	for _, atlas := range cases {
		if err := atlas.Validate(); err == nil {
			t.Fatalf("expected atlas validation failure for %+v", atlas)
		}
	}
}

func TestAnimationValidate(t *testing.T) {
	a := Animation{State: "idle", Frames: []string{"f1", "f2"}, FPS: 12}
	if err := a.Validate(); err != nil {
		t.Fatalf("expected valid animation, got %v", err)
	}
	bad := []Animation{{State: "", Frames: []string{"f"}, FPS: 12}, {State: "idle", Frames: []string{""}, FPS: 12}, {State: "idle", Frames: []string{"f"}, FPS: 0}}
	for _, x := range bad {
		if err := x.Validate(); err == nil {
			t.Fatalf("expected invalid animation: %+v", x)
		}
	}
}
