package nextstep

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

func TestBuildStatusView(t *testing.T) {
	root := t.TempDir()
	for _, dir := range []string{".claude", ".sdp", ".beads"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", dir, err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, ".beads", "issues.jsonl"), []byte("{}\n"), 0o644); err != nil {
		t.Fatalf("write beads issues: %v", err)
	}

	rec := &Recommendation{
		Command:    "sdp apply --ws 00-069-01",
		Reason:     "Ready to execute workstream 00-069-01 (F069)",
		Confidence: 0.95,
		Category:   CategoryExecution,
		Version:    ContractVersion,
		Metadata: map[string]any{
			"workstream_id": "00-069-01",
			"feature_id":    "F069",
		},
	}
	rec.enrich()

	view := BuildStatusView(root, ProjectState{
		Workstreams: []WorkstreamStatus{
			{ID: "00-069-01", Status: StatusReady, Priority: 0, Feature: "F069"},
			{ID: "00-069-02", Status: StatusReady, Priority: 1, Feature: "F069", BlockedBy: []string{"00-069-01"}},
		},
		GitStatus: GitStatusInfo{IsRepo: true, Branch: "feature/F069", MainBranch: "main"},
		Config:    ConfigInfo{HasSDPConfig: true, EvidenceEnabled: true},
		Session:   &SessionState{FeatureID: "F069", WorkstreamID: "00-069-01", ExpectedBranch: "feature/F069"},
	}, rec)

	if view.NextAction != "sdp apply --ws 00-069-01" {
		t.Fatalf("NextAction = %q, want sdp apply --ws 00-069-01", view.NextAction)
	}
	if len(view.Workstreams.Ready) != 1 {
		t.Fatalf("ready len = %d, want 1", len(view.Workstreams.Ready))
	}
	if len(view.Workstreams.Blocked) != 1 {
		t.Fatalf("blocked len = %d, want 1", len(view.Workstreams.Blocked))
	}
	if view.NextStep == nil || view.NextStep.ActionID == "" {
		t.Fatalf("NextStep should be enriched: %+v", view.NextStep)
	}
}

func TestStatusAndInstructionSchemasValidateContracts(t *testing.T) {
	schemaRoot := findSchemaRoot(t)
	instructionsPath := filepath.Join(schemaRoot, "contracts", "instructions.schema.json")
	statusViewPath := filepath.Join(schemaRoot, "contracts", "status-view.schema.json")

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("instructions.schema.json", bytes.NewReader(mustReadFile(t, instructionsPath))); err != nil {
		t.Fatalf("add instructions schema: %v", err)
	}
	if err := compiler.AddResource("status-view.schema.json", bytes.NewReader(mustReadFile(t, statusViewPath))); err != nil {
		t.Fatalf("add status-view schema: %v", err)
	}
	instructionsSchema, err := compiler.Compile("instructions.schema.json")
	if err != nil {
		t.Fatalf("compile instructions schema: %v", err)
	}
	statusViewSchema, err := compiler.Compile("status-view.schema.json")
	if err != nil {
		t.Fatalf("compile status-view schema: %v", err)
	}

	rec := &Recommendation{
		Command:    "sdp apply --ws 00-069-01",
		Reason:     "Ready to execute workstream 00-069-01 (F069)",
		Confidence: 0.95,
		Category:   CategoryExecution,
		Version:    ContractVersion,
		Metadata: map[string]any{
			"workstream_id": "00-069-01",
			"feature_id":    "F069",
		},
	}
	rec.enrich()

	view := BuildStatusView(t.TempDir(), ProjectState{
		Workstreams: []WorkstreamStatus{{ID: "00-069-01", Status: StatusReady, Priority: 0, Feature: "F069"}},
		GitStatus:   GitStatusInfo{IsRepo: true},
		Config:      ConfigInfo{HasSDPConfig: true},
	}, rec)

	instructionDoc := mustJSONDoc(t, rec)
	if err := instructionsSchema.Validate(instructionDoc); err != nil {
		t.Fatalf("instructions schema validation failed: %v", err)
	}
	statusDoc := mustJSONDoc(t, view)
	if err := statusViewSchema.Validate(statusDoc); err != nil {
		t.Fatalf("status-view schema validation failed: %v", err)
	}
}

func findSchemaRoot(t *testing.T) string {
	t.Helper()
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	for current := dir; current != filepath.Dir(current); current = filepath.Dir(current) {
		candidate := filepath.Join(current, "schema", "index.json")
		if _, err := os.Stat(candidate); err == nil {
			return filepath.Dir(candidate)
		}
	}
	t.Fatal("could not find schema root")
	return ""
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return content
}

func mustJSONDoc(t *testing.T, value any) any {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}
	var doc any
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("unmarshal json: %v", err)
	}
	return doc
}
