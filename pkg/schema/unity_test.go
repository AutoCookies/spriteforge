package schema

import (
	"encoding/json"
	"testing"
)

func TestUnityAtlasJSONMarshalDeterministic(t *testing.T) {
	atlas := UnityAtlasJSON{Frames: map[string]UnityFrame{}, Meta: UnityMeta{App: "pixelc", Version: "0.0.0-dev", Image: "atlas.png"}}
	atlas.Meta.Size.W = 16
	atlas.Meta.Size.H = 16

	fA := UnityFrame{}
	fA.Frame.X, fA.Frame.Y, fA.Frame.W, fA.Frame.H = 1, 2, 3, 4
	fA.Pivot.X, fA.Pivot.Y = 0.5, 1.0
	atlas.Frames["a"] = fA

	fB := UnityFrame{}
	fB.Frame.X, fB.Frame.Y, fB.Frame.W, fB.Frame.H = 5, 6, 7, 8
	fB.Pivot.X, fB.Pivot.Y = 0.5, 0.5
	atlas.Frames["b"] = fB

	b1, err := json.Marshal(atlas)
	if err != nil {
		t.Fatalf("marshal 1 failed: %v", err)
	}
	b2, err := json.Marshal(atlas)
	if err != nil {
		t.Fatalf("marshal 2 failed: %v", err)
	}
	if string(b1) != string(b2) {
		t.Fatalf("expected deterministic marshal, got %s vs %s", b1, b2)
	}
}
