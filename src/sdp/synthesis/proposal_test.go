package synthesis

import (
	"encoding/json"
	"testing"
	"time"
)

// TestNewProposal verifies proposal creation
func TestNewProposal(t *testing.T) {
	proposal := NewProposal("agent-1", "solution", 0.95, "test reasoning")

	if proposal.AgentID != "agent-1" {
		t.Error("AgentID not set")
	}
	if proposal.Solution != "solution" {
		t.Error("Solution not set")
	}
	if proposal.Confidence != 0.95 {
		t.Error("Confidence not set")
	}
	if proposal.Reasoning != "test reasoning" {
		t.Error("Reasoning not set")
	}
	if proposal.Timestamp.IsZero() {
		t.Error("Timestamp not set")
	}
}

// TestProposal_Equals_Same verifies equality for same proposals
func TestProposal_Equals_Same(t *testing.T) {
	p1 := &Proposal{
		AgentID:    "agent-1",
		Solution:   "solution",
		Confidence: 0.95,
		Reasoning:  "reasoning",
	}
	p2 := &Proposal{
		AgentID:    "agent-1",
		Solution:   "solution",
		Confidence: 0.95,
		Reasoning:  "reasoning",
	}

	if !p1.Equals(p2) {
		t.Error("proposals should be equal")
	}
}

// TestProposal_Equals_Different verifies inequality for different proposals
func TestProposal_Equals_Different(t *testing.T) {
	p1 := &Proposal{AgentID: "agent-1"}
	p2 := &Proposal{AgentID: "agent-2"}

	if p1.Equals(p2) {
		t.Error("proposals should not be equal")
	}
}

// TestProposal_Equals_Nil verifies nil handling
func TestProposal_Equals_Nil(t *testing.T) {
	var p1 *Proposal
	var p2 *Proposal

	if !p1.Equals(p2) {
		t.Error("nil proposals should be equal")
	}

	p1 = &Proposal{AgentID: "agent-1"}
	if p1.Equals(nil) {
		t.Error("non-nil should not equal nil")
	}
}

// TestProposal_Equals_DifferentConfidence verifies confidence comparison
func TestProposal_Equals_DifferentConfidence(t *testing.T) {
	p1 := &Proposal{AgentID: "agent-1", Confidence: 0.9}
	p2 := &Proposal{AgentID: "agent-1", Confidence: 0.8}

	if p1.Equals(p2) {
		t.Error("different confidence should not be equal")
	}
}

// TestProposal_Equals_DifferentReasoning verifies reasoning comparison
func TestProposal_Equals_DifferentReasoning(t *testing.T) {
	p1 := &Proposal{AgentID: "agent-1", Reasoning: "reason 1"}
	p2 := &Proposal{AgentID: "agent-1", Reasoning: "reason 2"}

	if p1.Equals(p2) {
		t.Error("different reasoning should not be equal")
	}
}

// TestProposal_Equals_ComplexSolution verifies complex solution comparison
func TestProposal_Equals_ComplexSolution(t *testing.T) {
	p1 := &Proposal{
		AgentID:  "agent-1",
		Solution: map[string]int{"a": 1, "b": 2},
	}
	p2 := &Proposal{
		AgentID:  "agent-1",
		Solution: map[string]int{"a": 1, "b": 2},
	}
	p3 := &Proposal{
		AgentID:  "agent-1",
		Solution: map[string]int{"a": 1, "b": 3},
	}

	if !p1.Equals(p2) {
		t.Error("same complex solutions should be equal")
	}
	if p1.Equals(p3) {
		t.Error("different complex solutions should not be equal")
	}
}

// TestProposal_IsHigherConfidenceThan verifies confidence comparison
func TestProposal_IsHigherConfidenceThan(t *testing.T) {
	p1 := &Proposal{Confidence: 0.9}
	p2 := &Proposal{Confidence: 0.8}

	if !p1.IsHigherConfidenceThan(p2) {
		t.Error("p1 should have higher confidence")
	}
	if p2.IsHigherConfidenceThan(p1) {
		t.Error("p2 should not have higher confidence")
	}
}

// TestProposal_IsHigherConfidenceThan_Nil verifies nil handling
func TestProposal_IsHigherConfidenceThan_Nil(t *testing.T) {
	p1 := &Proposal{Confidence: 0.9}

	if p1.IsHigherConfidenceThan(nil) {
		t.Error("should return false for nil")
	}

	var p2 *Proposal
	if p2.IsHigherConfidenceThan(p1) {
		t.Error("nil should return false")
	}
}

// TestProposal_MarshalJSON verifies JSON marshaling
func TestProposal_MarshalJSON(t *testing.T) {
	ts := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	p := &Proposal{
		AgentID:    "agent-1",
		Solution:   "solution",
		Confidence: 0.95,
		Reasoning:  "test",
		Timestamp:  ts,
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	if string(data) == "" {
		t.Error("marshal should produce output")
	}
}

// TestProposal_UnmarshalJSON verifies JSON unmarshaling
func TestProposal_UnmarshalJSON(t *testing.T) {
	jsonStr := `{"agent_id":"agent-1","solution":"test","confidence":0.9,"reasoning":"reason","timestamp":"2024-01-01T12:00:00Z"}`

	var p Proposal
	err := json.Unmarshal([]byte(jsonStr), &p)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if p.AgentID != "agent-1" {
		t.Error("AgentID not unmarshaled")
	}
	if p.Confidence != 0.9 {
		t.Error("Confidence not unmarshaled")
	}
}

// TestProposal_UnmarshalJSON_InvalidTimestamp verifies invalid timestamp handling
func TestProposal_UnmarshalJSON_InvalidTimestamp(t *testing.T) {
	jsonStr := `{"agent_id":"agent-1","timestamp":"invalid"}`

	var p Proposal
	err := json.Unmarshal([]byte(jsonStr), &p)
	if err == nil {
		t.Error("expected error for invalid timestamp")
	}
}

// TestProposal_UnmarshalJSON_EmptyTimestamp verifies empty timestamp handling
func TestProposal_UnmarshalJSON_EmptyTimestamp(t *testing.T) {
	jsonStr := `{"agent_id":"agent-1","timestamp":""}`

	var p Proposal
	err := json.Unmarshal([]byte(jsonStr), &p)
	if err != nil {
		t.Errorf("empty timestamp should not error: %v", err)
	}
}
