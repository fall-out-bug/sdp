package orchestrate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/orchestrate"
)

func TestLoadCheckpoint(t *testing.T) {
	dir := t.TempDir()
	cp := &orchestrate.Checkpoint{
		Schema:    "orchestrate.v1",
		FeatureID: "F016",
		Branch:    "feature/F016-oneshot",
		Phase:     orchestrate.PhaseBuild,
	}
	if err := orchestrate.SaveCheckpoint(dir, cp); err != nil {
		t.Fatal(err)
	}
	loaded, err := orchestrate.LoadCheckpoint(dir, "F016")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.FeatureID != "F016" || loaded.Phase != orchestrate.PhaseBuild {
		t.Errorf("loaded checkpoint mismatch: %+v", loaded)
	}
}

func TestLoadCheckpointNotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := orchestrate.LoadCheckpoint(dir, "F999")
	if err == nil {
		t.Fatal("expected error for missing checkpoint")
	}
}

func TestLoadCheckpointInvalidFeatureID(t *testing.T) {
	_, err := orchestrate.LoadCheckpoint("/tmp", "F016/../")
	if err == nil {
		t.Fatal("expected error for invalid feature_id")
	}
}

func TestSaveCheckpointInvalidFeatureID(t *testing.T) {
	cp := &orchestrate.Checkpoint{FeatureID: "F016/../x"}
	err := orchestrate.SaveCheckpoint(t.TempDir(), cp)
	if err == nil {
		t.Fatal("expected error for invalid feature_id")
	}
}

func TestSaveCheckpointInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "F016.json")
	if err := os.WriteFile(path, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := orchestrate.LoadCheckpoint(dir, "F016")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
