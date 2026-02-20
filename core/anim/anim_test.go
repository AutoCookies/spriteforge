package anim

import "testing"

func TestParseFrameName(t *testing.T) {
	cases := map[string]struct {
		state   string
		idx     int
		grouped bool
	}{
		"player_idle_01.png":   {"idle", 1, true},
		"player-run-12.png":    {"run", 12, true},
		"enemy attack 003.png": {"attack", 3, true},
		"icon.png":             {"", 0, false},
	}
	for in, want := range cases {
		got := ParseFrameName(in)
		if got.Grouped != want.grouped || got.State != want.state || got.Index != want.idx {
			t.Fatalf("parse %s got %+v", in, got)
		}
	}
}

func TestBuildAnimationsOrderingAndUngrouped(t *testing.T) {
	names := []string{"player_run_02", "player_idle_02", "player_idle_01", "misc", "player_run_01"}
	anims, ungrouped, err := BuildAnimations(names, 12)
	if err != nil {
		t.Fatalf("build failed: %v", err)
	}
	if len(anims) != 2 || anims[0].State != "idle" || anims[1].State != "run" {
		t.Fatalf("unexpected states: %+v", anims)
	}
	if anims[0].Frames[0] != "player_idle_01" || anims[1].Frames[1] != "player_run_02" {
		t.Fatalf("unexpected frame order")
	}
	if len(ungrouped) != 1 || ungrouped[0] != "misc" {
		t.Fatalf("unexpected ungrouped")
	}
}

func TestBuildAnimationsDuplicateIndex(t *testing.T) {
	_, _, err := BuildAnimations([]string{"p_idle_01", "p-idle-01"}, 12)
	if err == nil {
		t.Fatalf("expected duplicate index error")
	}
}
