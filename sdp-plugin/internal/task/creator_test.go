package task

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewCreator(t *testing.T) {
	t.Run("default config", func(t *testing.T) {
		c := NewCreator(CreatorConfig{})
		if c.config.WorkstreamDir != "docs/workstreams/backlog" {
			t.Errorf("expected default workstream dir, got %s", c.config.WorkstreamDir)
		}
		if c.config.IssuesDir != "docs/issues" {
			t.Errorf("expected default issues dir, got %s", c.config.IssuesDir)
		}
	})

	t.Run("custom config", func(t *testing.T) {
		c := NewCreator(CreatorConfig{
			WorkstreamDir: "custom/ws",
			IssuesDir:     "custom/issues",
			ProjectID:     "myproject",
		})
		if c.config.WorkstreamDir != "custom/ws" {
			t.Errorf("expected custom workstream dir, got %s", c.config.WorkstreamDir)
		}
		if c.config.ProjectID != "myproject" {
			t.Errorf("expected custom project ID, got %s", c.config.ProjectID)
		}
	})
}

func TestCreator_CreateWorkstream(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatal(err)
	}

	c := NewCreator(CreatorConfig{
		WorkstreamDir: wsDir,
		ProjectID:     "00",
	})

	t.Run("create bug workstream", func(t *testing.T) {
		task := &Task{
			Type:       TypeBug,
			Title:      "Fix CI Go version",
			Priority:   PriorityP1,
			FeatureID:  "F064",
			Goal:       "Fix the CI pipeline",
			Context:    "CI is failing due to Go version mismatch",
			ScopeFiles: []string{"cmd/ci/main.go"},
		}

		ws, err := c.CreateWorkstream(task)
		if err != nil {
			t.Fatalf("CreateWorkstream() error = %v", err)
		}

		// Verify WS ID format: 99-{FEATURE_NUM}-{SEQ} for bugs
		if !strings.HasPrefix(ws.WSID, "99-064-") {
			t.Errorf("expected bug WS ID prefix '99-064-', got %s", ws.WSID)
		}

		// Verify file exists
		if _, err := os.Stat(ws.Path); err != nil {
			t.Errorf("workstream file not created: %s", ws.Path)
		}

		// Verify content
		content, err := os.ReadFile(ws.Path)
		if err != nil {
			t.Fatal(err)
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, "ws_id: "+ws.WSID) {
			t.Error("workstream file missing ws_id frontmatter")
		}
		if !strings.Contains(contentStr, "type: bug") {
			t.Error("workstream file missing type frontmatter")
		}
		if !strings.Contains(contentStr, task.Title) {
			t.Error("workstream file missing title")
		}
	})

	t.Run("create task workstream", func(t *testing.T) {
		task := &Task{
			Type:      TypeTask,
			Title:     "Add user authentication",
			Priority:  PriorityP2,
			FeatureID: "F064",
			Goal:      "Implement OAuth2",
			DependsOn: []string{"00-064-01"},
		}

		ws, err := c.CreateWorkstream(task)
		if err != nil {
			t.Fatalf("CreateWorkstream() error = %v", err)
		}

		// Verify WS ID format: 00-{FEATURE}-{SEQ} for tasks
		if !strings.HasPrefix(ws.WSID, "00-064-") {
			t.Errorf("expected task WS ID prefix '00-064-', got %s", ws.WSID)
		}
	})

	t.Run("create hotfix workstream", func(t *testing.T) {
		task := &Task{
			Type:       TypeHotfix,
			Title:      "Production database connection fix",
			Priority:   PriorityP0,
			FeatureID:  "F064",
			Goal:       "Restore database connectivity",
			BranchBase: "main",
		}

		ws, err := c.CreateWorkstream(task)
		if err != nil {
			t.Fatalf("CreateWorkstream() error = %v", err)
		}

		// Verify branch_base in content
		content, err := os.ReadFile(ws.Path)
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(content), "branch_base: main") {
			t.Error("hotfix workstream missing branch_base: main")
		}
	})
}

