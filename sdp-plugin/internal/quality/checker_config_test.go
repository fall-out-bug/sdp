package quality

import (
	"os"
	"path/filepath"
	"testing"
)

func TestChecker_SetStrictMode(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a go.mod file for detection
	goModPath := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	checker, err := NewChecker(tmpDir)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}

	// Default should be false
	if checker.IsStrictMode() {
		t.Error("Default strict mode should be false")
	}

	// Enable strict mode
	checker.SetStrictMode(true)
	if !checker.IsStrictMode() {
		t.Error("Strict mode should be true after SetStrictMode(true)")
	}

	// Disable strict mode
	checker.SetStrictMode(false)
	if checker.IsStrictMode() {
		t.Error("Strict mode should be false after SetStrictMode(false)")
	}
}

func TestChecker_IsStrictMode(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a go.mod file for detection
	goModPath := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	checker, err := NewChecker(tmpDir)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}

	// Test default
	if checker.IsStrictMode() != false {
		t.Errorf("Expected false, got %v", checker.IsStrictMode())
	}

	// Test after setting
	checker.SetStrictMode(true)
	if checker.IsStrictMode() != true {
		t.Errorf("Expected true, got %v", checker.IsStrictMode())
	}
}

func TestNewChecker_ProjectPath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a go.mod file for detection
	goModPath := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	checker, err := NewChecker(tmpDir)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}

	if checker.projectPath != tmpDir {
		t.Errorf("Expected projectPath %s, got %s", tmpDir, checker.projectPath)
	}
}

func TestNewChecker_NoProjectFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Empty directory - should default to Python
	checker, err := NewChecker(tmpDir)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}

	if checker.projectType != Python {
		t.Errorf("Expected default type Python, got %v", checker.projectType)
	}
}

func TestDetectProjectType_Python_pyproject(t *testing.T) {
	tmpDir := t.TempDir()

	// Create pyproject.toml
	pyprojectPath := filepath.Join(tmpDir, "pyproject.toml")
	err := os.WriteFile(pyprojectPath, []byte("[tool.poetry]\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create pyproject.toml: %v", err)
	}

	checker := &Checker{projectPath: tmpDir}
	pt, err := checker.detectProjectType()
	if err != nil {
		t.Fatalf("detectProjectType failed: %v", err)
	}

	if pt != Python {
		t.Errorf("Expected Python, got %v", pt)
	}
}

func TestDetectProjectType_Python_setup(t *testing.T) {
	tmpDir := t.TempDir()

	// Create setup.py
	setupPath := filepath.Join(tmpDir, "setup.py")
	err := os.WriteFile(setupPath, []byte("from setuptools import setup\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create setup.py: %v", err)
	}

	checker := &Checker{projectPath: tmpDir}
	pt, err := checker.detectProjectType()
	if err != nil {
		t.Fatalf("detectProjectType failed: %v", err)
	}

	if pt != Python {
		t.Errorf("Expected Python, got %v", pt)
	}
}

func TestDetectProjectType_Python_requirements(t *testing.T) {
	tmpDir := t.TempDir()

	// Create requirements.txt
	reqPath := filepath.Join(tmpDir, "requirements.txt")
	err := os.WriteFile(reqPath, []byte("requests==2.28.0\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create requirements.txt: %v", err)
	}

	checker := &Checker{projectPath: tmpDir}
	pt, err := checker.detectProjectType()
	if err != nil {
		t.Fatalf("detectProjectType failed: %v", err)
	}

	if pt != Python {
		t.Errorf("Expected Python, got %v", pt)
	}
}

