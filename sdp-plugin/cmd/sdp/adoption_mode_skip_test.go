package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/config"
)

func TestAdoptionModeSkip_Off_ReturnsFalse(t *testing.T) {
	tmpDir := t.TempDir()
	cfgDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(cfgDir, "config.yml")
	if err := os.WriteFile(cfgPath, []byte("version: 1\nadoption_mode: false\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	origWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	if adoptionModeSkip("coverage", true) {
		t.Error("expected false when adoption mode is off")
	}
}

func TestAdoptionModeSkip_On_ReturnsTrue(t *testing.T) {
	tmpDir := t.TempDir()
	cfgDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(cfgDir, "config.yml")
	if err := os.WriteFile(cfgPath, []byte("version: 1\nadoption_mode: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	origWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	if !adoptionModeSkip("coverage", false) {
		t.Error("expected true when adoption mode is on")
	}
}

func TestAdoptionModeSkip_NoConfig_ReturnsFalse(t *testing.T) {
	tmpDir := t.TempDir()

	origWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	// No .sdp/config.yml → defaults → adoption_mode: false
	if adoptionModeSkip("coverage", true) {
		t.Error("expected false when no config exists (defaults to off)")
	}
}

func TestSetAdoptionMode_InvalidYAML_ReturnsError(t *testing.T) {
	tmpDir := t.TempDir()
	cfgDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(cfgDir, "config.yml")
	if err := os.WriteFile(cfgPath, []byte("version: 1\ninvalid: yaml: [\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := config.SetAdoptionMode(tmpDir, true)
	if err == nil {
		t.Error("expected error for invalid YAML in config")
	}
}
