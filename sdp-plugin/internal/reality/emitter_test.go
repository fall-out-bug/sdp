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

func seedProject(t *testing.T, root string) {
	t.Helper()

	writeFile(t, filepath.Join(root, "cmd", "app", "main.go"), "package main\n\nfunc main() {}\n")
	writeFile(t, filepath.Join(root, "internal", "pkg", "logic.go"), "package pkg\n\nfunc Sum(a, b int) int { return a + b }\n")
	writeFile(t, filepath.Join(root, "internal", "pkg", "logic_test.go"), "package pkg\n\nimport \"testing\"\n\nfunc TestSum(t *testing.T) {\n\tif Sum(2, 2) != 4 {\n\t\tt.Fatal(\"unexpected\")\n\t}\n}\n")
	writeFile(t, filepath.Join(root, "docs", "overview.md"), "# Overview\n")
	writeFile(t, filepath.Join(root, "docs", "specs", "reality", "ARTIFACT-CONTRACT.md"), "# contract\n")
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
