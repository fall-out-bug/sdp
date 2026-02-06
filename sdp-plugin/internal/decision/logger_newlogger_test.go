package decision_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/decision"
)

func TestNewLogger_CreatesDirectory(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")

	logger, err := decision.NewLogger(subDir)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	if logger == nil {
		t.Error("Logger should not be nil")
	}

	// Verify directory was created
	decisionsDir := filepath.Join(subDir, "docs", "decisions")
	if _, err := os.Stat(decisionsDir); os.IsNotExist(err) {
		t.Error("decisions directory not created")
	}
}

func TestNewLogger_CanLogAfterCreation(t *testing.T) {
	tempDir := t.TempDir()

	logger, err := decision.NewLogger(tempDir)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}

	d := decision.Decision{
		Question: "Test",
		Decision: "Yes",
	}

	if err := logger.Log(d); err != nil {
		t.Errorf("Log failed after NewLogger: %v", err)
	}
}
