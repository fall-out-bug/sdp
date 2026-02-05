package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ai-masters/sdp/internal/tdd"
	"github.com/spf13/cobra"
)

func tddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tdd <phase> [path]",
		Short: "Run TDD cycle (Red-Green-Refactor)",
		Long: `Execute TDD cycle phases with automated test running.

Phases:
  red      - Run tests, expect failure (write failing test)
  green    - Run tests, expect success (make test pass)
  refactor - Run tests, ensure no regression (improve code)

Arguments:
  phase   - TDD phase to run (red, green, refactor)
  path    - Package path to test (default: ./internal/parser)

The command automatically detects the project language (Go, Python, Java)
and runs the appropriate test runner.

Examples:
  sdp tdd green ./internal/parser
  sdp tdd red ./myapp
  sdp tdd refactor`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			phase := args[0]

			// Get path, default to parser package
			path := "./internal/parser"
			if len(args) == 2 {
				path = args[1]
			}

			// Validate phase
			var tddPhase tdd.Phase
			switch phase {
			case "red":
				tddPhase = tdd.Red
			case "green":
				tddPhase = tdd.Green
			case "refactor":
				tddPhase = tdd.Refactor
			default:
				return fmt.Errorf("invalid phase: %s (must be red, green, or refactor)", phase)
			}

			// Detect language and create runner
			runner := tdd.NewRunner(tdd.Go) // Default to Go, will be overridden by auto-detection

			// Create context with cancellation
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Handle Ctrl+C
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
			go func() {
				<-sigChan
				fmt.Println("\n🛑 TDD cycle interrupted")
				cancel()
			}()

			// Run the phase
			result, err := runner.RunPhase(ctx, tddPhase, path)

			// Print results
			fmt.Printf("\n📊 TDD Phase: %s\n", result.Phase)
			fmt.Printf("⏱️  Duration: %v\n", result.Duration)
			if result.Success {
				fmt.Printf("✅ Success\n")
			} else {
				fmt.Printf("❌ Failed\n")
			}

			if result.Stdout != "" {
				fmt.Printf("\n📤 Output:\n%s\n", result.Stdout)
			}

			if result.Stderr != "" {
				fmt.Printf("\n⚠️  Errors:\n%s\n", result.Stderr)
			}

			if err != nil {
				return fmt.Errorf("phase failed: %w", err)
			}

			return nil
		},
	}

	return cmd
}
