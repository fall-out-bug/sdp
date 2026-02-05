package main

import (
	"fmt"

	"github.com/ai-masters/sdp/internal/beads"
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
  sync      Synchronize Beads state`,
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
				return err
			}

			tasks, err := client.Ready()
			if err != nil {
				return err
			}

			if len(tasks) == 0 {
				fmt.Println("No available tasks")
				return nil
			}

			fmt.Printf("Found %d available task(s):\n\n", len(tasks))
			for _, task := range tasks {
				fmt.Printf("  • %s %s\n", task.ID, task.Title)
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
				return err
			}

			task, err := client.Show(args[0])
			if err != nil {
				return err
			}

			fmt.Printf("Task: %s\n", task.ID)
			fmt.Printf("  Title: %s\n", task.Title)
			fmt.Printf("  Status: %s\n", task.Status)
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
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if status == "" {
				return fmt.Errorf("--status flag is required")
			}

			client, err := beads.NewClient()
			if err != nil {
				return err
			}

			if err := client.Update(args[0], status); err != nil {
				return err
			}

			fmt.Printf("✓ Updated task %s to status: %s\n", args[0], status)

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
				return err
			}

			if err := client.Sync(); err != nil {
				return err
			}

			fmt.Println("✓ Beads synchronized")

			return nil
		},
	}
}
