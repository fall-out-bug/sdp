package doctor

import (
	"fmt"
	"os"
)

// RepairAction represents the result of a repair attempt
type RepairAction struct {
	Check   string // Check name that was repaired
	Status  string // "fixed", "manual", "skipped", "failed"
	Message string // Human-readable result
}

// RunWithRepair runs all checks and attempts to fix issues automatically
func RunWithRepair() []RepairAction {
	actions := make([]RepairAction, 0, 3)

	// Fix 1: File permissions
	actions = append(actions, repairFilePermissions())

	// Fix 2: Missing .claude subdirectories
	actions = append(actions, repairClaudeDirs())

	// Fix 3: Missing .sdp directory structure
	actions = append(actions, repairSDPDirs())

	return actions
}

// repairClaudeDirs creates missing .claude subdirectories
func repairClaudeDirs() RepairAction {
	requiredDirs := []string{
		".claude/skills",
		".claude/agents",
		".claude/validators",
	}

	created := []string{}
	failed := []string{}

	for _, dir := range requiredDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				failed = append(failed, dir)
			} else {
				created = append(created, dir)
			}
		}
	}

	if len(failed) > 0 {
		return RepairAction{
			Check:   ".claude/ Structure",
			Status:  "failed",
			Message: fmt.Sprintf("Failed to create: %v", failed),
		}
	}

	if len(created) == 0 {
		return RepairAction{
			Check:   ".claude/ Structure",
			Status:  "skipped",
			Message: "All directories already exist",
		}
	}

	return RepairAction{
		Check:   ".claude/ Structure",
		Status:  "fixed",
		Message: fmt.Sprintf("Created: %v", created),
	}
}

// repairSDPDirs creates missing .sdp subdirectories
func repairSDPDirs() RepairAction {
	requiredDirs := []string{
		".sdp/log",
	}

	created := []string{}
	failed := []string{}

	// Check if .sdp exists first
	if _, err := os.Stat(".sdp"); os.IsNotExist(err) {
		return RepairAction{
			Check:   ".sdp/ Structure",
			Status:  "manual",
			Message: "Run 'sdp init' to initialize .sdp directory",
		}
	}

	for _, dir := range requiredDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				failed = append(failed, dir)
			} else {
				created = append(created, dir)
			}
		}
	}

	if len(failed) > 0 {
		return RepairAction{
			Check:   ".sdp/ Structure",
			Status:  "failed",
			Message: fmt.Sprintf("Failed to create: %v", failed),
		}
	}

	if len(created) == 0 {
		return RepairAction{
			Check:   ".sdp/ Structure",
			Status:  "skipped",
			Message: "All directories already exist",
		}
	}

	return RepairAction{
		Check:   ".sdp/ Structure",
		Status:  "fixed",
		Message: fmt.Sprintf("Created: %v", created),
	}
}

// HasUnfixableErrors checks if any action requires manual intervention
func HasUnfixableErrors(actions []RepairAction) bool {
	for _, a := range actions {
		if a.Status == "manual" || a.Status == "failed" {
			return true
		}
	}
	return false
}
