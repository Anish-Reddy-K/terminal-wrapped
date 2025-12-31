package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ASCII art for the header - compact but impactful
const asciiArt = `
▀█▀ █▀▀ █▀█ █▀▄▀█ █ █▄ █ ▄▀█ █       █ █ █ █▀█ ▄▀█ █▀█ █▀█ █▀▀ █▀▄
 █  ██▄ █▀▄ █ ▀ █ █ █ ▀█ █▀█ █▄▄     ▀▄▀▄▀ █▀▄ █▀█ █▀▀ █▀▀ ██▄ █▄▀`

// RenderHeader renders the colorful ASCII art header
func RenderHeader() string {
	lines := strings.Split(strings.TrimSpace(asciiArt), "\n")
	var result strings.Builder

	// Create gradient effect
	colors := []lipgloss.Color{
		ColorPrimary,  // Coral
		ColorOrange,   // Orange
		ColorAccent,   // Yellow
		ColorGreen,    // Green
		ColorSecondary, // Teal
	}

	// Top border
	borderStyle := lipgloss.NewStyle().Foreground(ColorDim)
	result.WriteString(borderStyle.Render("╔" + strings.Repeat("═", 78) + "╗"))
	result.WriteString("\n")

	for _, line := range lines {
		result.WriteString(borderStyle.Render("║"))
		result.WriteString(" ")

		// Apply gradient to each character
		runes := []rune(line)
		for i, r := range runes {
			colorIdx := i * len(colors) / len(runes)
			if colorIdx >= len(colors) {
				colorIdx = len(colors) - 1
			}
			charStyle := lipgloss.NewStyle().Foreground(colors[colorIdx]).Bold(true)
			result.WriteString(charStyle.Render(string(r)))
		}

		// Pad to 76 chars
		padding := 76 - len(runes)
		if padding > 0 {
			result.WriteString(strings.Repeat(" ", padding))
		}

		result.WriteString(" ")
		result.WriteString(borderStyle.Render("║"))
		result.WriteString("\n")
	}

	// Bottom border
	result.WriteString(borderStyle.Render("╚" + strings.Repeat("═", 78) + "╝"))

	return result.String()
}

// RenderFooter renders the footer with sharing info
func RenderFooter() string {
	borderStyle := lipgloss.NewStyle().Foreground(ColorDim)
	linkStyle := lipgloss.NewStyle().Foreground(ColorSecondary)
	hashtagStyle := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)

	footer := borderStyle.Render("─") + " " +
		linkStyle.Render("github.com/anishreddy/terminal-wrapped") +
		SubtleStyle.Render("  •  ") +
		SubtleStyle.Render("Share on X: ") +
		hashtagStyle.Render("#TerminalWrapped") +
		" " + borderStyle.Render("─")

	return CenterText(footer, 80)
}

// RenderHistoryTip renders the tip about increasing history
func RenderHistoryTip() string {
	style := BoxStyle.Copy().
		BorderForeground(ColorDim).
		Width(78)

	tipStyle := lipgloss.NewStyle().Foreground(ColorAccent)
	codeStyle := lipgloss.NewStyle().Foreground(ColorGreen)

	content := tipStyle.Render("💡 Tip: ") +
		LabelStyle.Render("To store more history, add to your .zshrc: ") +
		codeStyle.Render("HISTSIZE=50000 SAVEHIST=50000")

	return style.Render(content)
}

