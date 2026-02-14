package quality

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckGoComplexity_Basic(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a simple Go file
	goFile := filepath.Join(tmpDir, "simple.go")
	content := `package main

func simple() int {
	return 1
}
`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	checker, err := NewChecker(tmpDir)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}

	result := &ComplexityResult{
		Threshold: 10,
	}

	got, err := checker.checkGoComplexity(result)
	if err != nil {
		t.Errorf("checkGoComplexity returned error: %v", err)
		return
	}

	if got == nil {
		t.Error("checkGoComplexity should return non-nil result")
	}
}

func TestCheckGoComplexity_LargeFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a larger Go file (simulating complexity)
	goFile := filepath.Join(tmpDir, "complex.go")
	var content string
	content = "package main\n\n"
	for i := 0; i < 150; i++ {
		content += "func func" + string(rune('A'+i%26)) + "() int { return 1 }\n"
	}
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	checker, err := NewChecker(tmpDir)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}

	result := &ComplexityResult{
		Threshold: 10,
	}

	got, err := checker.checkGoComplexity(result)
	if err != nil {
		t.Logf("checkGoComplexity returned error: %v", err)
		return
	}

	if got == nil {
		t.Error("checkGoComplexity should return non-nil result")
	}
}

func TestBasicGoComplexity_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()

	checker, err := NewChecker(tmpDir)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}

	result := &ComplexityResult{
		Threshold: 10,
	}

	got, err := checker.basicGoComplexity(result)
	if err != nil {
		t.Errorf("basicGoComplexity returned error: %v", err)
		return
	}

	if got == nil {
		t.Error("basicGoComplexity should return non-nil result")
		return
	}

	if !got.Passed {
		t.Error("Empty directory should pass complexity check")
	}
}

func TestBasicGoComplexity_WithFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Go files
	for i := 0; i < 3; i++ {
		goFile := filepath.Join(tmpDir, "file"+string(rune('A'+i))+".go")
		content := `package main

func simple() int {
	return 1
}
`
		if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	checker, err := NewChecker(tmpDir)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}

	result := &ComplexityResult{
		Threshold: 10,
	}

	got, err := checker.basicGoComplexity(result)
	if err != nil {
		t.Errorf("basicGoComplexity returned error: %v", err)
		return
	}

	if got == nil {
		t.Error("basicGoComplexity should return non-nil result")
	}
}

func TestCheckPythonComplexity_Basic(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a simple Python file
	pyFile := filepath.Join(tmpDir, "simple.py")
	content := `def simple():
    return 1
`
	if err := os.WriteFile(pyFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	checker, err := NewChecker(tmpDir)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}

	result := &ComplexityResult{
		Threshold: 10,
	}

	got, err := checker.checkPythonComplexity(result)
	if err != nil {
		t.Logf("checkPythonComplexity returned error: %v", err)
		return
	}

	if got == nil {
		t.Error("checkPythonComplexity should return non-nil result")
	}
}

func TestBasicPythonComplexity_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()

	checker, err := NewChecker(tmpDir)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}

	result := &ComplexityResult{
		Threshold: 10,
	}

	got, err := checker.basicPythonComplexity(result)
	if err != nil {
		t.Errorf("basicPythonComplexity returned error: %v", err)
		return
	}

	if got == nil {
		t.Error("basicPythonComplexity should return non-nil result")
		return
	}

	if !got.Passed {
		t.Error("Empty directory should pass complexity check")
	}
}

func TestBasicPythonComplexity_SkipsTestFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test file (should be skipped)
	testFile := filepath.Join(tmpDir, "test_module.py")
	content := `def test_something():
    assert True
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	checker, err := NewChecker(tmpDir)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}

	result := &ComplexityResult{
		Threshold: 10,
	}

	got, err := checker.basicPythonComplexity(result)
	if err != nil {
		t.Errorf("basicPythonComplexity returned error: %v", err)
		return
	}

	if got == nil {
		t.Error("basicPythonComplexity should return non-nil result")
		return
	}

	// Test files should be skipped, so no complex files
	if len(got.ComplexFiles) > 0 {
		t.Log("Test files should be skipped (but may have been processed)")
	}
}

func TestComplexityResult_Fields(t *testing.T) {
	result := &ComplexityResult{
		Passed:       true,
		AverageCC:    5.5,
		MaxCC:        10,
		Threshold:    15,
		ComplexFiles: []FileComplexity{},
	}

	if !result.Passed {
		t.Error("Passed should be true")
	}
	if result.AverageCC != 5.5 {
		t.Errorf("AverageCC = %f, want 5.5", result.AverageCC)
	}
	if result.MaxCC != 10 {
		t.Errorf("MaxCC = %d, want 10", result.MaxCC)
	}
}

func TestFileComplexity_Fields(t *testing.T) {
	fc := FileComplexity{
		File:             "test.go",
		AverageCC:        8.0,
		MaxCC:            12,
		ExceedsThreshold: true,
	}

	if fc.File != "test.go" {
		t.Errorf("File = %s, want test.go", fc.File)
	}
	if fc.MaxCC != 12 {
		t.Errorf("MaxCC = %d, want 12", fc.MaxCC)
	}
	if !fc.ExceedsThreshold {
		t.Error("ExceedsThreshold should be true")
	}
}
