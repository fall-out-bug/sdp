package dashboard

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleKeyPress handles keyboard input
func (a *App) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		a.quit = true
		return a, tea.Quit

	case "r":
		// Force refresh
		a.state.Loading = true
		return a, a.refreshCmd()

	case "1":
		a.state.CursorPos = 0 // Reset cursor when switching tabs
		return a, func() tea.Msg {
			return TabSelectMsg(TabWorkstreams)
		}

	case "2":
		a.state.CursorPos = 0
		return a, func() tea.Msg {
			return TabSelectMsg(TabIdeas)
		}

	case "3":
		a.state.CursorPos = 0
		return a, func() tea.Msg {
			return TabSelectMsg(TabTests)
		}

	case "4":
		a.state.CursorPos = 0
		return a, func() tea.Msg {
			return TabSelectMsg(TabActivity)
		}

	case "up", "k":
		// Move cursor up
		if a.state.CursorPos > 0 {
			a.state.CursorPos--
		}
		return a, nil // Return same model with updated state

	case "down", "j":
		// Move cursor down
		maxItems := a.maxCursorPos()
		if a.state.CursorPos < maxItems-1 {
			a.state.CursorPos++
		}
		return a, nil // Return same model with updated state

	case "enter", " ":
		// Open selected item
		return a, a.openSelectedItem()

	case "o":
		// Open selected item (alternative key)
		return a, a.openSelectedItem()
	}

	return a, nil
}

// maxCursorPos returns the maximum cursor position for current tab
func (a *App) maxCursorPos() int {
	switch TabType(a.state.ActiveTab) {
	case TabWorkstreams:
		count := 0
		for _, wsList := range a.state.Workstreams {
			count += len(wsList)
		}
		return count
	case TabIdeas:
		return len(a.state.Ideas)
	case TabTests:
		return len(a.state.TestResults.QualityGates)
	default:
		return 0
	}
}

// openSelectedItem opens the file for the selected item
func (a *App) openSelectedItem() tea.Cmd {
	return func() tea.Msg {
		// TODO: Implement file opening
		// For now, just print what would be opened
		return nil
	}
}
