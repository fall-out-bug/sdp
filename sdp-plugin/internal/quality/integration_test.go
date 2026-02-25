package quality

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// Test that we can create checkers for all project types
func TestNewCheckerAllTypes(t *testing.T) {
	// Python
	tmpPy, err := os.MkdirTemp("", "sdp-checker-py-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpPy)
	if err := os.WriteFile(filepath.Join(tmpPy, "pyproject.toml"), []byte("[tool.pytest]"), 0644); err != nil {
		t.Fatal(err)
	}
	checkerPy, err := NewChecker(tmpPy)
	if err != nil {
		t.Errorf("NewChecker Python failed: %v", err)
	}
	if checkerPy.projectType != Python {
		t.Errorf("Expected Python type, got %d", checkerPy.projectType)
	}

	// Go
	tmpGo, err := os.MkdirTemp("", "sdp-checker-go-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpGo)
	if err := os.WriteFile(filepath.Join(tmpGo, "go.mod"), []byte("module test\n\ngo 1.21"), 0644); err != nil {
		t.Fatal(err)
	}
	checkerGo, err := NewChecker(tmpGo)
	if err != nil {
		t.Errorf("NewChecker Go failed: %v", err)
	}
	if checkerGo.projectType != Go {
		t.Errorf("Expected Go type, got %d", checkerGo.projectType)
	}

	// Java
	tmpJava, err := os.MkdirTemp("", "sdp-checker-java-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpJava)
	if err := os.WriteFile(filepath.Join(tmpJava, "pom.xml"), []byte("<project></project>"), 0644); err != nil {
		t.Fatal(err)
	}
	checkerJava, err := NewChecker(tmpJava)
	if err != nil {
		t.Errorf("NewChecker Java failed: %v", err)
	}
	if checkerJava.projectType != Java {
		t.Errorf("Expected Java type, got %d", checkerJava.projectType)
	}
}

// Test all checker methods don't crash
func TestAllCheckersNoCrash(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-all-check-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	if err := os.WriteFile(filepath.Join(tmpDir, "pyproject.toml"), []byte("[tool.pytest]"), 0644); err != nil {
		t.Fatal(err)
	}

	checker, err := NewChecker(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	// Test all checkers - they should handle empty projects gracefully
	if _, err := checker.CheckCoverage(context.Background()); err != nil {
		t.Logf("CheckCoverage (empty project): %v", err)
	}
	if _, err := checker.CheckComplexity(); err != nil {
		t.Logf("CheckComplexity (empty project): %v", err)
	}
	if _, err := checker.CheckFileSize(); err != nil {
		t.Logf("CheckFileSize (empty project): %v", err)
	}
	if _, err := checker.CheckTypes(); err != nil {
		t.Logf("CheckTypes (empty project): %v", err)
	}
}

// Test coverage result structure
func TestCoverageResultStructure(t *testing.T) {
	result := &CoverageResult{
		ProjectType: "Python",
		Coverage:    85.5,
		Threshold:   80.0,
		Passed:      true,
		FilesBelow: []FileCoverage{
			{File: "test.py", Coverage: 75.0},
		},
	}

	if result.ProjectType != "Python" {
		t.Errorf("Expected Python, got %s", result.ProjectType)
	}

	if len(result.FilesBelow) != 1 {
		t.Errorf("Expected 1 file below threshold, got %d", len(result.FilesBelow))
	}

	if result.FilesBelow[0].File != "test.py" {
		t.Errorf("Expected test.py, got %s", result.FilesBelow[0].File)
	}
}

// Test complexity result structure
func TestComplexityResultStructure(t *testing.T) {
	result := &ComplexityResult{
		AverageCC: 5.5,
		MaxCC:     12,
		Threshold: 10,
		Passed:    false,
		ComplexFiles: []FileComplexity{
			{
				File:             "complex.py",
				AverageCC:        12.0,
				MaxCC:            12,
				ExceedsThreshold: true,
			},
		},
	}

	if result.AverageCC != 5.5 {
		t.Errorf("Expected 5.5, got %f", result.AverageCC)
	}

	if len(result.ComplexFiles) != 1 {
		t.Errorf("Expected 1 complex file, got %d", len(result.ComplexFiles))
	}

	if !result.ComplexFiles[0].ExceedsThreshold {
		t.Error("Expected ExceedsThreshold to be true")
	}
}

// Test file size result structure
func TestFileSizeResultStructure(t *testing.T) {
	result := &FileSizeResult{
		TotalFiles: 100,
		Violators: []FileViolation{
			{File: "large.py", LOC: 250},
		},
		Threshold:  200,
		Passed:     false,
		AverageLOC: 75,
	}

	if result.TotalFiles != 100 {
		t.Errorf("Expected 100 files, got %d", result.TotalFiles)
	}

	if len(result.Violators) != 1 {
		t.Errorf("Expected 1 violator, got %d", len(result.Violators))
	}

	if result.Violators[0].LOC != 250 {
		t.Errorf("Expected 250 LOC, got %d", result.Violators[0].LOC)
	}
}

// Test type result structure
func TestTypeResultStructure(t *testing.T) {
	result := &TypeResult{
		ProjectType: "Python",
		Passed:      false,
		Errors: []TypeError{
			{File: "test.py", Line: 10, Message: "type error"},
		},
		Warnings: []TypeError{
			{File: "", Line: 0, Message: "warning message"},
		},
	}

	if result.ProjectType != "Python" {
		t.Errorf("Expected Python, got %s", result.ProjectType)
	}

	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(result.Errors))
	}

	if result.Errors[0].Line != 10 {
		t.Errorf("Expected line 10, got %d", result.Errors[0].Line)
	}

	if len(result.Warnings) != 1 {
		t.Errorf("Expected 1 warning, got %d", len(result.Warnings))
	}
}

// Test type constants
func TestTypeConstants(t *testing.T) {
	if Python != 0 {
		t.Errorf("Expected Python to be 0, got %d", Python)
	}

	if Go != 1 {
		t.Errorf("Expected Go to be 1, got %d", Go)
	}

	if Java != 2 {
		t.Errorf("Expected Java to be 2, got %d", Java)
	}
}
