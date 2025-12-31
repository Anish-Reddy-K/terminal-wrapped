package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Command represents a single command from history
type Command struct {
	Raw       string
	Command   string    // First word (the actual command)
	Args      []string  // Arguments
	Timestamp time.Time // Zero value if no timestamp available
	HasTime   bool      // Whether we have timestamp data
}

// HistoryData contains all parsed history information
type HistoryData struct {
	Commands  []Command
	Shell     string // "zsh" or "bash"
	FilePath  string
	HasTimes  bool // Whether timestamps are available
	ParsedAt  time.Time
	LineCount int
}

// zsh extended history format: `: 1703961234:0;command`
var zshExtendedRegex = regexp.MustCompile(`^:\s*(\d+):\d+;(.*)$`)

// DetectShell tries to detect the current shell
func DetectShell() string {
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "zsh") {
		return "zsh"
	}
	if strings.Contains(shell, "bash") {
		return "bash"
	}
	// Default to zsh as it's more common on modern systems
	return "zsh"
}

// GetHistoryPath returns the default history file path for a shell
func GetHistoryPath(shell string) string {
	home, _ := os.UserHomeDir()
	
	switch shell {
	case "zsh":
		// Check HISTFILE env first
		if histFile := os.Getenv("HISTFILE"); histFile != "" {
			return histFile
		}
		return filepath.Join(home, ".zsh_history")
	case "bash":
		// Check HISTFILE env first
		if histFile := os.Getenv("HISTFILE"); histFile != "" {
			return histFile
		}
		return filepath.Join(home, ".bash_history")
	default:
		return filepath.Join(home, ".zsh_history")
	}
}

// Parse reads and parses a shell history file
func Parse(historyPath string, shell string) (*HistoryData, error) {
	file, err := os.Open(historyPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := &HistoryData{
		Commands: make([]Command, 0, 10000),
		Shell:    shell,
		FilePath: historyPath,
		ParsedAt: time.Now(),
	}

	scanner := bufio.NewScanner(file)
	// Increase buffer size for very long commands
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var multilineCmd strings.Builder
	inMultiline := false

	for scanner.Scan() {
		line := scanner.Text()
		data.LineCount++

		// Handle multiline commands (ending with \)
		if inMultiline {
			multilineCmd.WriteString("\n")
			multilineCmd.WriteString(line)
			if !strings.HasSuffix(line, "\\") {
				inMultiline = false
				cmd := parseCommand(multilineCmd.String(), shell)
				if cmd != nil {
					if cmd.HasTime {
						data.HasTimes = true
					}
					data.Commands = append(data.Commands, *cmd)
				}
				multilineCmd.Reset()
			}
			continue
		}

		if strings.HasSuffix(line, "\\") {
			inMultiline = true
			multilineCmd.WriteString(line)
			continue
		}

		cmd := parseCommand(line, shell)
		if cmd != nil {
			if cmd.HasTime {
				data.HasTimes = true
			}
			data.Commands = append(data.Commands, *cmd)
		}
	}

	// Handle any remaining multiline command
	if multilineCmd.Len() > 0 {
		cmd := parseCommand(multilineCmd.String(), shell)
		if cmd != nil {
			data.Commands = append(data.Commands, *cmd)
		}
	}

	return data, scanner.Err()
}

func parseCommand(line string, shell string) *Command {
	if len(line) == 0 {
		return nil
	}

	cmd := &Command{
		Raw: line,
	}

	// Try to parse ZSH extended history format
	if shell == "zsh" {
		if matches := zshExtendedRegex.FindStringSubmatch(line); matches != nil {
			timestamp, err := strconv.ParseInt(matches[1], 10, 64)
			if err == nil {
				cmd.Timestamp = time.Unix(timestamp, 0)
				cmd.HasTime = true
			}
			line = matches[2]
			cmd.Raw = line
		}
	}

	// Skip empty commands after parsing
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return nil
	}

	// Parse the command and arguments
	parts := parseCommandParts(line)
	if len(parts) == 0 {
		return nil
	}

	cmd.Command = parts[0]
	if len(parts) > 1 {
		cmd.Args = parts[1:]
	}

	return cmd
}

// parseCommandParts splits a command line into parts, respecting quotes
func parseCommandParts(line string) []string {
	var parts []string
	var current strings.Builder
	inQuote := false
	quoteChar := rune(0)

	for _, r := range line {
		switch {
		case r == '"' || r == '\'':
			if inQuote && r == quoteChar {
				inQuote = false
			} else if !inQuote {
				inQuote = true
				quoteChar = r
			}
			current.WriteRune(r)
		case r == ' ' || r == '\t':
			if inQuote {
				current.WriteRune(r)
			} else if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

// GetBaseCommand extracts the base command, handling sudo, env vars, etc.
func GetBaseCommand(cmd *Command) string {
	command := cmd.Command

	// Skip leading env var assignments (VAR=value cmd)
	if strings.Contains(command, "=") && len(cmd.Args) > 0 {
		command = cmd.Args[0]
		// Could be multiple env vars, keep going
		for i := 1; i < len(cmd.Args); i++ {
			if !strings.Contains(command, "=") {
				break
			}
			command = cmd.Args[i]
		}
	}

	// Handle sudo/doas
	if command == "sudo" || command == "doas" {
		if len(cmd.Args) > 0 {
			// Skip flags like -u, -i, etc.
			for _, arg := range cmd.Args {
				if !strings.HasPrefix(arg, "-") && !strings.Contains(arg, "=") {
					return arg
				}
			}
		}
	}

	// Handle common wrappers
	wrappers := []string{"time", "nice", "nohup", "strace", "ltrace"}
	for _, wrapper := range wrappers {
		if command == wrapper && len(cmd.Args) > 0 {
			return cmd.Args[0]
		}
	}

	return command
}

