package drift

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAnalyzeFileGo(t *testing.T) {
	// Create test Go file
	tmpFile, err := os.CreateTemp("", "*.go")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `package test

// This is a test file
// It contains test functions
// For testing purposes
func TestFunc() {
}
`
	if err := os.WriteFile(tmpFile.Name(), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	analyzer := NewAnalyzer(filepath.Dir(tmpFile.Name()))
	purpose, err := analyzer.AnalyzeFile(filepath.Base(tmpFile.Name()))
	if err != nil {
		t.Fatalf("Failed to analyze file: %v", err)
	}

	if !strings.Contains(purpose, "test file") {
		t.Errorf("Expected purpose to mention 'test file', got: %s", purpose)
	}
}

func TestAnalyzeFilePython(t *testing.T) {
	// Create test Python file
	tmpFile, err := os.CreateTemp("", "*.py")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `"""
This module provides utility functions
for data processing and analysis.
"""

def process_data(data):
    pass
`
	if err := os.WriteFile(tmpFile.Name(), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	analyzer := NewAnalyzer(filepath.Dir(tmpFile.Name()))
	purpose, err := analyzer.AnalyzeFile(filepath.Base(tmpFile.Name()))
	if err != nil {
		t.Fatalf("Failed to analyze file: %v", err)
	}

	if !strings.Contains(purpose, "utility functions") {
		t.Errorf("Expected purpose to mention 'utility functions', got: %s", purpose)
	}
}

func TestComparePurposeMatch(t *testing.T) {
	// Create test file
	tmpFile, err := os.CreateTemp("", "*.go")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `package test

// Parser implementation for workstreams
func Parse() {
}
`
	if err := os.WriteFile(tmpFile.Name(), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	analyzer := NewAnalyzer(filepath.Dir(tmpFile.Name()))
	matches, actual := analyzer.ComparePurpose(filepath.Base(tmpFile.Name()), "parser")

	if !matches {
		t.Errorf("Expected purpose to match, got: %s", actual)
	}
}

func TestComparePurposeNoMatch(t *testing.T) {
	// Create test file
	tmpFile, err := os.CreateTemp("", "*.go")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `package test

// This is a utility file
func Helper() {
}
`
	if err := os.WriteFile(tmpFile.Name(), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	analyzer := NewAnalyzer(filepath.Dir(tmpFile.Name()))
	matches, _ := analyzer.ComparePurpose(filepath.Base(tmpFile.Name()), "database")

	if matches {
		t.Error("Expected purpose not to match 'database'")
	}
}

func TestExtractKeyTerms(t *testing.T) {
	purpose := "The parser processes workstream files"
	terms := extractKeyTerms(purpose)

	if len(terms) == 0 {
		t.Error("Expected to extract key terms")
	}

	// Should contain meaningful terms
	hasMeaningfulTerm := false
	for _, term := range terms {
		if term == "parser" || term == "processes" || term == "workstream" {
			hasMeaningfulTerm = true
		}
	}

	if !hasMeaningfulTerm {
		t.Errorf("Expected meaningful terms, got: %v", terms)
	}
}

func TestExtractKeyTermsStopWords(t *testing.T) {
	purpose := "the and or but for of with"
	terms := extractKeyTerms(purpose)

	// Should filter out stop words
	if len(terms) > 0 {
		t.Errorf("Expected no terms (all stop words), got: %v", terms)
	}
}
