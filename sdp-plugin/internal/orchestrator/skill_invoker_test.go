package orchestrator

import (
	"strings"
	"testing"
)

// MockSkillTool mocks the Skill tool interface
type MockSkillTool struct {
	invokedSkill string
	invokedArgs  string
	results      map[string]string
	shouldFail   bool
}

func (m *MockSkillTool) Invoke(skillName string, args string) (string, error) {
	m.invokedSkill = skillName
	m.invokedArgs = args

	if m.shouldFail {
		return "", ErrSkillInvocationFailed
	}

	if result, ok := m.results[skillName]; ok {
		return result, nil
	}

	return "", nil
}

func (m *MockSkillTool) GetInvokedSkill() string {
	return m.invokedSkill
}

func (m *MockSkillTool) GetInvokedArgs() string {
	return m.invokedArgs
}

func TestSkillInvoker_NewSkillInvoker(t *testing.T) {
	tool := &MockSkillTool{
		results: make(map[string]string),
	}

	invoker := NewSkillInvoker(tool)

	if invoker == nil {
		t.Fatal("Expected non-nil invoker")
	}

	if invoker.skillTool == nil {
		t.Error("Expected skillTool to be initialized")
	}
}

func TestSkillInvoker_InvokeIdea_Success(t *testing.T) {
	// AC2: Orchestrator calls @idea for requirements gathering
	tool := &MockSkillTool{
		results: map[string]string{
			"idea": `{"problem": "User needs authentication", "users": ["developers"], "success_criteria": ["secure login"]}`,
		},
	}

	invoker := NewSkillInvoker(tool)

	result, err := invoker.InvokeIdea("F001", "Add user authentication")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify skill was called
	if tool.GetInvokedSkill() != "idea" {
		t.Errorf("Expected skill 'idea', got '%s'", tool.GetInvokedSkill())
	}

	// Verify args contain feature ID
	if tool.GetInvokedArgs() == "" {
		t.Error("Expected non-empty args")
	}

	// Verify result
	if result.Problem == "" {
		t.Error("Expected problem to be populated")
	}

	if len(result.Users) == 0 {
		t.Error("Expected users to be populated")
	}

	if len(result.SuccessCriteria) == 0 {
		t.Error("Expected success_criteria to be populated")
	}
}

func TestSkillInvoker_InvokeDesign_Success(t *testing.T) {
	// AC3: Orchestrator calls @design for workstream planning
	tool := &MockSkillTool{
		results: map[string]string{
			"design": `{"workstreams": ["WS-001", "WS-002", "WS-003"], "spec_path": "docs/drafts/idea-auth.md"}`,
		},
	}

	invoker := NewSkillInvoker(tool)

	result, err := invoker.InvokeDesign("docs/drafts/idea-auth.md")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify skill was called
	if tool.GetInvokedSkill() != "design" {
		t.Errorf("Expected skill 'design', got '%s'", tool.GetInvokedSkill())
	}

	// Verify result
	if len(result.Workstreams) == 0 {
		t.Error("Expected workstreams to be populated")
	}

	if result.SpecPath == "" {
		t.Error("Expected spec_path to be populated")
	}

	// Verify it creates docs/drafts/idea-{slug}.md
	if result.SpecPath[:12] != "docs/drafts/" {
		t.Errorf("Expected spec_path to start with 'docs/drafts/', got '%s'", result.SpecPath)
	}
}

func TestSkillInvoker_InvokeOneshot_Success(t *testing.T) {
	// AC4: Orchestrator calls @oneshot for autonomous execution
	tool := &MockSkillTool{
		results: map[string]string{
			"oneshot": `{"agent_id": "agent-123", "status": "started"}`,
		},
	}

	invoker := NewSkillInvoker(tool)

	err := invoker.InvokeOneshot("F001")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify skill was called
	if tool.GetInvokedSkill() != "oneshot" {
		t.Errorf("Expected skill 'oneshot', got '%s'", tool.GetInvokedSkill())
	}

	// Verify args contain feature ID
	args := tool.GetInvokedArgs()
	if args == "" {
		t.Error("Expected non-empty args")
	}
}

func TestSkillInvoker_InvokeIdea_Failure(t *testing.T) {
	tool := &MockSkillTool{
		shouldFail: true,
	}

	invoker := NewSkillInvoker(tool)

	_, err := invoker.InvokeIdea("F001", "Add user authentication")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Error should wrap ErrSkillInvocationFailed or contain "failed to invoke @idea"
	if err == nil || (!strings.Contains(err.Error(), "failed to invoke @idea")) {
		t.Errorf("Expected error to contain 'failed to invoke @idea', got %v", err)
	}
}

func TestSkillInvoker_InvokeDesign_Failure(t *testing.T) {
	tool := &MockSkillTool{
		shouldFail: true,
	}

	invoker := NewSkillInvoker(tool)

	_, err := invoker.InvokeDesign("docs/drafts/idea-auth.md")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestSkillInvoker_InvokeOneshot_Failure(t *testing.T) {
	tool := &MockSkillTool{
		shouldFail: true,
	}

	invoker := NewSkillInvoker(tool)

	err := invoker.InvokeOneshot("F001")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestIdeaResult_FromJSON(t *testing.T) {
	json := `{"problem": "Test problem", "users": ["user1", "user2"], "success_criteria": ["criteria1", "criteria2"]}`

	result, err := IdeaResultFromJSON(json)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Problem != "Test problem" {
		t.Errorf("Expected problem 'Test problem', got '%s'", result.Problem)
	}

	if len(result.Users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(result.Users))
	}

	if len(result.SuccessCriteria) != 2 {
		t.Errorf("Expected 2 success criteria, got %d", len(result.SuccessCriteria))
	}
}

func TestDesignResult_FromJSON(t *testing.T) {
	json := `{"workstreams": ["WS-001", "WS-002"], "spec_path": "docs/drafts/test.md"}`

	result, err := DesignResultFromJSON(json)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result.Workstreams) != 2 {
		t.Errorf("Expected 2 workstreams, got %d", len(result.Workstreams))
	}

	if result.SpecPath != "docs/drafts/test.md" {
		t.Errorf("Expected spec_path 'docs/drafts/test.md', got '%s'", result.SpecPath)
	}
}
