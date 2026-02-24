package ciloop_test

import (
	"errors"
	"testing"
	"time"

	"github.com/fall-out-bug/sdp/internal/ciloop"
)

// loopRunner simulates sequences of gh responses across calls.
type sequenceRunner struct {
	responses [][]byte
	errs      []error
	call      int
}

func (s *sequenceRunner) Run(_ string, _ ...string) ([]byte, error) {
	i := s.call
	if i >= len(s.responses) {
		i = len(s.responses) - 1
	}
	s.call++
	return s.responses[i], s.errs[i]
}

func newSequence(responses [][]byte) *sequenceRunner {
	errs := make([]error, len(responses))
	return &sequenceRunner{responses: responses, errs: errs}
}

func TestRunLoopGreenFirstTry(t *testing.T) {
	runner := newSequence([][]byte{greenJSON})
	opts := ciloop.LoopOptions{
		PRNumber:    42,
		MaxIter:     5,
		PollDelay:   0,
		RetryDelay:  0,
		Poller:      ciloop.NewPoller(runner),
		OnEscalate:  func(checks []ciloop.CheckResult) error { return nil },
	}
	result, err := ciloop.RunLoop(opts)
	if err != nil {
		t.Fatal(err)
	}
	if result != ciloop.ResultGreen {
		t.Errorf("expected Green, got %v", result)
	}
}

func TestRunLoopPendingThenGreen(t *testing.T) {
	runner := newSequence([][]byte{pendingJSON, greenJSON})
	opts := ciloop.LoopOptions{
		PRNumber:   42,
		MaxIter:    5,
		PollDelay:  0,
		RetryDelay: 0,
		Poller:     ciloop.NewPoller(runner),
		OnEscalate: func(checks []ciloop.CheckResult) error { return nil },
	}
	result, err := ciloop.RunLoop(opts)
	if err != nil {
		t.Fatal(err)
	}
	if result != ciloop.ResultGreen {
		t.Errorf("expected Green, got %v", result)
	}
}

func TestRunLoopEscalatesOnUnfixableFailure(t *testing.T) {
	secretsFailure := []byte(`[{"name":"secrets-scan","state":"FAILURE"}]`)
	runner := newSequence([][]byte{secretsFailure})
	escalated := false
	opts := ciloop.LoopOptions{
		PRNumber:   42,
		MaxIter:    5,
		PollDelay:  0,
		RetryDelay: 0,
		Poller:     ciloop.NewPoller(runner),
		OnEscalate: func(checks []ciloop.CheckResult) error {
			escalated = true
			return nil
		},
	}
	result, err := ciloop.RunLoop(opts)
	if err != nil {
		t.Fatal(err)
	}
	if result != ciloop.ResultEscalated {
		t.Errorf("expected Escalated, got %v", result)
	}
	if !escalated {
		t.Error("OnEscalate was not called")
	}
}

func TestRunLoopExceedsMaxIter(t *testing.T) {
	goTestFailure := []byte(`[{"name":"go-test","state":"FAILURE"}]`)
	responses := make([][]byte, 10)
	for i := range responses {
		responses[i] = goTestFailure
	}
	runner := newSequence(responses)
	// Use a fake Fixer that always succeeds so iterations are consumed.
	opts := ciloop.LoopOptions{
		PRNumber:   42,
		MaxIter:    3,
		PollDelay:  0,
		RetryDelay: 0,
		Poller:     ciloop.NewPoller(runner),
		OnEscalate: func(checks []ciloop.CheckResult) error { return nil },
		Fixer:      &fakeFixer{},
	}
	result, err := ciloop.RunLoop(opts)
	if err != nil {
		t.Fatal(err)
	}
	if result != ciloop.ResultMaxIter {
		t.Errorf("expected MaxIter, got %v", result)
	}
}

// fakeFixer is a Fixer that always succeeds without side effects.
type fakeFixer struct{}

func (f *fakeFixer) Fix(_ []ciloop.CheckResult) error { return nil }

func TestRunLoopNilFixerEscalatesAutoFixable(t *testing.T) {
	goTestFailure := []byte(`[{"name":"go-test","state":"FAILURE"}]`)
	runner := newSequence([][]byte{goTestFailure})
	escalated := false
	opts := ciloop.LoopOptions{
		PRNumber:   42,
		MaxIter:    5,
		PollDelay:  0,
		RetryDelay: 0,
		Poller:     ciloop.NewPoller(runner),
		OnEscalate: func(checks []ciloop.CheckResult) error {
			escalated = true
			return nil
		},
		Fixer: nil,
	}
	result, err := ciloop.RunLoop(opts)
	if err != nil {
		t.Fatal(err)
	}
	if result != ciloop.ResultEscalated {
		t.Errorf("expected Escalated when Fixer is nil, got %v", result)
	}
	if !escalated {
		t.Error("OnEscalate was not called")
	}
}

