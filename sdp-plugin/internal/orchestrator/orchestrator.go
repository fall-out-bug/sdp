package orchestrator

import (
	"errors"

	"github.com/fall-out-bug/sdp/internal/checkpoint"
)

// Common errors
var (
	ErrFeatureNotFound       = errors.New("feature not found")
	ErrCheckpointNotFound    = errors.New("checkpoint not found")
	ErrExecutionFailed       = errors.New("workstream execution failed")
	ErrCircularDependency    = errors.New("circular dependency detected")
	ErrMissingDependency     = errors.New("missing dependency")
	ErrWorkstreamNotFound    = errors.New("workstream not found")
	ErrSkillInvocationFailed = errors.New("skill invocation failed")
	ErrAgentSpawnFailed      = errors.New("agent spawn failed")
	ErrAgentNotFound         = errors.New("agent not found")
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
