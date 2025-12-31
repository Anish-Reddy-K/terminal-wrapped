package analyzer

// Archetype represents a developer personality type
type Archetype struct {
	Name    string
	Emoji   string
	Tagline string
	Score   float64 // Higher = stronger match
}

// All available archetypes
var archetypes = []struct {
	Name    string
	Emoji   string
	Tagline string
	Detect  func(*Stats) float64
}{
	{
		Name:    "THE SUDO SUMMONER",
		Emoji:   "⚡",
		Tagline: "With great power comes great responsibility",
		Detect: func(s *Stats) float64 {
			if s.SudoPct > 15 {
				return s.SudoPct * 2
			}
			return s.SudoPct
		},
	},
	{
		Name:    "THE GIT GLADIATOR",
		Emoji:   "⚔️",
		Tagline: "Commit early, commit often",
		Detect: func(s *Stats) float64 {
			gitPct := s.CategoryPct["Git"]
			if len(s.TopCommands) > 0 && s.TopCommands[0].Command == "git" {
				return gitPct * 3 // Bonus if git is #1
			}
			return gitPct * 1.5
		},
	},
	{
		Name:    "THE DOCKER CAPTAIN",
		Emoji:   "🐳",
		Tagline: "It works in my container",
		Detect: func(s *Stats) float64 {
			return s.CategoryPct["Containers"] * 2.5
		},
	},
	{
		Name:    "THE PACKAGE GOBLIN",
		Emoji:   "📦",
		Tagline: "Just one more dependency...",
		Detect: func(s *Stats) float64 {
			return s.CategoryPct["Packages"] * 2
		},
	},
	{
		Name:    "THE VIM WIZARD",
		Emoji:   "🧙",
		Tagline: "I use vim btw",
		Detect: func(s *Stats) float64 {
			if s.EditorChoice == "vim" || s.EditorChoice == "nvim" {
				// Check if in top 5
				for i, cmd := range s.TopCommands {
					if i >= 5 {
						break
					}
					if cmd.Command == "vim" || cmd.Command == "nvim" {
						return float64(cmd.Count) / float64(s.TotalCommands) * 100 * 5
					}
				}
				return float64(s.EditorCount) / float64(s.TotalCommands) * 100 * 3
			}
			return 0
		},
	},
	{
		Name:    "THE SSH NOMAD",
		Emoji:   "🌍",
		Tagline: "My servers miss me",
		Detect: func(s *Stats) float64 {
			return s.CategoryPct["Network"] * 2
		},
	},
	{
		Name:    "THE PIPE PLUMBER",
		Emoji:   "🔧",
		Tagline: "Data flows through me",
		Detect: func(s *Stats) float64 {
			if s.PipePct > 20 {
				return s.PipePct * 2
			}
			return s.PipePct
		},
	},
	{
		Name:    "THE SCRIPT SORCERER",
		Emoji:   "🪄",
		Tagline: "Why do it twice when you can automate?",
		Detect: func(s *Stats) float64 {
			// Look for script execution patterns
			var score float64
			for _, cmd := range s.TopCommands {
				if cmd.Command == "bash" || cmd.Command == "sh" || cmd.Command == "python" || cmd.Command == "python3" || cmd.Command == "node" {
					score += float64(cmd.Count) / float64(s.TotalCommands) * 100
				}
			}
			return score * 3
		},
	},
	{
		Name:    "THE DEBUG DETECTIVE",
		Emoji:   "🔍",
		Tagline: "The bug is in here somewhere...",
		Detect: func(s *Stats) float64 {
			return s.CategoryPct["Search"] * 3
		},
	},
	{
		Name:    "THE CLEAN FREAK",
		Emoji:   "🧹",
		Tagline: "Disk space is sacred",
		Detect: func(s *Stats) float64 {
			var score float64
			for _, cmd := range s.TopCommands {
				switch cmd.Command {
				case "rm", "rmdir", "clean", "prune", "gc":
					score += float64(cmd.Count) / float64(s.TotalCommands) * 100
				}
			}
			return score * 5
		},
	},
	{
		Name:    "THE NIGHT OWL",
		Emoji:   "🦉",
		Tagline: "Best code is written after midnight",
		Detect: func(s *Stats) float64 {
			if s.NightOwlPct > 15 {
				return s.NightOwlPct * 2
			}
			return s.NightOwlPct
		},
	},
	{
		Name:    "THE GENERALIST",
		Emoji:   "🎯",
		Tagline: "Jack of all trades, master of many",
		Detect: func(s *Stats) float64 {
			// Score based on diversity of categories used
			categoriesUsed := 0
			for _, pct := range s.CategoryPct {
				if pct > 2 {
					categoriesUsed++
				}
			}
			if categoriesUsed >= 6 {
				return float64(categoriesUsed) * 5
			}
			return 0
		},
	},
}

// DetectArchetype finds the best matching archetype for the user
func DetectArchetype(stats *Stats) *Archetype {
	var bestArchetype *Archetype
	var bestScore float64

	for _, arch := range archetypes {
		score := arch.Detect(stats)
		if score > bestScore {
			bestScore = score
			bestArchetype = &Archetype{
				Name:    arch.Name,
				Emoji:   arch.Emoji,
				Tagline: arch.Tagline,
				Score:   score,
			}
		}
	}

	// Default fallback
	if bestArchetype == nil || bestScore < 1 {
		return &Archetype{
			Name:    "THE TERMINAL WARRIOR",
			Emoji:   "⌨️",
			Tagline: "Command line is my home",
			Score:   0,
		}
	}

	return bestArchetype
}

// GetSecondaryArchetypes returns other notable archetypes
func GetSecondaryArchetypes(stats *Stats, primary *Archetype) []*Archetype {
	var secondary []*Archetype

	for _, arch := range archetypes {
		if arch.Name == primary.Name {
			continue
		}
		score := arch.Detect(stats)
		if score > 5 { // Threshold for notable
			secondary = append(secondary, &Archetype{
				Name:    arch.Name,
				Emoji:   arch.Emoji,
				Tagline: arch.Tagline,
				Score:   score,
			})
		}
	}

	return secondary
}

