package orchestrator

import (
	"errors"
	"fmt"

	"github.com/ai-masters/sdp/internal/checkpoint"
)

// Common errors
var (
	ErrFeatureNotFound     = errors.New("feature not found")
	ErrCheckpointNotFound  = errors.New("checkpoint not found")
	ErrExecutionFailed     = errors.New("workstream execution failed")
	ErrCircularDependency  = errors.New("circular dependency detected")
	ErrMissingDependency   = errors.New("missing dependency")
	ErrWorkstreamNotFound  = errors.New("workstream not found")
)

// WorkstreamNode represents a workstream in the dependency graph
type WorkstreamNode struct {
	ID           string
	Feature      string
	Status       string
	Dependencies []string
}

// DependencyNode represents a node in the dependency graph
type DependencyNode struct {
	Workstream   WorkstreamNode
	Dependencies []*DependencyNode
	InDegree     int // Number of incoming edges
}

// Graph represents the dependency graph
type Graph map[string]*DependencyNode

// WorkstreamLoader defines the interface for loading workstreams
type WorkstreamLoader interface {
	LoadWorkstreams(featureID string) ([]WorkstreamNode, error)
}

// WorkstreamExecutor defines the interface for executing workstreams
type WorkstreamExecutor interface {
	Execute(wsID string) error
}

// CheckpointSaver defines the interface for saving checkpoints
type CheckpointSaver interface {
	Save(cp checkpoint.Checkpoint) error
	Load(id string) (checkpoint.Checkpoint, error)
	Resume(id string) (checkpoint.Checkpoint, error)
}

// Orchestrator manages workstream execution with dependency tracking
type Orchestrator struct {
	loader     WorkstreamLoader
	executor   WorkstreamExecutor
	checkpoint CheckpointSaver
	maxRetries int
}

// NewOrchestrator creates a new orchestrator
func NewOrchestrator(
	loader WorkstreamLoader,
	executor WorkstreamExecutor,
	checkpoint CheckpointSaver,
	maxRetries int,
) *Orchestrator {
	return &Orchestrator{
		loader:     loader,
		executor:   executor,
		checkpoint: checkpoint,
		maxRetries: maxRetries,
	}
}

// BuildDependencyGraph builds a dependency graph from workstreams
func BuildDependencyGraph(workstreams []WorkstreamNode) (Graph, error) {
	graph := make(Graph)

	// Create nodes
	for _, ws := range workstreams {
		graph[ws.ID] = &DependencyNode{
			Workstream: ws,
			InDegree:   0,
		}
	}

	// Add edges
	for _, ws := range workstreams {
		node := graph[ws.ID]
		for _, depID := range ws.Dependencies {
			// Check for self-dependency
			if depID == ws.ID {
				return nil, fmt.Errorf("%w: workstream %s depends on itself",
					ErrCircularDependency, ws.ID)
			}

			// Check if dependency exists
			depNode, exists := graph[depID]
			if !exists {
				return nil, fmt.Errorf("%w: workstream %s depends on non-existent %s",
					ErrMissingDependency, ws.ID, depID)
			}

			// Add edge
			depNode.Dependencies = append(depNode.Dependencies, node)
			node.InDegree++
		}
	}

	// Check for cycles
	if hasCycle(graph) {
		return nil, ErrCircularDependency
	}

	return graph, nil
}

// TopologicalSort performs topological sort on the dependency graph
func TopologicalSort(graph Graph) ([]string, error) {
	if len(graph) == 0 {
		return []string{}, nil
	}

	// Kahn's algorithm
	inDegree := make(map[string]int)
	queue := []string{}

	// Initialize in-degree copy and find nodes with no dependencies
	for id, node := range graph {
		inDegree[id] = node.InDegree
		if node.InDegree == 0 {
			queue = append(queue, id)
		}
	}

	result := []string{}

	for len(queue) > 0 {
		// Dequeue
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// Reduce in-degree for all neighbors
		node := graph[current]
		for _, neighbor := range node.Dependencies {
			inDegree[neighbor.Workstream.ID]--
			if inDegree[neighbor.Workstream.ID] == 0 {
				queue = append(queue, neighbor.Workstream.ID)
			}
		}
	}

	// Check if all nodes were processed (cycle detection)
	if len(result) != len(graph) {
		return nil, ErrCircularDependency
	}

	return result, nil
}

// hasCycle checks if the graph has a cycle using DFS
func hasCycle(graph Graph) bool {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for id := range graph {
		if !visited[id] {
			if hasCycleDFS(graph, id, visited, recStack) {
				return true
			}
		}
	}

	return false
}

// hasCycleDFS performs DFS to detect cycles
func hasCycleDFS(graph Graph, nodeID string, visited, recStack map[string]bool) bool {
	visited[nodeID] = true
	recStack[nodeID] = true

	node := graph[nodeID]
	for _, neighbor := range node.Dependencies {
		neighborID := neighbor.Workstream.ID
		if !visited[neighborID] {
			if hasCycleDFS(graph, neighborID, visited, recStack) {
				return true
			}
		} else if recStack[neighborID] {
			return true
		}
	}

	recStack[nodeID] = false
	return false
}

