package graph

import (
	"fmt"
)

// WorkstreamFile represents a workstream file
type WorkstreamFile struct {
	ID        string
	DependsOn []string
}

// ExecuteResult represents the result of executing a workstream
type ExecuteResult struct {
	WorkstreamID string
	Success      bool
	Error        error
	Duration     int64 // Duration in milliseconds
}

// ExecuteFunc is a function that executes a single workstream
type ExecuteFunc func(wsID string) error

// BuildGraphFromWSFiles creates a dependency graph from workstream files
func BuildGraphFromWSFiles(workstreams []WorkstreamFile) (*DependencyGraph, error) {
	graph := NewDependencyGraph()

	// First pass: add all nodes
	for _, ws := range workstreams {
		err := graph.AddNode(ws.ID, ws.DependsOn)
		if err != nil {
			return nil, fmt.Errorf("failed to add workstream %s: %w", ws.ID, err)
		}
	}

	return graph, nil
}

// NewDispatcherWithCheckpoint creates a new dispatcher with checkpoint support
func NewDispatcherWithCheckpoint(g *DependencyGraph, concurrency int, featureID string, enableCheckpoint bool) *Dispatcher {
	d := NewDispatcher(g, concurrency)
	if enableCheckpoint {
		d.featureID = featureID
		d.enableCheckpoint = true
		d.checkpointManager = NewCheckpointManager(featureID)
	}
	return d
}
