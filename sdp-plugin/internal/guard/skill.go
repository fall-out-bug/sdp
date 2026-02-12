package guard

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/config"
)

// Skill implements the guard skill logic
type Skill struct {
	stateManager *StateManager
	activeWS     string
}

// NewSkill creates a new guard skill
func NewSkill(configDir string) *Skill {
	return &Skill{
		stateManager: NewStateManager(configDir),
		activeWS:     "",
	}
}

// Activate sets the active workstream
func (s *Skill) Activate(wsID string) error {
	s.activeWS = wsID

	// Load existing state to get scope files
	state, err := s.stateManager.Load()
	if err != nil {
		return fmt.Errorf("failed to load state: %w", err)
	}

	// Create new state with current WS
	newState := GuardState{
		ActiveWS:    wsID,
		ScopeFiles:  state.ScopeFiles,
		ActivatedAt: time.Now().Format(time.RFC3339),
		Timestamp:   "",
	}

	// Save state
	if err := s.stateManager.Save(newState); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	return nil
}

// CheckEdit checks if a file edit is allowed
func (s *Skill) CheckEdit(filePath string) (*GuardResult, error) {
	// Load current state
	state, err := s.stateManager.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load state: %w", err)
	}

	// Check if state is expired
	if state.IsExpired() {
		return &GuardResult{
			Allowed: false,
			Reason:  "No active WS (state expired). Run 'sdp guard activate <ws_id>' first.",
		}, nil
	}

	// No active WS
	if state.ActiveWS == "" {
		return &GuardResult{
			Allowed: false,
			Reason:  "No active WS. Run 'sdp guard activate <ws_id>' first.",
		}, nil
	}

	// No scope restrictions = all files allowed
	if len(state.ScopeFiles) == 0 {
		return &GuardResult{
			Allowed:    true,
			WSID:       state.ActiveWS,
			Reason:     "No scope restrictions",
			ScopeFiles: state.ScopeFiles,
		}, nil
	}

	// Check if file is in scope
	allowed := false
	for _, scopeFile := range state.ScopeFiles {
		if filePath == scopeFile {
			allowed = true
			break
		}
	}

	if !allowed {
		return &GuardResult{
			Allowed:    false,
			WSID:       state.ActiveWS,
			Reason:     fmt.Sprintf("File %s not in WS scope", filePath),
			ScopeFiles: state.ScopeFiles,
		}, nil
	}

	return &GuardResult{
		Allowed:    true,
		WSID:       state.ActiveWS,
		Reason:     "File in scope",
		ScopeFiles: state.ScopeFiles,
	}, nil
}

// GetActiveWS returns the currently active workstream ID
func (s *Skill) GetActiveWS() string {
	state, _ := s.stateManager.Load()
	return state.ActiveWS
}

// Deactivate deactivates the current workstream
func (s *Skill) Deactivate() error {
	if err := s.stateManager.Clear(); err != nil {
		return fmt.Errorf("failed to clear state: %w", err)
	}
	s.activeWS = ""
	return nil
}

// ResolvePath resolves a relative file path to absolute path
func ResolvePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	return filepath.Join(wd, path), nil
}

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

// applyGuardRules applies loaded guard rules to staged files (AC1)
func (s *Skill) applyGuardRules(files []string, rules *config.GuardRules) []Finding {
	var findings []Finding

	for _, file := range files {
		absPath, err := ResolvePath(file)
		if err != nil {
			continue
		}

		// Read file content
		content, err := os.ReadFile(absPath)
		if err != nil {
			continue
		}

		// Check each enabled rule
		for _, rule := range rules.Rules {
			if !rule.Enabled {
				continue
			}

			switch rule.ID {
			case "max-file-loc":
				findings = append(findings, s.checkMaxFileLOC(file, content, rule)...)
			case "coverage-threshold":
				// Coverage is checked at project level, not per file
				continue
			case "max-cyclomatic-complexity":
				// Complexity requires parsing - skip for now
				continue
			case "no-commented-code":
				findings = append(findings, s.checkCommentedCode(file, content, rule)...)
			case "no-orphaned-todos":
				findings = append(findings, s.checkOrphanedTODOs(file, content, rule)...)
			}
		}
	}

	return findings
}

// checkMaxFileLOC checks if file exceeds max lines of code (AC1, AC6)
func (s *Skill) checkMaxFileLOC(file string, content []byte, rule config.GuardRule) []Finding {
	var findings []Finding

	maxLines := 200 // default
	if maxLinesVal, ok := rule.Config["max_lines"]; ok {
		switch v := maxLinesVal.(type) {
		case int:
			maxLines = v
		case float64:
			maxLines = int(v)
		}
	}

	lines := strings.Split(string(content), "\n")
	loc := len(lines)

	if loc > maxLines {
		severity := SeverityWarning
		if rule.Severity == "error" {
			severity = SeverityError
		}

		findings = append(findings, Finding{
			Severity: severity,
			Rule:     rule.ID,
			File:     file,
			Message:  fmt.Sprintf("File exceeds maximum size: %d LOC (threshold: %d)", loc, maxLines),
		})
	}

	return findings
}

// checkCommentedCode checks for commented-out code (AC1)
func (s *Skill) checkCommentedCode(file string, content []byte, rule config.GuardRule) []Finding {
	var findings []Finding
	lines := strings.Split(string(content), "\n")

	commentPrefix := "//"
	if strings.HasSuffix(file, ".py") {
		commentPrefix = "#"
	}

	consecutiveComments := 0
	maxComments := 3 // Threshold for detecting commented code blocks

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, commentPrefix) && len(trimmed) > 10 {
			// Skip actual comments that look like documentation
			if strings.Contains(strings.ToLower(trimmed), "todo") ||
				strings.Contains(strings.ToLower(trimmed), "fixme") ||
				strings.Contains(strings.ToLower(trimmed), "note") {
				consecutiveComments = 0
				continue
			}

			consecutiveComments++
			if consecutiveComments >= maxComments {
				severity := SeverityWarning
				if rule.Severity == "error" {
					severity = SeverityError
				}

				findings = append(findings, Finding{
					Severity: severity,
					Rule:     rule.ID,
					File:     file,
					Line:     i + 1,
					Message:  "Possible commented-out code detected",
				})
				break
			}
		} else {
			consecutiveComments = 0
		}
	}

	return findings
}

// checkOrphanedTODOs checks for TODOs without workstream ID (AC1)
func (s *Skill) checkOrphanedTODOs(file string, content []byte, rule config.GuardRule) []Finding {
	var findings []Finding
	lines := strings.Split(string(content), "\n")

	todoPattern := "TODO"

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, todoPattern) {
			// Check if line contains workstream ID pattern
			hasWSID := false
			for _, wsPattern := range []string{"WS-", "ws-"} {
				if strings.Contains(trimmed, wsPattern) {
					// Verify format WS-XXX-YY
					idx := strings.Index(trimmed, wsPattern)
					if idx+len(wsPattern)+6 <= len(trimmed) {
						potentialID := trimmed[idx : idx+len(wsPattern)+6]
						// Check if it matches pattern like WS-063-03
						if len(potentialID) >= 7 && strings.Count(potentialID, "-") >= 2 {
							hasWSID = true
							break
						}
					}
				}
			}

			if !hasWSID {
				severity := SeverityWarning
				if rule.Severity == "error" {
					severity = SeverityError
				}

				findings = append(findings, Finding{
					Severity: severity,
					Rule:     rule.ID,
					File:     file,
					Line:     i + 1,
					Message:  "TODO without workstream ID (format: TODO(WS-XXX-YY))",
				})
			}
		}
	}

	return findings
}
