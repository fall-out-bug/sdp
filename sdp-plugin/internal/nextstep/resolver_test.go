package nextstep

import (
	"testing"
)

// TestResolverFreshProject tests AC2: Fresh project state.
func TestResolverFreshProject(t *testing.T) {
	resolver := NewResolver()
	state := ProjectState{
		Workstreams: []WorkstreamStatus{},
		Mode:        ModeDrive,
		GitStatus:   GitStatusInfo{IsRepo: true},
		Config:      ConfigInfo{HasSDPConfig: true},
	}

	rec, err := resolver.Recommend(state)
	if err != nil {
		t.Fatalf("Recommend() error: %v", err)
	}

	if err := rec.Validate(); err != nil {
		t.Errorf("Invalid recommendation: %v", err)
	}
	if rec.Category != CategorySetup && rec.Category != CategoryInformation {
		t.Errorf("Expected setup or information category for fresh project, got %s", rec.Category)
	}
	if rec.Confidence < 0.5 {
		t.Errorf("Expected confidence >= 0.5 for fresh project, got %.2f", rec.Confidence)
	}
}

// TestResolverInProgressWorkstream tests AC2: In-progress WS state.
func TestResolverInProgressWorkstream(t *testing.T) {
	resolver := NewResolver()
	state := ProjectState{
		Workstreams: []WorkstreamStatus{
			{ID: "00-069-01", Status: StatusInProgress, Priority: 0, Feature: "F069"},
		},
		ActiveWorkstream: "00-069-01",
		Mode:             ModeDrive,
		GitStatus:        GitStatusInfo{IsRepo: true},
	}

	rec, err := resolver.Recommend(state)
	if err != nil {
		t.Fatalf("Recommend() error: %v", err)
	}

	if err := rec.Validate(); err != nil {
		t.Errorf("Invalid recommendation: %v", err)
	}
	// Should recommend continuing or checking status
	if rec.Category != CategoryInformation && rec.Category != CategoryExecution {
		t.Errorf("Expected information or execution category, got %s", rec.Category)
	}
}

// TestResolverReadyWorkstream tests AC2: Ready workstream state.
func TestResolverReadyWorkstream(t *testing.T) {
	resolver := NewResolver()
	state := ProjectState{
		Workstreams: []WorkstreamStatus{
			{ID: "00-069-01", Status: StatusReady, Priority: 0, Feature: "F069", BlockedBy: nil},
		},
		Mode:      ModeDrive,
		GitStatus: GitStatusInfo{IsRepo: true},
	}

	rec, err := resolver.Recommend(state)
	if err != nil {
		t.Fatalf("Recommend() error: %v", err)
	}

	if err := rec.Validate(); err != nil {
		t.Errorf("Invalid recommendation: %v", err)
	}
	// Should recommend executing the ready workstream
	if rec.Category != CategoryExecution {
		t.Errorf("Expected execution category, got %s", rec.Category)
	}
	if rec.Confidence < 0.8 {
		t.Errorf("Expected high confidence for ready workstream, got %.2f", rec.Confidence)
	}
}

// TestResolverBlockedWorkstream tests AC2: Blocked WS state.
func TestResolverBlockedWorkstream(t *testing.T) {
	resolver := NewResolver()
	state := ProjectState{
		Workstreams: []WorkstreamStatus{
			{ID: "00-069-01", Status: StatusReady, Priority: 0, Feature: "F069", BlockedBy: nil},
			{ID: "00-069-02", Status: StatusReady, Priority: 0, Feature: "F069", BlockedBy: []string{"00-069-01"}},
		},
		Mode:      ModeDrive,
		GitStatus: GitStatusInfo{IsRepo: true},
	}

	rec, err := resolver.Recommend(state)
	if err != nil {
		t.Fatalf("Recommend() error: %v", err)
	}

	// Should recommend executing the unblocked workstream first
	if rec.Category != CategoryExecution {
		t.Errorf("Expected execution category, got %s", rec.Category)
	}
}

