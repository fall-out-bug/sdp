package executor

import (
	"fmt"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/parser"
)

// findReadyWorkstreams finds all workstreams that are ready to execute
// (no blockers or blockers already completed)
func (e *Executor) findReadyWorkstreams() ([]string, error) {
	var workstreams []string

	// List all workstream files in backlog
	pattern := filepath.Join(e.config.BacklogDir, "*.md")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob pattern failed: %w", err)
	}

	for _, match := range matches {
		// Parse workstream
		ws, err := parser.ParseWorkstream(match)
		if err != nil {
			continue // Skip unparseable files
		}

		// Check if workstream is in ready state
		if ws.Status == "pending" || ws.Status == "ready" {
			workstreams = append(workstreams, ws.ID)
		}
	}

	return workstreams, nil
}
