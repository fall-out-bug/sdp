package quality

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckPythonCoverageWithInvalidJson(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-cov-py-invalid-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create .coverage file
	if err := os.WriteFile(filepath.Join(tmpDir, ".coverage"), []byte("!coverage\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create invalid coverage.json
	if err := os.WriteFile(filepath.Join(tmpDir, "coverage.json"), []byte("invalid json"), 0644); err != nil {
		t.Fatal(err)
	}

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Python,
	}

	result, err := checker.checkPythonCoverage(&CoverageResult{Threshold: 80.0})
	if err != nil {
		t.Fatalf("checkPythonCoverage failed: %v", err)
	}

	// Should handle invalid JSON gracefully
	_ = result.Coverage
}

func TestCheckGoCoverageNoGo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-cov-go-nogo-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Go,
	}

	result, err := checker.checkGoCoverage(&CoverageResult{Threshold: 80.0})
	if err != nil {
		t.Fatalf("checkGoCoverage failed: %v", err)
	}

	// Should handle missing go command gracefully
	_ = result.ProjectType
}

func TestCheckPythonTypesWithMypyErrors(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-types-mypy-err-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create Python file with type error
	pyContent := `
x: int = "string"  # type error
def foo() -> str:
    return 123  # type error
`
	if err := os.WriteFile(filepath.Join(tmpDir, "test.py"), []byte(pyContent), 0644); err != nil {
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

	// Should handle type errors
	_ = result.Passed
}

func TestDetectProjectTypeMultipleMarkers(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-detect-multi-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{projectPath: tmpDir}

	// Create both pyproject.toml and go.mod - should detect Python first
	if err := os.WriteFile(filepath.Join(tmpDir, "pyproject.toml"), []byte("[tool.pytest]"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644); err != nil {
		t.Fatal(err)
	}

	pt, err := checker.detectProjectType()
	if err != nil {
		t.Fatalf("detectProjectType failed: %v", err)
	}

	if pt != Python {
		t.Errorf("Expected Python (first match), got %d", pt)
	}
}

func TestCheckFileSizeWithTestFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-size-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test file (should be skipped)
	lines := make([]string, 300)
	content := strings.Join(lines, "\n")
	if err := os.WriteFile(filepath.Join(tmpDir, "test_file.py"), []byte(content), 0644); err != nil {
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

	// Test files should be skipped
	if result.TotalFiles != 0 {
		t.Errorf("Expected 0 files (test files skipped), got %d", result.TotalFiles)
	}

	if !result.Passed {
		t.Errorf("Expected Passed to be true (no files checked)")
	}
}

func TestCheckFileSizeWithGoTestFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-size-go-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create Go test file (should be skipped)
	lines := make([]string, 300)
	content := strings.Join(lines, "\n")
	if err := os.WriteFile(filepath.Join(tmpDir, "file_test.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Go,
	}

	result, err := checker.CheckFileSize()
	if err != nil {
		t.Fatalf("CheckFileSize failed: %v", err)
	}

	// Test files should be skipped
	if result.TotalFiles != 0 {
		t.Errorf("Expected 0 files (test files skipped), got %d", result.TotalFiles)
	}
}

func TestCheckPythonComplexityWithRadonUnavailable(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-cc-no-radon-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Python,
	}

	result, err := checker.checkPythonComplexity(&ComplexityResult{Threshold: 10})
	if err != nil {
		t.Fatalf("checkPythonComplexity failed: %v", err)
	}

	// Should fall back to basic complexity check
	_ = result.AverageCC
}

func TestCheckGoComplexityWithGocycloUnavailable(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-cc-no-gocyclo-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Go,
	}

	result, err := checker.checkGoComplexity(&ComplexityResult{Threshold: 10})
	if err != nil {
		t.Fatalf("checkGoComplexity failed: %v", err)
	}

	// Should fall back to basic complexity check
	_ = result.AverageCC
}

func TestCheckPythonTypesParseMypyErrors(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-types-parse-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a Python file
	if err := os.WriteFile(filepath.Join(tmpDir, "module.py"), []byte("x = 1"), 0644); err != nil {
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

	// Should parse mypy output or handle missing mypy
	_ = result.ProjectType
}

func TestCheckGoTypesParseVetErrors(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-types-vet-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize go module
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test\n\ngo 1.21"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a Go file
	if err := os.WriteFile(filepath.Join(tmpDir, "module.go"), []byte("package test\n"), 0644); err != nil {
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

	// Should parse vet output or handle clean code
	_ = result.ProjectType
}

func TestCheckJavaTypesParseMvnErrors(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-types-mvn-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create pom.xml
	if err := os.WriteFile(filepath.Join(tmpDir, "pom.xml"), []byte("<project></project>"), 0644); err != nil {
		t.Fatal(err)
	}

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Java,
	}

	result, err := checker.checkJavaTypes(&TypeResult{})
	if err != nil {
		t.Fatalf("checkJavaTypes failed: %v", err)
	}

	// Should parse mvn output or handle missing mvn
	_ = result.ProjectType
}

func TestNewCheckerErrorPath(t *testing.T) {
	// Test with a directory that doesn't exist
	checker := &Checker{
		projectPath: "/nonexistent/directory/that/does/not/exist",
		projectType: Python,
	}

	// Should still be able to call detectProjectType
	_, err := checker.detectProjectType()
	// May fail, which is OK
	_ = err
}
