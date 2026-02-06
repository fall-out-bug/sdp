package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/fall-out-bug/sdp/internal/telemetry"
)

// TestTelemetryStatusCmd tests the telemetry status command
func TestTelemetryStatusCmd(t *testing.T) {
	// Create temp config directory
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "sdp")
	telemetryFile := filepath.Join(configDir, "telemetry.jsonl")

	// Create directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Create a collector and add some test events
	collector, err := telemetry.NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	// Track a test event
	collector.Record(telemetry.Event{
		Type:      "command_start",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"command": "test"},
	})

	// Capture output
	var out bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Set config dir env variable
	oldConfigDir := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", oldConfigDir)

	// Run command
	if err := telemetryStatusCmd.RunE(telemetryStatusCmd, []string{}); err != nil {
		t.Logf("Status command error: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout
	out.ReadFrom(r)

	output := out.String()

	// Check output contains expected fields
	if !strings.Contains(output, "Telemetry Status:") {
		t.Error("Status output should contain 'Telemetry Status:'")
	}
	if !strings.Contains(output, "Enabled:") {
		t.Error("Status output should contain 'Enabled:'")
	}
}

// TestTelemetryExportCmd tests the telemetry export command
func TestTelemetryExportCmd(t *testing.T) {
	// Create temp config directory
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "sdp")
	telemetryFile := filepath.Join(configDir, "telemetry.jsonl")

	// Create directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Create a collector and add test data
	collector, err := telemetry.NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	collector.Record(telemetry.Event{
		Type:      "test",
		Timestamp: time.Now(),
	})

	// Change to temp directory for export
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	// Set config dir env
	oldConfigDir := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", oldConfigDir)

	tests := []struct {
		name        string
		format      string
		expectError bool
	}{
		{
			name:        "export json",
			format:      "json",
			expectError: false,
		},
		{
			name:        "export csv",
			format:      "csv",
			expectError: false,
		},
		{
			name:        "invalid format",
			format:      "xml",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run export command
			err := telemetryExportCmd.RunE(telemetryExportCmd, []string{tt.format})

			if tt.expectError && err == nil {
				t.Errorf("Expected error for format %s but got none", tt.format)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for format %s: %v", tt.format, err)
			}

			// Check export file exists
			if !tt.expectError {
				exportPath := filepath.Join(tmpDir, "telemetry_export."+tt.format)
				if _, err := os.Stat(exportPath); os.IsNotExist(err) {
					t.Errorf("Export file %s was not created", exportPath)
				}
			}
		})
	}
}

// TestTelemetryEnableDisableCmd tests enable/disable commands
func TestTelemetryEnableDisableCmd(t *testing.T) {
	// Test enable - should not crash
	if err := telemetryEnableCmd.RunE(telemetryEnableCmd, []string{}); err != nil {
		t.Errorf("Enable command failed: %v", err)
	}

	// Test disable - should not crash
	if err := telemetryDisableCmd.RunE(telemetryDisableCmd, []string{}); err != nil {
		t.Errorf("Disable command failed: %v", err)
	}

	// Note: Actual file creation happens in real user config dir,
	// which varies by platform. We just verify commands don't crash.
}

// TestTelemetryAnalyzeCmd tests the analyze command
func TestTelemetryAnalyzeCmd(t *testing.T) {
	// Create temp config directory
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, "sdp")
	telemetryFile := filepath.Join(configDir, "telemetry.jsonl")

	// Create directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Create a collector and add test events
	collector, err := telemetry.NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	// Add test events
	events := []telemetry.Event{
		{
			Type:      "command_complete",
			Timestamp: time.Now(),
			Data:      map[string]interface{}{"command": "parse", "success": true},
		},
		{
			Type:      "command_complete",
			Timestamp: time.Now(),
			Data:      map[string]interface{}{"command": "parse", "success": false, "error": "test error"},
		},
	}

	for _, event := range events {
		collector.Record(event)
	}

	// Set config dir env
	oldConfigDir := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Setenv("XDG_CONFIG_HOME", oldConfigDir)

	// Capture output
	var out bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run analyze command
	if err := telemetryAnalyzeCmd.RunE(telemetryAnalyzeCmd, []string{}); err != nil {
		t.Logf("Analyze command error: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout
	out.ReadFrom(r)

	output := out.String()

	// Check output contains expected sections
	if !strings.Contains(output, "Telemetry Analysis Report") {
		t.Error("Analyze output should contain 'Telemetry Analysis Report'")
	}
	if !strings.Contains(output, "Total Events:") {
		t.Error("Analyze output should contain 'Total Events:'")
	}
}

// TestTelemetryUploadCmd tests the upload/packaging command
func TestTelemetryUploadCmd(t *testing.T) {
	// Use parent of TempDir as XDG_CONFIG_HOME to avoid duplicate paths
	tmpDir := t.TempDir()
	parentDir := filepath.Dir(tmpDir) // This avoids the duplicate /001/001 issue
	configDir := filepath.Join(parentDir, "sdp")
	telemetryFile := filepath.Join(configDir, "telemetry.jsonl")

	// Create directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Create test telemetry data
	collector, err := telemetry.NewCollector(telemetryFile, true)
	if err != nil {
		t.Fatalf("Failed to create collector: %v", err)
	}

	if err := collector.Record(telemetry.Event{
		Type:      telemetry.EventTypeCommandStart,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"command": "test",
		},
	}); err != nil {
		t.Fatalf("Failed to record event: %v", err)
	}

	// Set config dir env to parent directory
	oldConfigDir := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", parentDir)
	defer os.Setenv("XDG_CONFIG_HOME", oldConfigDir)

	// Change to tmpDir so upload files are created there
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	tests := []struct {
		name        string
		format      string
		expectError bool
	}{
		{
			name:        "upload json format",
			format:      "json",
			expectError: false,
		},
		{
			name:        "upload archive format",
			format:      "archive",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set format flag
			if err := telemetryUploadCmd.Flags().Set("format", tt.format); err != nil {
				t.Fatalf("Failed to set format flag: %v", err)
			}

			// Run upload command
			err := telemetryUploadCmd.RunE(telemetryUploadCmd, []string{})

			if tt.expectError && err == nil {
				t.Errorf("Expected error for format %s but got none", tt.format)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for format %s: %v", tt.format, err)
			}

			// Check upload file exists
			if !tt.expectError {
				files, _ := filepath.Glob(filepath.Join(tmpDir, "telemetry_upload_*"))
				if len(files) == 0 {
					t.Errorf("Upload file was not created for format %s", tt.format)
				}
			}
		})
	}
}

// TestTelemetryConsentCmd tests the consent info command
func TestTelemetryConsentCmd(t *testing.T) {
	// Capture output
	var out bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run consent command
	if err := telemetryConsentCmd.RunE(telemetryConsentCmd, []string{}); err != nil {
		t.Errorf("Consent command failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout
	out.ReadFrom(r)

	output := out.String()

	// Check output contains expected privacy info
	expectedStrings := []string{
		"Telemetry Consent:",
		"What's collected:",
		"What's NOT collected:",
		"No PII",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Consent output should contain %q", expected)
		}
	}
}
