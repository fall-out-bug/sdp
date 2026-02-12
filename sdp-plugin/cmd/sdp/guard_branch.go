package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/session"
	"github.com/spf13/cobra"
)

// guardBranchCmd returns the guard branch command group
func guardBranchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "branch",
		Short: "Branch validation commands",
		Long: `Validate branch naming and feature branch requirements.

These commands ensure that feature work is done on feature branches,
not on protected branches like dev or main.`,
	}
	cmd.AddCommand(guardBranchCheckCmd())
	cmd.AddCommand(guardBranchValidateCmd())
	return cmd
}

// guardBranchCheckCmd checks if current branch is valid for a feature
func guardBranchCheckCmd() *cobra.Command {
	var featureID string

	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check if current branch is valid for feature",
		Long: `Check that the current branch is valid for the specified feature.

Exit codes:
  0 - Branch is valid
  1 - On protected branch (dev/main)
  2 - Wrong feature branch
  3 - Not a feature branch`,
		Example: `  # Check if on correct branch for F065
  sdp guard branch check --feature=F065

  # Use exit code for scripting
  sdp guard branch check --feature=F065 && echo "OK"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := config.FindProjectRoot()
			if err != nil {
				return fmt.Errorf("find project root: %w", err)
			}

			// If no feature ID specified, try to get from session
			if featureID == "" {
				if session.Exists(root) {
					s, err := session.Load(root)
					if err == nil {
						featureID = s.FeatureID
					}
				}
			}

			if featureID == "" {
				return fmt.Errorf("--feature flag is required (or active session)")
			}

			// Get current branch
			currentBranch, err := getCurrentBranch()
			if err != nil {
				return fmt.Errorf("get current branch: %w", err)
			}

			// Check protected branches
			if currentBranch == "main" || currentBranch == "dev" {
				fmt.Fprintf(os.Stderr, "ERROR: Feature %s requires feature branch\n", featureID)
				fmt.Fprintf(os.Stderr, "  Current branch: %s\n", currentBranch)
				fmt.Fprintf(os.Stderr, "  Required branch: feature/%s\n", featureID)
				fmt.Fprintf(os.Stderr, "\nCreate the branch:\n")
				fmt.Fprintf(os.Stderr, "  git checkout -b feature/%s\n", featureID)
				fmt.Fprintf(os.Stderr, "\nOr if branch exists:\n")
				fmt.Fprintf(os.Stderr, "  git checkout feature/%s\n", featureID)
				os.Exit(1)
			}

			// Check if on correct feature branch
			expectedBranch := "feature/" + featureID
			if currentBranch == expectedBranch {
				fmt.Printf("OK: On correct feature branch %s\n", currentBranch)
				return nil
			}

			// Check if on a different feature branch
			if strings.HasPrefix(currentBranch, "feature/") {
				fmt.Fprintf(os.Stderr, "WARNING: On different feature branch\n")
				fmt.Fprintf(os.Stderr, "  Current branch: %s\n", currentBranch)
				fmt.Fprintf(os.Stderr, "  Expected: %s\n", expectedBranch)
				os.Exit(2)
			}

			// Check if on a bugfix/hotfix branch (allowed for some cases)
			if strings.HasPrefix(currentBranch, "bugfix/") || strings.HasPrefix(currentBranch, "hotfix/") {
				fmt.Printf("OK: On %s branch (non-feature work)\n", currentBranch)
				return nil
			}

			// Unknown branch type
			fmt.Fprintf(os.Stderr, "WARNING: Not on a feature branch\n")
			fmt.Fprintf(os.Stderr, "  Current branch: %s\n", currentBranch)
			fmt.Fprintf(os.Stderr, "  Expected: %s\n", expectedBranch)
			os.Exit(3)

			return nil
		},
	}

	cmd.Flags().StringVar(&featureID, "feature", "", "Feature ID to check branch for")

	return cmd
}

// guardBranchValidateCmd validates branch naming convention
func guardBranchValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <branch>",
		Short: "Validate branch naming convention",
		Long: `Validate that a branch name follows SDP conventions.

Valid prefixes:
- feature/  (e.g., feature/F065)
- bugfix/   (e.g., bugfix/sdp-1234)
- hotfix/   (e.g., hotfix/sdp-1234)

Exit codes:
  0 - Valid branch name
  1 - Invalid branch name`,
		Example: `  # Validate feature branch
  sdp guard branch validate feature/F065

  # Validate bugfix branch
  sdp guard branch validate bugfix/sdp-1234`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			branch := args[0]

			validPrefixes := []string{"feature/", "bugfix/", "hotfix/"}
			isValid := false
			var prefix string

			for _, p := range validPrefixes {
				if strings.HasPrefix(branch, p) {
					isValid = true
					prefix = p
					break
				}
			}

			// Also allow dev and main (protected but valid)
			if branch == "dev" || branch == "main" {
				fmt.Printf("OK: Protected branch %s (merge only)\n", branch)
				return nil
			}

			if !isValid {
				fmt.Fprintf(os.Stderr, "ERROR: Invalid branch name: %s\n", branch)
				fmt.Fprintf(os.Stderr, "\nValid prefixes:\n")
				fmt.Fprintf(os.Stderr, "  feature/F###   (feature implementation)\n")
				fmt.Fprintf(os.Stderr, "  bugfix/<id>    (bug fixes)\n")
				fmt.Fprintf(os.Stderr, "  hotfix/<id>    (emergency fixes)\n")
				os.Exit(1)
			}

			// Extract the ID portion
			id := strings.TrimPrefix(branch, prefix)
			if id == "" {
				fmt.Fprintf(os.Stderr, "ERROR: Missing identifier after prefix\n")
				os.Exit(1)
			}

			fmt.Printf("OK: Valid %s branch: %s\n", strings.TrimSuffix(prefix, "/"), branch)
			fmt.Printf("  ID: %s\n", id)

			return nil
		},
	}
}

