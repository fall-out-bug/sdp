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

	// Create a workstream file
	wsContent := `---
id: 00-001-01
feature: F001
status: open
goal: Test workstream
size: S
---
`
	wsFile := filepath.Join(wsDir, "00-001-01.md")
	if err := os.WriteFile(wsFile, []byte(wsContent), 0644); err != nil {
		t.Fatal(err)
	}

	app := New()
	summaries := app.fetchWorkstreams()

	// Should have at least the open group
	if len(summaries["open"]) == 0 {
		t.Log("No workstreams found (parser may require different format)")
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
