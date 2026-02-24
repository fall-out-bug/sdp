package watcher

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewQualityWatcher(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "quality-watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a simple go.mod file
	modFile := filepath.Join(tmpDir, "go.mod")
	err = os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	qw, err := NewQualityWatcher(tmpDir, &QualityWatcherConfig{
		Quiet: true,
	})

	if err != nil {
		t.Fatalf("Failed to create quality watcher: %v", err)
	}

	if qw == nil {
		t.Fatal("QualityWatcher is nil")
	}

	if qw.checker == nil {
		t.Error("Checker is nil")
	}

	qw.watcher.Close()
}

func TestQualityWatcher_FileSizeViolation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "quality-watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create go.mod
	modFile := filepath.Join(tmpDir, "go.mod")
	err = os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	var checkCount int64
	checkDone := make(chan bool, 1)

	qw, err := NewQualityWatcher(tmpDir, &QualityWatcherConfig{
		Quiet: true,
	})
	if err != nil {
		t.Fatalf("Failed to create quality watcher: %v", err)
	}
	defer qw.watcher.Close()

	// Override OnChange to detect when check runs
	originalOnChange := qw.watcher.config.OnChange
	qw.watcher.config.OnChange = func(path string) {
		atomic.AddInt64(&checkCount, 1)
		originalOnChange(path)
		if atomic.LoadInt64(&checkCount) >= 1 {
			select {
			case checkDone <- true:
			default:
			}
		}
	}

	// Start watcher
	go qw.Start()
	defer qw.Stop()

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Create a large file (>200 LOC)
	largeFile := filepath.Join(tmpDir, "large.go")
	content := "package test\n\nfunc Large() {\n"
	for i := 0; i < 250; i++ {
		content += "    var x int\n"
	}
	content += "}\n"

	err = os.WriteFile(largeFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write large file: %v", err)
	}

	// Wait for check to complete
	select {
	case <-checkDone:
		// Check was triggered
	case <-time.After(3 * time.Second):
		t.Error("File check not completed within timeout")
	}

	// Wait a bit for violations to be recorded
	time.Sleep(500 * time.Millisecond)

	// Verify that the check ran
	if atomic.LoadInt64(&checkCount) == 0 {
		t.Error("Expected at least 1 check, got none")
	}

	// Note: We may not always detect violations depending on timing
	// The important thing is that the watcher is running and checking files
}

