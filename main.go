package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/Anish-Reddy-K/terminal-wrapped/internal/analyzer"
	"github.com/Anish-Reddy-K/terminal-wrapped/internal/parser"
	"github.com/Anish-Reddy-K/terminal-wrapped/internal/ui"
	"github.com/muesli/termenv"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// CLI flags
	shellFlag := flag.String("shell", "", "Force shell type: zsh, bash (auto-detected if not set)")
	historyFlag := flag.String("history", "", "Custom history file path")
	noColorFlag := flag.Bool("no-color", false, "Disable colors")
	jsonFlag := flag.Bool("json", false, "Output raw stats as JSON")
	versionFlag := flag.Bool("version", false, "Show version")
	helpFlag := flag.Bool("help", false, "Show help")

	flag.Parse()

	if *helpFlag {
		printHelp()
		return
	}

	if *versionFlag {
		fmt.Printf("terminal-wrapped %s (%s) built on %s\n", version, commit, date)
		fmt.Printf("Go version: %s\n", runtime.Version())
		return
	}

	// Configure terminal output
	output := termenv.NewOutput(os.Stdout)
	if *noColorFlag {
		output = termenv.NewOutput(os.Stdout, termenv.WithProfile(termenv.Ascii))
	}
	_ = output // Used by lipgloss internally

	// Detect or use specified shell
	shell := *shellFlag
	if shell == "" {
		shell = parser.DetectShell()
	}

	// Get history file path
	historyPath := *historyFlag
	if historyPath == "" {
		historyPath = parser.GetHistoryPath(shell)
	}

	// Check if history file exists
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: History file not found: %s\n", historyPath)
		fmt.Fprintf(os.Stderr, "Try specifying the path with --history or shell with --shell\n")
		os.Exit(1)
	}

	// Parse history
	startTime := time.Now()
	data, err := parser.Parse(historyPath, shell)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing history: %v\n", err)
		os.Exit(1)
	}
	parseTime := time.Since(startTime)

	if len(data.Commands) == 0 {
		fmt.Fprintf(os.Stderr, "No commands found in history file: %s\n", historyPath)
		os.Exit(1)
	}

	// Analyze
	stats := analyzer.Analyze(data)

	// Detect archetype
	archetype := analyzer.DetectArchetype(stats)

	// Output
	if *jsonFlag {
		outputJSON(stats, archetype, parseTime)
	} else {
		fmt.Print(ui.Render(stats, archetype))
	}
}

func printHelp() {
	fmt.Println(`
Terminal Wrapped - Your Developer Stats Report Card

Usage:
  terminal-wrapped [flags]

Flags:
  --shell string    Force shell type: zsh, bash (auto-detected if not set)
  --history string  Custom history file path
  --no-color        Disable colors (for piping)
  --json            Output raw stats as JSON
  --version         Show version information
  --help            Show this help message

Examples:
  terminal-wrapped                    # Auto-detect and display
  terminal-wrapped --shell bash       # Force bash history
  terminal-wrapped --history ~/.histfile  # Custom history file
  terminal-wrapped --json             # Output as JSON

For more info: https://github.com/anishreddy/terminal-wrapped
`)
}

type jsonOutput struct {
	Stats     *analyzer.Stats     `json:"stats"`
	Archetype *analyzer.Archetype `json:"archetype"`
	Meta      struct {
		ParseTimeMs int64  `json:"parse_time_ms"`
		Version     string `json:"version"`
	} `json:"meta"`
}

func outputJSON(stats *analyzer.Stats, archetype *analyzer.Archetype, parseTime time.Duration) {
	out := jsonOutput{
		Stats:     stats,
		Archetype: archetype,
	}
	out.Meta.ParseTimeMs = parseTime.Milliseconds()
	out.Meta.Version = version

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(out)
}

