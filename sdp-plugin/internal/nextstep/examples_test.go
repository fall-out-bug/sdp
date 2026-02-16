package nextstep

import (
	"encoding/json"
	"testing"
)

// TestExampleFreshProject tests recommendation for a fresh project state.
// This is a documented example for the contract.
func TestExampleFreshProject(t *testing.T) {
	_ = ProjectState{
		Workstreams: []WorkstreamStatus{},
		LastCommand: "",
		Mode:        ModeDrive,
		GitStatus: GitStatusInfo{
			IsRepo:     true,
			Branch:     "main",
			MainBranch: "main",
		},
		Config: ConfigInfo{
			HasSDPConfig: true,
			Version:      "0.10.0",
		},
	}

	// A fresh project should recommend starting with vision/reality
	rec := Recommendation{
		Command:    "sdp doctor",
		Reason:     "Start by checking your environment setup",
		Confidence: 0.9,
		Category:   CategorySetup,
		Version:    ContractVersion,
		Alternatives: []Alternative{
			{Command: "sdp status", Reason: "View current project state"},
		},
		Metadata: map[string]any{
			"state_type": "fresh_project",
		},
	}

	if err := rec.Validate(); err != nil {
		t.Errorf("Fresh project recommendation invalid: %v", err)
	}

	// Verify JSON output
	data, _ := json.MarshalIndent(rec, "", "  ")
	t.Logf("Fresh project recommendation JSON:\n%s", string(data))
}

// TestExampleInProgressWorkstream tests recommendation for in-progress work.
func TestExampleInProgressWorkstream(t *testing.T) {
	_ = ProjectState{
		Workstreams: []WorkstreamStatus{
			{ID: "00-069-01", Status: StatusInProgress, Priority: 0, Feature: "F069"},
			{ID: "00-069-02", Status: StatusReady, Priority: 0, BlockedBy: []string{"00-069-01"}, Feature: "F069"},
		},
		LastCommand:      "sdp apply --ws 00-069-01",
		ActiveWorkstream: "00-069-01",
		Mode:             ModeDrive,
		GitStatus: GitStatusInfo{
			IsRepo:       true,
			Branch:       "feature/F069-next-step",
			Uncommitted:  false,
			MainBranch:   "main",
		},
		Config: ConfigInfo{
			HasSDPConfig:    true,
			Version:         "0.10.0",
			EvidenceEnabled: true,
		},
	}

	// When a workstream is in progress, recommend continuing or checking status
	rec := Recommendation{
		Command:    "sdp status",
		Reason:     "Workstream 00-069-01 is in progress - check current state",
		Confidence: 0.85,
		Category:   CategoryInformation,
		Version:    ContractVersion,
		Metadata: map[string]any{
			"active_workstream": "00-069-01",
			"feature_id":        "F069",
		},
	}

	if err := rec.Validate(); err != nil {
		t.Errorf("In-progress recommendation invalid: %v", err)
	}
}

// TestExampleBlockedWorkstream tests recommendation for blocked state.
func TestExampleBlockedWorkstream(t *testing.T) {
	_ = ProjectState{
		Workstreams: []WorkstreamStatus{
			{ID: "00-069-01", Status: StatusReady, Priority: 0, Feature: "F069"},
			{ID: "00-069-02", Status: StatusBlocked, Priority: 0, BlockedBy: []string{"00-069-01"}, Feature: "F069"},
		},
		Mode: ModeDrive,
		GitStatus: GitStatusInfo{
			IsRepo:     true,
			Branch:     "feature/F069-next-step",
			MainBranch: "main",
		},
	}

	// When next workstream is blocked, recommend completing blocker first
	rec := Recommendation{
		Command:    "sdp apply --ws 00-069-01",
		Reason:     "Complete blocking workstream first to unblock 00-069-02",
		Confidence: 0.95,
		Category:   CategoryExecution,
		Version:    ContractVersion,
		Metadata: map[string]any{
			"blocking":         "00-069-02",
			"blocker":          "00-069-01",
			"dependency_chain": []string{"00-069-01", "00-069-02"},
		},
	}

	if err := rec.Validate(); err != nil {
		t.Errorf("Blocked recommendation invalid: %v", err)
	}
}

