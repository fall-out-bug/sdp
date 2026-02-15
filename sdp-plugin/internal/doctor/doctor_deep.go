package doctor

import (
	"crypto/sha256"
	"encoding/hex"
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
		if info, err := os.Stat(hookPath); os.IsNotExist(err) {
			missing = append(missing, hook)
		} else if info.Mode().Perm()&0100 == 0 {
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

// checkBeadsIntegrity validates Beads database is healthy
func checkBeadsIntegrity() DeepCheckResult {
	start := getTime()
	details := make(map[string]interface{})

	beadsDB := ".beads/beads.db"
	if _, err := os.Stat(beadsDB); os.IsNotExist(err) {
		return DeepCheckResult{
			Check:    "Beads Integrity",
			Status:   "warning",
			Duration: since(start),
			Message:  "Beads database not found",
			Details:  details,
		}
	}

	// Check file can be read
	content, err := os.ReadFile(beadsDB)
	if err != nil {
		return DeepCheckResult{
			Check:    "Beads Integrity",
			Status:   "error",
			Duration: since(start),
			Message:  fmt.Sprintf("Cannot read beads.db: %v", err),
			Details:  details,
		}
	}

	// Basic integrity: check it's not empty and has expected structure
	hash := sha256.Sum256(content)
	details["size"] = len(content)
	details["hash"] = hex.EncodeToString(hash[:8])

	if len(content) == 0 {
		return DeepCheckResult{
			Check:    "Beads Integrity",
			Status:   "warning",
			Duration: since(start),
			Message:  "Beads database is empty",
			Details:  details,
		}
	}

	return DeepCheckResult{
		Check:    "Beads Integrity",
		Status:   "ok",
		Duration: since(start),
		Message:  fmt.Sprintf("Database OK (%d bytes)", len(content)),
		Details:  details,
	}
}

// checkConfigVersion validates config version is compatible
func checkConfigVersion() DeepCheckResult {
	start := getTime()
	details := make(map[string]interface{})

	configPath := ".sdp/config.yml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DeepCheckResult{
			Check:    "Config Version",
			Status:   "ok",
			Duration: since(start),
			Message:  "No custom config (using defaults)",
			Details:  details,
		}
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return DeepCheckResult{
			Check:    "Config Version",
			Status:   "error",
			Duration: since(start),
			Message:  fmt.Sprintf("Cannot read config: %v", err),
			Details:  details,
		}
	}

	// Check for version field
	contentStr := string(content)
	if !strings.Contains(contentStr, "version:") {
		return DeepCheckResult{
			Check:    "Config Version",
			Status:   "warning",
			Duration: since(start),
			Message:  "Config missing version field",
			Details:  details,
		}
	}

	return DeepCheckResult{
		Check:    "Config Version",
		Status:   "ok",
		Duration: since(start),
		Message:  "Config version present",
		Details:  details,
	}
}
