package exporter

import (
	"encoding/json"
	"fmt"
	"sort"

	"pixelc/internal/version"
	"pixelc/pkg/model"
	"pixelc/pkg/schema"
)

func ExportUnity(atlas model.Atlas, atlasImageName string, appVersion string) ([]byte, error) {
	if err := atlas.Validate(); err != nil {
		return nil, err
	}
	if atlasImageName == "" {
		return nil, fmt.Errorf("atlas image name is required")
	}
	if appVersion == "" {
		appVersion = version.Version
	}

	out := schema.UnityAtlasJSON{Frames: map[string]schema.UnityFrame{}}
	out.Meta.App = version.AppName
	out.Meta.Version = appVersion
	out.Meta.Image = atlasImageName
	out.Meta.Size.W = atlas.Width
	out.Meta.Size.H = atlas.Height

	ordered := make([]model.PlacedSprite, len(atlas.Sprites))
	copy(ordered, atlas.Sprites)
	sort.SliceStable(ordered, func(i, j int) bool {
		return ordered[i].Sprite.Name < ordered[j].Sprite.Name
	})

	for _, ps := range ordered {
		f := schema.UnityFrame{}
		f.Frame.X = ps.AtlasX
		f.Frame.Y = ps.AtlasY
		f.Frame.W = ps.Sprite.Width
		f.Frame.H = ps.Sprite.Height
		f.Pivot.X = ps.Sprite.PivotX
		f.Pivot.Y = ps.Sprite.PivotY
		out.Frames[ps.Sprite.Name] = f
	}

	b, err := json.Marshal(out)
	if err != nil {
		return nil, fmt.Errorf("marshal unity json: %w", err)
	}
	return b, nil
}
