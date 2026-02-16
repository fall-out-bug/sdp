package main

import (
	"os"
	"path/filepath"
	"testing"
)

// This file contains shared utility functions for integration tests.
// Integration tests require the sdp binary to be built separately.
// To build: go build -o sdp ./cmd/sdp

// skipIfBinaryNotBuilt skips the test if the sdp binary is not found.
// It returns the absolute path to the binary if it exists.
func skipIfBinaryNotBuilt(t *testing.T) string {
	// Tests run in sdp-plugin/cmd/sdp directory
	// Binary is in sdp-plugin/ directory (two levels up)
	// Need to go up TWO levels because of root-level cmd/ directory
	relativePath := filepath.Join("..", "..", "sdp")
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		t.Skip("Cannot resolve binary path")
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		t.Skip("sdp binary not found. Run: go build -o sdp ./cmd/sdp from sdp-plugin/ directory")
	}
	return absPath
}

// repoRoot returns the absolute path to the repository root
func repoRoot(t *testing.T) string {
	// From sdp-plugin/cmd/sdp, go up THREE levels to get to repo root
	// sdp-plugin/cmd/sdp → ../.. → sdp-plugin/ → ../../.. → sdp/
	path := filepath.Join("..", "..", "..")
	absPath, err := filepath.Abs(path)
	if err != nil {
		t.Fatalf("Cannot resolve repo root: %v", err)
	}
	return absPath
}
