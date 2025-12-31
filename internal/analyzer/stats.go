package analyzer

import (
	"sort"
	"strings"
	"time"

	"github.com/Anish-Reddy-K/terminal-wrapped/internal/parser"
)

// Stats contains all computed statistics
type Stats struct {
	// Basic counts
	TotalCommands  int
	UniqueCommands int

	// Time-based stats (only if timestamps available)
	HasTimeData    bool
	FirstCommand   time.Time
	LastCommand    time.Time
	HistorySpan    time.Duration
	CommandsPerDay float64
	LongestStreak  int // consecutive days
	BusiestDay     time.Time
	BusiestDayCount int

	// Hour/Day heatmap [day][hour] = count (0=Sunday, 1=Monday, ...)
	HeatMap [7][24]int

	// Peak time
	PeakHour    int
	PeakDay     int // 0=Sunday
	NightOwlPct float64 // % of commands after midnight (0-5am)
	WeekendPct  float64 // % on Saturday/Sunday

	// Top commands
	TopCommands []CommandCount
	
	// Category breakdown
	Categories map[string]int
	CategoryPct map[string]float64

	// Command analysis
	SudoCount       int
	SudoPct         float64
	PipeCount       int
	PipePct         float64
	AvgCommandLen   float64
	LongestCommand  string
	LongestCmdLen   int

	// Fun facts
	MostRepeated      string
	MostRepeatedCount int
	FavoriteDir       string
	FavoriteDirCount  int
	EditorChoice      string
	EditorCount       int
}

// CommandCount holds a command and its count
type CommandCount struct {
	Command string
	Count   int
}

// Category definitions
var categoryCommands = map[string][]string{
	"Git":        {"git", "gh", "hub", "tig", "lazygit", "gitui"},
	"Containers": {"docker", "podman", "kubectl", "k9s", "helm", "docker-compose", "minikube", "kind"},
	"Packages":   {"npm", "yarn", "pnpm", "pip", "pip3", "cargo", "brew", "apt", "apt-get", "yum", "dnf", "pacman", "go"},
	"Editors":    {"vim", "nvim", "nano", "emacs", "code", "cursor", "subl", "atom", "micro", "hx", "helix"},
	"Navigation": {"cd", "ls", "pwd", "tree", "z", "autojump", "j", "exa", "eza", "lsd", "ll", "la"},
	"Search":     {"grep", "rg", "ag", "find", "fd", "fzf", "ack", "locate"},
	"Network":    {"ssh", "scp", "curl", "wget", "rsync", "ping", "nc", "netcat", "telnet", "sftp"},
	"Files":      {"cat", "less", "head", "tail", "rm", "cp", "mv", "mkdir", "touch", "chmod", "chown", "ln", "bat"},
}