// TestResolverFailedWorkstream tests AC2: Failed WS state.
func TestResolverFailedWorkstream(t *testing.T) {
	resolver := NewResolver()
	state := ProjectState{
		Workstreams: []WorkstreamStatus{
			{
				ID:        "00-069-01",
				Status:    StatusFailed,
				Priority:  0,
				Feature:   "F069",
				LastError: "test failure",
			},
		},
		LastCommand:      "sdp apply --ws 00-069-01",
		LastCommandError: "test failure",
		ActiveWorkstream: "00-069-01",
		Mode:             ModeDrive,
		GitStatus:        GitStatusInfo{IsRepo: true},
	}

	rec, err := resolver.Recommend(state)
	if err != nil {
		t.Fatalf("Recommend() error: %v", err)
	}

	// Should recommend recovery
	if rec.Category != CategoryRecovery {
		t.Errorf("Expected recovery category for failed workstream, got %s", rec.Category)
	}
}

// TestResolverDeterministic tests AC3: Deterministic fallback.
func TestResolverDeterministic(t *testing.T) {
	resolver := NewResolver()
	state := ProjectState{
		Workstreams: []WorkstreamStatus{
			{ID: "00-069-02", Status: StatusReady, Priority: 0, Feature: "F069", BlockedBy: nil},
			{ID: "00-069-01", Status: StatusReady, Priority: 0, Feature: "F069", BlockedBy: nil},
		},
		Mode:      ModeDrive,
		GitStatus: GitStatusInfo{IsRepo: true},
	}

	// Run multiple times - should always get the same recommendation
	var lastCmd string
	for i := 0; i < 10; i++ {
		rec, err := resolver.Recommend(state)
		if err != nil {
			t.Fatalf("Recommend() error: %v", err)
		}
		if lastCmd != "" && rec.Command != lastCmd {
			t.Errorf("Non-deterministic output: got %q, want %q", rec.Command, lastCmd)
		}
		lastCmd = rec.Command
	}

	// Should prefer 00-069-01 (lower ID)
	if lastCmd != "sdp apply --ws 00-069-01" {
		t.Errorf("Expected recommendation for 00-069-01, got %s", lastCmd)
	}
}

// TestResolverPriorityOrdering tests AC4: Priority ordering.
func TestResolverPriorityOrdering(t *testing.T) {
	resolver := NewResolver()
	state := ProjectState{
		Workstreams: []WorkstreamStatus{
			{ID: "00-069-01", Status: StatusReady, Priority: 2, Feature: "F069", BlockedBy: nil},
			{ID: "00-069-02", Status: StatusReady, Priority: 0, Feature: "F069", BlockedBy: nil},
			{ID: "00-069-03", Status: StatusReady, Priority: 1, Feature: "F069", BlockedBy: nil},
		},
		Mode:      ModeDrive,
		GitStatus: GitStatusInfo{IsRepo: true},
	}

	rec, err := resolver.Recommend(state)
	if err != nil {
		t.Fatalf("Recommend() error: %v", err)
	}

	// Should recommend 00-069-02 (priority 0)
	if rec.Command != "sdp apply --ws 00-069-02" {
		t.Errorf("Expected recommendation for 00-069-02 (priority 0), got %s", rec.Command)
	}
}

// TestResolverAllComplete tests when all workstreams are done.
func TestResolverAllComplete(t *testing.T) {
	resolver := NewResolver()
	state := ProjectState{
		Workstreams: []WorkstreamStatus{
			{ID: "00-069-01", Status: StatusCompleted, Priority: 0, Feature: "F069"},
			{ID: "00-069-02", Status: StatusCompleted, Priority: 0, Feature: "F069"},
		},
		Mode:      ModeDrive,
		GitStatus: GitStatusInfo{IsRepo: true, Branch: "feature/F069"},
	}

	rec, err := resolver.Recommend(state)
	if err != nil {
		t.Fatalf("Recommend() error: %v", err)
	}

	// Should recommend review or deploy
	if rec.Category != CategoryPlanning && rec.Category != CategoryInformation {
		t.Errorf("Expected planning or information for completed feature, got %s", rec.Category)
	}
}

