package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/anishreddy/terminal-wrapped/internal/analyzer"
)

// Render produces the complete terminal output
func Render(stats *analyzer.Stats, archetype *analyzer.Archetype) string {
	var sb strings.Builder

	// Header
	sb.WriteString("\n")
	sb.WriteString(RenderHeader())
	sb.WriteString("\n\n")

	// Hero section: Total commands + Archetype (side by side)
	heroLeft := renderHeroStats(stats)
	heroRight := renderArchetype(archetype)
	sb.WriteString("  ")
	sb.WriteString(JoinHorizontal(2, heroLeft, heroRight))
	sb.WriteString("\n\n")

	// Quick stats row
	sb.WriteString(renderQuickStats(stats))
	sb.WriteString("\n\n")

	// Middle section: Top commands + Category mix + Heatmap
	topCmds := renderTopCommands(stats)
	rightPanel := renderRightPanel(stats)
	sb.WriteString("  ")
	sb.WriteString(JoinHorizontal(2, topCmds, rightPanel))
	sb.WriteString("\n\n")

	// Fun facts row
	sb.WriteString(renderFunFacts(stats))
	sb.WriteString("\n\n")

	// History tip
	sb.WriteString("  ")
	sb.WriteString(RenderHistoryTip())
	sb.WriteString("\n\n")

	// Footer
	sb.WriteString(RenderFooter())
	sb.WriteString("\n")

	return sb.String()
}

func renderHeroStats(stats *analyzer.Stats) string {
	style := BoxStyle.Copy().
		Width(37).
		Height(6).
		BorderForeground(ColorPrimary)

	bigNumber := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorAccent)

	var content strings.Builder

	// Big number
	content.WriteString(bigNumber.Render(fmt.Sprintf("  ██  %s COMMANDS", FormatNumber(stats.TotalCommands))))
	content.WriteString("\n")
	content.WriteString(SubtleStyle.Render("      " + strings.Repeat("━", 28)))
	content.WriteString("\n")

	// Time span
	if stats.HasTimeData && !stats.FirstCommand.IsZero() {
		span := fmt.Sprintf("  %s → %s (%s)",
			stats.FirstCommand.Format("Jan 2006"),
			stats.LastCommand.Format("Jan 2006"),
			analyzer.FormatDuration(stats.HistorySpan))
		content.WriteString(LabelStyle.Render(span))
		content.WriteString("\n")

		perDay := fmt.Sprintf("  ~%.0f commands/day", stats.CommandsPerDay)
		content.WriteString(LabelStyle.Render(perDay))
	} else {
		content.WriteString(LabelStyle.Render("  All-time history"))
		content.WriteString("\n")
		content.WriteString(SubtleStyle.Render("  (no timestamps available)"))
	}

	return style.Render(content.String())
}

func renderArchetype(arch *analyzer.Archetype) string {
	style := HighlightBoxStyle.Copy().
		Width(37).
		Height(6).
		BorderForeground(ColorPurple)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorBright)

	emojiStyle := lipgloss.NewStyle().
		Foreground(ColorAccent)

	taglineStyle := lipgloss.NewStyle().
		Foreground(ColorMuted).
		Italic(true)

	var content strings.Builder
	content.WriteString(LabelStyle.Render("  YOUR ARCHETYPE"))
	content.WriteString("\n\n")
	content.WriteString(emojiStyle.Render("   " + arch.Emoji + "  "))
	content.WriteString(titleStyle.Render(arch.Name))
	content.WriteString("\n")
	content.WriteString(taglineStyle.Render("   \"" + arch.Tagline + "\""))

	return style.Render(content.String())
}

