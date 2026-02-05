package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestWatcher_MatchesPatterns(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	watcher, err := NewWatcher(tmpDir, &WatcherConfig{
		IncludePatterns: []string{"*.go", "*.py"},
		ExcludePatterns: []string{"*_test.go"},
	})
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	tests := []struct {
		file     string
		expected bool
	}{
		{"test.go", true},
		{"test.py", true},
		{"test_test.go", false},
		{"test.txt", false},
		{"mock_test.go", false},
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

func TestWatcher_AddWatch(t *testing.T) {
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

	// Test adding existing directory
	err = watcher.addWatch(tmpDir)
	if err != nil {
		t.Errorf("Failed to add watch for existing directory: %v", err)
	}

	// Test adding non-existent path
	err = watcher.addWatch("/nonexistent/path")
	if err == nil {
		t.Error("Expected error for non-existent path, got nil")
	}
}

func TestWatcher_HandleEvent_Create(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	createCount := 0

	watcher, err := NewWatcher(tmpDir, &WatcherConfig{
		IncludePatterns: []string{"*.go"},
		OnChange: func(path string) {
			createCount++
		},
	})
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Test create event
	event := fsnotify.Event{
		Name: filepath.Join(tmpDir, "test.go"),
		Op:   fsnotify.Create,
	}

	watcher.handleEvent(event)

	// Give debounce time to trigger
	time.Sleep(200 * time.Millisecond)

	if createCount != 1 {
		t.Errorf("Expected 1 create event, got %d", createCount)
	}
}

func TestWatcher_HandleEvent_Write(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	writeCount := 0

	watcher, err := NewWatcher(tmpDir, &WatcherConfig{
		IncludePatterns: []string{"*.go"},
		OnChange: func(path string) {
			writeCount++
		},
	})
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Test write event
	event := fsnotify.Event{
		Name: filepath.Join(tmpDir, "test.go"),
		Op:   fsnotify.Write,
	}

	watcher.handleEvent(event)

	// Give debounce time to trigger
	time.Sleep(200 * time.Millisecond)

	if writeCount != 1 {
		t.Errorf("Expected 1 write event, got %d", writeCount)
	}
}

func TestWatcher_HandleEvent_Remove(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	removeCount := 0

	watcher, err := NewWatcher(tmpDir, &WatcherConfig{
		IncludePatterns: []string{"*.go"},
		OnChange: func(path string) {
			removeCount++
		},
	})
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Test remove event (should not trigger)
	event := fsnotify.Event{
		Name: filepath.Join(tmpDir, "test.go"),
		Op:   fsnotify.Remove,
	}

	watcher.handleEvent(event)

	time.Sleep(100 * time.Millisecond)

	if removeCount != 0 {
		t.Errorf("Expected 0 events for remove, got %d", removeCount)
	}
}

func TestWatcher_HandleEvent_Directory(t *testing.T) {
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

	// Create a subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Test directory create event
	event := fsnotify.Event{
		Name: subDir,
		Op:   fsnotify.Create,
	}

	// Should not panic
	watcher.handleEvent(event)
}

func TestWatcher_NoIncludePatterns(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// No include patterns - should match all
	watcher, err := NewWatcher(tmpDir, &WatcherConfig{})
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	tests := []struct {
		file     string
		expected bool
	}{
		{"test.go", true},
		{"test.py", true},
		{"test.txt", true},
		{"anyfile", true},
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

func TestQualityWatcher_PrintViolation(t *testing.T) {
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

	// Test printing violations (quiet mode should not panic)
	violation := Violation{
		File:     "test.go",
		Check:    "complexity",
		Message:  "Too complex",
		Severity: "warning",
	}

	// Should not panic even in quiet mode
	qw.printViolation(violation)
}