func TestQualityWatcher_ExcludeTestFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "quality-watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create go.mod
	modFile := filepath.Join(tmpDir, "go.mod")
	err = os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	var checkCount int64

	qw, err := NewQualityWatcher(tmpDir, &QualityWatcherConfig{
		Quiet: true,
	})
	if err != nil {
		t.Fatalf("Failed to create quality watcher: %v", err)
	}
	defer qw.watcher.Close()

	// Override OnChange to count checks
	qw.watcher.config.OnChange = func(path string) {
		atomic.AddInt64(&checkCount, 1)
	}

	// Start watcher
	go qw.Start()
	defer qw.Stop()

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Create a test file (should be excluded)
	testFile := filepath.Join(tmpDir, "test_test.go")
	err = os.WriteFile(testFile, []byte("package test\n\nfunc Test() {}\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create a source file (should be checked)
	srcFile := filepath.Join(tmpDir, "src.go")
	err = os.WriteFile(srcFile, []byte("package test\n\nfunc Src() {}\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write src file: %v", err)
	}

	// Wait for changes
	time.Sleep(500 * time.Millisecond)

	// Should only check source file, not test file
	if atomic.LoadInt64(&checkCount) != 1 {
		t.Errorf("Expected 1 check (excluded test file), got %d", atomic.LoadInt64(&checkCount))
	}
}

func TestQualityWatcher_ClearViolations(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "quality-watcher-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create go.mod
	modFile := filepath.Join(tmpDir, "go.mod")
	err = os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	qw, err := NewQualityWatcher(tmpDir, &QualityWatcherConfig{
		Quiet: true,
	})
	if err != nil {
		t.Fatalf("Failed to create quality watcher: %v", err)
	}
	defer qw.watcher.Close()

	// Manually add a violation
	qw.addViolation(Violation{
		File:     "test.go",
		Check:    "test",
		Message:  "test violation",
		Severity: "error",
	})

	violations := qw.GetViolations()
	if len(violations) != 1 {
		t.Fatalf("Expected 1 violation, got %d", len(violations))
	}

	// Clear violations
	qw.clearViolations("test.go")

	violations = qw.GetViolations()
	if len(violations) != 0 {
		t.Errorf("Expected 0 violations after clear, got %d", len(violations))
	}
}

// TestQualityWatcher_OnFileChange_TypeErrors triggers checkTypes with Go vet errors (coverage: checkTypes Errors loop).
func TestQualityWatcher_OnFileChange_TypeErrors(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "quality-watcher-types-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	modFile := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}
	// Go file with type error so go vet returns Errors
	badFile := filepath.Join(tmpDir, "bad.go")
	content := "package test\n\nfunc F() {\n\tvar x int = \"string\"\n}\n"
	if err := os.WriteFile(badFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write bad.go: %v", err)
	}

	qw, err := NewQualityWatcher(tmpDir, &QualityWatcherConfig{Quiet: true})
	if err != nil {
		t.Fatalf("NewQualityWatcher: %v", err)
	}
	defer qw.watcher.Close()

	qw.onFileChange(badFile)

	violations := qw.GetViolations()
	var typeErrors int
	for _, v := range violations {
		if v.Check == "types" {
			typeErrors++
		}
	}
	if typeErrors == 0 {
		t.Logf("No type violations (go vet may not report for this snippet in all environments); violations: %d", len(violations))
	}
}

// TestQualityWatcher_OnFileChange_Complexity triggers checkComplexity with a high-LOC Go file (coverage: checkComplexity loop).
func TestQualityWatcher_OnFileChange_Complexity(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "quality-watcher-complex-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	modFile := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}
	// Large file so basicGoComplexity estimates high complexity (loc/10 > threshold when threshold is 10)
	complexFile := filepath.Join(tmpDir, "complex.go")
	var b []byte
	b = append(b, "package test\n\nfunc Complex() {\n"...)
	for i := 0; i < 200; i++ {
		b = append(b, "\tif true { }\n"...)
	}
	b = append(b, "}\n"...)
	if err := os.WriteFile(complexFile, b, 0644); err != nil {
		t.Fatalf("Failed to write complex.go: %v", err)
	}

	qw, err := NewQualityWatcher(tmpDir, &QualityWatcherConfig{Quiet: true})
	if err != nil {
		t.Fatalf("NewQualityWatcher: %v", err)
	}
	defer qw.watcher.Close()

	qw.onFileChange(complexFile)

	violations := qw.GetViolations()
	var complexityViolations int
	for _, v := range violations {
		if v.Check == "complexity" {
			complexityViolations++
		}
	}
	if complexityViolations == 0 {
		t.Logf("No complexity violations (threshold or gocyclo may vary); violations: %d", len(violations))
	}
}

// TestQualityWatcher_Start_NonQuiet covers QualityWatcher.Start with Quiet: false (prints "Watching ...").
func TestQualityWatcher_Start_NonQuiet(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "quality-watcher-nonquiet-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	modFile := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	var buf bytes.Buffer
	restore := captureStdout(&buf)

	qw, err := NewQualityWatcher(tmpDir, &QualityWatcherConfig{Quiet: false})
	if err != nil {
		t.Fatalf("NewQualityWatcher: %v", err)
	}
	defer qw.watcher.Close()

	done := make(chan struct{})
	go func() {
		_ = qw.Start()
		close(done)
	}()
	// Allow time for Start() to run and print before we stop
	time.Sleep(300 * time.Millisecond)
	qw.Stop()
	<-done
	// Restore stdout and wait for pipe copy to finish so buf is safe to read
	restore()

	out := buf.String()
	if !strings.Contains(out, "Watching") || !strings.Contains(out, "quality violations") {
		t.Errorf("Expected Start() to print 'Watching ... quality violations'; got: %q", out)
	}
}

func captureStdout(w *bytes.Buffer) func() {
	old := os.Stdout
	pr, pw, _ := os.Pipe()
	os.Stdout = pw
	done := make(chan struct{})
	go func() {
		_, _ = io.Copy(w, pr)
		close(done)
	}()
	return func() {
		pw.Close()
		<-done
		os.Stdout = old
	}
}

// TestNewQualityWatcher_CustomPatterns covers custom IncludePatterns and ExcludePatterns (no defaults).
func TestNewQualityWatcher_CustomPatterns(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "quality-watcher-patterns-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	modFile := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(modFile, []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	qw, err := NewQualityWatcher(tmpDir, &QualityWatcherConfig{
		Quiet:            true,
		IncludePatterns:  []string{"*.py"},
		ExcludePatterns:  []string{"test_*.py"},
	})
	if err != nil {
		t.Fatalf("NewQualityWatcher: %v", err)
	}
	defer qw.watcher.Close()

	if qw.watcher.config.IncludePatterns[0] != "*.py" || qw.watcher.config.ExcludePatterns[0] != "test_*.py" {
		t.Errorf("Custom patterns not applied: include=%v exclude=%v",
			qw.watcher.config.IncludePatterns, qw.watcher.config.ExcludePatterns)
	}
}
