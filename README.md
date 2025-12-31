# Terminal Wrapped 🎁

Your personal "Spotify Wrapped" for the command line. Parses your shell history to generate a beautiful, screenshot-ready developer stats report card.

![Terminal Wrapped](https://raw.githubusercontent.com/anishreddy/terminal-wrapped/main/screenshot.png)

## Quick Start

Run instantly with a single command - no installation required:

```bash
curl -fsSL arkr.ca/wrapped | bash
```

Or if you prefer to build from source:

```bash
go install github.com/anishreddy/terminal-wrapped@latest
terminal-wrapped
```

## Features

- **Personality Archetypes** - Get assigned a fun developer archetype based on your habits
  - The Git Gladiator ⚔️
  - The Sudo Summoner ⚡
  - The Docker Captain 🐳
  - The Vim Wizard 🧙
  - The Night Owl 🦉
  - And more...

- **Rich Statistics**
  - Total commands & unique commands
  - Peak coding hours heatmap
  - Command category breakdown (Git, Docker, Packages, etc.)
  - Longest streak of active days
  - Busiest coding day
  - Favorite directories & editors

- **Screenshot Ready** - Designed specifically for sharing on X/Twitter

## Supported Shells

- **zsh** (with or without extended history timestamps)
- **bash**

Works with any terminal: iTerm2, Terminal.app, Warp, Ghostty, Alacritty, Kitty, etc.

## Usage

```bash
terminal-wrapped [flags]

Flags:
  --shell string    Force shell type: zsh, bash (auto-detected if not set)
  --history string  Custom history file path
  --no-color        Disable colors (for piping)
  --json            Output raw stats as JSON
  --version         Show version information
  --help            Show help message
```

### Examples

```bash
# Auto-detect everything
terminal-wrapped

# Force bash history
terminal-wrapped --shell bash

# Custom history file
terminal-wrapped --history ~/.histfile

# Output as JSON for programmatic use
terminal-wrapped --json
```

## Increase Your History

To get more accurate stats, increase your shell history limit:

### For zsh (add to ~/.zshrc):
```bash
HISTSIZE=50000
SAVEHIST=50000
setopt EXTENDED_HISTORY      # Save timestamps
setopt SHARE_HISTORY         # Share between sessions
```

### For bash (add to ~/.bashrc):
```bash
HISTSIZE=50000
HISTFILESIZE=50000
```

## Building from Source

```bash
git clone https://github.com/anishreddy/terminal-wrapped.git
cd terminal-wrapped
make build
./terminal-wrapped
```

### Cross-compile for all platforms:

```bash
make release
# Creates binaries in dist/ for:
# - darwin-amd64 (Intel Mac)
# - darwin-arm64 (Apple Silicon)
# - linux-amd64
# - linux-arm64
```

## How It Works

1. Detects your shell and history file location
2. Parses command history (handles both timestamped and plain formats)
3. Analyzes patterns: command frequency, timing, categories
4. Assigns a personality archetype based on your habits
5. Renders a beautiful TUI output optimized for screenshots

## Privacy

Terminal Wrapped runs 100% locally. Your history data never leaves your machine.

## Tech Stack

- **Go** - Fast, single binary, no runtime dependencies
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Beautiful terminal styling
- **[Termenv](https://github.com/muesli/termenv)** - Terminal capability detection

## Contributing

Contributions welcome! Ideas for new archetypes, stats, or improvements are appreciated.

## License

MIT License - See [LICENSE](LICENSE) for details.

---

**Share your wrapped!** Screenshot and post on X with **#TerminalWrapped**
