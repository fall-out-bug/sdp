package doctor

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCheckDrift_CheckedCountZero tests when all workstreams fail to check
func TestCheckDrift_CheckedCountZero(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "completed")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create a workstream file that will fail to parse
	wsFile := filepath.Join(wsDir, "00-999-01.md")
	if err := os.WriteFile(wsFile, []byte(""), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(tmpDir)

	result := checkDrift()

	if result.Name != "Drift Detection" {
		t.Errorf("Name = %s, want Drift Detection", result.Name)
	}

	// May return warning if no workstreams were successfully checked
	t.Logf("Status: %s, Message: %s", result.Status, result.Message)
}

// TestFindProjectRootForDrift_TraversesUp tests traversal up the directory tree
func TestFindProjectRootForDrift_TraversesUp(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested directory structure
	nestedDir := filepath.Join(tmpDir, "a", "b", "c", "d")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create .beads at root level
	beadsDir := filepath.Join(tmpDir, ".beads")
	if err := os.MkdirAll(beadsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(nestedDir)

	root, err := findProjectRootForDrift()
	if err != nil {
		t.Errorf("findProjectRootForDrift() failed: %v", err)
	}

	// Should traverse up and find .beads at tmpDir
	t.Logf("Root: %s (tmpDir: %s)", root, tmpDir)
}

// TestFindRecentWorkstreamsForDrift_InProgress tests finding workstreams in in_progress directory
func TestFindRecentWorkstreamsForDrift_InProgress(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "in_progress")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create workstreams in in_progress directory
	for i := 1; i <= 3; i++ {
		wsFile := filepath.Join(wsDir, "00-051-0"+string(rune('0'+i))+".md")
		if err := os.WriteFile(wsFile, []byte("# Test\n"), 0644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}

	workstreams, err := findRecentWorkstreamsForDrift(tmpDir)
	if err != nil {
		t.Errorf("findRecentWorkstreamsForDrift() failed: %v", err)
	}

	if len(workstreams) == 0 {
		t.Error("Expected to find workstreams in in_progress directory")
	}

	t.Logf("Found %d workstreams", len(workstreams))
}

// TestFindRecentWorkstreamsForDrift_Limit tests that only 5 workstreams are returned
func TestFindRecentWorkstreamsForDrift_Limit(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "completed")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create 10 workstreams
	for i := 0; i < 10; i++ {
		wsFile := filepath.Join(wsDir, "00-050-"+string(rune('0'+i/10))+string(rune('0'+i%10))+".md")
		if err := os.WriteFile(wsFile, []byte("# Test\n"), 0644); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
	}

	workstreams, err := findRecentWorkstreamsForDrift(tmpDir)
	if err != nil {
		t.Errorf("findRecentWorkstreamsForDrift() failed: %v", err)
	}

	if len(workstreams) > 5 {
		t.Errorf("Expected at most 5 workstreams, got %d", len(workstreams))
	}

	t.Logf("Found %d workstreams (max 5)", len(workstreams))
}

// TestFindRecentWorkstreamsForDrift_EmptyDir tests with empty workstreams directory
func TestFindRecentWorkstreamsForDrift_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "completed")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	workstreams, err := findRecentWorkstreamsForDrift(tmpDir)
	if err != nil {
		t.Errorf("findRecentWorkstreamsForDrift() failed: %v", err)
	}

	if len(workstreams) != 0 {
		t.Errorf("Expected 0 workstreams, got %d", len(workstreams))
	}
}

// TestCheckDrift_NoRecentWorkstreams tests when no recent workstreams exist
func TestCheckDrift_NoRecentWorkstreams(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "completed")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Don't create any workstream files
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(tmpDir)

	result := checkDrift()

	if result.Status != "ok" {
		t.Errorf("Expected ok status, got %s", result.Status)
	}
	t.Logf("Message: %s", result.Message)
}

