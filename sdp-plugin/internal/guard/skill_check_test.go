package guard

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewSkill(t *testing.T) {
	configDir := t.TempDir()
	skill := NewSkill(configDir)

	if skill == nil {
		t.Fatal("NewSkill returned nil")
	}
	if skill.stateManager == nil {
		t.Error("stateManager not initialized")
	}
	if skill.activeWS != "" {
		t.Errorf("activeWS should be empty, got %s", skill.activeWS)
	}
}

func TestActivate(t *testing.T) {
	configDir := t.TempDir()
	skill := NewSkill(configDir)

	// Test activation
	err := skill.Activate("00-001-01")
	if err != nil {
		t.Fatalf("Activate failed: %v", err)
	}

	// Verify state saved
	state, err := skill.stateManager.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if state.ActiveWS != "00-001-01" {
		t.Errorf("ActiveWS = %s, want 00-001-01", state.ActiveWS)
	}

	if state.ActivatedAt == "" {
		t.Error("ActivatedAt should be set")
	}

	// Verify timestamp is recent
	activatedAt, err := time.Parse(time.RFC3339, state.ActivatedAt)
	if err != nil {
		t.Fatalf("Failed to parse ActivatedAt: %v", err)
	}

	if time.Since(activatedAt) > 5*time.Second {
		t.Error("ActivatedAt should be recent")
	}
}

func TestCheckEditNoActiveWS(t *testing.T) {
	configDir := t.TempDir()
	skill := NewSkill(configDir)

	result, err := skill.CheckEdit("/some/file.go")
	if err != nil {
		t.Fatalf("CheckEdit failed: %v", err)
	}

	if result.Allowed {
		t.Error("Should not be allowed when no active WS")
	}

	if result.WSID != "" {
		t.Errorf("WSID should be empty, got %s", result.WSID)
	}
}

func TestCheckEditNoScope(t *testing.T) {
	configDir := t.TempDir()
	skill := NewSkill(configDir)

	// Activate WS without scope restrictions
	err := skill.Activate("00-001-01")
	if err != nil {
		t.Fatalf("Activate failed: %v", err)
	}

	// Any file should be allowed
	result, err := skill.CheckEdit("/any/file.go")
	if err != nil {
		t.Fatalf("CheckEdit failed: %v", err)
	}

	if !result.Allowed {
		t.Errorf("Should be allowed (no scope): %s", result.Reason)
	}

	if result.WSID != "00-001-01" {
		t.Errorf("WSID = %s, want 00-001-01", result.WSID)
	}
}

func TestCheckEditWithScope(t *testing.T) {
	configDir := t.TempDir()
	skill := NewSkill(configDir)

	// Activate WS with scope files
	err := skill.Activate("00-001-01")
	if err != nil {
		t.Fatalf("Activate failed: %v", err)
	}

	// Set scope files
	state, _ := skill.stateManager.Load()
	state.ScopeFiles = []string{
		"/allowed/file1.go",
		"/allowed/file2.go",
	}
	skill.stateManager.Save(*state)

	// Test allowed file
	result, err := skill.CheckEdit("/allowed/file1.go")
	if err != nil {
		t.Fatalf("CheckEdit failed: %v", err)
	}

	if !result.Allowed {
		t.Errorf("Should be allowed (in scope): %s", result.Reason)
	}

	// Test blocked file
	result, err = skill.CheckEdit("/blocked/file.go")
	if err != nil {
		t.Fatalf("CheckEdit failed: %v", err)
	}

	if result.Allowed {
		t.Error("Should not be allowed (not in scope)")
	}

	if len(result.ScopeFiles) != 2 {
		t.Errorf("ScopeFiles count = %d, want 2", len(result.ScopeFiles))
	}
}

func TestCheckEditExpired(t *testing.T) {
	configDir := t.TempDir()
	skill := NewSkill(configDir)

	// Activate WS
	err := skill.Activate("00-001-01")
	if err != nil {
		t.Fatalf("Activate failed: %v", err)
	}

	// Manually set old timestamp (>24 hours)
	state, _ := skill.stateManager.Load()
	oldTime := time.Now().Add(-25 * time.Hour)
	state.ActivatedAt = oldTime.Format(time.RFC3339)
	skill.stateManager.Save(*state)

	// Check should fail due to expiration
	result, err := skill.CheckEdit("/some/file.go")
	if err != nil {
		t.Fatalf("CheckEdit failed: %v", err)
	}

	if result.Allowed {
		t.Error("Should not be allowed (state expired)")
	}

	if result.WSID != "" {
		t.Errorf("WSID should be empty when expired, got %s", result.WSID)
	}
}

func TestGetActiveWS(t *testing.T) {
	configDir := t.TempDir()
	skill := NewSkill(configDir)

	// No active WS
	wsID := skill.GetActiveWS()
	if wsID != "" {
		t.Errorf("WSID should be empty, got %s", wsID)
	}

	// Activate WS
	err := skill.Activate("00-001-01")
	if err != nil {
		t.Fatalf("Activate failed: %v", err)
	}

	// Should return active WS
	wsID = skill.GetActiveWS()
	if wsID != "00-001-01" {
		t.Errorf("WSID = %s, want 00-001-01", wsID)
	}
}

func TestDeactivate(t *testing.T) {
	configDir := t.TempDir()
	skill := NewSkill(configDir)

	// Activate WS first
	err := skill.Activate("00-001-01")
	if err != nil {
		t.Fatalf("Activate failed: %v", err)
	}

	// Verify active
	if skill.GetActiveWS() != "00-001-01" {
		t.Error("WS should be active")
	}

	// Deactivate
	err = skill.Deactivate()
	if err != nil {
		t.Fatalf("Deactivate failed: %v", err)
	}

	// Verify inactive
	wsID := skill.GetActiveWS()
	if wsID != "" {
		t.Errorf("WSID should be empty after deactivate, got %s", wsID)
	}

	if skill.activeWS != "" {
		t.Errorf("activeWS should be empty, got %s", skill.activeWS)
	}
}

func TestResolvePathAbsolute(t *testing.T) {
	absPath := "/absolute/path/to/file.go"
	result, err := ResolvePath(absPath)
	if err != nil {
		t.Fatalf("ResolvePath failed: %v", err)
	}

	if result != absPath {
		t.Errorf("ResolvePath = %s, want %s", result, absPath)
	}
}

func TestResolvePathRelative(t *testing.T) {
	// Create temp directory and change to it
	tmpDir := t.TempDir()
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	os.Chdir(tmpDir)

	relPath := "relative/file.go"
	result, err := ResolvePath(relPath)
	if err != nil {
		t.Fatalf("ResolvePath failed: %v", err)
	}

	// Resolve symlinks for comparison (macOS /var -> /private/var)
	expected, _ := filepath.EvalSymlinks(filepath.Join(tmpDir, relPath))
	resultResolved, _ := filepath.EvalSymlinks(result)

	if resultResolved != expected {
		t.Errorf("ResolvePath = %s (resolved: %s), want %s", result, resultResolved, expected)
	}
}
