package sdpinit

import (
	"fmt"
)

// GuidedStep represents a single step in the guided setup flow
type GuidedStep struct {
	ID          string
	Name        string
	Description string
	Check       func() GuidedStepResult
	Fix         func() error
	FixCommand  string
}

// GuidedStepResult contains the result of a step check
type GuidedStepResult struct {
	Passed  bool
	Message string
	Details string
}

// GuidedConfig contains configuration for guided setup
type GuidedConfig struct {
	ProjectType string
	SkipBeads   bool
	AutoFix     bool
}

// GuidedResult contains the complete guided setup result
type GuidedResult struct {
	Steps     []GuidedStepResult
	AllPassed bool
	NextSteps []string
}

// RunGuided executes the guided setup flow
func RunGuided(cfg GuidedConfig) (*GuidedResult, error) {
	result := &GuidedResult{
		Steps:     []GuidedStepResult{},
		AllPassed: true,
		NextSteps: []string{},
	}

	steps := getGuidedSteps()

	fmt.Println("SDP Guided Setup")
	fmt.Println("================")
	fmt.Println("This wizard will guide you through SDP setup.")
	fmt.Println()

	for i, step := range steps {
		passed := processStep(i, len(steps), step, cfg, result)
		if !passed {
			result.AllPassed = false
		}
		fmt.Println()
	}

	setNextSteps(result)
	return result, nil
}

// processStep handles a single guided setup step
func processStep(index, total int, step GuidedStep, cfg GuidedConfig, result *GuidedResult) bool {
	fmt.Printf("Step %d/%d: %s\n", index+1, total, step.Name)
	fmt.Printf("  %s\n", step.Description)

	var stepResult GuidedStepResult
	if step.ID == "project-type" && cfg.ProjectType != "" {
		stepResult = CheckProjectTypeWithOverride(cfg.ProjectType)
	} else {
		stepResult = step.Check()
	}
	result.Steps = append(result.Steps, stepResult)

	if stepResult.Passed {
		fmt.Printf("  [PASS] %s\n", stepResult.Message)
		return true
	}

	return handleFailedStep(step, cfg, stepResult)
}

// handleFailedStep processes a failed step with optional auto-fix
func handleFailedStep(step GuidedStep, cfg GuidedConfig, stepResult GuidedStepResult) bool {
	fmt.Printf("  [FAIL] %s\n", stepResult.Message)

	if step.FixCommand != "" {
		fmt.Printf("  Fix: %s\n", step.FixCommand)
	}

	if !cfg.AutoFix || step.Fix == nil {
		return false
	}

	fmt.Println("  Attempting automatic fix...")
	if err := step.Fix(); err != nil {
		fmt.Printf("  Fix failed: %v\n", err)
		return false
	}

	stepResult = step.Check()
	if stepResult.Passed {
		fmt.Printf("  [FIXED] %s\n", stepResult.Message)
		return true
	}
	return false
}

// setNextSteps sets the appropriate next steps based on result
func setNextSteps(result *GuidedResult) {
	if result.AllPassed {
		result.NextSteps = []string{
			"Run 'sdp doctor' to verify your setup",
			"Run 'sdp status' to see project status",
			"Try 'sdp plan \"your feature idea\"' to plan a feature",
		}
	} else {
		result.NextSteps = []string{
			"Fix the issues above and run 'sdp init --guided' again",
			"Or run 'sdp doctor --repair' for automatic fixes",
		}
	}
}

// PrintGuidedResult prints the final result of guided setup
func PrintGuidedResult(result *GuidedResult) {
	fmt.Println()
	fmt.Println("Setup Summary")
	fmt.Println("=============")

	passCount := 0
	for _, step := range result.Steps {
		if step.Passed {
			passCount++
		}
	}

	fmt.Printf("Passed: %d/%d\n", passCount, len(result.Steps))

	if result.AllPassed {
		fmt.Println("\nSetup complete! Here's what to do next:")
	} else {
		fmt.Println("\nSetup incomplete. Next steps:")
	}

	for i, step := range result.NextSteps {
		fmt.Printf("  %d. %s\n", i+1, step)
	}
}

// getGuidedSteps returns the ordered list of setup steps
func getGuidedSteps() []GuidedStep {
	return []GuidedStep{
		{
			ID:          "git",
			Name:        "Git Installation",
			Description: "Check if Git is installed",
			Check:       checkGitStep,
			FixCommand:  "Install Git from https://git-scm.com",
		},
		{
			ID:          "git-repo",
			Name:        "Git Repository",
			Description: "Check if current directory is a Git repository",
			Check:       checkGitRepoStep,
			FixCommand:  "Run 'git init' to initialize a repository",
			Fix:         fixGitRepo,
		},
		{
			ID:          "project-type",
			Name:        "Project Detection",
			Description: "Detect project type for SDP configuration",
			Check:       checkProjectTypeStep,
			FixCommand:  "Run 'sdp init --project-type <type>' to specify manually",
		},
		{
			ID:          "claude-code",
			Name:        "Claude Code CLI",
			Description: "Check if Claude Code CLI is installed (optional)",
			Check:       checkClaudeCodeStep,
			FixCommand:  "Install Claude Code from https://claude.ai/code",
		},
		{
			ID:          "beads",
			Name:        "Beads CLI",
			Description: "Check if Beads CLI is installed for task tracking",
			Check:       checkBeadsStep,
			FixCommand:  "Run 'brew tap beads-dev/tap && brew install beads'",
		},
	}
}
