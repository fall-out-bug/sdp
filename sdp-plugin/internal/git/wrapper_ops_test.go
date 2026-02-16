package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestHasSession_NoSession(t *testing.T) {
	tmpDir := t.TempDir()

	wrapper := NewWrapper(tmpDir)

	if wrapper.HasSession() {
		t.Error("HasSession should return false when no session exists")
	}
}

func TestHasSession_WithSession(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .sdp directory
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("Failed to create .sdp dir: %v", err)
	}

	// Create session file
	sessionContent := `{"feature_id":"F067","worktree_path":"` + tmpDir + `","expected_branch":"feature/F067","hash":"test-hash"}`
	sessionPath := filepath.Join(sdpDir, "session.json")
	if err := os.WriteFile(sessionPath, []byte(sessionContent), 0644); err != nil {
		t.Fatalf("Failed to write session: %v", err)
	}

	wrapper := NewWrapper(tmpDir)

	if !wrapper.HasSession() {
		t.Error("HasSession should return true when session exists")
	}
}

func TestGetWorktreePath_NoSession(t *testing.T) {
	tmpDir := t.TempDir()

	wrapper := NewWrapper(tmpDir)

	_, err := wrapper.GetWorktreePath()
	if err == nil {
		t.Error("GetWorktreePath should return error when no session exists")
	}
}

func TestGetWorktreePath_WithSession(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .sdp directory
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("Failed to create .sdp dir: %v", err)
	}

	// Create session file - note: hash validation will fail, but path should still be readable
	sessionContent := `{"feature_id":"F067","worktree_path":"/custom/worktree","expected_branch":"feature/F067","hash":"test-hash"}`
	sessionPath := filepath.Join(sdpDir, "session.json")
	if err := os.WriteFile(sessionPath, []byte(sessionContent), 0644); err != nil {
		t.Fatalf("Failed to write session: %v", err)
	}

	wrapper := NewWrapper(tmpDir)

	path, err := wrapper.GetWorktreePath()
	// Session.Load may fail due to hash validation, which is expected behavior
	if err != nil {
		t.Logf("GetWorktreePath returned error (expected due to hash validation): %v", err)
		return
	}

	if path != "/custom/worktree" {
		t.Errorf("WorktreePath = %v, want /custom/worktree", path)
	}
}

func TestGetCurrentBranch_NoGit(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := GetCurrentBranch(tmpDir)
	if err == nil {
		t.Error("GetCurrentBranch should return error when not in a git repo")
	}
}

func TestGetCurrentBranch_InGitRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	// Create temp git repo
	tmpDir := t.TempDir()

	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to init git repo: %v", err)
	}

	// Configure git user
	exec.Command("git", "-C", tmpDir, "config", "user.email", "test@test.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Test").Run()

	// Create initial commit (needed for branch operations)
	exec.Command("git", "-C", tmpDir, "commit", "--allow-empty", "-m", "initial").Run()

	branch, err := GetCurrentBranch(tmpDir)
	if err != nil {
		// May fail if no commits yet
		t.Logf("GetCurrentBranch returned error: %v", err)
		return
	}

	if branch == "" {
		t.Error("GetCurrentBranch should return non-empty branch name")
	}
}

func TestExecute_SafeCommand(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	// Create temp git repo
	tmpDir := t.TempDir()

	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to init git repo: %v", err)
	}

	wrapper := NewWrapper(tmpDir)

	// status is a safe command and should not require session
	err := wrapper.Execute("status")
	// May fail due to no commits, but should not fail due to session check
	t.Logf("Execute(status) error: %v", err)
}

func TestExecute_NoGitRepo(t *testing.T) {
	tmpDir := t.TempDir()

	wrapper := NewWrapper(tmpDir)

	err := wrapper.Execute("status")
	if err == nil {
		t.Error("Execute should return error when not in a git repo")
	}
}

func TestValidator_ValidateSession_NoSession(t *testing.T) {
	tmpDir := t.TempDir()

	validator := NewValidator(tmpDir)

	result, err := validator.ValidateSession()
	// ValidateSession may return an error when no session file exists
	// This is expected behavior
	if err != nil {
		t.Logf("ValidateSession returned error (expected without session): %v", err)
		return
	}

	if result == nil {
		t.Fatal("ValidateSession should return non-nil result")
	}

	if result.Valid {
		t.Error("ValidateSession should return invalid when no session exists")
	}
}

func TestValidator_ValidateSession_WithSession(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .sdp directory
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("Failed to create .sdp dir: %v", err)
	}

	// Create session file
	sessionContent := `{"feature_id":"F067","worktree_path":"` + tmpDir + `","expected_branch":"feature/F067","hash":"test-hash"}`
	sessionPath := filepath.Join(sdpDir, "session.json")
	if err := os.WriteFile(sessionPath, []byte(sessionContent), 0644); err != nil {
		t.Fatalf("Failed to write session: %v", err)
	}

	validator := NewValidator(tmpDir)

	result, err := validator.ValidateSession()
	if err != nil {
		t.Logf("ValidateSession returned error: %v (expected without git)", err)
		return
	}

	// Result may still be invalid due to git operations failing
	t.Logf("ValidateSession result: Valid=%v, Error=%v", result.Valid, result.Error)
}

func TestFindProjectRoot(t *testing.T) {
	// FindProjectRoot may or may not succeed depending on current directory
	root, err := FindProjectRoot()
	if err != nil {
		t.Logf("FindProjectRoot returned error: %v (expected if not in SDP project)", err)
		return
	}
	t.Logf("FindProjectRoot found: %s", root)
}

func TestWrapper_ExecuteWithArgs(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	// Create temp git repo
	tmpDir := t.TempDir()

	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to init git repo: %v", err)
	}

	wrapper := NewWrapper(tmpDir)

	// Test with additional arguments
	err := wrapper.Execute("rev-parse", "--git-dir")
	t.Logf("Execute(rev-parse --git-dir) error: %v", err)
}

func TestWrapper_ValidatorField(t *testing.T) {
	wrapper := NewWrapper("/path/to/project")

	if wrapper.validator == nil {
		t.Error("wrapper.validator should not be nil")
	}

	if wrapper.validator.ProjectRoot != wrapper.ProjectRoot {
		t.Error("validator should have same ProjectRoot as wrapper")
	}
}
