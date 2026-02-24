package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMainMissingSkillExits(t *testing.T) {
	modRoot, _ := os.Getwd()
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
	bin := filepath.Join(t.TempDir(), "sdp-eval")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/sdp-eval")
	cmd.Dir = modRoot
	if err := cmd.Run(); err != nil {
		t.Skipf("build failed: %v", err)
	}
	out, err := exec.Command(bin).CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit when --skill and --all are missing")
	}
	s := string(out)
	if !strings.Contains(s, "skill") && !strings.Contains(s, "error") {
		t.Errorf("stderr should mention skill or error, got: %s", out)
	}
}
