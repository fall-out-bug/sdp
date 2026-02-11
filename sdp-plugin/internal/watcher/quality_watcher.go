package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fall-out-bug/sdp/internal/quality"
)

// QualityWatcher watches files and runs quality checks on changes
type QualityWatcher struct {
	watcher    *Watcher
	checker    *quality.Checker
	violations []Violation
	mu         sync.RWMutex
	quiet      bool
	watchPath  string
}

// Violation represents a quality violation
type Violation struct {
	File     string
	Check    string
	Message  string
	Severity string // "error", "warning"
}

// QualityWatcherConfig holds configuration for quality watcher
type QualityWatcherConfig struct {
	// IncludePatterns specifies glob patterns for files to include
	IncludePatterns []string

	// ExcludePatterns specifies glob patterns for files to exclude
	ExcludePatterns []string

	// Quiet suppresses output
	Quiet bool
}

// NewQualityWatcher creates a new quality watcher
func NewQualityWatcher(watchPath string, config *QualityWatcherConfig) (*QualityWatcher, error) {
	if config == nil {
		config = &QualityWatcherConfig{}
	}

	// Create quality checker
	checker, err := quality.NewChecker(watchPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create quality checker: %w", err)
	}

	// Default patterns if not specified
	includePatterns := config.IncludePatterns
	if len(includePatterns) == 0 {
		// Default to Go files
		includePatterns = []string{"*.go"}
	}

	excludePatterns := config.ExcludePatterns
	if len(excludePatterns) == 0 {
		// Exclude test files
		excludePatterns = []string{"*_test.go", "mock_*.go"}
	}

	qw := &QualityWatcher{
		checker:   checker,
		watchPath: watchPath,
		quiet:     config.Quiet,
	}

	// Create file watcher
	watcher, err := NewWatcher(watchPath, &WatcherConfig{
		IncludePatterns:  includePatterns,
		ExcludePatterns:  excludePatterns,
		DebounceInterval: 100 * time.Millisecond,
		OnChange:         qw.onFileChange,
		OnError: func(err error) {
			if !qw.quiet {
				fmt.Fprintf(os.Stderr, "Watch error: %v\n", err)
			}
		},
	})
	if err != nil {
		return nil, err
	}

	qw.watcher = watcher
	return qw, nil
}

// Start begins watching for file changes and running quality checks
func (qw *QualityWatcher) Start() error {
	if !qw.quiet {
		fmt.Printf("Watching %s for quality violations...\n", qw.watchPath)
		fmt.Println("Press Ctrl+C to stop")
	}

	return qw.watcher.Start()
}

// Stop stops the quality watcher
func (qw *QualityWatcher) Stop() {
	qw.watcher.Stop()
}

// Close closes the quality watcher and releases resources
func (qw *QualityWatcher) Close() {
	qw.watcher.Close()
}

// GetViolations returns all violations detected
func (qw *QualityWatcher) GetViolations() []Violation {
	qw.mu.RLock()
	defer qw.mu.RUnlock()

	violations := make([]Violation, len(qw.violations))
	copy(violations, qw.violations)
	return violations
}

func (qw *QualityWatcher) onFileChange(path string) {
	// Clear previous violations for this file
	qw.clearViolations(path)

	// Run quality checks on changed file
	qw.checkFile(path)
}

func (qw *QualityWatcher) checkFile(path string) {
	relPath, err := filepath.Rel(qw.watchPath, path)
	if err != nil {
		relPath = path
	}

	if !qw.quiet {
		fmt.Printf("\n\033[36mChecking: %s\033[0m\n", relPath)
	}

	// Check file size
	sizeResult, err := qw.checker.CheckFileSize()
	if err == nil {
		for _, violator := range sizeResult.Violators {
			if violator.File == path || violator.File == relPath {
				violation := Violation{
					File:     relPath,
					Check:    "file-size",
					Message:  fmt.Sprintf("File too large: %d LOC (max %d)", violator.LOC, sizeResult.Threshold),
					Severity: "error",
				}
				qw.addViolation(violation)
				if !qw.quiet {
					qw.printViolation(violation)
				}
			}
		}
	}

	// Check complexity
	complexityResult, err := qw.checker.CheckComplexity()
	if err == nil {
		for _, complexFile := range complexityResult.ComplexFiles {
			if complexFile.File == path || complexFile.File == relPath {
				violation := Violation{
					File:     relPath,
					Check:    "complexity",
					Message:  fmt.Sprintf("Cyclomatic complexity too high: %.1f avg, %d max (max %d)", complexFile.AverageCC, complexFile.MaxCC, complexityResult.Threshold),
					Severity: "warning",
				}
				qw.addViolation(violation)
				if !qw.quiet {
					qw.printViolation(violation)
				}
			}
		}
	}

	// Check types
	typeResult, err := qw.checker.CheckTypes()
	if err == nil {
		for _, typeErr := range typeResult.Errors {
			if typeErr.File == path || typeErr.File == relPath {
				violation := Violation{
					File:     relPath,
					Check:    "types",
					Message:  fmt.Sprintf("Line %d: %s", typeErr.Line, typeErr.Message),
					Severity: "error",
				}
				qw.addViolation(violation)
				if !qw.quiet {
					qw.printViolation(violation)
				}
			}
		}

		for _, typeWarn := range typeResult.Warnings {
			if typeWarn.File == path || typeWarn.File == relPath {
				violation := Violation{
					File:     relPath,
					Check:    "types",
					Message:  fmt.Sprintf("Line %d: %s", typeWarn.Line, typeWarn.Message),
					Severity: "warning",
				}
				qw.addViolation(violation)
				if !qw.quiet {
					qw.printViolation(violation)
				}
			}
		}
	}

	if !qw.quiet {
		fmt.Println("\033[90m────────────────────────────────────\033[0m")
	}
}

func (qw *QualityWatcher) addViolation(violation Violation) {
	qw.mu.Lock()
	defer qw.mu.Unlock()

	qw.violations = append(qw.violations, violation)
}

func (qw *QualityWatcher) clearViolations(path string) {
	qw.mu.Lock()
	defer qw.mu.Unlock()

	// Filter out violations for this file
	filtered := make([]Violation, 0, len(qw.violations))
	for _, v := range qw.violations {
		if v.File != path {
			filtered = append(filtered, v)
		}
	}
	qw.violations = filtered
}

func (qw *QualityWatcher) printViolation(violation Violation) {
	var icon, color string
	if violation.Severity == "error" {
		icon = "✖"
		color = "\033[31m" // red
	} else {
		icon = "⚠"
		color = "\033[33m" // yellow
	}

	fmt.Printf("%s%s %s\033[0m\n", color, icon, violation.Check)
	fmt.Printf("  \033[90mFile:\033[0m %s\n", violation.File)
	fmt.Printf("  \033[90mMessage:\033[0m %s\n", violation.Message)
}
