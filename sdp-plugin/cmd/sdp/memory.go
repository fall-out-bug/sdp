package main

import (
	"github.com/spf13/cobra"
)

func memoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "memory",
		Short: "Long-term memory management",
		Long: `Manage the SDP long-term memory system.

The memory system indexes project artifacts for fast search and
provides hybrid search capabilities (full-text + semantic).

Examples:
  sdp memory index              # Index all docs/ artifacts
  sdp memory search "API"       # Search for "API" in artifacts
  sdp memory stats              # Show memory statistics`,
	}

	cmd.AddCommand(memoryIndexCmd())
	cmd.AddCommand(memorySearchCmd())
	cmd.AddCommand(memoryStatsCmd())

	return cmd
}
