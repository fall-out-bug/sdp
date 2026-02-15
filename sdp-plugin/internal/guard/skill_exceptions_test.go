package guard

import (
	"testing"
	"time"
)

// TestAC1_ExceptionFields tests that exception entries support all required fields
// AC1: Exception entries support rule_id, path glob, reason, owner, expires_at
func TestAC1_ExceptionFields(t *testing.T) {
	now := time.Now()
	futureExpiry := now.Add(24 * time.Hour)

	exception := Exception{
		RuleID:    "max-file-loc",
		PathGlob:  "internal/large/*.go",
		Reason:    "Legacy file migration in progress",
		Owner:     "team-backend",
		ExpiresAt: futureExpiry.Format(time.RFC3339),
	}

	if exception.RuleID != "max-file-loc" {
		t.Errorf("RuleID = %s, want max-file-loc", exception.RuleID)
	}
	if exception.PathGlob != "internal/large/*.go" {
		t.Errorf("PathGlob = %s, want internal/large/*.go", exception.PathGlob)
	}
	if exception.Reason != "Legacy file migration in progress" {
		t.Errorf("Reason = %s, want 'Legacy file migration in progress'", exception.Reason)
	}
	if exception.Owner != "team-backend" {
		t.Errorf("Owner = %s, want team-backend", exception.Owner)
	}
	if exception.ExpiresAt == "" {
		t.Error("ExpiresAt should not be empty")
	}
}

// TestAC2_ExpiredExceptionIgnored tests that expired exceptions are ignored
// AC2: Expired exceptions are ignored automatically
func TestAC2_ExpiredExceptionIgnored(t *testing.T) {
	pastTime := time.Now().Add(-1 * time.Hour) // Expired 1 hour ago

	exception := Exception{
		RuleID:    "max-file-loc",
		PathGlob:  "test.go",
		Reason:    "Test exception",
		Owner:     "test",
		ExpiresAt: pastTime.Format(time.RFC3339),
	}

	if !exception.IsExpired() {
		t.Error("Exception should be expired")
	}

	// Future exception should not be expired
	futureTime := time.Now().Add(24 * time.Hour)
	futureException := Exception{
		RuleID:    "max-file-loc",
		PathGlob:  "test.go",
		Reason:    "Test exception",
		Owner:     "test",
		ExpiresAt: futureTime.Format(time.RFC3339),
	}

	if futureException.IsExpired() {
		t.Error("Future exception should not be expired")
	}
}

// TestAC2_ExceptionMatchesFile tests that exception path glob matches files
func TestAC2_ExceptionMatchesFile(t *testing.T) {
	tests := []struct {
		name      string
		pathGlob  string
		filePath  string
		wantMatch bool
	}{
		{
			name:      "Exact match",
			pathGlob:  "test.go",
			filePath:  "test.go",
			wantMatch: true,
		},
		{
			name:      "Glob match with star",
			pathGlob:  "internal/**/*.go",
			filePath:  "internal/guard/test.go",
			wantMatch: true,
		},
		{
			name:      "No match",
			pathGlob:  "internal/**/*.go",
			filePath:  "cmd/main.go",
			wantMatch: false,
		},
		{
			name:      "Star matches any file",
			pathGlob:  "*",
			filePath:  "anyfile.go",
			wantMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exception := Exception{PathGlob: tt.pathGlob}
			got := exception.MatchesFile(tt.filePath)
			if got != tt.wantMatch {
				t.Errorf("MatchesFile(%s) = %v, want %v", tt.filePath, got, tt.wantMatch)
			}
		})
	}
}

// TestAC3_ApplyExceptions tests that staged check applies active exceptions
// AC3: sdp guard check --staged applies active exceptions before final verdict
func TestAC3_ApplyExceptions(t *testing.T) {
	exceptions := []Exception{
		{
			RuleID:    "max-file-loc",
			PathGlob:  "large_file.go",
			Reason:    "Temporary exception for migration",
			Owner:     "team-xyz",
			ExpiresAt: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		},
	}

	findings := []Finding{
		{
			Severity: SeverityError,
			Rule:     "max-file-loc",
			File:     "large_file.go",
			Message:  "File exceeds maximum size",
		},
		{
			Severity: SeverityError,
			Rule:     "max-file-loc",
			File:     "other_file.go",
			Message:  "File exceeds maximum size",
		},
	}

	filtered := ApplyExceptions(findings, exceptions)

	// The finding for large_file.go should be filtered out
	if len(filtered) != 1 {
		t.Errorf("ApplyExceptions returned %d findings, want 1", len(filtered))
	}

	if len(filtered) > 0 && filtered[0].File != "other_file.go" {
		t.Errorf("Filtered finding file = %s, want other_file.go", filtered[0].File)
	}
}

// TestAC3_ExpiredExceptionNotApplied tests that expired exceptions don't filter findings
func TestAC3_ExpiredExceptionNotApplied(t *testing.T) {
	// Expired exception
	exceptions := []Exception{
		{
			RuleID:    "max-file-loc",
			PathGlob:  "large_file.go",
			Reason:    "Expired exception",
			Owner:     "team-xyz",
			ExpiresAt: time.Now().Add(-1 * time.Hour).Format(time.RFC3339), // Expired
		},
	}

	findings := []Finding{
		{
			Severity: SeverityError,
			Rule:     "max-file-loc",
			File:     "large_file.go",
			Message:  "File exceeds maximum size",
		},
	}

	filtered := ApplyExceptions(findings, exceptions)

	// Expired exception should not filter the finding
	if len(filtered) != 1 {
		t.Errorf("ApplyExceptions with expired exception returned %d findings, want 1", len(filtered))
	}
}

