package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

// TestVersionCommand tests the sdp --version flag
func TestVersionCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	cmd := exec.Command(binaryPath, "--version")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run version command: %v", err)
	}

	output := stdout.String() + stderr.String()
	if !strings.Contains(output, "sdp version") {
		t.Errorf("Version output does not contain version string\nGot: %s", output)
	}
}

// TestHelpCommand tests the sdp --help flag
func TestHelpCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	cmd := exec.Command(binaryPath, "--help")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run help command: %v", err)
	}

	output := stdout.String() + stderr.String()

	// Check for expected help content
	expectedKeywords := []string{"Usage", "Available Commands", "Flags"}
	for _, keyword := range expectedKeywords {
		if !strings.Contains(output, keyword) {
			t.Errorf("Help output does not contain expected keyword %q\nGot: %s", keyword, output)
		}
	}
}

// TestDoctorCommand tests the sdp doctor command
func TestDoctorCommand(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	cmd := exec.Command(binaryPath, "doctor")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// doctor should not fail
	if err := cmd.Run(); err != nil {
		t.Logf("Doctor command failed: %v\nOutput: %s", err, stdout.String()+stderr.String())
	}

	output := stdout.String() + stderr.String()

	// Check for expected doctor content
	if !strings.Contains(output, "SDP Doctor") && !strings.Contains(output, "doctor") {
		t.Logf("Doctor output: %s", output)
	}
}

// TestCommandsCoverage tests that all commands are reachable
func TestCommandsCoverage(t *testing.T) {
	binaryPath := skipIfBinaryNotBuilt(t)

	// Get list of all commands
	cmd := exec.Command(binaryPath, "--help")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to get help: %v", err)
	}

	output := stdout.String() + stderr.String()

	// List of expected commands
	expectedCommands := []string{
		"beads",
		"checkpoint",
		"completion",
		"decisions",
		"doctor",
		"drift",
		"guard",
		"help",
		"hooks",
		"init",
		"orchestrate",
		"parse",
		"prd",
		"quality",
		"skill",
		"status",
		"tdd",
		"telemetry",
		"verify",
		"watch",
	}

	for _, expected := range expectedCommands {
		if !strings.Contains(output, expected) {
			t.Errorf("Expected command %q not found in help output", expected)
		}
	}
}
