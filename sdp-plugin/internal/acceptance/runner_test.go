package acceptance

import (
	"context"
	"testing"
	"time"
)

func TestRunner_Run_ExitZero(t *testing.T) {
	r := &Runner{
		Command:  "true",
		Timeout:  5 * time.Second,
		Expected: "",
	}
	res, err := r.Run(context.Background())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected pass, got Passed=false: %s", res.Error)
	}
	if res.Duration == 0 {
		t.Error("Duration should be non-zero")
	}
}

func TestRunner_Run_ExitNonZero(t *testing.T) {
	r := &Runner{
		Command: "false",
		Timeout: 5 * time.Second,
	}
	res, err := r.Run(context.Background())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if res.Passed {
		t.Error("expected fail for exit code non-zero")
	}
}

func TestRunner_Run_ExpectedSubstring(t *testing.T) {
	r := &Runner{
		Command:  "echo PASS",
		Timeout:  5 * time.Second,
		Expected: "PASS",
	}
	res, err := r.Run(context.Background())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !res.Passed {
		t.Errorf("expected pass (PASS in output), got %s", res.Error)
	}
}

func TestRunner_Run_ExpectedMissing(t *testing.T) {
	r := &Runner{
		Command:  "echo fail",
		Timeout:  5 * time.Second,
		Expected: "PASS",
	}
	res, err := r.Run(context.Background())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if res.Passed {
		t.Error("expected fail when expected substring missing")
	}
}

func TestRunner_Run_Timeout(t *testing.T) {
	r := &Runner{
		Command: "sleep 10",
		Timeout: 50 * time.Millisecond,
	}
	res, err := r.Run(context.Background())
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if res.Passed {
		t.Error("expected fail on timeout")
	}
	if res.Duration < 50*time.Millisecond || res.Duration > 2*time.Second {
		t.Errorf("duration should be ~timeout, got %v", res.Duration)
	}
}

func TestParseTimeout(t *testing.T) {
	d, err := ParseTimeout("30s")
	if err != nil {
		t.Fatalf("ParseTimeout: %v", err)
	}
	if d != 30*time.Second {
		t.Errorf("want 30s, got %v", d)
	}
	d, err = ParseTimeout("")
	if err != nil {
		t.Fatalf("ParseTimeout empty: %v", err)
	}
	if d != 30*time.Second {
		t.Errorf("empty want 30s default, got %v", d)
	}
}
