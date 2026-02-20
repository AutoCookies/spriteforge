package anim

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"pixelc/pkg/model"
)

var framePattern = regexp.MustCompile(`^(?i)(.*?)[\s_-]+([0-9]{1,6})$`)

type Parsed struct {
	Name      string
	State     string
	Index     int
	Grouped   bool
	SpriteRef string
}

func ParseFrameName(filename string) Parsed {
	base := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	m := framePattern.FindStringSubmatch(base)
	if len(m) != 3 {
		return Parsed{Name: base, SpriteRef: base, Grouped: false}
	}
	prefix := strings.TrimSpace(m[1])
	index := parseInt(m[2])
	tokens := splitTokens(prefix)
	if len(tokens) == 0 {
		return Parsed{Name: base, SpriteRef: base, Grouped: false}
	}
	state := tokens[len(tokens)-1]
	return Parsed{Name: base, State: state, Index: index, Grouped: true, SpriteRef: base}
}

func BuildAnimations(spriteNames []string, fps int) ([]model.Animation, []string, error) {
	if fps <= 0 {
		fps = 12
	}
	groups := map[string][]Parsed{}
	ungrouped := make([]string, 0)
	for _, name := range spriteNames {
		p := ParseFrameName(name)
		if !p.Grouped {
			ungrouped = append(ungrouped, p.SpriteRef)
			continue
		}
		groups[p.State] = append(groups[p.State], p)
	}
	states := make([]string, 0, len(groups))
	for s := range groups {
		states = append(states, s)
	}
	sort.Strings(states)

	anims := make([]model.Animation, 0, len(states))
	for _, state := range states {
		frames := groups[state]
		sort.SliceStable(frames, func(i, j int) bool {
			if frames[i].Index != frames[j].Index {
				return frames[i].Index < frames[j].Index
			}
			return frames[i].SpriteRef < frames[j].SpriteRef
		})
		for i := 1; i < len(frames); i++ {
			if frames[i-1].Index == frames[i].Index {
				return nil, nil, fmt.Errorf("duplicate frame index for state %s: %d", state, frames[i].Index)
			}
		}
		names := make([]string, 0, len(frames))
		for _, f := range frames {
			names = append(names, f.SpriteRef)
		}
		a := model.Animation{State: state, Frames: names, FPS: fps}
		if err := a.Validate(); err != nil {
			return nil, nil, err
		}
		anims = append(anims, a)
	}
	sort.Strings(ungrouped)
	return anims, ungrouped, nil
}

func splitTokens(s string) []string {
	parts := regexp.MustCompile(`[\s_-]+`).Split(strings.TrimSpace(s), -1)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, strings.ToLower(p))
		}
	}
	return out
}

func parseInt(s string) int {
	n := 0
	for _, ch := range s {
		n = n*10 + int(ch-'0')
	}
	return n
}
