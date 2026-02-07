package synthesis

// SynthesisRule defines the interface for synthesis rules
type SynthesisRule interface {
	// Name returns the rule name
	Name() string

	// Priority returns the rule priority (1=highest, 5=lowest)
	Priority() int

	// CanApply checks if this rule can be applied to the given proposals
	CanApply(proposals []*Proposal) bool

	// Apply executes the rule and returns a synthesis result
	Apply(proposals []*Proposal) (*SynthesisResult, error)
}

// UnanimousRule applies when all agents propose the same solution
type UnanimousRule struct{}

// NewUnanimousRule creates a new unanimous rule
func NewUnanimousRule() *UnanimousRule {
	return &UnanimousRule{}
}

// Name returns the rule name
func (r *UnanimousRule) Name() string {
	return "unanimous"
}

// Priority returns 1 (highest priority)
func (r *UnanimousRule) Priority() int {
	return 1
}

// CanApply checks if all proposals have the same solution
func (r *UnanimousRule) CanApply(proposals []*Proposal) bool {
	if len(proposals) < 2 {
		return false
	}

	// Get first solution as reference
	var firstSolution interface{}
	for _, p := range proposals {
		firstSolution = p.Solution
		break
	}

	// Check if all solutions match
	for _, p := range proposals {
		if !solutionsEqual(p.Solution, firstSolution) {
			return false
		}
	}

	return true
}

// Apply returns the unanimous solution
func (r *UnanimousRule) Apply(proposals []*Proposal) (*SynthesisResult, error) {
	// Get first solution
	var firstSolution interface{}
	for _, p := range proposals {
		firstSolution = p.Solution
		break
	}

	return &SynthesisResult{
		Solution:  firstSolution,
		Rule:      r.Name(),
		Reasoning: "All agents agreed on the same solution",
		Proposals: proposals,
	}, nil
}

// DomainExpertiseRule applies when one agent has highest confidence
type DomainExpertiseRule struct{}

// NewDomainExpertiseRule creates a new domain expertise rule
func NewDomainExpertiseRule() *DomainExpertiseRule {
	return &DomainExpertiseRule{}
}

// Name returns the rule name
func (r *DomainExpertiseRule) Name() string {
	return "domain_expertise"
}

// Priority returns 2
func (r *DomainExpertiseRule) Priority() int {
	return 2
}

// CanApply checks if there's a unique agent with highest confidence
func (r *DomainExpertiseRule) CanApply(proposals []*Proposal) bool {
	if len(proposals) == 0 {
		return false
	}

	maxConfidence := 0.0
	count := 0
	for _, p := range proposals {
		if p.Confidence > maxConfidence {
			maxConfidence = p.Confidence
			count = 1
		} else if p.Confidence == maxConfidence {
			count++
		}
	}

	// Must have unique highest confidence
	return count == 1
}

// Apply returns solution from agent with highest confidence
func (r *DomainExpertiseRule) Apply(proposals []*Proposal) (*SynthesisResult, error) {
	var best *Proposal
	for _, p := range proposals {
		if best == nil || p.Confidence > best.Confidence {
			best = p
		}
	}

	return &SynthesisResult{
		Solution:     best.Solution,
		Rule:         r.Name(),
		WinningAgent: best.AgentID,
		Reasoning:    "Agent has highest confidence",
		Proposals:    proposals,
	}, nil
}