// executeSequential executes workstreams sequentially with retry
func (o *Orchestrator) executeSequential(workstreams []string) error {
	for _, wsID := range workstreams {
		err := o.executeWithRetry(wsID)
		if err != nil {
			return fmt.Errorf("%w: %s: %v", ErrExecutionFailed, wsID, err)
		}
	}
	return nil
}

// executeWithRetry executes a workstream with retry logic
func (o *Orchestrator) executeWithRetry(wsID string) error {
	var lastErr error
	for attempt := 0; attempt <= o.maxRetries; attempt++ {
		err := o.executor.Execute(wsID)
		if err == nil {
			return nil
		}
		lastErr = err
	}
	return fmt.Errorf("failed after %d retries: %w", o.maxRetries, lastErr)
}

// executeWithCheckpoint executes workstreams with checkpointing
func (o *Orchestrator) executeWithCheckpoint(workstreams []string, cp *checkpoint.Checkpoint) error {
	for _, wsID := range workstreams {
		// Update current workstream
		cp.CurrentWorkstream = wsID
		cp.Status = checkpoint.StatusInProgress

		// Save checkpoint before execution
		if err := o.checkpoint.Save(*cp); err != nil {
			return fmt.Errorf("failed to save checkpoint: %w", err)
		}

		// Execute workstream
		err := o.executeWithRetry(wsID)
		if err != nil {
			// Update checkpoint with failure
			cp.Status = checkpoint.StatusFailed
			o.checkpoint.Save(*cp)
			return err
		}

		// Mark as completed
		cp.CompletedWorkstreams = append(cp.CompletedWorkstreams, wsID)

		// Save checkpoint after successful execution
		if err := o.checkpoint.Save(*cp); err != nil {
			return fmt.Errorf("failed to save checkpoint: %w", err)
		}
	}

	// Mark checkpoint as completed
	cp.Status = checkpoint.StatusCompleted
	cp.CurrentWorkstream = ""
	if err := o.checkpoint.Save(*cp); err != nil {
		return fmt.Errorf("failed to save final checkpoint: %w", err)
	}

	return nil
}

// getRemainingWorkstreams determines which workstreams still need to execute
func (o *Orchestrator) getRemainingWorkstreams(allWorkstreams []string, cp *checkpoint.Checkpoint) ([]string, error) {
	// If already completed, nothing to do
	if cp.Status == checkpoint.StatusCompleted {
		return []string{}, nil
	}

	// Find position of current workstream
	startIndex := 0
	if cp.CurrentWorkstream != "" {
		for i, wsID := range allWorkstreams {
			if wsID == cp.CurrentWorkstream {
				startIndex = i
				break
			}
		}
	}

	// Return remaining workstreams
	return allWorkstreams[startIndex:], nil
}

// Run executes all workstreams for a feature
func (o *Orchestrator) Run(featureID string) error {
	// Load workstreams
	workstreams, err := o.loader.LoadWorkstreams(featureID)
	if err != nil {
		return fmt.Errorf("failed to load workstreams: %w", err)
	}

	if len(workstreams) == 0 {
		return fmt.Errorf("%w: feature %s has no workstreams", ErrFeatureNotFound, featureID)
	}

	// Build dependency graph
	graph, err := BuildDependencyGraph(workstreams)
	if err != nil {
		return fmt.Errorf("failed to build dependency graph: %w", err)
	}

	// Get execution order
	order, err := TopologicalSort(graph)
	if err != nil {
		return fmt.Errorf("failed to determine execution order: %w", err)
	}

	// Create initial checkpoint
	cp := &checkpoint.Checkpoint{
		ID:                    featureID,
		FeatureID:             featureID,
		Status:                checkpoint.StatusPending,
		CompletedWorkstreams:  []string{},
		CurrentWorkstream:     "",
	}

	// Execute with checkpointing
	return o.executeWithCheckpoint(order, cp)
}

// Resume resumes execution from a checkpoint
func (o *Orchestrator) Resume(checkpointID string) error {
	// Load checkpoint
	cp, err := o.checkpoint.Resume(checkpointID)
	if err != nil {
		return fmt.Errorf("failed to resume checkpoint: %w", err)
	}

	// If already completed, nothing to do
	if cp.Status == checkpoint.StatusCompleted {
		return nil
	}

	// Load workstreams
	workstreams, err := o.loader.LoadWorkstreams(cp.FeatureID)
	if err != nil {
		return fmt.Errorf("failed to load workstreams: %w", err)
	}

	// Build dependency graph
	graph, err := BuildDependencyGraph(workstreams)
	if err != nil {
		return fmt.Errorf("failed to build dependency graph: %w", err)
	}

	// Get full execution order
	fullOrder, err := TopologicalSort(graph)
	if err != nil {
		return fmt.Errorf("failed to determine execution order: %w", err)
	}

	// Get remaining workstreams
	remaining, err := o.getRemainingWorkstreams(fullOrder, &cp)
	if err != nil {
		return fmt.Errorf("failed to determine remaining workstreams: %w", err)
	}

	// Execute remaining
	return o.executeWithCheckpoint(remaining, &cp)
}