func TestCreator_CreateIssue(t *testing.T) {
	tmpDir := t.TempDir()
	issuesDir := filepath.Join(tmpDir, "docs", "issues")
	indexFile := filepath.Join(tmpDir, ".sdp", "issues-index.jsonl")

	if err := os.MkdirAll(issuesDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(indexFile), 0755); err != nil {
		t.Fatal(err)
	}

	c := NewCreator(CreatorConfig{
		IssuesDir: issuesDir,
		IndexFile: indexFile,
	})

	t.Run("create issue without feature", func(t *testing.T) {
		task := &Task{
			Type:     TypeBug,
			Title:    "Authentication fails",
			Priority: PriorityP1,
			Goal:     "Fix login issue",
			Context:  "Users cannot log in",
		}

		issue, err := c.CreateIssue(task)
		if err != nil {
			t.Fatalf("CreateIssue() error = %v", err)
		}

		// Verify issue ID format
		if !strings.HasPrefix(issue.IssueID, "ISSUE-") {
			t.Errorf("expected issue ID prefix 'ISSUE-', got %s", issue.IssueID)
		}

		// Verify file exists
		if _, err := os.Stat(issue.Path); err != nil {
			t.Errorf("issue file not created: %s", issue.Path)
		}

		// Verify index updated
		indexContent, err := os.ReadFile(indexFile)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(indexContent), issue.IssueID) {
			t.Error("issue not added to index")
		}
	})

	t.Run("sequencing issues", func(t *testing.T) {
		issue1, err := c.CreateIssue(&Task{Title: "First"})
		if err != nil {
			t.Fatal(err)
		}

		issue2, err := c.CreateIssue(&Task{Title: "Second"})
		if err != nil {
			t.Fatal(err)
		}

		// Verify sequential IDs
		if issue1.IssueID >= issue2.IssueID {
			t.Errorf("expected sequential issue IDs, got %s then %s", issue1.IssueID, issue2.IssueID)
		}
	})
}

func TestCreator_GenerateWSID(t *testing.T) {
	c := NewCreator(CreatorConfig{ProjectID: "00"})

	tests := []struct {
		name      string
		taskType  Type
		featureID string
		seq       int
		expected  string
	}{
		{"bug task", TypeBug, "F064", 1, "99-064-01"},
		{"bug task high seq", TypeBug, "F064", 9, "99-064-09"},
		{"regular task", TypeTask, "F064", 1, "00-064-01"},
		{"hotfix", TypeHotfix, "F064", 1, "99-064-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := c.generateWSID(tt.taskType, tt.featureID, tt.seq)
			if result != tt.expected {
				t.Errorf("generateWSID() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestCreator_NextSequence(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create existing workstream files
	existingFiles := []string{
		"99-F064-0001.md",
		"99-F064-0002.md",
		"00-064-0001.md",
	}
	for _, f := range existingFiles {
		path := filepath.Join(wsDir, f)
		if err := os.WriteFile(path, []byte("---\nws_id: test\n---"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	c := NewCreator(CreatorConfig{WorkstreamDir: wsDir})

	t.Run("next bug sequence", func(t *testing.T) {
		seq := c.nextSequence("99-F064")
		if seq != 3 {
			t.Errorf("expected next sequence 3, got %d", seq)
		}
	})

	t.Run("next task sequence", func(t *testing.T) {
		seq := c.nextSequence("00-064")
		if seq != 2 {
			t.Errorf("expected next sequence 2, got %d", seq)
		}
	})

	t.Run("new prefix sequence", func(t *testing.T) {
		seq := c.nextSequence("00-999")
		if seq != 1 {
			t.Errorf("expected next sequence 1, got %d", seq)
		}
	})
}

func TestCreator_ValidationError(t *testing.T) {
	c := NewCreator(CreatorConfig{})

	t.Run("empty title", func(t *testing.T) {
		_, err := c.CreateWorkstream(&Task{
			Type:      TypeBug,
			Title:     "",
			FeatureID: "F064",
		})
		if err == nil {
			t.Error("expected error for empty title")
		}
	})

	t.Run("missing feature for workstream", func(t *testing.T) {
		_, err := c.CreateWorkstream(&Task{
			Type:  TypeTask,
			Title: "Test",
			// No FeatureID
		})
		if err == nil {
			t.Error("expected error for missing feature ID")
		}
	})
}
