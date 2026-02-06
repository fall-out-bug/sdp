package telemetry

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAskForConsentYesInput tests consent granted with 'y' input
func TestAskForConsentYesInput(t *testing.T) {
	// Mock stdin with 'y' input
	input := "y\n"
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	w.WriteString(input)
	w.Close()
	defer func() {
		os.Stdin = oldStdin
		r.Close()
	}()

	granted, err := AskForConsent()
	if err != nil {
		t.Fatalf("AskForConsent failed: %v", err)
	}

	if !granted {
		t.Error("Expected consent to be granted with 'y' input")
	}
}

// TestAskForConsentNoInput tests consent denied with 'n' input
func TestAskForConsentNoInput(t *testing.T) {
	// Mock stdin with 'n' input
	input := "n\n"
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	w.WriteString(input)
	w.Close()
	defer func() {
		os.Stdin = oldStdin
		r.Close()
	}()

	granted, err := AskForConsent()
	if err != nil {
		t.Fatalf("AskForConsent failed: %v", err)
	}

	if granted {
		t.Error("Expected consent to be denied with 'n' input")
	}
}

// TestAskForConsentNonInteractive tests behavior in non-interactive mode
func TestAskForConsentNonInteractive(t *testing.T) {
	// Create pipe that will return EOF immediately (non-interactive)
	r, w, _ := os.Pipe()
	w.Close()
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() {
		os.Stdin = oldStdin
		r.Close()
	}()

	// Should return false (disabled) without error in non-interactive mode
	granted, err := AskForConsent()
	if err != nil {
		t.Fatalf("AskForConsent should not error in non-interactive mode: %v", err)
	}

	if granted {
		t.Error("Expected consent to be denied in non-interactive mode")
	}
}

// TestConsentPersistenceAfterGrant tests that consent persists across sessions
func TestConsentPersistenceAfterGrant(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "telemetry.json")

	// First session: grant consent
	GrantConsent(configFile, true)

	// Simulate new session: read consent back
	granted, err := CheckConsent(configFile)
	if err != nil {
		t.Fatalf("CheckConsent failed: %v", err)
	}

	if !granted {
		t.Error("Consent should persist after granting")
	}

	// Verify it's not considered first run anymore
	if IsFirstRun(configFile) {
		t.Error("Should not be first run after consent is saved")
	}
}

// TestConsentPersistenceAfterDenial tests that denial persists
func TestConsentPersistenceAfterDenial(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "telemetry.json")

	// First session: deny consent
	GrantConsent(configFile, false)

	// Second session: check it's still denied
	granted, err := CheckConsent(configFile)
	if err != nil {
		t.Fatalf("CheckConsent failed: %v", err)
	}

	if granted {
		t.Error("Consent denial should persist")
	}
}

// TestConsentConfigOverwrite tests that consent can be changed
func TestConsentConfigOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "telemetry.json")

	// Grant consent
	GrantConsent(configFile, true)

	// Verify granted
	granted, _ := CheckConsent(configFile)
	if !granted {
		t.Error("Consent should be granted")
	}

	// Change mind: deny consent
	GrantConsent(configFile, false)

	// Verify denied
	granted, _ = CheckConsent(configFile)
	if granted {
		t.Error("Consent should be denied after overwrite")
	}

	// Change mind again: grant consent
	GrantConsent(configFile, true)

	// Verify granted again
	granted, _ = CheckConsent(configFile)
	if !granted {
		t.Error("Consent should be granted after second overwrite")
	}
}

// TestConsentWithInvalidConfigPath tests behavior with invalid path
func TestConsentWithInvalidConfigPath(t *testing.T) {
	// Use a path that cannot be created (e.g., /dev/null/telemetry.json)
	configFile := "/dev/null/invalid/telemetry.json"

	// GrantConsent should fail
	err := GrantConsent(configFile, true)
	if err == nil {
		t.Error("GrantConsent should fail with invalid path")
	}
}

// TestMultipleConsentCalls tests that multiple calls work correctly
func TestMultipleConsentCalls(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "telemetry.json")

	// Multiple calls should be idempotent
	GrantConsent(configFile, true)
	GrantConsent(configFile, true)
	GrantConsent(configFile, true)

	granted, _ := CheckConsent(configFile)
	if !granted {
		t.Error("Consent should still be granted after multiple calls")
	}
}

// TestConsentFormat tests JSON format of consent file
func TestConsentFormat(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "telemetry.json")

	GrantConsent(configFile, true)

	// Read file
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	// Verify it's valid JSON and contains expected fields
	dataStr := string(data)
	expectedFields := []string{
		`"enabled": true`,
		`{`,
		`}`,
	}

	for _, field := range expectedFields {
		if !strings.Contains(dataStr, field) {
			t.Errorf("Config file missing expected field: %s\nActual: %s", field, dataStr)
		}
	}
}
