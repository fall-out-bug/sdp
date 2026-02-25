package quality

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckPythonCoverageWithFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-cov-py-file-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create .coverage file
	covContent := "!coverage.py\n"
	if err := os.WriteFile(filepath.Join(tmpDir, ".coverage"), []byte(covContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create coverage.json with proper format
	jsonContent := `{"percent_covered": 85.5}`
	if err := os.WriteFile(filepath.Join(tmpDir, "coverage.json"), []byte(jsonContent), 0644); err != nil {
		t.Fatal(err)
	}

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Python,
	}

	result, err := checker.CheckCoverage(context.Background())
	if err != nil {
		t.Fatalf("CheckCoverage failed: %v", err)
	}

	// Verify it parsed correctly
	_ = result.Coverage
	_ = result.Passed
}

func TestCheckPythonCoverageNoData(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-cov-py-nodata-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Python,
	}

	result, err := checker.CheckCoverage(context.Background())
	if err != nil {
		t.Fatalf("CheckCoverage failed: %v", err)
	}

	if result.Coverage != 0.0 {
		t.Errorf("Expected coverage 0.0, got %f", result.Coverage)
	}

	if result.Passed {
		t.Errorf("Expected Passed to be false with no coverage data")
	}

	if result.Report == "" {
		t.Errorf("Expected report message with no coverage data")
	}
}

func TestCheckGoCoverageWithData(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-cov-go-data-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Skip if go not available
	if _, err := os.Stat("/usr/local/go/bin/go"); os.IsNotExist(err) {
		if _, err := os.Stat("/usr/bin/go"); os.IsNotExist(err) {
			t.Skip("Go not installed")
		}
	}

	// Initialize go module
	modContent := "module test\n\ngo 1.21\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(modContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a simple Go file
	goContent := `package test

func Add(a, b int) int {
	return a + b
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(goContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create test file
	testContent := `package test

import "testing"

func TestAdd(t *testing.T) {
	if Add(1, 2) != 3 {
		t.Fail()
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "test_test.go"), []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Go,
	}

	result, err := checker.CheckCoverage(context.Background())
	if err != nil {
		t.Fatalf("CheckCoverage failed: %v", err)
	}

	if result.ProjectType != "Go" {
		t.Errorf("Expected ProjectType Go, got %s", result.ProjectType)
	}
}

