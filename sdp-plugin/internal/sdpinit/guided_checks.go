package sdpinit

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// CheckProjectTypeWithOverride checks project type, allowing explicit override
func CheckProjectTypeWithOverride(explicitType string) GuidedStepResult {
	if explicitType != "" {
		return GuidedStepResult{
			Passed:  true,
			Message: fmt.Sprintf("Using specified type: %s", explicitType),
			Details: explicitType,
		}
	}
	return checkProjectTypeStep()
}

func checkGitStep() GuidedStepResult {
	if _, err := exec.LookPath("git"); err != nil {
		return GuidedStepResult{
			Passed:  false,
			Message: "Git is not installed",
			Details: "Git is required for SDP version control features",
		}
	}

	cmd := exec.Command("git", "--version")
	output, err := cmd.Output()
	if err != nil {
		return GuidedStepResult{
			Passed:  true,
			Message: "Git installed (version unknown)",
		}
	}

	version := strings.TrimSpace(string(output))
	return GuidedStepResult{
		Passed:  true,
		Message: fmt.Sprintf("Git installed (%s)", version),
	}
}

func checkGitRepoStep() GuidedStepResult {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return GuidedStepResult{
			Passed:  false,
			Message: "Not a Git repository",
			Details: "SDP requires a Git repository for workstream tracking",
		}
	}
	return GuidedStepResult{
		Passed:  true,
		Message: "Git repository detected",
	}
}

func fixGitRepo() error {
	cmd := exec.Command("git", "init")
	return cmd.Run()
}

func checkProjectTypeStep() GuidedStepResult {
	projectType := DetectProjectType()
	if projectType == "unknown" {
		return GuidedStepResult{
			Passed:  false,
			Message: "Could not detect project type",
			Details: "Specify project type with --project-type flag",
		}
	}
	return GuidedStepResult{
		Passed:  true,
		Message: fmt.Sprintf("Detected: %s", projectType),
		Details: projectType,
	}
}

func checkClaudeCodeStep() GuidedStepResult {
	if _, err := exec.LookPath("claude"); err != nil {
		return GuidedStepResult{
			Passed:  true,
			Message: "Claude Code CLI not installed (optional)",
			Details: "Install for enhanced Claude Code integration",
		}
	}

	cmd := exec.Command("claude", "--version")
	output, err := cmd.Output()
	if err != nil {
		return GuidedStepResult{
			Passed:  true,
			Message: "Claude Code CLI installed (version unknown)",
		}
	}

	version := strings.TrimSpace(string(output))
	return GuidedStepResult{
		Passed:  true,
		Message: fmt.Sprintf("Claude Code CLI installed (%s)", version),
	}
}

func checkBeadsStep() GuidedStepResult {
	if _, err := exec.LookPath("bd"); err != nil {
		return GuidedStepResult{
			Passed:  true,
			Message: "Beads CLI not installed (optional)",
			Details: "Install for task tracking integration",
		}
	}

	cmd := exec.Command("bd", "--version")
	output, err := cmd.Output()
	if err != nil {
		return GuidedStepResult{
			Passed:  true,
			Message: "Beads CLI installed (version unknown)",
		}
	}

	version := strings.TrimSpace(string(output))
	return GuidedStepResult{
		Passed:  true,
		Message: fmt.Sprintf("Beads CLI installed (%s)", version),
	}
}
