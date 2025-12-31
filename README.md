# Terminal Wrapped

Your "Spotify Wrapped" for the command line. See your developer stats in a beautiful, screenshot-ready report.

## Quick Start

Run instantly - no installation required:

```bash
curl -fsSL arkr.ca/terminal-wrapped | bash
```

That's it. It auto-detects your shell and history file.

<img width="510" height="671" alt="terminal-wrapped" src="https://github.com/user-attachments/assets/d7340901-ca6f-45b6-a113-98f054702990" />

## What You Get

- **Your Developer Archetype** - Are you a Git Gladiator? Night Owl? Sudo Summoner?
- **Top Commands** - Your most-used commands with visual bars
- **Activity Heatmap** - When do you code most?
- **Quick Stats** - Streaks, sudo usage, pipe complexity
- **Category Breakdown** - Git, Docker, Packages, and more

## Supported

- **Shells**: zsh, bash
- **OS**: macOS, Linux
- **Terminals**: All (iTerm2, Terminal.app, Warp, Ghostty, Kitty, etc.)

## Save More History

To get better stats, increase your history limit:

```bash
echo 'HISTSIZE=100000' >> ~/.zshrc && exec zsh
```
This will increase your history capacity by 50x while using only ~3 MB of extra disk space (less than a single song).

## Building from Source

```bash
git clone https://github.com/Anish-Reddy-K/terminal-wrapped.git
cd terminal-wrapped
go build -o terminal-wrapped .
./terminal-wrapped
```

## Privacy

Runs 100% locally. Your data never leaves your machine.

---

**Share your stats!** Screenshot and post on X with **#TerminalWrapped**

*by Anish Reddy ([arkr.ca](https://arkr.ca))*
