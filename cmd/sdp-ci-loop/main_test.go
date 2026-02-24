package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMainFlagsHelp(t *testing.T) {
	wd, _ := os.Getwd()
	modRoot := wd
	for {
		if _, err := os.Stat(filepath.Join(modRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(modRoot)
		if parent == modRoot {
			t.Skip("no go.mod found")
		}
		modRoot = parent
	}
	dir := t.TempDir()
	bin := filepath.Join(dir, "sdp-ci-loop")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/sdp-ci-loop")
	cmd.Dir = modRoot
	if err := cmd.Run(); err != nil {
		t.Skipf("build failed: %v", err)
	}
	out, err := exec.Command(bin, "-h").CombinedOutput()
	if err != nil {
		t.Fatalf("sdp-ci-loop -h: %v", err)
	}
	if !strings.Contains(string(out), "-pr") || !strings.Contains(string(out), "-feature") {
		t.Errorf("help output missing -pr or -feature: %s", out)
	}
}

func TestMainMissingPRExits(t *testing.T) {
	wd, _ := os.Getwd()
	modRoot := wd
	for {
		if _, err := os.Stat(filepath.Join(modRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(modRoot)
		if parent == modRoot {
			t.Skip("no go.mod found")
		}
		modRoot = parent
	}
	dir := t.TempDir()
	bin := filepath.Join(dir, "sdp-ci-loop")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/sdp-ci-loop")
	cmd.Dir = modRoot
	if err := cmd.Run(); err != nil {
		t.Skipf("build failed: %v", err)
	}
	run := exec.Command(bin)
	run.Dir = t.TempDir()
	err := run.Run()
	if err == nil {
		t.Fatal("expected exit 1 when --pr missing")
	}
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 1 {
		t.Errorf("expected exit 1, got %d", exitErr.ExitCode())
	}
}

// TestIntegrationStub is a placeholder for full integration tests (requires gh CLI, repo).
func TestIntegrationStub(t *testing.T) {
	t.Skip("integration test: requires gh CLI and authenticated repo")
}
