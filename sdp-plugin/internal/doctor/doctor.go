package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CheckResult struct {
	Name    string
	Status  string // "ok", "warning", "error"
	Message string
}

// RunOptions controls which checks are run
type RunOptions struct {
	DriftCheck bool // Run drift detection on recent workstreams
}

// RunWithOptions runs doctor checks with the given options
func RunWithOptions(opts RunOptions) []CheckResult {
	results := []CheckResult{}

	// Check 1: Git
	results = append(results, checkGit())

	// Check 2: Claude Code
	results = append(results, checkClaudeCode())

	// Check 3: Go (for building binary)
	results = append(results, checkGo())

	// Check 4: IDE integration
	results = append(results, checkIDEIntegration())

	// Check 5: File permissions on sensitive data
	results = append(results, checkFilePermissions())

	// Check 6: .sdp/config.yml validation (AC6: if present)
	results = append(results, checkProjectConfig())

	// Check 7: Drift detection (if requested)
	if opts.DriftCheck {
		results = append(results, checkDrift())
	}

	return results
}

// Run runs all standard doctor checks
func Run() []CheckResult {
	return RunWithOptions(RunOptions{})
}

func checkGit() CheckResult {
	if _, err := exec.LookPath("git"); err != nil {
		return CheckResult{
			Name:    "Git",
			Status:  "error",
			Message: "Git not found. Install from https://git-scm.com",
		}
	}

	// Get version
	cmd := exec.Command("git", "--version")
	output, err := cmd.Output()
	if err != nil {
		return CheckResult{
			Name:    "Git",
			Status:  "error",
			Message: "Failed to get version",
		}
	}
	version := strings.TrimSpace(string(output))

	return CheckResult{
		Name:    "Git",
		Status:  "ok",
		Message: fmt.Sprintf("Installed (%s)", version),
	}
}

func checkClaudeCode() CheckResult {
	if _, err := exec.LookPath("claude"); err != nil {
		return CheckResult{
			Name:    "Claude Code",
			Status:  "warning",
			Message: "Claude Code CLI not found. Plugin will work in Claude Desktop app.",
		}
	}

	// Get version
	cmd := exec.Command("claude", "--version")
	output, err := cmd.Output()
	if err != nil {
		return CheckResult{
			Name:    "Claude Code",
			Status:  "ok", // Don't fail if version check fails
			Message: "Installed (version unknown)",
		}
	}
	version := strings.TrimSpace(string(output))

	return CheckResult{
		Name:    "Claude Code",
		Status:  "ok",
		Message: fmt.Sprintf("Installed (%s)", version),
	}
}

func checkGo() CheckResult {
	if _, err := exec.LookPath("go"); err != nil {
		return CheckResult{
			Name:    "Go",
			Status:  "warning",
			Message: "Go not found. Required only for building SDP binary.",
		}
	}

	// Get version
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return CheckResult{
			Name:    "Go",
			Status:  "error",
			Message: "Failed to get version",
		}
	}
	version := strings.TrimSpace(string(output))

	return CheckResult{
		Name:    "Go",
		Status:  "ok",
		Message: fmt.Sprintf("Installed (%s)", version),
	}
}

func checkClaudeDir() CheckResult {
	return checkIDEIntegration()
}

func checkIDEIntegration() CheckResult {
	checks := []struct {
		label    string
		basePath string
		required []string
	}{
		{label: "Claude", basePath: ".claude", required: []string{"skills", "agents", "validators"}},
		{label: "Cursor", basePath: ".cursor", required: []string{"skills", "agents"}},
		{label: "OpenCode", basePath: ".opencode", required: []string{"skills", "agents"}},
		{label: "Codex", basePath: ".codex", required: []string{"skills", "agents", "INSTALL.md"}},
	}

	found := []string{}
	incomplete := []string{}

	for _, check := range checks {
		info, err := os.Stat(check.basePath)
		if err != nil || !info.IsDir() {
			continue
		}

		found = append(found, check.label)
		missing := []string{}
		for _, required := range check.required {
			if _, err := os.Stat(filepath.Join(check.basePath, required)); err != nil {
				missing = append(missing, required)
			}
		}
		if len(missing) > 0 {
			incomplete = append(incomplete, fmt.Sprintf("%s missing %s", check.label, strings.Join(missing, ", ")))
		}
	}

	if len(found) == 0 {
		return CheckResult{
			Name:    "IDE integration",
			Status:  "error",
			Message: "Not found. Run the SDP installer or set up one of .claude/, .cursor/, .opencode/, or .codex/",
		}
	}

	if len(incomplete) > 0 {
		return CheckResult{
			Name:    "IDE integration",
			Status:  "warning",
			Message: fmt.Sprintf("Found: %s. Incomplete: %s", strings.Join(found, ", "), strings.Join(incomplete, "; ")),
		}
	}

	return CheckResult{
		Name:    "IDE integration",
		Status:  "ok",
		Message: fmt.Sprintf("Found: %s", strings.Join(found, ", ")),
	}
}
