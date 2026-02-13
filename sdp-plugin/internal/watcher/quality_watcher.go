package watcher

import (
	"fmt"
	"os"
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
