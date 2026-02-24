package ciloop_test

import (
	"errors"
	"testing"

	"github.com/fall-out-bug/sdp/internal/ciloop"
)

type fakeRunner struct {
	output []byte
	err    error
}

func (f *fakeRunner) Run(_ string, _ ...string) ([]byte, error) {
	return f.output, f.err
}

var greenJSON = []byte(`[
  {"name": "go-test",    "state": "SUCCESS"},
  {"name": "go-build",   "state": "SUCCESS"}
]`)

var pendingJSON = []byte(`[
  {"name": "go-test",  "state": "PENDING"},
  {"name": "go-build", "state": "SUCCESS"}
]`)

var failureJSON = []byte(`[
  {"name": "go-test",  "state": "FAILURE"},
  {"name": "go-build", "state": "SUCCESS"}
]`)

var mixedJSON = []byte(`[
  {"name": "go-test",    "state": "SUCCESS"},
  {"name": "secrets",    "state": "FAILURE"},
  {"name": "k8s-validate", "state": "IN_PROGRESS"}
]`)

func TestGetChecksGreen(t *testing.T) {
	p := ciloop.NewPoller(&fakeRunner{output: greenJSON})
	checks, err := p.GetChecks(42)
	if err != nil {
		t.Fatal(err)
	}
	if len(checks) != 2 {
		t.Fatalf("expected 2 checks, got %d", len(checks))
	}
	for _, c := range checks {
		if c.State != ciloop.StateSuccess {
			t.Errorf("expected SUCCESS for %q, got %q", c.Name, c.State)
		}
	}
}

func TestGetChecksPending(t *testing.T) {
	p := ciloop.NewPoller(&fakeRunner{output: pendingJSON})
	checks, err := p.GetChecks(42)
	if err != nil {
		t.Fatal(err)
	}
	pending := ciloop.FilterByState(checks, ciloop.StatePending)
	if len(pending) != 1 || pending[0].Name != "go-test" {
		t.Errorf("expected 1 pending check named go-test, got %v", pending)
	}
}

func TestGetChecksFailure(t *testing.T) {
	p := ciloop.NewPoller(&fakeRunner{output: failureJSON})
	checks, err := p.GetChecks(42)
	if err != nil {
		t.Fatal(err)
	}
	failing := ciloop.FilterByState(checks, ciloop.StateFailure)
	if len(failing) != 1 || failing[0].Name != "go-test" {
		t.Errorf("expected 1 failure check named go-test, got %v", failing)
	}
}

func TestGetChecksCommandError(t *testing.T) {
	p := ciloop.NewPoller(&fakeRunner{err: errors.New("gh: not found")})
	_, err := p.GetChecks(42)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetChecksMixed(t *testing.T) {
	p := ciloop.NewPoller(&fakeRunner{output: mixedJSON})
	checks, err := p.GetChecks(42)
	if err != nil {
		t.Fatal(err)
	}
	if len(checks) != 3 {
		t.Fatalf("expected 3 checks, got %d", len(checks))
	}
	inProgress := ciloop.FilterByState(checks, ciloop.StateInProgress)
	if len(inProgress) != 1 {
		t.Errorf("expected 1 IN_PROGRESS check, got %d", len(inProgress))
	}
}

func TestIsAllGreen(t *testing.T) {
	green := []ciloop.CheckResult{
		{Name: "a", State: ciloop.StateSuccess},
		{Name: "b", State: ciloop.StateSuccess},
	}
	if !ciloop.IsAllGreen(green) {
		t.Error("expected all green")
	}

	mixed := []ciloop.CheckResult{
		{Name: "a", State: ciloop.StateSuccess},
		{Name: "b", State: ciloop.StatePending},
	}
	if ciloop.IsAllGreen(mixed) {
		t.Error("expected not all green when pending present")
	}
}
