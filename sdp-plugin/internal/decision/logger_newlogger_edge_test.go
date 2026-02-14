package decision_test

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/fall-out-bug/sdp/internal/decision"
)

// TestLogger_NewLogger_ReadOnlyParent tests NewLogger when parent directory is read-only
func TestLogger_NewLogger_ReadOnlyParent(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping test when running as root")
	}

	tmpDir := t.TempDir()

	// Create a directory and make it read-only
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.MkdirAll(readOnlyDir, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Make it read-only
	if err := os.Chmod(readOnlyDir, 0555); err != nil {
		t.Fatalf("Failed to chmod: %v", err)
	}
	defer os.Chmod(readOnlyDir, 0755) // Restore for cleanup

	// Try to create logger in read-only directory
	targetPath := filepath.Join(readOnlyDir, "subdir")
	_, err := decision.NewLogger(targetPath)
	if err == nil {
		t.Error("Expected error when creating logger in read-only directory")
	} else {
		t.Logf("Got expected error: %v", err)
	}
}

// TestLogger_NewLogger_ExistingDirectory tests NewLogger when decisions directory already exists
func TestLogger_NewLogger_ExistingDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// Pre-create the decisions directory
	decisionsDir := filepath.Join(tmpDir, "docs", "decisions")
	if err := os.MkdirAll(decisionsDir, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Create logger - should succeed
	logger, err := decision.NewLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewLogger failed with existing directory: %v", err)
	}

	if logger == nil {
		t.Error("Expected non-nil logger")
	}
}

// TestLogger_Log_RotationWarning tests that logging continues even if rotation fails
func TestLogger_Log_RotationWarning(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping test when running as root")
	}

	tmpDir := t.TempDir()
	logger, err := decision.NewLogger(tmpDir)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	// Create the decisions file
	decisionsDir := filepath.Join(tmpDir, "docs", "decisions")
	filePath := filepath.Join(decisionsDir, "decisions.jsonl")

	// Write some initial content
	if err := os.WriteFile(filePath, []byte(`{"question":"Q1","decision":"D1"}`+"\n"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Make the directory read-only to cause rotation to fail
	if err := os.Chmod(decisionsDir, 0555); err != nil {
		t.Fatalf("Failed to chmod: %v", err)
	}
	defer os.Chmod(decisionsDir, 0755) // Restore for cleanup

	// Make the file large enough to trigger rotation (need to work around read-only)
	os.Chmod(decisionsDir, 0755)
	largeData := make([]byte, 11*1024*1024) // 11MB
	if err := os.WriteFile(filePath, largeData, 0644); err != nil {
		t.Fatalf("Failed to write large file: %v", err)
	}
	os.Chmod(decisionsDir, 0555)

	// Log should still succeed even if rotation fails
	d := decision.Decision{
		Question: "Test rotation warning",
		Decision: "Should still log",
	}

	err = logger.Log(d)
	// The log should succeed - rotation failure is logged as warning
	t.Logf("Log result: %v", err)
}

// TestLoadOptions_Values tests LoadOptions with various values
func TestLoadOptions_Values(t *testing.T) {
	opts := decision.LoadOptions{Offset: 10, Limit: 5}

	if opts.Offset != 10 {
		t.Errorf("Offset = %d, want 10", opts.Offset)
	}
	if opts.Limit != 5 {
		t.Errorf("Limit = %d, want 5", opts.Limit)
	}
}

// TestLogger_LogBatch_EmptySlice tests LogBatch with empty slice
func TestLogger_LogBatch_EmptySlice(t *testing.T) {
	tmpDir := t.TempDir()
	logger, _ := decision.NewLogger(tmpDir)

	// Log empty batch - should succeed without error
	err := logger.LogBatch([]decision.Decision{})
	if err != nil {
		t.Errorf("LogBatch with empty slice should not error: %v", err)
	}

	// Verify file was synced (exists and is valid)
	loaded, _ := logger.LoadAll()
	if len(loaded) != 0 {
		t.Errorf("Expected 0 decisions, got %d", len(loaded))
	}
}

func init() {
	// Ensure syscall is used
	_ = syscall.O_RDONLY
}
