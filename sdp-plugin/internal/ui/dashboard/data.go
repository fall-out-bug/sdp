package dashboard

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/beads"
	"github.com/fall-out-bug/sdp/internal/nextstep"
	"github.com/fall-out-bug/sdp/internal/parser"
)

// fetchWorkstreams fetches workstream data from Beads and docs/workstreams/
func (a *App) fetchWorkstreams() map[string][]WorkstreamSummary {
	// Try to fetch from Beads first
	summaries := make(map[string][]WorkstreamSummary)

	// Initialize all status groups
	summaries["open"] = []WorkstreamSummary{}
	summaries["in_progress"] = []WorkstreamSummary{}
	summaries["completed"] = []WorkstreamSummary{}
	summaries["blocked"] = []WorkstreamSummary{}

	// Try to get data from Beads (if available)
	client, err := beads.NewClient()
	if err == nil {
		// Beads is available, fetch tasks
		tasks, err := client.Ready()
		if err == nil && len(tasks) > 0 {
			for _, task := range tasks {
				status := strings.ToLower(task.Status)
				if status == "ready" || status == "open" {
					status = "open"
				}

				// Map Beads priority to our format
				priority := task.Priority
				if strings.HasPrefix(priority, "P") {
					// Already in P0-P4 format
				} else {
					// Convert numeric to P format
					switch priority {
					case "0":
						priority = "P0"
					case "1":
						priority = "P1"
					case "2":
						priority = "P2"
					case "3":
						priority = "P3"
					case "4":
						priority = "P4"
					default:
						priority = "P2" // Default
					}
				}

				ws := WorkstreamSummary{
					ID:       task.ID,
					Title:    task.Title,
					Status:   status,
					Priority: priority,
					Assignee: "", // Task doesn't have Assignee field
				}

				if group, ok := summaries[status]; ok {
					summaries[status] = append(group, ws)
				}
			}
		}
	}

	// Always check docs/workstreams/ as additional source
	// (even if Beads returned data, this might have different items)
	wsFiles, err := filepath.Glob("docs/workstreams/*/backlog/*.md")
	if err == nil && len(wsFiles) > 0 {
		for _, wsFile := range wsFiles {
			ws, err := parser.ParseWorkstream(wsFile)
			if err != nil {
				continue
			}

			status := strings.ToLower(ws.Status)
			wsSummary := WorkstreamSummary{
				ID:       ws.ID,
				Title:    ws.Goal,
				Status:   status,
				Priority: "P2", // Default priority
				Size:     ws.Size,
			}

			// Add to appropriate status group (avoiding duplicates by ID)
			group := summaries[status]
			alreadyExists := false
			for _, existing := range group {
				if existing.ID == wsSummary.ID {
					alreadyExists = true
					break
				}
			}
			if !alreadyExists {
				summaries[status] = append(group, wsSummary)
			}
		}
	}

	// Sort each group by ID
	for status := range summaries {
		sort.Slice(summaries[status], func(i, j int) bool {
			return summaries[status][i].ID < summaries[status][j].ID
		})
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

// fetchTestResults fetches test results
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

// fetchNextStep fetches the next step recommendation
func (a *App) fetchNextStep() NextStepInfo {
	// Get project root
	projectRoot, err := os.Getwd()
	if err != nil {
		return NextStepInfo{
			Command:    "sdp doctor",
			Reason:     "Check environment setup",
			Confidence: 0.5,
			Category:   "setup",
		}
	}

	// Collect state
	collector := nextstep.NewStateCollector(projectRoot)
	state, err := collector.Collect()
	if err != nil {
		return NextStepInfo{
			Command:    "sdp doctor",
			Reason:     "Unable to collect project state",
			Confidence: 0.5,
			Category:   "setup",
		}
	}

	// Get recommendation
	resolver := nextstep.NewResolver()
	rec, err := resolver.Recommend(state)
	if err != nil {
		return NextStepInfo{
			Command:    "sdp status",
			Reason:     "Check current project state",
			Confidence: 0.5,
			Category:   "information",
		}
	}

	return NextStepInfo{
		Command:    rec.Command,
		Reason:     rec.Reason,
		Confidence: rec.Confidence,
		Category:   string(rec.Category),
	}
}
