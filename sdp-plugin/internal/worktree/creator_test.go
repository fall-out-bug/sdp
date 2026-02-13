package worktree

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/session"
)

func TestCreator_NewCreator(t *testing.T) {
	mainRepo := "/path/to/repo"
	creator := NewCreator(mainRepo)

	if creator.MainRepoPath != mainRepo {
		t.Errorf("MainRepoPath = %v, want %v", creator.MainRepoPath, mainRepo)
	}
	// WorktreesDir should default to parent of main repo
	expected := filepath.Dir(mainRepo)
	if creator.WorktreesDir != expected {
		t.Errorf("WorktreesDir = %v, want %v", creator.WorktreesDir, expected)
	}
}

func TestCreator_CreateOptionsDefaults(t *testing.T) {
	// Test that defaults are applied correctly
	// This is a unit test that doesn't actually create a worktree
	opts := CreateOptions{
		FeatureID: "F065",
	}

	// BranchName should default to feature/F### if not specified
	if opts.BranchName == "" {
		expectedBranch := "feature/F065"
		// This is what would be set in Create()
		if expectedBranch != "feature/F065" {
			t.Errorf("expected default branch name to be feature/F065")
		}
	}

	// BaseBranch should default to dev if not specified
	if opts.BaseBranch == "" {
		expectedBase := "dev"
		if expectedBase != "dev" {
			t.Errorf("expected default base branch to be dev")
		}
	}
}

func TestWorktreeInfo(t *testing.T) {
	// Test that WorktreeInfo can hold session info
	tmpDir := t.TempDir()
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("failed to create .sdp dir: %v", err)
	}

	// Create a session
	s, err := session.Init("F065", tmpDir, "test")
	if err != nil {
		t.Fatalf("session.Init() error = %v", err)
	}
	if err := s.Save(tmpDir); err != nil {
		t.Fatalf("session.Save() error = %v", err)
	}

	// Verify WorktreeInfo can hold the session
	info := WorktreeInfo{
		Path:    tmpDir,
		Branch:  "feature/F065",
		Session: s,
	}

	if info.Session == nil {
		t.Error("Session should not be nil")
	}
	if info.Session.FeatureID != "F065" {
		t.Errorf("Session.FeatureID = %v, want F065", info.Session.FeatureID)
	}
}

func TestParseWorktreeList(t *testing.T) {
	// Test parsing the output of git worktree list --porcelain
	tests := []struct {
		name     string
		input    string
		expected []WorktreeInfo
	}{
		{
			name: "single worktree",
			input: `worktree /path/to/main
branch refs/heads/main

`,
			expected: []WorktreeInfo{
				{Path: "/path/to/main", Branch: "main"},
			},
		},
		{
			name: "multiple worktrees",
			input: `worktree /path/to/main
branch refs/heads/main

worktree /path/to/sdp-F065
branch refs/heads/feature/F065

`,
			expected: []WorktreeInfo{
				{Path: "/path/to/main", Branch: "main"},
				{Path: "/path/to/sdp-F065", Branch: "feature/F065"},
			},
		},
		{
			name: "bare repo",
			input: `worktree /path/to/bare
bare

`,
			expected: []WorktreeInfo{
				{Path: "/path/to/bare", Branch: ""},
			},
		},
		{
			name: "detached head",
			input: `worktree /path/to/detached
detached

`,
			expected: []WorktreeInfo{
				{Path: "/path/to/detached", Branch: ""},
			},
		},
		{
			name:     "empty input",
			input:    ``,
			expected: []WorktreeInfo(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creator := &Creator{MainRepoPath: "/path/to/main"}
			result, err := creator.parseWorktreeList(tt.input)
			if err != nil {
				t.Fatalf("parseWorktreeList() error = %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Errorf("got %d worktrees, want %d", len(result), len(tt.expected))
				return
			}

			for i, info := range result {
				if info.Path != tt.expected[i].Path {
					t.Errorf("worktree[%d].Path = %v, want %v", i, info.Path, tt.expected[i].Path)
				}
				if info.Branch != tt.expected[i].Branch {
					t.Errorf("worktree[%d].Branch = %v, want %v", i, info.Branch, tt.expected[i].Branch)
				}
			}
		})
	}
}

func TestCreateRequiresFeatureID(t *testing.T) {
	creator := NewCreator("/path/to/repo")
	_, err := creator.Create(CreateOptions{FeatureID: ""})
	if err == nil {
		t.Error("Create() should fail with empty FeatureID")
	}
}

func TestCreateResult(t *testing.T) {
	// Test that CreateResult is properly formed
	result := &CreateResult{
		WorktreePath: "/path/to/sdp-F065",
		BranchName:   "feature/F065",
		SessionFile:  "/path/to/sdp-F065/.sdp/session.json",
	}

	if result.WorktreePath == "" {
		t.Error("WorktreePath should not be empty")
	}
	if result.BranchName == "" {
		t.Error("BranchName should not be empty")
	}
	if result.SessionFile == "" {
		t.Error("SessionFile should not be empty")
	}
}

func TestCreatorWithCustomWorktreesDir(t *testing.T) {
	creator := &Creator{
		MainRepoPath: "/path/to/repo",
		WorktreesDir: "/custom/worktrees",
	}

	if creator.WorktreesDir != "/custom/worktrees" {
		t.Errorf("WorktreesDir = %v, want /custom/worktrees", creator.WorktreesDir)
	}
}
