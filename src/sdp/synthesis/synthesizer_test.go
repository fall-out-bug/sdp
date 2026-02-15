package synthesis

import (
	"errors"
	"testing"
)

// TestNewSynthesizer verifies synthesizer creation
func TestNewSynthesizer(t *testing.T) {
	s := NewSynthesizer()
	if s == nil {
		t.Fatal("NewSynthesizer returned nil")
	}
	if s.proposals == nil {
		t.Error("proposals map not initialized")
	}
}

// TestSynthesizer_AddProposal verifies adding proposals
func TestSynthesizer_AddProposal(t *testing.T) {
	s := NewSynthesizer()
	p := NewProposal("agent-1", "solution", 0.9, "test")

	s.AddProposal(p)

	if len(s.proposals) != 1 {
		t.Error("proposal not added")
	}
}

// TestSynthesizer_AddProposal_Nil verifies nil handling
func TestSynthesizer_AddProposal_Nil(t *testing.T) {
	s := NewSynthesizer()
	s.AddProposal(nil)

	if len(s.proposals) != 0 {
		t.Error("nil proposal should not be added")
	}
}

// TestSynthesizer_AddProposal_Overwrite verifies overwrite by agent ID
func TestSynthesizer_AddProposal_Overwrite(t *testing.T) {
	s := NewSynthesizer()
	s.AddProposal(NewProposal("agent-1", "solution1", 0.9, ""))
	s.AddProposal(NewProposal("agent-1", "solution2", 0.8, ""))

	if len(s.proposals) != 1 {
		t.Error("should have 1 proposal (overwritten)")
	}
	if s.proposals["agent-1"].Solution != "solution2" {
		t.Error("proposal should be overwritten")
	}
}

// TestSynthesizer_GetProposals verifies getting proposals
func TestSynthesizer_GetProposals(t *testing.T) {
	s := NewSynthesizer()
	s.AddProposal(NewProposal("agent-1", "s1", 0.9, ""))
	s.AddProposal(NewProposal("agent-2", "s2", 0.8, ""))

	proposals := s.GetProposals()
	if len(proposals) != 2 {
		t.Errorf("expected 2 proposals, got %d", len(proposals))
	}
}

// TestSynthesizer_GetProposals_Empty verifies empty case
func TestSynthesizer_GetProposals_Empty(t *testing.T) {
	s := NewSynthesizer()
	proposals := s.GetProposals()
	if len(proposals) != 0 {
		t.Error("expected empty slice")
	}
}

// TestSynthesizer_Clear verifies clearing proposals
func TestSynthesizer_Clear(t *testing.T) {
	s := NewSynthesizer()
	s.AddProposal(NewProposal("agent-1", "s1", 0.9, ""))
	s.Clear()

	if len(s.proposals) != 0 {
		t.Error("proposals should be cleared")
	}
}

// TestSynthesizer_Synthesize_Unanimous verifies unanimous synthesis
func TestSynthesizer_Synthesize_Unanimous(t *testing.T) {
	s := NewSynthesizer()
	s.AddProposal(NewProposal("agent-1", "same-solution", 0.9, ""))
	s.AddProposal(NewProposal("agent-2", "same-solution", 0.8, ""))

	result, err := s.Synthesize()
	if err != nil {
		t.Fatalf("synthesis failed: %v", err)
	}
	if result.Rule != "unanimous" {
		t.Errorf("expected unanimous rule, got %s", result.Rule)
	}
}

// TestSynthesizer_Synthesize_DomainExpertise verifies domain expertise synthesis
func TestSynthesizer_Synthesize_DomainExpertise(t *testing.T) {
	s := NewSynthesizer()
	s.AddProposal(NewProposal("agent-1", "solution-1", 0.9, ""))
	s.AddProposal(NewProposal("agent-2", "solution-2", 0.7, ""))

	result, err := s.Synthesize()
	if err != nil {
		t.Fatalf("synthesis failed: %v", err)
	}
	if result.Rule != "domain_expertise" {
		t.Errorf("expected domain_expertise rule, got %s", result.Rule)
	}
	if result.WinningAgent != "agent-1" {
		t.Errorf("expected agent-1 to win, got %s", result.WinningAgent)
	}
}