// TestAC4_GovernanceMetaCheck tests that policy file edits are blocked without approval
// AC4: Governance meta-check blocks direct edits to guarded policy files
func TestAC4_GovernanceMetaCheck(t *testing.T) {
	tests := []struct {
		name         string
		file         string
		hasApproval  bool
		wantBlocked  bool
		wantSeverity Severity
	}{
		{
			name:         "Policy file without approval - blocked",
			file:         ".sdp/guard-rules.yml",
			hasApproval:  false,
			wantBlocked:  true,
			wantSeverity: SeverityError,
		},
		{
			name:         "Policy file with approval - allowed",
			file:         ".sdp/guard-rules.yml",
			hasApproval:  true,
			wantBlocked:  false,
			wantSeverity: "",
		},
		{
			name:         "Non-policy file - allowed",
			file:         "internal/guard/test.go",
			hasApproval:  false,
			wantBlocked:  false,
			wantSeverity: "",
		},
		{
			name:         "Config file without approval - blocked",
			file:         ".sdp/config.yml",
			hasApproval:  false,
			wantBlocked:  true,
			wantSeverity: SeverityError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckPolicyFileEdit(tt.file, tt.hasApproval)

			if tt.wantBlocked {
				if result == nil {
					t.Error("CheckPolicyFileEdit returned nil, expected finding")
					return
				}
				if result.Severity != tt.wantSeverity {
					t.Errorf("Severity = %s, want %s", result.Severity, tt.wantSeverity)
				}
			} else {
				if result != nil {
					t.Errorf("CheckPolicyFileEdit returned finding %v, expected nil", result)
				}
			}
		})
	}
}

// TestAC4_GuardedPolicyFiles tests that correct files are identified as policy files
func TestAC4_GuardedPolicyFiles(t *testing.T) {
	tests := []struct {
		file     string
		isPolicy bool
	}{
		{".sdp/guard-rules.yml", true},
		{".sdp/config.yml", true},
		{".sdp/other.yml", false},
		{"internal/guard/rules.go", false},
		{".github/workflows/ci.yml", false},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			got := IsGuardedPolicyFile(tt.file)
			if got != tt.isPolicy {
				t.Errorf("IsGuardedPolicyFile(%s) = %v, want %v", tt.file, got, tt.isPolicy)
			}
		})
	}
}

// TestAC5_ExceptionSummary tests that summaries show exception counts
// AC5: WARNING and ERROR summaries show applied exception counts
func TestAC5_ExceptionSummary(t *testing.T) {
	now := time.Now()
	exceptions := []Exception{
		{
			RuleID:    "max-file-loc",
			PathGlob:  "large_file.go",
			Reason:    "Test",
			Owner:     "test",
			ExpiresAt: now.Add(24 * time.Hour).Format(time.RFC3339),
		},
	}

	findings := []Finding{
		{Severity: SeverityError, Rule: "max-file-loc", File: "large_file.go", Message: "Too large"},
		{Severity: SeverityWarning, Rule: "no-todos", File: "test.go", Message: "Has TODO"},
	}

	result := BuildCheckResultWithExceptions(findings, exceptions)

	// Check that summary includes exception count
	if result.Summary.AppliedExceptions != 1 {
		t.Errorf("AppliedExceptions = %d, want 1", result.Summary.AppliedExceptions)
	}

	// Check that the finding was suppressed
	if result.Summary.Errors != 0 {
		t.Errorf("Errors after exceptions = %d, want 0", result.Summary.Errors)
	}
}

// TestAC6_JSONOutputIncludesExceptions tests that JSON output includes exception metadata
// AC6: JSON output includes exception metadata for CI auditability
func TestAC6_JSONOutputIncludesExceptions(t *testing.T) {
	now := time.Now()
	exceptions := []Exception{
		{
			RuleID:    "max-file-loc",
			PathGlob:  "large_file.go",
			Reason:    "Test exception",
			Owner:     "test-owner",
			ExpiresAt: now.Add(24 * time.Hour).Format(time.RFC3339),
		},
	}

	result := BuildCheckResultWithExceptions([]Finding{
		{Severity: SeverityError, Rule: "max-file-loc", File: "large_file.go", Message: "Too large"},
	}, exceptions)

	// Verify AppliedExceptionsInfo is populated
	if len(result.AppliedExceptions) == 0 {
		t.Error("AppliedExceptions should not be empty")
		return
	}

	info := result.AppliedExceptions[0]
	if info.RuleID != "max-file-loc" {
		t.Errorf("Exception RuleID = %s, want max-file-loc", info.RuleID)
	}
	if info.Reason != "Test exception" {
		t.Errorf("Exception Reason = %s, want 'Test exception'", info.Reason)
	}
	if info.Owner != "test-owner" {
		t.Errorf("Exception Owner = %s, want test-owner", info.Owner)
	}
}
