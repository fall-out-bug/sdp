package watcher

import (
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewWatcher(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	watcher, err := NewWatcher(tmpDir, &WatcherConfig{
		IncludePatterns: []string{"*.go"},
		ExcludePatterns: []string{"*_test.go"},
	})

	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}

	if watcher == nil {
		t.Fatal("Watcher is nil")
	}

	watcher.Close()
}

func TestWatcher_StartStop(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	watcher, err := NewWatcher(tmpDir, &WatcherConfig{
		IncludePatterns: []string{"*.go"},
	})
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer watcher.Close()

	// Start watcher
	done := make(chan error)
	go func() {
		done <- watcher.Start()
	}()

	// Give it time to start
	time.Sleep(100 * time.Millisecond)

	// Stop watcher
	watcher.Stop()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Watch returned error: %v", err)
		}
	case <-time.After(time.Second):
		t.Error("Watcher did not stop within timeout")
	}
}

func TestWatcher_FileChange(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	var changeCount int64
	changeChan := make(chan string, 10)

	watcher, err := NewWatcher(tmpDir, &WatcherConfig{
		IncludePatterns: []string{"*.go"},
		OnChange: func(path string) {
			atomic.AddInt64(&changeCount, 1)
			changeChan <- path
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

	// Create a file
	testFile := filepath.Join(tmpDir, "test.go")
	err = os.WriteFile(testFile, []byte("package test\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Wait for change event
	select {
	case path := <-changeChan:
		if path != testFile {
			t.Errorf("Expected path %s, got %s", testFile, path)
		}
	case <-time.After(2 * time.Second):
		t.Error("Did not receive file change event within timeout")
	}

	if atomic.LoadInt64(&changeCount) != 1 {
		t.Errorf("Expected 1 change, got %d", changeCount)
	}
}

func TestWatcher_ExcludePatterns(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	var changeCount int64

	watcher, err := NewWatcher(tmpDir, &WatcherConfig{
		IncludePatterns: []string{"*.go"},
		ExcludePatterns: []string{"*_test.go"},
		OnChange: func(path string) {
			atomic.AddInt64(&changeCount, 1)
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

	// Create excluded file
	testFile := filepath.Join(tmpDir, "test_test.go")
	err = os.WriteFile(testFile, []byte("package test\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create included file
	srcFile := filepath.Join(tmpDir, "src.go")
	err = os.WriteFile(srcFile, []byte("package test\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write src file: %v", err)
	}

	// Wait for change event
	time.Sleep(500 * time.Millisecond)

	if atomic.LoadInt64(&changeCount) != 1 {
		t.Errorf("Expected 1 change (excluded test file), got %d", atomic.LoadInt64(&changeCount))
	}
}

func TestWatcher_MultipleChanges(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	changeChan := make(chan string, 10)

	watcher, err := NewWatcher(tmpDir, &WatcherConfig{
		IncludePatterns: []string{"*.go"},
		OnChange: func(path string) {
			changeChan <- path
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

	// Create multiple files
	files := []string{"file1.go", "file2.go", "file3.go"}
	for _, fname := range files {
		fpath := filepath.Join(tmpDir, fname)
		err = os.WriteFile(fpath, []byte("package test\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write %s: %v", fname, err)
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Collect changes
	changes := make(map[string]bool)
	timeout := time.After(2 * time.Second)
	for i := 0; i < len(files); i++ {
		select {
		case path := <-changeChan:
			changes[path] = true
		case <-timeout:
			t.Errorf("Timeout waiting for changes, got %d/%d", len(changes), len(files))
			return
		}
	}

	if len(changes) != len(files) {
		t.Errorf("Expected %d changes, got %d", len(files), len(changes))
	}
}

func TestWatcher_Debounce(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	var changeCount int64
	var mu sync.Mutex

	watcher, err := NewWatcher(tmpDir, &WatcherConfig{
		IncludePatterns: []string{"*.go"},
		DebounceInterval: 100 * time.Millisecond,
		OnChange: func(path string) {
			mu.Lock()
			atomic.AddInt64(&changeCount, 1)
			mu.Unlock()
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

	// Write to same file multiple times rapidly
	testFile := filepath.Join(tmpDir, "test.go")
	for i := 0; i < 5; i++ {
		err = os.WriteFile(testFile, []byte("package test\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	// Wait for debounced changes
	time.Sleep(500 * time.Millisecond)

	// Should have fewer changes than writes due to debouncing
	mu.Lock()
	finalCount := changeCount
	mu.Unlock()

	if finalCount > 3 {
		t.Errorf("Debouncing not effective, got %d changes from 5 rapid writes", finalCount)
	}
}
