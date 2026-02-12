package coordination

import (
	"context"
	"testing"
)

func TestACCoverageRule_AllCovered(t *testing.T) {
	rule := NewACCoverageRule()

	spec := &WorkstreamSpec{
		ID:                 "00-051-06",
		AcceptanceCriteria: []string{"AC1: Parse requirements", "AC2: Verify requirements"},
	}
	code := &CodeSnapshot{
		Files:    []string{"verifier.go", "verifier_test.go"},
		Entities: []string{"TestACCoverageRule_AllCovered", "TestACCoverageRule_NotCovered"},
	}

	result, err := rule.Verify(context.Background(), spec, code)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	// Rule should pass when tests exist
	if result.Status != "pass" {
		t.Errorf("Expected pass, got %s - %s", result.Status, result.Message)
	}
}

func TestACCoverageRule_NotCovered(t *testing.T) {
	rule := NewACCoverageRule()

	spec := &WorkstreamSpec{
		ID:                 "00-051-06",
		AcceptanceCriteria: []string{"AC1: Parse requirements", "AC2: Verify requirements"},
	}
	code := &CodeSnapshot{
		Files:    []string{"verifier.go"},
		Entities: []string{}, // No test entities
	}

	result, err := rule.Verify(context.Background(), spec, code)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	// Rule should warn when no tests found
	if result.Status != "fail" {
		t.Error("Expected fail when no test entities")
	}
}

func TestScopeBoundariesRule_InScope(t *testing.T) {
	rule := NewScopeBoundariesRule()

	spec := &WorkstreamSpec{
		ID:         "00-051-06",
		ScopeFiles: []string{"verifier.go", "verifier_test.go"},
	}
	code := &CodeSnapshot{
		Files: []string{"verifier.go", "verifier_test.go"},
	}

	result, err := rule.Verify(context.Background(), spec, code)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	if result.Status != "pass" {
		t.Errorf("Expected pass for in-scope files, got %s", result.Status)
	}
}

func TestScopeBoundariesRule_OutOfScope(t *testing.T) {
	rule := NewScopeBoundariesRule()

	spec := &WorkstreamSpec{
		ID:         "00-051-06",
		ScopeFiles: []string{"verifier.go"},
	}
	code := &CodeSnapshot{
		Files: []string{"verifier.go", "unrelated.go"}, // unrelated.go not in scope
	}

	result, err := rule.Verify(context.Background(), spec, code)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	if result.Status != "fail" {
		t.Error("Expected fail for out-of-scope files")
	}
}

func TestLOCLimitRule_UnderLimit(t *testing.T) {
	rule := NewLOCLimitRule(200)

	spec := &WorkstreamSpec{ID: "00-051-06"}
	code := &CodeSnapshot{
		LOCMetric: map[string]int{
			"verifier.go": 150,
			"rules.go":    180,
		},
	}

	result, err := rule.Verify(context.Background(), spec, code)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	if result.Status != "pass" {
		t.Errorf("Expected pass for files under 200 LOC, got %s", result.Status)
	}
}

func TestLOCLimitRule_OverLimit(t *testing.T) {
	rule := NewLOCLimitRule(200)

	spec := &WorkstreamSpec{ID: "00-051-06"}
	code := &CodeSnapshot{
		LOCMetric: map[string]int{
			"verifier.go": 250, // Over limit
		},
	}

	result, err := rule.Verify(context.Background(), spec, code)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	if result.Status != "fail" {
		t.Error("Expected fail for file over 200 LOC")
	}
}

func TestDependencyCheckRule_NoNewDeps(t *testing.T) {
	rule := NewDependencyCheckRule()

	spec := &WorkstreamSpec{ID: "00-051-06"}
	code := &CodeSnapshot{
		Files: []string{"verifier.go"},
	}

	result, err := rule.Verify(context.Background(), spec, code)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	// Should pass when no new dependencies detected
	if result.Status != "pass" {
		t.Errorf("Expected pass, got %s", result.Status)
	}
}
