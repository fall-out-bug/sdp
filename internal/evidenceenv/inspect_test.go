package evidenceenv

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInspectValid(t *testing.T) {
	// Use template with requirePRURL=false (specs at repo root)
	wd, _ := os.Getwd()
	repoRoot := filepath.Dir(filepath.Dir(wd)) // internal/evidence -> repo
	template := filepath.Join(repoRoot, "specs", "strict-evidence-template.json")
	summary, res, err := Inspect(template, false)
	if err != nil {
		t.Fatalf("Inspect: %v", err)
	}
	if !res.OK {
		t.Fatalf("expected OK, got %v", res)
	}
	if !strings.Contains(summary, "intent") {
		t.Error("summary should include intent")
	}
	if !strings.Contains(summary, "plan") {
		t.Error("summary should include plan")
	}
	if !strings.Contains(summary, "boundary_compliance") {
		t.Error("summary should include boundary_compliance")
	}
	if !strings.Contains(summary, "provenance") {
		t.Error("summary should include provenance")
	}
}

func TestInspectInvalidFile(t *testing.T) {
	_, _, err := Inspect("/nonexistent/path.json", false)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestInspectPromptProvenance(t *testing.T) {
	// Envelope with prompt_hash and context_sources should display in inspect output
	tmp := t.TempDir()
	f := filepath.Join(tmp, "evidence.json")
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
	summary, res, err := Inspect(f, false)
	if err != nil {
		t.Fatalf("Inspect: %v", err)
	}
	if !res.OK {
		t.Fatalf("expected OK: %s", res.Reason)
	}
	if !strings.Contains(summary, "prompt_hash") {
		t.Error("inspect output should include prompt_hash when present")
	}
	if !strings.Contains(summary, "context_sources") {
		t.Error("inspect output should include context_sources when present")
	}
	if !strings.Contains(summary, "workstream_spec") {
		t.Error("inspect output should include context source type")
	}
}

func TestInspectInvalidEvidence(t *testing.T) {
	tmp := t.TempDir()
	bad := filepath.Join(tmp, "bad.json")
	os.WriteFile(bad, []byte(`{"intent":{}}`), 0644)
	summary, res, err := Inspect(bad, false)
	if err != nil {
		t.Fatalf("Inspect should not return error for invalid evidence: %v", err)
	}
	if res.OK {
		t.Fatal("expected !res.OK for invalid evidence")
	}
	if summary != "" {
		t.Error("summary should be empty for invalid evidence")
	}
}
