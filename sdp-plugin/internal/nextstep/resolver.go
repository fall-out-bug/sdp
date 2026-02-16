package nextstep

// Resolver converts project state into next-step recommendations.
// It implements deterministic logic for "what should happen next".
type Resolver struct {
	rules []RecommendationRule
}

// NewResolver creates a new resolver with default rules.
func NewResolver() *Resolver {
	return &Resolver{
		rules: defaultRules(),
	}
}

// Recommend generates a next-step recommendation based on project state.
// The recommendation is deterministic given the same input state.
func (r *Resolver) Recommend(state ProjectState) (*Recommendation, error) {
	for _, rule := range r.rules {
		if rec := rule.Evaluate(state); rec != nil {
			rec.Version = ContractVersion
			return rec, nil
		}
	}
	return r.fallbackRecommendation(state), nil
}

// fallbackRecommendation provides a low-confidence fallback when no rule matches.
func (r *Resolver) fallbackRecommendation(state ProjectState) *Recommendation {
	return &Recommendation{
		Command:    "sdp status",
		Reason:     "Check current project state to determine next steps",
		Confidence: 0.5,
		Category:   CategoryInformation,
		Version:    ContractVersion,
		Alternatives: []Alternative{
			{Command: "sdp doctor", Reason: "Run diagnostics if issues suspected"},
			{Command: "sdp --help", Reason: "View available commands"},
		},
	}
}
