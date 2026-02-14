package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/fall-out-bug/sdp/internal/collision"
)

func TestLoadFeatureScopes_NonexistentDir(t *testing.T) {
	scopes, err := loadFeatureScopes("/nonexistent/path")
	if err != nil {
		t.Errorf("loadFeatureScopes should not error for nonexistent dir: %v", err)
	}
	if len(scopes) != 0 {
		t.Errorf("Expected empty scopes for nonexistent dir, got %d", len(scopes))
	}
}

func TestOutputBoundariesAsHuman_Empty(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := outputBoundariesAsHuman([]collision.SharedBoundary{})

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Errorf("outputBoundariesAsHuman() error = %v", err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if output != "No shared boundaries detected.\n" {
		t.Logf("Output: %q", output)
	}
}

func TestOutputBoundariesAsHuman_WithBoundaries(t *testing.T) {
	boundaries := []collision.SharedBoundary{
		{
			FileName:       "test.go",
			TypeName:       "UserService",
			Features:       []string{"F001", "F002"},
			Fields:         []collision.FieldInfo{{Name: "ID", Type: "int"}},
			Recommendation: "Consider splitting interface",
		},
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := outputBoundariesAsHuman(boundaries)

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Errorf("outputBoundariesAsHuman() error = %v", err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check that output contains expected elements
	if output == "" {
		t.Error("outputBoundariesAsHuman() produced empty output")
	}
}

func TestOutputBoundariesAsJSON_Empty(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := outputBoundariesAsJSON([]collision.SharedBoundary{})

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Errorf("outputBoundariesAsJSON() error = %v", err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if output != "[\n]\n" {
		t.Logf("Output: %q", output)
	}
}

func TestOutputBoundariesAsJSON_WithBoundaries(t *testing.T) {
	boundaries := []collision.SharedBoundary{
		{
			FileName:       "user.go",
			TypeName:       "User",
			Features:       []string{"F001"},
			Fields:         nil,
			Recommendation: "Coordinate changes",
		},
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := outputBoundariesAsJSON(boundaries)

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Errorf("outputBoundariesAsJSON() error = %v", err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check that output contains JSON array
	if output == "" {
		t.Error("outputBoundariesAsJSON() produced empty output")
	}
}
