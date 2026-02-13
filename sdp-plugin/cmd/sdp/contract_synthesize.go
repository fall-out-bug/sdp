package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fall-out-bug/sdp/internal/collision"
	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/spf13/cobra"
)

func runContractSynthesize(cmd *cobra.Command, args []string) error {
	featureName, err := cmd.Flags().GetString("feature")
	if err != nil {
		return fmt.Errorf("failed to get feature flag: %w", err)
	}
	requirementsPath, err := cmd.Flags().GetString("requirements")
	if err != nil {
		return fmt.Errorf("failed to get requirements flag: %w", err)
	}
	outputPath, err := cmd.Flags().GetString("output")
	if err != nil {
		return fmt.Errorf("failed to get output flag: %w", err)
	}

	// Set default requirements path if not provided
	if requirementsPath == "" {
		requirementsPath = fmt.Sprintf("docs/drafts/%s-idea.md", featureName)
	}

	// Set default output path if not provided
	if outputPath == "" {
		outputPath = fmt.Sprintf(".contracts/%s.yaml", featureName)
	}

	fmt.Printf("✓ Generating contract for feature: %s\n", featureName)
	fmt.Printf("  Requirements: %s\n", requirementsPath)
	fmt.Printf("  Output: %s\n", outputPath)
	fmt.Printf("\n⚠️  Contract synthesis not yet implemented\n")
	fmt.Printf("   This will require integration with multi-agent synthesis system\n")

	return nil
}

func runContractGenerate(cmd *cobra.Command, args []string) error {
	featuresFlag, err := cmd.Flags().GetString("features")
	if err != nil {
		return fmt.Errorf("failed to get features flag: %w", err)
	}

	// Parse feature IDs
	var featureIDs []string
	if featuresFlag != "" {
		for _, f := range strings.Split(featuresFlag, ",") {
			featureIDs = append(featureIDs, strings.TrimSpace(f))
		}
	}

	root, err := config.FindProjectRoot()
	if err != nil {
		return fmt.Errorf("find project root: %w", err)
	}

	// Load feature scopes
	featureScopes, err := loadFeatureScopes(root)
	if err != nil {
		return fmt.Errorf("load feature scopes: %w", err)
	}

	// Filter by specified features if provided
	if len(featureIDs) > 0 {
		filtered := make([]collision.FeatureScope, 0)
		for _, fs := range featureScopes {
			for _, fid := range featureIDs {
				if fs.FeatureID == fid {
					filtered = append(filtered, fs)
					break
				}
			}
		}
		featureScopes = filtered
	}

	// Detect boundaries using the collision package
	boundaries := collision.DetectBoundaries(featureScopes)

	if len(boundaries) == 0 {
		fmt.Println("No shared boundaries detected.")
		fmt.Println("  Run 'sdp collision detect' to find shared interfaces.")
		return nil
	}

	// Generate contracts
	contractsDir := filepath.Join(root, ".contracts")
	contracts, err := collision.GenerateContracts(boundaries, contractsDir)
	if err != nil {
		return fmt.Errorf("generate contracts: %w", err)
	}

	fmt.Printf("✓ Generated %d contract(s)\n", len(contracts))
	for _, c := range contracts {
		fmt.Printf("  - %s.yaml (required by: %v)\n", c.TypeName, c.RequiredBy)
	}
	fmt.Printf("\n  Output directory: %s\n", contractsDir)

	// AC2: Run validation post-generation (timing fix for sdp-ubdr)
	// Check if implementation directory exists and validate
	implDir := filepath.Join(root, "internal")
	if _, err := os.Stat(implDir); err == nil {
		fmt.Println("\n→ Running post-generation validation...")
		if err := validateImplementation(contractsDir, implDir); err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Post-generation validation: %v\n", err)
			// Note: In P1, violations are warnings, not blockers
		}
	}

	fmt.Printf("\n  Next step: sdp contract lock .contracts/<type>.yaml\n")

	return nil
}
