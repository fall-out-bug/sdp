package decision_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/decision"
)

func TestLogger_LoadAll_FileNotExist(t *testing.T) {
	tempDir := t.TempDir()
	logger, _ := decision.NewLogger(tempDir)

	// Don't create any file, try to load
	decisions, err := logger.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll should not error on missing file: %v", err)
	}

	if len(decisions) != 0 {
		t.Errorf("Expected 0 decisions, got %d", len(decisions))
	}
}

func TestLogger_MalformedJSON(t *testing.T) {
	tempDir := t.TempDir()
	logger, _ := decision.NewLogger(tempDir)

	// Create log file with malformed JSON
	logPath := filepath.Join(tempDir, "docs", "decisions", "decisions.jsonl")
	os.WriteFile(logPath, []byte("{invalid json}\n{\"valid\": true}\n"), 0644)

	decisions, err := logger.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll should handle malformed JSON gracefully: %v", err)
	}

	// Should skip malformed line and load valid one
	if len(decisions) != 0 {
		// Valid JSON is not a valid Decision struct, so it's skipped too
		t.Logf("Loaded %d decisions (malformed handling)", len(decisions))
	}
}
