package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRealityCmd(t *testing.T) {
	cmd := realityCmd()

	if cmd.Use != "reality" {
		t.Fatalf("realityCmd() has wrong use: %s", cmd.Use)
	}

	foundEmitOSS := false
	for _, sub := range cmd.Commands() {
		if sub.Name() == "emit-oss" {
			foundEmitOSS = true
			break
		}
	}
	if !foundEmitOSS {
		t.Fatal("realityCmd() missing emit-oss subcommand")
	}
}

func TestRealityEmitOSSCmd_WithRoot(t *testing.T) {
	tmpDir := t.TempDir()
	seedRealityProject(t, tmpDir)

	cmd := realityCmd()
	cmd.SetArgs([]string{"emit-oss", "--root", tmpDir})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("reality emit-oss failed: %v", err)
	}

	expected := []string{
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
	for _, rel := range expected {
		abs := filepath.Join(tmpDir, rel)
		if _, err := os.Stat(abs); err != nil {
			t.Fatalf("expected artifact missing: %s (%v)", rel, err)
		}
	}
}

func seedRealityProject(t *testing.T, root string) {
	t.Helper()

	writeRealityFile(t, filepath.Join(root, "cmd", "demo", "main.go"), "package main\n\nfunc main() {}\n")
	writeRealityFile(t, filepath.Join(root, "internal", "module", "service.go"), "package module\n\nfunc Enabled() bool { return true }\n")
	writeRealityFile(t, filepath.Join(root, "internal", "module", "service_test.go"), "package module\n\nimport \"testing\"\n\nfunc TestEnabled(t *testing.T) {\n\tif !Enabled() {\n\t\tt.Fatal(\"expected true\")\n\t}\n}\n")
	writeRealityFile(t, filepath.Join(root, "docs", "specs", "reality", "ARTIFACT-CONTRACT.md"), "# contract\n")
}

func writeRealityFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
