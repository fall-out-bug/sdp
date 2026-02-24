package evidenceenv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateStrictFile_missing(t *testing.T) {
	_, err := ValidateStrictFile("/nonexistent", false)
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestValidateStrictFile_invalidJSON(t *testing.T) {
	f := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(f, []byte(`{invalid`), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := ValidateStrictFile(f, false)
	if err == nil {
		t.Error("invalid JSON should return error")
	}
}

func TestValidateStrictFile_missingSections(t *testing.T) {
	f := filepath.Join(t.TempDir(), "partial.json")
	if err := os.WriteFile(f, []byte(`{"intent":{}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	r, err := ValidateStrictFile(f, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.OK {
		t.Error("missing sections should not be OK")
	}
	if len(r.Missing) == 0 {
		t.Error("expected missing sections")
	}
}

func TestFormatMissing(t *testing.T) {
	got := FormatMissing([]string{"a", "b"})
	if got != "missing: a, b" {
		t.Errorf("FormatMissing = %q", got)
	}
	got = FormatMissing(nil)
	if got != "" {
		t.Errorf("FormatMissing(nil) = %q", got)
	}
}

func TestValidateStrictFile_promptProvenance(t *testing.T) {
	// Envelope with prompt_hash and context_sources (F026) should validate
	f := filepath.Join(t.TempDir(), "evidence.json")
	payload := `{
		"intent": {"issue_id": "sdp_dev-abc", "trigger": "user", "acceptance": [], "risk_class": "low"},
		"plan": {"workstreams": [], "ordering_rationale": ""},
		"execution": {"claimed_issue_ids": [], "branch": "main", "changed_files": []},
		"verification": {"tests": [], "lint": [], "contracts": [], "coverage": {"value": 80, "threshold": 80}},
		"review": {"self_review": [], "adversarial_review": []},
		"risk_notes": {"residual_risks": [], "out_of_scope": []},
		"boundary": {
			"declared": {"allowed_path_prefixes": [], "control_path_prefixes": [], "forbidden_path_prefixes": [], "role": "", "lane": ""},
			"observed": {"touched_paths": [], "out_of_boundary_paths": []},
			"compliance": {"ok": true, "reason": ""}
		},
		"provenance": {
			"run_id": "run-1",
			"orchestrator": "test",
			"runtime": "local",
			"model": "test",
			"gate_results": [],
			"phase": "execute",
			"role": "coder",
			"captured_at": "2026-01-01T00:00:00Z",
			"source_issue_id": "sdp_dev-abc",
			"artifact_id": "art-1",
			"contract_version": "artifact-provenance/v1",
			"hash_algorithm": "sha256",
			"sequence": 0,
			"payload_digest": "",
			"hash": "",
			"hash_prev": "",
			"prompt_hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			"context_sources": [
				{"type": "workstream_spec", "path": "docs/workstreams/backlog/00-026-01.md", "hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}
			]
		},
		"trace": {"beads_ids": [], "branch": "main", "commits": [], "pr_url": "https://github.com/org/repo/pull/1"}
	}`
	if err := os.WriteFile(f, []byte(payload), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := ValidateStrictFile(f, false)
	if err != nil {
		t.Fatalf("ValidateStrictFile: %v", err)
	}
	if !res.OK {
		t.Errorf("envelope with prompt provenance should validate: %s", res.Reason)
	}
}
