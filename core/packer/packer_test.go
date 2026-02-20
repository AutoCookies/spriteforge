package packer

import (
	"fmt"
	"image"
	"image/color"
	"testing"

	"pixelc/pkg/model"
)

func TestPackBasicAndConstraints(t *testing.T) {
	cfg := model.Config{Connectivity: 4, Padding: 1, PivotMode: "center", Preset: "unity"}
	sprites := []model.Sprite{makeSprite("b", 3, 2, color.RGBA{R: 255, A: 255}), makeSprite("a", 2, 3, color.RGBA{G: 255, A: 255})}

	atlas, img, err := Pack(sprites, cfg)
	if err != nil {
		t.Fatalf("pack failed: %v", err)
	}
	if atlas.Width <= 0 || atlas.Height <= 0 || img == nil {
		t.Fatalf("invalid atlas output")
	}
	if len(atlas.Sprites) != 2 {
		t.Fatalf("expected 2 sprites")
	}
	assertNoOverlap(t, atlas, cfg.Padding)
	for _, ps := range atlas.Sprites {
		if ps.AtlasX < cfg.Padding || ps.AtlasY < cfg.Padding {
			t.Fatalf("edge padding not respected: %+v", ps)
		}
	}
}

func TestPackDeterministic(t *testing.T) {
	cfg := model.Config{Connectivity: 4, Padding: 2, PivotMode: "center", Preset: "unity"}
	sprites := []model.Sprite{
		makeSprite("s1", 4, 3, color.RGBA{R: 255, A: 255}),
		makeSprite("s2", 3, 4, color.RGBA{G: 255, A: 255}),
		makeSprite("s3", 2, 2, color.RGBA{B: 255, A: 255}),
	}
	a1, _, err := Pack(sprites, cfg)
	if err != nil {
		t.Fatalf("pack1 failed: %v", err)
	}
	a2, _, err := Pack(sprites, cfg)
	if err != nil {
		t.Fatalf("pack2 failed: %v", err)
	}
	if a1.Width != a2.Width || a1.Height != a2.Height || len(a1.Sprites) != len(a2.Sprites) {
		t.Fatalf("atlas mismatch")
	}
	for i := range a1.Sprites {
		if a1.Sprites[i].AtlasX != a2.Sprites[i].AtlasX || a1.Sprites[i].AtlasY != a2.Sprites[i].AtlasY {
			t.Fatalf("placement mismatch at %d", i)
		}
	}
}

func TestPackPowerOfTwoAndEdgeCases(t *testing.T) {
	cfg := model.Config{Connectivity: 4, Padding: 0, PivotMode: "center", Preset: "unity", PowerOfTwo: true}
	a, _, err := Pack([]model.Sprite{makeSprite("a", 13, 7, color.RGBA{A: 255})}, cfg)
	if err != nil {
		t.Fatalf("pack failed: %v", err)
	}
	if !isPow2(a.Width) || !isPow2(a.Height) {
		t.Fatalf("dimensions not power-of-two: %dx%d", a.Width, a.Height)
	}

	empty, img, err := Pack(nil, cfg)
	if err != nil || empty.Width != 0 || empty.Height != 0 || img != nil {
		t.Fatalf("empty pack invalid")
	}

	_, _, err = Pack([]model.Sprite{{Name: "bad", Width: 0, Height: 1}}, cfg)
	if err == nil {
		t.Fatalf("expected error for invalid sprite")
	}
}

func TestPackRendering(t *testing.T) {
	cfg := model.Config{Connectivity: 4, Padding: 1, PivotMode: "center", Preset: "unity"}
	sprites := []model.Sprite{makeSprite("red", 2, 2, color.RGBA{R: 255, A: 255}), makeSprite("green", 2, 2, color.RGBA{G: 255, A: 255})}
	atlas, img, err := Pack(sprites, cfg)
	if err != nil {
		t.Fatalf("pack failed: %v", err)
	}
	for _, ps := range atlas.Sprites {
		c := img.RGBAAt(ps.AtlasX, ps.AtlasY)
		sc := ps.Sprite.Image.RGBAAt(0, 0)
		if c != sc {
			t.Fatalf("render mismatch")
		}
	}
	if img.RGBAAt(0, 0).A != 0 {
		t.Fatalf("expected transparent background")
	}
}

func BenchmarkPacker_500Sprites(b *testing.B) {
	cfg := model.Config{Connectivity: 4, Padding: 1, PivotMode: "center", Preset: "unity", PowerOfTwo: true}
	sprites := make([]model.Sprite, 0, 500)
	for i := 0; i < 500; i++ {
		w := 4 + (i % 5)
		h := 4 + ((i / 5) % 5)
		sprites = append(sprites, makeSprite("s"+string(rune(i%26+'a')), w, h, color.RGBA{R: uint8(i), A: 255}))
	}
	one, _, err := Pack(sprites, cfg)
	if err != nil {
		b.Fatalf("pack failed: %v", err)
	}
	fmt.Printf("PACKER_BENCH_ATLAS_PX=%d\n", one.Width*one.Height)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, _, err := Pack(sprites, cfg); err != nil {
			b.Fatalf("pack failed: %v", err)
		}
	}
}

func makeSprite(name string, w, h int, c color.RGBA) model.Sprite {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, c)
		}
	}
	return model.Sprite{Name: name, Image: img, Width: w, Height: h}
}

func assertNoOverlap(t *testing.T, atlas model.Atlas, padding int) {
	t.Helper()
	for i := 0; i < len(atlas.Sprites); i++ {
		a := atlas.Sprites[i]
		ar := rect{x: a.AtlasX - padding, y: a.AtlasY - padding, w: a.Sprite.Width + padding*2, h: a.Sprite.Height + padding*2}
		for j := i + 1; j < len(atlas.Sprites); j++ {
			b := atlas.Sprites[j]
			br := rect{x: b.AtlasX - padding, y: b.AtlasY - padding, w: b.Sprite.Width + padding*2, h: b.Sprite.Height + padding*2}
			if intersects(ar, br) {
				t.Fatalf("overlap between %d and %d", i, j)
			}
		}
	}
}

func intersects(a, b rect) bool {
	return a.x < b.x+b.w && a.x+a.w > b.x && a.y < b.y+b.h && a.y+a.h > b.y
}

func isPow2(v int) bool { return v > 0 && (v&(v-1)) == 0 }
