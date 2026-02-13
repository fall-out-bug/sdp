package coordination

import (
	"context"
	"strings"
)

// ACCoverageRule verifies all acceptance criteria have tests (AC2)
type ACCoverageRule struct{}

// NewACCoverageRule creates a new AC coverage rule
func NewACCoverageRule() *ACCoverageRule {
	return &ACCoverageRule{}
}

// Name returns the rule name
func (r *ACCoverageRule) Name() string {
	return "ac_coverage"
}

// Verify checks if all ACs are covered by tests
func (r *ACCoverageRule) Verify(ctx context.Context, spec *WorkstreamSpec, code *CodeSnapshot) (*VerificationResult, error) {
	if len(spec.AcceptanceCriteria) == 0 {
		return &VerificationResult{
			RuleName: r.Name(),
			Status:   "skip",
			Message:  "No acceptance criteria defined",
		}, nil
	}

	// Check if test entities exist
	hasTests := false
	for _, entity := range code.Entities {
		if strings.HasPrefix(entity, "Test") {
			hasTests = true
			break
		}
	}

	if !hasTests {
		return &VerificationResult{
			RuleName: r.Name(),
			Status:   "fail",
			Message:  "No test functions found for acceptance criteria",
			Severity: "warning",
		}, nil
	}

	return &VerificationResult{
		RuleName: r.Name(),
		Status:   "pass",
		Message:  "Tests found for acceptance criteria",
	}, nil
}

// ScopeBoundariesRule verifies code changes only touch scope_files
type ScopeBoundariesRule struct{}

// NewScopeBoundariesRule creates a new scope boundaries rule
func NewScopeBoundariesRule() *ScopeBoundariesRule {
	return &ScopeBoundariesRule{}
}

// Name returns the rule name
func (r *ScopeBoundariesRule) Name() string {
	return "scope_boundaries"
}

// Verify checks if all changed files are in scope
func (r *ScopeBoundariesRule) Verify(ctx context.Context, spec *WorkstreamSpec, code *CodeSnapshot) (*VerificationResult, error) {
	scopeSet := make(map[string]bool)
	for _, f := range spec.ScopeFiles {
		scopeSet[f] = true
	}

	var outOfScope []string
	for _, f := range code.Files {
		if !scopeSet[f] {
			outOfScope = append(outOfScope, f)
		}
	}

	if len(outOfScope) > 0 {
		return &VerificationResult{
			RuleName: r.Name(),
			Status:   "fail",
			Message:  "Files modified outside scope",
			Details:  outOfScope,
			Severity: "warning",
		}, nil
	}

	return &VerificationResult{
		RuleName: r.Name(),
		Status:   "pass",
		Message:  "All files within scope",
	}, nil
}

// LOCLimitRule verifies files don't exceed LOC limit
type LOCLimitRule struct {
	limit int
}

// NewLOCLimitRule creates a new LOC limit rule
func NewLOCLimitRule(limit int) *LOCLimitRule {
	return &LOCLimitRule{limit: limit}
}

// Name returns the rule name
func (r *LOCLimitRule) Name() string {
	return "loc_limit"
}

// Verify checks if all files are under LOC limit
func (r *LOCLimitRule) Verify(ctx context.Context, spec *WorkstreamSpec, code *CodeSnapshot) (*VerificationResult, error) {
	var overLimit []string

	for file, loc := range code.LOCMetric {
		if loc > r.limit {
			overLimit = append(overLimit, file)
		}
	}

	if len(overLimit) > 0 {
		return &VerificationResult{
			RuleName: r.Name(),
			Status:   "fail",
			Message:  "Files exceed LOC limit",
			Details:  overLimit,
			Severity: "error",
		}, nil
	}

	return &VerificationResult{
		RuleName: r.Name(),
		Status:   "pass",
		Message:  "All files within LOC limit",
	}, nil
}

// DependencyCheckRule checks for new dependencies
type DependencyCheckRule struct{}

// NewDependencyCheckRule creates a new dependency check rule
func NewDependencyCheckRule() *DependencyCheckRule {
	return &DependencyCheckRule{}
}

// Name returns the rule name
func (r *DependencyCheckRule) Name() string {
	return "dependency_check"
}

// Verify checks for new dependencies
func (r *DependencyCheckRule) Verify(ctx context.Context, spec *WorkstreamSpec, code *CodeSnapshot) (*VerificationResult, error) {
	// Check if go.mod or package.json was modified
	for _, f := range code.Files {
		if f == "go.mod" || f == "go.sum" || f == "package.json" || f == "package-lock.json" {
			return &VerificationResult{
				RuleName: r.Name(),
				Status:   "fail",
				Message:  "Dependency file modified - requires approval",
				Severity: "warning",
			}, nil
		}
	}

	return &VerificationResult{
		RuleName: r.Name(),
		Status:   "pass",
		Message:  "No new dependencies",
	}, nil
}
