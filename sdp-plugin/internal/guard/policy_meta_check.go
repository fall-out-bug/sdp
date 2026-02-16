package guard

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// PolicyApprovalRecord represents an approved policy change
type PolicyApprovalRecord struct {
	File        string   `yaml:"file"`          // File that was approved for edit
	ApprovedAt  string   `yaml:"approved_at"`   // When it was approved
	ApprovedBy  string   `yaml:"approved_by"`   // Who approved it
	Reason      string   `yaml:"reason"`        // Justification
	ExpiresAt   string   `yaml:"expires_at"`    // When approval expires
	CommitHash  string   `yaml:"commit_hash"`   // Expected commit hash after edit
	WorkItemIDs []string `yaml:"work_item_ids"` // Related work items/issues
}

// PolicyApprovalsFile is the structure of .sdp/policy-approvals.yml
type PolicyApprovalsFile struct {
	Version   int                    `yaml:"version"`
	Approvals []PolicyApprovalRecord `yaml:"approvals"`
}

// GovernanceMetaCheck performs governance meta-checks on staged files
// AC4: Governance meta-check blocks direct edits to guarded policy files
func (s *Skill) GovernanceMetaCheck(stagedFiles []string, approvalPath string) []Finding {
	var findings []Finding

	// Load approvals if file exists
	approvals := s.loadPolicyApprovals(approvalPath)

	for _, file := range stagedFiles {
		if IsGuardedPolicyFile(file) {
			// Check if there's an active approval for this file
			if !s.hasActiveApproval(file, approvals) {
				findings = append(findings, Finding{
					Severity: SeverityError,
					Rule:     "governance-policy-edit",
					File:     file,
					Message: fmt.Sprintf(
						"Direct edit to policy file '%s' requires governance approval. "+
							"Submit a proposal and get approval first.", file),
				})
			}
		}
	}

	return findings
}

// loadPolicyApprovals loads policy approvals from file
func (s *Skill) loadPolicyApprovals(path string) []PolicyApprovalRecord {
	if path == "" {
		// Default path
		path = ".sdp/policy-approvals.yml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var approvalsFile PolicyApprovalsFile
	if err := yaml.Unmarshal(data, &approvalsFile); err != nil {
		return nil
	}

	return approvalsFile.Approvals
}

// hasActiveApproval checks if there's a valid, non-expired approval for a file
func (s *Skill) hasActiveApproval(file string, approvals []PolicyApprovalRecord) bool {
	// Normalize file path
	normalizedFile := filepath.ToSlash(filepath.Clean(file))

	for _, approval := range approvals {
		// Normalize approval file path
		approvalFile := filepath.ToSlash(filepath.Clean(approval.File))

		// Check if file matches (exact or glob)
		matched := approvalFile == normalizedFile
		if !matched && strings.Contains(approvalFile, "*") {
			if m, err := filepath.Match(approvalFile, normalizedFile); err == nil {
				matched = m
			}
		}

		if matched {
			// Check if approval has expired
			if approval.ExpiresAt != "" && isApprovalExpired(approval.ExpiresAt) {
				continue
			}
			return true
		}
	}

	return false
}

// isApprovalExpired checks if an approval has expired
func isApprovalExpired(expiresAt string) bool {
	// Simple check - parse as RFC3339 and compare
	// This is a simplified implementation
	return strings.Contains(expiresAt, "expired") ||
		(len(expiresAt) > 0 && expiresAt != "")
	// Full implementation would parse the timestamp and compare with now
}

// GetGuardedPolicyFiles returns the list of guarded policy files
func GetGuardedPolicyFiles() []string {
	files := make([]string, 0, len(GuardedPolicyFiles))
	for f := range GuardedPolicyFiles {
		files = append(files, f)
	}
	return files
}
