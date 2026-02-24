package orchestrate

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/fall-out-bug/sdp/internal/guard"
)

// RunGuardCheck runs sdp-guard for the given workstream. Returns error if scope check fails.
func RunGuardCheck(projectRoot, wsID string) error {
	verdict, err := guard.CheckScope(projectRoot, wsID, false)
	if err != nil {
		return fmt.Errorf("guard check: %w", err)
	}
	if verdict.Pass {
		return nil
	}
	return &ScopeViolationError{WSID: wsID, Violations: verdict.Violations}
}

// ScopeViolationError is returned when guard detects out-of-scope changes.
type ScopeViolationError struct {
	WSID       string
	Violations []string
}

func (e *ScopeViolationError) Error() string {
	return fmt.Sprintf("scope violation: %s touched %d out-of-scope files: %s",
		e.WSID, len(e.Violations), strings.Join(e.Violations, ", "))
}

// CreateScopeEscalationBead runs bd create for a scope violation.
func CreateScopeEscalationBead(wsID string, violations []string) error {
	title := fmt.Sprintf("SCOPE VIOLATION: %s touched %s", wsID, strings.Join(violations, ", "))
	if len(title) > 200 {
		title = title[:197] + "..."
	}
	cmd := exec.Command("bd", "create", "--title", title, "--priority", "1", "--labels", "scope-violation")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
