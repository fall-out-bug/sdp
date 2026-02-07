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
  update    Update task status
  sync      Synchronize Beads state

Examples:
  sdp beads ready
  sdp beads show sdp-abc
  sdp beads update sdp-abc --status in_progress
  sdp beads sync`,
	}

	cmd.AddCommand(beadsReadyCmd())
	cmd.AddCommand(beadsShowCmd())
	cmd.AddCommand(beadsUpdateCmd())
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

			ui.SuccessLine("Updated task %s to status: %s", args[0], ui.BoldText(status))

			return nil
		},
	}

	cmd.Flags().StringVar(&status, "status", "", "New status (in_progress, completed, blocked)")

	return cmd
}

func beadsSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Synchronize Beads state",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := beads.NewClient()
			if err != nil {
				return fmt.Errorf("failed to create beads client: %w", err)
			}

			if err := client.Sync(); err != nil {
				return fmt.Errorf("failed to synchronize: %w", err)
			}

			ui.SuccessLine("Beads synchronized")

			return nil
		},
	}
}
