package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall-out-bug/sdp/internal/nextstep"
)

func TestStatusCmdTextMode(t *testing.T) {
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	createStatusFixture(t, tmpDir)

	cmd := statusCmd()
	if err := cmd.Flags().Set("text", "true"); err != nil {
		t.Fatalf("set text flag: %v", err)
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := cmd.RunE(cmd, []string{})
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("statusCmd() error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	for _, snippet := range []string{"SDP Project Status", "Environment:", "Workstreams:", "Next Action:", "Command:"} {
		if !strings.Contains(output, snippet) {
			t.Fatalf("output missing %q:\n%s", snippet, output)
		}
	}
	if !strings.Contains(output, "sdp apply --ws 00-069-01") {
		t.Fatalf("output should include recommended command, got:\n%s", output)
	}
}

func TestStatusCmdJSONMode(t *testing.T) {
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	createStatusFixture(t, tmpDir)

	cmd := statusCmd()
	if err := cmd.Flags().Set("json", "true"); err != nil {
		t.Fatalf("set json flag: %v", err)
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := cmd.RunE(cmd, []string{})
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("statusCmd() error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	var view nextstep.StatusView
	if err := json.Unmarshal(buf.Bytes(), &view); err != nil {
		t.Fatalf("status JSON invalid: %v\n%s", err, buf.String())
	}
	if !view.HasGit || !view.HasClaude || !view.HasSDP {
		t.Fatalf("environment flags not populated: %+v", view)
	}
	if view.NextAction == "" || view.NextStep == nil {
		t.Fatalf("next-step contract missing: %+v", view)
	}
	if view.NextStep.ActionID == "" {
		t.Fatalf("action_id should be populated: %+v", view.NextStep)
	}
	if len(view.Workstreams.Ready) != 1 {
		t.Fatalf("ready workstreams = %d, want 1", len(view.Workstreams.Ready))
	}
	if len(view.Workstreams.Blocked) != 1 {
		t.Fatalf("blocked workstreams = %d, want 1", len(view.Workstreams.Blocked))
	}
	if view.ActiveSession == nil || view.ActiveSession.WorkstreamID != "00-069-01" {
		t.Fatalf("active session not surfaced: %+v", view.ActiveSession)
	}
}

func TestPrintStatusJSON(t *testing.T) {
	view := &nextstep.StatusView{
		Version:    nextstep.ContractVersion,
		HasGit:     true,
		HasClaude:  true,
		HasSDP:     true,
		HasBeads:   false,
		NextAction: "sdp status",
		NextStep: &nextstep.Recommendation{
			ActionID:   "sdp.status",
			Command:    "sdp status",
			Reason:     "Check current project state",
			Confidence: 0.7,
			Category:   nextstep.CategoryInformation,
			Version:    nextstep.ContractVersion,
		},
	}

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := printStatusJSON(view)
	w.Close()
	os.Stdout = old
	if err != nil {
		t.Fatalf("printStatusJSON error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	if !strings.Contains(output, `"next_action": "sdp status"`) {
		t.Fatalf("json output missing next_action: %s", output)
	}
	if !strings.Contains(output, `"action_id": "sdp.status"`) {
		t.Fatalf("json output missing action_id: %s", output)
	}
}

func TestBoolIcon(t *testing.T) {
	if boolIcon(true) != "[OK]" {
		t.Fatal("boolIcon(true) should be [OK]")
	}
	if boolIcon(false) != "[MISSING]" {
		t.Fatal("boolIcon(false) should be [MISSING]")
	}
}

func createStatusFixture(t *testing.T, root string) {
	t.Helper()
	for _, dir := range []string{".git", ".claude", ".sdp", "docs/workstreams/backlog"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, ".beads"), 0o755); err != nil {
		t.Fatalf("mkdir .beads: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".beads", "issues.jsonl"), []byte("{}\n"), 0o644); err != nil {
		t.Fatalf("write beads issues: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".sdp", "config.yml"), []byte("version: 1\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".sdp", "session.json"), []byte(`{"workstream_id":"00-069-01","feature_id":"F069","branch":"feature/F069"}`), 0o644); err != nil {
		t.Fatalf("write session: %v", err)
	}
	ready := `---
ws_id: 00-069-01
feature_id: F069
status: ready
priority: 0
size: M
depends_on: []
---

## Goal

Ready workstream.

## Acceptance Criteria

- [ ] AC1: Execute
`
	blocked := `---
ws_id: 00-069-02
feature_id: F069
status: ready
priority: 1
size: S
depends_on: ["00-069-01"]
---

## Goal

Blocked workstream.

## Acceptance Criteria

- [ ] AC1: Wait
`
	if err := os.WriteFile(filepath.Join(root, "docs", "workstreams", "backlog", "00-069-01.md"), []byte(ready), 0o644); err != nil {
		t.Fatalf("write ready workstream: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "workstreams", "backlog", "00-069-02.md"), []byte(blocked), 0o644); err != nil {
		t.Fatalf("write blocked workstream: %v", err)
	}
}
