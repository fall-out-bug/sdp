package executor

import (
	"fmt"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/parser"
)

// ParseDependencies parses dependencies from a workstream file
func (e *Executor) ParseDependencies(wsID string) ([]string, error) {
	// Find workstream file
	pattern := filepath.Join(e.config.BacklogDir, fmt.Sprintf("%s-*.md", wsID))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob pattern failed: %w", err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("workstream file not found: %s", wsID)
	}

	// Parse workstream to extract dependencies
	ws, err := parser.ParseWorkstream(matches[0])
	if err != nil {
		return nil, fmt.Errorf("parse workstream: %w", err)
	}

	var deps []string

	// Check if workstream has a parent (dependency)
	if ws.Parent != "" {
		deps = append(deps, ws.Parent)
	}

	// NOTE: Parent dependencies are currently supported via ws.Parent field.
	// Enhancement to parse explicit "Dependencies:" section is tracked in beads issue sdp-fe6q.
	//
	// Future implementation:
	// - Add Dependencies []string to Workstream schema
	// - Parse "## Dependencies" section from workstream markdown
	// - Extract WS IDs and merge with parent dependencies

	return deps, nil
}

// TopologicalSort performs topological sort on workstreams by dependencies
// Uses Kahn's algorithm
func (e *Executor) TopologicalSort(workstreams []string, dependencies map[string][]string) ([]string, error) {
	// Build adjacency list and in-degree count
	inDegree := make(map[string]int)
	adjList := make(map[string][]string)

	// Initialize in-degree for all nodes
	for _, ws := range workstreams {
		inDegree[ws] = 0
		adjList[ws] = []string{}
	}

	// Build graph
	for _, ws := range workstreams {
		for _, dep := range dependencies[ws] {
			// Check if dependency is in our workstream list
			if _, exists := inDegree[dep]; exists {
				adjList[dep] = append(adjList[dep], ws)
				inDegree[ws]++
			}
		}
	}

	// Find all nodes with zero in-degree
	var queue []string
	for ws := range inDegree {
		if inDegree[ws] == 0 {
			queue = append(queue, ws)
		}
	}

	// Process nodes
	var result []string
	for len(queue) > 0 {
		// Remove first node
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// Reduce in-degree for all neighbors
		for _, neighbor := range adjList[current] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// Check for cycles
	if len(result) != len(workstreams) {
		return nil, fmt.Errorf("cycle detected in workstream dependencies")
	}

	return result, nil
}
