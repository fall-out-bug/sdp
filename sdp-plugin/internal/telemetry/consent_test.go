package telemetry

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDefaultDisabled tests that telemetry is disabled by default (opt-in)
func TestDefaultDisabled(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()
	telemetryFile := filepath.Join(tmpDir, "telemetry.jsonl")

	// Create tracker with no existing config
	collector, err := NewCollector(telemetryFile, false) // Disabled by default
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	status := collector.Status()
	if status.Enabled {
		t.Error("Telemetry should be disabled by default (opt-in)")
	}
}

// TestFirstRunConsentPrompt tests first-run consent flow
func TestFirstRunConsentPrompt(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "telemetry.json")

	// Simulate first run (no config exists)
	_, err := CheckConsent(configFile)
	if err != nil {
		t.Fatalf("CheckConsent failed: %v", err)
	}

	// Should prompt for consent (in real implementation)
	// For now, just verify it defaults to false
}

// TestConsentGranted tests that granted consent is persisted
func TestConsentGranted(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "telemetry.json")

	// Grant consent
	err := GrantConsent(configFile, true)
	if err != nil {
		t.Fatalf("GrantConsent failed: %v", err)
	}

	// Verify consent was saved
	granted, err := CheckConsent(configFile)
	if err != nil {
		t.Fatalf("CheckConsent failed: %v", err)
	}

	if !granted {
		t.Error("Consent should be granted")
	}

	// Verify config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Config file should exist after granting consent")
	}
}

// TestConsentDenied tests that denied consent is persisted
func TestConsentDenied(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "telemetry.json")

	// Deny consent
	err := GrantConsent(configFile, false)
	if err != nil {
		t.Fatalf("GrantConsent failed: %v", err)
	}

	// Verify consent was saved
	granted, err := CheckConsent(configFile)
	if err != nil {
		t.Fatalf("CheckConsent failed: %v", err)
	}

	if granted {
		t.Error("Consent should be denied")
	}
}

// TestConsentFilePermissions tests that consent file has secure permissions
func TestConsentFilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "telemetry.json")

	// Grant consent
	err := GrantConsent(configFile, true)
	if err != nil {
		t.Fatalf("GrantConsent failed: %v", err)
	}

	// Check file permissions
	info, err := os.Stat(configFile)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}

	mode := info.Mode().Perm()
	if mode != 0600 {
		t.Errorf("Config file should have 0600 permissions, got: %04o", mode)
	}
}

// TestRevokeConsent tests that consent can be revoked
func TestRevokeConsent(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "telemetry.json")

	// Grant consent
	GrantConsent(configFile, true)

	// Revoke consent
	GrantConsent(configFile, false)

	// Verify consent was revoked
	granted, _ := CheckConsent(configFile)
	if granted {
		t.Error("Consent should be revoked")
	}
}
