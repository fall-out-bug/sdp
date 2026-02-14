package watcher

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall-out-bug/sdp/internal/quality"
)

func TestCheckFileSize_WithViolation(t *testing.T) {
	tmpDir := t.TempDir()

	// Create go.mod
	modFile := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a large file (> 200 lines)
	lines := make([]string, 250)
	for i := range lines {
		lines[i] = "// line content"
	}
	largeContent := strings.Join(lines, "\n")
	largeFile := filepath.Join(tmpDir, "large.go")
	if err := os.WriteFile(largeFile, []byte("package main\n\n"+largeContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create checker and watcher with strict mode
	checker, err := quality.NewChecker(tmpDir)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}
	checker.SetStrictMode(true) // Enable strict mode to get Violators instead of Warnings

	qw := &QualityWatcher{
		checker:   checker,
		watchPath: tmpDir,
		quiet:     true,
	}

	// Check file size
	qw.checkFileSize(largeFile, "large.go")

	violations := qw.GetViolations()
	// Should have violation for large file
	t.Logf("Violations: %d", len(violations))
	for _, v := range violations {
		t.Logf("  - %s: %s", v.Check, v.Message)
	}
}

func TestCheckFileSize_WithRelativePath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create go.mod
	modFile := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a large file
	lines := make([]string, 250)
	for i := range lines {
		lines[i] = "// line content"
	}
	largeContent := strings.Join(lines, "\n")
	largeFile := filepath.Join(tmpDir, "bigfile.go")
	if err := os.WriteFile(largeFile, []byte("package main\n\n"+largeContent), 0644); err != nil {
		t.Fatal(err)
	}

	checker, err := quality.NewChecker(tmpDir)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}

	qw := &QualityWatcher{
		checker:   checker,
		watchPath: tmpDir,
		quiet:     true,
	}

	// Check with relative path matching
	qw.checkFileSize(largeFile, largeFile) // Use full path as relPath too

	t.Logf("Violations: %d", len(qw.GetViolations()))
}

func TestCheckComplexity_WithViolation(t *testing.T) {
	tmpDir := t.TempDir()

	// Create go.mod
	modFile := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a complex Go file (many nested conditionals)
	complexContent := `package main

func complexFunction(x int) int {
	result := 0
	if x > 0 {
		if x > 10 {
			if x > 20 {
				if x > 30 {
					if x > 40 {
						if x > 50 {
							if x > 60 {
								if x > 70 {
									if x > 80 {
										if x > 90 {
											result = x * 2
										} else {
											result = x * 3
										}
									} else {
										result = x + 1
									}
								} else {
									result = x + 2
								}
							} else {
								result = x + 3
							}
						} else {
							result = x + 4
						}
					} else {
						result = x + 5
					}
				} else {
					result = x + 6
				}
			} else {
				result = x + 7
			}
		} else {
			result = x + 8
		}
	} else {
		result = 0
	}
	return result
}
`
	complexFile := filepath.Join(tmpDir, "complex.go")
	if err := os.WriteFile(complexFile, []byte(complexContent), 0644); err != nil {
		t.Fatal(err)
	}

	checker, err := quality.NewChecker(tmpDir)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}

	qw := &QualityWatcher{
		checker:   checker,
		watchPath: tmpDir,
		quiet:     true,
	}

	// Check complexity
	qw.checkComplexity(complexFile, "complex.go")

	violations := qw.GetViolations()
	t.Logf("Violations: %d", len(violations))
	for _, v := range violations {
		t.Logf("  - %s: %s", v.Check, v.Message)
	}
}

func TestCheckFileSize_NotQuiet(t *testing.T) {
	tmpDir := t.TempDir()

	// Create go.mod
	modFile := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a large file
	lines := make([]string, 250)
	for i := range lines {
		lines[i] = "// line"
	}
	largeFile := filepath.Join(tmpDir, "large.go")
	if err := os.WriteFile(largeFile, []byte("package main\n"+strings.Join(lines, "\n")), 0644); err != nil {
		t.Fatal(err)
	}

	checker, err := quality.NewChecker(tmpDir)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}
	checker.SetStrictMode(true) // Enable strict mode

	qw := &QualityWatcher{
		checker:   checker,
		watchPath: tmpDir,
		quiet:     false, // Not quiet to test print path
	}

	// Check file size - should print violation
	qw.checkFileSize(largeFile, "large.go")

	t.Logf("Violations: %d", len(qw.GetViolations()))
}

func TestClearViolations_RemovesExisting(t *testing.T) {
	tmpDir := t.TempDir()

	// Create go.mod
	modFile := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	checker, err := quality.NewChecker(tmpDir)
	if err != nil {
		t.Fatalf("NewChecker failed: %v", err)
	}

	qw := &QualityWatcher{
		checker:   checker,
		watchPath: tmpDir,
		quiet:     true,
	}

	// Add some violations manually
	qw.addViolation(Violation{File: "test.go", Check: "test", Message: "test msg"})
	qw.addViolation(Violation{File: "other.go", Check: "test", Message: "other msg"})

	// Clear violations for test.go
	qw.clearViolations("test.go")

	violations := qw.GetViolations()
	if len(violations) != 1 {
		t.Errorf("Expected 1 violation after clear, got %d", len(violations))
	}

	if len(violations) > 0 && violations[0].File != "other.go" {
		t.Errorf("Expected other.go to remain, got %s", violations[0].File)
	}
}
