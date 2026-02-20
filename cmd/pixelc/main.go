package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"pixelc/core/compiler"
	"pixelc/pkg/model"
)

type cliConfigFile struct {
	Connectivity int      `json:"connectivity"`
	Padding      int      `json:"padding"`
	PivotMode    string   `json:"pivotMode"`
	PowerOfTwo   bool     `json:"powerOfTwo"`
	Preset       string   `json:"preset"`
	FPS          int      `json:"fps"`
	Ignore       []string `json:"ignore"`
}

type stringList []string

func (s *stringList) String() string { return fmt.Sprintf("%v", []string(*s)) }
func (s *stringList) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printHelp(stdout)
		return 0
	}
	if args[0] != "compile" {
		fmt.Fprintf(stderr, "unknown command: %s\n", args[0])
		printHelp(stderr)
		return 1
	}
	return runCompile(args[1:], stdout, stderr)
}

func runCompile(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "compile requires exactly one input path")
		return 1
	}
	inputPath := args[0]
	args = args[1:]

	cfgFilePath := detectConfigPath(args)
	fileCfg := cliConfigFile{Connectivity: 4, PivotMode: "center", Preset: "unity", FPS: 12}
	if cfgFilePath != "" {
		var err error
		fileCfg, err = loadCLIConfig(cfgFilePath)
		if err != nil {
			fmt.Fprintf(stderr, "config load error: %v\n", err)
			return 1
		}
	}

	fs := flag.NewFlagSet("compile", flag.ContinueOnError)
	fs.SetOutput(stderr)
	outDir := fs.String("out", "", "output directory")
	preset := fs.String("preset", fileCfg.Preset, "output preset")
	padding := fs.Int("padding", fileCfg.Padding, "atlas padding")
	connectivity := fs.Int("connectivity", fileCfg.Connectivity, "pixel connectivity (4 or 8)")
	pivot := fs.String("pivot", fileCfg.PivotMode, "pivot mode")
	power2 := fs.Bool("power2", fileCfg.PowerOfTwo, "power-of-two atlas dimensions")
	fps := fs.Int("fps", fileCfg.FPS, "animation fps")
	batch := fs.Bool("batch", false, "batch compile recursive directories")
	dryRun := fs.Bool("dry-run", false, "plan outputs without writing files")
	report := fs.Bool("report", false, "write report.json")
	configPath := fs.String("config", "", "config file path")
	ignores := stringList{}
	ignores = append(ignores, fileCfg.Ignore...)
	fs.Var(&ignores, "ignore", "ignore glob pattern (repeatable)")
	if err := fs.Parse(args); err != nil {
		return 1
	}
	_ = configPath
	if fs.NArg() != 0 {
		fmt.Fprintln(stderr, "unexpected extra compile arguments")
		return 1
	}
	if *outDir == "" {
		fmt.Fprintln(stderr, "--out is required")
		return 1
	}

	cfg := model.Config{Connectivity: *connectivity, Padding: *padding, PivotMode: *pivot, PowerOfTwo: *power2, Preset: *preset, FPS: *fps}
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(stderr, "config validation error: %v\n", err)
		return 1
	}

	if *batch {
		return runBatchCompile(inputPath, *outDir, cfg, []string(ignores), *dryRun, *report, stdout, stderr)
	}
	return runSingleCompile(inputPath, *outDir, cfg, *dryRun, *report, stdout, stderr)
}

func runSingleCompile(inputPath, outDir string, cfg model.Config, dryRun, writeReport bool, stdout, stderr io.Writer) int {
	atlas, atlasImg, presetJSON, err := compiler.Compile(inputPath, cfg)
	if err != nil {
		fmt.Fprintf(stderr, "compile failed: %v\n", err)
		return 1
	}
	if dryRun {
		fmt.Fprintf(stdout, "dry-run sprites=%d atlas=%dx%d out=%s\n", len(atlas.Sprites), atlas.Width, atlas.Height, outDir)
		return 0
	}
	if err := compiler.WriteOutputs(outDir, atlasImg, presetJSON); err != nil {
		fmt.Fprintf(stderr, "compile failed: %v\n", err)
		return 1
	}
	if writeReport {
		reportData, err := compiler.WriteSingleReport(outDir, filepath.Base(inputPath), *atlas, atlasImg, presetJSON)
		if err != nil {
			fmt.Fprintf(stderr, "compile failed: %v\n", err)
			return 1
		}
		_ = reportData
	}
	fmt.Fprintf(stdout, "compiled sprites=%d atlas=%dx%d wrote=atlas.png,atlas.json\n", len(atlas.Sprites), atlas.Width, atlas.Height)
	return 0
}

func runBatchCompile(inputPath, outDir string, cfg model.Config, ignores []string, dryRun, writeReport bool, stdout, stderr io.Writer) int {
	res, err := compiler.CompileBatch(inputPath, cfg, compiler.BatchOptions{OutDir: outDir, IgnorePatterns: ignores, DryRun: dryRun, WriteReport: writeReport})
	if err != nil {
		fmt.Fprintf(stderr, "batch compile failed: %v\n", err)
		return 1
	}
	fmt.Fprintf(stdout, "batch units=%d out=%s dry_run=%v\n", len(res.Units), outDir, dryRun)
	for _, u := range res.Units {
		fmt.Fprintf(stdout, "unit=%s sprites=%d atlas=%dx%d\n", u.UnitName, len(u.Atlas.Sprites), u.Atlas.Width, u.Atlas.Height)
	}
	return 0
}

func detectConfigPath(args []string) string {
	for i := 0; i < len(args); i++ {
		if args[i] == "--config" && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}

func loadCLIConfig(path string) (cliConfigFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return cliConfigFile{}, err
	}
	cfg := cliConfigFile{Connectivity: 4, PivotMode: "center", Preset: "unity", FPS: 12}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cliConfigFile{}, err
	}
	if cfg.Preset == "" {
		cfg.Preset = "unity"
	}
	if cfg.PivotMode == "" {
		cfg.PivotMode = "center"
	}
	if cfg.Connectivity == 0 {
		cfg.Connectivity = 4
	}
	if cfg.FPS == 0 {
		cfg.FPS = 12
	}
	return cfg, nil
}

func printHelp(w io.Writer) {
	fmt.Fprintln(w, "pixelc compile <input> --out <dir> --preset unity --padding 2 --connectivity 4 --pivot bottom-center --power2 [--batch] [--dry-run] [--config cfg.json] [--ignore pattern] [--fps 12] [--report]")
}
