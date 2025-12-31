package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ASCII art for the header - clean and impactful
const asciiArt = `
 _____                   _             _  __      __                              _ 
|_   _|__ _ __ _ __ ___ (_)_ __   __ _| | \ \    / / __ __ _ _ __  _ __   ___  __| |
  | |/ _ \ '__| '_ ' _ \| | '_ \ / _' | |  \ \/\/ / '__/ _' | '_ \| '_ \ / _ \/ _' |
  | |  __/ |  | | | | | | | | | | (_| | |   \    /| | | (_| | |_) | |_) |  __/ (_| |
  |_|\___|_|  |_| |_| |_|_|_| |_|\__,_|_|    \/\/ |_|  \__,_| .__/| .__/ \___|\__,_|
                                                            |_|   |_|              `

// renderHeader renders the colorful ASCII art header
func RenderHeader() string {
	lines := strings.Split(strings.TrimPrefix(asciiArt, "\n"), "\n")
	var result strings.Builder

	// create gradient effect
	colors := []lipgloss.Color{
		ColorPrimary,   // Coral
		ColorOrange,    // Orange
		ColorAccent,    // Yellow
		ColorGreen,     // Green
		ColorSecondary, // Teal
		ColorBlue,      // Blue
	}

	for _, line := range lines {
		// apply gradient to each character
		runes := []rune(line)
		for i, r := range runes {
			if r == ' ' {
				result.WriteRune(' ')
				continue
			}
			colorIdx := i * len(colors) / max(len(runes), 1)
			if colorIdx >= len(colors) {
				colorIdx = len(colors) - 1
			}
			charStyle := lipgloss.NewStyle().Foreground(colors[colorIdx]).Bold(true)
			result.WriteString(charStyle.Render(string(r)))
		}
		result.WriteString("\n")
	}

	return strings.TrimSuffix(result.String(), "\n")
}

// renderFooter renders the footer with sharing info and branding
func RenderFooter() string {
	lineStyle := lipgloss.NewStyle().Foreground(ColorDim)
	linkStyle := lipgloss.NewStyle().Foreground(ColorSecondary)
	hashtagStyle := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)
	brandStyle := lipgloss.NewStyle().Foreground(ColorMuted)

	line := lineStyle.Render(strings.Repeat("-", 76))
	
	// line 1: GitHub link
	link := linkStyle.Render("github.com/Anish-Reddy-K/terminal-wrapped")
	line1 := CenterText(link, 76)
	
	// line 2: Share + branding
	share := SubtleStyle.Render("Share on X with ") + hashtagStyle.Render("#TerminalWrapped")
	brand := brandStyle.Render("by Anish Reddy (arkr.ca)")
	line2 := "  " + share + strings.Repeat(" ", 76-lipgloss.Width(share)-lipgloss.Width(brand)-4) + brand

	return line + "\n" + line1 + "\n" + line2
}

// RenderHistoryTip renders the tip about increasing history
func RenderHistoryTip() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorDim).
		Padding(0, 1).
		Width(76)

	tipIcon := AccentStyle.Render("[i]")
	tipText := LabelStyle.Render(" Save more history: ")
	
	code := lipgloss.NewStyle().Foreground(ColorGreen).Render(
		"echo 'HISTSIZE=100000' >> ~/.zshrc && exec zsh")

	return style.Render(tipIcon + tipText + code)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
