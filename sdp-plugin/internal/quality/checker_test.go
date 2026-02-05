package quality

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewChecker(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-new-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create pyproject.toml to detect as Python
	tomlContent := "[tool.pytest]\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "pyproject.toml"), []byte(tomlContent), 0644); err != nil {
		t.Fatal(err)
	}

	checker, err := NewChecker(tmpDir)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}

	if checker == nil {
		t.Fatal("Expected non-nil checker")
	}

	if checker.projectPath != tmpDir {
		t.Errorf("Expected projectPath %s, got %s", tmpDir, checker.projectPath)
	}

	if checker.projectType != Python {
		t.Errorf("Expected projectType Python, got %d", checker.projectType)
	}
}

func TestDetectProjectTypePython(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-detect-py-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{projectPath: tmpDir}

	// Test with pyproject.toml
	tomlContent := "[tool.pytest]\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "pyproject.toml"), []byte(tomlContent), 0644); err != nil {
		t.Fatal(err)
	}

	pt, err := checker.detectProjectType()
	if err != nil {
		t.Fatalf("detectProjectType failed: %v", err)
	}

	if pt != Python {
		t.Errorf("Expected Python, got %d", pt)
	}
}

func TestDetectProjectTypeGo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-detect-go-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{projectPath: tmpDir}

	// Test with go.mod
	modContent := "module test\n\ngo 1.21\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(modContent), 0644); err != nil {
		t.Fatal(err)
	}

	pt, err := checker.detectProjectType()
	if err != nil {
		t.Fatalf("detectProjectType failed: %v", err)
	}

	if pt != Go {
		t.Errorf("Expected Go, got %d", pt)
	}
}

func TestCheckCoveragePython(t *testing.T) {
	// Create temp Python project
	tmpDir, err := os.MkdirTemp("", "sdp-test-python-*")
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

	result, err := checker.CheckCoverage()
	if err != nil {
		t.Fatalf("CheckCoverage failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.ProjectType != "Python" {
		t.Errorf("Expected ProjectType Python, got %s", result.ProjectType)
	}

	if result.Threshold != 80.0 {
		t.Errorf("Expected threshold 80.0, got %f", result.Threshold)
	}
}

func TestCheckCoverageGo(t *testing.T) {
	// Create temp Go project
	tmpDir, err := os.MkdirTemp("", "sdp-test-go-*")
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

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Go,
	}

	result, err := checker.CheckCoverage()
	if err != nil {
		t.Fatalf("CheckCoverage failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.ProjectType != "Go" {
		t.Errorf("Expected ProjectType Go, got %s", result.ProjectType)
	}

	if result.Threshold != 80.0 {
		t.Errorf("Expected threshold 80.0, got %f", result.Threshold)
	}
}

func TestCheckCoverageJava(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-java-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Java,
	}

	result, err := checker.CheckCoverage()
	if err != nil {
		t.Fatalf("CheckCoverage failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.ProjectType != "Java" {
		t.Errorf("Expected ProjectType Java, got %s", result.ProjectType)
	}
}

func TestCheckComplexity(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-complexity-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Python,
	}

	result, err := checker.CheckComplexity()
	if err != nil {
		t.Fatalf("CheckComplexity failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.Threshold != 10 {
		t.Errorf("Expected threshold 10, got %d", result.Threshold)
	}
}

func TestCheckComplexityGo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-complexity-go-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Go,
	}

	result, err := checker.CheckComplexity()
	if err != nil {
		t.Fatalf("CheckComplexity failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.Threshold != 10 {
		t.Errorf("Expected threshold 10, got %d", result.Threshold)
	}
}

func TestCheckFileSize(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-size-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Python,
	}

	result, err := checker.CheckFileSize()
	if err != nil {
		t.Fatalf("CheckFileSize failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.Threshold != 200 {
		t.Errorf("Expected threshold 200, got %d", result.Threshold)
	}

	if result.TotalFiles != 0 {
		t.Errorf("Expected 0 files, got %d", result.TotalFiles)
	}

	if !result.Passed {
		t.Errorf("Expected Passed to be true for empty directory")
	}
}

func TestCheckFileSizeWithViolations(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-size-viol-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a file with 300 lines
	lines := make([]string, 300)
	for i := range lines {
		lines[i] = "line " + string(rune('0'+i%10))
	}
	content := strings.Join(lines, "\n")
	if err := os.WriteFile(filepath.Join(tmpDir, "large.py"), []byte(content), 0644); err != nil {
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

	if result.Passed {
		t.Errorf("Expected Passed to be false with violations")
	}

	if len(result.Violators) != 1 {
		t.Errorf("Expected 1 violator, got %d", len(result.Violators))
	}
}

func TestCheckTypes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-types-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Python,
	}

	result, err := checker.CheckTypes()
	if err != nil {
		t.Fatalf("CheckTypes failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.ProjectType != "Python" {
		t.Errorf("Expected ProjectType Python, got %s", result.ProjectType)
	}
}

func TestCheckTypesGo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-types-go-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Go,
	}

	result, err := checker.CheckTypes()
	if err != nil {
		t.Fatalf("CheckTypes failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.ProjectType != "Go" {
		t.Errorf("Expected ProjectType Go, got %s", result.ProjectType)
	}
}

func TestCheckTypesJava(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-types-java-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{
		projectPath: tmpDir,
		projectType: Java,
	}

	result, err := checker.CheckTypes()
	if err != nil {
		t.Fatalf("CheckTypes failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.ProjectType != "Java" {
		t.Errorf("Expected ProjectType Java, got %s", result.ProjectType)
	}
}

func TestDetectProjectTypeByExtensions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "sdp-test-ext-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	checker := &Checker{projectPath: tmpDir}

	// Create .py files
	for i := 0; i < 3; i++ {
		if err := os.WriteFile(filepath.Join(tmpDir, "test"+string(rune('0'+i))+".py"), []byte("print('hello')"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	pt, err := checker.detectProjectType()
	if err != nil {
		t.Fatalf("detectProjectType failed: %v", err)
	}

	if pt != Python {
		t.Errorf("Expected Python by extension, got %d", pt)
	}
}
