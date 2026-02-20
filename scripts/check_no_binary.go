package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"
)

func main() {
	cmd := exec.Command("git", "ls-files")
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "list tracked files: %v\n", err)
		os.Exit(1)
	}

	files := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, f := range files {
		if f == "" {
			continue
		}
		data, err := os.ReadFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read file %s: %v\n", f, err)
			os.Exit(1)
		}
		if isBinary(data) {
			fmt.Fprintf(os.Stderr, "binary file detected: %s\n", f)
			os.Exit(1)
		}
	}
	fmt.Println("no binary files detected")
}

func isBinary(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	if bytes.IndexByte(data, 0) >= 0 {
		return true
	}
	if !utf8.Valid(data) {
		return true
	}
	return false
}
