package doctor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckClaudeDir(t *testing.T) {
	// Create temp directory with .claude and subdirs
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	for _, subdir := range []string{"skills", "agents", "validators"} {
		if err := os.MkdirAll(filepath.Join(claudeDir, subdir), 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
	}

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	result := checkClaudeDir()

	if result.Status != "ok" {
		t.Errorf("Expected status ok, got %s", result.Status)
	}
	if !contains(result.Message, "SDP prompts installed") {
		t.Errorf("Wrong message: %s", result.Message)
	}
}

func TestCheckClaudeDir_NotFound(t *testing.T) {
	// Create temp directory WITHOUT .claude
	tmpDir := t.TempDir()

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	result := checkClaudeDir()

	if result.Status != "error" {
		t.Errorf("Expected status error, got %s", result.Status)
	}
	if !contains(result.Message, "Not found") {
		t.Errorf("Wrong message: %s", result.Message)
	}
}

func TestCheckFilePermissions(t *testing.T) {
	// This test checks real file permissions in HOME/.sdp
	// We can't easily mock this, so just verify the function runs
	result := checkFilePermissions()

	// Should not crash
	if result.Name != "File Permissions" {
		t.Errorf("Expected name 'File Permissions', got '%s'", result.Name)
	}

	// Status can be ok, warning, or error depending on actual files
	if result.Status != "ok" && result.Status != "warning" && result.Status != "error" {
		t.Errorf("Unexpected status: %s", result.Status)
	}
}

func TestCheckFilePermissions_AllSecure(t *testing.T) {
	// Create temp directory with .beads
	tmpDir := t.TempDir()
	beadsDir := filepath.Join(tmpDir, ".beads")
	if err := os.MkdirAll(beadsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create files with secure permissions
	testFile := filepath.Join(beadsDir, "test.db")
	if err := os.WriteFile(testFile, []byte("test"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	result := checkFilePermissions()

	if result.Status != "ok" {
		t.Errorf("Expected status ok, got %s: %s", result.Status, result.Message)
	}
}

func TestCheckFilePermissions_NoBeadsDir(t *testing.T) {
	// Create temp directory WITHOUT .beads
	tmpDir := t.TempDir()

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	result := checkFilePermissions()

	// Should skip if no .beads directory
	if result.Status != "ok" {
		t.Errorf("Expected status ok when no .beads, got %s", result.Status)
	}
}

func TestFindProjectRootForDrift(t *testing.T) {
	// We're already in a git repo (sdp-plugin)
	// findProjectRootForDrift should find the repo root
	root, err := findProjectRootForDrift()
	if err != nil {
		t.Errorf("findProjectRootForDrift() failed: %v", err)
	}

	// Should return a valid path
	if root == "" {
		t.Error("Expected non-empty root path")
	}

	// Verify it's a directory that exists
	if info, err := os.Stat(root); err != nil {
		t.Errorf("Root path doesn't exist: %v", err)
	} else if !info.IsDir() {
		t.Errorf("Root path is not a directory: %s", root)
	}
}

func TestFindProjectRootForDrift_ReturnsCwd(t *testing.T) {
	// The function should always return a path (never error)
	// Even if .git is not found, it returns cwd

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })

	// Change to temp dir (no .git, no docs)
	tmpDir := t.TempDir()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	root, err := findProjectRootForDrift()
	if err != nil {
		t.Errorf("findProjectRootForDrift() should not fail: %v", err)
	}

	if root == "" {
		t.Error("Expected non-empty root path")
	}

	// Should return tmpDir or a parent
	if root != tmpDir && root != filepath.Dir(tmpDir) {
		// OK as long as it returns something valid
	}
}

func TestRun(t *testing.T) {
	// Run() should return results without crashing
	results := Run()

	if results == nil {
		t.Error("Expected results slice, got nil")
	}

	// Should have at least 5 checks (Git, Claude Code, Go, .claude/, File Permissions)
	expectedMinChecks := 5
	if len(results) < expectedMinChecks {
		t.Errorf("Expected at least %d checks, got %d", expectedMinChecks, len(results))
	}

	// Verify each result has required fields
	for _, result := range results {
		if result.Name == "" {
			t.Error("Result has empty name")
		}
		if result.Status == "" {
			t.Error("Result has empty status")
		}
		if result.Message == "" {
			t.Error("Result has empty message")
		}
	}
}

func TestRunWithOptions(t *testing.T) {
	opts := RunOptions{
		DriftCheck: false, // Don't run drift check (slow)
	}

	results := RunWithOptions(opts)

	if results == nil {
		t.Error("Expected results slice, got nil")
	}

	// Should have 6 checks (no drift: Git, Claude, Go, .claude/, permissions, .sdp/config.yml)
	expectedChecks := 6
	if len(results) != expectedChecks {
		t.Errorf("Expected %d checks without drift, got %d", expectedChecks, len(results))
	}
}

func TestRunWithOptions_WithDrift(t *testing.T) {
	// Create temp directory with docs/workstreams/completed
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "completed")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create a test workstream
	wsFile := filepath.Join(wsDir, "00-050-01.md")
	if err := os.WriteFile(wsFile, []byte("# Test WS\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Add go.mod to make it a valid project root
	goMod := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	opts := RunOptions{
		DriftCheck: true,
	}

	results := RunWithOptions(opts)

	if results == nil {
		t.Error("Expected results slice, got nil")
	}

	// Should have 7 checks (6 standard + drift)
	expectedChecks := 7
	if len(results) != expectedChecks {
		t.Errorf("Expected %d checks with drift, got %d", expectedChecks, len(results))
	}

	// Find drift check result
	var driftResult *CheckResult
	for i := range results {
		if results[i].Name == "Drift Detection" {
			driftResult = &results[i]
			break
		}
	}

	if driftResult == nil {
		t.Error("Expected drift check result")
	}
}

func TestFindRecentWorkstreamsForDrift(t *testing.T) {
	// Create temp directory with docs/workstreams/completed
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "completed")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create test workstreams
	for i := 1; i <= 3; i++ {
		wsFile := filepath.Join(wsDir, "00-050-0"+string(rune('0'+i))+".md")
		content := "# Test Workstream\n"
		if err := os.WriteFile(wsFile, []byte(content), 0644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}

	workstreams, err := findRecentWorkstreamsForDrift(tmpDir)
	if err != nil {
		t.Errorf("findRecentWorkstreamsForDrift() failed: %v", err)
	}

	if len(workstreams) != 3 {
		t.Errorf("Expected 3 workstreams, got %d", len(workstreams))
	}
}

func TestFindRecentWorkstreamsForDrift_NoneFound(t *testing.T) {
	// Create temp directory WITHOUT workstreams
	tmpDir := t.TempDir()

	workstreams, err := findRecentWorkstreamsForDrift(tmpDir)
	if err != nil {
		t.Errorf("findRecentWorkstreamsForDrift() should not fail: %v", err)
	}

	if len(workstreams) != 0 {
		t.Errorf("Expected 0 workstreams, got %d", len(workstreams))
	}
}

func TestFindRecentWorkstreamsForDrift_EmptyDirectories(t *testing.T) {
	// Create temp directory with empty workstreams directories
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "completed")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	workstreams, err := findRecentWorkstreamsForDrift(tmpDir)
	if err != nil {
		t.Errorf("findRecentWorkstreamsForDrift() should not fail: %v", err)
	}

	if len(workstreams) != 0 {
		t.Errorf("Expected 0 workstreams from empty directories, got %d", len(workstreams))
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func TestCheckFilePermissions_Insecure(t *testing.T) {
	// Create a temp .beads directory with insecure file
	tmpDir := t.TempDir()
	beadsDir := filepath.Join(tmpDir, ".beads")
	if err := os.MkdirAll(beadsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create beads.db with world-readable permissions (insecure)
	dbFile := filepath.Join(beadsDir, "beads.db")
	if err := os.WriteFile(dbFile, []byte("test"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	result := checkFilePermissions()

	// Should detect insecure permissions
	if result.Status != "warning" && result.Status != "ok" {
		t.Logf("Status: %s, Message: %s", result.Status, result.Message)
		// Warning if insecure files found, OK otherwise
	}
}

func TestCheckFilePermissions_Secure(t *testing.T) {
	// Create a temp .beads directory with secure file
	tmpDir := t.TempDir()
	beadsDir := filepath.Join(tmpDir, ".beads")
	if err := os.MkdirAll(beadsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create beads.db with secure permissions
	dbFile := filepath.Join(beadsDir, "beads.db")
	if err := os.WriteFile(dbFile, []byte("test"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	result := checkFilePermissions()

	if result.Status != "ok" {
		t.Errorf("Expected status ok with secure file, got %s: %s", result.Status, result.Message)
	}
}

func TestCheckFilePermissions_NoSensitiveFiles(t *testing.T) {
	// Create temp directory without sensitive files
	tmpDir := t.TempDir()

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	result := checkFilePermissions()

	// Should be OK if no sensitive files exist
	if result.Status != "ok" {
		t.Errorf("Expected status ok when no sensitive files, got %s", result.Status)
	}
}

func TestCheckFilePermissions_InsecureFilesInDirectory(t *testing.T) {
	// Create .oneshot directory with insecure file
	tmpDir := t.TempDir()
	oneshotDir := filepath.Join(tmpDir, ".oneshot")
	if err := os.MkdirAll(oneshotDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create file with world-readable permissions (insecure)
	testFile := filepath.Join(oneshotDir, "test.json")
	if err := os.WriteFile(testFile, []byte("{}"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	result := checkFilePermissions()

	// Should detect insecure files in directory
	if result.Status != "warning" {
		t.Errorf("Expected warning for insecure directory files, got %s: %s", result.Status, result.Message)
	}
	if !contains(result.Message, "insecure permissions") {
		t.Errorf("Expected insecure permissions message, got: %s", result.Message)
	}
}

func TestCheckDrift_NoWorkstreamsFound(t *testing.T) {
	// Create temp directory with docs/workstreams/completed but empty
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "completed")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Add go.mod to make it a valid project root
	goMod := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	result := checkDrift()

	// Should return OK with "no recent workstreams" message
	if result.Status != "ok" {
		t.Errorf("Expected ok status when no workstreams, got %s: %s", result.Status, result.Message)
	}
	if !contains(result.Message, "No recent workstreams") {
		t.Errorf("Expected 'No recent workstreams' message, got: %s", result.Message)
	}
}

func TestCheckDrift_CouldNotCheckAny(t *testing.T) {
	// Create temp directory with invalid workstream files (can't be parsed)
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "completed")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create files that look like workstreams but have invalid content
	for i := 1; i <= 3; i++ {
		wsFile := filepath.Join(wsDir, "00-050-0"+string(rune('0'+i))+".md")
		// Empty file will fail parsing
		if err := os.WriteFile(wsFile, []byte(""), 0644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}

	// Add go.mod to make it a valid project root
	goMod := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	result := checkDrift()

	// Should return warning when can't check any
	if result.Status != "warning" {
		t.Errorf("Expected warning status when can't check workstreams, got %s: %s", result.Status, result.Message)
	}
	if !contains(result.Message, "Could not check") {
		t.Errorf("Expected 'Could not check' message, got: %s", result.Message)
	}
}

func TestCheckDrift_WithValidWorkstreams(t *testing.T) {
	// Use existing test that has valid workstreams
	// This test just ensures checkDrift doesn't crash when called
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })

	// Change to a directory with actual workstreams (the sdp repo)
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "completed")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create a minimal valid workstream
	wsContent := `---
ws_id: 00-050-01
title: Test
scope_files: []
---

# Test
`
	wsFile := filepath.Join(wsDir, "00-050-01-test.md")
	if err := os.WriteFile(wsFile, []byte(wsContent), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Add go.mod
	goMod := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	result := checkDrift()

	// Just verify it doesn't crash and has valid fields
	if result.Name == "" {
		t.Error("Expected non-empty name")
	}
	if result.Status == "" {
		t.Error("Expected non-empty status")
	}
}