func TestRunLoopMaxPendingRetriesEscalates(t *testing.T) {
	runner := newSequence([][]byte{pendingJSON, pendingJSON, pendingJSON, pendingJSON})
	escalated := false
	opts := ciloop.LoopOptions{
		PRNumber:          42,
		MaxIter:           5,
		MaxPendingRetries: 2,
		PollDelay:         0,
		RetryDelay:        0,
		Poller:            ciloop.NewPoller(runner),
		OnEscalate: func(checks []ciloop.CheckResult) error {
			escalated = true
			return nil
		},
	}
	result, err := ciloop.RunLoop(opts)
	if err != nil {
		t.Fatal(err)
	}
	if result != ciloop.ResultEscalated {
		t.Errorf("expected Escalated after MaxPendingRetries, got %v", result)
	}
	if !escalated {
		t.Error("OnEscalate was not called for max pending retries")
	}
}

func TestLoopOptionsPollDelayIsRespected(t *testing.T) {
	runner := newSequence([][]byte{greenJSON})
	start := time.Now()
	opts := ciloop.LoopOptions{
		PRNumber:   42,
		MaxIter:    5,
		PollDelay:  10 * time.Millisecond,
		RetryDelay: 0,
		Poller:     ciloop.NewPoller(runner),
		OnEscalate: func(checks []ciloop.CheckResult) error { return nil },
	}
	ciloop.RunLoop(opts)
	elapsed := time.Since(start)
	if elapsed < 10*time.Millisecond {
		t.Errorf("expected poll delay of at least 10ms, elapsed: %v", elapsed)
	}
}

// TestOnEscalateErrorPath verifies that when OnEscalate returns an error, RunLoop propagates it (028g).
func TestOnEscalateErrorPath(t *testing.T) {
	secretsFailure := []byte(`[{"name":"secrets-scan","state":"FAILURE"}]`)
	runner := newSequence([][]byte{secretsFailure})
	wantErr := errors.New("escalation callback failed")
	opts := ciloop.LoopOptions{
		PRNumber:   42,
		MaxIter:    5,
		PollDelay:  0,
		RetryDelay: 0,
		Poller:     ciloop.NewPoller(runner),
		OnEscalate: func(checks []ciloop.CheckResult) error { return wantErr },
	}
	result, err := ciloop.RunLoop(opts)
	if err != wantErr {
		t.Errorf("expected OnEscalate error, got %v", err)
	}
	if result != ciloop.ResultEscalated {
		t.Errorf("expected Escalated, got %v", result)
	}
}

// TestFixerFixFailureEscalates verifies that when Fixer.Fix returns error, RunLoop escalates and propagates it (850r).
func TestFixerFixFailureEscalates(t *testing.T) {
	goTestFailure := []byte(`[{"name":"go-test","state":"FAILURE"}]`)
	runner := newSequence([][]byte{goTestFailure})
	wantErr := errors.New("commit failed")
	opts := ciloop.LoopOptions{
		PRNumber:   42,
		MaxIter:    5,
		PollDelay:  0,
		RetryDelay: 0,
		Poller:     ciloop.NewPoller(runner),
		OnEscalate: func(checks []ciloop.CheckResult) error { return nil },
		Fixer:      &breakingFixer{err: wantErr},
	}
	result, err := ciloop.RunLoop(opts)
	if err != wantErr {
		t.Errorf("expected Fixer error, got %v", err)
	}
	if result != ciloop.ResultEscalated {
		t.Errorf("expected Escalated, got %v", result)
	}
}

type breakingFixer struct{ err error }

func (f *breakingFixer) Fix(_ []ciloop.CheckResult) error { return f.err }

// TestFixPushStillFailingMaxIter verifies fix->push->still failing->max iter path (65dj).
func TestFixPushStillFailingMaxIter(t *testing.T) {
	goTestFailure := []byte(`[{"name":"go-test","state":"FAILURE"}]`)
	responses := make([][]byte, 5)
	for i := range responses {
		responses[i] = goTestFailure
	}
	runner := newSequence(responses)
	opts := ciloop.LoopOptions{
		PRNumber:   3,
		MaxIter:    3,
		PollDelay:  0,
		RetryDelay: 0,
		Poller:     ciloop.NewPoller(runner),
		OnEscalate: func(checks []ciloop.CheckResult) error { return nil },
		Fixer:      &fakeFixer{},
	}
	result, err := ciloop.RunLoop(opts)
	if err != nil {
		t.Fatal(err)
	}
	if result != ciloop.ResultMaxIter {
		t.Errorf("expected MaxIter, got %v", result)
	}
}
