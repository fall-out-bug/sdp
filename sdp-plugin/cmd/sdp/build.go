package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/guard"
	"github.com/spf13/cobra"
)

func buildCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "build <ws-id>",
		Short: "Execute a single workstream (guard + go test)",
		Long: `Run workstream execution: activate guard, run pre-build hook, go test, post-build hook.
For full TDD cycle with agent, use @build or sdp-orchestrate.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()

			wsID := args[0]
			root, err := config.FindProjectRoot()
			if err != nil {
				return fmt.Errorf("find project root: %w", err)
			}
			if err := os.Chdir(root); err != nil {
				return fmt.Errorf("chdir: %w", err)
			}
			configDir := os.Getenv("XDG_CONFIG_HOME")
			if configDir == "" {
				configDir, _ = os.UserConfigDir()
			}
			skill := guard.NewSkill(filepath.Join(configDir, "sdp"))
			if err := skill.Activate(wsID); err != nil {
				return fmt.Errorf("guard activate: %w", err)
			}
			preHook := filepath.Join(root, "scripts", "hooks", "pre-build.sh")
			if _, err := os.Stat(preHook); err == nil {
				c := exec.CommandContext(ctx, preHook, wsID)
				c.Dir = root
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				if err := c.Run(); err != nil {
					return fmt.Errorf("pre-build hook: %w", err)
				}
			}
			goTest := exec.CommandContext(ctx, "go", "test", "./...")
			goTest.Dir = root
			goTest.Stdout = os.Stdout
			goTest.Stderr = os.Stderr
			if err := goTest.Run(); err != nil {
				return fmt.Errorf("go test: %w", err)
			}
			postHook := filepath.Join(root, "scripts", "hooks", "post-build.sh")
			if _, err := os.Stat(postHook); err == nil {
				c := exec.CommandContext(ctx, postHook, wsID, "completed")
				c.Dir = root
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				_ = c.Run()
			}
			return nil
		},
	}
}
