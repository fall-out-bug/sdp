package dashboard

import (
	"time"
)

// DashboardState represents the current state of the dashboard
type DashboardState struct {
	// Current active tab (0-indexed)
	ActiveTab int

	// Workstreams grouped by status
	Workstreams map[string][]WorkstreamSummary

	// Ideas from docs/drafts/
	Ideas []IdeaSummary

	// Test coverage and results
	TestResults TestSummary

	// Last update timestamp
	LastUpdate time.Time

	// Loading state
	Loading bool

	// Error state
	Error error
}

// WorkstreamSummary represents a workstream for display
type WorkstreamSummary struct {
	ID       string
	Title    string
	Status   string // open, in_progress, completed, blocked
	Priority string // P0, P1, P2, P3
	Assignee string
	Size     string // SMALL, MEDIUM, LARGE
}

// IdeaSummary represents an idea draft for display
type IdeaSummary struct {
	Title string
	Path  string
	Date  time.Time
}

// TestSummary represents test results for display
type TestSummary struct {
	Coverage      string
	Passing       int
	Total         int
	LastRun       time.Time
	QualityGates  []GateStatus
}

// GateStatus represents a quality gate status
type GateStatus struct {
	Name   string
	Passed bool
}

// TabType represents different dashboard tabs
type TabType int

const (
	TabWorkstreams TabType = iota
	TabIdeas
	TabTests
	TabActivity
)

// String returns the tab name
func (t TabType) String() string {
	switch t {
	case TabWorkstreams:
		return "Workstreams"
	case TabIdeas:
		return "Ideas"
	case TabTests:
		return "Tests"
	case TabActivity:
		return "Activity"
	default:
		return "Unknown"
	}
}

// TickMsg is sent on a timer to trigger data refresh
type TickMsg time.Time

// RefreshMsg is sent to force a refresh
type RefreshMsg struct{}

// TabSelectMsg is sent when user switches tabs
type TabSelectMsg int
