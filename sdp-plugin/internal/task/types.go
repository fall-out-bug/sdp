package task

import "time"

// Type represents the type of task
type Type string

const (
	TypeBug    Type = "bug"
	TypeTask   Type = "task"
	TypeHotfix Type = "hotfix"
)

// Priority represents task priority (P0=0, P1=1, P2=2, P3=3)
type Priority int

const (
	PriorityP0 Priority = 0 // Critical/Production down
	PriorityP1 Priority = 1 // High
	PriorityP2 Priority = 2 // Medium
	PriorityP3 Priority = 3 // Low
)

// Status represents task status
type Status string

const (
	StatusBacklog   Status = "backlog"
	StatusActive    Status = "active"
	StatusBlocked   Status = "blocked"
	StatusCompleted Status = "completed"
)

// Task represents a task to be created
type Task struct {
	Type       Type
	Title      string
	Priority   Priority
	FeatureID  string // Optional: parent feature
	BranchBase string // "dev" (default) or "main" (hotfix)
	Goal       string
	Context    string
	ScopeFiles []string
	BeadsID    string // Optional: linked beads ID
	DependsOn  []string
}

// Workstream represents the created workstream artifact
type Workstream struct {
	WSID      string
	Path      string
	FeatureID string
	BeadsID   string
	CreatedAt time.Time
}

// Issue represents the created issue artifact (Beads-free fallback)
type Issue struct {
	IssueID   string
	Path      string
	CreatedAt time.Time
}

// CreatorConfig configures the task creator
type CreatorConfig struct {
	WorkstreamDir string
	IssuesDir     string
	IndexFile     string
	ProjectID     string
	BeadsEnabled  bool
}

// Creator creates tasks (bugs, tasks, hotfixes) with dual-track artifact creation
type Creator struct {
	config CreatorConfig
}

// NewCreator creates a new task creator
func NewCreator(config CreatorConfig) *Creator {
	if config.WorkstreamDir == "" {
		config.WorkstreamDir = "docs/workstreams/backlog"
	}
	if config.IssuesDir == "" {
		config.IssuesDir = "docs/issues"
	}
	if config.IndexFile == "" {
		config.IndexFile = ".sdp/issues-index.jsonl"
	}
	if config.ProjectID == "" {
		config.ProjectID = "00"
	}

	return &Creator{config: config}
}
