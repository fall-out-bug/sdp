package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fall-out-bug/sdp/internal/drift"
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

	// Check 4: .claude/ directory
	results = append(results, checkClaudeDir())

	// Check 5: File permissions on sensitive data
	results = append(results, checkFilePermissions())

	// Check 6: Drift detection (if requested)
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
	if _, err := os.Stat(".claude"); os.IsNotExist(err) {
		return CheckResult{
			Name:    ".claude/ directory",
			Status:  "error",
			Message: "Not found. Run 'sdp init' to initialize",
		}
	}

	// Check if it has expected structure
	dirs := []string{"skills", "agents", "validators"}
	missing := []string{}
	for _, dir := range dirs {
		if _, err := os.Stat(".claude/" + dir); os.IsNotExist(err) {
			missing = append(missing, dir)
		}
	}

	if len(missing) > 0 {
		return CheckResult{
			Name:    ".claude/ directory",
			Status:  "warning",
			Message: fmt.Sprintf("Incomplete (missing: %s)", strings.Join(missing, ", ")),
		}
	}

	return CheckResult{
		Name:    ".claude/ directory",
		Status:  "ok",
		Message: "SDP prompts installed",
	}
}

func checkDrift() CheckResult {
	// Find project root
	projectRoot, err := findProjectRootForDrift()
	if err != nil {
		return CheckResult{
			Name:    "Drift Detection",
			Status:  "warning",
			Message: fmt.Sprintf("Could not find project root: %v", err),
		}
	}

	// Find recent workstreams
	recentWorkstreams, err := findRecentWorkstreamsForDrift(projectRoot)
	if err != nil {
		return CheckResult{
			Name:    "Drift Detection",
			Status:  "warning",
			Message: fmt.Sprintf("Could not find workstreams: %v", err),
		}
	}

	if len(recentWorkstreams) == 0 {
		return CheckResult{
			Name:    "Drift Detection",
			Status:  "ok",
			Message: "No recent workstreams to check",
		}
	}

	// Check drift for each workstream (limit to 5 for performance)
	detector := drift.NewDetector(projectRoot)
	totalErrors := 0
	totalWarnings := 0
	checkedCount := 0
	maxToCheck := 5

	for _, wsPath := range recentWorkstreams {
		if checkedCount >= maxToCheck {
			break
		}

		// Detect drift
		report, err := detector.DetectDrift(wsPath)
		if err != nil {
			continue // Skip workstreams with errors
		}

		checkedCount++

		// Count issues
		for _, issue := range report.Issues {
			if issue.Status == drift.StatusError {
				totalErrors++
			} else if issue.Status == drift.StatusWarning {
				totalWarnings++
			}
		}
	}

	// Generate result message
	if checkedCount == 0 {
		return CheckResult{
			Name:    "Drift Detection",
			Status:  "warning",
			Message: "Could not check any workstreams",
		}
	}

	message := fmt.Sprintf("Checked %d recent workstream(s)", checkedCount)
	if totalErrors > 0 || totalWarnings > 0 {
		message += fmt.Sprintf(" - %d error(s), %d warning(s) found", totalErrors, totalWarnings)
		if totalErrors > 0 {
			message += ". Run 'sdp drift detect <ws-id>' for details"
		}
	}

	status := "ok"
	if totalErrors > 0 {
		status = "error"
	} else if totalWarnings > 0 {
		status = "warning"
	}

	return CheckResult{
		Name:    "Drift Detection",
		Status:  status,
		Message: message,
	}
}

// findProjectRootForDrift finds the project root by looking for docs or .beads directory
func findProjectRootForDrift() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Check if we're in sdp-plugin directory
	if _, err := os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
		// We're in sdp-plugin, go up one level
		parent := filepath.Dir(cwd)
		if _, err := os.Stat(filepath.Join(parent, "docs")); err == nil {
			return parent, nil
		}
	}

	// Check if we're already in project root
	if _, err := os.Stat(filepath.Join(cwd, "docs")); err == nil {
		return cwd, nil
	}

	// Check for .beads directory
	if _, err := os.Stat(filepath.Join(cwd, ".beads")); err == nil {
		return cwd, nil
	}

	// Traverse up looking for project root
	current := cwd
	for {
		if _, err := os.Stat(filepath.Join(current, "docs")); err == nil {
			return current, nil
		}
		if _, err := os.Stat(filepath.Join(current, ".beads")); err == nil {
			return current, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			// Reached root, return current directory
			return cwd, nil
		}
		current = parent
	}
}

// findRecentWorkstreamsForDrift finds recent workstreams to check for drift
func findRecentWorkstreamsForDrift(projectRoot string) ([]string, error) {
	workstreams := []string{}

	// Directories to check
	dirs := []string{
		filepath.Join(projectRoot, "docs", "workstreams", "in_progress"),
		filepath.Join(projectRoot, "docs", "workstreams", "completed"),
	}

	maxTotal := 5 // Maximum total workstreams to return

	for _, dir := range dirs {
		// Check if directory exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		// Read directory
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		// Skip if no entries
		if len(entries) == 0 {
			continue
		}

		// Add workstreams (limit to 5 most recent total)
		count := 0
		for i := len(entries) - 1; i >= 0 && len(workstreams) < maxTotal; i-- {
			if i < 0 || i >= len(entries) {
				continue
			}
			entry := entries[i]
			if entry.IsDir() {
				continue
			}

			// Check if it's a markdown file
			if strings.HasSuffix(entry.Name(), ".md") {
				wsPath := filepath.Join(dir, entry.Name())
				workstreams = append(workstreams, wsPath)
				count++
			}
		}

		// Stop if we have enough workstreams
		if len(workstreams) >= maxTotal {
			break
		}
	}

	return workstreams, nil
}

func checkFilePermissions() CheckResult {
	// List of sensitive files to check
	sensitiveFiles := []string{
		filepath.Join(os.Getenv("HOME"), ".sdp", "telemetry.jsonl"),
		".beads/beads.db",
		".oneshot",
	}

	insecureFiles := []string{}
	for _, path := range sensitiveFiles {
		info, err := os.Stat(path)
		if err != nil {
			// File doesn't exist, skip
			continue
		}

		// Check if file or directory
		if info.IsDir() {
			// Check files in directory
			entries, err := os.ReadDir(path)
			if err != nil {
				continue
			}

			for _, entry := range entries {
				if entry.IsDir() {
					continue
				}

				fullPath := filepath.Join(path, entry.Name())
				fileInfo, err := os.Stat(fullPath)
				if err != nil {
					continue
				}

				// Check permissions (should be 0600 for files)
				if fileInfo.Mode().Perm()&0077 != 0 {
					insecureFiles = append(insecureFiles, fullPath)
				}
			}
		} else {
			// Check single file permissions
			if info.Mode().Perm()&0077 != 0 {
				insecureFiles = append(insecureFiles, path)
			}
		}
	}

	if len(insecureFiles) > 0 {
		return CheckResult{
			Name:    "File Permissions",
			Status:  "warning",
			Message: fmt.Sprintf("Sensitive files have insecure permissions: %s (run 'chmod 0600 <file>' to fix)", strings.Join(insecureFiles, ", ")),
		}
	}

	return CheckResult{
		Name:    "File Permissions",
		Status:  "ok",
		Message: "All sensitive files have secure permissions",
	}
}
