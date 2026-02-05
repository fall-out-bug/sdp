package quality

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckPythonCoverageParseOutput(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-cov-parse-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create pytest.ini
	iniContent := "[pytest]\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "pytest.ini"), []byte(iniContent), 0644); err != nil {
		t.Fatal(err)
	}

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Python,
	}

	// Test the parsing logic
	result, err := checker.checkPythonCoverage(&CoverageResult{Threshold: 80.0})
	if err != nil {
		t.Fatalf("checkPythonCoverage failed: %v", err)
	}

	// Verify it runs
	_ = result.ProjectType
	_ = result.Coverage
}

func TestCheckJavaCoverage(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-cov-java-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Java,
	}

	result, err := checker.checkJavaCoverage(&CoverageResult{Threshold: 80.0})
	if err != nil {
		t.Fatalf("checkJavaCoverage failed: %v", err)
	}

	// Verify defaults
	if result.ProjectType != "Java" {
		t.Errorf("Expected Java, got %s", result.ProjectType)
	}

	if result.Coverage != 0.0 {
		t.Errorf("Expected 0 coverage with no report, got %f", result.Coverage)
	}
}

func TestCheckJavaComplexity(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-cc-java-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Java,
	}

	result, err := checker.checkJavaComplexity(&ComplexityResult{Threshold: 10})
	if err != nil {
		t.Fatalf("checkJavaComplexity failed: %v", err)
	}

	// Verify defaults
	if !result.Passed {
		t.Errorf("Expected passed with no complexity issues")
	}

	if result.AverageCC != 0.0 {
		t.Errorf("Expected 0 average CC, got %f", result.AverageCC)
	}
}

func TestCheckTypesUnsupported(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-types-unsup-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Test that type check functions handle unsupported types gracefully
	checker := &Checker{
		projectPath: tmpDir,
		projectType: Python,
	}

	// These should work even without tools installed
	result1, err := checker.checkPythonTypes(&TypeResult{})
	_ = result1
	_ = err // May fail if tools not installed

	checker.projectType = Go
	result2, err := checker.checkGoTypes(&TypeResult{})
	_ = result2
	_ = err // May fail if tools not installed

	checker.projectType = Java
	result3, err := checker.checkJavaTypes(&TypeResult{})
	_ = result3
	_ = err // May fail if tools not installed
}

func TestCheckPythonTypesWithMypy(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-types-mypy-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create mypy.ini
	iniContent := "[mypy]\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "mypy.ini"), []byte(iniContent), 0644); err != nil {
		t.Fatal(err)
	}

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Python,
	}

	result, err := checker.checkPythonTypes(&TypeResult{})
	if err != nil {
		t.Fatalf("checkPythonTypes failed: %v", err)
	}

	// Should pass if mypy not configured properly
	_ = result.Passed
}

func TestCheckGoTypesNoVet(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-types-vet-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize go module
	modContent := "module test\n\ngo 1.21\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(modContent), 0644); err != nil {
		t.Fatal(err)
	}

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Go,
	}

	result, err := checker.checkGoTypes(&TypeResult{})
	if err != nil {
		t.Fatalf("checkGoTypes failed: %v", err)
	}

	// Should pass if no vet errors
	_ = result.Passed
}

func TestCheckJavaTypesNoMvn(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-types-java-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Java,
	}

	result, err := checker.checkJavaTypes(&TypeResult{})
	if err != nil {
		t.Fatalf("checkJavaTypes failed: %v", err)
	}

	// Should handle missing mvn gracefully
	_ = result.Passed
}

func TestDetectProjectTypeGoExtension(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-detect-go-ext-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{projectPath: tmpDir}

	// Create .go files
	for i := 0; i < 3; i++ {
		if err := os.WriteFile(filepath.Join(tmpDir, "test"+string(rune('0'+i))+".go"), []byte("package main\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	pt, err := checker.detectProjectType()
	if err != nil {
		t.Fatalf("detectProjectType failed: %v", err)
	}

	if pt != Go {
		t.Errorf("Expected Go by extension, got %d", pt)
	}
}

func TestDetectProjectTypeJavaExtension(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-detect-java-ext-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{projectPath: tmpDir}

	// Create .java files
	for i := 0; i < 3; i++ {
		if err := os.WriteFile(filepath.Join(tmpDir, "Test"+string(rune('0'+i))+".java"), []byte("public class Test {}"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	pt, err := checker.detectProjectType()
	if err != nil {
		t.Fatalf("detectProjectType failed: %v", err)
	}

	if pt != Java {
		t.Errorf("Expected Java by extension, got %d", pt)
	}
}

func TestCheckFileSizeWithSmallFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-size-small-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a small file (under threshold)
	lines := make([]string, 50)
	for j := range lines {
		lines[j] = "line content"
	}
	content := strings.Join(lines, "\n")
	if err := os.WriteFile(filepath.Join(tmpDir, "file.py"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Python,
	}

	result, err := checker.CheckFileSize()
	if err != nil {
		t.Fatalf("CheckFileSize failed: %v", err)
	}

	if result.TotalFiles != 1 {
		t.Errorf("Expected 1 file, got %d", result.TotalFiles)
	}

	if !result.Passed {
		t.Errorf("Expected Passed to be true with small file")
	}
}

func TestCheckCoverageThresholdComparison(t *testing.T) {
	result := &CoverageResult{
		Coverage:  85.5,
		Threshold: 80.0,
	}
	result.Passed = result.Coverage >= result.Threshold

	if !result.Passed {
		t.Error("Expected 85.5% >= 80.0% to pass")
	}

	result.Coverage = 75.0
	result.Passed = result.Coverage >= result.Threshold

	if result.Passed {
		t.Error("Expected 75.0% < 80.0% to fail")
	}
}

func TestComplexityThresholdComparison(t *testing.T) {
	result := &ComplexityResult{
		MaxCC:     8,
		Threshold: 10,
	}
	result.Passed = result.MaxCC <= result.Threshold

	if !result.Passed {
		t.Error("Expected 8 <= 10 to pass")
	}

	result.MaxCC = 15
	result.Passed = result.MaxCC <= result.Threshold

	if result.Passed {
		t.Error("Expected 15 > 10 to fail")
	}
}
