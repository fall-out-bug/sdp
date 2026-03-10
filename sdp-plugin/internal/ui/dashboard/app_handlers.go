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
		a.state.Loading = true
		return a, a.refreshCmdWithForce(true)

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
		if a.state.CursorPos > 0 {
			a.state.CursorPos--
		}
		return a, nil

	case "down", "j":
		maxItems := a.maxCursorPos()
		if a.state.CursorPos < maxItems-1 {
			a.state.CursorPos++
		}
		return a, nil

	case "enter", " ":
		return a, a.openSelectedItem()

	case "o":
		return a, a.openSelectedItem()
	}

	return a, nil
}

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

func (a *App) openSelectedItem() tea.Cmd {
	return func() tea.Msg {
		return nil
	}
}
