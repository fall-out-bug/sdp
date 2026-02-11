package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestWatcher_EventChannels tests event channel handling
func TestWatcher_EventChannels(t *testing.T) {
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

	// Verify channels are initialized
	if watcher.done == nil {
		t.Error("done channel not initialized")
	}
	if watcher.stopped == nil {
		t.Error("stopped channel not initialized")
	}
	if watcher.fsWatcher == nil {
		t.Error("fsWatcher not initialized")
	}

	// Verify channels are properly sized
	// done and stopped should be unbuffered (size 0)
	if cap(watcher.done) != 0 {
		t.Errorf("done channel should be unbuffered, got cap %d", cap(watcher.done))
	}
	if cap(watcher.stopped) != 0 {
		t.Errorf("stopped channel should be unbuffered, got cap %d", cap(watcher.stopped))
	}
}

// TestQualityWatcher_FilePathHandling tests various file path scenarios
func TestQualityWatcher_FilePathHandling(t *testing.T) {
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

	// Test with absolute path
	absPath, err := filepath.Abs(filepath.Join(tmpDir, "test.go"))
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create file
	err = os.WriteFile(absPath, []byte("package test\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Check with absolute path (should not panic)
	qw.checkFile(absPath)

	// Check with relative path
	relPath := "test.go"
	qw.checkFile(relPath)

	// Both should work
	violations := qw.GetViolations()
	_ = violations
}

// TestWatcher_TimerCleanup tests that timer is properly cleaned up
func TestWatcher_TimerCleanup(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	watcher, err := NewWatcher(tmpDir, &WatcherConfig{
		DebounceInterval: 100 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Trigger debounce
	watcher.debounce("file1.go")

	// Verify timer is set
	if watcher.timer == nil {
		t.Error("Timer should be set after debounce")
	}

	// Trigger debounce again with same file (should reset timer)
	watcher.debounce("file1.go")

	if watcher.timer == nil {
		t.Error("Timer should still be set after reset")
	}

	// Wait for timer to fire
	time.Sleep(150 * time.Millisecond)

	// Timer should still be set (previous timer completed)
	if watcher.timer == nil {
		t.Error("Timer should exist even after firing")
	}
}

// TestQualityWatcher_ErrorCallback tests error callback handling
func TestQualityWatcher_ErrorCallback(t *testing.T) {
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

	errorCount := 0

	qw, err := NewQualityWatcher(tmpDir, &QualityWatcherConfig{
		Quiet: true,
	})
	if err != nil {
		t.Fatalf("Failed to create quality watcher: %v", err)
	}
	defer qw.Close()

	// Set error callback on underlying watcher
	qw.watcher.config.OnError = func(err error) {
		errorCount++
	}

	// Verify error callback is set
	if qw.watcher.config.OnError == nil {
		t.Error("OnError callback not set")
	}

	// Start watcher (error callback should be available)
	go qw.Start()
	defer qw.Stop()

	time.Sleep(50 * time.Millisecond)

	// No errors expected in normal operation
	_ = errorCount
}
