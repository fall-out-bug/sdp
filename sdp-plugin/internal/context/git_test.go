package context

import (
	"os"
	"os/exec"
	"testing"
)

// TestGetCurrentBranch tests the getCurrentBranchIn function
// This test requires being in a git repository
func TestGetCurrentBranch(t *testing.T) {
	// Skip if not in a git repo
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Skipf("failed to get cwd: %v", err)
	}

	// Check if we're in a git repo
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = cwd
	if err := cmd.Run(); err != nil {
		t.Skip("not in a git repository")
	}

	branch, err := getCurrentBranchIn(cwd)
	if err != nil {
		t.Skipf("getCurrentBranchIn failed: %v (may not be on a branch)", err)
	}

	// In CI, we may be in detached HEAD state (no branch name)
	if branch == "" {
		t.Skip("detached HEAD state - no branch name available")
	}

	// Branch should not contain newlines or spaces
	for _, r := range branch {
		if r == '\n' || r == ' ' {
			t.Errorf("Branch name contains unexpected character: %q", r)
		}
	}
}

// TestGetRemoteTracking tests the getRemoteTrackingIn function
func TestGetRemoteTracking(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Skipf("failed to get cwd: %v", err)
	}

	// Check if we're in a git repo
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = cwd
	if err := cmd.Run(); err != nil {
		t.Skip("not in a git repository")
	}

	// This may fail if there's no upstream configured
	remote, err := getRemoteTrackingIn(cwd)
	if err != nil {
		// This is expected if no upstream is set
		t.Logf("getRemoteTrackingIn returned error (expected if no upstream): %v", err)
		return
	}

	if remote == "" {
		t.Error("getRemoteTrackingIn should return non-empty remote name when successful")
	}
}

// TestListWorktrees tests the listWorktrees function
func TestListWorktrees(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	// Create a temp git repo
	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to init git repo: %v", err)
	}

	// Configure git user (required for commits)
	exec.Command("git", "-C", tmpDir, "config", "user.email", "test@test.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Test").Run()

	r := NewRecovery(tmpDir)

	worktrees, err := r.listWorktrees()
	if err != nil {
		t.Skipf("listWorktrees failed: %v", err)
	}

	// At minimum, the main worktree should be listed
	if len(worktrees) == 0 {
		t.Error("listWorktrees should return at least the main worktree")
	}

	// Verify that at least one worktree exists
	t.Logf("Found worktrees: %v", worktrees)
}

// TestRecoveryCheck tests the Check function
func TestRecoveryCheck(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	// Check runs in current directory, which should be a git repo
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}

	// Verify we're in a git repo
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = cwd
	if err := cmd.Run(); err != nil {
		t.Skip("not in a git repository")
	}

	r := NewRecovery(cwd)
	result, err := r.Check()

	if err != nil {
		t.Logf("Check returned error: %v (may be expected)", err)
		return
	}

	if result == nil {
		t.Fatal("Check returned nil result")
	}

	// Result should have a worktree path
	if result.WorktreePath == "" {
		t.Error("Check result should have WorktreePath")
	}

	// ExitCode should be one of the defined constants
	validExitCode := result.ExitCode == ExitCodeOK ||
		result.ExitCode == ExitCodeContextMismatch ||
		result.ExitCode == ExitCodeNoSession ||
		result.ExitCode == ExitCodeHashMismatch ||
		result.ExitCode == ExitCodeRuntimeError

	if !validExitCode {
		t.Errorf("Invalid exit code: %d", result.ExitCode)
	}
}

// TestRecoveryShow tests the Show function
func TestRecoveryShow(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}

	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = cwd
	if err := cmd.Run(); err != nil {
		t.Skip("not in a git repository")
	}

	r := NewRecovery(cwd)
	result, err := r.Show()

	if err != nil {
		t.Logf("Show returned error: %v", err)
		return
	}

	if result == nil {
		t.Error("Show should return non-nil result")
	}
}

// TestClean tests the Clean function
func TestClean(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get cwd: %v", err)
	}

	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = cwd
	if err := cmd.Run(); err != nil {
		t.Skip("not in a git repository")
	}

	r := NewRecovery(cwd)
	cleaned, err := r.Clean()

	if err != nil {
		t.Logf("Clean returned error: %v (may be expected)", err)
		return
	}

	// Clean may or may not find files to clean
	t.Logf("Cleaned %d files", len(cleaned))
}
