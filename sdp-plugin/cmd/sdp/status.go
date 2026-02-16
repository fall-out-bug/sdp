package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/fall-out-bug/sdp/internal/ui/dashboard"
)

// statusCmd returns the status command
func statusCmd() *cobra.Command {
	var textMode bool
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show project status",
		Long: `Display project status including workstreams, active work, and environment health.

Output Modes:
  --text    Simple text output (no TUI, useful for scripts)
  --json    Machine-readable JSON output
  (default) Rich TUI dashboard

The TUI dashboard shows:
  - Workstreams (grouped by status: open, in-progress, completed, blocked)
  - Ideas (from docs/drafts/)
  - Test results and coverage
  - Activity log

Keyboard shortcuts (TUI only):
  [1-4] - Switch tabs
  [r]   - Refresh data
  [q]   - Quit dashboard`,
		Example: `  # Interactive TUI dashboard
  sdp status

  # Quick text status
  sdp status --text

  # JSON output for scripts
  sdp status --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if textMode || jsonOutput {
				return runTextStatus(jsonOutput)
			}

			app := dashboard.New()
			p := tea.NewProgram(
				app,
				tea.WithAltScreen(),
				tea.WithMouseCellMotion(),
			)

			_, err := p.Run()
			return err
		},
	}

	cmd.Flags().BoolVar(&textMode, "text", false, "Text output (no TUI)")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "JSON output")

	return cmd
}
