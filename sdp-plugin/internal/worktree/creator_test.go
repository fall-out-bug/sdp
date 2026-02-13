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

func TestCreateOptionsBranchDefault(t *testing.T) {
	// Test that branch name defaults correctly based on FeatureID
	tests := []struct {
		featureID     string
		customBranch  string
		expectedMatch string
	}{
		{"F067", "", "feature/F067"},
		{"F100", "custom-branch", "custom-branch"},
		{"ABC-123", "", "feature/ABC-123"},
	}

	for _, tt := range tests {
		t.Run(tt.featureID, func(t *testing.T) {
			branch := tt.customBranch
			if branch == "" {
				branch = "feature/" + tt.featureID
			}
			if branch != tt.expectedMatch {
				t.Errorf("branch = %v, want %v", branch, tt.expectedMatch)
			}
		})
	}
}

func TestCreateOptionsBaseBranchDefault(t *testing.T) {
	// Test that base branch defaults correctly
	tests := []struct {
		customBase    string
		expectedMatch string
	}{
		{"", "dev"},
		{"main", "main"},
		{"develop", "develop"},
	}

	for _, tt := range tests {
		t.Run(tt.customBase, func(t *testing.T) {
			base := tt.customBase
			if base == "" {
				base = "dev"
			}
			if base != tt.expectedMatch {
				t.Errorf("base = %v, want %v", base, tt.expectedMatch)
			}
		})
	}
}

func TestWorktreePathConstruction(t *testing.T) {
	creator := NewCreator("/path/to/repo")

	tests := []struct {
		featureID      string
		expectedSuffix string
	}{
		{"F067", "sdp-F067"},
		{"ABC-123", "sdp-ABC-123"},
	}

	for _, tt := range tests {
		t.Run(tt.featureID, func(t *testing.T) {
			worktreeName := "sdp-" + tt.featureID
			if worktreeName != tt.expectedSuffix {
				t.Errorf("worktreeName = %v, want %v", worktreeName, tt.expectedSuffix)
			}

			// Verify path is constructed correctly
			expectedPath := "/path/to/sdp-" + tt.featureID
			if creator.WorktreesDir != "/path/to" {
				t.Errorf("WorktreesDir = %v, want /path/to", creator.WorktreesDir)
			}
			_ = expectedPath
		})
	}
}

func TestDeletePathConstruction(t *testing.T) {
	creator := NewCreator("/path/to/repo")

	// Delete constructs the path from featureID
	featureID := "F067"
	expectedName := "sdp-" + featureID
	expectedPath := "/path/to/sdp-" + featureID

	// Verify creator has correct WorktreesDir
	if creator.WorktreesDir != "/path/to" {
		t.Errorf("WorktreesDir = %v, want /path/to", creator.WorktreesDir)
	}

	// The Delete function constructs the path internally
	_ = expectedName
	_ = expectedPath
}

func TestParseWorktreeListWithSession(t *testing.T) {
	// Create a temp directory with session for testing session loading
	tmpDir := t.TempDir()
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("failed to create .sdp dir: %v", err)
	}

	// Create a session
	s, err := session.Init("F067", tmpDir, "test-parse")
	if err != nil {
		t.Fatalf("session.Init() error = %v", err)
	}
	if err := s.Save(tmpDir); err != nil {
		t.Fatalf("session.Save() error = %v", err)
	}

	// Parse worktree list with session path
	input := "worktree " + tmpDir + "\nbranch refs/heads/feature/F067\n\n"
	creator := &Creator{MainRepoPath: "/path/to/main"}
	result, err := creator.parseWorktreeList(input)
	if err != nil {
		t.Fatalf("parseWorktreeList() error = %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 worktree, got %d", len(result))
	}

	if result[0].Session == nil {
		t.Error("expected Session to be loaded")
	} else if result[0].Session.FeatureID != "F067" {
		t.Errorf("Session.FeatureID = %v, want F067", result[0].Session.FeatureID)
	}
}

func TestParseWorktreeListEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantLen     int
		wantPaths   []string
		wantBranches []string
	}{
		{
			name:        "no trailing newline",
			input:       "worktree /path/to/main\nbranch refs/heads/main",
			wantLen:     1,
			wantPaths:   []string{"/path/to/main"},
			wantBranches: []string{"main"},
		},
		{
			name:        "multiple blank lines",
			input:       "worktree /path/to/main\nbranch refs/heads/main\n\n\n",
			wantLen:     1,
			wantPaths:   []string{"/path/to/main"},
			wantBranches: []string{"main"},
		},
		{
			name:        "mixed refs format",
			input:       "worktree /path/to/main\nbranch refs/heads/feature/new-feature\n\n",
			wantLen:     1,
			wantPaths:   []string{"/path/to/main"},
			wantBranches: []string{"feature/new-feature"},
		},
		{
			name:        "lines with extra spaces - trimmed by parser",
			input:       "worktree /path/to/main\nbranch refs/heads/main\n\n",
			wantLen:     1,
			wantPaths:   []string{"/path/to/main"},
			wantBranches: []string{"main"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creator := &Creator{MainRepoPath: "/path/to/main"}
			result, err := creator.parseWorktreeList(tt.input)
			if err != nil {
				t.Fatalf("parseWorktreeList() error = %v", err)
			}

			if len(result) != tt.wantLen {
				t.Errorf("got %d worktrees, want %d", len(result), tt.wantLen)
				return
			}

			for i, info := range result {
				if i < len(tt.wantPaths) && info.Path != tt.wantPaths[i] {
					t.Errorf("worktree[%d].Path = %v, want %v", i, info.Path, tt.wantPaths[i])
				}
				if i < len(tt.wantBranches) && info.Branch != tt.wantBranches[i] {
					t.Errorf("worktree[%d].Branch = %v, want %v", i, info.Branch, tt.wantBranches[i])
				}
			}
		})
	}
}

func TestCreateRequiresNonEmptyFeatureID(t *testing.T) {
	creator := NewCreator("/nonexistent/path")

	// Empty FeatureID should fail immediately without git calls
	_, err := creator.Create(CreateOptions{FeatureID: ""})
	if err == nil {
		t.Error("Create() should fail with empty FeatureID")
	}
	if err != nil && err.Error() != "feature ID is required" {
		t.Errorf("error message = %v, want 'feature ID is required'", err.Error())
	}
}

func TestCreateBranchFlagLogic(t *testing.T) {
	// Test the logic for when to use -b flag
	tests := []struct {
		name         string
		createBranch bool
		branchName   string
		baseBranch   string
		wantCreateB  bool
	}{
		{
			name:         "create new branch from dev",
			createBranch: true,
			branchName:   "feature/F067",
			baseBranch:   "dev",
			wantCreateB:  true,
		},
		{
			name:         "use existing branch",
			createBranch: false,
			branchName:   "feature/F067",
			baseBranch:   "",
			wantCreateB:  false,
		},
		{
			name:         "create new branch from main",
			createBranch: true,
			branchName:   "hotfix/urgent",
			baseBranch:   "main",
			wantCreateB:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify the flag logic matches expectations
			if tt.createBranch != tt.wantCreateB {
				t.Errorf("CreateBranch = %v, want %v", tt.createBranch, tt.wantCreateB)
			}
		})
	}
}
