package main

import (
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/checkpoint"
	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/fall-out-bug/sdp/internal/orchestrator"
	"github.com/spf13/cobra"
)

var orchestrateCmd = &cobra.Command{
	Use:   "orchestrate <feature-id>",
	Short: "Orchestrate workstream execution for a feature",
	Long: `Execute all workstreams for a feature in dependency order.

This command:
1. Loads workstreams from docs/workstreams/backlog/
2. Builds dependency graph
3. Executes workstreams in topological order
4. Creates checkpoints after each workstream
5. Supports resume via 'sdp checkpoint resume'

Example:
  sdp orchestrate F050    # Execute all workstreams for F050`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		featureID := args[0]

		// Get checkpoint directory
		dir, err := cmd.Flags().GetString("checkpoint-dir")
		if err != nil {
			return fmt.Errorf("failed to get checkpoint-dir flag: %w", err)
		}

		if dir == "" {
			dir, err = checkpoint.GetDefaultDir()
			if err != nil {
				return fmt.Errorf("failed to get default checkpoint directory: %w", err)
			}
		}

		// Get workstream directory
		wsDir, err := cmd.Flags().GetString("workstream-dir")
		if err != nil {
			return fmt.Errorf("failed to get workstream-dir flag: %w", err)
		}

		if wsDir == "" {
			wsDir = "docs/workstreams/backlog"
		}

		// Get retry count
		maxRetries, err := cmd.Flags().GetInt("retry")
		if err != nil {
			return fmt.Errorf("failed to get retry flag: %w", err)
		}

		// Create components
		loader := orchestrator.NewBeadsLoader(wsDir, ".beads-sdp-mapping.jsonl")
		executor := orchestrator.NewCLIExecutor("sdp")
		checkpointMgr := checkpoint.NewManager(dir)

		// Create orchestrator
		orch := orchestrator.NewOrchestrator(loader, executor, checkpointMgr, maxRetries)

		// Run orchestration
		fmt.Printf("🚀 Orchestrating feature %s\n", featureID)
		fmt.Printf("   Workstream dir: %s\n", wsDir)
		fmt.Printf("   Checkpoint dir: %s\n", dir)
		fmt.Printf("   Max retries: %d\n\n", maxRetries)

		// F056-03: plan event (non-blocking)
		evidence.Emit(evidence.SkillEvent("oneshot", "plan", "00-000-00", map[string]interface{}{"feature_id": featureID}))

		err = orch.Run(featureID)
		if err != nil {
			return fmt.Errorf("orchestration failed: %w", err)
		}

		if err := evidence.EmitSync(evidence.SkillEvent("oneshot", "approval", "00-000-00", map[string]interface{}{"feature_id": featureID})); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "warning: evidence emit: %v\n", err)
		}

		fmt.Printf("\n✅ Feature %s completed successfully\n", featureID)
		return nil
	},
}

var orchestrateResumeCmd = &cobra.Command{
	Use:   "resume <checkpoint-id>",
	Short: "Resume orchestration from a checkpoint",
	Long: `Resume workstream execution from a saved checkpoint.

Example:
  sdp orchestrate resume F050    # Resume execution of F050`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		checkpointID := args[0]

		// Get checkpoint directory
		dir, err := cmd.Flags().GetString("checkpoint-dir")
		if err != nil {
			return fmt.Errorf("failed to get checkpoint-dir flag: %w", err)
		}

		if dir == "" {
			dir, err = checkpoint.GetDefaultDir()
			if err != nil {
				return fmt.Errorf("failed to get default checkpoint directory: %w", err)
			}
		}

		// Get workstream directory
		wsDir, err := cmd.Flags().GetString("workstream-dir")
		if err != nil {
			return fmt.Errorf("failed to get workstream-dir flag: %w", err)
		}

		if wsDir == "" {
			wsDir = "docs/workstreams/backlog"
		}

		// Get retry count
		maxRetries, err := cmd.Flags().GetInt("retry")
		if err != nil {
			return fmt.Errorf("failed to get retry flag: %w", err)
		}

		// Create components
		loader := orchestrator.NewBeadsLoader(wsDir, ".beads-sdp-mapping.jsonl")
		executor := orchestrator.NewCLIExecutor("sdp")
		checkpointMgr := checkpoint.NewManager(dir)

		// Create orchestrator
		orch := orchestrator.NewOrchestrator(loader, executor, checkpointMgr, maxRetries)

		// Resume orchestration
		fmt.Printf("🔄 Resuming from checkpoint %s\n", checkpointID)
		fmt.Printf("   Workstream dir: %s\n", wsDir)
		fmt.Printf("   Checkpoint dir: %s\n\n", dir)

		err = orch.Resume(checkpointID)
		if err != nil {
			return fmt.Errorf("resume failed: %w", err)
		}

		fmt.Printf("\n✅ Checkpoint %s completed successfully\n", checkpointID)
		return nil
	},
}

func init() {
	// Add persistent flags to parent orchestrate command
	orchestrateCmd.PersistentFlags().String("checkpoint-dir", "", "Checkpoint directory (default: .sdp/checkpoints)")
	orchestrateCmd.PersistentFlags().String("workstream-dir", "docs/workstreams/backlog", "Workstream directory")
	orchestrateCmd.PersistentFlags().Int("retry", 2, "Max retries per workstream (default: 2)")

	// Add resume subcommand
	orchestrateCmd.AddCommand(orchestrateResumeCmd)
}
