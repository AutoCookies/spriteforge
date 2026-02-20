package version

import (
	"strings"
	"testing"
)

func TestFullVersion(t *testing.T) {
	oldV, oldC, oldD := Version, Commit, BuildDate
	Version, Commit, BuildDate = "1.0.0", "abc123", "2026-02-20T00:00:00Z"
	t.Cleanup(func() { Version, Commit, BuildDate = oldV, oldC, oldD })
	out := FullVersion()
	for _, x := range []string{"pixelc", "1.0.0", "abc123", "2026-02-20T00:00:00Z"} {
		if !strings.Contains(out, x) {
			t.Fatalf("missing %q in %s", x, out)
		}
	}
}
