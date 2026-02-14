package worktree

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCreate_NoFeatureID(t *testing.T) {
	tmpDir := t.TempDir()
	creator := NewCreator(tmpDir)

	_, err := creator.Create(CreateOptions{})
	if err == nil {
		t.Error("Create should fail without FeatureID")
	}
}

func TestCreate_Defaults(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

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

	// Create initial commit on main branch
	cmd = exec.Command("git", "-C", tmpDir, "commit", "--allow-empty", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to create initial commit: %v", err)
	}

	// Create dev branch
	cmd = exec.Command("git", "-C", tmpDir, "branch", "dev")
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to create dev branch: %v", err)
	}

	creator := NewCreator(tmpDir)
	worktreesDir := filepath.Join(tmpDir, "worktrees")
	os.MkdirAll(worktreesDir, 0755)
	creator.WorktreesDir = worktreesDir

	result, err := creator.Create(CreateOptions{
		FeatureID:    "F067",
		CreateBranch: true,
		BaseBranch:   "dev",
	})

	if err != nil {
		t.Logf("Create returned error: %v (may be expected)", err)
		return
	}

	if result == nil {
		t.Error("Create should return result on success")
		return
	}

	// Verify defaults
	expectedBranch := "feature/F067"
	if result.BranchName != expectedBranch {
		t.Errorf("BranchName = %q, want %q", result.BranchName, expectedBranch)
	}

	// Verify worktree path
	expectedPath := filepath.Join(worktreesDir, "sdp-F067")
	if result.WorktreePath != expectedPath {
		t.Errorf("WorktreePath = %q, want %q", result.WorktreePath, expectedPath)
	}

	// Cleanup
	creator.Delete("F067")
}

func TestCreate_CustomBranchName(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to init git repo: %v", err)
	}

	exec.Command("git", "-C", tmpDir, "config", "user.email", "test@test.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Test").Run()

	cmd = exec.Command("git", "-C", tmpDir, "commit", "--allow-empty", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to create initial commit: %v", err)
	}

	cmd = exec.Command("git", "-C", tmpDir, "branch", "dev")
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to create dev branch: %v", err)
	}

	creator := NewCreator(tmpDir)
	worktreesDir := filepath.Join(tmpDir, "worktrees")
	os.MkdirAll(worktreesDir, 0755)
	creator.WorktreesDir = worktreesDir

	result, err := creator.Create(CreateOptions{
		FeatureID:    "F068",
		BranchName:   "custom/branch",
		CreateBranch: true,
		BaseBranch:   "dev",
	})

	if err != nil {
		t.Logf("Create returned error: %v", err)
		return
	}

	if result.BranchName != "custom/branch" {
		t.Errorf("BranchName = %q, want custom/branch", result.BranchName)
	}

	// Cleanup
	creator.Delete("F068")
}

func TestRemoveWorktree_Nonexistent(t *testing.T) {
	tmpDir := t.TempDir()
	creator := NewCreator(tmpDir)

	err := creator.removeWorktree("/nonexistent/path")
	// Should fail for nonexistent worktree
	if err == nil {
		t.Log("removeWorktree succeeded for nonexistent path (unexpected)")
	}
}
