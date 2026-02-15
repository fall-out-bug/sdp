package guard

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGovernanceMetaCheck_BlockedFile tests blocking edits to policy files
func TestGovernanceMetaCheck_BlockedFile(t *testing.T) {
	tmpDir := t.TempDir()

	skill := NewSkill(tmpDir)
	stagedFiles := []string{".sdp/guard-rules.yml", "internal/guard/test.go"}

	findings := skill.GovernanceMetaCheck(stagedFiles, "")

	// Should have 1 finding for the policy file
	if len(findings) != 1 {
		t.Errorf("Expected 1 finding, got %d", len(findings))
		return
	}

	if findings[0].Rule != "governance-policy-edit" {
		t.Errorf("Rule = %s, want governance-policy-edit", findings[0].Rule)
	}

	if findings[0].File != ".sdp/guard-rules.yml" {
		t.Errorf("File = %s, want .sdp/guard-rules.yml", findings[0].File)
	}
}

// TestGovernanceMetaCheck_NoPolicyFiles tests no findings when no policy files edited
func TestGovernanceMetaCheck_NoPolicyFiles(t *testing.T) {
	tmpDir := t.TempDir()

	skill := NewSkill(tmpDir)
	stagedFiles := []string{"internal/guard/test.go", "cmd/main.go"}

	findings := skill.GovernanceMetaCheck(stagedFiles, "")

	// Should have no findings
	if len(findings) != 0 {
		t.Errorf("Expected 0 findings, got %d", len(findings))
	}
}

// TestGovernanceMetaCheck_WithApproval tests that approved edits pass
func TestGovernanceMetaCheck_WithApproval(t *testing.T) {
	tmpDir := t.TempDir()

	// Create approval file
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("Failed to create .sdp dir: %v", err)
	}

	approvalContent := `version: 1
approvals:
  - file: ".sdp/guard-rules.yml"
    approved_at: "2025-01-01T00:00:00Z"
    approved_by: "admin"
    reason: "Updating rules for new feature"
    expires_at: ""
`
	approvalPath := filepath.Join(sdpDir, "policy-approvals.yml")
	if err := os.WriteFile(approvalPath, []byte(approvalContent), 0644); err != nil {
		t.Fatalf("Failed to write approval file: %v", err)
	}

	// Change to tmpDir so the approval file can be found
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	skill := NewSkill(tmpDir)
	stagedFiles := []string{".sdp/guard-rules.yml"}

	findings := skill.GovernanceMetaCheck(stagedFiles, "")

	// Should have no findings because file is approved
	if len(findings) != 0 {
		t.Errorf("Expected 0 findings with approval, got %d", len(findings))
	}
}

// TestGetGuardedPolicyFiles tests getting list of guarded files
func TestGetGuardedPolicyFiles(t *testing.T) {
	files := GetGuardedPolicyFiles()

	if len(files) == 0 {
		t.Error("GetGuardedPolicyFiles returned empty list")
	}

	// Check that expected files are in the list
	expected := map[string]bool{
		".sdp/guard-rules.yml": false,
		".sdp/config.yml":      false,
	}

	for _, f := range files {
		if _, ok := expected[f]; ok {
			expected[f] = true
		}
	}

	for f, found := range expected {
		if !found {
			t.Errorf("Expected %s in guarded policy files", f)
		}
	}
}
