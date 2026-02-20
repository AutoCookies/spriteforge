package model

import "image"

type Sprite struct {
	Name   string
	Image  *image.RGBA
	X      int
	Y      int
	Width  int
	Height int
	PivotX float64
	PivotY float64
}

type PlacedSprite struct {
	Sprite Sprite
	AtlasX int
	AtlasY int
}

type Atlas struct {
	Width   int
	Height  int
	Sprites []PlacedSprite
}

type Animation struct {
	State  string
	Frames []string
	FPS    int
}

type Config struct {
	Connectivity int    // 4 or 8
	Padding      int    // >=0
	PivotMode    string // "center" | "bottom-center"
	PowerOfTwo   bool
	Preset       string // "unity" (v1), others later
	FPS          int    // >0 defaults to 12 when zero
}
