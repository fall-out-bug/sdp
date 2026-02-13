package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/collision"
	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/parser"
	"github.com/spf13/cobra"
)

func collisionDetectCmd() *cobra.Command {
	var outputJSON bool
	cmd := &cobra.Command{
		Use:   "detect",
		Short: "Deep analysis: detect shared interface boundaries across features",
		Long:  `Analyze not just file overlaps, but shared types, structs, and interfaces that need coordination.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCollisionDetect(cmd, args, outputJSON)
		},
	}
	cmd.Flags().BoolVar(&outputJSON, "output-json", false, "Output in JSON format")
	return cmd
}

// runCollisionDetect runs deep boundary detection.
func runCollisionDetect(cmd *cobra.Command, args []string, outputJSON bool) error {
	root, err := config.FindProjectRoot()
	if err != nil {
		return fmt.Errorf("find project root: %w", err)
	}
	featureScopes, err := loadFeatureScopes(root)
	if err != nil {
		return err
	}
	boundaries := collision.DetectBoundaries(featureScopes)
	if len(boundaries) == 0 {
		fmt.Println("No shared boundaries detected.")
		return nil
	}
	if outputJSON {
		return outputBoundariesAsJSON(boundaries)
	}
	return outputBoundariesAsHuman(boundaries)
}

// loadFeatureScopes loads feature scopes from in-progress workstreams.
func loadFeatureScopes(projectRoot string) ([]collision.FeatureScope, error) {
	inProgressDir := filepath.Join(projectRoot, "docs", "workstreams", "in_progress")
	entries, err := os.ReadDir(inProgressDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []collision.FeatureScope{}, nil
		}
		return nil, fmt.Errorf("read in_progress dir: %w", err)
	}
	scopes := make([]collision.FeatureScope, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !hasSuffix(e.Name(), ".md") {
			continue
		}
		path := filepath.Join(inProgressDir, e.Name())
		ws, err := parser.ParseWorkstream(path)
		if err != nil {
			continue
		}
		files := append([]string{}, ws.Scope.Implementation...)
		files = append(files, ws.Scope.Tests...)
		scopes = append(scopes, collision.FeatureScope{
			FeatureID:  ws.Feature,
			ScopeFiles: files,
		})
	}
	return scopes, nil
}

// outputBoundariesAsHuman prints boundaries in human-readable format.
func outputBoundariesAsHuman(boundaries []collision.SharedBoundary) error {
	fmt.Println("🔗 Shared boundaries detected:")
	fmt.Println()
	for _, b := range boundaries {
		fmt.Printf("  File: %s\n", b.FileName)
		fmt.Printf("  Type: %s\n", b.TypeName)
		fmt.Printf("  Features: %v\n", b.Features)
		if len(b.Fields) > 0 {
			fmt.Printf("  Fields:\n")
			for _, f := range b.Fields {
				fmt.Printf("    - %s: %s\n", f.Name, f.Type)
			}
		}
		fmt.Printf("  Recommendation: %s\n", b.Recommendation)
		fmt.Println()
	}
	fmt.Printf("  %d shared boundary(ies)\n", len(boundaries))
	return nil
}

// outputBoundariesAsJSON prints boundaries in JSON format.
func outputBoundariesAsJSON(boundaries []collision.SharedBoundary) error {
	fmt.Println("[")
	for i, b := range boundaries {
		jsonStr, err := collision.BoundaryToJSON(b)
		if err != nil {
			continue
		}
		fmt.Print(jsonStr)
		if i < len(boundaries)-1 {
			fmt.Println(",")
		} else {
			fmt.Println()
		}
	}
	fmt.Println("]")
	return nil
}
