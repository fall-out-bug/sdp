package doctor

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCheckProjectConfig_WithValidationErrors tests config validation failures
func TestCheckProjectConfig_WithValidationErrors(t *testing.T) {
	tmpDir := t.TempDir()
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create config with invalid version (should trigger validation error)
	configContent := `version: 0
`
	configPath := filepath.Join(sdpDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	result := checkProjectConfig()

	if result.Status != "error" {
		t.Errorf("Expected error status for invalid config, got %s: %s", result.Status, result.Message)
	}

	t.Logf("Status: %s, Message: %s", result.Status, result.Message)
}

// TestCheckProjectConfig_WithValidFullConfig tests a complete valid config
func TestCheckProjectConfig_WithValidFullConfig(t *testing.T) {
	tmpDir := t.TempDir()
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a complete valid config
	configContent := `version: 1
acceptance:
  command: "go test ./..."
  enabled: true
evidence:
  enabled: true
  log_path: ".sdp/log/events.jsonl"
`
	configPath := filepath.Join(sdpDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	result := checkProjectConfig()

	if result.Name != ".sdp/config.yml" {
		t.Errorf("Name = %s, want .sdp/config.yml", result.Name)
	}

	t.Logf("Status: %s, Message: %s", result.Status, result.Message)
}