// Analyze computes all statistics from history data
func Analyze(data *parser.HistoryData) *Stats {
	stats := &Stats{
		TotalCommands: len(data.Commands),
		HasTimeData:   data.HasTimes,
		Categories:    make(map[string]int),
		CategoryPct:   make(map[string]float64),
	}

	if stats.TotalCommands == 0 {
		return stats
	}

	// Count commands
	commandCounts := make(map[string]int)
	// Track consecutive repeats
	consecutiveRepeats := make(map[string]int)
	lastCmd := ""
	currentRepeat := 0

	// Track directories
	dirCounts := make(map[string]int)
	
	// Track editors
	editorCounts := make(map[string]int)
	editors := map[string]bool{"vim": true, "nvim": true, "nano": true, "emacs": true, "code": true, "cursor": true, "subl": true, "micro": true, "hx": true, "helix": true}

	// Track days for streak calculation
	activeDays := make(map[string]bool)

	// Total command length for average
	totalCmdLen := 0

	for _, cmd := range data.Commands {
		baseCmd := parser.GetBaseCommand(&cmd)
		commandCounts[baseCmd]++

		// Track raw command length (use parsed command, not raw with timestamp)
		cmdStr := cmd.Raw
		totalCmdLen += len(cmdStr)
		if len(cmdStr) > stats.LongestCmdLen {
			stats.LongestCmdLen = len(cmdStr)
			stats.LongestCommand = cmdStr
		}

		// Track consecutive repeats
		if cmd.Raw == lastCmd {
			currentRepeat++
		} else {
			if currentRepeat > stats.MostRepeatedCount {
				stats.MostRepeatedCount = currentRepeat
				stats.MostRepeated = lastCmd
			}
			currentRepeat = 1
		}
		lastCmd = cmd.Raw
		consecutiveRepeats[cmd.Raw] = max(consecutiveRepeats[cmd.Raw], currentRepeat)

		// Track sudo
		if cmd.Command == "sudo" || cmd.Command == "doas" {
			stats.SudoCount++
		}

		// Track pipes
		if strings.Contains(cmd.Raw, "|") {
			stats.PipeCount++
		}

		// Track directories (cd commands)
		if baseCmd == "cd" && len(cmd.Args) > 0 {
			dir := cmd.Args[0]
			// Normalize some common patterns
			if strings.HasPrefix(dir, "~/") || dir == "~" {
				dir = "~" + dir[1:]
			}
			dirCounts[dir]++
		}

		// Track editors
		if editors[baseCmd] {
			editorCounts[baseCmd]++
		}

		// Categorize
		for category, commands := range categoryCommands {
			for _, c := range commands {
				if baseCmd == c {
					stats.Categories[category]++
					break
				}
			}
		}

		// Time-based analysis
		if cmd.HasTime {
			if stats.FirstCommand.IsZero() || cmd.Timestamp.Before(stats.FirstCommand) {
				stats.FirstCommand = cmd.Timestamp
			}
			if cmd.Timestamp.After(stats.LastCommand) {
				stats.LastCommand = cmd.Timestamp
			}

			// Track active days for streak
			dayKey := cmd.Timestamp.Format("2006-01-02")
			activeDays[dayKey] = true

			// Heatmap
			hour := cmd.Timestamp.Hour()
			day := int(cmd.Timestamp.Weekday())
			stats.HeatMap[day][hour]++

			// Night owl (midnight to 5am)
			if hour >= 0 && hour < 5 {
				stats.NightOwlPct++
			}

			// Weekend
			if day == 0 || day == 6 {
				stats.WeekendPct++
			}
		}
	}

	// Check last repeat
	if currentRepeat > stats.MostRepeatedCount {
		stats.MostRepeatedCount = currentRepeat
		stats.MostRepeated = lastCmd
	}

	// Unique commands
	stats.UniqueCommands = len(commandCounts)

	// Calculate percentages
	total := float64(stats.TotalCommands)
	stats.SudoPct = float64(stats.SudoCount) / total * 100
	stats.PipePct = float64(stats.PipeCount) / total * 100
	stats.AvgCommandLen = float64(totalCmdLen) / total

	if stats.HasTimeData {
		stats.NightOwlPct = stats.NightOwlPct / total * 100
		stats.WeekendPct = stats.WeekendPct / total * 100
		stats.HistorySpan = stats.LastCommand.Sub(stats.FirstCommand)
		
		days := stats.HistorySpan.Hours() / 24
		if days > 0 {
			stats.CommandsPerDay = total / days
		}
	}

	// Top commands
	stats.TopCommands = topN(commandCounts, 10)

	// Category percentages
	for cat, count := range stats.Categories {
		stats.CategoryPct[cat] = float64(count) / total * 100
	}

	// Find peak hour/day
	maxHeatmap := 0
	for day := 0; day < 7; day++ {
		for hour := 0; hour < 24; hour++ {
			if stats.HeatMap[day][hour] > maxHeatmap {
				maxHeatmap = stats.HeatMap[day][hour]
				stats.PeakDay = day
				stats.PeakHour = hour
			}
		}
	}

	// Longest streak
	if len(activeDays) > 0 {
		stats.LongestStreak, stats.BusiestDay, stats.BusiestDayCount = calculateStreak(activeDays, data.Commands)
	}

	// Favorite directory
	maxDir := 0
	for dir, count := range dirCounts {
		if count > maxDir {
			maxDir = count
			stats.FavoriteDir = dir
			stats.FavoriteDirCount = count
		}
	}

	// Editor choice
	maxEditor := 0
	for editor, count := range editorCounts {
		if count > maxEditor {
			maxEditor = count
			stats.EditorChoice = editor
			stats.EditorCount = count
		}
	}

	return stats
}

