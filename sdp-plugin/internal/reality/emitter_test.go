package reality

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

var requiredOSSArtifacts = []string{
	".sdp/reality/reality-summary.json",
	".sdp/reality/feature-inventory.json",
	".sdp/reality/architecture-map.json",
	".sdp/reality/integration-map.json",
	".sdp/reality/quality-report.json",
	".sdp/reality/drift-report.json",
	".sdp/reality/readiness-report.json",
	"docs/reality/summary.md",
	"docs/reality/architecture.md",
	"docs/reality/quality.md",
	"docs/reality/bootstrap.md",
}

func TestEmitOSS_WritesRequiredArtifactsAndIsDeterministic(t *testing.T) {
	projectRoot := t.TempDir()
	seedProject(t, projectRoot)

	paths, err := EmitOSS(projectRoot)
	if err != nil {
		t.Fatalf("EmitOSS first run failed: %v", err)
	}
	if !slices.Equal(paths, requiredOSSArtifacts) {
		t.Fatalf("EmitOSS returned unexpected artifact list: %v", paths)
	}

	firstRun := make(map[string][]byte, len(requiredOSSArtifacts))
	for _, rel := range requiredOSSArtifacts {
		abs := filepath.Join(projectRoot, rel)
		data, err := os.ReadFile(abs)
		if err != nil {
			t.Fatalf("expected artifact missing: %s (%v)", rel, err)
		}
		if len(data) == 0 {
			t.Fatalf("artifact is empty: %s", rel)
		}
		firstRun[rel] = data
	}

	readinessPath := filepath.Join(projectRoot, ".sdp/reality/readiness-report.json")
	readinessData, err := os.ReadFile(readinessPath)
	if err != nil {
		t.Fatalf("read readiness report: %v", err)
	}
	var readiness map[string]any
	if err := json.Unmarshal(readinessData, &readiness); err != nil {
		t.Fatalf("parse readiness report: %v", err)
	}
	if readiness["verdict"] != "ready" {
		t.Fatalf("unexpected readiness verdict: %v", readiness["verdict"])
	}

	if _, err := EmitOSS(projectRoot); err != nil {
		t.Fatalf("EmitOSS second run failed: %v", err)
	}
	for _, rel := range requiredOSSArtifacts {
		abs := filepath.Join(projectRoot, rel)
		data, err := os.ReadFile(abs)
		if err != nil {
			t.Fatalf("read artifact after second run (%s): %v", rel, err)
		}
		if string(firstRun[rel]) != string(data) {
			t.Fatalf("artifact changed between deterministic runs: %s", rel)
		}
	}
}

func TestEmitOSSWithOptions_AnnotatesModeFocusAndIntegrations(t *testing.T) {
	projectRoot := t.TempDir()
	seedProject(t, projectRoot)
	writeFile(t, filepath.Join(projectRoot, "configs", "app.yaml"), "database: postgres\ncache: redis\n")

	if _, err := EmitOSSWithOptions(projectRoot, Options{Mode: ModeQuick, Focus: "docs"}); err != nil {
		t.Fatalf("EmitOSSWithOptions failed: %v", err)
	}

	summaryData, err := os.ReadFile(filepath.Join(projectRoot, ".sdp", "reality", "reality-summary.json"))
	if err != nil {
		t.Fatalf("read summary: %v", err)
	}
	var summary map[string]any
	if err := json.Unmarshal(summaryData, &summary); err != nil {
		t.Fatalf("parse summary: %v", err)
	}
	scope, ok := summary["scope"].(map[string]any)
	if !ok {
		t.Fatalf("scope missing or invalid: %#v", summary["scope"])
	}
	if scope["mode"] != "quick" {
		t.Fatalf("unexpected mode: %v", scope["mode"])
	}
	if scope["focus"] != "docs" {
		t.Fatalf("unexpected focus: %v", scope["focus"])
	}

	integrationData, err := os.ReadFile(filepath.Join(projectRoot, ".sdp", "reality", "integration-map.json"))
	if err != nil {
		t.Fatalf("read integration map: %v", err)
	}
	var integrationMap map[string]any
	if err := json.Unmarshal(integrationData, &integrationMap); err != nil {
		t.Fatalf("parse integration map: %v", err)
	}
	integrations, ok := integrationMap["integrations"].([]any)
	if !ok || len(integrations) == 0 {
		t.Fatalf("expected integrations to be detected: %#v", integrationMap["integrations"])
	}
}

func TestEmitOSSWithOptions_DetectsDocumentationDriftAndBootstrapRecommendations(t *testing.T) {
	projectRoot := t.TempDir()
	seedProject(t, projectRoot)
	writeFile(t, filepath.Join(projectRoot, "docs", "architecture.md"), "# Architecture\nUses `internal/missing/service.go`.\n")

	if _, err := EmitOSSWithOptions(projectRoot, Options{Mode: ModeBootstrapSDP, Focus: "architecture"}); err != nil {
		t.Fatalf("EmitOSSWithOptions bootstrap failed: %v", err)
	}

	readinessData, err := os.ReadFile(filepath.Join(projectRoot, ".sdp", "reality", "readiness-report.json"))
	if err != nil {
		t.Fatalf("read readiness report: %v", err)
	}
	var readiness map[string]any
	if err := json.Unmarshal(readinessData, &readiness); err != nil {
		t.Fatalf("parse readiness report: %v", err)
	}
	if readiness["verdict"] != "ready_with_constraints" {
		t.Fatalf("unexpected readiness verdict: %v", readiness["verdict"])
	}
	recommendations, ok := readiness["suggested_workstreams"].([]any)
	if !ok || len(recommendations) == 0 {
		t.Fatalf("expected bootstrap recommendations: %#v", readiness["suggested_workstreams"])
	}

	driftData, err := os.ReadFile(filepath.Join(projectRoot, ".sdp", "reality", "drift-report.json"))
	if err != nil {
		t.Fatalf("read drift report: %v", err)
	}
	var drift map[string]any
	if err := json.Unmarshal(driftData, &drift); err != nil {
		t.Fatalf("parse drift report: %v", err)
	}
	contradictions, ok := drift["contradictions"].([]any)
	if !ok || len(contradictions) == 0 {
		t.Fatalf("expected documentation drift contradictions: %#v", drift["contradictions"])
	}
}

func seedProject(t *testing.T, root string) {
	t.Helper()

	writeFile(t, filepath.Join(root, "cmd", "app", "main.go"), "package main\n\nfunc main() {}\n")
	writeFile(t, filepath.Join(root, "internal", "pkg", "logic.go"), "package pkg\n\nfunc Sum(a, b int) int { return a + b }\n")
	writeFile(t, filepath.Join(root, "internal", "pkg", "logic_test.go"), "package pkg\n\nimport \"testing\"\n\nfunc TestSum(t *testing.T) {\n\tif Sum(2, 2) != 4 {\n\t\tt.Fatal(\"unexpected\")\n\t}\n}\n")
	writeFile(t, filepath.Join(root, "docs", "overview.md"), "# Overview\n")
	writeFile(t, filepath.Join(root, "docs", "specs", "reality", "ARTIFACT-CONTRACT.md"), "# contract\nSee `internal/pkg/logic.go`.\n")
}

func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
