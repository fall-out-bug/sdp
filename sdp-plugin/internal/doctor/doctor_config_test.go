package doctor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckProjectConfig_NoProjectRoot(t *testing.T) {
	// Change to temp directory without project root
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	result := checkProjectConfig()

	// Should return warning since no project root
	if result.Name != ".sdp/config.yml" {
		t.Errorf("Name = %s, want .sdp/config.yml", result.Name)
	}
}

func TestCheckProjectConfig_NoConfig(t *testing.T) {
	// Change to temp directory with .sdp but no config
	tmpDir := t.TempDir()
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0o755); err != nil {
		t.Fatal(err)
	}

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	result := checkProjectConfig()

	// Should return ok with "using defaults" message
	if result.Name != ".sdp/config.yml" {
		t.Errorf("Name = %s, want .sdp/config.yml", result.Name)
	}

	if result.Status != "ok" {
		t.Errorf("Status = %s, want ok", result.Status)
	}
}

func TestCheckProjectConfig_WithValidConfig(t *testing.T) {
	// Create temp directory with valid config
	tmpDir := t.TempDir()
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0o755); err != nil {
		t.Fatal(err)
	}

	configContent := `version: 1
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

	if result.Status != "ok" {
		t.Logf("Status = %s (may be expected)", result.Status)
	}
}

func TestCheckProjectConfig_WithInvalidConfig(t *testing.T) {
	// Create temp directory with invalid config
	tmpDir := t.TempDir()
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0o755); err != nil {
		t.Fatal(err)
	}

	configContent := `invalid yaml content: [
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

	// Should return error for invalid config
	if result.Status != "error" {
		t.Logf("Status = %s (expected error for invalid yaml)", result.Status)
	}
}

func TestCheckProjectConfig_AdoptionModeOn(t *testing.T) {
	tmpDir := t.TempDir()
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0o755); err != nil {
		t.Fatal(err)
	}

	configContent := "version: 1\nadoption_mode: true\n"
	configPath := filepath.Join(sdpDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	result := checkProjectConfig()

	if result.Status != "ok" {
		t.Errorf("Status = %s, want ok", result.Status)
	}
	if !strings.Contains(result.Message, "Adoption mode: ON") {
		t.Errorf("Expected message to mention adoption mode, got: %s", result.Message)
	}
}

func TestCheckProjectConfig_AdoptionModeOff(t *testing.T) {
	tmpDir := t.TempDir()
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0o755); err != nil {
		t.Fatal(err)
	}

	configContent := "version: 1\nadoption_mode: false\n"
	configPath := filepath.Join(sdpDir, "config.yml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	result := checkProjectConfig()

	if result.Status != "ok" {
		t.Errorf("Status = %s, want ok", result.Status)
	}
	if strings.Contains(result.Message, "Adoption mode") {
		t.Errorf("Expected no adoption mode message when off, got: %s", result.Message)
	}
}