// TestSynthesizer_Synthesize_NoProposals verifies no proposals error
func TestSynthesizer_Synthesize_NoProposals(t *testing.T) {
	s := NewSynthesizer()
	_, err := s.Synthesize()
	if err == nil {
		t.Error("expected error for no proposals")
	}
}

// TestSynthesizer_Synthesize_CannotSynthesize verifies escalation
func TestSynthesizer_Synthesize_CannotSynthesize(t *testing.T) {
	s := NewSynthesizer()
	// Add proposals with same confidence (tie) and different solutions
	s.AddProposal(NewProposal("agent-1", "solution-1", 0.8, ""))
	s.AddProposal(NewProposal("agent-2", "solution-2", 0.8, ""))

	_, err := s.Synthesize()
	if err == nil {
		t.Error("expected error for tie")
	}
	if !errors.Is(err, ErrCannotSynthesize) {
		t.Errorf("expected ErrCannotSynthesize, got %v", err)
	}
}

// TestSynthesizer_DetectConflict_NoConflict verifies no conflict detection
func TestSynthesizer_DetectConflict_NoConflict(t *testing.T) {
	s := NewSynthesizer()
	s.AddProposal(NewProposal("agent-1", "same", 0.9, ""))
	s.AddProposal(NewProposal("agent-2", "same", 0.8, ""))

	conflict := s.DetectConflict()
	if conflict != NoConflict {
		t.Errorf("expected NoConflict, got %d", conflict)
	}
}

// TestSynthesizer_DetectConflict_MajorConflict verifies major conflict detection
func TestSynthesizer_DetectConflict_MajorConflict(t *testing.T) {
	s := NewSynthesizer()
	s.AddProposal(NewProposal("agent-1", "solution-1", 0.9, ""))
	s.AddProposal(NewProposal("agent-2", "solution-2", 0.8, ""))

	conflict := s.DetectConflict()
	if conflict != MajorConflict {
		t.Errorf("expected MajorConflict, got %d", conflict)
	}
}

// TestSynthesizer_DetectConflict_SingleProposal verifies single proposal case
func TestSynthesizer_DetectConflict_SingleProposal(t *testing.T) {
	s := NewSynthesizer()
	s.AddProposal(NewProposal("agent-1", "solution", 0.9, ""))

	conflict := s.DetectConflict()
	if conflict != NoConflict {
		t.Errorf("expected NoConflict for single proposal, got %d", conflict)
	}
}

// TestSynthesizer_DetectConflict_NoProposals verifies no proposals case
func TestSynthesizer_DetectConflict_NoProposals(t *testing.T) {
	s := NewSynthesizer()
	conflict := s.DetectConflict()
	if conflict != NoConflict {
		t.Errorf("expected NoConflict for no proposals, got %d", conflict)
	}
}

// TestSolutionsEqual verifies solution equality
func TestSolutionsEqual(t *testing.T) {
	tests := []struct {
		a, b     interface{}
		expected bool
	}{
		{"same", "same", true},
		{"a", "b", false},
		{123, 123, true},
		{123, 456, false},
		{[]string{"a", "b"}, []string{"a", "b"}, true},
		{[]string{"a"}, []string{"b"}, false},
		{nil, nil, true},
	}

	for _, tt := range tests {
		result := solutionsEqual(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("solutionsEqual(%v, %v) = %v, expected %v", tt.a, tt.b, result, tt.expected)
		}
	}
}

// TestSynthesisResult verifies result structure
func TestSynthesisResult(t *testing.T) {
	result := &SynthesisResult{
		Solution:     "test-solution",
		Rule:         "unanimous",
		WinningAgent: "agent-1",
		Reasoning:    "all agreed",
		Proposals:    []*Proposal{{AgentID: "agent-1"}},
	}

	if result.Solution != "test-solution" {
		t.Error("Solution not set")
	}
	if result.Rule != "unanimous" {
		t.Error("Rule not set")
	}
}

// TestConflictType verifies conflict type constants
func TestConflictType(t *testing.T) {
	if NoConflict != 0 {
		t.Error("NoConflict should be 0")
	}
	if MinorConflict != 1 {
		t.Error("MinorConflict should be 1")
	}
	if MediumConflict != 2 {
		t.Error("MediumConflict should be 2")
	}
	if MajorConflict != 3 {
		t.Error("MajorConflict should be 3")
	}
}
