package model

import "testing"

func TestConfigValidate(t *testing.T) {
	valid := Config{Connectivity: 4, Padding: 0, PivotMode: "center", Preset: "unity"}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid config, got %v", err)
	}

	cases := []Config{
		{Connectivity: 5, Padding: 0, PivotMode: "center", Preset: "unity"},
		{Connectivity: 4, Padding: -1, PivotMode: "center", Preset: "unity"},
		{Connectivity: 4, Padding: 0, PivotMode: "top", Preset: "unity"},
		{Connectivity: 4, Padding: 0, PivotMode: "center", Preset: "json"},
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
