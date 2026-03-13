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

func TestHydrateForReview_AggregatesAcceptanceCriteriaAndScopeFiles(t *testing.T) {
	root := writeHydrateProjectRoot(t)
	second := "---\nws_id: 00-022-02\nfeature_id: F022\ndepends_on: [\"00-022-01\"]\n---\n\n" +
		"# 00-022-02: Review aggregation fixture\n\n" +
		"## Acceptance Criteria\n\n" +
		"- [ ] Third criterion\n" +
		"- [ ] First criterion\n\n" +
		"## Scope Files\n\n" +
		"- `internal/orchestrate/invoke_opencode.go`\n" +
		"- `internal/orchestrate/state_machine.go`\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "workstreams", "backlog", "00-022-02.md"), []byte(second), 0o644); err != nil {
		t.Fatal(err)
	}

	pkt, err := HydrateForReview(root, "F022", &Checkpoint{FeatureID: "F022", Phase: PhaseReview}, []string{"00-022-01", "00-022-02"})
	if err != nil {
		t.Fatalf("HydrateForReview: %v", err)
	}
	if len(pkt.AcceptanceCriteria) != 3 {
		t.Fatalf("acceptance_criteria = %v, want 3 unique entries", pkt.AcceptanceCriteria)
	}
	if len(pkt.ScopeFiles) != 3 {
		t.Fatalf("scope_files = %v, want 3 unique entries", pkt.ScopeFiles)
	}
	if !strings.Contains(pkt.Workstream, "00-022-01") || !strings.Contains(pkt.Workstream, "00-022-02") {
		t.Fatalf("expected aggregated workstream content, got %q", pkt.Workstream)
	}
}

func TestHydrate_FailsWhenQualityGateSourceMissing(t *testing.T) {
	root := t.TempDir()
	wsDir := filepath.Join(root, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	wsContent := "---\nws_id: 00-022-01\n---\n\n## Scope Files\n\n- `internal/orchestrate/hydrate.go`\n\n## Acceptance Criteria\n\n- [ ] First criterion\n"
	if err := os.WriteFile(filepath.Join(wsDir, "00-022-01.md"), []byte(wsContent), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Hydrate(root, "F022", "00-022-01", &Checkpoint{FeatureID: "F022", Phase: PhaseBuild})
	if err == nil {
		t.Fatal("expected error when AGENTS.md is missing")
	}
	if !strings.Contains(err.Error(), "read quality gates source") {
		t.Fatalf("expected quality-gates read error, got %v", err)
	}
}

func TestHydrate_RecordsDriftStatusError(t *testing.T) {
	root := t.TempDir()
	writeHydrateFixture(t, root, "00-022-01", false)

	pkt, err := Hydrate(root, "F022", "00-022-01", &Checkpoint{FeatureID: "F022", Phase: PhaseBuild})
	if err != nil {
		t.Fatalf("Hydrate: %v", err)
	}
	if !strings.Contains(pkt.DriftStatus, "ERROR: collect drift status") {
		t.Fatalf("expected drift status error marker, got %q", pkt.DriftStatus)
	}
}

func TestHydrate_RecordsDependencyLookupError(t *testing.T) {
	root := t.TempDir()
	writeHydrateFixture(t, root, "00-022-01", true)
	mapping := `{"sdp_id":"00-016-04","beads_id":"sdp-missing"}` + "\n"
	if err := os.WriteFile(filepath.Join(root, ".beads-sdp-mapping.jsonl"), []byte(mapping), 0o644); err != nil {
		t.Fatal(err)
	}

	pkt, err := Hydrate(root, "F022", "00-022-01", &Checkpoint{FeatureID: "F022", Phase: PhaseBuild})
	if err != nil {
		t.Fatalf("Hydrate: %v", err)
	}
	msg, ok := pkt.Dependencies["00-016-04"]
	if !ok {
		t.Fatal("expected dependency entry for 00-016-04")
	}
	if !strings.Contains(msg, "ERROR: read dependency 00-016-04 (sdp-missing)") {
		t.Fatalf("expected dependency error marker, got %q", msg)
	}
}

func TestParseWorkstreamSections(t *testing.T) {
	content := "---\nws_id: 00-022-01\ndepends_on: [\"00-016-04\"]\n---\n\n" +
		"## Scope Files\n\n" +
		"- `internal/orchestrate/hydrate.go` — new\n" +
		"- `internal/orchestrate/state_machine.go` — wire\n\n" +
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

func writeHydrateFixture(t *testing.T, root, wsID string, withDependsOn bool) {
	t.Helper()
	wsDir := filepath.Join(root, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	depends := ""
	if withDependsOn {
		depends = "depends_on: [\"00-016-04\"]\n"
	}
	wsContent := "---\nws_id: " + wsID + "\n" + depends + "---\n\n" +
		"## Scope Files\n\n" +
		"- `internal/orchestrate/hydrate.go`\n" +
		"- `internal/orchestrate/invoke_opencode.go`\n\n" +
		"## Acceptance Criteria\n\n" +
		"- [ ] First criterion\n" +
		"- [x] Second criterion\n"
	if err := os.WriteFile(filepath.Join(wsDir, wsID+".md"), []byte(wsContent), 0o644); err != nil {
		t.Fatal(err)
	}
	agents := "# Agents\n\n## Quality Gates\n\nBefore pushing:\n\n```bash\ngo build ./...\n```\n"
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte(agents), 0o644); err != nil {
		t.Fatal(err)
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
	workstream := "---\nws_id: 00-022-01\nfeature_id: F022\ndepends_on: []\n---\n\n" +
		"# 00-022-01: Test hydrate fixture\n\n" +
		"## Acceptance Criteria\n\n" +
		"- [ ] First criterion\n" +
		"- [x] Second criterion\n\n" +
		"## Scope Files\n\n" +
		"- `internal/orchestrate/hydrate.go`\n" +
		"- `internal/orchestrate/state_machine.go`\n"
	if err := os.WriteFile(filepath.Join(wsDir, "00-022-01.md"), []byte(workstream), 0o644); err != nil {
		t.Fatal(err)
	}
	agents := "# Agent Instructions\n\n## Quality Gates\n\nBefore pushing:\n\n```bash\ngo build ./...\n```\n"
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte(agents), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}
