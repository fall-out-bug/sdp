package dashboard

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFetchWorkstreams_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	app := New()
	summaries := app.fetchWorkstreams()

	// Should have all status groups initialized
	requiredGroups := []string{"open", "in_progress", "completed", "blocked"}
	for _, group := range requiredGroups {
		if _, ok := summaries[group]; !ok {
			t.Errorf("Missing status group: %s", group)
		}
	}
}

func TestFetchWorkstreams_WithWorkstreamFiles(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	// Create workstream directory structure
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "F001", "backlog")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatal(err)
	}

	wsContent := `---
ws_id: 00-001-01
feature_id: F001
status: completed
priority: 0
size: S
depends_on: []
---

## Goal

Completed workstream

## Acceptance Criteria

- [x] Completed item
`
	blockedContent := `---
ws_id: 00-001-03
feature_id: F001
status: ready
priority: 2
size: M
depends_on: ["00-001-02"]
---

## Goal

Blocked workstream

## Acceptance Criteria

- [ ] Blocked item
`
	readyAfterDependencyContent := `---
ws_id: 00-001-02
feature_id: F001
status: ready
priority: 1
size: M
depends_on: ["00-001-01"]
---

## Goal

Ready after dependency

## Acceptance Criteria

- [ ] Ready item after dependency
`
	wsFile := filepath.Join(wsDir, "00-001-01.md")
	if err := os.WriteFile(wsFile, []byte(wsContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(wsDir, "00-001-02.md"), []byte(readyAfterDependencyContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(wsDir, "00-001-03.md"), []byte(blockedContent), 0644); err != nil {
		t.Fatal(err)
	}

	app := New()
	summaries := app.fetchWorkstreams()

	if len(summaries["open"]) != 1 {
		t.Fatalf("expected 1 open workstream, got %d", len(summaries["open"]))
	}
	if summaries["open"][0].ID != "00-001-02" {
		t.Fatalf("expected open workstream 00-001-02, got %s", summaries["open"][0].ID)
	}
	if summaries["open"][0].Title != "Ready after dependency" {
		t.Fatalf("expected title from parsed workstream, got %q", summaries["open"][0].Title)
	}
	if len(summaries["blocked"]) != 1 {
		t.Fatalf("expected 1 blocked workstream, got %d", len(summaries["blocked"]))
	}
	if summaries["blocked"][0].ID != "00-001-03" {
		t.Fatalf("expected blocked workstream 00-001-03, got %s", summaries["blocked"][0].ID)
	}
}

func TestFetchIdeas_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	app := New()
	ideas := app.fetchIdeas()

	if len(ideas) != 0 {
		t.Errorf("Expected 0 ideas in empty directory, got %d", len(ideas))
	}
}

