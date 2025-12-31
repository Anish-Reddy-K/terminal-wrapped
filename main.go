package main

import (
	"fmt"
	"os"

	"github.com/Anish-Reddy-K/terminal-wrapped/internal/analyzer"
	"github.com/Anish-Reddy-K/terminal-wrapped/internal/parser"
	"github.com/Anish-Reddy-K/terminal-wrapped/internal/ui"
)

func main() {
	// auto-detect shell
	shell := parser.DetectShell()

	// get history file path
	historyPath := parser.GetHistoryPath(shell)

	// check if history file exists
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: History file not found: %s\n", historyPath)
		fmt.Fprintf(os.Stderr, "Make sure you're using zsh or bash.\n")
		os.Exit(1)
	}

	// parse history
	data, err := parser.Parse(historyPath, shell)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing history: %v\n", err)
		os.Exit(1)
	}

	if len(data.Commands) == 0 {
		fmt.Fprintf(os.Stderr, "No commands found in history file: %s\n", historyPath)
		os.Exit(1)
	}

	// analyze
	stats := analyzer.Analyze(data)

	// detect archetype
	archetype := analyzer.DetectArchetype(stats)

	// render output
	fmt.Print(ui.Render(stats, archetype))
}
