package main

import (
	"fmt"

	"github.com/fabriqaai/llm-cli/internal/cmd"
)

// Version is set by goreleaser at build time
var Version = "dev"

// Commit is set by goreleaser at build time
var Commit = "none"

// Date is set by goreleaser at build time
var Date = "unknown"

func main() {
	cmd.Execute(Version, Commit, Date)
}

// PrintVersion prints version information
func PrintVersion() {
	if Version == "dev" {
		fmt.Printf("llm-cli version %s\n", Version)
	} else {
		fmt.Printf("llm-cli version %s (commit: %s, built: %s)\n", Version, Commit, Date)
	}
}
