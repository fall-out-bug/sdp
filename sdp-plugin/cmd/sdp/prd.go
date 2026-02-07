package main

import (
	"fmt"

	"github.com/fall-out-bug/sdp/internal/prd"
	"github.com/spf13/cobra"
)

func prdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prd",
		Short: "PRD (Product Requirements Document) operations",
		Long: `PRD operations for project analysis and validation.

Subcommands:
  detect-type - Detect project type from file structure`,
	}

	cmd.AddCommand(prdDetectType())

	return cmd
}

func prdDetectType() *cobra.Command {
	var projectPath string

	cmd := &cobra.Command{
		Use:   "detect-type [project-path]",
		Short: "Detect project type from file structure",
		Long: `Detect project type by analyzing file structure and configuration.

Detection strategy:
  1. go.mod → Go project
  2. pyproject.toml/setup.py → Python project
  3. pom.xml/build.gradle → Java project
  4. docker-compose.yml → Service
  5. Default → Library

Project types:
  - service: Web service/API
  - cli: Command-line tool
  - library: Code library/package
  - go/python/java: Language-specific project`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Use provided path or current directory
			if len(args) > 0 {
				projectPath = args[0]
			}
			if projectPath == "" {
				projectPath = "."
			}

			// Create detector
			detector := prd.NewDetector(projectPath)

			// Detect type
			projectType := detector.DetectType()

			// Print result
			fmt.Printf("Detected project type: %s\n", projectType.String())

			// Provide additional context
			switch projectType {
			case prd.Service:
				fmt.Println("  → Web service or API")
			case prd.CLI:
				fmt.Println("  → Command-line interface tool")
			case prd.Library:
				fmt.Println("  → Code library or package")
			case prd.Go:
				fmt.Println("  → Go module")
			case prd.Python:
				fmt.Println("  → Python project")
			case prd.Java:
				fmt.Println("  → Java project")
			default:
				fmt.Println("  → Unknown project type")
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&projectPath, "path", "", "Path to project root (default: current directory)")

	return cmd
}
