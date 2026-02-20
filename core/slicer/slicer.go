package slicer

import (
	"fmt"
	"image"
	"sort"

	"pixelc/pkg/model"
)

type point struct {
	x int
	y int
}

type component struct {
	minX  int
	minY  int
	maxX  int
	maxY  int
	count int
	px    []point
}

func SliceSpritesheet(img *image.RGBA, cfg model.Config) ([]model.Sprite, error) {
	if img == nil {
		return nil, fmt.Errorf("image is nil")
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w == 0 || h == 0 {
		return []model.Sprite{}, nil
	}

	visited := make([]bool, w*h)
	components := make([]component, 0)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if visited[idx(x, y, w)] || !opaqueAt(img, bounds.Min.X+x, bounds.Min.Y+y) {
				continue
			}
			c := bfsComponent(img, bounds, x, y, visited, cfg.Connectivity)
			if c.count >= 2 {
				components = append(components, c)
			}
		}
	}

	sort.Slice(components, func(i, j int) bool {
		if components[i].minY != components[j].minY {
			return components[i].minY < components[j].minY
		}
		return components[i].minX < components[j].minX
	})

	sprites := make([]model.Sprite, 0, len(components))
	for i, c := range components {
		s := componentToSprite(img, bounds, c, i)
		sprites = append(sprites, s)
	}
	return sprites, nil
}

func bfsComponent(img *image.RGBA, bounds image.Rectangle, sx, sy int, visited []bool, connectivity int) component {
	w := bounds.Dx()
	queue := make([]point, 0, 32)
	queue = append(queue, point{x: sx, y: sy})
	visited[idx(sx, sy, w)] = true

	c := component{minX: sx, minY: sy, maxX: sx, maxY: sy}
	for head := 0; head < len(queue); head++ {
		p := queue[head]
		c.count++
		c.px = append(c.px, p)
		if p.x < c.minX {
			c.minX = p.x
		}
		if p.y < c.minY {
			c.minY = p.y
		}
		if p.x > c.maxX {
			c.maxX = p.x
		}
		if p.y > c.maxY {
			c.maxY = p.y
		}

		for _, n := range neighbors(p, connectivity) {
			if n.x < 0 || n.y < 0 || n.x >= bounds.Dx() || n.y >= bounds.Dy() {
				continue
			}
			nidx := idx(n.x, n.y, w)
			if visited[nidx] {
				continue
			}
			if !opaqueAt(img, bounds.Min.X+n.x, bounds.Min.Y+n.y) {
				continue
			}
			visited[nidx] = true
			queue = append(queue, n)
		}
	}
	return c
}

func componentToSprite(img *image.RGBA, bounds image.Rectangle, c component, i int) model.Sprite {
	width := c.maxX - c.minX + 1
	height := c.maxY - c.minY + 1
	spriteImg := image.NewRGBA(image.Rect(0, 0, width, height))
	for _, p := range c.px {
		spriteImg.SetRGBA(p.x-c.minX, p.y-c.minY, img.RGBAAt(bounds.Min.X+p.x, bounds.Min.Y+p.y))
	}
	return model.Sprite{
		Name:   fmt.Sprintf("sprite_%04d", i),
		Image:  spriteImg,
		X:      c.minX,
		Y:      c.minY,
		Width:  width,
		Height: height,
	}
}

func neighbors(p point, connectivity int) []point {
	if connectivity == 8 {
		return []point{{p.x + 1, p.y}, {p.x - 1, p.y}, {p.x, p.y + 1}, {p.x, p.y - 1}, {p.x + 1, p.y + 1}, {p.x + 1, p.y - 1}, {p.x - 1, p.y + 1}, {p.x - 1, p.y - 1}}
	}
	return []point{{p.x + 1, p.y}, {p.x - 1, p.y}, {p.x, p.y + 1}, {p.x, p.y - 1}}
}

func opaqueAt(img *image.RGBA, x, y int) bool {
	return img.RGBAAt(x, y).A > 0
}

func idx(x, y, w int) int {
	return y*w + x
}
