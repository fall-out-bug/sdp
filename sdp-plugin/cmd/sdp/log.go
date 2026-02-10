package main

import (
	"fmt"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/spf13/cobra"
)

const defaultLogPath = ".sdp/log/events.jsonl"
const defaultRecentN = 20

func logCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log",
		Short: "Evidence log: show, export, and trace events",
		Long:  `Show recent events, export logs, or trace evidence chain by commit/workstream.`,
		RunE:  func(cmd *cobra.Command, args []string) error { return runLogShow("", "") },
	}
	cmd.AddCommand(logShowCmd())
	cmd.AddCommand(logExportCmd())
	cmd.AddCommand(logStatsCmd())
	cmd.AddCommand(logTraceCmd())
	return cmd
}

func evidenceLogPath() (string, error) {
	root, err := config.FindProjectRoot()
	if err != nil {
		return "", fmt.Errorf("find project root: %w", err)
	}
	cfg, _ := config.Load(root)
	logPath := defaultLogPath
	if cfg != nil && cfg.Evidence.LogPath != "" {
		logPath = cfg.Evidence.LogPath
	}
	return filepath.Join(root, logPath), nil
}

