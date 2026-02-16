package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// demoCmd returns the demo command
func demoCmd() *cobra.Command {
	var template string
	var cleanup bool
	var verbose bool

	cmd := &cobra.Command{
		Use:   "demo",
		Short: "Run interactive demo of SDP workflow",
		Long: `Run an interactive demonstration of the SDP workflow.

This command creates a temporary project and walks through:
1. Project initialization (sdp init --guided)
2. Environment verification (sdp doctor)
3. Status check (sdp status --text)
4. Cleanup (optional)

The demo is designed to give new users a "first success" experience
in under 5 minutes.`,
		Example: `  # Run demo with default template
  sdp demo

  # Run with verbose output
  sdp demo --verbose

  # Keep demo project for exploration
  sdp demo --no-cleanup

  # Use specific template
  sdp demo --template minimal-go`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDemo(template, cleanup, verbose)
		},
	}

	cmd.Flags().StringVar(&template, "template", "minimal-go", "Template to use")
	cmd.Flags().BoolVar(&cleanup, "cleanup", true, "Clean up demo project after completion")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "Show detailed output")

	return cmd
}

// DemoStep represents a step in the demo workflow
type DemoStep struct {
	Name        string
	Description string
	Command     string
	Expected    string
	Optional    bool
}

// runDemo executes the demo workflow
func runDemo(template string, shouldCleanup bool, verbose bool) error {
	fmt.Println("SDP Interactive Demo")
	fmt.Println("===================")
	fmt.Println()
	fmt.Println("This demo will walk you through a complete SDP workflow.")
	fmt.Println("Estimated time: 2-5 minutes")
	fmt.Println()

	templateDir := filepath.Join("templates", template)
	if _, err := os.Stat(templateDir); os.IsNotExist(err) {
		return fmt.Errorf("template not found: %s", templateDir)
	}

	demoDir, err := os.MkdirTemp("", "sdp-demo-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}

	if shouldCleanup {
		defer os.RemoveAll(demoDir)
	}

	fmt.Printf("Demo directory: %s\n", demoDir)
	fmt.Println()

	fmt.Println("Step 0: Setting up demo project...")
	if err := copyTemplate(templateDir, demoDir); err != nil {
		return fmt.Errorf("copy template: %w", err)
	}
	fmt.Println("  [OK] Template copied")
	fmt.Println()

	originalWd, _ := os.Getwd()
	if err := os.Chdir(demoDir); err != nil {
		return fmt.Errorf("change directory: %w", err)
	}
	defer os.Chdir(originalWd)

	steps := []DemoStep{
		{
			Name:        "Initialize Git",
			Description: "Initialize Git repository for version control",
			Command:     "git init",
			Expected:    "Initialized empty Git repository",
		},
		{
			Name:        "SDP Guided Setup",
			Description: "Run guided SDP initialization",
			Command:     "sdp init --guided --auto-fix",
			Expected:    "Setup complete",
		},
		{
			Name:        "Verify Environment",
			Description: "Check that all prerequisites are met",
			Command:     "sdp doctor",
			Expected:    "All required checks passed",
		},
		{
			Name:        "Check Status",
			Description: "View project status",
			Command:     "sdp status --text",
			Expected:    "Next Action:",
		},
		{
			Name:        "Run Tests",
			Description: "Execute project tests",
			Command:     "go test ./...",
			Expected:    "PASS",
		},
	}

	for i, step := range steps {
		if !executeDemoStep(i+1, step, verbose) {
			if !step.Optional {
				return fmt.Errorf("step %d failed", i+1)
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	printDemoSummary(shouldCleanup, demoDir)
	return nil
}

// executeDemoStep runs a single demo step and returns success
func executeDemoStep(stepNum int, step DemoStep, verbose bool) bool {
	fmt.Printf("Step %d: %s\n", stepNum, step.Name)
	fmt.Printf("  %s\n", step.Description)
	fmt.Printf("  Running: %s\n", step.Command)

	output, err := runDemoCommand(step.Command, verbose)
	if err != nil && !step.Optional {
		fmt.Printf("  [FAIL] %v\n", err)
		if verbose {
			fmt.Printf("  Output: %s\n", output)
		}
		return false
	}

	if strings.Contains(output, step.Expected) || step.Optional {
		fmt.Printf("  [OK] %s\n", step.Name)
	} else {
		fmt.Printf("  [WARN] Expected '%s' in output\n", step.Expected)
	}

	if verbose {
		fmt.Printf("  Output:\n%s\n", indent(output, "    "))
	}

	fmt.Println()
	return true
}

// printDemoSummary prints the final demo summary
func printDemoSummary(shouldCleanup bool, demoDir string) {
	fmt.Println("Demo Complete!")
	fmt.Println("==============")
	fmt.Println()
	fmt.Println("You have successfully:")
	fmt.Println("  1. Set up an SDP-compatible project")
	fmt.Println("  2. Initialized SDP with guided setup")
	fmt.Println("  3. Verified your environment")
	fmt.Println("  4. Checked project status")
	fmt.Println("  5. Run project tests")
	fmt.Println()

	if shouldCleanup {
		fmt.Println("Demo project has been cleaned up.")
	} else {
		fmt.Printf("Demo project preserved at: %s\n", demoDir)
		fmt.Println("You can explore it or delete it manually.")
	}

	fmt.Println()
	fmt.Println("Next Steps:")
	fmt.Println("  - Try 'sdp plan \"your feature idea\"' in your own project")
	fmt.Println("  - Run 'sdp demo --help' for more options")
}
