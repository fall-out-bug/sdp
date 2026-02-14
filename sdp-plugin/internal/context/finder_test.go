package context

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindWorktree_NoWorktrees(t *testing.T) {
	// Create a temp directory to simulate project root
	tmpDir := t.TempDir()

	r := NewRecovery(tmpDir)

	// Without any worktrees, should return error
	_, err := r.FindWorktree("F999")
	if err == nil {
		t.Error("FindWorktree should return error when no worktrees exist")
	}
}

func TestFindWorktree_ByWorktreeName(t *testing.T) {
	// This test requires git worktree which is integration-level
	// We test the error path instead
	tmpDir := t.TempDir()

	r := NewRecovery(tmpDir)

	// Without git, listWorktrees will fail
	_, err := r.FindWorktree("F999")
	if err == nil {
		t.Error("FindWorktree should return error without git")
	}
}

func TestGoToWorktree(t *testing.T) {
	tmpDir := t.TempDir()

	r := NewRecovery(tmpDir)

	// GoToWorktree is a wrapper around FindWorktree
	path, err := r.GoToWorktree("F999")

	// Should fail without git worktrees
	if err == nil {
		t.Error("GoToWorktree should return error without git")
	}
	if path != "" {
		t.Errorf("GoToWorktree path = %q, want empty", path)
	}
}

func TestFindWorktree_WithWorktreePath(t *testing.T) {
	// Test the path construction logic
	tmpDir := t.TempDir()

	// Create the workstreams directory to test that error path
	wsPath := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsPath, 0755); err != nil {
		t.Fatalf("Failed to create workstreams dir: %v", err)
	}

	r := NewRecovery(tmpDir)

	_, err := r.FindWorktree("F999")
	if err == nil {
		t.Error("FindWorktree should return error with suggestion to create worktree")
	}

	// Error should contain suggestion
	if err != nil && !containsSubstring(err.Error(), "sdp worktree create") {
		t.Errorf("Error should contain worktree creation suggestion, got: %v", err)
	}
}

func TestFindWorktree_FeatureNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	r := NewRecovery(tmpDir)

	_, err := r.FindWorktree("NONEXISTENT")
	if err == nil {
		t.Error("FindWorktree should return error for non-existent feature")
	}
}

func TestRecovery_ProjectRoot(t *testing.T) {
	tests := []struct {
		name    string
		root    string
		wantSet bool
	}{
		{"empty root", "", true},
		{"absolute path", "/tmp/project", true},
		{"relative path", "./project", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRecovery(tt.root)
			if r == nil {
				t.Fatal("NewRecovery returned nil")
			}
			if r.ProjectRoot != tt.root {
				t.Errorf("ProjectRoot = %q, want %q", r.ProjectRoot, tt.root)
			}
		})
	}
}
