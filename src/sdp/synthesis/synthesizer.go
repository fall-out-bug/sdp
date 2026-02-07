package synthesis

import (
	"errors"
	"fmt"
)

// ErrCannotSynthesize is returned when synthesis cannot resolve conflict
var ErrCannotSynthesize = errors.New("cannot synthesize: conflicting proposals")

// ConflictType represents the type of conflict between proposals
type ConflictType int

const (
	NoConflict ConflictType = iota
	MinorConflict  // Formatting, style
	MediumConflict // Approach, implementation
	MajorConflict  // Architecture, design
)

// SynthesisResult represents the result of synthesis
type SynthesisResult struct {
	Solution      interface{}
	Rule          string
	WinningAgent  string
	Reasoning     string
	Proposals     []*Proposal
}

// Synthesizer manages agent proposals and applies synthesis rules
type Synthesizer struct {
	proposals map[string]*Proposal
}

// NewSynthesizer creates a new synthesizer
func NewSynthesizer() *Synthesizer {
	return &Synthesizer{
		proposals: make(map[string]*Proposal),
	}
}

// AddProposal adds a proposal from an agent
func (s *Synthesizer) AddProposal(proposal *Proposal) {
	if proposal == nil {
		return
	}
	s.proposals[proposal.AgentID] = proposal
}

// GetProposals returns all proposals
func (s *Synthesizer) GetProposals() []*Proposal {
	result := make([]*Proposal, 0, len(s.proposals))
	for _, p := range s.proposals {
		result = append(result, p)
	}
	return result
}

// Clear removes all proposals
func (s *Synthesizer) Clear() {
	s.proposals = make(map[string]*Proposal)
}

// Synthesize applies synthesis rules in priority order
func (s *Synthesizer) Synthesize() (*SynthesisResult, error) {
	if len(s.proposals) == 0 {
		return nil, errors.New("no proposals to synthesize")
	}

	// Rule 1: Check for unanimous agreement
	if result := s.unanimousResult(); result != nil {
		return result, nil
	}

	// Rule 2: Domain expertise (highest confidence)
	if result := s.domainExpertiseResult(); result != nil {
		return result, nil
	}

	// Rule 3: Quality gate (not implemented yet - skip to next)
	// Rule 4: Merge solutions (not implemented yet - skip to next)

	// Rule 5: Escalate to human
	return nil, ErrCannotSynthesize
}

// DetectConflict detects the type of conflict between proposals
func (s *Synthesizer) DetectConflict() ConflictType {
	if len(s.proposals) < 2 {
		return NoConflict
	}

	// Get all unique solutions
	solutions := make(map[string]bool)
	for _, p := range s.proposals {
		key := fmt.Sprintf("%v", p.Solution)
		solutions[key] = true
	}

	// If all solutions are the same, no conflict
	if len(solutions) == 1 {
		return NoConflict
	}

	// Multiple different solutions - major conflict
	// (TODO: implement minor/medium conflict detection)
	return MajorConflict
}

// unanimousResult checks if all agents agree
func (s *Synthesizer) unanimousResult() *SynthesisResult {
	if len(s.proposals) == 0 {
		return nil
	}

	// Get first solution as reference
	var firstSolution interface{}
	for _, p := range s.proposals {
		firstSolution = p.Solution
		break
	}

	// Check if all solutions match
	for _, p := range s.proposals {
		if !solutionsEqual(p.Solution, firstSolution) {
			return nil // Not unanimous
		}
	}

	// All agree
	return &SynthesisResult{
		Solution:  firstSolution,
		Rule:      "unanimous",
		Reasoning: "All agents agreed on the same solution",
		Proposals: s.GetProposals(),
	}
}

// domainExpertiseResult returns solution from agent with highest confidence
// Returns nil if there's a tie in confidence or multiple solutions with same confidence
func (s *Synthesizer) domainExpertiseResult() *SynthesisResult {
	if len(s.proposals) == 0 {
		return nil
	}

	var bestProposal *Proposal
	maxConfidence := 0.0
	agentsWithMaxConfidence := 0

	for _, p := range s.proposals {
		if p.Confidence > maxConfidence {
			maxConfidence = p.Confidence
			bestProposal = p
			agentsWithMaxConfidence = 1
		} else if p.Confidence == maxConfidence {
			agentsWithMaxConfidence++
		}
	}

	// If there's a tie in confidence, cannot use domain expertise rule
	if agentsWithMaxConfidence > 1 {
		return nil
	}

	return &SynthesisResult{
		Solution:     bestProposal.Solution,
		Rule:         "domain_expertise",
		WinningAgent: bestProposal.AgentID,
		Reasoning:    fmt.Sprintf("Agent %s has highest confidence (%.2f)", bestProposal.AgentID, bestProposal.Confidence),
		Proposals:    s.GetProposals(),
	}
}

// solutionsEqual compares two solutions for equality
func solutionsEqual(a, b interface{}) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}
