package ciloop_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall-out-bug/sdp/internal/ciloop"
)

// fakeCommitter records calls to commit+push.
type fakeCommitter struct {
	commits []string
	pushes  []string
	err     error
}

func (f *fakeCommitter) Commit(ctx context.Context, msg string) error {
	if f.err != nil {
		return f.err
	}
	f.commits = append(f.commits, msg)
	return nil
}

func (f *fakeCommitter) Push(ctx context.Context) error {
	if f.err != nil {
		return f.err
	}
	f.pushes = append(f.pushes, "push")
	return nil
}

// fakeLogFetcher returns pre-set failure logs per run ID.
type fakeLogFetcher struct {
	logs map[string]string
	err  error
}

func (f *fakeLogFetcher) FailedLogs(prNumber int) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	if f.logs != nil {
		for _, v := range f.logs {
			return v, nil
		}
	}
	return "", nil
}

const goTestFailureLog = `
--- FAIL: TestFoo (0.00s)
    foo_test.go:12: assertion failed
FAIL	sdp_dev/internal/foo	1.234s
`

const goBuildFailureLog = `
./internal/bar/bar.go:42:5: undefined: SomeFunc
`

const goBuildNoPkgLog = `
./cmd/foo/main.go:5:2: cannot find package "github.com/example/missing"
`

const k8sFailureLog = `
Error: yaml: line 5: did not find expected key
`

