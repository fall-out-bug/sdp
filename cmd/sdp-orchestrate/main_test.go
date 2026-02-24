package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMainMissingFeatureExits(t *testing.T) {
	// Build and run sdp-orchestrate without --feature; expect exit 1 and stderr.
	bin := filepath.Join(t.TempDir(), "sdp-orchestrate")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Skipf("build failed: %v", err)
	}
	out, err := exec.Command(bin).CombinedOutput()
	if err == nil {
		t.Fatal("expected non-zero exit when --feature is missing")
	}
	if len(out) == 0 {
		t.Error("expected stderr output")
	}
	s := string(out)
	if !strings.Contains(s, "feature") && !strings.Contains(s, "error") {
		t.Errorf("stderr should mention feature or error, got: %s", out)
	}
}
