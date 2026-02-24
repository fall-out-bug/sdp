package orchestrate_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/orchestrate"
)

func TestErrNoPR(t *testing.T) {
	if orchestrate.ErrNoPR == nil {
		t.Fatal("ErrNoPR must be non-nil")
	}
	if !errors.Is(orchestrate.ErrNoPR, orchestrate.ErrNoPR) {
		t.Error("errors.Is(err, ErrNoPR) should be true for ErrNoPR")
	}
	if orchestrate.ErrNoPR.Error() != "no PR found for current branch" {
		t.Errorf("ErrNoPR message: got %q", orchestrate.ErrNoPR.Error())
	}
}

func TestEnsureRunFile(t *testing.T) {
	dir := t.TempDir()
	if err := orchestrate.EnsureRunFile(dir, "F016", "feature/F016-oneshot"); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 run file, got %d", len(entries))
	}
	name := filepath.Base(entries[0].Name())
	if len(name) < 10 || name[:10] != "oneshot-F0" {
		t.Errorf("unexpected run file name: %s", name)
	}
	data, err := os.ReadFile(filepath.Join(dir, entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}
	var rf struct {
		RunID     string `json:"run_id"`
		FeatureID string `json:"feature_id"`
		Branch    string `json:"branch"`
	}
	if err := json.Unmarshal(data, &rf); err != nil {
		t.Fatal(err)
	}
	if rf.FeatureID != "F016" || rf.Branch != "feature/F016-oneshot" {
		t.Errorf("run file content mismatch: %+v", rf)
	}
}

func TestEnsureRunFileInvalidFeatureID(t *testing.T) {
	dir := t.TempDir()
	err := orchestrate.EnsureRunFile(dir, "", "branch")
	if err == nil {
		t.Fatal("expected error for empty featureID")
	}
	err = orchestrate.EnsureRunFile(dir, "F016/../x", "branch")
	if err == nil {
		t.Fatal("expected error for path-traversal featureID")
	}
}

func TestEnsureRunFileMkdirFails(t *testing.T) {
	// Use a path that would fail MkdirAll (e.g. parent is a file)
	dir := t.TempDir()
	filePath := filepath.Join(dir, "blocker")
	if err := os.WriteFile(filePath, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	badDir := filepath.Join(filePath, "runs")
	err := orchestrate.EnsureRunFile(badDir, "F016", "branch")
	if err == nil {
		t.Fatal("expected error when parent is file")
	}
}