func topN(counts map[string]int, n int) []CommandCount {
	result := make([]CommandCount, 0, len(counts))
	for cmd, count := range counts {
		result = append(result, CommandCount{Command: cmd, Count: count})
	}
	
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	if len(result) > n {
		result = result[:n]
	}
	return result
}

func calculateStreak(activeDays map[string]bool, commands []parser.Command) (longestStreak int, busiestDay time.Time, busiestCount int) {
	if len(activeDays) == 0 {
		return 0, time.Time{}, 0
	}

	// Count commands per day
	dayCounts := make(map[string]int)
	for _, cmd := range commands {
		if cmd.HasTime {
			dayKey := cmd.Timestamp.Format("2006-01-02")
			dayCounts[dayKey]++
		}
	}

	// Find busiest day
	for dayStr, count := range dayCounts {
		if count > busiestCount {
			busiestCount = count
			busiestDay, _ = time.Parse("2006-01-02", dayStr)
		}
	}

	// Sort days
	days := make([]string, 0, len(activeDays))
	for day := range activeDays {
		days = append(days, day)
	}
	sort.Strings(days)

	// Calculate streak
	currentStreak := 1
	for i := 1; i < len(days); i++ {
		prevDay, _ := time.Parse("2006-01-02", days[i-1])
		currDay, _ := time.Parse("2006-01-02", days[i])
		
		if currDay.Sub(prevDay).Hours() <= 24 {
			currentStreak++
		} else {
			if currentStreak > longestStreak {
				longestStreak = currentStreak
			}
			currentStreak = 1
		}
	}
	
	if currentStreak > longestStreak {
		longestStreak = currentStreak
	}

	return longestStreak, busiestDay, busiestCount
}

// GetSudoLevel returns a fun label for sudo usage
func GetSudoLevel(pct float64) string {
	switch {
	case pct < 1:
		return "Peasant"
	case pct < 5:
		return "Apprentice"
	case pct < 10:
		return "Elevated"
	case pct < 20:
		return "Power User"
	default:
		return "Root God"
	}
}

// GetDayName returns the day name from weekday int
func GetDayName(day int) string {
	days := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	if day >= 0 && day < 7 {
		return days[day]
	}
	return "Unknown"
}

// FormatDuration formats a duration nicely
func FormatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	years := days / 365
	months := (days % 365) / 30
	
	if years > 0 {
		if months > 0 {
			return strings.TrimSuffix(strings.TrimSuffix(
				formatPlural(years, "year") + " " + formatPlural(months, "month"),
				" 0 months"), " ")
		}
		return formatPlural(years, "year")
	}
	if months > 0 {
		return formatPlural(months, "month")
	}
	if days > 0 {
		return formatPlural(days, "day")
	}
	return "< 1 day"
}

func formatPlural(n int, unit string) string {
	if n == 1 {
		return "1 " + unit
	}
	return strings.TrimSpace(strings.Replace(string(rune(n+'0'))+" "+unit+"s", string(rune(n+'0')), itoa(n), 1))
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	
	var result strings.Builder
	negative := n < 0
	if negative {
		n = -n
	}
	
	digits := make([]byte, 0, 10)
	for n > 0 {
		digits = append(digits, byte(n%10)+'0')
		n /= 10
	}
	
	if negative {
		result.WriteByte('-')
	}
	
	for i := len(digits) - 1; i >= 0; i-- {
		result.WriteByte(digits[i])
	}
	
	return result.String()
}

