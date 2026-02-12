package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/task"
	"github.com/spf13/cobra"
)

// taskCmd represents the task command
func taskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Task management commands",
		Long:  `Create and manage tasks (bugs, features, hotfixes).`,
	}

	cmd.AddCommand(taskCreateCmd())

	return cmd
}

// taskCreateCmd represents the task create command
func taskCreateCmd() *cobra.Command {
	var (
		taskType      string
		title         string
		priority      int
		featureID     string
		goal          string
		context       string
		branchBase    string
		outputJSON    bool
		createIssue   bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new task (bug, task, or hotfix)",
		Long: `Create a new task with dual-track artifact creation.

If a feature ID is provided, creates a workstream file.
If no feature ID (or --issue flag), creates an issue file.

Types:
  - bug: P1/P2 quality fixes (branches from dev/feature)
  - task: Feature workstreams (branches from dev/feature)
  - hotfix: P0 production fixes (branches from main)

Examples:
  # Create a bug workstream
  sdp task create --type=bug --title="Fix CI" --feature=F064 --priority=1

  # Create a standalone issue
  sdp task create --type=bug --title="Auth error" --issue

  # Create a hotfix
  sdp task create --type=hotfix --title="DB connection" --feature=F064 --priority=0 --branch=main`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate type
			t, err := parseTaskType(taskType)
			if err != nil {
				return err
			}

			// Validate priority
			if priority < 0 || priority > 3 {
				return fmt.Errorf("priority must be 0-3")
			}

			// Create task definition
			taskDef := &task.Task{
				Type:       t,
				Title:      title,
				Priority:   task.Priority(priority),
				FeatureID:  featureID,
				Goal:       goal,
				Context:    context,
				BranchBase: branchBase,
			}

			// Determine creation mode
			creator := task.NewCreator(task.CreatorConfig{
				WorkstreamDir: "docs/workstreams/backlog",
				IssuesDir:     "docs/issues",
				IndexFile:     ".sdp/issues-index.jsonl",
			})

			if createIssue || featureID == "" {
				// Create issue (Beads-free fallback)
				issue, err := creator.CreateIssue(taskDef)
				if err != nil {
					return fmt.Errorf("failed to create issue: %w", err)
				}

				if outputJSON {
					encoder := json.NewEncoder(os.Stdout)
					encoder.SetIndent("", "  ")
					return encoder.Encode(map[string]interface{}{
						"type":     "issue",
						"issue_id": issue.IssueID,
						"path":     issue.Path,
					})
				}

				fmt.Printf("Created issue: %s\n", issue.IssueID)
				fmt.Printf("Path: %s\n", issue.Path)
				return nil
			}

			// Create workstream
			ws, err := creator.CreateWorkstream(taskDef)
			if err != nil {
				return fmt.Errorf("failed to create workstream: %w", err)
			}

			if outputJSON {
				encoder := json.NewEncoder(os.Stdout)
				encoder.SetIndent("", "  ")
				return encoder.Encode(map[string]interface{}{
					"type":       "workstream",
					"ws_id":      ws.WSID,
					"feature_id": ws.FeatureID,
					"path":       ws.Path,
				})
			}

			fmt.Printf("Created workstream: %s\n", ws.WSID)
			fmt.Printf("Feature: %s\n", ws.FeatureID)
			fmt.Printf("Path: %s\n", ws.Path)
			return nil
		},
	}

	cmd.Flags().StringVar(&taskType, "type", "task", "Task type: bug, task, hotfix")
	cmd.Flags().StringVarP(&title, "title", "t", "", "Task title (required)")
	cmd.Flags().IntVarP(&priority, "priority", "p", 2, "Priority: 0=P0, 1=P1, 2=P2, 3=P3")
	cmd.Flags().StringVarP(&featureID, "feature", "f", "", "Parent feature ID (e.g., F064)")
	cmd.Flags().StringVar(&goal, "goal", "", "Goal description")
	cmd.Flags().StringVar(&context, "context", "", "Context/symptom description")
	cmd.Flags().StringVar(&branchBase, "branch", "", "Base branch (main for hotfix, dev otherwise)")
	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&createIssue, "issue", false, "Create issue file instead of workstream")

	_ = cmd.MarkFlagRequired("title")

	return cmd
}

// parseTaskType parses a task type string
func parseTaskType(s string) (task.Type, error) {
	switch s {
	case "bug":
		return task.TypeBug, nil
	case "task":
		return task.TypeTask, nil
	case "hotfix":
		return task.TypeHotfix, nil
	default:
		return "", fmt.Errorf("invalid task type: %s (valid: bug, task, hotfix)", s)
	}
}