func renderQuickStats(stats *analyzer.Stats) string {
	style := BoxStyle.Copy().
		Width(76)

	// Create stat items
	items := []struct {
		label string
		value string
	}{
		{"Unique Cmds", FormatNumber(stats.UniqueCommands)},
		{"Longest Streak", fmt.Sprintf("%d days", stats.LongestStreak)},
		{"Busiest Day", formatBusiestDay(stats)},
		{"sudo Level", fmt.Sprintf("%d %s", stats.SudoCount, sudoMeter(stats.SudoPct))},
		{"Pipes Used", FormatNumber(stats.PipeCount)},
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(ColorSecondary).
		Bold(true)

	var header strings.Builder
	header.WriteString(headerStyle.Render("─ QUICK STATS "))
	header.WriteString(SubtleStyle.Render(strings.Repeat("─", 60)))

	var values strings.Builder
	for i, item := range items {
		labelStyle := LabelStyle
		valueStyle := ValueStyle

		values.WriteString("  ")
		values.WriteString(labelStyle.Render(item.label))
		values.WriteString("\n  ")
		values.WriteString(valueStyle.Render(item.value))

		if i < len(items)-1 {
			values.WriteString("         ")
		}
	}

	// Arrange horizontally
	cols := make([]string, len(items))
	for i, item := range items {
		cols[i] = fmt.Sprintf("%s\n%s",
			LabelStyle.Render(item.label),
			ValueStyle.Render(item.value))
	}

	content := "  " + JoinHorizontal(4, cols...)
	return style.Render(header.String() + "\n" + content)
}

func formatBusiestDay(stats *analyzer.Stats) string {
	if stats.BusiestDay.IsZero() {
		return "N/A"
	}
	return fmt.Sprintf("%s (%d)", stats.BusiestDay.Format("Jan 2"), stats.BusiestDayCount)
}

func sudoMeter(pct float64) string {
	blocks := int(pct / 5)
	if blocks > 5 {
		blocks = 5
	}
	return strings.Repeat("█", blocks) + strings.Repeat("░", 5-blocks)
}

func renderTopCommands(stats *analyzer.Stats) string {
	style := BoxStyle.Copy().
		Width(38)

	headerStyle := lipgloss.NewStyle().
		Foreground(ColorSecondary).
		Bold(true)

	var content strings.Builder
	content.WriteString(headerStyle.Render("─ TOP COMMANDS "))
	content.WriteString(SubtleStyle.Render(strings.Repeat("─", 21)))
	content.WriteString("\n")

	maxCount := 0
	if len(stats.TopCommands) > 0 {
		maxCount = stats.TopCommands[0].Count
	}

	// Show top 8
	for i, cmd := range stats.TopCommands {
		if i >= 8 {
			break
		}

		numStyle := SubtleStyle
		cmdStyle := ValueStyle
		countStyle := LabelStyle

		bar := ProgressBar(cmd.Count, maxCount, 18, getCmdColor(cmd.Command))
		content.WriteString(fmt.Sprintf(" %s %-7s %s %s\n",
			numStyle.Render(fmt.Sprintf("%d.", i+1)),
			cmdStyle.Render(TruncateString(cmd.Command, 7)),
			bar,
			countStyle.Render(FormatNumber(cmd.Count))))
	}

	return style.Render(content.String())
}

func getCmdColor(cmd string) lipgloss.Color {
	// Assign colors based on category
	switch cmd {
	case "git", "gh", "hub":
		return ColorOrange
	case "docker", "kubectl", "podman":
		return ColorBlue
	case "npm", "yarn", "pip", "cargo", "brew":
		return ColorGreen
	case "vim", "nvim", "code":
		return ColorPurple
	case "ssh", "curl", "wget":
		return ColorAccent
	case "grep", "find", "rg":
		return ColorPink
	default:
		return ColorSecondary
	}
}

func renderRightPanel(stats *analyzer.Stats) string {
	// Category mix + Heatmap stacked
	catMix := renderCategoryMix(stats)
	heatmap := renderHeatmapSection(stats)

	return catMix + "\n" + heatmap
}

func renderCategoryMix(stats *analyzer.Stats) string {
	style := BoxStyle.Copy().
		Width(35)

	headerStyle := lipgloss.NewStyle().
		Foreground(ColorSecondary).
		Bold(true)

	var content strings.Builder
	content.WriteString(headerStyle.Render("─ CATEGORY MIX "))
	content.WriteString(SubtleStyle.Render(strings.Repeat("─", 18)))
	content.WriteString("\n")

	// Sort categories by usage
	type catItem struct {
		name string
		pct  float64
	}
	var cats []catItem
	for name, pct := range stats.CategoryPct {
		cats = append(cats, catItem{name, pct})
	}
	sort.Slice(cats, func(i, j int) bool {
		return cats[i].pct > cats[j].pct
	})

	// Display in 2 columns, up to 8 categories
	row1 := make([]string, 0, 4)
	row2 := make([]string, 0, 4)

	for i, cat := range cats {
		if i >= 8 {
			break
		}
		color := CategoryColors[cat.name]
		if color == "" {
			color = ColorMuted
		}
		item := fmt.Sprintf("%s %s %.0f%%",
			MiniBar(cat.pct, color),
			lipgloss.NewStyle().Foreground(color).Render(TruncateString(cat.name, 6)),
			cat.pct)

		if i%2 == 0 {
			row1 = append(row1, item)
		} else {
			row2 = append(row2, item)
		}
	}

	for i := 0; i < len(row1); i++ {
		content.WriteString(" ")
		content.WriteString(row1[i])
		if i < len(row2) {
			content.WriteString("  ")
			content.WriteString(row2[i])
		}
		content.WriteString("\n")
	}

	return style.Render(content.String())
}

func renderHeatmapSection(stats *analyzer.Stats) string {
	style := BoxStyle.Copy().
		Width(35)

	headerStyle := lipgloss.NewStyle().
		Foreground(ColorSecondary).
		Bold(true)

	var content strings.Builder
	content.WriteString(headerStyle.Render("─ PEAK HOURS "))
	content.WriteString(SubtleStyle.Render(strings.Repeat("─", 20)))
	content.WriteString("\n")

	if stats.HasTimeData {
		content.WriteString(Heatmap(stats.HeatMap))

		// Peak callout
		peakStyle := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)
		peak := fmt.Sprintf(" ⚡ Peak: %s %d:00-%d:00",
			analyzer.GetDayName(stats.PeakDay),
			stats.PeakHour,
			(stats.PeakHour+1)%24)
		content.WriteString(peakStyle.Render(peak))
	} else {
		content.WriteString(SubtleStyle.Render("\n  No timestamp data available\n"))
		content.WriteString(SubtleStyle.Render("  (enable EXTENDED_HISTORY in zsh)\n"))
	}

	return style.Render(content.String())
}

