package nextstep

import (
	"testing"
)

// TestInteractiveLoop tests the accept/refine/reject loop.
func TestInteractiveLoop(t *testing.T) {
	// Create a loop with a recommendation
	rec := &Recommendation{
		Command:    "sdp apply --ws 00-069-01",
		Reason:     "Ready to execute",
		Confidence: 0.9,
		Category:   CategoryExecution,
		Version:    ContractVersion,
	}

	loop := NewInteractiveLoop(rec)

	if loop == nil {
		t.Fatal("Expected non-nil loop")
	}
	if loop.CurrentRecommendation() == nil {
		t.Error("Expected current recommendation")
	}
}

// TestInteractiveLoopAccept tests accepting a recommendation.
func TestInteractiveLoopAccept(t *testing.T) {
	rec := &Recommendation{
		Command:    "sdp apply --ws 00-069-01",
		Reason:     "Ready to execute",
		Confidence: 0.9,
		Category:   CategoryExecution,
		Version:    ContractVersion,
	}

	loop := NewInteractiveLoop(rec)
	result := loop.Accept()

	if result.Action != ActionAccepted {
		t.Errorf("Expected ActionAccepted, got %v", result.Action)
	}
	if result.Command != rec.Command {
		t.Errorf("Expected command %s, got %s", rec.Command, result.Command)
	}
}

// TestInteractiveLoopReject tests rejecting a recommendation.
func TestInteractiveLoopReject(t *testing.T) {
	rec := &Recommendation{
		Command:    "sdp apply --ws 00-069-01",
		Reason:     "Ready to execute",
		Confidence: 0.9,
		Category:   CategoryExecution,
		Version:    ContractVersion,
		Alternatives: []Alternative{
			{Command: "sdp status", Reason: "Check state first"},
		},
	}

	loop := NewInteractiveLoop(rec)

	// Reject should move to next alternative
	result := loop.Reject()

	if result.Action != ActionAlternative {
		t.Errorf("Expected ActionAlternative, got %v", result.Action)
	}
	if result.Command != "sdp status" {
		t.Errorf("Expected alternative command, got %s", result.Command)
	}
}

// TestInteractiveLoopRejectAll tests rejecting all alternatives.
func TestInteractiveLoopRejectAll(t *testing.T) {
	rec := &Recommendation{
		Command:    "sdp apply --ws 00-069-01",
		Reason:     "Ready to execute",
		Confidence: 0.9,
		Category:   CategoryExecution,
		Alternatives: []Alternative{
			{Command: "sdp status", Reason: "Check state"},
		},
	}

	loop := NewInteractiveLoop(rec)

	// Reject primary
	loop.Reject()
	// Reject alternative
	result := loop.Reject()

	if result.Action != ActionRejected {
		t.Errorf("Expected ActionRejected after rejecting all, got %v", result.Action)
	}
}

// TestInteractiveLoopRefine tests refining a recommendation.
func TestInteractiveLoopRefine(t *testing.T) {
	rec := &Recommendation{
		Command:    "sdp apply --ws 00-069-01",
		Reason:     "Ready to execute",
		Confidence: 0.9,
		Category:   CategoryExecution,
		Version:    ContractVersion,
	}

	loop := NewInteractiveLoop(rec)

	// Refine with user input
	result := loop.Refine("use --dry-run")

	if result.Action != ActionRefined {
		t.Errorf("Expected ActionRefined, got %v", result.Action)
	}
	// Refined command should include the user's modification
	if result.Command == "" {
		t.Error("Expected non-empty refined command")
	}
}

// TestInteractiveLoopContext tests context preservation.
func TestInteractiveLoopContext(t *testing.T) {
	rec := &Recommendation{
		Command:    "sdp apply --ws 00-069-01",
		Reason:     "Ready to execute",
		Confidence: 0.9,
		Category:   CategoryExecution,
		Version:    ContractVersion,
		Metadata: map[string]any{
			"workstream_id": "00-069-01",
			"feature_id":    "F069",
		},
	}

	loop := NewInteractiveLoop(rec)

	// Context should be preserved
	ctx := loop.Context()
	if ctx["workstream_id"] != "00-069-01" {
		t.Error("Expected workstream_id to be preserved in context")
	}
}

// TestInteractiveLoopWhy tests the "why" explanation.
func TestInteractiveLoopWhy(t *testing.T) {
	rec := &Recommendation{
		Command:    "sdp apply --ws 00-069-01",
		Reason:     "Ready to execute",
		Confidence: 0.9,
		Category:   CategoryExecution,
		Version:    ContractVersion,
		Metadata: map[string]any{
			"priority":      0,
			"workstream_id": "00-069-01",
		},
	}

	loop := NewInteractiveLoop(rec)
	explanation := loop.Why()

	if explanation == "" {
		t.Error("Expected non-empty explanation")
	}
}

// TestInteractiveLoopState tests loop state tracking.
func TestInteractiveLoopState(t *testing.T) {
	rec := &Recommendation{
		Command:    "sdp apply --ws 00-069-01",
		Reason:     "Ready to execute",
		Confidence: 0.9,
		Category:   CategoryExecution,
		Alternatives: []Alternative{
			{Command: "sdp status", Reason: "Check state"},
		},
	}

	loop := NewInteractiveLoop(rec)

	// Initial state
	if loop.CurrentIndex() != 0 {
		t.Error("Expected initial index 0")
	}

	// After reject, should move to alternative
	loop.Reject()
	if loop.CurrentIndex() != 1 {
		t.Errorf("Expected index 1 after one reject, got %d", loop.CurrentIndex())
	}
}

// TestInteractiveLoopExit tests safe exit and resume.
func TestInteractiveLoopExit(t *testing.T) {
	rec := &Recommendation{
		Command:    "sdp apply --ws 00-069-01",
		Reason:     "Ready to execute",
		Confidence: 0.9,
		Category:   CategoryExecution,
		Version:    ContractVersion,
	}

	loop := NewInteractiveLoop(rec)

	// Create checkpoint for resume
	checkpoint := loop.CreateCheckpoint()

	if checkpoint == nil {
		t.Error("Expected non-nil checkpoint")
	}

	// Resume from checkpoint
	loop2 := ResumeFromCheckpoint(checkpoint)

	if loop2.CurrentRecommendation().Command != rec.Command {
		t.Error("Expected same command after resume")
	}
}
