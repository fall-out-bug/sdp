package watcher

import (
	"fmt"
	"path/filepath"
)

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
	qw.checkFileSize(path, relPath)

	// Check complexity
	qw.checkComplexity(path, relPath)

	// Check types
	qw.checkTypes(path, relPath)

	if !qw.quiet {
		fmt.Println("\033[90m────────────────────────────────────\033[0m")
	}
}

func (qw *QualityWatcher) checkFileSize(path, relPath string) {
	sizeResult, err := qw.checker.CheckFileSize()
	if err != nil {
		return
	}
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

func (qw *QualityWatcher) checkComplexity(path, relPath string) {
	complexityResult, err := qw.checker.CheckComplexity()
	if err != nil {
		return
	}
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

func (qw *QualityWatcher) checkTypes(path, relPath string) {
	typeResult, err := qw.checker.CheckTypes()
	if err != nil {
		return
	}
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