func renderFunFacts(stats *analyzer.Stats) string {
	style := BoxStyle.Copy().
		Width(76)

	headerStyle := lipgloss.NewStyle().
		Foreground(ColorSecondary).
		Bold(true)

	var content strings.Builder
	content.WriteString(headerStyle.Render("─ FUN FACTS "))
	content.WriteString(SubtleStyle.Render(strings.Repeat("─", 62)))
	content.WriteString("\n")

	facts := []string{}

	// Night owl
	if stats.HasTimeData {
		facts = append(facts, fmt.Sprintf("🦉 Night Owl: %.0f%% after midnight", stats.NightOwlPct))
		facts = append(facts, fmt.Sprintf("📅 Weekend: %.0f%% on Sat/Sun", stats.WeekendPct))
	}

	// Home base
	if stats.FavoriteDir != "" {
		facts = append(facts, fmt.Sprintf("🏠 Home: %s", TruncateString(stats.FavoriteDir, 20)))
	}

	// Editor
	if stats.EditorChoice != "" {
		facts = append(facts, fmt.Sprintf("⌨️  Editor: %s (%s opens)", stats.EditorChoice, FormatNumber(stats.EditorCount)))
	}

	// Most repeated
	if stats.MostRepeated != "" && stats.MostRepeatedCount > 2 {
		facts = append(facts, fmt.Sprintf("🔁 Repeated: \"%s\" (%dx)", TruncateString(stats.MostRepeated, 15), stats.MostRepeatedCount))
	}

	// Longest command (clean up any timestamp prefixes for display)
	if stats.LongestCmdLen > 50 {
		longestDisplay := stats.LongestCommand
		// Skip display if it looks like raw zsh format
		if len(longestDisplay) > 0 && longestDisplay[0] != ':' {
			facts = append(facts, fmt.Sprintf("📏 Longest: \"%s\" (%d ch)", TruncateString(longestDisplay, 15), stats.LongestCmdLen))
		}
	}

	// Average command length
	facts = append(facts, fmt.Sprintf("📊 Avg length: %.0f chars", stats.AvgCommandLen))

	// Display in 2 columns
	for i := 0; i < len(facts); i += 2 {
		content.WriteString(" ")
		content.WriteString(LabelStyle.Render(facts[i]))
		if i+1 < len(facts) {
			// Pad first column
			padding := 38 - lipgloss.Width(facts[i])
			if padding > 0 {
				content.WriteString(strings.Repeat(" ", padding))
			}
			content.WriteString(LabelStyle.Render(facts[i+1]))
		}
		content.WriteString("\n")
	}

	return style.Render(content.String())
}

