package dashboard

import (
	"testing"
)

func TestNewApp(t *testing.T) {
	app := New()

	if app == nil {
		t.Fatal("New() returned nil")
	}

	if app.state.ActiveTab != 0 {
		t.Errorf("Expected ActiveTab to be 0, got %d", app.state.ActiveTab)
	}

	if !app.state.Loading {
		t.Error("Expected initial Loading state to be true")
	}

	if app.quit {
		t.Error("Expected initial quit state to be false")
	}
}

func TestTabTypeString(t *testing.T) {
	tests := []struct {
		tab      TabType
		expected string
	}{
		{TabWorkstreams, "Workstreams"},
		{TabIdeas, "Ideas"},
		{TabTests, "Tests"},
		{TabActivity, "Activity"},
		{TabType(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.tab.String()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestDashboardState(t *testing.T) {
	state := DashboardState{
		ActiveTab:   1,
		Workstreams: make(map[string][]WorkstreamSummary),
		Ideas:       []IdeaSummary{},
		Loading:     false,
	}

	if state.ActiveTab != 1 {
		t.Errorf("Expected ActiveTab 1, got %d", state.ActiveTab)
	}

	if state.Workstreams == nil {
		t.Error("Expected Workstreams map to be initialized")
	}

	if state.Ideas == nil {
		t.Error("Expected Ideas slice to be initialized")
	}

	if state.Loading {
		t.Error("Expected Loading to be false")
	}
}

func TestWorkstreamSummary(t *testing.T) {
	ws := WorkstreamSummary{
		ID:       "sdp-abc",
		Title:    "Test workstream",
		Status:   "open",
		Priority: "P1",
		Assignee: "user",
		Size:     "MEDIUM",
	}

	if ws.ID != "sdp-abc" {
		t.Errorf("Expected ID 'sdp-abc', got '%s'", ws.ID)
	}

	if ws.Title != "Test workstream" {
		t.Errorf("Expected Title 'Test workstream', got '%s'", ws.Title)
	}

	if ws.Status != "open" {
		t.Errorf("Expected Status 'open', got '%s'", ws.Status)
	}

	if ws.Priority != "P1" {
		t.Errorf("Expected Priority 'P1', got '%s'", ws.Priority)
	}

	if ws.Assignee != "user" {
		t.Errorf("Expected Assignee 'user', got '%s'", ws.Assignee)
	}

	if ws.Size != "MEDIUM" {
		t.Errorf("Expected Size 'MEDIUM', got '%s'", ws.Size)
	}
}

func TestIdeaSummary(t *testing.T) {
	idea := IdeaSummary{
		Title: "Test Idea",
		Path:  "docs/drafts/test-idea.md",
	}

	if idea.Title != "Test Idea" {
		t.Errorf("Expected Title 'Test Idea', got '%s'", idea.Title)
	}

	if idea.Path != "docs/drafts/test-idea.md" {
		t.Errorf("Expected Path 'docs/drafts/test-idea.md', got '%s'", idea.Path)
	}
}

func TestTestSummary(t *testing.T) {
	summary := TestSummary{
		Coverage: "85.5%",
		Passing:  42,
		Total:    50,
		QualityGates: []GateStatus{
			{Name: "Coverage", Passed: true},
			{Name: "Linting", Passed: false},
		},
	}

	if summary.Coverage != "85.5%" {
		t.Errorf("Expected Coverage '85.5%%', got '%s'", summary.Coverage)
	}

	if summary.Passing != 42 {
		t.Errorf("Expected Passing 42, got %d", summary.Passing)
	}

	if summary.Total != 50 {
		t.Errorf("Expected Total 50, got %d", summary.Total)
	}

	if len(summary.QualityGates) != 2 {
		t.Errorf("Expected 2 quality gates, got %d", len(summary.QualityGates))
	}
}

func TestGateStatus(t *testing.T) {
	gate := GateStatus{
		Name:   "Test Gate",
		Passed: true,
	}

	if gate.Name != "Test Gate" {
		t.Errorf("Expected Name 'Test Gate', got '%s'", gate.Name)
	}

	if !gate.Passed {
		t.Error("Expected Passed to be true")
	}
}

// Test DashboardState initialization with data
func TestDashboardStateWithData(t *testing.T) {
	state := DashboardState{
		ActiveTab: 2, //nolint:unused // Test fixture
		Workstreams: map[string][]WorkstreamSummary{
			"open": {
				{ID: "sdp-001", Title: "First task", Status: "open", Priority: "P1"},
			},
		},
		Ideas: []IdeaSummary{
			{Title: "Idea 1", Path: "docs/drafts/idea1.md"},
		},
		TestResults: TestSummary{
			Coverage: "90%",
			Passing:  10,
			Total:    10,
		},
		Loading: false, //nolint:unused // Test fixture
	}

	// Verify workstreams
	openWS, ok := state.Workstreams["open"]
	if !ok {
		t.Fatal("Expected 'open' workstreams to exist")
	}

	if len(openWS) != 1 {
		t.Errorf("Expected 1 open workstream, got %d", len(openWS))
	}

	if openWS[0].ID != "sdp-001" {
		t.Errorf("Expected workstream ID 'sdp-001', got '%s'", openWS[0].ID)
	}

	// Verify ideas
	if len(state.Ideas) != 1 {
		t.Errorf("Expected 1 idea, got %d", len(state.Ideas))
	}

	// Verify test results
	if state.TestResults.Coverage != "90%" {
		t.Errorf("Expected coverage '90%%', got '%s'", state.TestResults.Coverage)
	}
}