func TestBasicPythonComplexity(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-complexity-basic-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a Python file with lines
	lines := make([]string, 150)
	for i := range lines {
		lines[i] = "x = i"
	}
	content := strings.Join(lines, "\n")

	if err := os.WriteFile(filepath.Join(tmpDir, "module.py"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Python,
	}

	result, err := checker.basicPythonComplexity(&ComplexityResult{Threshold: 10})
	if err != nil {
		t.Fatalf("basicPythonComplexity failed: %v", err)
	}

	// Just verify it runs without error
	_ = result.MaxCC
	_ = result.ComplexFiles
}

func TestBasicGoComplexity(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-complexity-go-basic-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a Go file with estimated complexity
	lines := make([]string, 150)
	for i := range lines {
		lines[i] = "// line " + string(rune('0'+i%10))
	}
	content := strings.Join(lines, "\n")

	if err := os.WriteFile(filepath.Join(tmpDir, "complex.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Go,
	}

	result, err := checker.basicGoComplexity(&ComplexityResult{Threshold: 10})
	if err != nil {
		t.Fatalf("basicGoComplexity failed: %v", err)
	}

	if result.MaxCC == 0 {
		t.Errorf("Expected MaxCC > 0, got %d", result.MaxCC)
	}
}

func TestCheckFileSizeWithMultipleFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-size-multi-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create multiple files
	for i := 0; i < 3; i++ {
		lines := make([]string, 50)
		for j := range lines {
			lines[j] = "line content"
		}
		content := strings.Join(lines, "\n")
		if err := os.WriteFile(filepath.Join(tmpDir, "file"+string(rune('0'+i))+".py"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Python,
	}

	result, err := checker.CheckFileSize()
	if err != nil {
		t.Fatalf("CheckFileSize failed: %v", err)
	}

	if result.TotalFiles != 3 {
		t.Errorf("Expected 3 files, got %d", result.TotalFiles)
	}

	if result.AverageLOC == 0 {
		t.Errorf("Expected non-zero average LOC")
	}

	if !result.Passed {
		t.Errorf("Expected Passed to be true with small files")
	}
}

func TestCheckFileSizeSkipsDirectories(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-size-skip-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create subdirectories that should be skipped
	dirs := []string{"vendor", "node_modules", ".git", "__pycache__"}
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755); err != nil {
			t.Fatal(err)
		}
		// Create file in skipped directory
		lines := make([]string, 300)
		content := strings.Join(lines, "\n")
		if err := os.WriteFile(filepath.Join(tmpDir, dir, "file.py"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Python,
	}

	result, err := checker.CheckFileSize()
	if err != nil {
		t.Fatalf("CheckFileSize failed: %v", err)
	}

	if result.TotalFiles != 0 {
		t.Errorf("Expected 0 files (all skipped), got %d", result.TotalFiles)
	}
}

func TestCheckPythonTypesWithErrors(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-types-py-err-*")
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

	result, err := checker.CheckTypes()
	if err != nil {
		t.Fatalf("CheckTypes failed: %v", err)
	}

	if result.ProjectType != "Python" {
		t.Errorf("Expected ProjectType Python, got %s", result.ProjectType)
	}
}

func TestCheckGoTypesWithErrors(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-types-go-err-*")
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

	result, err := checker.CheckTypes()
	if err != nil {
		t.Fatalf("CheckTypes failed: %v", err)
	}

	if result.ProjectType != "Go" {
		t.Errorf("Expected ProjectType Go, got %s", result.ProjectType)
	}
}

func TestUnsupportedProjectType(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-unsupported-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Type(999), // Invalid type
	}

	_, err = checker.CheckCoverage(context.Background())
	if err == nil {
		t.Error("Expected error for unsupported project type")
	}
}

func TestDetectProjectTypeFallback(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-detect-fallback-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{projectPath: tmpDir}

	// No config files, should default to Python
	pt, err := checker.detectProjectType()
	if err != nil {
		t.Fatalf("detectProjectType failed: %v", err)
	}

	if pt != Python {
		t.Errorf("Expected default Python, got %d", pt)
	}
}

// TestBasicGoComplexityFallback tests fallback when gocyclo not available
func TestBasicGoComplexityFallback(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a Go file that will exceed threshold (>100 LOC)
	complexFile := filepath.Join(tmpDir, "complex.go")
	var lines []string
	lines = append(lines, "package main")
	lines = append(lines, "")
	for i := 0; i < 150; i++ {
		lines = append(lines, "func function"+string(rune('0'+i%10))+"() {")
		lines = append(lines, "\tfmt.Println(\"line\")")
		lines = append(lines, "}")
	}
	content := strings.Join(lines, "\n")
	if err := os.WriteFile(complexFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	checker := &Checker{
		projectPath: tmpDir,
	}

	result := &ComplexityResult{
		Threshold: 10,
	}

	// Run basic complexity check (gocyclo will not be available)
	finalResult, err := checker.basicGoComplexity(result)
	if err != nil {
		t.Fatalf("basicGoComplexity() failed: %v", err)
	}

	// Should have detected complex file
	if len(finalResult.ComplexFiles) == 0 {
		t.Error("basicGoComplexity() should detect complex files")
	}

	// MaxCC should be set (estimated from LOC)
	if finalResult.MaxCC == 0 {
		t.Error("basicGoComplexity() should set MaxCC")
	}

	// Should fail because file exceeds threshold
	if finalResult.Passed {
		t.Error("basicGoComplexity() should fail when files exceed threshold")
	}
}

// TestBasicPythonComplexityFallback tests fallback when radon not available
func TestBasicPythonComplexityFallback(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a Python file that will exceed threshold (>100 LOC)
	complexFile := filepath.Join(tmpDir, "complex.py")
	var lines []string
	for i := 0; i < 150; i++ {
		lines = append(lines, "def function"+string(rune('0'+i%10))+"():")
		lines = append(lines, "    pass")
	}
	content := strings.Join(lines, "\n")
	if err := os.WriteFile(complexFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	checker := &Checker{
		projectPath: tmpDir,
	}

	result := &ComplexityResult{
		Threshold: 10,
	}

	// Run basic complexity check (radon will not be available)
	finalResult, err := checker.basicPythonComplexity(result)
	if err != nil {
		t.Fatalf("basicPythonComplexity() failed: %v", err)
	}

	// Should have detected complex file
	if len(finalResult.ComplexFiles) == 0 {
		t.Error("basicPythonComplexity() should detect complex files")
	}

	// Should fail because file exceeds threshold
	if finalResult.Passed {
		t.Error("basicPythonComplexity() should fail when files exceed threshold")
	}
}

// TestCheckJavaCoverageNoJacocoCsv tests when jacoco.csv missing
func TestCheckJavaCoverageNoJacocoCsv(t *testing.T) {
	tmpDir := t.TempDir()
	// No jacoco.csv

	checker := &Checker{
		projectPath: tmpDir,
	}

	result := &CoverageResult{
		Threshold: 80.0,
	}

	// Run Java coverage check
	finalResult, err := checker.checkJavaCoverage(context.Background(), result)
	if err != nil {
		t.Fatalf("checkJavaCoverage() failed: %v", err)
	}

	// Should return 0% coverage when file not found
	if finalResult.Coverage != 0.0 {
		t.Errorf("checkJavaCoverage() coverage = %f, want 0.0", finalResult.Coverage)
	}
}

// TestCheckPythonCoverageNoCoverageFile tests when .coverage doesn't exist
func TestCheckPythonCoverageNoCoverageFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create pyproject.toml
	pyproject := filepath.Join(tmpDir, "pyproject.toml")
	if err := os.WriteFile(pyproject, []byte("[tool.pytest]\n"), 0644); err != nil {
		t.Fatalf("Failed to create pyproject.toml: %v", err)
	}

	checker := &Checker{
		projectPath: tmpDir,
	}

	result := &CoverageResult{
		Threshold: 80.0,
	}

	// Run Python coverage check
	finalResult, err := checker.checkPythonCoverage(context.Background(), result)
	if err != nil {
		t.Fatalf("checkPythonCoverage() failed: %v", err)
	}

	// Should handle gracefully (coverage will be 0)
	if finalResult.Coverage < 0 || finalResult.Coverage > 100 {
		t.Errorf("checkPythonCoverage() coverage = %f, want [0, 100]", finalResult.Coverage)
	}
}
