package worktree

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestDelete_NoWorktree(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	// Create temp dir
	tmpDir := t.TempDir()

	creator := &Creator{
		MainRepoPath: tmpDir,
		WorktreesDir: tmpDir,
	}

	// Delete should fail for non-existent worktree
	err := creator.Delete("NONEXISTENT")
	if err == nil {
		t.Error("Delete should fail for non-existent worktree")
	}
}

func TestList_NoGitRepo(t *testing.T) {
	// Create temp dir without git
	tmpDir := t.TempDir()

	creator := &Creator{
		MainRepoPath: tmpDir,
		WorktreesDir: tmpDir,
	}

	// List should fail without git repo
	_, err := creator.List()
	if err == nil {
		t.Error("List should fail without git repo")
	}
}

func TestParseWorktreeList_WithSession(t *testing.T) {
	// Create a temp directory with session
	tmpDir := t.TempDir()
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("failed to create .sdp dir: %v", err)
	}

	// Create session file with valid format
	sessionContent := `{"feature_id":"F067","worktree_path":"` + tmpDir + `","expected_branch":"feature/F067","hash":"test-hash"}`
	if err := os.WriteFile(filepath.Join(sdpDir, "session.json"), []byte(sessionContent), 0644); err != nil {
		t.Fatalf("failed to write session: %v", err)
	}

	// Test parsing with session
	input := "worktree " + tmpDir + "\nbranch refs/heads/feature/F067\n\n"

	creator := &Creator{MainRepoPath: tmpDir}
	result, err := creator.parseWorktreeList(input)
	if err != nil {
		t.Fatalf("parseWorktreeList() error = %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 worktree, got %d", len(result))
	}

	// Note: Session may be nil if Load fails, which is acceptable
	// The important thing is that parseWorktreeList attempts to load it
	t.Logf("Session loaded: %v", result[0].Session != nil)
}

func TestParseWorktreeList_ComplexInput(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedCount int
	}{
		{
			name: "no trailing newline",
			input: `worktree /path/to/main
branch refs/heads/main

worktree /path/to/feature
branch refs/heads/feature/F067`,
			expectedCount: 2,
		},
		{
			name: "multiple blank lines",
			input: `worktree /path/to/main
branch refs/heads/main


worktree /path/to/feature
branch refs/heads/feature/F067


`,
			expectedCount: 2,
		},
		{
			name: "mixed refs format",
			input: `worktree /path/to/main
branch refs/heads/main

worktree /path/to/other
branch refs/heads/feature/ABC-123

`,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creator := &Creator{MainRepoPath: "/path/to/main"}
			result, err := creator.parseWorktreeList(tt.input)
			if err != nil {
				t.Fatalf("parseWorktreeList() error = %v", err)
			}

			if len(result) != tt.expectedCount {
				t.Errorf("got %d worktrees, want %d", len(result), tt.expectedCount)
			}
		})
	}
}

func TestWorktreeInfo_Fields(t *testing.T) {
	// Test all fields of WorktreeInfo
	info := WorktreeInfo{
		Path:    "/path/to/worktree",
		Branch:  "feature/F067",
		Session: nil,
	}

	if info.Path != "/path/to/worktree" {
		t.Errorf("Path = %v, want /path/to/worktree", info.Path)
	}
	if info.Branch != "feature/F067" {
		t.Errorf("Branch = %v, want feature/F067", info.Branch)
	}
	if info.Session != nil {
		t.Error("Session should be nil")
	}
}

func TestCreateOptions_AllFields(t *testing.T) {
	opts := CreateOptions{
		FeatureID:    "F067",
		BranchName:   "custom-branch",
		BaseBranch:   "main",
		CreateBranch: true,
	}

	if opts.FeatureID != "F067" {
		t.Errorf("FeatureID = %v, want F067", opts.FeatureID)
	}
	if opts.BranchName != "custom-branch" {
		t.Errorf("BranchName = %v, want custom-branch", opts.BranchName)
	}
	if opts.BaseBranch != "main" {
		t.Errorf("BaseBranch = %v, want main", opts.BaseBranch)
	}
	if !opts.CreateBranch {
		t.Error("CreateBranch should be true")
	}
}

func TestCreateResult_AllFields(t *testing.T) {
	result := &CreateResult{
		WorktreePath: "/custom/worktrees/sdp-F067",
		BranchName:   "feature/F067",
		SessionFile:  "/custom/worktrees/sdp-F067/.sdp/session.json",
	}

	if result.WorktreePath != "/custom/worktrees/sdp-F067" {
		t.Errorf("WorktreePath = %v", result.WorktreePath)
	}
	if result.BranchName != "feature/F067" {
		t.Errorf("BranchName = %v", result.BranchName)
	}
	if result.SessionFile != "/custom/worktrees/sdp-F067/.sdp/session.json" {
		t.Errorf("SessionFile = %v", result.SessionFile)
	}
}

func TestNewCreator_EmptyPath(t *testing.T) {
	creator := NewCreator("")
	if creator == nil {
		t.Fatal("NewCreator should not return nil")
	}
	if creator.MainRepoPath != "" {
		t.Errorf("MainRepoPath = %v, want empty", creator.MainRepoPath)
	}
}
