package orchestrate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestComputePromptHash(t *testing.T) {
	// Empty string has known SHA-256
	got := ComputePromptHash("")
	if len(got) != 64 {
		t.Errorf("hash length = %d, want 64", len(got))
	}
	// Deterministic
	if got != ComputePromptHash("") {
		t.Error("hash should be deterministic")
	}
}

func TestBuildContextSources(t *testing.T) {
	dir := t.TempDir()
	// Create minimal files
	wsDir := filepath.Join(dir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	wsPath := filepath.Join(wsDir, "00-026-01.md")
	if err := os.WriteFile(wsPath, []byte("# test"), 0o644); err != nil {
		t.Fatal(err)
	}
	sdpDir := filepath.Join(dir, ".sdp", "checkpoints")
	if err := os.MkdirAll(sdpDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cpPath := filepath.Join(sdpDir, "F026.json")
	if err := os.WriteFile(cpPath, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	sources := BuildContextSources(dir, "F026", "00-026-01", nil)
	if len(sources) == 0 {
		t.Error("expected at least workstream_spec and checkpoint")
	}
	for _, s := range sources {
		if s.Type == "" || s.Path == "" || s.Hash == "" {
			t.Errorf("invalid source: %+v", s)
		}
		if len(s.Hash) != 64 {
			t.Errorf("hash length = %d for %s", len(s.Hash), s.Type)
		}
	}
}

func TestWritePromptProvenance(t *testing.T) {
	dir := t.TempDir()
	sources := []ContextSource{
		{Type: "workstream_spec", Path: "docs/ws.md", Hash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
	}
	if err := WritePromptProvenance(dir, "abc123", sources); err != nil {
		t.Fatalf("WritePromptProvenance: %v", err)
	}
	path := filepath.Join(dir, ".sdp", "prompt-provenance.json")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(b) == 0 {
		t.Error("expected non-empty file")
	}
}
