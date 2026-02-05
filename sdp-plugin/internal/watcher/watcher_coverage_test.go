package watcher

import (
	"os"
	"testing"
	"time"
)

// TestWatcher_ErrorHandling tests error handling in the watcher
func TestWatcher_ErrorHandling(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	errorCount := 0
	var lastError error

	watcher, err := NewWatcher(tmpDir, &WatcherConfig{
		OnError: func(err error) {
			errorCount++
			lastError = err
		},
	})
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Start watcher
	go watcher.Start()
	defer watcher.Stop()

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// No errors expected initially
	if errorCount != 0 {
		t.Errorf("Expected 0 errors initially, got %d", errorCount)
	}

	// Test OnError is set
	if watcher.config.OnError == nil {
		t.Error("OnError handler not set")
	}

	_ = lastError // Use the variable to avoid lint errors
}

// TestWatcher_DebounceReset tests that debounce timer resets correctly
func TestWatcher_DebounceReset(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	changeCount := 0

	watcher, err := NewWatcher(tmpDir, &WatcherConfig{
		IncludePatterns:  []string{"*.go"},
		DebounceInterval: 200 * time.Millisecond,
		OnChange: func(path string) {
			changeCount++
		},
	})
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Manually test debounce reset
	watcher.debounce("file1.go")
	time.Sleep(50 * time.Millisecond)

	// Reset timer with same file
	watcher.debounce("file1.go")
	time.Sleep(50 * time.Millisecond)

	// Reset again
	watcher.debounce("file1.go")

	// Wait for debounce to complete
	time.Sleep(300 * time.Millisecond)

	// Should only trigger once due to debouncing
	if changeCount != 1 {
		t.Errorf("Expected 1 change after debounce reset, got %d", changeCount)
	}
}

// TestWatcher_PatternMatchEdgeCases tests edge cases in pattern matching
func TestWatcher_PatternMatchEdgeCases(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	watcher, err := NewWatcher(tmpDir, &WatcherConfig{
		IncludePatterns: []string{"*.go"},
		ExcludePatterns: []string{"test*"},
	})
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	tests := []struct {
		file     string
		expected bool
	}{
		{"my.go", true},
		{"test.go", false},           // Excluded by test*
		{"test_file.go", false},      // Excluded by test*
		{"mytest.go", true},          // Not starting with test
		{"README", false},            // Not *.go
		{".git", false},              // Directory
		{"", false},                  // Empty string
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			result := watcher.matchesPatterns(tt.file)
			if result != tt.expected {
				t.Errorf("matchesPatterns(%q) = %v, want %v", tt.file, result, tt.expected)
			}
		})
	}
}
