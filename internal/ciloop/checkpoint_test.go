package ciloop_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/ciloop"
)

func TestLoadCheckpoint(t *testing.T) {
	dir := t.TempDir()
	content := `{
		"schema": "1.0",
		"feature_id": "F014",
		"branch": "feature/F014-ci-loop-cli",
		"pr_number": 42,
		"pr_url": "https://github.com/org/repo/pull/42",
		"phase": "build"
	}`
	if err := os.WriteFile(filepath.Join(dir, "F014.json"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cp, err := ciloop.LoadCheckpoint(dir, "F014")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cp.FeatureID != "F014" {
		t.Errorf("expected feature_id F014, got %q", cp.FeatureID)
	}
	if cp.PRNumber == nil || *cp.PRNumber != 42 {
		t.Errorf("expected pr_number 42, got %v", cp.PRNumber)
	}
	if cp.Branch != "feature/F014-ci-loop-cli" {
		t.Errorf("expected branch feature/F014-ci-loop-cli, got %q", cp.Branch)
	}
}

func TestLoadCheckpointNotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := ciloop.LoadCheckpoint(dir, "F999")
	if err == nil {
		t.Fatal("expected error for missing checkpoint, got nil")
	}
}

func TestLoadCheckpointPathTraversalRejected(t *testing.T) {
	dir := t.TempDir()
	_, err := ciloop.LoadCheckpoint(dir, "../../../etc/passwd")
	if err == nil {
		t.Fatal("expected error for path traversal featureID, got nil")
	}
}

func TestSaveCheckpointPathTraversalRejected(t *testing.T) {
	dir := t.TempDir()
	cp := &ciloop.Checkpoint{FeatureID: "../../../etc/passwd"}
	err := ciloop.SaveCheckpoint(dir, cp)
	if err == nil {
		t.Fatal("expected error for path traversal featureID in save, got nil")
	}
}

func TestLoadCheckpointInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "F014.json"), []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := ciloop.LoadCheckpoint(dir, "F014")
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestSaveCheckpoint(t *testing.T) {
	dir := t.TempDir()
	prNum := 42
	cp := &ciloop.Checkpoint{
		Schema:    "1.0",
		FeatureID: "F014",
		Branch:    "feature/F014-ci-loop-cli",
		PRNumber:  &prNum,
		PRURL:     "https://github.com/org/repo/pull/42",
		Phase:     "build",
	}
	if err := ciloop.SaveCheckpoint(dir, cp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Read back and verify.
	loaded, err := ciloop.LoadCheckpoint(dir, "F014")
	if err != nil {
		t.Fatalf("load after save: %v", err)
	}
	if loaded.Phase != "build" {
		t.Errorf("expected phase=build (saved as given), got %q", loaded.Phase)
	}
	if loaded.UpdatedAt == "" {
		t.Error("expected updated_at to be set")
	}
}
