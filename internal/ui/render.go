package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Anish-Reddy-K/terminal-wrapped/internal/analyzer"
)

const (
	// Fixed widths for consistent layout
	TotalWidth    = 76
	LeftColWidth  = 36
	RightColWidth = 36
	ColGap        = 2
)

// Render produces the complete terminal output
func Render(stats *analyzer.Stats, archetype *analyzer.Archetype) string {
	var sb strings.Builder

	// Header
	sb.WriteString("\n")
	sb.WriteString(RenderHeader())
	sb.WriteString("\n\n")

	// Hero section: Total commands + Archetype (side by side, same height)
	heroLeft := renderHeroStats(stats)
	heroRight := renderArchetype(archetype)
	
	// Ensure same height
	heroLeftStyled := lipgloss.NewStyle().Height(8).Render(heroLeft)
	heroRightStyled := lipgloss.NewStyle().Height(8).Render(heroRight)
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, heroLeftStyled, "  ", heroRightStyled))
	sb.WriteString("\n\n")

	// Quick stats row
	sb.WriteString(renderQuickStats(stats))
	sb.WriteString("\n\n")

	// Middle section: Top commands + Right panel (aligned heights)
	topCmds := renderTopCommands(stats)
	rightPanel := renderRightPanel(stats)
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, topCmds, "  ", rightPanel))
	sb.WriteString("\n\n")

	// Fun facts row
	sb.WriteString(renderFunFacts(stats))
	sb.WriteString("\n\n")

	// History tip
	sb.WriteString(RenderHistoryTip())
	sb.WriteString("\n\n")

	// Footer
	sb.WriteString(RenderFooter())
	sb.WriteString("\n")

	return sb.String()
}

func renderHeroStats(stats *analyzer.Stats) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(0, 1).
		Width(LeftColWidth).
		Height(6)

	numberStyle := lipgloss.NewStyle().Bold(true).Foreground(ColorAccent)
	
	var lines []string
	lines = append(lines, LabelStyle.Render(" TOTAL COMMANDS"))
	lines = append(lines, "")
	lines = append(lines, numberStyle.Render(fmt.Sprintf(" [#] %s", FormatNumber(stats.TotalCommands))))
	lines = append(lines, SubtleStyle.Render("     "+strings.Repeat("-", 26)))

	if stats.HasTimeData && !stats.FirstCommand.IsZero() {
		span := fmt.Sprintf(" %s -> %s (%s)",
			stats.FirstCommand.Format("Jan 2006"),
			stats.LastCommand.Format("Jan 2006"),
			analyzer.FormatDuration(stats.HistorySpan))
		lines = append(lines, LabelStyle.Render(span))
		lines = append(lines, LabelStyle.Render(fmt.Sprintf(" ~%.0f commands/day", stats.CommandsPerDay)))
	} else {
		lines = append(lines, LabelStyle.Render(" All-time history"))
		lines = append(lines, SubtleStyle.Render(" (no timestamps)"))
	}

	return style.Render(strings.Join(lines, "\n"))
}

func renderArchetype(arch *analyzer.Archetype) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPurple).
		Padding(0, 1).
		Width(RightColWidth).
		Height(6)

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(ColorBright)
	taglineStyle := lipgloss.NewStyle().Foreground(ColorMuted).Italic(true)

	// Wrap tagline if needed (max ~28 chars per line)
	tagline := arch.Tagline
	tagline1 := tagline
	tagline2 := ""
	if len(tagline) > 28 {
		breakPoint := 28
		for i := 28; i > 15; i-- {
			if tagline[i] == ' ' {
				breakPoint = i
				break
			}
		}
		tagline1 = tagline[:breakPoint]
		tagline2 = tagline[breakPoint+1:]
	}

	var lines []string
	lines = append(lines, LabelStyle.Render(" YOUR ARCHETYPE"))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf(" %s  %s", AccentStyle.Render(arch.Icon), titleStyle.Render(arch.Name)))
	if tagline2 != "" {
		lines = append(lines, taglineStyle.Render(fmt.Sprintf(" \"%s", tagline1)))
		lines = append(lines, taglineStyle.Render(fmt.Sprintf("  %s\"", tagline2)))
	} else {
		lines = append(lines, taglineStyle.Render(fmt.Sprintf(" \"%s\"", tagline1)))
	}

	return style.Render(strings.Join(lines, "\n"))
}

