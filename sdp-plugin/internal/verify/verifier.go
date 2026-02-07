package verify

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Verifier handles workstream completion verification
type Verifier struct {
	parser *Parser
}

// NewVerifier creates a new workstream verifier
func NewVerifier(wsDir string) *Verifier {
	return &Verifier{
		parser: NewParser(wsDir),
	}
}

// VerifyOutputFiles checks all scope_files exist
func (v *Verifier) VerifyOutputFiles(wsData *WorkstreamData) []CheckResult {
	checks := []CheckResult{}

	for _, filePath := range wsData.ScopeFiles {
		check := CheckResult{
			Name: fmt.Sprintf("File: %s", filePath),
		}

		// Check if file exists
		if _, err := os.Stat(filePath); err == nil {
			// File exists
			check.Passed = true
			check.Message = filePath
			absPath, err := filepath.Abs(filePath)
			if err != nil {
				absPath = filePath // Fall back to original path
			}
			check.Evidence = absPath
		} else {
			// File doesn't exist
			check.Passed = false
			check.Message = fmt.Sprintf("Missing: %s", filePath)
		}

		checks = append(checks, check)
	}

	return checks
}

// VerifyCommands runs verification commands
func (v *Verifier) VerifyCommands(wsData *WorkstreamData) []CheckResult {
	checks := []CheckResult{}

	for _, cmd := range wsData.VerificationCommands {
		check := CheckResult{
			Name: fmt.Sprintf("Command: %s", truncate(cmd, 50)),
		}

		// Run command with timeout
		cmdParts := strings.Fields(cmd)
		if len(cmdParts) == 0 {
			check.Passed = false
			check.Message = "Empty command"
			checks = append(checks, check)
			continue
		}

		// Create context with 60s timeout
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		command := exec.CommandContext(ctx, cmdParts[0], cmdParts[1:]...)
		output, err := command.CombinedOutput()

		if err != nil {
			check.Passed = false
			check.Message = fmt.Sprintf("Exit code: %v", err)
			check.Evidence = truncate(string(output), 500)
		} else {
			check.Passed = true
			check.Message = "Exit code: 0"
			check.Evidence = truncate(string(output), 500)
		}

		checks = append(checks, check)
	}

	return checks
}

// VerifyCoverage checks test coverage (placeholder for now)
func (v *Verifier) VerifyCoverage(wsData *WorkstreamData) *CheckResult {
	if wsData.CoverageThreshold == 0 {
		return nil
	}

	return &CheckResult{
		Name:    "Coverage Check",
		Passed:  true, // Placeholder - would run actual coverage check
		Message: fmt.Sprintf("Coverage threshold: %.1f%%", wsData.CoverageThreshold),
	}
}

// Verify runs all verification checks
func (v *Verifier) Verify(wsID string) *VerificationResult {
	start := time.Now()

	result := &VerificationResult{
		WSID:           wsID,
		Checks:         []CheckResult{},
		MissingFiles:   []string{},
		FailedCommands: []string{},
	}

	// Find WS file
	wsPath, err := v.parser.FindWSFile(wsID)
	if err != nil {
		result.Passed = false
		result.Checks = append(result.Checks, CheckResult{
			Name:    "Find WS",
			Passed:  false,
			Message: err.Error(),
		})
		result.Duration = time.Since(start)
		return result
	}

	// Parse WS file
	wsData, parseErr := v.parser.ParseWSFile(wsPath)
	if parseErr != nil {
		result.Passed = false
		result.Checks = append(result.Checks, CheckResult{
			Name:    "Parse WS",
			Passed:  false,
			Message: parseErr.Error(),
		})
		result.Duration = time.Since(start)
		return result
	}

	// Check 1: Verify output files
	fileChecks := v.VerifyOutputFiles(wsData)
	result.Checks = append(result.Checks, fileChecks...)
	for _, check := range fileChecks {
		if !check.Passed {
			result.MissingFiles = append(result.MissingFiles, check.Message)
		}
	}

	// Check 2: Run verification commands
	cmdChecks := v.VerifyCommands(wsData)
	result.Checks = append(result.Checks, cmdChecks...)
	for _, check := range cmdChecks {
		if !check.Passed {
			result.FailedCommands = append(result.FailedCommands, check.Name)
		}
	}

	// Check 3: Verify coverage
	coverageCheck := v.VerifyCoverage(wsData)
	if coverageCheck != nil {
		result.Checks = append(result.Checks, *coverageCheck)
	}

	// Determine overall pass/fail
	result.Passed = true
	for _, check := range result.Checks {
		if !check.Passed {
			result.Passed = false
			break
		}
	}

	result.Duration = time.Since(start)
	return result
}

// truncate truncates a string to max length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
