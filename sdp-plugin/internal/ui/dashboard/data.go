package dashboard

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/controltower"
	"github.com/fall-out-bug/sdp/internal/nextstep"
	"github.com/fall-out-bug/sdp/internal/parser"
)

func (a *App) fetchWorkstreams() map[string][]WorkstreamSummary {
	return a.workstreamsFromControlTower(a.fetchControlTower(false))
}

func (a *App) fetchControlTower(force bool) *controltower.Data {
	if !force && a.controlTowerData != nil && time.Since(a.controlTowerLoaded) < refreshInterval {
		return a.controlTowerData
	}

	data, err := controltower.Collect()
	if err != nil {
		return nil
	}
	a.controlTowerData = data
	a.controlTowerLoaded = time.Now()
	return data
}

func (a *App) workstreamsFromControlTower(data *controltower.Data) map[string][]WorkstreamSummary {
	// Try to fetch from Beads first
	summaries := make(map[string][]WorkstreamSummary)

	// Initialize all status groups
	summaries["open"] = []WorkstreamSummary{}
	summaries["in_progress"] = []WorkstreamSummary{}
	summaries["completed"] = []WorkstreamSummary{}
	summaries["blocked"] = []WorkstreamSummary{}
	if data == nil {
		return summaries
	}

	titles := loadWorkstreamTitles(data.ProjectRoot)
	grouped := map[string][]nextstep.WorkstreamStatus{
		"open":        {},
		"in_progress": {},
		"completed":   {},
		"blocked":     {},
	}

	for _, ws := range data.State.Workstreams {
		group := dashboardGroup(ws)
		grouped[group] = append(grouped[group], ws)
	}

	for group, items := range grouped {
		sort.Slice(items, func(i, j int) bool {
			return nextstep.ComparePriority(items[i], items[j]) < 0
		})
		for _, ws := range items {
			summaries[group] = append(summaries[group], WorkstreamSummary{
				ID:       ws.ID,
				Title:    workstreamTitle(ws, titles),
				Status:   group,
				Priority: formatPriority(ws.Priority),
				Size:     ws.Size,
			})
		}
	}

	return summaries
}

// fetchIdeas fetches ideas from docs/drafts/
func (a *App) fetchIdeas() []IdeaSummary {
	ideas := []IdeaSummary{}

	// Find all markdown files in docs/drafts/
	ideaFiles, err := filepath.Glob("docs/drafts/*.md")
	if err != nil {
		return ideas
	}

	for _, ideaFile := range ideaFiles {
		// Get file info for modification time
		info, err := os.Stat(ideaFile)
		if err != nil {
			continue
		}

		// Extract title from filename
		filename := filepath.Base(ideaFile)
		title := strings.TrimSuffix(filename, ".md")
		title = strings.ReplaceAll(title, "-", " ")
		// Capitalize first letter (strings.Title is deprecated)
		if len(title) > 0 {
			title = strings.ToUpper(title[:1]) + title[1:]
		}

		ideas = append(ideas, IdeaSummary{
			Title: title,
			Path:  ideaFile,
			Date:  info.ModTime(),
		})
	}

	// Sort by date (newest first)
	sort.Slice(ideas, func(i, j int) bool {
		return ideas[i].Date.After(ideas[j].Date)
	})

	return ideas
}

func (a *App) fetchTestResults() TestSummary {
	// Future: Integrate with go test -cover output
	// For now, return placeholder data (dashboard MVP)
	return TestSummary{
		Coverage: "N/A",
		Passing:  0,
		Total:    0,
		LastRun:  time.Now(),
		QualityGates: []GateStatus{
			{Name: "Coverage", Passed: false},
			{Name: "Type Hints", Passed: false},
			{Name: "Linting", Passed: false},
		},
	}
}

func (a *App) fetchNextStep() NextStepInfo {
	return a.nextStepFromControlTower(a.fetchControlTower(false))
}

func (a *App) nextStepFromControlTower(data *controltower.Data) NextStepInfo {
	if data == nil || data.NextStep == nil {
		return NextStepInfo{
			Command:    "sdp doctor",
			Reason:     "Check environment setup",
			Confidence: 0.5,
			Category:   "setup",
		}
	}

	rec := data.NextStep
	return NextStepInfo{
		Command:    rec.Command,
		Reason:     rec.Reason,
		Confidence: rec.Confidence,
		Category:   string(rec.Category),
	}
}

func dashboardGroup(ws nextstep.WorkstreamStatus) string {
	switch ws.Status {
	case nextstep.StatusInProgress:
		return "in_progress"
	case nextstep.StatusCompleted:
		return "completed"
	case nextstep.StatusBlocked, nextstep.StatusFailed:
		return "blocked"
	case nextstep.StatusReady:
		if len(ws.BlockedBy) > 0 {
			return "blocked"
		}
		return "open"
	default:
		return "open"
	}
}

func formatPriority(priority int) string {
	if priority < 0 {
		priority = 0
	}
	return fmt.Sprintf("P%d", priority)
}

func workstreamTitle(ws nextstep.WorkstreamStatus, titles map[string]string) string {
	if title := titles[ws.ID]; title != "" {
		return title
	}
	if ws.Feature != "" {
		return ws.Feature
	}
	if ws.LastError != "" {
		return ws.LastError
	}
	return "Workstream"
}

func loadWorkstreamTitles(projectRoot string) map[string]string {
	titles := map[string]string{}
	patterns := []string{
		filepath.Join(projectRoot, "docs", "workstreams", "*", "backlog", "*.md"),
		filepath.Join(projectRoot, "docs", "workstreams", "backlog", "*.md"),
	}
	for _, pattern := range patterns {
		files, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		for _, file := range files {
			ws, err := parser.ParseWorkstream(file)
			if err != nil {
				continue
			}
			if _, exists := titles[ws.ID]; exists {
				continue
			}
			titles[ws.ID] = ws.Goal
		}
	}
	return titles
}
