package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestQualityWatcher_DebounceBehavior tests debounce timing behavior
func TestQualityWatcher_DebounceBehavior(t *testing.T) {
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

	changeCount := 0

	qw, err := NewQualityWatcher(tmpDir, &QualityWatcherConfig{
		Quiet: true,
	})
	if err != nil {
		t.Fatalf("Failed to create quality watcher: %v", err)
	}
	defer qw.Close()

	// Track changes
	originalOnChange := qw.watcher.config.OnChange
	qw.watcher.config.OnChange = func(path string) {
		changeCount++
		originalOnChange(path)
	}

	// Start watcher
	go qw.Start()
	defer qw.Stop()

	time.Sleep(100 * time.Millisecond)

	// Create file
	testFile := filepath.Join(tmpDir, "test.go")
	err = os.WriteFile(testFile, []byte("package test\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Wait for debounce
	time.Sleep(200 * time.Millisecond)

	if changeCount < 1 {
		t.Errorf("Expected at least 1 change, got %d", changeCount)
	}
}

// TestQualityWatcher_ViolationStates tests violation state transitions
func TestQualityWatcher_ViolationStates(t *testing.T) {
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
	defer qw.Close()

	// Test empty state
	violations := qw.GetViolations()
	if len(violations) != 0 {
		t.Errorf("Expected 0 violations initially, got %d", len(violations))
	}

	// Add violation
	testViolation := Violation{
		File:     "test.go",
		Check:    "test-check",
		Message:  "test message",
		Severity: "warning",
	}
	qw.addViolation(testViolation)

	violations = qw.GetViolations()
	if len(violations) != 1 {
		t.Errorf("Expected 1 violation, got %d", len(violations))
	}

	// Verify violation content
	if violations[0].File != "test.go" {
		t.Errorf("Expected file 'test.go', got '%s'", violations[0].File)
	}
	if violations[0].Check != "test-check" {
		t.Errorf("Expected check 'test-check', got '%s'", violations[0].Check)
	}
	if violations[0].Severity != "warning" {
		t.Errorf("Expected severity 'warning', got '%s'", violations[0].Severity)
	}

	// Clear the violation
	qw.clearViolations("test.go")

	violations = qw.GetViolations()
	if len(violations) != 0 {
		t.Errorf("Expected 0 violations after clear, got %d", len(violations))
	}

	// Add multiple violations for different files
	for i := 0; i < 5; i++ {
		qw.addViolation(Violation{
			File:     fmt.Sprintf("file%d.go", i),
			Check:    "test",
			Message:  "test",
			Severity: "error",
		})
	}

	violations = qw.GetViolations()
	if len(violations) != 5 {
		t.Errorf("Expected 5 violations, got %d", len(violations))
	}

	// Clear one file
	qw.clearViolations("file2.go")

	violations = qw.GetViolations()
	if len(violations) != 4 {
		t.Errorf("Expected 4 violations after clearing one, got %d", len(violations))
	}
}

// TestWatcher_WatchPath tests watch path handling
func TestWatcher_WatchPath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	watcher, err := NewWatcher(tmpDir, &WatcherConfig{})
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Verify watch path is set
	if watcher.watchPath != tmpDir {
		t.Errorf("Expected watch path %s, got %s", tmpDir, watcher.watchPath)
	}

	// Add watch should work
	err = watcher.addWatch(tmpDir)
	if err != nil {
		t.Errorf("Failed to add watch for tmpDir: %v", err)
	}

	// Try to add watch for non-existent path
	err = watcher.addWatch("/nonexistent/path")
	if err == nil {
		t.Error("Expected error for non-existent path")
	}
}
