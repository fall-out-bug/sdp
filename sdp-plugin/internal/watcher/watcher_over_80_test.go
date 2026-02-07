package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestWatcher_StopSequence tests the stop sequence
func TestWatcher_StopSequence(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	watcher, err := NewWatcher(tmpDir, &WatcherConfig{})
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}

	// Start watcher
	go watcher.Start()

	// Give it time to start
	time.Sleep(50 * time.Millisecond)

	// Stop watcher (this should trigger the done channel)
	watcher.Stop()

	// Verify stopped channel is closed
	select {
	case <-watcher.stopped:
		// Expected
	case <-time.After(time.Second):
		t.Error("stopped channel not closed within timeout")
	}

	// Multiple stops should be safe (idempotent)
	// Note: After Stop, we shouldn't call Close as it closes the same channels
}

// TestQualityWatcher_MultipleViolationsForSameFile tests handling multiple violations for one file
func TestQualityWatcher_MultipleViolationsForSameFile(t *testing.T) {
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

	// Add multiple violations for the same file
	violations := []Violation{
		{
			File:     "test.go",
			Check:    "size",
			Message:  "File too large",
			Severity: "error",
		},
		{
			File:     "test.go",
			Check:    "complexity",
			Message:  "Too complex",
			Severity: "warning",
		},
		{
			File:     "test.go",
			Check:    "types",
			Message:  "Type error",
			Severity: "error",
		},
	}

	for _, v := range violations {
		qw.addViolation(v)
	}

	// Should have all 3 violations
	allViolations := qw.GetViolations()
	if len(allViolations) != 3 {
		t.Errorf("Expected 3 violations, got %d", len(allViolations))
	}

	// Clear violations for the file
	qw.clearViolations("test.go")

	// Should have 0 violations
	allViolations = qw.GetViolations()
	if len(allViolations) != 0 {
		t.Errorf("Expected 0 violations after clear, got %d", len(allViolations))
	}

	// Add violations back
	for _, v := range violations {
		qw.addViolation(v)
	}

	// Should have 3 violations again
	allViolations = qw.GetViolations()
	if len(allViolations) != 3 {
		t.Errorf("Expected 3 violations after re-adding, got %d", len(allViolations))
	}

	// Clear non-existent file (should not affect count)
	qw.clearViolations("nonexistent.go")

	allViolations = qw.GetViolations()
	if len(allViolations) != 3 {
		t.Errorf("Expected 3 violations after clearing non-existent file, got %d", len(allViolations))
	}
}

// TestWatcher_FsnotifyWatcherCreation tests fsnotify watcher creation
func TestWatcher_FsnotifyWatcherCreation(t *testing.T) {
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

	// Verify fsnotify watcher is created
	if watcher.fsWatcher == nil {
		t.Error("fsnotify watcher should be created")
	}

	// Verify we can close the fsnotify watcher
	err = watcher.fsWatcher.Close()
	if err != nil {
		t.Errorf("Failed to close fsnotify watcher: %v", err)
	}
}

// TestQualityWatcher_ConfigOptions tests various config options
func TestQualityWatcher_ConfigOptions(t *testing.T) {
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

	// Test with custom patterns
	qw, err := NewQualityWatcher(tmpDir, &QualityWatcherConfig{
		IncludePatterns: []string{"*.py", "*.go"},
		ExcludePatterns: []string{"*_test.py", "*_test.go"},
		Quiet:           true,
	})
	if err != nil {
		t.Fatalf("Failed to create quality watcher: %v", err)
	}
	defer qw.Close()

	// Verify patterns are set
	if len(qw.watcher.config.IncludePatterns) != 2 {
		t.Errorf("Expected 2 include patterns, got %d", len(qw.watcher.config.IncludePatterns))
	}

	if len(qw.watcher.config.ExcludePatterns) != 2 {
		t.Errorf("Expected 2 exclude patterns, got %d", len(qw.watcher.config.ExcludePatterns))
	}

	// Verify quiet mode
	if !qw.quiet {
		t.Error("Quiet mode should be set")
	}
}
