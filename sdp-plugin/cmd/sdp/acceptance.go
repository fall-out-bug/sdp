package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fall-out-bug/sdp/internal/acceptance"
	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/spf13/cobra"
)

func acceptanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acceptance",
		Short: "Run acceptance (smoke) test gate",
		Long:  `Runs the command from .sdp/config.yml acceptance section. Exit 0 = pass.`,
	}
	cmd.AddCommand(acceptanceRunCmd())
	return cmd
}

func acceptanceRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Execute acceptance test command",
		RunE:  runAcceptanceRun,
	}
}

func runAcceptanceRun(cmd *cobra.Command, args []string) error {
	root, err := config.FindProjectRoot()
	if err != nil {
		return fmt.Errorf("find project root: %w", err)
	}
	cfg, err := config.Load(root)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	timeout, err := acceptance.ParseTimeout(cfg.Acceptance.Timeout)
	if err != nil {
		timeout = 30 * time.Second
	}
	r := &acceptance.Runner{
		Command:  cfg.Acceptance.Command,
		Timeout:  timeout,
		Expected: cfg.Acceptance.Expected,
	}
	res, err := r.Run(context.Background())
	if err != nil {
		return err
	}
	if res.Passed {
		fmt.Printf("Acceptance: PASS (%.1fs)\n", res.Duration.Seconds())
		return nil
	}
	fmt.Fprintf(os.Stderr, "Acceptance: FAIL — %s (%.1fs)\n", res.Error, res.Duration.Seconds())
	if res.Output != "" {
		fmt.Fprintf(os.Stderr, "%s\n", res.Output)
	}
	os.Exit(1)
	return nil
}

// runAcceptanceFromConfig runs acceptance if config is present; returns (passed, skipped, error). AC6: skip when no config.
func runAcceptanceFromConfig(projectRoot string) (passed bool, skipped bool, err error) {
	cfgPath := filepath.Join(projectRoot, ".sdp", "config.yml")
	if _, err := os.Stat(cfgPath); err != nil {
		return false, true, nil // no config file — skip
	}
	cfg, err := config.Load(projectRoot)
	if err != nil {
		return false, true, err
	}
	if cfg.Acceptance.Command == "" {
		return false, true, nil
	}
	timeout, _ := acceptance.ParseTimeout(cfg.Acceptance.Timeout)
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	r := &acceptance.Runner{
		Command:  cfg.Acceptance.Command,
		Timeout:  timeout,
		Expected: cfg.Acceptance.Expected,
	}
	res, err := r.Run(context.Background())
	if err != nil {
		return false, false, err
	}
	return res.Passed, false, nil
}