// TestResolverNoGitRepo tests when not in a git repo.
func TestResolverNoGitRepo(t *testing.T) {
	resolver := NewResolver()
	state := ProjectState{
		Workstreams: []WorkstreamStatus{},
		Mode:        ModeDrive,
		GitStatus:   GitStatusInfo{IsRepo: false},
	}

	rec, err := resolver.Recommend(state)
	if err != nil {
		t.Fatalf("Recommend() error: %v", err)
	}

	// Should recommend initializing git
	if rec.Category != CategorySetup {
		t.Errorf("Expected setup category for non-git project, got %s", rec.Category)
	}
}

// TestResolverUncommittedChanges tests when there are uncommitted changes.
func TestResolverUncommittedChanges(t *testing.T) {
	resolver := NewResolver()
	state := ProjectState{
		Workstreams: []WorkstreamStatus{
			{ID: "00-069-01", Status: StatusReady, Priority: 0, Feature: "F069", BlockedBy: nil},
		},
		Mode: ModeDrive,
		GitStatus: GitStatusInfo{
			IsRepo:      true,
			Uncommitted: true,
			Branch:      "feature/F069",
		},
	}

	rec, err := resolver.Recommend(state)
	if err != nil {
		t.Fatalf("Recommend() error: %v", err)
	}

	// Should recommend committing first
	if rec.Category != CategoryRecovery && rec.Category != CategoryInformation {
		t.Errorf("Expected recovery or information for uncommitted changes, got %s", rec.Category)
	}
}

// TestResolverMachineReadableOutput tests AC5: Machine-readable output.
func TestResolverMachineReadableOutput(t *testing.T) {
	resolver := NewResolver()
	state := ProjectState{
		Workstreams: []WorkstreamStatus{
			{ID: "00-069-01", Status: StatusReady, Priority: 0, Feature: "F069", BlockedBy: nil},
		},
		Mode:      ModeDrive,
		GitStatus: GitStatusInfo{IsRepo: true},
	}

	rec, err := resolver.Recommend(state)
	if err != nil {
		t.Fatalf("Recommend() error: %v", err)
	}

	// Verify JSON serialization works
	jsonData, err := rec.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error: %v", err)
	}
	if len(jsonData) == 0 {
		t.Error("ToJSON() returned empty data")
	}

	// Verify we can parse it back
	parsed, err := FromJSON(jsonData)
	if err != nil {
		t.Fatalf("FromJSON() error: %v", err)
	}
	if parsed.Command != rec.Command {
		t.Errorf("Parsed command mismatch: got %q, want %q", parsed.Command, rec.Command)
	}
}

// TestResolverLowConfidenceFallback tests AC3: Low confidence fallback.
func TestResolverLowConfidenceFallback(t *testing.T) {
	resolver := NewResolver()
	// Ambiguous state with multiple signals
	state := ProjectState{
		Workstreams: []WorkstreamStatus{
			{ID: "00-069-01", Status: StatusBacklog, Priority: 0, Feature: "F069", BlockedBy: []string{"00-068-99"}},
			{ID: "00-069-02", Status: StatusBacklog, Priority: 0, Feature: "F069", BlockedBy: []string{"00-068-99"}},
		},
		Mode:      ModeManual,
		GitStatus: GitStatusInfo{IsRepo: true},
	}

	rec, err := resolver.Recommend(state)
	if err != nil {
		t.Fatalf("Recommend() error: %v", err)
	}

	// Should still provide a valid recommendation with alternatives
	if err := rec.Validate(); err != nil {
		t.Errorf("Invalid recommendation: %v", err)
	}
	// When nothing is ready, should have alternatives
	if len(rec.Alternatives) == 0 && rec.Confidence < 0.7 {
		t.Error("Low confidence recommendation should have alternatives")
	}
}
