package ciloop_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/ciloop"
)

func TestMatchingFixersImportError(t *testing.T) {
	dir := t.TempDir()
	reg := ciloop.NewAutofixerRegistry(dir)
	log := "internal/foo/bar.go:5:2: imported and not used: \"fmt\""
	matching := reg.MatchingFixers(log)
	if len(matching) == 0 {
		t.Fatal("expected matching fixers for import error, got none")
	}
	names := make([]string, len(matching))
	for i, f := range matching {
		names[i] = f.Name
	}
	if !contains(names, "goimports") {
		t.Errorf("expected goimports to match, got %v", names)
	}
}

func TestMatchingFixersGoModTidy(t *testing.T) {
	dir := t.TempDir()
	reg := ciloop.NewAutofixerRegistry(dir)
	log := "cannot find package \"github.com/example/missing\""
	matching := reg.MatchingFixers(log)
	names := make([]string, len(matching))
	for i, f := range matching {
		names[i] = f.Name
	}
	if !contains(names, "go-mod-tidy") {
		t.Errorf("expected go-mod-tidy to match missing package, got %v", names)
	}
}

func TestMatchingFixersNoMatch(t *testing.T) {
	dir := t.TempDir()
	reg := ciloop.NewAutofixerRegistry(dir)
	log := "secrets detected in file xyz"
	matching := reg.MatchingFixers(log)
	if len(matching) != 0 {
		t.Errorf("expected no matching fixers for secrets log, got %v", matching)
	}
}

func TestDeterministicFirstFixerFallsThroughToInnerWhenNoDeterministicHelp(t *testing.T) {
	dir := t.TempDir()
	reg := ciloop.NewAutofixerRegistry(dir)
	committer := &fakeCommitter{}
	fetcher := &fakeLogFetcher{logs: map[string]string{"run1": goTestFailureLog}}
	inner := ciloop.NewFixer(ciloop.FixerOptions{
		PRNumber:       42,
		FeatureID:      "F027",
		DiagnosticsDir: filepath.Join(dir, ".sdp", "ci-fixes"),
		Ctx:            context.Background(),
		Committer:      committer,
		LogFetcher:     fetcher,
		DecisionLogger: func(_, _ string) error { return nil },
	})
	wrapper := &ciloop.DeterministicFirstFixer{
		ProjectRoot: dir,
		Registry:    reg,
		Runner:      &autofixerRunner{},
		Committer:   &fakeCommitter{}, // separate committer for deterministic path
		LogFetcher:  fetcher,
		Inner:       inner,
		PRNumber:    42,
	}
	// Log matches goimports but we use a dir with no .go files - deterministic won't change anything.
	// The inner fixer will run (go-test pattern matches) and commit diagnostics.
	checks := []ciloop.CheckResult{{Name: "go-test", State: ciloop.StateFailure}}
	err := wrapper.Fix(checks)
	if err != nil {
		t.Fatalf("Fix: %v", err)
	}
	// Inner fixer should have committed (diagnostics file)
	if len(committer.commits) != 1 {
		t.Errorf("expected inner fixer to commit, got %d commits", len(committer.commits))
	}
}

func contains(s []string, x string) bool {
	for _, v := range s {
		if v == x {
			return true
		}
	}
	return false
}

func TestSplitCommand(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{"goimports -w .", []string{"goimports", "-w", "."}},
		{"go mod tidy", []string{"go", "mod", "tidy"}},
		{"go fmt ./...", []string{"go", "fmt", "./..."}},
		{"single", []string{"single"}},
	}
	for _, tt := range tests {
		got := ciloop.SplitCommand(tt.in)
		if len(got) != len(tt.want) {
			t.Errorf("splitCommand(%q): got %v, want %v", tt.in, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("splitCommand(%q)[%d]: got %q, want %q", tt.in, i, got[i], tt.want[i])
			}
		}
	}
}

type autofixerRunner struct{}

func (f *autofixerRunner) Run(_ string, _ ...string) ([]byte, error) {
	return nil, nil
}

func TestParseAutoFixersYAML(t *testing.T) {
	valid := `
fixers:
  - name: custom
    command: "echo hello"
    applies_to: "some pattern"
    timeout: 10
`
	fixers, err := ciloop.ParseAutoFixersYAML([]byte(valid))
	if err != nil {
		t.Fatalf("parse valid YAML: %v", err)
	}
	if len(fixers) != 1 {
		t.Fatalf("expected 1 fixer, got %d", len(fixers))
	}
	if fixers[0].Name != "custom" || fixers[0].Command != "echo hello" || fixers[0].Timeout != 10 {
		t.Errorf("got %+v", fixers[0])
	}

	invalid := "not: valid: yaml"
	_, err = ciloop.ParseAutoFixersYAML([]byte(invalid))
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}
