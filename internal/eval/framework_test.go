package eval

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCases_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	cases, err := LoadCases(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(cases) != 0 {
		t.Errorf("expected 0 cases, got %d", len(cases))
	}
}

func TestLoadCases_MalformedYAML(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(f, []byte("not: valid: yaml: here"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadCases(dir, "")
	if err == nil {
		t.Fatal("expected error for malformed YAML")
	}
}

func TestExtractAgentOutput(t *testing.T) {
	// Simple format: role + content
	data := []byte(`{"role":"user","content":"hello"}
{"role":"assistant","content":"agent says hi"}`)
	out := extractAgentOutput(data)
	if out != "agent says hi\n" {
		t.Errorf("got %q", out)
	}
}

func TestRunCase_KnownBad(t *testing.T) {
	tmp := t.TempDir()
	// Transcript with forbidden patterns; verdict FAIL = we expect to catch it
	os.WriteFile(filepath.Join(tmp, "bad.jsonl"), []byte(`{"role":"assistant","content":"Next steps: 1. approve and merge"}`), 0o644)
	c := &Case{
		Name:              "bad",
		InputTranscript:   "bad.jsonl",
		ForbiddenPatterns: []string{"Next steps", "approve and merge"},
		RequiredPatterns:  []string{},
		Verdict:           "FAIL",
	}
	r := RunCase(c, tmp)
	if !r.Pass {
		t.Error("expected PASS for known-bad transcript (correctly flagged)")
	}
}

func TestRunCase_KnownGood(t *testing.T) {
	tmp := t.TempDir()
	os.WriteFile(filepath.Join(tmp, "good.jsonl"), []byte(`{"role":"assistant","content":"CI GREEN - @oneshot complete"}`), 0o644)
	c := &Case{
		Name:              "good",
		InputTranscript:   "good.jsonl",
		ForbiddenPatterns: []string{"Next steps"},
		RequiredPatterns:  []string{"CI GREEN"},
		Verdict:           "PASS",
	}
	r := RunCase(c, tmp)
	if !r.Pass {
		t.Errorf("expected PASS for known-good transcript: %s", r.Reason)
	}
}

func TestRun_OneshotEvals(t *testing.T) {
	// Run from project root so testdata paths resolve
	root, _ := os.Getwd()
	for _, d := range []string{"internal/eval", "eval"} {
		if _, err := os.Stat(filepath.Join(root, d)); err == nil {
			root = filepath.Dir(root)
			break
		}
	}
	// Find project root (has testdata/eval)
	for {
		if _, err := os.Stat(filepath.Join(root, "testdata", "eval")); err == nil {
			break
		}
		parent := filepath.Dir(root)
		if parent == root {
			t.Skip("project root not found")
		}
		root = parent
	}
	casesDir := filepath.Join(root, "internal", "eval", "cases")
	results, err := Run(root, casesDir, "oneshot")
	if err != nil {
		t.Fatal(err)
	}
	passed := 0
	for _, r := range results {
		if r.Pass {
			passed++
		}
	}
	// We expect: 5 cases, all pass (3 verdict FAIL correctly flag bad transcripts, 2 verdict PASS)
	if len(results) != 5 {
		t.Errorf("expected 5 cases, got %d", len(results))
	}
	if passed != 5 {
		t.Errorf("expected all 5 to pass, got %d", passed)
	}
}
