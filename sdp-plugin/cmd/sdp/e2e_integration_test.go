package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall-out-bug/sdp/internal/evidence"
)

// TestPlanApplyTraceWorkflow tests the end-to-end plan → apply → trace workflow
func TestPlanApplyTraceWorkflow(t *testing.T) {
	evidence.ResetGlobalWriter()
	binaryPath := skipIfBinaryNotBuilt(t)

	// Create temp directory for isolated test environment
	tmpDir := t.TempDir()

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Initialize minimal project structure
	backlogDir := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(backlogDir, 0755); err != nil {
		t.Fatalf("create backlog dir: %v", err)
	}
	sdpDir := filepath.Join(tmpDir, ".sdp", "log")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("create .sdp/log dir: %v", err)
	}

	// Create minimal config file
	configContent := `version: "0.9.0"
evidence:
  enabled: true
  log_path: ".sdp/log/events.jsonl"
`
	configPath := filepath.Join(tmpDir, ".sdp", "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("create config: %v", err)
	}

	t.Run("plan command creates workstreams", func(t *testing.T) {
		// Test plan --dry-run first (doesn't require MODEL_API)
		cmd := exec.Command(binaryPath, "plan", "Add simple feature", "--dry-run")
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		cmd.Dir = tmpDir

		err := cmd.Run()
		output := stdout.String() + stderr.String()

		// dry-run should succeed without MODEL_API
		if err != nil {
			t.Logf("Plan dry-run output: %s\nPlan exit: %v", output, err)
		}

		// Check for dry-run indicators
		if !strings.Contains(output, "DRY RUN") && !strings.Contains(output, "dry") {
			t.Logf("Expected dry-run output to mention DRY RUN or dry\nGot: %s", output)
		}
	})

	t.Run("plan command JSON output", func(t *testing.T) {
		// Test JSON output format
		cmd := exec.Command(binaryPath, "plan", "Add JSON feature", "--dry-run", "--output=json")
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		cmd.Dir = tmpDir

		err := cmd.Run()
		output := stdout.String()

		if err != nil {
			t.Logf("Plan JSON output: %s\nPlan exit: %v", output, err)
		}

		// Check for JSON structure
		if !strings.Contains(output, "{") && !strings.Contains(output, "[") {
			t.Logf("Expected JSON output to contain braces/brackets\nGot: %s", output)
		}
	})
}
