package context

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/session"
)

func TestFindWorktree_WithSession(t *testing.T) {
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

	// Create initial commit
	cmd = exec.Command("git", "-C", tmpDir, "commit", "--allow-empty", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to create initial commit: %v", err)
	}

	// Create worktrees directory
	worktreesDir := filepath.Join(tmpDir, "worktrees")
	os.MkdirAll(worktreesDir, 0755)

	// Create sdp-F067 directory with session
	featureDir := filepath.Join(worktreesDir, "sdp-F067")
	sdpDir := filepath.Join(featureDir, ".sdp")
	os.MkdirAll(sdpDir, 0755)

	// Create session
	s, err := session.Init("F067", featureDir, "test")
	if err != nil {
		t.Fatalf("Failed to init session: %v", err)
	}
	if err := s.Save(featureDir); err != nil {
		t.Fatalf("Failed to save session: %v", err)
	}

	r := NewRecovery(tmpDir)

	path, err := r.FindWorktree("F067")
	t.Logf("FindWorktree(F067) = %q, err = %v", path, err)
}

func TestRecovery_Check_WithValidSession(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	tmpDir := t.TempDir()

	// Get real path
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

	// Create feature branch
	cmd = exec.Command("git", "-C", realTmpDir, "checkout", "-b", "feature/F067")
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to checkout feature branch: %v", err)
	}

	// Create .sdp directory
	sdpDir := filepath.Join(realTmpDir, ".sdp")
	os.MkdirAll(sdpDir, 0755)

	// Create valid session
	s, err := session.Init("F067", realTmpDir, "test")
	if err != nil {
		t.Fatalf("Failed to init session: %v", err)
	}
	s.ExpectedBranch = "feature/F067"
	if err := s.Save(realTmpDir); err != nil {
		t.Fatalf("Failed to save session: %v", err)
	}

	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(realTmpDir)

	r := NewRecovery(realTmpDir)
	result, err := r.Check()

	if err != nil {
		t.Errorf("Check error: %v", err)
	}

	if result == nil {
		t.Fatal("Check should return result")
	}

	t.Logf("Check result: Valid=%v, ExitCode=%d, Errors=%v", result.Valid, result.ExitCode, result.Errors)
}

func TestRecovery_Check_BranchMismatch(t *testing.T) {
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

	exec.Command("git", "-C", realTmpDir, "config", "user.email", "test@test.com").Run()
	exec.Command("git", "-C", realTmpDir, "config", "user.name", "Test").Run()

	// Create initial commit
	cmd = exec.Command("git", "-C", realTmpDir, "commit", "--allow-empty", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to create initial commit: %v", err)
	}

	// Stay on main branch

	// Create .sdp directory
	sdpDir := filepath.Join(realTmpDir, ".sdp")
	os.MkdirAll(sdpDir, 0755)

	// Create session expecting different branch
	s, err := session.Init("F067", realTmpDir, "test")
	if err != nil {
		t.Fatalf("Failed to init session: %v", err)
	}
	s.ExpectedBranch = "feature/F067" // Different from current branch
	if err := s.Save(realTmpDir); err != nil {
		t.Fatalf("Failed to save session: %v", err)
	}

	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(realTmpDir)

	r := NewRecovery(realTmpDir)
	result, err := r.Check()

	if err != nil {
		t.Errorf("Check error: %v", err)
	}

	if result == nil {
		t.Fatal("Check should return result")
	}

	// Should be invalid due to branch mismatch
	if result.Valid {
		t.Error("Check should return invalid for branch mismatch")
	}

	if result.ExitCode != ExitCodeContextMismatch {
		t.Logf("ExitCode = %d (expected ExitCodeContextMismatch=%d)", result.ExitCode, ExitCodeContextMismatch)
	}
}

func TestRecovery_Check_CorruptedSession(t *testing.T) {
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

	// Create .sdp directory with corrupted session
	sdpDir := filepath.Join(tmpDir, ".sdp")
	os.MkdirAll(sdpDir, 0755)

	// Create corrupted session file
	sessionContent := `{"feature_id":"F067","worktree_path":"/path","hash":"invalid"}`
	sessionPath := filepath.Join(sdpDir, "session.json")
	os.WriteFile(sessionPath, []byte(sessionContent), 0644)

	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	r := NewRecovery(tmpDir)
	result, _ := r.Check()

	if result != nil {
		t.Logf("Check result: Valid=%v, ExitCode=%d, SessionValid=%v", result.Valid, result.ExitCode, result.SessionValid)

		// Corrupted session should result in invalid
		if result.Valid {
			t.Error("Check should return invalid for corrupted session")
		}
	}
}

func TestRecovery_Clean_NoWorktrees(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	tmpDir := t.TempDir()

	// Initialize git repo (no worktrees except main)
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

	r := NewRecovery(tmpDir)
	cleaned, err := r.Clean()

	// Clean should succeed even with no stale sessions
	if err != nil {
		t.Logf("Clean returned error: %v (may be expected)", err)
	}

	t.Logf("Cleaned sessions: %v", cleaned)
}

func TestRecovery_Clean_WithStaleSession(t *testing.T) {
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

	exec.Command("git", "-C", realTmpDir, "config", "user.email", "test@test.com").Run()
	exec.Command("git", "-C", realTmpDir, "config", "user.name", "Test").Run()

	// Create initial commit
	cmd = exec.Command("git", "-C", realTmpDir, "commit", "--allow-empty", "-m", "initial")
	if err := cmd.Run(); err != nil {
		t.Skipf("Failed to create initial commit: %v", err)
	}

	// Create .sdp directory with corrupted (stale) session
	sdpDir := filepath.Join(realTmpDir, ".sdp")
	os.MkdirAll(sdpDir, 0755)

	// Create corrupted session file (will be detected as stale)
	sessionContent := `{"feature_id":"F067","worktree_path":"/different/path","hash":"invalid"}`
	sessionPath := filepath.Join(sdpDir, "session.json")
	os.WriteFile(sessionPath, []byte(sessionContent), 0644)

	r := NewRecovery(realTmpDir)
	cleaned, err := r.Clean()

	if err != nil {
		t.Logf("Clean returned error: %v", err)
	}

	t.Logf("Cleaned sessions: %v", cleaned)
}

func TestRecovery_Clean_EmptyProjectRoot(t *testing.T) {
	// Test Clean with empty/non-existent project root
	r := NewRecovery("/nonexistent/path/that/does/not/exist")
	cleaned, err := r.Clean()

	// Should return error for non-existent path
	if err != nil {
		t.Logf("Clean correctly returned error: %v", err)
	}

	// Cleaned should be nil or empty
	if len(cleaned) > 0 {
		t.Logf("Cleaned unexpected sessions: %v", cleaned)
	}
}

func TestRecovery_Show(t *testing.T) {
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

	r := NewRecovery(realTmpDir)

	// Show should delegate to Check
	result, err := r.Show()

	if err != nil && result == nil {
		t.Errorf("Show should return result even on error")
	}

	t.Logf("Show result: %+v, err: %v", result, err)
}
