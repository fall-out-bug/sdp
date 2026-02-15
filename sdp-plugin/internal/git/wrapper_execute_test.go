package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/session"
)

func TestExecute_NeedsSessionCheck_NoSession(t *testing.T) {
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

	wrapper := NewWrapper(tmpDir)

	// 'commit' needs session check
	err := wrapper.Execute("commit", "-m", "test")
	if err == nil {
		t.Error("Execute should fail when session check needed but no session exists")
	}
}

func TestExecute_SafeCommand_NoSession(t *testing.T) {
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

	wrapper := NewWrapper(tmpDir)

	// 'status' is safe and doesn't need session check
	err := wrapper.Execute("status")
	// May fail due to no commits, but shouldn't fail due to session
	t.Logf("Execute(status) error: %v", err)
}

func TestExecute_WithValidSession(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	tmpDir := t.TempDir()

	// Get real path (resolve symlinks on macOS)
	realTmpDir, err := filepath.EvalSymlinks(tmpDir)
	if err != nil {
		t.Fatalf("Failed to eval symlinks: %v", err)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = realTmpDir
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to init git repo: %v", err)
	}

	exec.Command("git", "-C", realTmpDir, "config", "user.email", "test@test.com").Run()
	exec.Command("git", "-C", realTmpDir, "config", "user.name", "Test").Run()

	// Create initial commit
	cmd = exec.Command("git", "-C", realTmpDir, "commit", "--allow-empty", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to create initial commit: %v", err)
	}

	// Create and checkout feature branch
	cmd = exec.Command("git", "-C", realTmpDir, "checkout", "-b", "feature/F067")
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to checkout feature branch: %v", err)
	}

	// Create .sdp directory
	sdpDir := filepath.Join(realTmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Initialize valid session
	s, err := session.Init("F067", realTmpDir, "test")
	if err != nil {
		t.Fatalf("Failed to init session: %v", err)
	}

	s.ExpectedBranch = "feature/F067"
	if err := s.Save(realTmpDir); err != nil {
		t.Fatalf("Failed to save session: %v", err)
	}

	// Change to temp directory (use real path)
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(realTmpDir)

	wrapper := NewWrapper(realTmpDir)

	// Now execute should work with session check
	err = wrapper.Execute("status")
	t.Logf("Execute(status) with session error: %v", err)
}

func TestExecute_NonexistentGitCommand(t *testing.T) {
	tmpDir := t.TempDir()

	wrapper := NewWrapper(tmpDir)

	err := wrapper.Execute("nonexistent-command")
	if err == nil {
		t.Error("Execute should fail for nonexistent git command")
	}
}

func TestExecute_PostCheck(t *testing.T) {
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

	wrapper := NewWrapper(tmpDir)

	// 'checkout' needs post-check
	// This will fail since there's no commits, but tests the post-check path
	err := wrapper.Execute("checkout", "main")
	t.Logf("Execute(checkout) error: %v", err)
}

func TestWrapper_ValidatorIntegration(t *testing.T) {
	tmpDir := t.TempDir()

	wrapper := NewWrapper(tmpDir)

	// Wrapper should have a validator
	if wrapper.validator == nil {
		t.Error("wrapper.validator should not be nil")
	}

	if wrapper.validator.ProjectRoot != tmpDir {
		t.Errorf("validator.ProjectRoot = %s, want %s", wrapper.validator.ProjectRoot, tmpDir)
	}
}

func TestWrapper_GetWorktreePath_NoSession(t *testing.T) {
	tmpDir := t.TempDir()

	wrapper := NewWrapper(tmpDir)

	_, err := wrapper.GetWorktreePath()
	if err == nil {
		t.Error("GetWorktreePath should fail when no session exists")
	}

	t.Logf("Expected error: %v", err)
}

func TestWrapper_GetWorktreePath_WithSession(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	tmpDir := t.TempDir()

	realTmpDir, err := filepath.EvalSymlinks(tmpDir)
	if err != nil {
		t.Fatalf("Failed to eval symlinks: %v", err)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = realTmpDir
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to init git repo: %v", err)
	}

	// Create .sdp directory
	sdpDir := filepath.Join(realTmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Initialize session
	s, err := session.Init("F067", realTmpDir, "test")
	if err != nil {
		t.Fatalf("Failed to init session: %v", err)
	}

	if err := s.Save(realTmpDir); err != nil {
		t.Fatalf("Failed to save session: %v", err)
	}

	wrapper := NewWrapper(realTmpDir)

	path, err := wrapper.GetWorktreePath()
	if err != nil {
		t.Errorf("GetWorktreePath error: %v", err)
	}

	t.Logf("GetWorktreePath = %s", path)
}

func TestWrapper_HasSession_False(t *testing.T) {
	tmpDir := t.TempDir()

	wrapper := NewWrapper(tmpDir)

	if wrapper.HasSession() {
		t.Error("HasSession should return false when no session exists")
	}
}

func TestWrapper_HasSession_True(t *testing.T) {
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

	// Create .sdp directory
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Initialize session
	s, err := session.Init("F067", tmpDir, "test")
	if err != nil {
		t.Fatalf("Failed to init session: %v", err)
	}

	if err := s.Save(tmpDir); err != nil {
		t.Fatalf("Failed to save session: %v", err)
	}

	wrapper := NewWrapper(tmpDir)

	if !wrapper.HasSession() {
		t.Error("HasSession should return true when session exists")
	}
}

func TestExecute_CommandFailure(t *testing.T) {
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

	wrapper := NewWrapper(tmpDir)

	// Try to push to nonexistent remote - should fail
	err := wrapper.Execute("push", "nonexistent-remote", "main")
	if err == nil {
		t.Error("Execute should fail for invalid git command")
	}

	t.Logf("Expected error: %v", err)
}