func TestDiagnosticsFileNoRawLog(t *testing.T) {
	// Security: diagnostics file must not contain raw CI log (secrets, tokens).
	dir := t.TempDir()
	committer := &fakeCommitter{}
	fetcher := &fakeLogFetcher{logs: map[string]string{"run1": goTestFailureLog}}
	fixer := ciloop.NewFixer(ciloop.FixerOptions{
		PRNumber:       42,
		FeatureID:      "F014",
		DiagnosticsDir: dir,
		Committer:      committer,
		LogFetcher:     fetcher,
		DecisionLogger: func(decision, rationale string) error { return nil },
	})
	checks := []ciloop.CheckResult{{Name: "go-test", State: ciloop.StateFailure}}
	if err := fixer.Fix(checks); err != nil {
		t.Fatalf("Fix: %v", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 diagnostics file, got %d", len(entries))
	}
	data, err := os.ReadFile(filepath.Join(dir, entries[0].Name()))
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	content := string(data)
	// Raw log contains "assertion failed", "FAIL", "foo_test.go" — must not appear.
	for _, forbidden := range []string{"assertion failed", "foo_test.go", "FAIL\t"} {
		if strings.Contains(content, forbidden) {
			t.Errorf("diagnostics file must not contain raw log; found %q", forbidden)
		}
	}
	// Must contain sanitized fix type.
	if !strings.Contains(content, "go-test") {
		t.Errorf("diagnostics file should contain fix type go-test")
	}
}

func TestFixerGoTestFailure(t *testing.T) {
	committer := &fakeCommitter{}
	fetcher := &fakeLogFetcher{logs: map[string]string{"run1": goTestFailureLog}}
	fixer := ciloop.NewFixer(ciloop.FixerOptions{
		PRNumber:       42,
		FeatureID:      "F014",
		DiagnosticsDir: t.TempDir(),
		Committer:      committer,
		LogFetcher:     fetcher,
		DecisionLogger: func(decision, rationale string) error { return nil },
	})
	checks := []ciloop.CheckResult{{Name: "go-test", State: ciloop.StateFailure}}
	err := fixer.Fix(checks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(committer.commits) != 1 {
		t.Errorf("expected 1 commit, got %d", len(committer.commits))
	}
	if len(committer.pushes) != 1 {
		t.Errorf("expected 1 push, got %d", len(committer.pushes))
	}
}

func TestFixerGoBuildFailure(t *testing.T) {
	committer := &fakeCommitter{}
	fetcher := &fakeLogFetcher{logs: map[string]string{"run1": goBuildFailureLog}}
	fixer := ciloop.NewFixer(ciloop.FixerOptions{
		PRNumber:       42,
		FeatureID:      "F014",
		DiagnosticsDir: t.TempDir(),
		Committer:      committer,
		LogFetcher:     fetcher,
		DecisionLogger: func(decision, rationale string) error { return nil },
	})
	checks := []ciloop.CheckResult{{Name: "go-build", State: ciloop.StateFailure}}
	err := fixer.Fix(checks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(committer.commits) != 1 {
		t.Errorf("expected 1 commit, got %d", len(committer.commits))
	}
}

func TestFixerGoBuildNoPkgFailure(t *testing.T) {
	committer := &fakeCommitter{}
	fetcher := &fakeLogFetcher{logs: map[string]string{"run1": goBuildNoPkgLog}}
	fixer := ciloop.NewFixer(ciloop.FixerOptions{
		PRNumber:       42,
		FeatureID:      "F014",
		DiagnosticsDir: t.TempDir(),
		Committer:      committer,
		LogFetcher:     fetcher,
		DecisionLogger: func(decision, rationale string) error { return nil },
	})
	checks := []ciloop.CheckResult{{Name: "go-build", State: ciloop.StateFailure}}
	err := fixer.Fix(checks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(committer.commits) != 1 {
		t.Errorf("expected 1 commit, got %d", len(committer.commits))
	}
}

func TestFixerK8sValidateFailure(t *testing.T) {
	committer := &fakeCommitter{}
	fetcher := &fakeLogFetcher{logs: map[string]string{"run1": k8sFailureLog}}
	fixer := ciloop.NewFixer(ciloop.FixerOptions{
		PRNumber:       42,
		FeatureID:      "F014",
		DiagnosticsDir: t.TempDir(),
		Committer:      committer,
		LogFetcher:     fetcher,
		DecisionLogger: func(decision, rationale string) error { return nil },
	})
	checks := []ciloop.CheckResult{{Name: "k8s-validate", State: ciloop.StateFailure}}
	err := fixer.Fix(checks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(committer.commits) != 1 {
		t.Errorf("expected 1 commit, got %d", len(committer.commits))
	}
}

func TestFixerUnparsableLogEscalates(t *testing.T) {
	committer := &fakeCommitter{}
	fetcher := &fakeLogFetcher{logs: map[string]string{"run1": "some unparseable noise with no pattern"}}
	fixer := ciloop.NewFixer(ciloop.FixerOptions{
		PRNumber:       42,
		FeatureID:      "F014",
		DiagnosticsDir: t.TempDir(),
		Committer:      committer,
		LogFetcher:     fetcher,
		DecisionLogger: func(decision, rationale string) error { return nil },
	})
	checks := []ciloop.CheckResult{{Name: "go-test", State: ciloop.StateFailure}}
	err := fixer.Fix(checks)
	if err == nil {
		t.Fatal("expected error for unparseable log, got nil")
	}
}

func TestFixerLogFetchError(t *testing.T) {
	committer := &fakeCommitter{}
	fetcher := &fakeLogFetcher{err: errors.New("gh: auth error")}
	fixer := ciloop.NewFixer(ciloop.FixerOptions{
		PRNumber:       42,
		FeatureID:      "F014",
		DiagnosticsDir: t.TempDir(),
		Committer:      committer,
		LogFetcher:     fetcher,
		DecisionLogger: func(decision, rationale string) error { return nil },
	})
	checks := []ciloop.CheckResult{{Name: "go-test", State: ciloop.StateFailure}}
	err := fixer.Fix(checks)
	if err == nil {
		t.Fatal("expected error from log fetch failure, got nil")
	}
}

func TestFixerCommitMessageContainsFixCi(t *testing.T) {
	committer := &fakeCommitter{}
	fetcher := &fakeLogFetcher{logs: map[string]string{"run1": goTestFailureLog}}
	fixer := ciloop.NewFixer(ciloop.FixerOptions{
		PRNumber:       42,
		FeatureID:      "F014",
		DiagnosticsDir: t.TempDir(),
		Committer:      committer,
		LogFetcher:     fetcher,
		DecisionLogger: func(decision, rationale string) error { return nil },
	})
	checks := []ciloop.CheckResult{{Name: "go-test", State: ciloop.StateFailure}}
	fixer.Fix(checks)
	if len(committer.commits) == 0 {
		t.Fatal("no commit made")
	}
	msg := committer.commits[0]
	if len(msg) < 5 || msg[:4] != "fix(" {
		t.Errorf("commit message should start with fix(...), got: %q", msg)
	}
}

func TestFixerLogsDecision(t *testing.T) {
	committer := &fakeCommitter{}
	fetcher := &fakeLogFetcher{logs: map[string]string{"run1": goTestFailureLog}}
	logged := false
	fixer := ciloop.NewFixer(ciloop.FixerOptions{
		PRNumber:       42,
		FeatureID:      "F014",
		DiagnosticsDir: t.TempDir(),
		Committer:      committer,
		LogFetcher:     fetcher,
		DecisionLogger: func(decision, rationale string) error {
			logged = true
			return nil
		},
	})
	checks := []ciloop.CheckResult{{Name: "go-test", State: ciloop.StateFailure}}
	fixer.Fix(checks)
	if !logged {
		t.Error("DecisionLogger was not called")
	}
}

// Integration: RunLoop with Fixer wired - go-test failure → fix → green
func TestRunLoopWithFixerGreenAfterFix(t *testing.T) {
	committer := &fakeCommitter{}
	fetcher := &fakeLogFetcher{logs: map[string]string{"run1": goTestFailureLog}}
	fixer := ciloop.NewFixer(ciloop.FixerOptions{
		PRNumber:       42,
		FeatureID:      "F014",
		DiagnosticsDir: t.TempDir(),
		Committer:      committer,
		LogFetcher:     fetcher,
		DecisionLogger: func(decision, rationale string) error { return nil },
	})

	// First poll: go-test fails. Second poll: green.
	runner := newSequence([][]byte{failureJSON, greenJSON})
	opts := ciloop.LoopOptions{
		PRNumber:   42,
		MaxIter:    5,
		PollDelay:  0,
		RetryDelay: 0,
		Poller:     ciloop.NewPoller(runner),
		OnEscalate: func(checks []ciloop.CheckResult) error { return nil },
		Fixer:      fixer,
	}
	result, err := ciloop.RunLoop(opts)
	if err != nil {
		t.Fatal(err)
	}
	if result != ciloop.ResultGreen {
		t.Errorf("expected Green after fix, got %v", result)
	}
	if len(committer.commits) != 1 {
		t.Errorf("expected 1 auto-fix commit, got %d", len(committer.commits))
	}
}