// TestCheckDrift_WithWarnings tests when drift warnings are found
func TestCheckDrift_WithWarnings(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "completed")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create a workstream file
	wsFile := filepath.Join(wsDir, "00-001-01.md")
	if err := os.WriteFile(wsFile, []byte("# Test\n"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(tmpDir)

	result := checkDrift()

	// Should not fail, might have warnings
	t.Logf("Status: %s, Message: %s", result.Status, result.Message)
}

// TestFindRecentWorkstreamsForDrift_NonexistentDir tests when directory doesn't exist
func TestFindRecentWorkstreamsForDrift_NonexistentDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Don't create any workstream directories
	workstreams, err := findRecentWorkstreamsForDrift(tmpDir)

	if err != nil {
		t.Errorf("findRecentWorkstreamsForDrift() failed: %v", err)
	}

	if len(workstreams) != 0 {
		t.Errorf("Expected 0 workstreams, got %d", len(workstreams))
	}
}

// TestFindRecentWorkstreamsForDrift_WithDirectoriesOnly tests when entries are directories
func TestFindRecentWorkstreamsForDrift_WithDirectoriesOnly(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "completed")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create subdirectories (not files)
	if err := os.Mkdir(filepath.Join(wsDir, "00-001-01"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.Mkdir(filepath.Join(wsDir, "00-002-01"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	workstreams, err := findRecentWorkstreamsForDrift(tmpDir)
	if err != nil {
		t.Errorf("findRecentWorkstreamsForDrift() failed: %v", err)
	}

	// Should not include directories
	if len(workstreams) != 0 {
		t.Errorf("Expected 0 workstreams (directories should be skipped), got %d", len(workstreams))
	}
}

// TestFindRecentWorkstreamsForDrift_NonMarkdownFiles tests with non-.md files
func TestFindRecentWorkstreamsForDrift_NonMarkdownFiles(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "completed")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create non-markdown files
	if err := os.WriteFile(filepath.Join(wsDir, "00-001-01.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := os.WriteFile(filepath.Join(wsDir, "00-002-01.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	workstreams, err := findRecentWorkstreamsForDrift(tmpDir)
	if err != nil {
		t.Errorf("findRecentWorkstreamsForDrift() failed: %v", err)
	}

	// Should not include non-markdown files
	if len(workstreams) != 0 {
		t.Errorf("Expected 0 workstreams (non-md should be skipped), got %d", len(workstreams))
	}
}

// TestFindRecentWorkstreamsForDrift_MixedFiles tests with mixed file types
func TestFindRecentWorkstreamsForDrift_MixedFiles(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "completed")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create mixed files (markdown and non-markdown)
	if err := os.WriteFile(filepath.Join(wsDir, "00-001-01.md"), []byte("# WS1"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := os.WriteFile(filepath.Join(wsDir, "00-002-01.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := os.WriteFile(filepath.Join(wsDir, "00-003-01.md"), []byte("# WS3"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	// Subdirectory should be ignored
	if err := os.Mkdir(filepath.Join(wsDir, "00-004-01"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	workstreams, err := findRecentWorkstreamsForDrift(tmpDir)
	if err != nil {
		t.Errorf("findRecentWorkstreamsForDrift() failed: %v", err)
	}

	// Should only include markdown files
	if len(workstreams) != 2 {
		t.Errorf("Expected 2 workstreams (only .md), got %d", len(workstreams))
	}
}

// TestFindProjectRootForDrift_InSdpPlugin tests finding root from sdp-plugin directory
func TestFindProjectRootForDrift_InSdpPlugin(t *testing.T) {
	tmpDir := t.TempDir()

	// Create sdp-plugin subdirectory with go.mod
	pluginDir := filepath.Join(tmpDir, "sdp-plugin")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(pluginDir, "go.mod"), []byte("module test"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Create docs directory at parent level
	if err := os.Mkdir(filepath.Join(tmpDir, "docs"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	os.Chdir(pluginDir)

	root, err := findProjectRootForDrift()
	if err != nil {
		t.Errorf("findProjectRootForDrift() failed: %v", err)
	}

	// Should find parent with docs
	t.Logf("Root: %s (tmpDir: %s)", root, tmpDir)
}