func renderQuickStats(stats *analyzer.Stats) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorDim).
		Padding(0, 1).
		Width(TotalWidth)

	headerStyle := lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true)

	// Fixed-width columns for alignment
	colWidth := 14

	// Create stat items with fixed widths
	items := []struct {
		label string
		value string
	}{
		{"Unique Cmds", FormatNumber(stats.UniqueCommands)},
		{"Streak", fmt.Sprintf("%d days", stats.LongestStreak)},
		{"Busiest", formatBusiestDay(stats)},
		{"sudo", sudoMeter(stats.SudoCount)},
		{"Pipes", FormatNumber(stats.PipeCount)},
	}

	header := headerStyle.Render("-- QUICK STATS ") + SubtleStyle.Render(strings.Repeat("-", 57))

	// Build columns with fixed width
	var labelRow strings.Builder
	var valueRow strings.Builder

	for _, item := range items {
		labelRow.WriteString(padRight(LabelStyle.Render(item.label), colWidth))
		valueRow.WriteString(padRight(ValueStyle.Render(item.value), colWidth))
	}

	content := header + "\n" + labelRow.String() + "\n" + valueRow.String()
	return style.Render(content)
}

func padRight(s string, width int) string {
	visibleLen := lipgloss.Width(s)
	if visibleLen >= width {
		return s
	}
	return s + strings.Repeat(" ", width-visibleLen)
}

func formatBusiestDay(stats *analyzer.Stats) string {
	if stats.BusiestDay.IsZero() {
		return "N/A"
	}
	return fmt.Sprintf("%s (%d)", stats.BusiestDay.Format("Jan 2"), stats.BusiestDayCount)
}

func sudoMeter(count int) string {
	// Scale based on count, not percentage
	var level string
	var blocks int
	switch {
	case count == 0:
		level = "none"
		blocks = 0
	case count < 10:
		level = "low"
		blocks = 1
	case count < 50:
		level = "med"
		blocks = 2
	case count < 100:
		level = "high"
		blocks = 3
	case count < 500:
		level = "power"
		blocks = 4
	default:
		level = "god"
		blocks = 5
	}
	filled := lipgloss.NewStyle().Foreground(ColorPrimary).Render(strings.Repeat("#", blocks))
	empty := SubtleStyle.Render(strings.Repeat("-", 5-blocks))
	return "[" + filled + empty + "] " + level
}

func renderTopCommands(stats *analyzer.Stats) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorDim).
		Padding(0, 1).
		Width(LeftColWidth).
		Height(11)

	headerStyle := lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true)

	var lines []string
	lines = append(lines, headerStyle.Render("-- TOP COMMANDS ")+SubtleStyle.Render(strings.Repeat("-", 17)))

	maxCount := 0
	if len(stats.TopCommands) > 0 {
		maxCount = stats.TopCommands[0].Count
	}

	// Show top 8
	for i, cmd := range stats.TopCommands {
		if i >= 8 {
			break
		}

		num := SubtleStyle.Render(fmt.Sprintf("%d.", i+1))
		name := ValueStyle.Render(padRight(cmd.Command, 8))
		bar := ProgressBar(cmd.Count, maxCount, 14, getCmdColor(cmd.Command))
		count := LabelStyle.Render(fmt.Sprintf("%5s", FormatNumber(cmd.Count)))

		lines = append(lines, fmt.Sprintf("%s %s%s %s", num, name, bar, count))
	}

	return style.Render(strings.Join(lines, "\n"))
}

func getCmdColor(cmd string) lipgloss.Color {
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
	catMix := renderCategoryMix(stats)
	heatmap := renderHeatmapSection(stats)
	return lipgloss.JoinVertical(lipgloss.Left, catMix, heatmap)
}

