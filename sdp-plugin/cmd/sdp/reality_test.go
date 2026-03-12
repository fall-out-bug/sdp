package main

import (
	"encoding/json"
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
	foundValidate := false
	for _, sub := range cmd.Commands() {
		if sub.Name() == "emit-oss" {
			foundEmitOSS = true
		}
		if sub.Name() == "validate" {
			foundValidate = true
		}
	}
	if !foundEmitOSS {
		t.Fatal("realityCmd() missing emit-oss subcommand")
	}
	if !foundValidate {
		t.Fatal("realityCmd() missing validate subcommand")
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

func TestRealityEmitOSSCmd_WithQuickFocus(t *testing.T) {
	tmpDir := t.TempDir()
	seedRealityProject(t, tmpDir)
	writeRealityFile(t, filepath.Join(tmpDir, "configs", "app.yaml"), "database: postgres\ncache: redis\n")

	cmd := realityCmd()
	cmd.SetArgs([]string{"emit-oss", "--root", tmpDir, "--quick", "--focus", "docs"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("reality emit-oss quick/docs failed: %v", err)
	}

	summaryPath := filepath.Join(tmpDir, ".sdp", "reality", "reality-summary.json")
	data, err := os.ReadFile(summaryPath)
	if err != nil {
		t.Fatalf("read summary: %v", err)
	}

	var summary map[string]any
	if err := json.Unmarshal(data, &summary); err != nil {
		t.Fatalf("parse summary: %v", err)
	}

	scope, ok := summary["scope"].(map[string]any)
	if !ok {
		t.Fatalf("summary scope missing or invalid: %#v", summary["scope"])
	}
	if scope["mode"] != "quick" {
		t.Fatalf("unexpected mode: %v", scope["mode"])
	}
	if scope["focus"] != "docs" {
		t.Fatalf("unexpected focus: %v", scope["focus"])
	}
}

func TestRealityEmitOSSCmd_RejectsConflictingModes(t *testing.T) {
	tmpDir := t.TempDir()
	seedRealityProject(t, tmpDir)

	cmd := realityCmd()
	cmd.SetArgs([]string{"emit-oss", "--root", tmpDir, "--quick", "--deep"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("expected conflicting modes to fail")
	}
}

func TestRealityValidateCmd_WithRoot(t *testing.T) {
	tmpDir := t.TempDir()
	seedRealityProject(t, tmpDir)

	emit := realityCmd()
	emit.SetArgs([]string{"emit-oss", "--root", tmpDir})
	if err := emit.Execute(); err != nil {
		t.Fatalf("reality emit-oss failed: %v", err)
	}

	validate := realityCmd()
	validate.SetArgs([]string{"validate", "--root", tmpDir})
	if err := validate.Execute(); err != nil {
		t.Fatalf("reality validate failed: %v", err)
	}
}

func TestRealityValidateCmd_FailsOnCorruptArtifact(t *testing.T) {
	tmpDir := t.TempDir()
	seedRealityProject(t, tmpDir)

	emit := realityCmd()
	emit.SetArgs([]string{"emit-oss", "--root", tmpDir})
	if err := emit.Execute(); err != nil {
		t.Fatalf("reality emit-oss failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, ".sdp", "reality", "quality-report.json"), []byte("{bad json"), 0o644); err != nil {
		t.Fatalf("corrupt quality report: %v", err)
	}

	validate := realityCmd()
	validate.SetArgs([]string{"validate", "--root", tmpDir})
	if err := validate.Execute(); err == nil {
		t.Fatal("expected reality validate to fail on corrupt artifact")
	}
}

func seedRealityProject(t *testing.T, root string) {
	t.Helper()

	writeRealityFile(t, filepath.Join(root, "cmd", "demo", "main.go"), "package main\n\nfunc main() {}\n")
	writeRealityFile(t, filepath.Join(root, "internal", "module", "service.go"), "package module\n\nfunc Enabled() bool { return true }\n")
	writeRealityFile(t, filepath.Join(root, "internal", "module", "service_test.go"), "package module\n\nimport \"testing\"\n\nfunc TestEnabled(t *testing.T) {\n\tif !Enabled() {\n\t\tt.Fatal(\"expected true\")\n\t}\n}\n")
	writeRealityFile(t, filepath.Join(root, "docs", "specs", "reality", "ARTIFACT-CONTRACT.md"), "# contract\nSee `internal/module/service.go`.\n")
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
