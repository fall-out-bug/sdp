package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/fall-out-bug/sdp/internal/planner"
	"github.com/spf13/cobra"
)

const defaultBacklogDir = "docs/workstreams/backlog"

func planCmd() *cobra.Command {
	var interactive, autoApply, dryRun bool
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "plan <description>",
		Short: "Decompose feature into workstreams",
		Long: `Decompose a feature description into workstreams from the terminal.

Modes:
  - Default (drive mode): Shows decomposition, waits for confirmation
  - --interactive: Asks questions to refine requirements
  - --auto-apply: Plans then executes immediately (ship mode)
  - --dry-run: Shows what would be created without writing files

Output:
  - Default: Human-readable table of workstreams and dependencies
  - --output=json: Machine-readable JSON format

Examples:
  sdp plan "Add OAuth2"
  sdp plan "Add OAuth2" --interactive
  sdp plan "Add OAuth2" --auto-apply
  sdp plan "Add OAuth2" --dry-run
  sdp plan "Add OAuth2" --output=json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			description := args[0]

			// Find project root
			root, err := config.FindProjectRoot()
			if err != nil {
				return fmt.Errorf("find project root: %w", err)
			}

			// Setup paths
			backlogDir := filepath.Join(root, defaultBacklogDir)
			logPath := filepath.Join(root, ".sdp/log/events.jsonl")

			// Create evidence writer
			evWriter, err := evidence.NewWriter(logPath)
			if err != nil {
				return fmt.Errorf("create evidence writer: %w", err)
			}

			// Get model API from environment (AC8)
			modelAPI := os.Getenv("MODEL_API")
			// Don't set a default - require explicit configuration
			// This ensures users are aware they need to configure the model

			// Create planner
			p := &planner.Planner{
				BacklogDir:     backlogDir,
				Description:    description,
				Interactive:    interactive,
				AutoApply:      autoApply,
				DryRun:         dryRun,
				OutputFormat:   outputFormat,
				ModelAPI:       modelAPI,
				EvidenceWriter: evWriter,
			}

			// Run interactive questions if in interactive mode
			if p.Interactive {
				fmt.Println("Interactive mode: Answer questions to refine the plan")
				if err := p.PromptForInteractive(); err != nil {
					return fmt.Errorf("interactive prompt failed: %w", err)
				}
			}

			// Perform decomposition
			fmt.Printf("Decomposing: %s\n", description)
			result, err := p.Decompose()
			if err != nil {
				// AC8: Clear error message when no model configured
				return fmt.Errorf("decomposition failed: %w\n\nHint: Set MODEL_API environment variable or configure model endpoint in .sdp/config.json", err)
			}

			// Output the plan
			output, err := p.FormatOutput(result)
			if err != nil {
				return fmt.Errorf("format output: %w", err)
			}
			fmt.Println(output)

			// AC7: Dry-run mode - show what would be created
			if dryRun {
				fmt.Println("\n[DRY RUN] Would create the following files:")
				for _, ws := range result.Workstreams {
					filename := ws.Filename()
					fmt.Printf("  - %s\n", filepath.Join(backlogDir, filename))
				}
				return nil
			}

			// Create workstream files (AC4)
			fmt.Printf("\nCreating %d workstream files...\n", len(result.Workstreams))
			if err := p.CreateWorkstreamFiles(result); err != nil {
				return fmt.Errorf("create workstream files: %w", err)
			}

			// Emit plan event (AC5)
			if err := p.EmitPlanEvent(result); err != nil {
				// Non-fatal: log warning but continue
				fmt.Fprintf(os.Stderr, "warning: failed to emit plan event: %v\n", err)
			}

			// AC3: Auto-apply mode - trigger execution
			if autoApply {
				fmt.Println("\nAuto-apply mode: Triggering execution...")
				if err := p.ExecuteAutoApply(result); err != nil {
					return fmt.Errorf("auto-apply failed: %w", err)
				}
				fmt.Printf("Execution started for %d workstreams\n", len(result.Workstreams))
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&interactive, "interactive", false, "Ask questions to refine requirements (drive mode)")
	cmd.Flags().BoolVar(&autoApply, "auto-apply", false, "Execute plan after creation (ship mode)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be created without writing files")
	cmd.Flags().StringVar(&outputFormat, "output", "human", "Output format: human or json")

	return cmd
}
