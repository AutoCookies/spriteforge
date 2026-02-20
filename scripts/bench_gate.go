//go:build benchgate

package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var allocLimit = map[string]float64{
	"BenchmarkSlicer_100Sprites": 2000,
	"BenchmarkPacker_500Sprites": 2500,
}

const maxAtlasPx = 4_000_000

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run -tags benchgate ./scripts/bench_gate.go <bench_output.txt>")
		os.Exit(1)
	}
	b, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	text := string(b)
	fail := false
	for name, limit := range allocLimit {
		alloc, ok := findMetric(text, name, `([0-9]+)\s+allocs/op`)
		if !ok {
			fmt.Fprintf(os.Stderr, "missing benchmark metric for %s\n", name)
			fail = true
			continue
		}
		if alloc > limit {
			fmt.Fprintf(os.Stderr, "%s allocs/op %.0f exceeds limit %.0f\n", name, alloc, limit)
			fail = true
		}
	}
	atlasPx, ok := findGlobalMetric(text, `PACKER_BENCH_ATLAS_PX=([0-9]+)`)
	if !ok {
		fmt.Fprintln(os.Stderr, "missing PACKER_BENCH_ATLAS_PX marker")
		fail = true
	} else if atlasPx > maxAtlasPx {
		fmt.Fprintf(os.Stderr, "atlas size metric too large %.0f > %d\n", atlasPx, maxAtlasPx)
		fail = true
	}
	if fail {
		os.Exit(1)
	}
	fmt.Println("benchmark gate passed")
}

func findMetric(text, benchName, pattern string) (float64, bool) {
	re := regexp.MustCompile(benchName + `[\s\S]*?` + pattern)
	m := re.FindStringSubmatch(text)
	if len(m) < 2 {
		return 0, false
	}
	v, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

func findGlobalMetric(text, pattern string) (float64, bool) {
	re := regexp.MustCompile(pattern)
	m := re.FindStringSubmatch(text)
	if len(m) < 2 {
		return 0, false
	}
	v, err := strconv.ParseFloat(m[1], 64)
	if err != nil {
		return 0, false
	}
	return v, true
}
