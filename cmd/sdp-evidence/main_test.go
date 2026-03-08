package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var validEvidenceFixture = []byte(`{
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
		"hash_prev": ""
	},
	"trace": {"beads_ids": [], "branch": "main", "commits": [], "pr_url": "https://github.com/org/repo/pull/1"}
}`)

func writeValidEvidenceFixture(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "strict-evidence-template.json")
	if err := os.WriteFile(path, validEvidenceFixture, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestValidateValid(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "sdp-evidence")
	if err := exec.Command("go", "build", "-o", bin, ".").Run(); err != nil {
		t.Fatalf("build: %v", err)
	}
	evidencePath := writeValidEvidenceFixture(t)
	wd, _ := os.Getwd()
	root := filepath.Dir(filepath.Dir(wd))
	cmd := exec.Command(bin, "validate", "--evidence", evidencePath, "--require-pr-url=false")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("validate should succeed: %v\n%s", err, out)
	}
	if string(out) != "valid\n" {
		t.Errorf("expected 'valid', got %q", out)
	}
}

func TestValidateInvalidMissingFile(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "sdp-evidence")
	if err := exec.Command("go", "build", "-o", bin, ".").Run(); err != nil {
		t.Fatalf("build: %v", err)
	}
	wd, _ := os.Getwd()
	root := filepath.Dir(filepath.Dir(wd))
	cmd := exec.Command(bin, "validate", "--evidence", ".sdp/evidence/nonexistent.json")
	cmd.Dir = root
	err := cmd.Run()
	if err == nil {
		t.Fatal("validate should fail for missing file")
	}
}

func TestValidateInvalidEvidence(t *testing.T) {
	tmp := t.TempDir()
	bad := filepath.Join(tmp, "bad.json")
	os.WriteFile(bad, []byte(`{"intent":{}}`), 0644)

	bin := filepath.Join(t.TempDir(), "sdp-evidence")
	if err := exec.Command("go", "build", "-o", bin, ".").Run(); err != nil {
		t.Fatalf("build: %v", err)
	}
	wd, _ := os.Getwd()
	root := filepath.Dir(filepath.Dir(wd))
	cmd := exec.Command(bin, "validate", "--evidence", bad)
	cmd.Dir = root
	err := cmd.Run()
	if err == nil {
		t.Fatal("validate should fail for invalid evidence")
	}
}

func TestInspectValid(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "sdp-evidence")
	if err := exec.Command("go", "build", "-o", bin, ".").Run(); err != nil {
		t.Fatalf("build: %v", err)
	}
	evidencePath := writeValidEvidenceFixture(t)
	wd, _ := os.Getwd()
	root := filepath.Dir(filepath.Dir(wd))
	cmd := exec.Command(bin, "inspect", "--evidence", evidencePath, "--require-pr-url=false")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("inspect should succeed: %v\n%s", err, out)
	}
	if len(out) == 0 {
		t.Error("inspect should print summary")
	}
	if !strings.Contains(string(out), "intent") || !strings.Contains(string(out), "plan") {
		t.Errorf("inspect output should include intent and plan: %s", out)
	}
}

func TestInspectInvalidExitsNonZero(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "sdp-evidence")
	if err := exec.Command("go", "build", "-o", bin, ".").Run(); err != nil {
		t.Fatalf("build: %v", err)
	}
	tmp := t.TempDir()
	bad := filepath.Join(tmp, "bad.json")
	os.WriteFile(bad, []byte(`{"intent":{}}`), 0644)
	wd, _ := os.Getwd()
	root := filepath.Dir(filepath.Dir(wd))
	cmd := exec.Command(bin, "inspect", "--evidence", bad)
	cmd.Dir = root
	err := cmd.Run()
	if err == nil {
		t.Fatal("inspect should fail for invalid evidence")
	}
}
