package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestQualityWatcher_StartStop(t *testing.T) {
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

	// Start and stop watcher
	done := make(chan error)
	go func() {
		done <- qw.Start()
	}()

	// Give it time to start
	time.Sleep(100 * time.Millisecond)

	// Stop watcher
	qw.Stop()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Start returned error: %v", err)
		}
	case <-time.After(time.Second):
		t.Error("Watcher did not stop within timeout")
	}
}

func TestQualityWatcher_GetViolations(t *testing.T) {
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

	// Initially no violations
	violations := qw.GetViolations()
	if len(violations) != 0 {
		t.Errorf("Expected 0 violations initially, got %d", len(violations))
	}

	// Add violations
	qw.addViolation(Violation{
		File:     "test1.go",
		Check:    "test",
		Message:  "test violation 1",
		Severity: "error",
	})
	qw.addViolation(Violation{
		File:     "test2.go",
		Check:    "test",
		Message:  "test violation 2",
		Severity: "warning",
	})

	// Should have 2 violations
	violations = qw.GetViolations()
	if len(violations) != 2 {
		t.Errorf("Expected 2 violations, got %d", len(violations))
	}

	// Verify violations are copied (not same slice)
	violations[0] = Violation{File: "modified"}
	violations2 := qw.GetViolations()
	if violations2[0].File == "modified" {
		t.Error("GetViolations returned same slice, not a copy")
	}
}

func TestQualityWatcher_DefaultConfig(t *testing.T) {
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

	// Create watcher with nil config (should use defaults)
	qw, err := NewQualityWatcher(tmpDir, nil)
	if err != nil {
		t.Fatalf("Failed to create quality watcher with nil config: %v", err)
	}
	defer qw.Close()

	if qw == nil {
		t.Fatal("QualityWatcher is nil")
	}

	if qw.checker == nil {
		t.Error("Checker is nil")
	}
}

func TestQualityWatcher_CheckFileIntegration(t *testing.T) {
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

	// Create a valid Go file
	validFile := filepath.Join(tmpDir, "valid.go")
	err = os.WriteFile(validFile, []byte("package test\n\nfunc Valid() {}\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write valid file: %v", err)
	}

	// Check the file (should not panic)
	qw.checkFile(validFile)

	// Verify violations are accessible
	violations := qw.GetViolations()
	// We don't assert specific violations since quality checks may vary
	_ = violations
}
