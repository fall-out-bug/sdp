package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/fall-out-bug/sdp/internal/ui/dashboard"
)

// statusCmd returns the status command
func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show TUI dashboard with project status",
		Long: `Display a rich terminal UI (TUI) dashboard showing:

- Workstreams (grouped by status: open, in-progress, completed, blocked)
- Ideas (from docs/drafts/)
- Test results and coverage
- Activity log

Keyboard shortcuts:
  [1-4] - Switch tabs
  [r]   - Refresh data
  [q]   - Quit dashboard`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create and run the dashboard
			app := dashboard.New()

			// Use full screen with alt screen + mouse support
			p := tea.NewProgram(
				app,
				tea.WithAltScreen(),
				tea.WithMouseCellMotion(),
			)

			if _, err := p.Run(); err != nil {
				return err
			}

			return nil
		},
	}
}
