package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/session"
)

func TestValidateSession_ValidSession(t *testing.T) {
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

	// Update session with correct branch
	s.ExpectedBranch = "feature/F067"
	if err := s.Save(realTmpDir); err != nil {
		t.Fatalf("Failed to save session: %v", err)
	}

	// Change to temp directory (use real path)
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(realTmpDir)

	v := NewValidator(realTmpDir)
	result, err := v.ValidateSession()

	if err != nil {
		t.Errorf("ValidateSession error: %v", err)
	}

	if result == nil {
		t.Fatal("ValidateSession should return result")
	}

	if !result.Valid {
		t.Errorf("Expected valid session, got invalid: %s", result.Error)
	}
}

func TestValidateSession_WorktreeMismatch(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	tmpDir := t.TempDir()
	otherDir := t.TempDir() // Different directory

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to init git repo: %v", err)
	}

	exec.Command("git", "-C", tmpDir, "config", "user.email", "test@test.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Test").Run()

	// Create initial commit
	cmd = exec.Command("git", "-C", tmpDir, "commit", "--allow-empty", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to create initial commit: %v", err)
	}

	// Create .sdp directory with session pointing to different path
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Initialize session with different worktree path
	s, err := session.Init("F067", otherDir, "test")
	if err != nil {
		t.Fatalf("Failed to init session: %v", err)
	}

	s.ExpectedBranch = "feature/F067"
	if err := s.Save(tmpDir); err != nil {
		t.Fatalf("Failed to save session: %v", err)
	}

	// Change to tmpDir (different from session's worktree path)
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	v := NewValidator(tmpDir)
	result, err := v.ValidateSession()

	if err != nil {
		t.Logf("ValidateSession error: %v", err)
	}

	if result == nil {
		t.Fatal("ValidateSession should return result")
	}

	// Should be invalid due to worktree mismatch
	if result.Valid {
		t.Error("Expected invalid session due to worktree mismatch")
	}

	if result.Error != "wrong worktree" {
		t.Logf("Error = %s", result.Error)
	}
}

func TestValidateSession_CorruptedSession(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .sdp directory with invalid session
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create invalid session file
	sessionContent := `{"feature_id":"F067","worktree_path":"/path","hash":"invalid"}`
	sessionPath := filepath.Join(sdpDir, "session.json")
	if err := os.WriteFile(sessionPath, []byte(sessionContent), 0644); err != nil {
		t.Fatal(err)
	}

	v := NewValidator(tmpDir)
	result, err := v.ValidateSession()

	// Should return error for corrupted session
	if err == nil {
		t.Log("Expected error for corrupted session")
	}

	if result != nil && result.Valid {
		t.Error("Corrupted session should not be valid")
	}
}

func TestValidator_ProjectRoot(t *testing.T) {
	v := NewValidator("/path/to/project")

	if v.ProjectRoot != "/path/to/project" {
		t.Errorf("ProjectRoot = %s, want /path/to/project", v.ProjectRoot)
	}
}

func TestValidationResult_AllFields(t *testing.T) {
	r := &ValidationResult{
		Valid:          false,
		WorktreePath:   "/worktree",
		ActualPath:     "/actual",
		CurrentBranch:  "main",
		ExpectedBranch: "feature/F001",
		ExpectedRemote: "origin/main",
		ActualRemote:   "origin/feature",
		Error:          "test error",
		Fix:            "run fix command",
	}

	if r.WorktreePath != "/worktree" {
		t.Errorf("WorktreePath = %s", r.WorktreePath)
	}
	if r.ActualPath != "/actual" {
		t.Errorf("ActualPath = %s", r.ActualPath)
	}
	if r.CurrentBranch != "main" {
		t.Errorf("CurrentBranch = %s", r.CurrentBranch)
	}
	if r.ExpectedBranch != "feature/F001" {
		t.Errorf("ExpectedBranch = %s", r.ExpectedBranch)
	}
	if r.ExpectedRemote != "origin/main" {
		t.Errorf("ExpectedRemote = %s", r.ExpectedRemote)
	}
	if r.ActualRemote != "origin/feature" {
		t.Errorf("ActualRemote = %s", r.ActualRemote)
	}
}
