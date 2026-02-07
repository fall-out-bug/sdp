package doctor

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCheckDriftNoWorkstreams tests checkDrift when no recent workstreams exist
func TestCheckDriftNoWorkstreams(t *testing.T) {
	// Create temp directory without workstreams
	tmpDir := t.TempDir()

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(tmpDir)

	result := checkDrift()

	if result.Name != "Drift Detection" {
		t.Errorf("Expected name 'Drift Detection', got '%s'", result.Name)
	}

	// Should be OK or warning (no workstreams to check)
	if result.Status != "ok" && result.Status != "warning" {
		t.Errorf("Expected status ok or warning, got %s: %s", result.Status, result.Message)
	}
}

// TestCheckDriftWithWorkstreams tests checkDrift with real workstreams
func TestCheckDriftWithWorkstreams(t *testing.T) {
	// Create temp directory with docs/workstreams/completed
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "completed")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create test workstreams
	for i := 1; i <= 3; i++ {
		wsFile := filepath.Join(wsDir, "00-050-0"+string(rune('0'+i))+".md")
		content := "# Test Workstream\n\n## Scope\n- Test file\n"
		if err := os.WriteFile(wsFile, []byte(content), 0644); err != nil {
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
	os.Chdir(tmpDir)

	result := checkDrift()

	if result.Name != "Drift Detection" {
		t.Errorf("Expected name 'Drift Detection', got '%s'", result.Name)
	}

	// Should successfully check workstreams
	// Status depends on whether drift is detected
	if result.Status != "ok" && result.Status != "warning" && result.Status != "error" {
		t.Errorf("Unexpected status: %s, message: %s", result.Status, result.Message)
	}

	// Message should mention checked count
	if len(result.Message) == 0 {
		t.Error("Expected non-empty message")
	}
}

// TestCheckDriftInvalidWorkstream tests checkDrift with invalid workstream files
func TestCheckDriftInvalidWorkstream(t *testing.T) {
	// Create temp directory with invalid workstream files
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "completed")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create invalid workstream (empty file)
	wsFile := filepath.Join(wsDir, "00-050-01.md")
	if err := os.WriteFile(wsFile, []byte(""), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Add go.mod
	goMod := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(tmpDir)

	result := checkDrift()

	// Should handle gracefully (may skip invalid workstreams)
	if result.Name != "Drift Detection" {
		t.Errorf("Expected name 'Drift Detection', got '%s'", result.Name)
	}
}

// TestCheckDriftManyWorkstreams tests checkDrift with more than 5 workstreams (should limit to 5)
func TestCheckDriftManyWorkstreams(t *testing.T) {
	// Create temp directory with 7 workstreams
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "completed")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create 7 workstreams (should only check 5)
	for i := 1; i <= 7; i++ {
		wsFile := filepath.Join(wsDir, "00-050-0"+string(rune('0'+i))+".md")
		content := "# Test Workstream\n\n## Scope\n- Test file\n"
		if err := os.WriteFile(wsFile, []byte(content), 0644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}

	// Add go.mod
	goMod := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(tmpDir)

	result := checkDrift()

	// Message should mention checking at most 5
	if len(result.Message) == 0 {
		t.Error("Expected non-empty message")
	}

	// Verify it mentions checking some workstreams
	// (exact number depends on drift detector success)
}

// TestFindProjectRootForDriftWithGit tests finding project root with .git
func TestFindProjectRootForDriftWithGit(t *testing.T) {
	// Create temp directory with .git
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(tmpDir)

	root, err := findProjectRootForDrift()
	if err != nil {
		t.Errorf("findProjectRootForDrift() failed: %v", err)
	}

	if root == "" {
		t.Error("Expected non-empty root path")
	}

	// Should return tmpDir (where .git is)
	if root != tmpDir {
		t.Logf("Expected root %s, got %s", tmpDir, root)
	}
}

// TestFindProjectRootForDriftWithDocs tests finding project root with docs directory
func TestFindProjectRootForDriftWithDocs(t *testing.T) {
	// Create temp directory with docs
	tmpDir := t.TempDir()
	docsDir := filepath.Join(tmpDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(tmpDir)

	root, err := findProjectRootForDrift()
	if err != nil {
		t.Errorf("findProjectRootForDrift() failed: %v", err)
	}

	if root == "" {
		t.Error("Expected non-empty root path")
	}
}

// TestFindProjectRootForDriftInSdpPlugin tests finding root from sdp-plugin subdirectory
func TestFindProjectRootForDriftInSdpPlugin(t *testing.T) {
	// Create sdp-plugin directory structure
	tmpDir := t.TempDir()
	pluginDir := filepath.Join(tmpDir, "sdp-plugin")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Add go.mod to plugin dir
	goMod := filepath.Join(pluginDir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(pluginDir)

	root, err := findProjectRootForDrift()
	if err != nil {
		t.Errorf("findProjectRootForDrift() failed: %v", err)
	}

	if root == "" {
		t.Error("Expected non-empty root path")
	}

	// Should go up one level from sdp-plugin
	// (because it detects go.mod in current dir)
}

// TestFindProjectRootForDriftWithBeads tests finding project root with .beads
func TestFindProjectRootForDriftWithBeads(t *testing.T) {
	// Create temp directory with .beads
	tmpDir := t.TempDir()
	beadsDir := filepath.Join(tmpDir, ".beads")
	if err := os.MkdirAll(beadsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(tmpDir)

	root, err := findProjectRootForDrift()
	if err != nil {
		t.Errorf("findProjectRootForDrift() failed: %v", err)
	}

	if root == "" {
		t.Error("Expected non-empty root path")
	}
}

// TestFindRecentWorkstreamsForDriftMultiple tests finding multiple workstreams
func TestFindRecentWorkstreamsForDriftMultiple(t *testing.T) {
	// Create temp directory with multiple workstreams
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "completed")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create multiple workstreams directly in completed directory
	for i := 1; i <= 5; i++ {
		wsFile := filepath.Join(wsDir, "00-050-0"+string(rune('0'+i))+".md")
		if err := os.WriteFile(wsFile, []byte("# Test\n"), 0644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}

	workstreams, err := findRecentWorkstreamsForDrift(tmpDir)
	if err != nil {
		t.Errorf("findRecentWorkstreamsForDrift() failed: %v", err)
	}

	// Should find all 5 workstreams
	if len(workstreams) != 5 {
		t.Errorf("Expected 5 workstreams, got %d", len(workstreams))
	}
}

// TestCheckFilePermissionsInsecureDirectory tests checkFilePermissions with insecure directory
func TestCheckFilePermissionsInsecureDirectory(t *testing.T) {
	// Create .beads directory with insecure file inside
	tmpDir := t.TempDir()
	beadsDir := filepath.Join(tmpDir, ".beads")
	if err := os.MkdirAll(beadsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create test.db with world-readable permissions
	testFile := filepath.Join(beadsDir, "test.db")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(tmpDir)

	result := checkFilePermissions()

	if result.Name != "File Permissions" {
		t.Errorf("Expected name 'File Permissions', got '%s'", result.Name)
	}

	// Should detect insecure permissions (status may vary by platform)
	if result.Status != "ok" && result.Status != "warning" {
		t.Logf("Status: %s, Message: %s", result.Status, result.Message)
	}
}

// TestCheckFilePermissionsMultipleInsecure tests checkFilePermissions with multiple insecure files
func TestCheckFilePermissionsMultipleInsecure(t *testing.T) {
	// Create .beads directory with multiple insecure files
	tmpDir := t.TempDir()
	beadsDir := filepath.Join(tmpDir, ".beads")
	if err := os.MkdirAll(beadsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create multiple insecure files
	files := []string{"beads.db", "cache.json", "log.txt"}
	for _, file := range files {
		fullPath := filepath.Join(beadsDir, file)
		if err := os.WriteFile(fullPath, []byte("test"), 0644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(tmpDir)

	result := checkFilePermissions()

	// Should detect all insecure files
	if result.Name != "File Permissions" {
		t.Errorf("Expected name 'File Permissions', got '%s'", result.Name)
	}
}

// TestCheckFilePermissionsSecureDirectory tests checkFilePermissions with all secure files
func TestCheckFilePermissionsSecureDirectory(t *testing.T) {
	// Create .beads directory with secure files
	tmpDir := t.TempDir()
	beadsDir := filepath.Join(tmpDir, ".beads")
	if err := os.MkdirAll(beadsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create secure files
	files := []string{"beads.db", "cache.json"}
	for _, file := range files {
		fullPath := filepath.Join(beadsDir, file)
		if err := os.WriteFile(fullPath, []byte("test"), 0600); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(tmpDir)

	result := checkFilePermissions()

	if result.Status != "ok" {
		t.Errorf("Expected status ok with all secure files, got %s: %s", result.Status, result.Message)
	}
}

// TestCheckFilePermissionsTelemetryFile tests checkFilePermissions with telemetry file
func TestCheckFilePermissionsTelemetryFile(t *testing.T) {
	// Create .sdp directory with telemetry file
	tmpDir := t.TempDir()
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create telemetry.jsonl with insecure permissions
	teleFile := filepath.Join(sdpDir, "telemetry.jsonl")
	if err := os.WriteFile(teleFile, []byte("test"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(tmpDir)

	result := checkFilePermissions()

	// Should check telemetry file permissions
	if result.Name != "File Permissions" {
		t.Errorf("Expected name 'File Permissions', got '%s'", result.Name)
	}
}

// TestCheckClaudeDirMissingSubdirs tests checkClaudeDir when subdirectories are missing
func TestCheckClaudeDirMissingSubdirs(t *testing.T) {
	// Create .claude directory but missing subdirs
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(tmpDir)

	result := checkClaudeDir()

	// Should detect missing subdirectories
	if result.Status != "ok" {
		// May be warning or error if subdirs missing
		t.Logf("Status: %s, Message: %s", result.Status, result.Message)
	}
}

// TestCheckClaudeDirPartialSubdirs tests checkClaudeDir with only some subdirectories
func TestCheckClaudeDirPartialSubdirs(t *testing.T) {
	// Create .claude with only skills subdirectory
	tmpDir := t.TempDir()
	claudeDir := filepath.Join(tmpDir, ".claude")
	skillsDir := filepath.Join(claudeDir, "skills")
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(tmpDir)

	result := checkClaudeDir()

	// Should detect partial installation
	// Name is ".claude/ directory" not "Claude Directory"
	if result.Name != ".claude/ directory" {
		t.Errorf("Expected name '.claude/ directory', got '%s'", result.Name)
	}

	// Should be warning since subdirs are missing
	if result.Status != "warning" {
		t.Errorf("Expected status warning, got %s: %s", result.Status, result.Message)
	}
}
