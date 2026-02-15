package reality

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProjectScanner_New(t *testing.T) {
	scanner := NewProjectScanner("/tmp")
	if scanner == nil {
		t.Fatal("scanner should not be nil")
	}
}

func TestProjectScanner_Scan(t *testing.T) {
	// Create temp project structure
	tempDir := t.TempDir()

	// Create files
	os.WriteFile(filepath.Join(tempDir, "main.go"), []byte("package main\n\nfunc main() {}"), 0644)
	os.WriteFile(filepath.Join(tempDir, "main_test.go"), []byte("package main\n\nfunc TestMain(t *testing.T) {}"), 0644)
	os.MkdirAll(filepath.Join(tempDir, "pkg"), 0755)
	os.WriteFile(filepath.Join(tempDir, "pkg", "handler.go"), []byte("package pkg\n\nfunc Handler() {}"), 0644)

	scanner := NewProjectScanner(tempDir)
	result, err := scanner.Scan()

	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	if result.Language == "" {
		t.Error("language should be detected")
	}
	if result.TotalFiles < 3 {
		t.Errorf("expected at least 3 files, got %d", result.TotalFiles)
	}
	if result.TestFiles < 1 {
		t.Errorf("expected at least 1 test file, got %d", result.TestFiles)
	}
	if result.LinesOfCode < 1 {
		t.Error("expected lines of code to be counted")
	}
}

func TestProjectScanner_SkipHidden(t *testing.T) {
	tempDir := t.TempDir()

	// Create visible file
	os.WriteFile(filepath.Join(tempDir, "visible.go"), []byte("package main"), 0644)

	// Create hidden directory with file
	os.MkdirAll(filepath.Join(tempDir, ".hidden"), 0755)
	os.WriteFile(filepath.Join(tempDir, ".hidden", "secret.go"), []byte("package hidden"), 0644)

	scanner := NewProjectScanner(tempDir)
	result, err := scanner.Scan()

	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Should only count visible.go, not .hidden/secret.go
	if result.TotalFiles > 2 {
		t.Errorf("hidden files should be skipped, got %d files", result.TotalFiles)
	}
}

func TestProjectScanner_SkipNodeModules(t *testing.T) {
	tempDir := t.TempDir()

	// Create file
	os.WriteFile(filepath.Join(tempDir, "index.js"), []byte("console.log('hi')"), 0644)

	// Create node_modules with file (should be skipped)
	os.MkdirAll(filepath.Join(tempDir, "node_modules", "pkg"), 0755)
	os.WriteFile(filepath.Join(tempDir, "node_modules", "pkg", "index.js"), []byte("module.exports = {}"), 0644)

	scanner := NewProjectScanner(tempDir)
	result, err := scanner.Scan()

	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// node_modules should be skipped
	for _, dir := range result.Directories {
		if filepath.Base(dir) == "node_modules" {
			t.Error("node_modules should be skipped in directories")
		}
	}
}

func TestCountLines(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.go")

	content := `package main

// This is a comment
func main() {
	println("hello")
}
`
	os.WriteFile(filePath, []byte(content), 0644)

	lines, err := countLines(filePath)
	if err != nil {
		t.Fatalf("countLines failed: %v", err)
	}

	// Should count: package main, func main() {, println("hello"), }
	// Should skip: empty lines and // comment
	if lines < 3 {
		t.Errorf("expected at least 3 non-empty, non-comment lines, got %d", lines)
	}
}

func TestCountLines_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "empty.go")
	os.WriteFile(filePath, []byte(""), 0644)

	lines, err := countLines(filePath)
	if err != nil {
		t.Fatalf("countLines failed: %v", err)
	}

	if lines != 0 {
		t.Errorf("empty file should have 0 lines, got %d", lines)
	}
}

func TestProjectScanner_NonexistentPath(t *testing.T) {
	scanner := NewProjectScanner("/nonexistent/path/12345")
	_, err := scanner.Scan()

	if err == nil {
		t.Error("should fail on nonexistent path")
	}
}

func TestScanResult_Fields(t *testing.T) {
	result := &ScanResult{
		Language:    "go",
		Framework:   "gin",
		TotalFiles:  10,
		LinesOfCode: 500,
		TestFiles:   3,
	}

	if result.Language != "go" {
		t.Error("language not set")
	}
	if result.Framework != "gin" {
		t.Error("framework not set")
	}
	if result.TestFiles != 3 {
		t.Error("test files not set")
	}
}
