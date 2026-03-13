package orchestrate

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunBuildPhase_WithFakeInvoker(t *testing.T) {
	dir := t.TempDir()
	sdpDir := filepath.Join(dir, ".sdp", "checkpoints")
	if err := os.MkdirAll(sdpDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cpPath := filepath.Join(sdpDir, "F053.json")
	if err := os.WriteFile(cpPath, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	wsDir := filepath.Join(dir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(wsDir, "00-053-34.md"), []byte("# test"), 0o644); err != nil {
		t.Fatal(err)
	}

	fake := &fakeLLMInvoker{
		output:   "abc123def4567890123456789012345678901234",
		exitCode: 0,
	}
	commit, err := RunBuildPhase(context.Background(), dir, "F053", "00-053-34", fake)
	if err != nil {
		t.Fatalf("RunBuildPhase: %v", err)
	}
	if commit != "abc123def4567890123456789012345678901234" {
		t.Errorf("commit = %q, want %q", commit, "abc123def4567890123456789012345678901234")
	}
	if !fake.invoked {
		t.Error("fake invoker was not invoked")
	}
}

func TestRunBuildPhase_WithFakeInvoker_NonZeroExit(t *testing.T) {
	dir := t.TempDir()
	sdpDir := filepath.Join(dir, ".sdp", "checkpoints")
	if err := os.MkdirAll(sdpDir, 0o755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(sdpDir, "F053.json"), []byte("{}"), 0o644)
	wsDir := filepath.Join(dir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(filepath.Join(wsDir, "00-053-34.md"), []byte("# test"), 0o644)

	fake := &fakeLLMInvoker{output: "build failed", exitCode: 1}
	_, err := RunBuildPhase(context.Background(), dir, "F053", "00-053-34", fake)
	if err == nil {
		t.Fatal("expected error from non-zero exit")
	}
}

type fakeLLMInvoker struct {
	output   string
	exitCode int
	invoked  bool
	agent    string
	prompt   string
}

func (f *fakeLLMInvoker) Invoke(ctx context.Context, dir, agent, prompt string) (string, int, error) {
	f.invoked = true
	f.agent = agent
	f.prompt = prompt
	return f.output, f.exitCode, nil
}

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

func TestRunReviewPhase_UsesReadOnlyPrompt(t *testing.T) {
	fake := &fakeLLMInvoker{
		output:   "APPROVED",
		exitCode: 0,
	}
	approved, err := RunReviewPhase(context.Background(), t.TempDir(), "F053", fake)
	if err != nil {
		t.Fatalf("RunReviewPhase: %v", err)
	}
	if !approved {
		t.Fatal("expected approved review result")
	}
	if fake.agent != "reviewer" {
		t.Fatalf("agent = %q, want reviewer", fake.agent)
	}
	if fake.prompt == "" {
		t.Fatal("expected prompt to be captured")
	}
	if prompt := fake.prompt; strings.Contains(prompt, "Fix P0/P1 findings") {
		t.Fatalf("prompt should not ask reviewer to fix findings: %q", prompt)
	}
	if !strings.Contains(fake.prompt, "Do not modify code, docs, prompts, or workstream files") {
		t.Fatalf("prompt should enforce read-only review: %q", fake.prompt)
	}
}

func TestRunReviewPhase_ChangesRequestedIsNotApproved(t *testing.T) {
	fake := &fakeLLMInvoker{
		output:   "CHANGES_REQUESTED\nmissing test coverage",
		exitCode: 0,
	}
	approved, err := RunReviewPhase(context.Background(), t.TempDir(), "F053", fake)
	if err != nil {
		t.Fatalf("RunReviewPhase: %v", err)
	}
	if approved {
		t.Fatal("expected CHANGES_REQUESTED output to be treated as not approved")
	}
}
