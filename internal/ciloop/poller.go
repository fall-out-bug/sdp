package ciloop

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/sdputil"
)

// CheckState represents the state of a CI check.
type CheckState string

const (
	StatePending    CheckState = "PENDING"
	StateSuccess    CheckState = "SUCCESS"
	StateFailure    CheckState = "FAILURE"
	StateError      CheckState = "ERROR"
	StateInProgress CheckState = "IN_PROGRESS"
)

// CheckResult holds the name and state of a single CI check.
type CheckResult struct {
	Name  string     `json:"name"`
	State CheckState `json:"state"`
}

// CommandRunner executes an external command and returns its stdout.
type CommandRunner interface {
	Run(name string, args ...string) ([]byte, error)
}

// Poller polls GitHub PR checks via the gh CLI.
type Poller struct {
	runner CommandRunner
}

// NewPoller creates a Poller backed by the given runner.
func NewPoller(runner CommandRunner) *Poller {
	return &Poller{runner: runner}
}

// GetChecks fetches current check states for the given PR number.
// Retries with exponential backoff (2s, 4s, 8s) on transient failures, max 3 retries.
func (p *Poller) GetChecks(prNumber int) ([]CheckResult, error) {
	delays := []time.Duration{2 * time.Second, 4 * time.Second, 8 * time.Second}
	var out []byte
	var err error
	for attempt := 0; attempt <= len(delays); attempt++ {
		out, err = p.runner.Run("gh", "pr", "checks", strconv.Itoa(prNumber), "--json", "name,state")
		if err == nil {
			break
		}
		if attempt < len(delays) {
			time.Sleep(delays[attempt])
		} else {
			return nil, fmt.Errorf("gh pr checks: %w", err)
		}
	}
	var raw []map[string]string
	if err := json.NewDecoder(io.LimitReader(bytes.NewReader(out), sdputil.MaxJSONDecodeBytes)).Decode(&raw); err != nil {
		return nil, fmt.Errorf("parse checks JSON: %w", err)
	}
	results := make([]CheckResult, 0, len(raw))
	for _, r := range raw {
		results = append(results, CheckResult{
			Name:  r["name"],
			State: CheckState(strings.ToUpper(r["state"])),
		})
	}
	return results, nil
}

// FilterByState returns checks matching the given state.
func FilterByState(checks []CheckResult, state CheckState) []CheckResult {
	var out []CheckResult
	for _, c := range checks {
		if c.State == state {
			out = append(out, c)
		}
	}
	return out
}

// IsAllGreen returns true when all checks are in SUCCESS state.
func IsAllGreen(checks []CheckResult) bool {
	if len(checks) == 0 {
		return false
	}
	for _, c := range checks {
		if c.State != StateSuccess {
			return false
		}
	}
	return true
}
