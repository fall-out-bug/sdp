package trial

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Trial represents a trial execution session
type Trial struct {
	BranchName      string
	OriginalBranch  string
	ProjectPath     string
	TaskDescription string
	StartTime       time.Time
}

// TrialResult represents the result of a trial execution
type TrialResult struct {
	Success  bool
	Message  string
	Changes  []string
	Duration time.Duration
}

// NewTrial creates a new trial session
func NewTrial(projectPath, taskDescription string) (*Trial, error) {
	originalBranch, err := getCurrentBranch(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get current branch: %w", err)
	}

	timestamp := time.Now().Format("20060102-150405")
	branchName := fmt.Sprintf("sdp-try-%s", timestamp)

	return &Trial{
		BranchName:      branchName,
		OriginalBranch:  originalBranch,
		ProjectPath:     projectPath,
		TaskDescription: taskDescription,
		StartTime:       time.Now(),
	}, nil
}

// Start creates the temporary branch and sets up the trial
func (t *Trial) Start() error {
	if err := t.createBranch(); err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}
	return nil
}

// Execute runs the task as a dry-run planner
// It parses the task description and returns a structured execution plan
// without making actual changes to the codebase
func (t *Trial) Execute() (*TrialResult, error) {
	// Validate task description
	if t.TaskDescription == "" {
		return &TrialResult{
			Success:  false,
			Message:  "Task description cannot be empty",
			Changes:  []string{},
			Duration: time.Since(t.StartTime),
		}, nil
	}

	// Analyze task description and create execution plan
	plan := t.createExecutionPlan()

	result := &TrialResult{
		Success: true,
		Message: fmt.Sprintf("Dry-run plan created for: %s\n\n%s", t.TaskDescription, plan),
		Changes: []string{
			fmt.Sprintf("Branch: %s", t.BranchName),
			fmt.Sprintf("Task: %s", t.TaskDescription),
			fmt.Sprintf("Plan:\n%s", plan),
			fmt.Sprintf("Duration: %v", time.Since(t.StartTime).Round(time.Millisecond)),
		},
		Duration: time.Since(t.StartTime),
	}

	return result, nil
}

// createExecutionPlan analyzes the task and creates a structured execution plan
func (t *Trial) createExecutionPlan() string {
	// Detect task type from description
	taskDesc := strings.ToLower(t.TaskDescription)

	var taskType string
	var steps []string

	// Simple pattern matching for common task types
	switch {
	case strings.Contains(taskDesc, "add") || strings.Contains(taskDesc, "create") || strings.Contains(taskDesc, "implement"):
		taskType = "Feature Addition"
		steps = []string{
			"1. Analyze existing codebase structure",
			"2. Identify relevant files and dependencies",
			"3. Create/modify implementation files",
			"4. Add/update tests",
			"5. Update documentation",
		}
	case strings.Contains(taskDesc, "fix") || strings.Contains(taskDesc, "bug"):
		taskType = "Bug Fix"
		steps = []string{
			"1. Reproduce and identify the issue",
			"2. Locate problematic code",
			"3. Implement fix",
			"4. Add regression tests",
			"5. Verify fix resolves issue",
		}
	case strings.Contains(taskDesc, "refactor") || strings.Contains(taskDesc, "clean"):
		taskType = "Refactoring"
		steps = []string{
			"1. Analyze current implementation",
			"2. Identify refactoring opportunities",
			"3. Apply refactoring changes",
			"4. Ensure tests pass",
			"5. Update documentation if needed",
		}
	case strings.Contains(taskDesc, "test"):
		taskType = "Test Addition"
		steps = []string{
			"1. Identify untested code paths",
			"2. Design test cases",
			"3. Implement tests",
			"4. Verify coverage",
			"5. Document test scenarios",
		}
	default:
		taskType = "General Task"
		steps = []string{
			"1. Understand requirements",
			"2. Analyze affected components",
			"3. Implement changes",
			"4. Test and verify",
			"5. Document changes",
		}
	}

	// Build plan string
	var plan strings.Builder
	plan.WriteString(fmt.Sprintf("Task Type: %s\n", taskType))
	plan.WriteString("Proposed Execution Steps:\n")
	for _, step := range steps {
		plan.WriteString(fmt.Sprintf("  %s\n", step))
	}
	plan.WriteString("\nNote: This is a dry-run plan. No actual changes have been made.")
	plan.WriteString("\nUse 'sdp adopt' to convert this trial into a full SDP setup.")

	return plan.String()
}

// Accept keeps the branch and suggests adoption
func (t *Trial) Accept() error {
	fmt.Printf("✓ Trial accepted. Branch '%s' kept for adoption.\n", t.BranchName)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review the changes: git diff")
	fmt.Println("  2. Run 'sdp adopt' to convert to full SDP setup")
	fmt.Printf("  3. Or merge manually: git checkout <target-branch> && git merge %s\n", t.BranchName)
	return nil
}

// Discard deletes the branch and returns to original state
func (t *Trial) Discard() error {
	if err := checkoutBranch(t.ProjectPath, t.OriginalBranch); err != nil {
		return fmt.Errorf("failed to checkout original branch: %w", err)
	}

	if err := deleteBranch(t.ProjectPath, t.BranchName); err != nil {
		return fmt.Errorf("failed to delete trial branch: %w", err)
	}

	fmt.Printf("✓ Trial discarded. Repository restored to original state.\n")
	return nil
}

// VerifyClean checks if repository is clean (no uncommitted changes)
func (t *Trial) VerifyClean() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = t.ProjectPath
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %w", err)
	}

	return len(strings.TrimSpace(string(output))) == 0, nil
}

// getCurrentBranch returns the current git branch name
func getCurrentBranch(projectPath string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = projectPath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// createBranch creates and checks out a new branch
func (t *Trial) createBranch() error {
	cmd := exec.Command("git", "checkout", "-b", t.BranchName)
	cmd.Dir = t.ProjectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create branch: %s", string(output))
	}
	return nil
}

// checkoutBranch checks out the specified branch
func checkoutBranch(projectPath, branchName string) error {
	cmd := exec.Command("git", "checkout", branchName)
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to checkout branch: %s", string(output))
	}
	return nil
}

// deleteBranch deletes the specified branch
func deleteBranch(projectPath, branchName string) error {
	cmd := exec.Command("git", "branch", "-D", branchName)
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete branch: %s", string(output))
	}
	return nil
}
