package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fall-out-bug/sdp/internal/drift"
)

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
			}
		}

		// Stop if we have enough workstreams
		if len(workstreams) >= maxTotal {
			break
		}
	}

	return workstreams, nil
}
