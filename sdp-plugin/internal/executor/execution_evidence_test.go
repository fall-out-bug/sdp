package executor

import (
	"bytes"
	"context"
	"testing"
)

// TestExecutor_EvidenceChain tests AC8: emits full evidence chain
func TestExecutor_EvidenceChain(t *testing.T) {
	exec := NewExecutor(ExecutorConfig{
		BacklogDir:      "testdata/backlog",
		DryRun:          false,
		RetryCount:      1,
		EvidenceLogPath: "testdata/evidence.jsonl",
	})

	ctx := context.Background()
	var output bytes.Buffer

	result, err := exec.Execute(ctx, &output, ExecuteOptions{
		All:        true,
		SpecificWS: "",
		Retry:      0,
		Output:     "human",
	})

	if err != nil {
		t.Fatalf("Execute() failed: %v", err)
	}

	if len(result.EvidenceEvents) == 0 {
		t.Error("Expected evidence events to be emitted")
	}

	// Check for required evidence types
	hasPlanEvent := false
	hasGenEvent := false
	hasVerifyEvent := false
	hasApprovalEvent := false

	for _, event := range result.EvidenceEvents {
		switch event.Type {
		case "plan":
			hasPlanEvent = true
		case "generation":
			hasGenEvent = true
		case "verification":
			hasVerifyEvent = true
		case "approval":
			hasApprovalEvent = true
		}
	}

	if !hasPlanEvent {
		t.Error("Expected 'plan' evidence event")
	}
	if !hasGenEvent {
		t.Error("Expected 'generation' evidence event")
	}
	if !hasVerifyEvent {
		t.Error("Expected 'verification' evidence event")
	}
	if !hasApprovalEvent {
		t.Error("Expected 'approval' evidence event")
	}
}
