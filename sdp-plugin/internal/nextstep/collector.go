package nextstep

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fall-out-bug/sdp/internal/parser"
)

// StateCollector gathers project state from various sources.
type StateCollector struct {
	projectRoot string
}

// NewStateCollector creates a new state collector.
func NewStateCollector(projectRoot string) *StateCollector {
	return &StateCollector{projectRoot: projectRoot}
}

// Collect gathers the current project state for recommendation.
func (c *StateCollector) Collect() (ProjectState, error) {
	state := ProjectState{
		Workstreams: c.collectWorkstreams(),
		GitStatus:   c.collectGitStatus(),
		Config:      c.collectConfig(),
		Mode:        ModeDrive, // Default mode
	}

	return state, nil
}

// collectWorkstreams gathers workstream status from docs/workstreams/.
func (c *StateCollector) collectWorkstreams() []WorkstreamStatus {
	workstreams := []WorkstreamStatus{}

	// Search for workstream files
	patterns := []string{
		filepath.Join(c.projectRoot, "docs", "workstreams", "*", "backlog", "*.md"),
		filepath.Join(c.projectRoot, "docs", "workstreams", "backlog", "*.md"),
	}

	seen := make(map[string]bool)

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

			// Avoid duplicates
			if seen[ws.ID] {
				continue
			}
			seen[ws.ID] = true

			status := WorkstreamStatus{
				ID:        ws.ID,
				Feature:   ws.Feature,
				Status:    mapStatus(ws.Status),
				Priority:  0, // Default priority
				Size:      ws.Size,
				BlockedBy: extractBlockers(ws),
			}

			workstreams = append(workstreams, status)
		}
	}

	return workstreams
}

// mapStatus converts string status to WorkstreamState.
func mapStatus(status string) WorkstreamState {
	switch strings.ToLower(status) {
	case "backlog":
		return StatusBacklog
	case "ready", "open":
		return StatusReady
	case "in_progress", "in-progress", "started":
		return StatusInProgress
	case "blocked":
		return StatusBlocked
	case "completed", "done":
		return StatusCompleted
	case "failed", "error":
		return StatusFailed
	default:
		return StatusBacklog
	}
}

// extractBlockers extracts blocking dependencies from a workstream.
func extractBlockers(ws *parser.Workstream) []string {
	// Parse depends_on from frontmatter if available
	// For now, return empty - will be enhanced when frontmatter supports deps
	return nil
}

// collectGitStatus gathers git repository status.
func (c *StateCollector) collectGitStatus() GitStatusInfo {
	info := GitStatusInfo{
		IsRepo:     false,
		MainBranch: "main",
	}

	gitDir := filepath.Join(c.projectRoot, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		info.IsRepo = true
	}

	// Check for uncommitted changes
	if info.IsRepo {
		// Simple check: if git status --porcelain returns anything, there are changes
		// This is a simplified implementation; the real one would use git package
		info.Uncommitted = false // Will be populated by git package integration
	}

	return info
}

// collectConfig gathers SDP configuration.
func (c *StateCollector) collectConfig() ConfigInfo {
	info := ConfigInfo{
		HasSDPConfig: false,
	}

	configPath := filepath.Join(c.projectRoot, ".sdp", "config.yml")
	if _, err := os.Stat(configPath); err == nil {
		info.HasSDPConfig = true
	}

	configJSONPath := filepath.Join(c.projectRoot, ".sdp", "config.json")
	if _, err := os.Stat(configJSONPath); err == nil {
		info.HasSDPConfig = true
	}

	// Check for evidence log
	evidencePath := filepath.Join(c.projectRoot, ".sdp", "log", "events.jsonl")
	if _, err := os.Stat(evidencePath); err == nil {
		info.EvidenceEnabled = true
	}

	info.ProjectRoot = c.projectRoot

	return info
}
