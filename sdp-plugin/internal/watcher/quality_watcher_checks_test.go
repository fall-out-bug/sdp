package watcher

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/quality"
)

func TestQualityWatcher_OnFileChange(t *testing.T) {
	tmpDir := t.TempDir()

	// Create go.mod for valid module
	modFile := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a Go file
	goFile := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(goFile, []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}

	qw, err := NewQualityWatcher(tmpDir, &QualityWatcherConfig{Quiet: true})
	if err != nil {
		t.Fatalf("NewQualityWatcher failed: %v", err)
	}
	defer qw.watcher.Close()

	// Trigger file change
	qw.onFileChange(goFile)

	// Should have processed the file (violations depend on checker results)
	t.Logf("Violations: %d", len(qw.GetViolations()))
}

func TestQualityWatcher_CheckFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create go.mod for valid module
	modFile := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a Go file
	goFile := filepath.Join(tmpDir, "test.go")
	if err := os.WriteFile(goFile, []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}

	qw, err := NewQualityWatcher(tmpDir, &QualityWatcherConfig{Quiet: true})
	if err != nil {
		t.Fatalf("NewQualityWatcher failed: %v", err)
	}
	defer qw.watcher.Close()

	// Check file directly
	qw.checkFile(goFile)
}

func TestQualityWatcher_AddViolation(t *testing.T) {
	tmpDir := t.TempDir()

	// Create go.mod for valid module
	modFile := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	qw, err := NewQualityWatcher(tmpDir, &QualityWatcherConfig{Quiet: true})
	if err != nil {
		t.Fatalf("NewQualityWatcher failed: %v", err)
	}
	defer qw.watcher.Close()

	// Add violation
	v := Violation{
		File:     "test.go",
		Check:    "file-size",
		Message:  "File too large",
		Severity: "error",
	}
	qw.addViolation(v)

	violations := qw.GetViolations()
	if len(violations) != 1 {
		t.Errorf("Expected 1 violation, got %d", len(violations))
	}

	if violations[0].File != "test.go" {
		t.Errorf("File = %s, want test.go", violations[0].File)
	}
}

func TestViolation_Fields(t *testing.T) {
	v := Violation{
		File:     "path/to/file.go",
		Check:    "complexity",
		Message:  "CC too high",
		Severity: "warning",
	}

	if v.File != "path/to/file.go" {
		t.Errorf("File = %s", v.File)
	}
	if v.Check != "complexity" {
		t.Errorf("Check = %s", v.Check)
	}
	if v.Message != "CC too high" {
		t.Errorf("Message = %s", v.Message)
	}
	if v.Severity != "warning" {
		t.Errorf("Severity = %s", v.Severity)
	}
}

func TestQualityWatcherConfig_Defaults(t *testing.T) {
	tmpDir := t.TempDir()

	// Create go.mod for valid module
	modFile := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create with nil config (should use defaults)
	qw, err := NewQualityWatcher(tmpDir, nil)
	if err != nil {
		t.Fatalf("NewQualityWatcher failed: %v", err)
	}
	defer qw.watcher.Close()

	if qw.watchPath != tmpDir {
		t.Errorf("watchPath = %s, want %s", qw.watchPath, tmpDir)
	}
}

func TestQualityWatcher_checkFileSize(t *testing.T) {
	tmpDir := t.TempDir()

	// Create go.mod for valid module
	modFile := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create checker and watcher
	checker, err := quality.NewChecker(tmpDir)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}

	qw := &QualityWatcher{
		checker:   checker,
		watchPath: tmpDir,
		quiet:     true,
	}

	// Create a small file
	goFile := filepath.Join(tmpDir, "small.go")
	if err := os.WriteFile(goFile, []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Check file size (should not add violations for small file)
	qw.checkFileSize(goFile, "small.go")
}

func TestQualityWatcher_checkTypes(t *testing.T) {
	tmpDir := t.TempDir()

	// Create go.mod for valid module
	modFile := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	checker, err := quality.NewChecker(tmpDir)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}

	qw := &QualityWatcher{
		checker:   checker,
		watchPath: tmpDir,
		quiet:     true,
	}

	// Create a Go file
	goFile := filepath.Join(tmpDir, "typed.go")
	if err := os.WriteFile(goFile, []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Check types
	qw.checkTypes(goFile, "typed.go")
}