// TestExampleFailedWorkstream tests recommendation after failure.
func TestExampleFailedWorkstream(t *testing.T) {
	_ = ProjectState{
		Workstreams: []WorkstreamStatus{
			{
				ID:        "00-069-01",
				Status:    StatusFailed,
				Priority:  0,
				Feature:   "F069",
				LastError: "test failure: expected X got Y",
			},
		},
		LastCommand:      "sdp apply --ws 00-069-01",
		LastCommandError: "test failure: expected X got Y",
		Mode:             ModeDrive,
		ActiveWorkstream: "00-069-01",
	}

	// After failure, recommend recovery actions
	rec := Recommendation{
		Command:    "sdp debug --ws 00-069-01",
		Reason:     "Workstream failed with test error - use debug to investigate",
		Confidence: 0.9,
		Category:   CategoryRecovery,
		Version:    ContractVersion,
		Alternatives: []Alternative{
			{Command: "sdp retry --ws 00-069-01", Reason: "Retry the workstream"},
			{Command: "sdp status", Reason: "Check current state before debugging"},
		},
		Metadata: map[string]any{
			"failed_workstream": "00-069-01",
			"error_type":        "test_failure",
		},
	}

	if err := rec.Validate(); err != nil {
		t.Errorf("Failed recommendation invalid: %v", err)
	}
}

// TestExampleAllWorkstreamsComplete tests recommendation when all WS done.
func TestExampleAllWorkstreamsComplete(t *testing.T) {
	_ = ProjectState{
		Workstreams: []WorkstreamStatus{
			{ID: "00-069-01", Status: StatusCompleted, Priority: 0, Feature: "F069"},
			{ID: "00-069-02", Status: StatusCompleted, Priority: 0, Feature: "F069"},
			{ID: "00-069-03", Status: StatusCompleted, Priority: 0, Feature: "F069"},
			{ID: "00-069-04", Status: StatusCompleted, Priority: 0, Feature: "F069"},
			{ID: "00-069-05", Status: StatusCompleted, Priority: 0, Feature: "F069"},
		},
		Mode: ModeDrive,
		GitStatus: GitStatusInfo{
			IsRepo:      true,
			Branch:      "feature/F069-next-step",
			Uncommitted: false,
			MainBranch:  "main",
		},
	}

	// When all workstreams complete, recommend review and deploy
	rec := Recommendation{
		Command:    "sdp review F069",
		Reason:     "All workstreams complete - review feature before deployment",
		Confidence: 0.95,
		Category:   CategoryPlanning,
		Version:    ContractVersion,
		Alternatives: []Alternative{
			{Command: "sdp deploy F069", Reason: "Deploy if review passed"},
			{Command: "sdp status", Reason: "Check overall project state"},
		},
		Metadata: map[string]any{
			"feature_id":      "F069",
			"completed_count": 5,
			"total_count":     5,
		},
	}

	if err := rec.Validate(); err != nil {
		t.Errorf("Complete recommendation invalid: %v", err)
	}
}

// TestContractBackwardCompatibility ensures old versions can still be parsed.
func TestContractBackwardCompatibility(t *testing.T) {
	// Simulate a v1.0.0 recommendation
	jsonV1 := `{
		"command": "sdp apply --ws 00-069-01",
		"reason": "Test backward compat",
		"confidence": 0.8,
		"category": "execution",
		"version": "1.0.0"
	}`

	var rec Recommendation
	if err := json.Unmarshal([]byte(jsonV1), &rec); err != nil {
		t.Fatalf("Failed to parse v1.0.0 recommendation: %v", err)
	}

	if rec.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", rec.Version)
	}
	if rec.Command != "sdp apply --ws 00-069-01" {
		t.Errorf("Unexpected command: %s", rec.Command)
	}
}