func TestDetectProjectType_Go(t *testing.T) {
	tmpDir := t.TempDir()

	// Create go.mod
	goModPath := filepath.Join(tmpDir, "go.mod")
	err := os.WriteFile(goModPath, []byte("module test\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	checker := &Checker{projectPath: tmpDir}
	pt, err := checker.detectProjectType()
	if err != nil {
		t.Fatalf("detectProjectType failed: %v", err)
	}

	if pt != Go {
		t.Errorf("Expected Go, got %v", pt)
	}
}

func TestDetectProjectType_Java(t *testing.T) {
	tmpDir := t.TempDir()

	// Create pom.xml
	pomPath := filepath.Join(tmpDir, "pom.xml")
	err := os.WriteFile(pomPath, []byte("<project></project>\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create pom.xml: %v", err)
	}

	checker := &Checker{projectPath: tmpDir}
	pt, err := checker.detectProjectType()
	if err != nil {
		t.Fatalf("detectProjectType failed: %v", err)
	}

	if pt != Java {
		t.Errorf("Expected Java, got %v", pt)
	}
}

func TestDetectProjectType_ByExtension_Python(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Python files
	for i := 0; i < 3; i++ {
		filename := filepath.Join(tmpDir, "test"+string(rune('0'+i))+".py")
		err := os.WriteFile(filename, []byte("print('hello')\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to create %s: %v", filename, err)
		}
	}

	checker := &Checker{projectPath: tmpDir}
	pt, err := checker.detectProjectType()
	if err != nil {
		t.Fatalf("detectProjectType failed: %v", err)
	}

	if pt != Python {
		t.Errorf("Expected Python, got %v", pt)
	}
}

func TestDetectProjectType_ByExtension_Go(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Go files
	for i := 0; i < 3; i++ {
		filename := filepath.Join(tmpDir, "test"+string(rune('0'+i))+".go")
		err := os.WriteFile(filename, []byte("package main\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to create %s: %v", filename, err)
		}
	}

	checker := &Checker{projectPath: tmpDir}
	pt, err := checker.detectProjectType()
	if err != nil {
		t.Fatalf("detectProjectType failed: %v", err)
	}

	if pt != Go {
		t.Errorf("Expected Go, got %v", pt)
	}
}

func TestDetectProjectType_ByExtension_Java(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Java files
	for i := 0; i < 3; i++ {
		filename := filepath.Join(tmpDir, "Test"+string(rune('0'+i))+".java")
		err := os.WriteFile(filename, []byte("public class Test {}\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to create %s: %v", filename, err)
		}
	}

	checker := &Checker{projectPath: tmpDir}
	pt, err := checker.detectProjectType()
	if err != nil {
		t.Fatalf("detectProjectType failed: %v", err)
	}

	if pt != Java {
		t.Errorf("Expected Java, got %v", pt)
	}
}

func TestDetectProjectType_MixedExtensions(t *testing.T) {
	tmpDir := t.TempDir()

	// Create more Python files than others
	for i := 0; i < 5; i++ {
		filename := filepath.Join(tmpDir, "test"+string(rune('0'+i))+".py")
		err := os.WriteFile(filename, []byte("print('hello')\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to create %s: %v", filename, err)
		}
	}

	// Create fewer Go files
	for i := 0; i < 2; i++ {
		filename := filepath.Join(tmpDir, "test"+string(rune('0'+i))+".go")
		err := os.WriteFile(filename, []byte("package main\n"), 0644)
		if err != nil {
			t.Fatalf("Failed to create %s: %v", filename, err)
		}
	}

	checker := &Checker{projectPath: tmpDir}
	pt, err := checker.detectProjectType()
	if err != nil {
		t.Fatalf("detectProjectType failed: %v", err)
	}

	if pt != Python {
		t.Errorf("Expected Python (most files), got %v", pt)
	}
}

func TestType_Values(t *testing.T) {
	tests := []struct {
		name  string
		typ   Type
		value int
	}{
		{"Python", Python, 0},
		{"Go", Go, 1},
		{"Java", Java, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.typ) != tt.value {
				t.Errorf("Expected %d, got %d", tt.value, int(tt.typ))
			}
		})
	}
}

func TestFileCoverage(t *testing.T) {
	fc := FileCoverage{
		File:     "test.go",
		Coverage: 85.5,
	}

	if fc.File != "test.go" {
		t.Errorf("Expected File 'test.go', got '%s'", fc.File)
	}

	if fc.Coverage != 85.5 {
		t.Errorf("Expected Coverage 85.5, got %f", fc.Coverage)
	}
}

