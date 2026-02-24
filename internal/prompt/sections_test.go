package prompt

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTaskSection(t *testing.T) {
	ws := WorkstreamSpec{
		ID:                 "00-025-01",
		Title:              "Prompt Consolidation",
		Description:        "Consolidate 5 scattered prompt-building functions.",
		AcceptanceCriteria: []string{"All prompt-building logic consolidated", "TaskSection pure function"},
		SpecID:             "sdp_dev-h7qu",
	}
	got := TaskSection(ws)
	goldenPath := filepath.Join("testdata", "task_section.golden")
	want := readGolden(t, goldenPath)
	if got != want {
		t.Errorf("TaskSection mismatch:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestTaskSectionForReview(t *testing.T) {
	ws := WorkstreamSpec{
		ID:          "sdp_dev-4pg",
		Title:       "QA: Test coverage",
		Description: "Raise coverage to 80%",
	}
	got := TaskSectionForReview(ws)
	goldenPath := filepath.Join("testdata", "task_section_review.golden")
	want := readGolden(t, goldenPath)
	if got != want {
		t.Errorf("TaskSectionForReview mismatch:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestBoundarySection(t *testing.T) {
	in := BoundaryInput{
		AllowedPathPrefixes:   []string{"internal/", "cmd/"},
		ForbiddenPathPrefixes: []string{".git/"},
		ControlPathPrefixes:   []string{".beads/", ".sdp/"},
	}
	got := BoundarySection(in)
	goldenPath := filepath.Join("testdata", "boundary_section.golden")
	want := readGolden(t, goldenPath)
	if got != want {
		t.Errorf("BoundarySection mismatch:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestEvidenceSection(t *testing.T) {
	in := EvidenceInput{
		Content:      `{"verdict":"approve","comments":[]}`,
		CompletedWS:  []string{"00-025-01 (abc123)"},
		ReviewStatus: "pending",
	}
	got := EvidenceSection(in)
	goldenPath := filepath.Join("testdata", "evidence_section.golden")
	want := readGolden(t, goldenPath)
	if got != want {
		t.Errorf("EvidenceSection mismatch:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestAcceptanceCriteriaSection(t *testing.T) {
	items := []string{"Criterion one", "Criterion two"}
	got := AcceptanceCriteriaSection(items)
	goldenPath := filepath.Join("testdata", "acceptance_criteria_section.golden")
	want := readGolden(t, goldenPath)
	if got != want {
		t.Errorf("AcceptanceCriteriaSection mismatch:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestScopeFilesSection(t *testing.T) {
	files := []string{"internal/prompt/sections.go", "internal/llm/prompt.go"}
	got := ScopeFilesSection(files)
	goldenPath := filepath.Join("testdata", "scope_files_section.golden")
	want := readGolden(t, goldenPath)
	if got != want {
		t.Errorf("ScopeFilesSection mismatch:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func readGolden(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v", path, err)
	}
	return string(b)
}
