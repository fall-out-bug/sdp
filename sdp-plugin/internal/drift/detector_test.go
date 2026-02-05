package drift

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGenerateDriftReport(t *testing.T) {
	report := &DriftReport{
		WorkstreamID: "00-050-01",
		Timestamp:    time.Now(),
		Issues: []DriftIssue{
			{
				File:     "internal/parser/workstream.go",
				Status:   StatusOK,
				Expected: "File exists",
				Actual:   "Found 5 entities",
			},
			{
				File:          "internal/parser/schema.go",
				Status:        StatusError,
				Expected:      "File exists",
				Actual:        "File not found",
				Recommendation: "Create file: internal/parser/schema.go",
			},
		},
	}

	// Generate verdict
	report.Verdict = report.GenerateVerdict()

	if report.Verdict != "FAIL" {
		t.Errorf("Expected verdict FAIL, got %s", report.Verdict)
	}

	// Check string output
	output := report.String()
	if !strings.Contains(output, "00-050-01") {
		t.Error("Report should contain workstream ID")
	}

	if !strings.Contains(output, "internal/parser/workstream.go") {
		t.Error("Report should list files")
	}

	if !strings.Contains(output, "FAIL") {
		t.Error("Report should show FAIL verdict")
	}
}

func TestDriftReportPass(t *testing.T) {
	report := &DriftReport{
		WorkstreamID: "00-050-01",
		Timestamp:    time.Now(),
		Issues: []DriftIssue{
			{
				File:     "internal/parser/validator.go",
				Status:   StatusOK,
				Expected: "File exists",
				Actual:   "Found 3 entities",
			},
		},
	}

	report.Verdict = report.GenerateVerdict()

	if report.Verdict != "PASS" {
		t.Errorf("Expected verdict PASS, got %s", report.Verdict)
	}
}

func TestDriftReportWarning(t *testing.T) {
	report := &DriftReport{
		WorkstreamID: "00-050-01",
		Timestamp:    time.Now(),
		Issues: []DriftIssue{
			{
				File:          "internal/parser/helper.go",
				Status:        StatusWarning,
				Expected:      "Contains functions",
				Actual:        "No entities found",
				Recommendation: "Add implementation",
			},
		},
	}

	report.Verdict = report.GenerateVerdict()

	if report.Verdict != "WARNING" {
		t.Errorf("Expected verdict WARNING, got %s", report.Verdict)
	}
}

func TestDetectMissingFile(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Create a fake workstream file
	wsContent := `---
ws_id: 00-050-01
feature: F050
---
# Test Workstream

## Scope Files

**Implementation:**
- internal/test/file1.go
- internal/test/file2.go
`
	wsPath := filepath.Join(tmpDir, "00-050-01.md")
	if err := os.WriteFile(wsPath, []byte(wsContent), 0644); err != nil {
		t.Fatalf("Failed to create workstream: %v", err)
	}

	// Create detector
	detector := NewDetector(tmpDir)

	// Detect drift
	report, err := detector.DetectDrift(wsPath)
	if err != nil {
		t.Fatalf("Failed to detect drift: %v", err)
	}

	// Should have 2 error issues (both files missing)
	errorCount := 0
	for _, issue := range report.Issues {
		if issue.Status == StatusError {
			errorCount++
		}
	}

	if errorCount != 2 {
		t.Errorf("Expected 2 errors, got %d", errorCount)
	}

	if report.Verdict != "FAIL" {
		t.Errorf("Expected verdict FAIL, got %s", report.Verdict)
	}
}

func TestDetectExistingFile(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Create test file
	testDir := filepath.Join(tmpDir, "internal", "test")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}

	testFile := filepath.Join(testDir, "file1.go")
	testContent := `package test

func TestFunc() {
	// Test function
}

type TestStruct struct {
	Field string
}
`
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Create workstream
	wsContent := `---
ws_id: 00-050-01
feature: F050
---
# Test Workstream

## Scope Files

**Implementation:**
- internal/test/file1.go
`
	wsPath := filepath.Join(tmpDir, "00-050-01.md")
	if err := os.WriteFile(wsPath, []byte(wsContent), 0644); err != nil {
		t.Fatalf("Failed to create workstream: %v", err)
	}

	// Create detector
	detector := NewDetector(tmpDir)

	// Detect drift
	report, err := detector.DetectDrift(wsPath)
	if err != nil {
		t.Fatalf("Failed to detect drift: %v", err)
	}

	// Should have 1 OK issue
	if len(report.Issues) != 1 {
		t.Errorf("Expected 1 issue, got %d", len(report.Issues))
	}

	if report.Issues[0].Status != StatusOK {
		t.Errorf("Expected status OK, got %s", report.Issues[0].Status)
	}

	if report.Verdict != "PASS" {
		t.Errorf("Expected verdict PASS, got %s", report.Verdict)
	}
}

func TestExtractEntitiesGo(t *testing.T) {
	// Create test Go file
	tmpFile, err := os.CreateTemp("", "*.go")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `package test

func FunctionOne() {
}

func FunctionTwo() {
}

type StructOne struct {
	Field string
}

type InterfaceOne interface {
	Method()
}
`
	if err := os.WriteFile(tmpFile.Name(), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	entities := extractEntities(tmpFile.Name())

	// Should find 2 functions and 2 types
	if len(entities) != 4 {
		t.Errorf("Expected 4 entities, got %d: %v", len(entities), entities)
	}
}

func TestExtractEntitiesPython(t *testing.T) {
	// Create test Python file
	tmpFile, err := os.CreateTemp("", "*.py")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `def function_one():
    pass

def function_two():
    pass

class ClassOne:
    pass

class ClassTwo:
    pass
`
	if err := os.WriteFile(tmpFile.Name(), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	entities := extractEntities(tmpFile.Name())

	// Should find 2 functions and 2 classes
	if len(entities) != 4 {
		t.Errorf("Expected 4 entities, got %d: %v", len(entities), entities)
	}
}

func TestAddIssue(t *testing.T) {
	report := &DriftReport{
		WorkstreamID: "00-050-01",
		Timestamp:    time.Now(),
		Issues:       []DriftIssue{},
	}

	// Add OK issue
	report.AddIssue(DriftIssue{
		File:   "test.go",
		Status: StatusOK,
	})

	if report.Verdict != "PASS" {
		t.Errorf("Expected verdict PASS, got %s", report.Verdict)
	}

	// Add ERROR issue
	report.AddIssue(DriftIssue{
		File:   "missing.go",
		Status: StatusError,
	})

	if report.Verdict != "FAIL" {
		t.Errorf("Expected verdict FAIL, got %s", report.Verdict)
	}
}
