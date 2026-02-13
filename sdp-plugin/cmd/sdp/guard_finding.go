package main

import (
	"github.com/spf13/cobra"
)

func guardFindingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "finding",
		Short:   "Manage review findings",
		Long:    `Register and manage findings from @review skill.`,
		Aliases: []string{"findings"},
	}
	cmd.AddCommand(guardFindingAddCmd())
	cmd.AddCommand(guardFindingListCmd())
	cmd.AddCommand(guardFindingResolveCmd())
	cmd.AddCommand(guardFindingClearCmd())
	return cmd
}
