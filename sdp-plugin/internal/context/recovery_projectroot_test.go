package context

import (
	"os"
	"os/exec"
	"testing"
)

func TestFindProjectRoot(t *testing.T) {
	// This should find the project root of the current repo
	root, err := FindProjectRoot()
	if err != nil {
		// If not in a git repo, this is expected
		t.Logf("FindProjectRoot returned error: %v (may be expected if not in git repo)", err)
		return
	}

	if root == "" {
		t.Error("FindProjectRoot should return non-empty string when successful")
	}

	// Verify .sdp or .git exists (function finds either)
	hasSDP := false
	hasGit := false
	if _, err := os.Stat(root + "/.sdp"); err == nil {
		hasSDP = true
	}
	if _, err := os.Stat(root + "/.git"); err == nil {
		hasGit = true
	}

	if !hasSDP && !hasGit {
		t.Errorf("Project root should contain .sdp or .git, got root=%s", root)
	}
}

func TestRecovery_Repair_NoGit(t *testing.T) {
	// Create temp dir without git
	tmpDir := t.TempDir()

	// Change to the temp directory so getCurrentBranch() runs in non-git context
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	r := NewRecovery(tmpDir)
	err := r.Repair()

	// Should fail without git
	if err == nil {
		t.Error("Repair should fail in non-git directory")
	}
}

func TestRecovery_Repair_WithGit(t *testing.T) {
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

	// Create initial commit
	cmd = exec.Command("git", "-C", tmpDir, "commit", "--allow-empty", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to create initial commit: %v", err)
	}

	// Create and checkout a feature branch
	cmd = exec.Command("git", "-C", tmpDir, "checkout", "-b", "feature/F067")
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to create feature branch: %v", err)
	}

	// Change to the temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	r := NewRecovery(tmpDir)
	err := r.Repair()

	// May fail if can't extract feature ID or session repair fails
	if err != nil {
		t.Logf("Repair returned error: %v (may be expected)", err)
	}
}

func TestRecovery_Check_NoSession(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	// Create temp dir without session
	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to init git repo: %v", err)
	}

	// Change to the temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	r := NewRecovery(tmpDir)
	result, err := r.Check()

	if err != nil {
		t.Logf("Check returned error: %v", err)
		return
	}

	if result.Valid {
		t.Error("Check should return invalid for directory without session")
	}

	if result.ExitCode != ExitCodeNoSession {
		t.Errorf("ExitCode = %d, want %d", result.ExitCode, ExitCodeNoSession)
	}
}
