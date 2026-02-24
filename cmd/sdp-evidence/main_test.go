package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateValid(t *testing.T) {
	// Build and run: sdp-evidence validate --evidence specs/strict-evidence-template.json --require-pr-url=false
	bin := filepath.Join(t.TempDir(), "sdp-evidence")
	if err := exec.Command("go", "build", "-o", bin, ".").Run(); err != nil {
		t.Fatalf("build: %v", err)
	}
	wd, _ := os.Getwd()
	root := filepath.Dir(filepath.Dir(wd))
	cmd := exec.Command(bin, "validate", "--evidence", "specs/strict-evidence-template.json", "--require-pr-url=false")
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
	wd, _ := os.Getwd()
	root := filepath.Dir(filepath.Dir(wd))
	cmd := exec.Command(bin, "inspect", "--evidence", "specs/strict-evidence-template.json", "--require-pr-url=false")
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
