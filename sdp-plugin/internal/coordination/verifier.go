package coordination

import (
	"context"
	"slices"
)

// VerificationResult represents the result of a verification rule (AC3)
type VerificationResult struct {
	RuleName string
	Status   string // pass, fail, skip
	Message  string
	Details  []string
	Severity string // error, warning, info
}

// WorkstreamSpec represents a workstream specification (AC1)
type WorkstreamSpec struct {
	ID                 string
	AcceptanceCriteria []string
	ScopeFiles         []string
}

// CodeSnapshot represents current code state
type CodeSnapshot struct {
	Files     []string
	Entities  []string
	LOCMetric map[string]int
}

// VerificationRule is the interface for verification rules (AC4)
type VerificationRule interface {
	Name() string
	Verify(ctx context.Context, spec *WorkstreamSpec, code *CodeSnapshot) (*VerificationResult, error)
}

// VerificationPipeline runs all verification rules (AC2)
type VerificationPipeline struct {
	rules []VerificationRule
}

// NewVerificationPipeline creates a new pipeline
func NewVerificationPipeline() *VerificationPipeline {
	return &VerificationPipeline{
		rules: []VerificationRule{},
	}
}

// AddRule adds a rule to the pipeline
func (p *VerificationPipeline) AddRule(rule VerificationRule) {
	p.rules = append(p.rules, rule)
}

// Run executes all rules and returns results
func (p *VerificationPipeline) Run(ctx context.Context, spec *WorkstreamSpec, code *CodeSnapshot) ([]VerificationResult, error) {
	var results []VerificationResult

	for _, rule := range p.rules {
		result, err := rule.Verify(ctx, spec, code)
		if err != nil {
			result = &VerificationResult{
				RuleName: rule.Name(),
				Status:   "fail",
				Message:  err.Error(),
				Severity: "error",
			}
		}
		results = append(results, *result)
	}

	return results, nil
}

// HasErrors returns true if any rule failed with error severity
func (p *VerificationPipeline) HasErrors(results []VerificationResult) bool {
	return slices.ContainsFunc(results, func(r VerificationResult) bool {
		return r.Status == "fail" && r.Severity == "error"
	})
}

// HasWarnings returns true if any rule failed with warning severity
func (p *VerificationPipeline) HasWarnings(results []VerificationResult) bool {
	return slices.ContainsFunc(results, func(r VerificationResult) bool {
		return r.Status == "fail" && r.Severity == "warning"
	})
}
