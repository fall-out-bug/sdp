package orchestrate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHydrate(t *testing.T) {
	root := writeHydrateProjectRoot(t)
	cp := &Checkpoint{
		Schema:      "1.0",
		FeatureID:   "F022",
		Branch:      "feature/F022-context-pre-hydration",
		Phase:       PhaseBuild,
		Workstreams: []WSStatus{{ID: "00-022-01", Status: "pending"}},
	}
	pkt, err := Hydrate(root, "F022", "00-022-01", cp)
	if err != nil {
		t.Fatalf("Hydrate: %v", err)
	}
	if pkt.Workstream == "" {
		t.Error("workstream should not be empty")
	}
	if !strings.Contains(pkt.Workstream, "00-022-01") {
		t.Error("workstream should contain 00-022-01")
	}
	if len(pkt.AcceptanceCriteria) == 0 {
		t.Error("acceptance_criteria should not be empty")
	}
	if len(pkt.ScopeFiles) == 0 {
		t.Error("scope_files should not be empty")
	}
	if pkt.Checkpoint == nil {
		t.Error("checkpoint should not be nil")
	}
	if pkt.QualityGates == "" {
		t.Error("quality_gates should not be empty")
	}
	// Validate required fields
	if err := pkt.Validate(); err != nil {
		t.Errorf("Validate: %v", err)
	}
}

func TestHydrate_WritesFile(t *testing.T) {
	root := writeHydrateProjectRoot(t)
	cp := &Checkpoint{FeatureID: "F022", Phase: PhaseBuild}
	pkt, err := Hydrate(root, "F022", "00-022-01", cp)
	if err != nil {
		t.Fatalf("Hydrate: %v", err)
	}
	path := filepath.Join(root, contextPacketPath)
	defer os.Remove(path)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	var loaded ContextPacket
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if loaded.Workstream != pkt.Workstream {
		t.Error("loaded workstream should match")
	}
}

func TestParseWorkstreamSections(t *testing.T) {
	bt := "`" // backtick for path wrapping in markdown
	content := "---\nws_id: 00-022-01\ndepends_on: [\"00-016-04\"]\n---\n\n" +
		"## Scope Files\n\n" +
		"- " + bt + "internal/orchestrate/hydrate.go" + bt + " — new\n" +
		"- " + bt + "internal/orchestrate/state_machine.go" + bt + " — wire\n\n" +
		"## Acceptance Criteria\n\n" +
		"- [ ] First criterion\n" +
		"- [x] Second criterion\n"
	ac, sf := parseWorkstreamSections(content)
	if len(ac) != 2 {
		t.Errorf("acceptance criteria: want 2, got %d: %v", len(ac), ac)
	}
	if len(sf) != 2 {
		t.Errorf("scope files: want 2, got %d: %v", len(sf), sf)
	}
	if sf[0] != "internal/orchestrate/hydrate.go" {
		t.Errorf("scope_files[0] = %q", sf[0])
	}
}

func TestParseQualityGates(t *testing.T) {
	content := "# Agents\n\n## Quality Gates\n\nBefore pushing:\n\n```bash\ngo build ./...\n```\n\n## Other\n"
	got := parseQualityGates(content)
	if !strings.Contains(got, "Quality Gates") {
		t.Errorf("parseQualityGates: want Quality Gates section, got %q", got)
	}
}

func findProjectRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for d := dir; d != "" && d != "/"; d = filepath.Dir(d) {
		if _, err := os.Stat(filepath.Join(d, "docs", "workstreams", "backlog")); err == nil {
			return d
		}
	}
	t.Fatal("project root not found")
	return ""
}

func writeHydrateProjectRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	wsDir := filepath.Join(root, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	workstream := `---
ws_id: 00-022-01
feature_id: F022
depends_on: []
---

# 00-022-01: Test hydrate fixture

## Acceptance Criteria

- [ ] First criterion
- [x] Second criterion

## Scope Files

- ` + "`internal/orchestrate/hydrate.go`" + `
- ` + "`internal/orchestrate/state_machine.go`" + `
`
	if err := os.WriteFile(filepath.Join(wsDir, "00-022-01.md"), []byte(workstream), 0o644); err != nil {
		t.Fatal(err)
	}
	agents := "# Agent Instructions\n\n## Quality Gates\n\nBefore pushing:\n\n```bash\ngo build ./...\n```\n"
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte(agents), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}
