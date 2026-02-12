package main

import (
	"github.com/spf13/cobra"
)

func guardCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "guard",
		Short: "Pre-edit guard for quality gate enforcement",
		Long: `Guard commands for managing workstream editing scope.

Prevents editing files outside of active workstream's scope.
This is part of TDD discipline - one workstream at a time.

Examples:
  sdp guard activate 00-001-01
  sdp guard check internal/file.go
  sdp guard status
  sdp guard context check
  sdp guard branch check --feature=F065`,
	}

	cmd.AddCommand(guardActivate())
	cmd.AddCommand(guardCheck())
	cmd.AddCommand(guardStatus())
	cmd.AddCommand(guardDeactivate())
	cmd.AddCommand(guardContextCmd())
	cmd.AddCommand(guardBranchCmd())

	return cmd
}
