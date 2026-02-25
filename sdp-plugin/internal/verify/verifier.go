package verify

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/config"
)

// VerifierOption configures optional dependencies for Verifier.
type VerifierOption func(*Verifier)

// WithCoverageChecker injects a CoverageChecker. Default: quality.Checker.
func WithCoverageChecker(c CoverageChecker) VerifierOption {
	return func(v *Verifier) { v.coverageChecker = c }
}

// WithPathValidator injects a PathValidator. Default: security.ValidatePathInDirectory.
func WithPathValidator(p PathValidator) VerifierOption {
	return func(v *Verifier) { v.pathValidator = p }
}

// WithCommandRunner injects a CommandRunner. Default: security.SafeCommand.
func WithCommandRunner(r CommandRunner) VerifierOption {
	return func(v *Verifier) { v.commandRunner = r }
}

// Verifier handles workstream completion verification
type Verifier struct {
	parser           *Parser
	coverageChecker  CoverageChecker
	pathValidator    PathValidator
	commandRunner    CommandRunner
}

// NewVerifier creates a new workstream verifier with default implementations.
func NewVerifier(wsDir string) *Verifier {
	return NewVerifierWithOptions(wsDir)
}

// NewVerifierWithOptions creates a verifier with optional injected dependencies.
func NewVerifierWithOptions(wsDir string, opts ...VerifierOption) *Verifier {
	v := &Verifier{
		parser: NewParser(wsDir),
	}
	for _, opt := range opts {
		opt(v)
	}
	return v
}

// VerifyOutputFiles checks all scope_files exist and are within project root (path traversal safety).
func (v *Verifier) VerifyOutputFiles(wsData *WorkstreamData) []CheckResult {
	checks := make([]CheckResult, 0, len(wsData.ScopeFiles))
	projectRoot, rootErr := config.FindProjectRoot()
	if rootErr != nil {
		projectRoot = ""
	}

	pv := v.pathValidator
	if pv == nil {
		pv = defaultPathValidator()
	}

	for _, filePath := range wsData.ScopeFiles {
		check := CheckResult{
			Name: fmt.Sprintf("File: %s", filePath),
		}

		// Path traversal: ensure path is within project root
		if projectRoot != "" {
			if err := pv.ValidatePathInDirectory(projectRoot, filePath); err != nil {
				check.Passed = false
				check.Message = fmt.Sprintf("Path outside project: %v", err)
				checks = append(checks, check)
				continue
			}
		}

		// Check if file exists
		if _, err := os.Stat(filePath); err == nil {
			// File exists
			check.Passed = true
			check.Message = filePath
			absPath, err := filepath.Abs(filePath)
			if err != nil {
				slog.Debug("filepath.Abs failed, using relative path", "path", filePath, "error", err)
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

// VerifyCommands runs verification commands with security validation (sdp-5ho2).
// Derives per-command timeouts from the parent context; respects ctx cancellation for graceful shutdown.
func (v *Verifier) VerifyCommands(ctx context.Context, wsData *WorkstreamData) []CheckResult {
	checks := []CheckResult{}

	for _, cmd := range wsData.VerificationCommands {
		check := CheckResult{
			Name: fmt.Sprintf("Command: %s", truncate(cmd, 50)),
		}

		cmdParts := strings.Fields(cmd)
		if len(cmdParts) == 0 {
			check.Passed = false
			check.Message = "Empty command"
			checks = append(checks, check)
			continue
		}

		timeout := verificationTimeout()
		cmdCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		cr := v.commandRunner
		if cr == nil {
			cr = defaultCommandRunner()
		}
		command, err := cr.SafeCommand(cmdCtx, cmdParts[0], cmdParts[1:]...)
		if err != nil {
			check.Passed = false
			check.Message = fmt.Sprintf("Security validation: %v", err)
			checks = append(checks, check)
			continue
		}

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

// VerifyCoverage runs actual coverage check via CoverageChecker. ctx is used for cancellation.
func (v *Verifier) VerifyCoverage(ctx context.Context, wsData *WorkstreamData) *CheckResult {
	if wsData.CoverageThreshold == 0 {
		return nil
	}

	cc := v.coverageChecker
	if cc == nil {
		projectRoot, rootErr := config.FindProjectRoot()
		if rootErr != nil {
			return &CheckResult{
				Name:    "Coverage Check",
				Passed:  false,
				Message: fmt.Sprintf("project root: %v", rootErr),
			}
		}
		var err error
		cc, err = defaultCoverageChecker(projectRoot)
		if err != nil {
			return &CheckResult{
				Name:    "Coverage Check",
				Passed:  false,
				Message: fmt.Sprintf("checker init: %v", err),
			}
		}
	}

	result, err := cc.CheckCoverage(ctx)
	if err != nil {
		return &CheckResult{
			Name:    "Coverage Check",
			Passed:  false,
			Message: fmt.Sprintf("coverage check: %v", err),
		}
	}

	threshold := result.Threshold
	if wsData.CoverageThreshold > 0 {
		threshold = wsData.CoverageThreshold
	}
	passed := result.Coverage >= threshold

	return &CheckResult{
		Name:     "Coverage Check",
		Passed:   passed,
		Message:  fmt.Sprintf("Coverage: %.1f%% (threshold: %.1f%%)", result.Coverage, threshold),
		Evidence: result.Report,
	}
}

// Verify runs all verification checks. ctx is used for command timeouts and cancellation.
func (v *Verifier) Verify(ctx context.Context, wsID string) *VerificationResult {
	start := time.Now()

	result := &VerificationResult{
		WSID:           wsID,
		Checks:         []CheckResult{},
		MissingFiles:   []string{},
		FailedCommands: []string{},
	}

	if ctx == nil {
		ctx = context.TODO()
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
	cmdChecks := v.VerifyCommands(ctx, wsData)
	result.Checks = append(result.Checks, cmdChecks...)
	for _, check := range cmdChecks {
		if !check.Passed {
			result.FailedCommands = append(result.FailedCommands, check.Name)
		}
	}

	// Check 3: Verify coverage
	coverageCheck := v.VerifyCoverage(ctx, wsData)
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

// verificationTimeout returns verification command timeout from config (or env, or default).
func verificationTimeout() time.Duration {
	root, err := config.FindProjectRoot()
	if err != nil {
		return config.TimeoutFromEnv("SDP_TIMEOUT_VERIFICATION", 60*time.Second)
	}
	cfg, err := config.Load(root)
	if err != nil || cfg == nil {
		return config.TimeoutFromEnv("SDP_TIMEOUT_VERIFICATION", 60*time.Second)
	}
	return config.TimeoutFromConfigOrEnv(cfg.Timeouts.VerificationCommand, "SDP_TIMEOUT_VERIFICATION", 60*time.Second)
}

// truncate truncates a string to max length
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
