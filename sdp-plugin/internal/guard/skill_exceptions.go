package guard

import (
	"strings"
)

// GuardedPolicyFiles is the list of files requiring governance approval to edit.
// AC4: Governance meta-check blocks direct edits to guarded policy files
var GuardedPolicyFiles = map[string]bool{
	".sdp/guard-rules.yml": true,
	".sdp/config.yml":      true,
}

// ApplyExceptions filters findings based on active exceptions.
// AC3: sdp guard check --staged applies active exceptions before final verdict
func ApplyExceptions(findings []Finding, exceptions []Exception) []Finding {
	var filtered []Finding

	for _, finding := range findings {
		if !isFindingExcepted(finding, exceptions) {
			filtered = append(filtered, finding)
		}
	}

	return filtered
}

// isFindingExcepted checks if a finding is covered by an active exception
func isFindingExcepted(finding Finding, exceptions []Exception) bool {
	for _, exception := range exceptions {
		if exception.IsExpired() {
			continue
		}

		if exception.RuleID != "" && exception.RuleID != finding.Rule {
			continue
		}

		if exception.MatchesFile(finding.File) {
			return true
		}
	}

	return false
}

// GetAppliedExceptions returns info about which exceptions were applied.
// AC6: JSON output includes exception metadata for CI auditability
func GetAppliedExceptions(findings []Finding, exceptions []Exception) []AppliedExceptionInfo {
	var applied []AppliedExceptionInfo

	for _, finding := range findings {
		for _, exception := range exceptions {
			if exception.IsExpired() {
				continue
			}

			if exception.RuleID != "" && exception.RuleID != finding.Rule {
				continue
			}

			if exception.MatchesFile(finding.File) {
				applied = append(applied, AppliedExceptionInfo{
					RuleID: exception.RuleID,
					File:   finding.File,
					Reason: exception.Reason,
					Owner:  exception.Owner,
				})
				break
			}
		}
	}

	return applied
}

// IsGuardedPolicyFile checks if a file is a guarded policy file.
// AC4: Governance meta-check blocks direct edits to guarded policy files
func IsGuardedPolicyFile(filePath string) bool {
	normalized := strings.TrimPrefix(filePath, "./")
	return GuardedPolicyFiles[normalized]
}

// CheckPolicyFileEdit checks if editing a policy file is allowed.
// Returns nil if allowed, or a Finding if blocked.
// AC4: Governance meta-check blocks direct edits to guarded policy files
func CheckPolicyFileEdit(filePath string, hasApproval bool) *Finding {
	if !IsGuardedPolicyFile(filePath) {
		return nil
	}

	if hasApproval {
		return nil
	}

	return &Finding{
		Severity: SeverityError,
		Rule:     "governance-policy-edit",
		File:     filePath,
		Message:  "Direct edit to policy file requires governance approval. Submit a proposal first.",
	}
}

// BuildCheckResultWithExceptions builds a check result with exceptions applied.
// AC5: WARNING and ERROR summaries show applied exception counts
// AC6: JSON output includes exception metadata
func BuildCheckResultWithExceptions(findings []Finding, exceptions []Exception) *CheckResult {
	appliedInfo := GetAppliedExceptions(findings, exceptions)
	appliedCount := len(appliedInfo)

	filteredFindings := ApplyExceptions(findings, exceptions)

	errors := 0
	warnings := 0
	for _, f := range filteredFindings {
		switch f.Severity {
		case SeverityError:
			errors++
		case SeverityWarning:
			warnings++
		}
	}

	exitCode := ExitCodePass
	success := true

	if errors > 0 {
		exitCode = ExitCodeViolation
		success = false
	}

	return &CheckResult{
		Success:  success,
		ExitCode: exitCode,
		Findings: filteredFindings,
		Summary: CheckSummary{
			Total:             len(filteredFindings),
			Errors:            errors,
			Warnings:          warnings,
			AppliedExceptions: appliedCount,
		},
		AppliedExceptions: appliedInfo,
	}
}