func TestFetchIdeas_WithDrafts(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	// Create drafts directory
	draftsDir := filepath.Join(tmpDir, "docs", "drafts")
	if err := os.MkdirAll(draftsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create draft files
	draft1 := filepath.Join(draftsDir, "my-feature-idea.md")
	draft2 := filepath.Join(draftsDir, "another-idea.md")

	if err := os.WriteFile(draft1, []byte("# My Feature Idea"), 0644); err != nil {
		t.Fatal(err)
	}
	// Small delay to ensure different modification times
	time.Sleep(10 * time.Millisecond)
	if err := os.WriteFile(draft2, []byte("# Another Idea"), 0644); err != nil {
		t.Fatal(err)
	}

	app := New()
	ideas := app.fetchIdeas()

	if len(ideas) != 2 {
		t.Errorf("Expected 2 ideas, got %d", len(ideas))
	}

	// Should be sorted by date (newest first)
	if len(ideas) >= 2 {
		if !ideas[0].Date.After(ideas[1].Date) && !ideas[0].Date.Equal(ideas[1].Date) {
			t.Error("Ideas should be sorted by date (newest first)")
		}
	}
}

func TestFetchTestResults(t *testing.T) {
	app := New()
	result := app.fetchTestResults()

	// Should return placeholder data
	if result.Coverage == "" {
		t.Error("Coverage should not be empty")
	}

	if len(result.QualityGates) == 0 {
		t.Error("Should have quality gates")
	}
}

func TestWorkstreamSummary_Fields(t *testing.T) {
	ws := WorkstreamSummary{
		ID:       "00-001-01",
		Title:    "Test Workstream",
		Status:   "open",
		Priority: "P2",
		Assignee: "user",
		Size:     "S",
	}

	if ws.ID != "00-001-01" {
		t.Errorf("ID = %s, want 00-001-01", ws.ID)
	}
	if ws.Title != "Test Workstream" {
		t.Errorf("Title = %s, want Test Workstream", ws.Title)
	}
}

func TestIdeaSummary_Fields(t *testing.T) {
	now := time.Now()
	idea := IdeaSummary{
		Title: "Test Idea",
		Path:  "/path/to/idea.md",
		Date:  now,
	}

	if idea.Title != "Test Idea" {
		t.Errorf("Title = %s, want Test Idea", idea.Title)
	}
	if idea.Path != "/path/to/idea.md" {
		t.Errorf("Path = %s", idea.Path)
	}
}

func TestTestSummary_Fields(t *testing.T) {
	ts := TestSummary{
		Coverage: "80%",
		Passing:  10,
		Total:    12,
		LastRun:  time.Now(),
		QualityGates: []GateStatus{
			{Name: "Coverage", Passed: true},
		},
	}

	if ts.Coverage != "80%" {
		t.Errorf("Coverage = %s, want 80%%", ts.Coverage)
	}
	if ts.Passing != 10 {
		t.Errorf("Passing = %d, want 10", ts.Passing)
	}
	if ts.Total != 12 {
		t.Errorf("Total = %d, want 12", ts.Total)
	}
	if len(ts.QualityGates) != 1 {
		t.Errorf("QualityGates count = %d, want 1", len(ts.QualityGates))
	}
}

func TestGateStatus_Fields(t *testing.T) {
	gate := GateStatus{
		Name:   "Coverage",
		Passed: true,
	}

	if gate.Name != "Coverage" {
		t.Errorf("Name = %s, want Coverage", gate.Name)
	}
	if !gate.Passed {
		t.Error("Passed should be true")
	}
}

// TestFetchNextStep tests the fetchNextStep function.
func TestFetchNextStep(t *testing.T) {
	app := New()
	nextStep := app.fetchNextStep()

	// Should return a valid recommendation
	if nextStep.Command == "" {
		t.Error("Command should not be empty")
	}
	if nextStep.Reason == "" {
		t.Error("Reason should not be empty")
	}
	if nextStep.Confidence < 0 || nextStep.Confidence > 1 {
		t.Errorf("Confidence should be between 0 and 1, got %f", nextStep.Confidence)
	}
	if nextStep.Category == "" {
		t.Error("Category should not be empty")
	}
}

// TestRefreshCmd tests the refreshCmd function.
func TestRefreshCmd(t *testing.T) {
	app := New()
	cmd := app.refreshCmd()

	if cmd == nil {
		t.Error("refreshCmd should return a non-nil command")
	}
}

// TestTickCmd tests the tickCmd function.
func TestTickCmd(t *testing.T) {
	app := New()
	cmd := app.tickCmd()

	if cmd == nil {
		t.Error("tickCmd should return a non-nil command")
	}
}

// TestOpenSelectedItem tests the openSelectedItem function.
func TestOpenSelectedItem(t *testing.T) {
	app := New()
	app.state.Workstreams = map[string][]WorkstreamSummary{
		"open": {
			{ID: "00-001-01", Title: "Test", Status: "open"},
		},
	}
	app.state.CursorPos = 0
	app.state.ActiveTab = 0 // TabWorkstreams

	// Should not panic
	app.openSelectedItem()
}

// TestOpenSelectedItem_EmptyList tests openSelectedItem with empty list.
func TestOpenSelectedItem_EmptyList(t *testing.T) {
	app := New()
	app.state.Workstreams = map[string][]WorkstreamSummary{
		"open": {},
	}
	app.state.CursorPos = 0
	app.state.ActiveTab = 0 // TabWorkstreams

	// Should not panic with empty list
	app.openSelectedItem()
}
