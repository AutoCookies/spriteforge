package packer

import (
	"fmt"
	"image"
	"sort"

	"pixelc/internal/imageutil"
	"pixelc/pkg/model"
)

type rect struct {
	x int
	y int
	w int
	h int
}

type sortableSprite struct {
	idx int
	s   model.Sprite
}

func Pack(sprites []model.Sprite, cfg model.Config) (model.Atlas, *image.RGBA, error) {
	if err := cfg.Validate(); err != nil {
		return model.Atlas{}, nil, err
	}
	if len(sprites) == 0 {
		return model.Atlas{Width: 0, Height: 0, Sprites: nil}, nil, nil
	}
	for i, s := range sprites {
		if s.Width <= 0 || s.Height <= 0 || s.Image == nil {
			return model.Atlas{}, nil, fmt.Errorf("invalid sprite at index %d", i)
		}
	}

	sorted := make([]sortableSprite, len(sprites))
	for i := range sprites {
		sorted[i] = sortableSprite{idx: i, s: sprites[i]}
	}
	sort.SliceStable(sorted, func(i, j int) bool {
		a, b := sorted[i].s, sorted[j].s
		if a.Height != b.Height {
			return a.Height > b.Height
		}
		if a.Width != b.Width {
			return a.Width > b.Width
		}
		return a.Name < b.Name
	})

	w, h := initialDimensions(sorted, cfg.Padding)
	if cfg.PowerOfTwo {
		w = nextPowerOfTwo(w)
		h = nextPowerOfTwo(h)
	}

	for {
		atlas, img, ok, err := tryPack(sorted, sprites, cfg, w, h)
		if err != nil {
			return model.Atlas{}, nil, err
		}
		if ok {
			return atlas, img, nil
		}
		if cfg.PowerOfTwo {
			if w <= h {
				w *= 2
			} else {
				h *= 2
			}
		} else {
			if w <= h {
				w += max(1, w/4)
			} else {
				h += max(1, h/4)
			}
		}
	}
}

func tryPack(sorted []sortableSprite, original []model.Sprite, cfg model.Config, atlasW, atlasH int) (model.Atlas, *image.RGBA, bool, error) {
	free := []rect{{x: 0, y: 0, w: atlasW, h: atlasH}}
	placedByOriginal := make([]model.PlacedSprite, len(original))

	for _, item := range sorted {
		paddedW := item.s.Width + cfg.Padding*2
		paddedH := item.s.Height + cfg.Padding*2

		bestIdx, bestNode := bestFreeRect(free, paddedW, paddedH)
		if bestIdx < 0 {
			return model.Atlas{}, nil, false, nil
		}

		free = splitFreeRects(free, bestIdx, bestNode)
		free = pruneFreeRects(free)

		placedByOriginal[item.idx] = model.PlacedSprite{
			Sprite: item.s,
			AtlasX: bestNode.x + cfg.Padding,
			AtlasY: bestNode.y + cfg.Padding,
		}
	}

	atlasSprites := make([]model.PlacedSprite, len(original))
	copy(atlasSprites, placedByOriginal)
	atlas := model.Atlas{Width: atlasW, Height: atlasH, Sprites: atlasSprites}
	atlasImg := image.NewRGBA(image.Rect(0, 0, atlasW, atlasH))
	for _, ps := range atlas.Sprites {
		if err := imageutil.Blit(atlasImg, ps.Sprite.Image, ps.AtlasX, ps.AtlasY); err != nil {
			return model.Atlas{}, nil, false, err
		}
	}
	return atlas, atlasImg, true, nil
}

// Tie-break rules: short side fit, then long side fit, then top-most (y), then left-most (x).
func bestFreeRect(free []rect, w, h int) (int, rect) {
	bestIdx := -1
	best := rect{}
	bestShort, bestLong := int(^uint(0)>>1), int(^uint(0)>>1)
	for i, fr := range free {
		if w > fr.w || h > fr.h {
			continue
		}
		leftoverH := fr.h - h
		leftoverW := fr.w - w
		short := min(leftoverW, leftoverH)
		long := max(leftoverW, leftoverH)
		cand := rect{x: fr.x, y: fr.y, w: w, h: h}
		if short < bestShort ||
			(short == bestShort && (long < bestLong ||
				(long == bestLong && (cand.y < best.y || (cand.y == best.y && cand.x < best.x))))) {
			bestIdx = i
			best = cand
			bestShort = short
			bestLong = long
		}
	}
	return bestIdx, best
}

func splitFreeRects(free []rect, usedIdx int, used rect) []rect {
	selected := free[usedIdx]
	result := make([]rect, 0, len(free)+4)
	for i, fr := range free {
		if i == usedIdx {
			continue
		}
		result = append(result, fr)
	}

	if selected.y < used.y {
		result = append(result, rect{x: selected.x, y: selected.y, w: selected.w, h: used.y - selected.y})
	}
	if used.y+used.h < selected.y+selected.h {
		result = append(result, rect{x: selected.x, y: used.y + used.h, w: selected.w, h: selected.y + selected.h - (used.y + used.h)})
	}
	if selected.x < used.x {
		result = append(result, rect{x: selected.x, y: selected.y, w: used.x - selected.x, h: selected.h})
	}
	if used.x+used.w < selected.x+selected.w {
		result = append(result, rect{x: used.x + used.w, y: selected.y, w: selected.x + selected.w - (used.x + used.w), h: selected.h})
	}
	return result
}

func pruneFreeRects(free []rect) []rect {
	out := make([]rect, 0, len(free))
	for i := range free {
		if free[i].w <= 0 || free[i].h <= 0 {
			continue
		}
		contained := false
		for j := range free {
			if i == j {
				continue
			}
			if contains(free[j], free[i]) {
				contained = true
				break
			}
		}
		if !contained {
			out = append(out, free[i])
		}
	}
	return out
}

func contains(a, b rect) bool {
	return b.x >= a.x && b.y >= a.y && b.x+b.w <= a.x+a.w && b.y+b.h <= a.y+a.h
}

func initialDimensions(sorted []sortableSprite, padding int) (int, int) {
	maxW, maxH := 1, 1
	area := 0
	for _, s := range sorted {
		w := s.s.Width + padding*2
		h := s.s.Height + padding*2
		if w > maxW {
			maxW = w
		}
		if h > maxH {
			maxH = h
		}
		area += w * h
	}
	side := 1
	for side*side < area {
		side *= 2
	}
	if side < maxW {
		side = maxW
	}
	if side < maxH {
		side = maxH
	}
	return side, side
}

func nextPowerOfTwo(v int) int {
	if v <= 1 {
		return 1
	}
	p := 1
	for p < v {
		p <<= 1
	}
	return p
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
