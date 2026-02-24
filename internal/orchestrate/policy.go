package orchestrate

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PolicyResult holds the output of OPA policy evaluation.
type PolicyResult struct {
	Denials  []string
	Warnings []string
	Level    string // "advisory" or "blocking"
}

// PolicyInput is the data passed to OPA for evaluation.
type PolicyInput struct {
	Phase                   string   `json:"phase"`
	FeatureID               string   `json:"feature_id"`
	WorkstreamID            string   `json:"workstream_id,omitempty"`
	ChangedFiles            []string `json:"changed_files"`
	ScopeViolationsCount    int      `json:"scope_violations_count"`
	EvidenceFilesCount      int      `json:"evidence_files_count"`
	EvidenceValidationPassed bool    `json:"evidence_validation_passed"`
	HasWorkstreamChanges    bool     `json:"has_workstream_changes"`
	HasFeatureChanges       bool     `json:"has_feature_changes"`
	BeadsReferenced         bool     `json:"beads_referenced"`
	P0Findings              int      `json:"p0_findings"`
	P1Findings              int      `json:"p1_findings"`
	P2Findings              int      `json:"p2_findings"`
}

// EvaluatePolicies evaluates .sdp/policies/*.rego against the given input.
// Returns PolicyResult. If OPA is not installed, returns empty result (graceful degradation).
func EvaluatePolicies(projectRoot string, input PolicyInput) (PolicyResult, error) {
	policiesDir := filepath.Join(projectRoot, ".sdp", "policies")
	if _, err := os.Stat(policiesDir); os.IsNotExist(err) {
		return PolicyResult{Level: "advisory"}, nil
	}

	// Check if opa is available
	opaPath, err := exec.LookPath("opa")
	if err != nil {
		// OPA not installed — skip policy evaluation silently
		return PolicyResult{Level: "advisory"}, nil
	}

	// Write input to temp file
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return PolicyResult{}, fmt.Errorf("marshal policy input: %w", err)
	}
	tmpInput, err := os.CreateTemp("", "sdp-policy-input-*.json")
	if err != nil {
		return PolicyResult{}, fmt.Errorf("create temp input: %w", err)
	}
	defer os.Remove(tmpInput.Name())
	if _, err := tmpInput.Write(inputJSON); err != nil {
		tmpInput.Close()
		return PolicyResult{}, fmt.Errorf("write temp input: %w", err)
	}
	tmpInput.Close()

	result := PolicyResult{}

	// Query enforcement level
	level := queryOPAString(opaPath, policiesDir, tmpInput.Name(), "data.sdp.policies.enforcement_level")
	if level == "" {
		level = "advisory"
	}
	result.Level = level

	// Query effective denials
	denials := queryOPAStringSet(opaPath, policiesDir, tmpInput.Name(), "data.sdp.policies.effective_deny")
	result.Denials = denials

	// Query advisory warnings
	warnings := queryOPAStringSet(opaPath, policiesDir, tmpInput.Name(), "data.sdp.policies.advisory_warn")
	result.Warnings = warnings

	return result, nil
}

func queryOPAString(opaPath, policiesDir, inputFile, query string) string {
	cmd := exec.Command(opaPath, "eval",
		"--data", policiesDir,
		"--input", inputFile,
		"--format", "raw",
		query,
	)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.Trim(strings.TrimSpace(string(out)), `"`)
}

func queryOPAStringSet(opaPath, policiesDir, inputFile, query string) []string {
	cmd := exec.Command(opaPath, "eval",
		"--data", policiesDir,
		"--input", inputFile,
		"--format", "raw",
		query,
	)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	s := strings.TrimSpace(string(out))
	if s == "[]" || s == "" || s == "null" {
		return nil
	}
	var msgs []string
	if json.Unmarshal([]byte(s), &msgs) != nil {
		return nil
	}
	return msgs
}

// BuildPolicyInput constructs a PolicyInput from a checkpoint and scope info.
func BuildPolicyInput(cp *Checkpoint, scopeViolations int, changedFiles []string) PolicyInput {
	wsID := CurrentBuildWS(cp)

	// Check if workstream files changed
	hasWS := false
	hasFeature := false
	for _, f := range changedFiles {
		if strings.HasPrefix(f, "docs/workstreams/") {
			hasWS = true
		}
		if strings.HasPrefix(f, "internal/") || strings.HasPrefix(f, "cmd/") {
			hasFeature = true
		}
	}

	// Check if evidence exists for this feature
	evidenceCount := 0
	evidencePath := fmt.Sprintf(".sdp/evidence/%s.json", cp.FeatureID)
	if _, err := os.Stat(evidencePath); err == nil {
		evidenceCount = 1
	}

	return PolicyInput{
		Phase:                    cp.Phase,
		FeatureID:                cp.FeatureID,
		WorkstreamID:             wsID,
		ChangedFiles:             changedFiles,
		ScopeViolationsCount:     scopeViolations,
		EvidenceFilesCount:       evidenceCount,
		EvidenceValidationPassed: evidenceCount > 0,
		HasWorkstreamChanges:     hasWS,
		HasFeatureChanges:        hasFeature,
		BeadsReferenced:          len(lookupBeadsIDsForFeature(".", cp.FeatureID)) > 0,
	}
}