func renderCategoryMix(stats *analyzer.Stats) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorDim).
		Padding(0, 1).
		Width(RightColWidth)

	headerStyle := lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true)

	var lines []string
	lines = append(lines, headerStyle.Render("-- CATEGORIES ")+SubtleStyle.Render(strings.Repeat("-", 19)))

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
	for i := 0; i < len(cats) && i < 8; i += 2 {
		cat1 := cats[i]
		color1 := CategoryColors[cat1.name]
		if color1 == "" {
			color1 = ColorMuted
		}

		// Fixed width columns: bar(4) + space + name(7) + space + pct(3) = 16
		col1 := fmt.Sprintf("%s %-7s%3.0f%%",
			MiniBar(cat1.pct, color1),
			TruncateString(cat1.name, 7),
			cat1.pct)

		col2 := ""
		if i+1 < len(cats) {
			cat2 := cats[i+1]
			color2 := CategoryColors[cat2.name]
			if color2 == "" {
				color2 = ColorMuted
			}
			col2 = fmt.Sprintf("%s %-7s%3.0f%%",
				MiniBar(cat2.pct, color2),
				TruncateString(cat2.name, 7),
				cat2.pct)
		}

		lines = append(lines, padRight(col1, 16)+" "+col2)
	}

	return style.Render(strings.Join(lines, "\n"))
}

func renderHeatmapSection(stats *analyzer.Stats) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorDim).
		Padding(0, 1).
		Width(RightColWidth)

	headerStyle := lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true)

	var lines []string
	lines = append(lines, headerStyle.Render("-- ACTIVITY ")+SubtleStyle.Render(strings.Repeat("-", 21)))

	if stats.HasTimeData {
		heatmapStr := Heatmap(stats.HeatMap)
		lines = append(lines, strings.Split(strings.TrimSuffix(heatmapStr, "\n"), "\n")...)

		peakStyle := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true)
		peak := fmt.Sprintf(">> Peak: %s %02d:00",
			analyzer.GetDayName(stats.PeakDay),
			stats.PeakHour)
		lines = append(lines, peakStyle.Render(peak))
	} else {
		lines = append(lines, "")
		lines = append(lines, SubtleStyle.Render(" No timestamp data"))
		lines = append(lines, SubtleStyle.Render(" Enable EXTENDED_HISTORY"))
		lines = append(lines, "")
	}

	return style.Render(strings.Join(lines, "\n"))
}

func renderFunFacts(stats *analyzer.Stats) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorDim).
		Padding(0, 1).
		Width(TotalWidth)

	headerStyle := lipgloss.NewStyle().Foreground(ColorSecondary).Bold(true)

	var lines []string
	lines = append(lines, headerStyle.Render("-- INSIGHTS ")+SubtleStyle.Render(strings.Repeat("-", 61)))

	// Collect facts as pairs for 2-column layout
	type fact struct {
		icon  string
		label string
		value string
	}
	var facts []fact

	if stats.HasTimeData {
		facts = append(facts, fact{"(O)", "Night Owl", fmt.Sprintf("%.0f%% after midnight", stats.NightOwlPct)})
		facts = append(facts, fact{"[S]", "Weekend", fmt.Sprintf("%.0f%% on Sat/Sun", stats.WeekendPct)})
	}

	if stats.FavoriteDir != "" {
		facts = append(facts, fact{"~/", "Home Dir", TruncateString(stats.FavoriteDir, 18)})
	}

	if stats.EditorChoice != "" {
		facts = append(facts, fact{":w", "Editor", fmt.Sprintf("%s (%s)", stats.EditorChoice, FormatNumber(stats.EditorCount))})
	}

	facts = append(facts, fact{"##", "Avg Length", fmt.Sprintf("%.0f chars", stats.AvgCommandLen)})

	if stats.PipeCount > 0 {
		facts = append(facts, fact{"|>", "Complexity", fmt.Sprintf("%.1f%% use pipes", stats.PipePct)})
	}

	// Render in 2 columns
	colWidth := 36
	for i := 0; i < len(facts); i += 2 {
		f1 := facts[i]
		col1 := fmt.Sprintf("%s %-11s %s",
			AccentStyle.Render(f1.icon),
			LabelStyle.Render(f1.label+":"),
			ValueStyle.Render(f1.value))

		col2 := ""
		if i+1 < len(facts) {
			f2 := facts[i+1]
			col2 = fmt.Sprintf("%s %-11s %s",
				AccentStyle.Render(f2.icon),
				LabelStyle.Render(f2.label+":"),
				ValueStyle.Render(f2.value))
		}

		lines = append(lines, padRight(col1, colWidth)+col2)
	}

	return style.Render(strings.Join(lines, "\n"))
}
