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

func collisionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collision",
		Short: "Scope collision detection for parallel workstreams",
		Long:  `Detect when in-progress workstreams touch the same files or directories.`,
	}
	cmd.AddCommand(collisionCheckCmd())
	return cmd
}

func collisionCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "List scope overlaps across in-progress workstreams",
		RunE:  runCollisionCheck,
	}
}

func runCollisionCheck(cmd *cobra.Command, args []string) error {
	root, err := config.FindProjectRoot()
	if err != nil {
		return fmt.Errorf("find project root: %w", err)
	}
	scopes, err := loadInProgressScopes(root)
	if err != nil {
		return err
	}
	overlaps := collision.DetectOverlaps(scopes)
	if len(overlaps) == 0 {
		fmt.Println("No scope overlaps detected.")
		return nil
	}
	fmt.Println("⚠️  Scope overlaps detected:")
	fmt.Println()
	for _, o := range overlaps {
		fmt.Printf("  %s\n", o.File)
		for _, wsID := range o.Workstreams {
			fmt.Printf("    → %s\n", wsID)
		}
		fmt.Println()
	}
	fmt.Printf("  %d overlap(s) across workstreams\n", len(overlaps))
	fmt.Println("  Recommendation: coordinate or sequence these workstreams")
	return nil
}

func loadInProgressScopes(projectRoot string) ([]collision.WorkstreamScope, error) {
	inProgressDir := filepath.Join(projectRoot, "docs", "workstreams", "in_progress")
	entries, err := os.ReadDir(inProgressDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read in_progress dir: %w", err)
	}
	var scopes []collision.WorkstreamScope
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
		scopes = append(scopes, collision.WorkstreamScope{
			ID:         ws.ID,
			Status:     "in_progress",
			ScopeFiles: files,
		})
	}
	return scopes, nil
}

func hasSuffix(s, suf string) bool {
	return len(s) >= len(suf) && s[len(s)-len(suf):] == suf
}

// loadScopesForGuard loads in-progress scopes and optionally includes the WS being activated.
func loadScopesForGuard(projectRoot, activatingWSID string) ([]collision.WorkstreamScope, error) {
	scopes, err := loadInProgressScopes(projectRoot)
	if err != nil {
		return nil, err
	}
	if activatingWSID == "" {
		return scopes, nil
	}
	// Try to parse the activating WS from backlog or in_progress
	for _, sub := range []string{"backlog", "in_progress"} {
		p := filepath.Join(projectRoot, "docs", "workstreams", sub, activatingWSID+".md")
		ws, err := parser.ParseWorkstream(p)
		if err != nil {
			continue
		}
		files := append([]string{}, ws.Scope.Implementation...)
		files = append(files, ws.Scope.Tests...)
		scopes = append(scopes, collision.WorkstreamScope{
			ID:         ws.ID,
			Status:     "in_progress",
			ScopeFiles: files,
		})
		break
	}
	return scopes, nil
}

// scopeFilesForWS returns scope files for the given workstream ID (for evidence plan event).
func scopeFilesForWS(wsID string) []string {
	root, err := config.FindProjectRoot()
	if err != nil {
		return nil
	}
	scopes, err := loadScopesForGuard(root, wsID)
	if err != nil {
		return nil
	}
	for _, s := range scopes {
		if s.ID == wsID {
			return s.ScopeFiles
		}
	}
	return nil
}

// warnCollisionIfAny finds project root, loads scopes (including activatingWSID), and prints overlap warning if any.
func warnCollisionIfAny(activatingWSID string) {
	root, err := config.FindProjectRoot()
	if err != nil {
		return
	}
	scopes, err := loadScopesForGuard(root, activatingWSID)
	if err != nil || len(scopes) == 0 {
		return
	}
	overlaps := collision.DetectOverlaps(scopes)
	if len(overlaps) == 0 {
		return
	}
	fmt.Fprintln(os.Stderr, "⚠️  Scope overlaps detected with other in-progress workstreams:")
	for _, o := range overlaps {
		fmt.Fprintf(os.Stderr, "  %s → %v\n", o.File, o.Workstreams)
	}
	fmt.Fprintln(os.Stderr, "  Run 'sdp collision check' for full report.")
}
