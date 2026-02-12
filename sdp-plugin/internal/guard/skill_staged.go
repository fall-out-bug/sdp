package guard

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fall-out-bug/sdp/internal/config"
)

// ParseCheckOptions parses check options from environment variables (AC6: CI diff-range auto-detection)
func ParseCheckOptions() CheckOptions {
	options := CheckOptions{
		Base: os.Getenv("CI_BASE_SHA"),
		Head: os.Getenv("CI_HEAD_SHA"),
	}

	return options
}

// StagedCheck performs staged guard checks (AC1: staged checks, AC6: CI diff-range, AC8: hybrid mode enforcement)
func (s *Skill) StagedCheck(opts CheckOptions) (*CheckResult, error) {
	// Load guard rules from configuration (AC1)
	projectRoot, err := config.FindProjectRoot()
	if err != nil {
		return &CheckResult{
			Success:  false,
			ExitCode: ExitCodeRuntimeError,
		}, fmt.Errorf("failed to find project root: %w", err)
	}

	guardRules, err := config.LoadGuardRules(projectRoot)
	if err != nil {
		return &CheckResult{
			Success:  false,
			ExitCode: ExitCodeRuntimeError,
		}, fmt.Errorf("failed to load guard rules: %w", err)
	}

	// Get staged files (pass options for CI diff-range support)
	stagedFiles, err := getStagedFiles(opts)
	if err != nil {
		return &CheckResult{
			Success:  false,
			ExitCode: ExitCodeRuntimeError,
		}, fmt.Errorf("failed to get staged files: %w", err)
	}

	// No staged files - pass with empty result
	if len(stagedFiles) == 0 {
		return &CheckResult{
			Success:  true,
			ExitCode: ExitCodePass,
			Findings: []Finding{},
			Summary:  CheckSummary{},
		}, nil
	}

	// Load active workstream scope
	state, err := s.stateManager.Load()
	if err != nil {
		return &CheckResult{
			Success:  false,
			ExitCode: ExitCodeRuntimeError,
		}, fmt.Errorf("failed to load state: %w", err)
	}

	// Check if there's an active workstream
	var findings []Finding

	// If active WS exists, check scope violations
	if state.ActiveWS != "" && !state.IsExpired() && len(state.ScopeFiles) > 0 {
		for _, stagedFile := range stagedFiles {
			// Normalize path
			absPath, err := ResolvePath(stagedFile)
			if err != nil {
				continue // Skip files that can't be resolved
			}

			// Check if file is in scope
			inScope := false
			for _, scopeFile := range state.ScopeFiles {
				if absPath == scopeFile || strings.Contains(absPath, scopeFile) {
					inScope = true
					break
				}
			}

			if !inScope {
				// File is not in scope - this is a WARNING (not blocking)
				// AC8: Hybrid mode - warnings don't block
				findings = append(findings, Finding{
					Severity: SeverityWarning,
					Rule:     "scope-check",
					File:     stagedFile,
					Message:  fmt.Sprintf("File not in active workstream scope (%s)", state.ActiveWS),
				})
			}
		}
	}

	// Apply guard rules to staged files (AC1)
	findings = append(findings, s.applyGuardRules(stagedFiles, guardRules)...)

	// Build result with hybrid mode enforcement (AC8)
	return BuildCheckResult(findings), nil
}

// BuildCheckResult builds a check result with hybrid mode enforcement (AC8: block on ERROR, warn on WARNING)
func BuildCheckResult(findings []Finding) *CheckResult {
	errors := 0
	warnings := 0

	for _, f := range findings {
		switch f.Severity {
		case SeverityError:
			errors++
		case SeverityWarning:
			warnings++
		}
	}

	// AC8: Hybrid mode - only ERROR causes failure
	exitCode := ExitCodePass
	success := true

	if errors > 0 {
		exitCode = ExitCodeViolation
		success = false
	}

	return &CheckResult{
		Success:  success,
		ExitCode: exitCode,
		Findings: findings,
		Summary: CheckSummary{
			Total:    len(findings),
			Errors:   errors,
			Warnings: warnings,
		},
	}
}

// getStagedFiles returns a list of staged files using git diff --cached
// AC6: CI diff-range auto-detection - uses CI_BASE_SHA and CI_HEAD_SHA when available
func getStagedFiles(opts CheckOptions) ([]string, error) {
	var cmd *exec.Cmd

	// Use CI diff range if provided (AC6)
	if opts.Base != "" && opts.Head != "" {
		cmd = exec.Command("git", "diff", "--name-only", "--diff-filter=ACM", opts.Base, opts.Head)
	} else {
		// Use staged files (local)
		cmd = exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=ACM")
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// If git command fails (not in a git repo), return empty list
		return []string{}, nil
	}

	// Parse output
	output := stdout.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	var files []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, line)
		}
	}

	return files, nil
}
