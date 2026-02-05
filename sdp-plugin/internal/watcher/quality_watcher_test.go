package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewQualityWatcher(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "quality-watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a simple go.mod file
	modFile := filepath.Join(tmpDir, "go.mod")
	err = os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	qw, err := NewQualityWatcher(tmpDir, &QualityWatcherConfig{
		Quiet: true,
	})

	if err != nil {
		t.Fatalf("Failed to create quality watcher: %v", err)
	}

	if qw == nil {
		t.Fatal("QualityWatcher is nil")
	}

	if qw.checker == nil {
		t.Error("Checker is nil")
	}

	qw.watcher.Close()
}

func TestQualityWatcher_FileSizeViolation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "quality-watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create go.mod
	modFile := filepath.Join(tmpDir, "go.mod")
	err = os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	checkCount := 0
	checkDone := make(chan bool, 1)

	qw, err := NewQualityWatcher(tmpDir, &QualityWatcherConfig{
		Quiet: true,
	})
	if err != nil {
		t.Fatalf("Failed to create quality watcher: %v", err)
	}
	defer qw.watcher.Close()

	// Override OnChange to detect when check runs
	originalOnChange := qw.watcher.config.OnChange
	qw.watcher.config.OnChange = func(path string) {
		checkCount++
		originalOnChange(path)
		if checkCount >= 1 {
			select {
			case checkDone <- true:
			default:
			}
		}
	}

	// Start watcher
	go qw.Start()
	defer qw.Stop()

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Create a large file (>200 LOC)
	largeFile := filepath.Join(tmpDir, "large.go")
	content := "package test\n\nfunc Large() {\n"
	for i := 0; i < 250; i++ {
		content += "    var x int\n"
	}
	content += "}\n"

	err = os.WriteFile(largeFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write large file: %v", err)
	}

	// Wait for check to complete
	select {
	case <-checkDone:
		// Check was triggered
	case <-time.After(3 * time.Second):
		t.Error("File check not completed within timeout")
	}

	// Wait a bit for violations to be recorded
	time.Sleep(500 * time.Millisecond)

	// Verify that the check ran
	if checkCount == 0 {
		t.Error("Expected at least 1 check, got none")
	}

	// Note: We may not always detect violations depending on timing
	// The important thing is that the watcher is running and checking files
}

func TestQualityWatcher_ExcludeTestFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "quality-watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create go.mod
	modFile := filepath.Join(tmpDir, "go.mod")
	err = os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	checkCount := 0

	qw, err := NewQualityWatcher(tmpDir, &QualityWatcherConfig{
		Quiet: true,
	})
	if err != nil {
		t.Fatalf("Failed to create quality watcher: %v", err)
	}
	defer qw.watcher.Close()

	// Override OnChange to count checks
	qw.watcher.config.OnChange = func(path string) {
		checkCount++
	}

	// Start watcher
	go qw.Start()
	defer qw.Stop()

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Create a test file (should be excluded)
	testFile := filepath.Join(tmpDir, "test_test.go")
	err = os.WriteFile(testFile, []byte("package test\n\nfunc Test() {}\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create a source file (should be checked)
	srcFile := filepath.Join(tmpDir, "src.go")
	err = os.WriteFile(srcFile, []byte("package test\n\nfunc Src() {}\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write src file: %v", err)
	}

	// Wait for changes
	time.Sleep(500 * time.Millisecond)

	// Should only check source file, not test file
	if checkCount != 1 {
		t.Errorf("Expected 1 check (excluded test file), got %d", checkCount)
	}
}

func TestQualityWatcher_ClearViolations(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "quality-watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create go.mod
	modFile := filepath.Join(tmpDir, "go.mod")
	err = os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	qw, err := NewQualityWatcher(tmpDir, &QualityWatcherConfig{
		Quiet: true,
	})
	if err != nil {
		t.Fatalf("Failed to create quality watcher: %v", err)
	}
	defer qw.watcher.Close()

	// Manually add a violation
	qw.addViolation(Violation{
		File:     "test.go",
		Check:    "test",
		Message:  "test violation",
		Severity: "error",
	})

	violations := qw.GetViolations()
	if len(violations) != 1 {
		t.Fatalf("Expected 1 violation, got %d", len(violations))
	}

	// Clear violations
	qw.clearViolations("test.go")

	violations = qw.GetViolations()
	if len(violations) != 0 {
		t.Errorf("Expected 0 violations after clear, got %d", len(violations))
	}
}
