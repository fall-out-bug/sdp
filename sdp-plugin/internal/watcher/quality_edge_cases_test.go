package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestQualityWatcher_QualityCheckIntegration tests full quality check integration
func TestQualityWatcher_QualityCheckIntegration(t *testing.T) {
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

	// Test with various file scenarios
	testCases := []struct {
		name    string
		content string
	}{
		{
			name:    "simple.go",
			content: "package test\n\nfunc Simple() {}\n",
		},
		{
			name:    "with_imports.go",
			content: "package test\n\nimport \"fmt\"\n\nfunc WithImports() {\n\tfmt.Println(\"test\")\n}\n",
		},
		{
			name:    "multi_func.go",
			content: "package test\n\nfunc Func1() {}\n\nfunc Func2() {}\n\nfunc Func3() {}\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, tc.name)
			err = os.WriteFile(testFile, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("Failed to write %s: %v", tc.name, err)
			}

			// Run checkFile (should not panic)
			qw.checkFile(testFile)
		})
	}
}

// TestQualityWatcher_ErrorPaths tests error handling paths
func TestQualityWatcher_ErrorPaths(t *testing.T) {
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

	t.Run("InvalidProjectPath", func(t *testing.T) {
		// Try to create watcher with non-existent path
		// Note: NewQualityWatcher may succeed even with invalid path
		// because file watcher is created lazily
		_, err := NewQualityWatcher("/nonexistent/path/xyz", &QualityWatcherConfig{
			Quiet: true,
		})
		// We just verify it doesn't panic
		_ = err
	})

	t.Run("NoGoMod", func(t *testing.T) {
		// Empty directory without go.mod
		emptyDir, err := os.MkdirTemp("", "empty-test")
		if err != nil {
			t.Fatalf("Failed to create empty temp dir: %v", err)
		}
		defer os.RemoveAll(emptyDir)

		// Should still create watcher (will default to Python)
		qw, err := NewQualityWatcher(emptyDir, &QualityWatcherConfig{
			Quiet: true,
		})
		if err != nil {
			t.Fatalf("Failed to create watcher for empty dir: %v", err)
		}
		defer qw.Close()

		if qw == nil {
			t.Error("Watcher is nil for empty directory")
		}
	})
}

// TestQualityWatcher_ViolationManagement tests violation tracking
func TestQualityWatcher_ViolationManagement(t *testing.T) {
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

	// Add multiple violations for same file
	for i := 0; i < 3; i++ {
		qw.addViolation(Violation{
			File:     "test.go",
			Check:    "test",
			Message:  "test violation",
			Severity: "error",
		})
	}

	violations := qw.GetViolations()
	if len(violations) != 3 {
		t.Errorf("Expected 3 violations, got %d", len(violations))
	}

	// Clear violations
	qw.clearViolations("test.go")

	violations = qw.GetViolations()
	if len(violations) != 0 {
		t.Errorf("Expected 0 violations after clear, got %d", len(violations))
	}

	// Clear non-existent file (should not panic)
	qw.clearViolations("nonexistent.go")
}

// TestWatcher_StartImmediately tests starting watcher immediately after creation
func TestWatcher_StartImmediately(t *testing.T) {
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

	// Start immediately without delay
	done := make(chan error)
	go func() {
		done <- watcher.Start()
	}()

	time.Sleep(50 * time.Millisecond)

	watcher.Stop()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Start returned error: %v", err)
		}
	case <-time.After(time.Second):
		t.Error("Watcher did not stop within timeout")
	}
}
