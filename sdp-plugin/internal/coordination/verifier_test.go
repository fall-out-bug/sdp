package coordination

import (
	"context"
	"testing"
)

func TestVerificationPipeline_Run(t *testing.T) {
	pipeline := NewVerificationPipeline()
	pipeline.AddRule(&MockRule{name: "test_rule", status: "pass"})

	spec := &WorkstreamSpec{ID: "00-051-06", AcceptanceCriteria: []string{"AC1", "AC2"}}
	code := &CodeSnapshot{Files: []string{"test.go"}}

	results, err := pipeline.Run(context.Background(), spec, code)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
}

func TestVerificationPipeline_AllRulesPass(t *testing.T) {
	pipeline := NewVerificationPipeline()
	pipeline.AddRule(&MockRule{name: "rule1", status: "pass"})
	pipeline.AddRule(&MockRule{name: "rule2", status: "pass"})
	pipeline.AddRule(&MockRule{name: "rule3", status: "pass"})

	spec := &WorkstreamSpec{ID: "test", AcceptanceCriteria: []string{}}
	code := &CodeSnapshot{Files: []string{}}

	results, _ := pipeline.Run(context.Background(), spec, code)

	allPass := true
	for _, r := range results {
		if r.Status != "pass" {
			allPass = false
		}
	}

	if !allPass {
		t.Error("Expected all rules to pass")
	}
}

func TestVerificationPipeline_RuleFails(t *testing.T) {
	pipeline := NewVerificationPipeline()
	pipeline.AddRule(&MockRule{name: "failing_rule", status: "fail", message: "coverage too low"})

	spec := &WorkstreamSpec{ID: "test", AcceptanceCriteria: []string{}}
	code := &CodeSnapshot{Files: []string{}}

	results, _ := pipeline.Run(context.Background(), spec, code)

	if results[0].Status != "fail" {
		t.Error("Expected rule to fail")
	}
	if results[0].Message != "coverage too low" {
		t.Errorf("Expected 'coverage too low', got %s", results[0].Message)
	}
}

func TestVerificationPipeline_HasErrors(t *testing.T) {
	pipeline := NewVerificationPipeline()

	results := []VerificationResult{
		{RuleName: "rule1", Status: "pass"},
		{RuleName: "rule2", Status: "fail", Severity: "error"},
	}

	if !pipeline.HasErrors(results) {
		t.Error("Expected HasErrors to return true")
	}
}

func TestVerificationPipeline_HasNoErrors(t *testing.T) {
	pipeline := NewVerificationPipeline()

	results := []VerificationResult{
		{RuleName: "rule1", Status: "pass"},
		{RuleName: "rule2", Status: "fail", Severity: "warning"},
	}

	if pipeline.HasErrors(results) {
		t.Error("Expected HasErrors to return false")
	}
}

func TestVerificationPipeline_HasWarnings(t *testing.T) {
	pipeline := NewVerificationPipeline()

	results := []VerificationResult{
		{RuleName: "rule1", Status: "pass"},
		{RuleName: "rule2", Status: "fail", Severity: "warning"},
	}

	if !pipeline.HasWarnings(results) {
		t.Error("Expected HasWarnings to return true")
	}
}

func TestVerificationPipeline_HasNoWarnings(t *testing.T) {
	pipeline := NewVerificationPipeline()

	results := []VerificationResult{
		{RuleName: "rule1", Status: "pass"},
		{RuleName: "rule2", Status: "pass"},
	}

	if pipeline.HasWarnings(results) {
		t.Error("Expected HasWarnings to return false")
	}
}

// MockRule for testing
type MockRule struct {
	name    string
	status  string
	message string
}

func (r *MockRule) Name() string {
	return r.name
}

func (r *MockRule) Verify(ctx context.Context, spec *WorkstreamSpec, code *CodeSnapshot) (*VerificationResult, error) {
	return &VerificationResult{
		RuleName: r.name,
		Status:   r.status,
		Message:  r.message,
	}, nil
}
