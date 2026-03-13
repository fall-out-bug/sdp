package main

import (
	"fmt"

	"github.com/fall-out-bug/sdp/internal/beads"
	"github.com/fall-out-bug/sdp/internal/ui"
	"github.com/spf13/cobra"
)

func beadsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "beads",
		Short: "Interact with Beads task tracker",
		Long: `Interact with Beads CLI for task tracking.

Commands:
  ready     List available tasks
  show      Show task details
  create    Create a task and persist repo snapshot
  update    Update task status
  close     Close a task and persist repo snapshot
  sync      Export Beads state back to repo snapshot

Examples:
  sdp beads ready
  sdp beads create --title "WS 00-001-01: Add parser" --type task
  sdp beads show sdp-abc
  sdp beads update sdp-abc --status in_progress
  sdp beads close sdp-abc --reason "WS completed"
  sdp beads sync`,
	}

	cmd.AddCommand(beadsReadyCmd())
	cmd.AddCommand(beadsCreateCmd())
	cmd.AddCommand(beadsShowCmd())
	cmd.AddCommand(beadsUpdateCmd())
	cmd.AddCommand(beadsCloseCmd())
	cmd.AddCommand(beadsSyncCmd())

	return cmd
}

func beadsReadyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ready",
		Short: "List available tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := beads.NewClient()
			if err != nil {
				return fmt.Errorf("failed to create beads client: %w", err)
			}

			tasks, err := client.Ready()
			if err != nil {
				return fmt.Errorf("failed to list tasks: %w", err)
			}

			if len(tasks) == 0 {
				ui.InfoLine("No available tasks")
				return nil
			}

			ui.Header(fmt.Sprintf("Found %d available task(s)", len(tasks)))
			for _, task := range tasks {
				fmt.Printf("  • %s %s\n", ui.BoldText(task.ID), task.Title)
			}

			return nil
		},
	}
}

func beadsShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <beads-id>",
		Short: "Show task details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := beads.NewClient()
			if err != nil {
				return fmt.Errorf("failed to create beads client: %w", err)
			}

			task, err := client.Show(args[0])
			if err != nil {
				return fmt.Errorf("failed to show task: %w", err)
			}

			ui.Subheader("Task Details")
			fmt.Printf("  ID:       %s\n", ui.BoldText(task.ID))
			fmt.Printf("  Title:    %s\n", task.Title)
			fmt.Printf("  Status:   %s\n", ui.Info(task.Status))
			fmt.Printf("  Priority: %s\n", task.Priority)

			return nil
		},
	}
}

func beadsCreateCmd() *cobra.Command {
	var (
		title       string
		issueType   string
		priority    string
		labels      []string
		description string
		parent      string
		silent      bool
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a task and persist repo snapshot",
		Example: `  sdp beads create --title "WS 00-001-01: Add parser" --type task
  sdp beads create --title "Security: missing auth check" --type bug --priority 0 --labels review-finding,F051,security --silent`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if title == "" {
				return fmt.Errorf("%s: --title flag is required", ui.Error("Error"))
			}

			client, err := beads.NewClient()
			if err != nil {
				return fmt.Errorf("failed to create beads client: %w", err)
			}

			beadsID, err := client.Create(title, beads.CreateOptions{
				Type:        issueType,
				Priority:    priority,
				Labels:      labels,
				Description: description,
				Parent:      parent,
				Silent:      silent,
			})
			if err != nil {
				return fmt.Errorf("failed to create task: %w", err)
			}
			if err := client.Sync(); err != nil {
				return fmt.Errorf("created task %s but failed to persist snapshot: %w", beadsID, err)
			}

			ui.SuccessLine("Created task %s", ui.BoldText(beadsID))
			return nil
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Task title (required)")
	cmd.Flags().StringVar(&issueType, "type", "task", "Task type (task, bug, feature, ...)")
	cmd.Flags().StringVar(&priority, "priority", "", "Priority value understood by Beads")
	cmd.Flags().StringSliceVar(&labels, "labels", nil, "Comma-separated labels")
	cmd.Flags().StringVar(&description, "description", "", "Task description")
	cmd.Flags().StringVar(&parent, "parent", "", "Parent Beads issue ID")
	cmd.Flags().BoolVar(&silent, "silent", false, "Pass --silent to bd create")

	return cmd
}

func beadsUpdateCmd() *cobra.Command {
	var status string

	cmd := &cobra.Command{
		Use:   "update <beads-id>",
		Short: "Update task status",
		Long: `Update task status in Beads.

Valid statuses:
  in_progress  Task is currently being worked on
  completed    Task is finished
  blocked      Task is blocked and cannot proceed`,
		Example: `  sdp beads update sdp-abc --status in_progress
  sdp beads update sdp-abc --status completed`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if status == "" {
				return fmt.Errorf("%s: --status flag is required (valid: in_progress, completed, blocked)", ui.Error("Error"))
			}

			client, err := beads.NewClient()
			if err != nil {
				return fmt.Errorf("failed to create beads client: %w", err)
			}

			if err := client.Update(args[0], status); err != nil {
				return fmt.Errorf("failed to update task: %w", err)
			}
			if err := client.Sync(); err != nil {
				return fmt.Errorf("updated task %s but failed to persist snapshot: %w", args[0], err)
			}

			ui.SuccessLine("Updated task %s to status: %s", args[0], ui.BoldText(status))

			return nil
		},
	}

	cmd.Flags().StringVar(&status, "status", "", "New status (in_progress, completed, blocked)")

	return cmd
}

func beadsCloseCmd() *cobra.Command {
	var reason string

	cmd := &cobra.Command{
		Use:   "close <beads-id>",
		Short: "Close a task and persist repo snapshot",
		Args:  cobra.ExactArgs(1),
		Example: `  sdp beads close sdp-abc --reason "WS completed"
  sdp beads close sdp-def --reason "Duplicate finding"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := beads.NewClient()
			if err != nil {
				return fmt.Errorf("failed to create beads client: %w", err)
			}

			if err := client.Close(args[0], reason); err != nil {
				return fmt.Errorf("failed to close task: %w", err)
			}
			if err := client.Sync(); err != nil {
				return fmt.Errorf("closed task %s but failed to persist snapshot: %w", args[0], err)
			}

			ui.SuccessLine("Closed task %s", ui.BoldText(args[0]))
			return nil
		},
	}

	cmd.Flags().StringVar(&reason, "reason", "", "Close reason")

	return cmd
}

func beadsSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Export Beads state back to repo snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := beads.NewClient()
			if err != nil {
				return fmt.Errorf("failed to create beads client: %w", err)
			}

			if err := client.Sync(); err != nil {
				return fmt.Errorf("failed to synchronize: %w", err)
			}

			ui.SuccessLine("Beads snapshot exported")

			return nil
		},
	}
}
