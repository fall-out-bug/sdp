package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DeepCheckResult represents the result of a deep diagnostic check
type DeepCheckResult struct {
	Check    string
	Status   string // "ok", "warning", "error"
	Duration time.Duration
	Message  string
	Details  map[string]interface{}
}

// getTime returns current time for duration tracking
func getTime() time.Time {
	return time.Now()
}

// since returns duration since start
func since(start time.Time) time.Duration {
	return time.Since(start)
}

// RunDeepChecks runs extended diagnostic checks beyond standard
func RunDeepChecks() []DeepCheckResult {
	results := make([]DeepCheckResult, 0, 5)

	// Check 1: Git hooks integrity
	results = append(results, checkGitHooks())

	// Check 2: Skill files syntax
	results = append(results, checkSkillsSyntax())

	// Check 3: Workstream circular dependencies
	results = append(results, checkWorkstreamCircularDeps())

	// Check 4: Beads database integrity
	results = append(results, checkBeadsIntegrity())

	// Check 5: Config version compatibility
	results = append(results, checkConfigVersion())

	return results
}

// checkGitHooks validates git hooks are properly configured
func checkGitHooks() DeepCheckResult {
	start := getTime()
	details := make(map[string]interface{})

	hooksDir := ".git/hooks"
	if _, err := os.Stat(hooksDir); os.IsNotExist(err) {
		return DeepCheckResult{
			Check:    "Git Hooks",
			Status:   "warning",
			Duration: since(start),
			Message:  "Not in a git repository",
			Details:  details,
		}
	}

	requiredHooks := []string{"pre-commit", "pre-push"}
	missing := []string{}
	invalid := []string{}

	for _, hook := range requiredHooks {
		hookPath := filepath.Join(hooksDir, hook)
		info, err := os.Stat(hookPath)
		if err != nil {
			if os.IsNotExist(err) {
				missing = append(missing, hook)
			}
			// Other errors (permission denied, etc.) - skip
			continue
		}
		if info.Mode().Perm()&0100 == 0 {
			invalid = append(invalid, hook+" (not executable)")
		}
	}

	details["missing"] = missing
	details["invalid"] = invalid

	if len(missing) > 0 || len(invalid) > 0 {
		msg := "Issues found"
		if len(missing) > 0 {
			msg = fmt.Sprintf("Missing: %s", strings.Join(missing, ", "))
		}
		if len(invalid) > 0 {
			msg += fmt.Sprintf("; %s", strings.Join(invalid, ", "))
		}
		return DeepCheckResult{
			Check:    "Git Hooks",
			Status:   "warning",
			Duration: since(start),
			Message:  msg,
			Details:  details,
		}
	}

	return DeepCheckResult{
		Check:    "Git Hooks",
		Status:   "ok",
		Duration: since(start),
		Message:  "All hooks present and executable",
		Details:  details,
	}
}