func TestFileComplexity(t *testing.T) {
	fc := FileComplexity{
		File:             "test.go",
		AverageCC:        2.5,
		MaxCC:            5,
		ExceedsThreshold: true,
	}

	if fc.File != "test.go" {
		t.Errorf("Expected File 'test.go', got '%s'", fc.File)
	}

	if fc.AverageCC != 2.5 {
		t.Errorf("Expected AverageCC 2.5, got %f", fc.AverageCC)
	}

	if fc.MaxCC != 5 {
		t.Errorf("Expected MaxCC 5, got %d", fc.MaxCC)
	}

	if !fc.ExceedsThreshold {
		t.Error("Expected ExceedsThreshold true")
	}
}

func TestFileViolation(t *testing.T) {
	fv := FileViolation{
		File: "test.go",
		LOC:  250,
	}

	if fv.File != "test.go" {
		t.Errorf("Expected File 'test.go', got '%s'", fv.File)
	}

	if fv.LOC != 250 {
		t.Errorf("Expected LOC 250, got %d", fv.LOC)
	}
}

func TestTypeError(t *testing.T) {
	te := TypeError{
		File:    "test.go",
		Line:    42,
		Message: "type error",
	}

	if te.File != "test.go" {
		t.Errorf("Expected File 'test.go', got '%s'", te.File)
	}

	if te.Line != 42 {
		t.Errorf("Expected Line 42, got %d", te.Line)
	}

	if te.Message != "type error" {
		t.Errorf("Expected Message 'type error', got '%s'", te.Message)
	}
}

func TestCoverageResult(t *testing.T) {
	cr := CoverageResult{
		ProjectType: "Go",
		Coverage:    75.5,
		Threshold:   80.0,
		Passed:      false,
		Report:      "Coverage below threshold",
		FilesBelow: []FileCoverage{
			{File: "test.go", Coverage: 50.0},
		},
	}

	if cr.ProjectType != "Go" {
		t.Errorf("Expected ProjectType 'Go', got '%s'", cr.ProjectType)
	}

	if cr.Coverage != 75.5 {
		t.Errorf("Expected Coverage 75.5, got %f", cr.Coverage)
	}

	if cr.Passed {
		t.Error("Expected Passed false")
	}

	if len(cr.FilesBelow) != 1 {
		t.Errorf("Expected 1 file below threshold, got %d", len(cr.FilesBelow))
	}
}

func TestComplexityResult(t *testing.T) {
	ctr := ComplexityResult{
		AverageCC: 3.5,
		MaxCC:     10,
		Threshold: 5,
		Passed:    false,
		ComplexFiles: []FileComplexity{
			{File: "complex.go", AverageCC: 8.0, MaxCC: 10, ExceedsThreshold: true},
		},
	}

	if ctr.AverageCC != 3.5 {
		t.Errorf("Expected AverageCC 3.5, got %f", ctr.AverageCC)
	}

	if ctr.MaxCC != 10 {
		t.Errorf("Expected MaxCC 10, got %d", ctr.MaxCC)
	}

	if ctr.Passed {
		t.Error("Expected Passed false")
	}
}

func TestFileSizeResult(t *testing.T) {
	fsr := FileSizeResult{
		TotalFiles: 10,
		Violators: []FileViolation{
			{File: "large.go", LOC: 250},
		},
		Warnings: []FileViolation{
			{File: "medium.go", LOC: 180},
		},
		Threshold:  200,
		Passed:     false,
		AverageLOC: 150,
		Strict:     true,
	}

	if fsr.TotalFiles != 10 {
		t.Errorf("Expected TotalFiles 10, got %d", fsr.TotalFiles)
	}

	if fsr.Threshold != 200 {
		t.Errorf("Expected Threshold 200, got %d", fsr.Threshold)
	}

	if !fsr.Strict {
		t.Error("Expected Strict true")
	}
}

func TestTypeResult(t *testing.T) {
	tr := TypeResult{
		ProjectType: "Go",
		Passed:      true,
		Errors: []TypeError{
			{File: "test.go", Line: 10, Message: "error"},
		},
		Warnings: []TypeError{
			{File: "test.go", Line: 20, Message: "warning"},
		},
	}

	if tr.ProjectType != "Go" {
		t.Errorf("Expected ProjectType 'Go', got '%s'", tr.ProjectType)
	}

	if !tr.Passed {
		t.Error("Expected Passed true")
	}

	if len(tr.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(tr.Errors))
	}
}
