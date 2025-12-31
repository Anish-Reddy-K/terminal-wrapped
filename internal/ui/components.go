package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Color palette - vibrant and eye-catching
var (
	// Primary colors
	ColorPrimary   = lipgloss.Color("#FF6B6B") // Coral red
	ColorSecondary = lipgloss.Color("#4ECDC4") // Teal
	ColorAccent    = lipgloss.Color("#FFE66D") // Yellow
	ColorPurple    = lipgloss.Color("#A855F7") // Purple
	ColorBlue      = lipgloss.Color("#3B82F6") // Blue
	ColorGreen     = lipgloss.Color("#10B981") // Green
	ColorOrange    = lipgloss.Color("#F97316") // Orange
	ColorPink      = lipgloss.Color("#EC4899") // Pink

	// Neutral colors
	ColorDim    = lipgloss.Color("#6B7280") // Gray
	ColorMuted  = lipgloss.Color("#9CA3AF") // Light gray
	ColorBright = lipgloss.Color("#F9FAFB") // White

	// Heatmap colors (low to high intensity) - more contrast
	HeatmapColors = []lipgloss.Color{
		lipgloss.Color("#2d2d2d"), // Empty/very low
		lipgloss.Color("#0e4429"), // Low green
		lipgloss.Color("#006d32"), // Medium-low
		lipgloss.Color("#26a641"), // Medium
		lipgloss.Color("#39d353"), // Medium-high
		lipgloss.Color("#4ae168"), // High
		lipgloss.Color("#73e87c"), // Very high
		lipgloss.Color("#a6f5a6"), // Max
	}
)

// Styles
var (
	// Box styles
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorDim).
			Padding(0, 1)

	HighlightBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary).
				Padding(0, 1)

	// Text styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorBright)

	LabelStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	ValueStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorBright)

	AccentStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent)

	SubtleStyle = lipgloss.NewStyle().
			Foreground(ColorDim)

	// Category colors
	CategoryColors = map[string]lipgloss.Color{
		"Git":        ColorOrange,
		"Containers": ColorBlue,
		"Packages":   ColorGreen,
		"Editors":    ColorPurple,
		"Navigation": ColorSecondary,
		"Search":     ColorPink,
		"Network":    ColorAccent,
		"Files":      ColorPrimary,
	}
)

// ProgressBar creates a horizontal progress bar
func ProgressBar(value, max int, width int, color lipgloss.Color) string {
	if max == 0 {
		max = 1
	}
	filled := int(float64(value) / float64(max) * float64(width))
	if filled > width {
		filled = width
	}

	filledStyle := lipgloss.NewStyle().Foreground(color)
	emptyStyle := lipgloss.NewStyle().Foreground(ColorDim)

	bar := filledStyle.Render(strings.Repeat("█", filled))
	bar += emptyStyle.Render(strings.Repeat("░", width-filled))
	return bar
}

// MiniBar creates a compact 4-block progress bar for categories
func MiniBar(pct float64, color lipgloss.Color) string {
	blocks := int(pct / 10) // 10% per block, max 10 blocks but we cap at 4 for display
	if blocks > 4 {
		blocks = 4
	}
	if blocks < 0 {
		blocks = 0
	}

	filledStyle := lipgloss.NewStyle().Foreground(color)
	emptyStyle := lipgloss.NewStyle().Foreground(ColorDim)

	return filledStyle.Render(strings.Repeat("█", blocks)) + emptyStyle.Render(strings.Repeat("░", 4-blocks))
}

// Heatmap creates a 7x8 heatmap (7 days, 8 time blocks of 3 hours each)
func Heatmap(data [7][24]int) string {
	var sb strings.Builder

	// Aggregate into 3-hour blocks first and find max
	blocks := [7][8]int{}
	maxVal := 1

	for day := 0; day < 7; day++ {
		for block := 0; block < 8; block++ {
			sum := 0
			for h := 0; h < 3; h++ {
				hour := block*3 + h
				sum += data[day][hour]
			}
			blocks[day][block] = sum
			if sum > maxVal {
				maxVal = sum
			}
		}
	}

	// Header
	header := "    0  3  6  9  12 15 18 21"
	sb.WriteString(SubtleStyle.Render(header))
	sb.WriteString("\n")

	days := []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}

	for day := 0; day < 7; day++ {
		sb.WriteString(SubtleStyle.Render(fmt.Sprintf(" %s ", days[day])))

		for block := 0; block < 8; block++ {
			sum := blocks[day][block]

			// Normalize and get color (0 = empty, higher = more intense)
			var colorIdx int
			if sum == 0 {
				colorIdx = 0
			} else {
				intensity := float64(sum) / float64(maxVal)
				colorIdx = 1 + int(intensity*float64(len(HeatmapColors)-2))
				if colorIdx >= len(HeatmapColors) {
					colorIdx = len(HeatmapColors) - 1
				}
			}

			blockStyle := lipgloss.NewStyle().Foreground(HeatmapColors[colorIdx])
			sb.WriteString(blockStyle.Render("██ "))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// FormatNumber formats a number with commas
func FormatNumber(n int) string {
	if n < 0 {
		return "-" + FormatNumber(-n)
	}
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}

	str := fmt.Sprintf("%d", n)
	var result strings.Builder

	for i, ch := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(ch)
	}

	return result.String()
}

// TruncateString truncates a string to max length with ellipsis
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// CenterText centers text within a given width
func CenterText(s string, width int) string {
	if len(s) >= width {
		return s
	}
	padding := (width - len(s)) / 2
	return strings.Repeat(" ", padding) + s + strings.Repeat(" ", width-len(s)-padding)
}

// RightAlign right-aligns text within a given width
func RightAlign(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return strings.Repeat(" ", width-len(s)) + s
}

// JoinHorizontal joins multiple strings horizontally with a gap
func JoinHorizontal(gap int, blocks ...string) string {
	if len(blocks) == 0 {
		return ""
	}

	// Split each block into lines
	blockLines := make([][]string, len(blocks))
	maxLines := 0

	for i, block := range blocks {
		blockLines[i] = strings.Split(block, "\n")
		if len(blockLines[i]) > maxLines {
			maxLines = len(blockLines[i])
		}
	}

	// Find width of each block
	blockWidths := make([]int, len(blocks))
	for i, lines := range blockLines {
		for _, line := range lines {
			lineLen := lipgloss.Width(line)
			if lineLen > blockWidths[i] {
				blockWidths[i] = lineLen
			}
		}
	}

	// Build output
	var result strings.Builder
	gapStr := strings.Repeat(" ", gap)

	for lineIdx := 0; lineIdx < maxLines; lineIdx++ {
		for blockIdx, lines := range blockLines {
			var line string
			if lineIdx < len(lines) {
				line = lines[lineIdx]
			}
			// Pad to block width
			lineLen := lipgloss.Width(line)
			if lineLen < blockWidths[blockIdx] {
				line += strings.Repeat(" ", blockWidths[blockIdx]-lineLen)
			}
			result.WriteString(line)
			if blockIdx < len(blockLines)-1 {
				result.WriteString(gapStr)
			}
		}
		if lineIdx < maxLines-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// SectionHeader creates a styled section header
func SectionHeader(title string, width int) string {
	style := lipgloss.NewStyle().
		Foreground(ColorSecondary).
		Bold(true)

	line := strings.Repeat("─", width-len(title)-4)
	return style.Render("─ " + title + " " + line)
}

