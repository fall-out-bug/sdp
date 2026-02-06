package dashboard

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// App represents the TUI dashboard application
type App struct {
	state DashboardState
	quit  bool
}

// New creates a new dashboard app
func New() *App {
	return &App{
		state: DashboardState{
			ActiveTab:   0,
			CursorPos:   0,
			Workstreams: make(map[string][]WorkstreamSummary),
			Ideas:       []IdeaSummary{},
			Loading:     true,
		},
		quit: false,
	}
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	// Start ticker for auto-refresh (every 500ms - faster!)
	return tea.Batch(
		a.tickCmd(),
		a.refreshCmd(),
	)
}

// tickCmd returns a command that ticks every 500ms
func (a *App) tickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// refreshCmd returns a command that refreshes data
func (a *App) refreshCmd() tea.Cmd {
	return func() tea.Msg {
		// Fetch real data
		a.state.Workstreams = a.fetchWorkstreams()
		a.state.Ideas = a.fetchIdeas()
		a.state.TestResults = a.fetchTestResults()

		return RefreshMsg{}
	}
}

// Update handles messages and updates the state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return a.handleKeyPress(msg)

	case TickMsg:
		// Auto-refresh tick
		return a, tea.Batch(a.tickCmd(), a.refreshCmd())

	case RefreshMsg:
		// Data refreshed
		a.state.Loading = false
		a.state.LastUpdate = time.Now()
		return a, nil

	case TabSelectMsg:
		// Tab changed
		a.state.ActiveTab = int(msg)
		return a, nil
	}

	return a, nil
}

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

// View renders the UI
// renderHeader renders the dashboard header
func (a *App) renderHeader() string {
	return matrixHeaderStyle.Render("🚀 SDP Dashboard [MATRIX MODE]")
}

// renderTabs renders the tab bar
func (a *App) renderTabs() string {
	tabs := []TabType{TabWorkstreams, TabIdeas, TabTests, TabActivity}
	var rendered string

	for i, tab := range tabs {
		tabName := fmt.Sprintf("%d. %s", i+1, tab.String())
		if i == a.state.ActiveTab {
			rendered += matrixActiveTabStyle.Render(tabName) + " "
		} else {
			rendered += matrixInactiveTabStyle.Render(tabName) + " "
		}
	}

	return rendered
}

// renderContent renders the active tab content
func (a *App) renderContent() string {
	if a.state.Loading {
		return matrixDimStyle.Render("Loading...")
	}

	switch TabType(a.state.ActiveTab) {
	case TabWorkstreams:
		return a.renderWorkstreams()
	case TabIdeas:
		return a.renderIdeas()
	case TabTests:
		return a.renderTests()
	case TabActivity:
		return a.renderActivity()
	default:
		return "Unknown tab"
	}
}

// renderWorkstreams renders the workstreams tab
func (a *App) renderWorkstreams() string {
	if len(a.state.Workstreams) == 0 {
		return "Workstreams\n\nNo workstreams found"
	}

	var content string
	content += matrixBaseStyle.Render("Workstreams\n\n")

	statusOrder := []string{"open", "in_progress", "completed", "blocked"}
	statusLabels := map[string]string{
		"open":         "Open",
		"in_progress":  "In Progress",
		"completed":    "Completed",
		"blocked":      "Blocked",
	}

	totalCount := 0
	globalIndex := 0 // Global index for cursor tracking

	for _, status := range statusOrder {
		wss, ok := a.state.Workstreams[status]
		if !ok || len(wss) == 0 {
			continue
		}

		label := statusLabels[status]

		// Use matrix style for status header
		var statusHeader string
		switch status {
		case "open":
			statusHeader = statusOpenMatrixStyle.Render(fmt.Sprintf("%s (%d)", label, len(wss)))
		case "in_progress":
			statusHeader = statusInProgressMatrixStyle.Render(fmt.Sprintf("%s (%d)", label, len(wss)))
		case "completed":
			statusHeader = statusCompletedMatrixStyle.Render(fmt.Sprintf("%s (%d)", label, len(wss)))
		case "blocked":
			statusHeader = statusBlockedMatrixStyle.Render(fmt.Sprintf("%s (%d)", label, len(wss)))
		}

		content += statusHeader + "\n"

		for _, ws := range wss {
			priority := ws.Priority
			if priority == "" {
				priority = "P2"
			}

			assignee := ""
			if ws.Assignee != "" {
				assignee = " @" + ws.Assignee
			}

			size := ""
			if ws.Size != "" {
				size = " [" + ws.Size + "]"
			}

			// Check if this item is selected
			isSelected := (globalIndex == a.state.CursorPos)

			// Style the workstream line
			var wsLine string
			if isSelected {
				// Add cursor indicator
				wsLine = matrixSelectedStyle.Render("► ") + ws.ID + ": " + ws.Title + assignee + size + " "
				wsLine += a.renderPriorityMatrix(priority)
			} else {
				wsLine = "  " + ws.ID + ": " + ws.Title + assignee + size + " "
				wsLine += a.renderPriorityMatrix(priority)
			}

			content += wsLine + "\n"
			globalIndex++
		}

		content += "\n"
		totalCount += len(wss)
	}

	content += matrixDimStyle.Render(fmt.Sprintf("Total: %d workstream(s)\n", totalCount))

	return content
}

// renderPriorityMatrix renders priority with matrix colors
func (a *App) renderPriorityMatrix(priority string) string {
	switch priority {
	case "P0":
		return priorityP0MatrixStyle.Render("["+priority+"]")
	case "P1":
		return priorityP1MatrixStyle.Render("["+priority+"]")
	case "P2":
		return priorityP2MatrixStyle.Render("["+priority+"]")
	case "P3":
		return priorityP3MatrixStyle.Render("["+priority+"]")
	default:
		return matrixBaseStyle.Render("["+priority+"]")
	}
}

// renderIdeas renders the ideas tab
func (a *App) renderIdeas() string {
	if len(a.state.Ideas) == 0 {
		return matrixBaseStyle.Render("Ideas\n\nNo ideas found")
	}

	var content string
	content += matrixBaseStyle.Render(fmt.Sprintf("Ideas (%d)\n\n", len(a.state.Ideas)))

	for i, idea := range a.state.Ideas {
		// Format date
		dateStr := idea.Date.Format("2006-01-02")

		// Check if selected
		isSelected := (i == a.state.CursorPos)

		var prefix string
		if isSelected {
			prefix = matrixSelectedStyle.Render("► ")
		} else {
			prefix = "  "
		}

		content += prefix + idea.Title + "\n"
		content += "    " + idea.Path + "\n"
		content += "    " + matrixDimStyle.Render("Last modified: "+dateStr) + "\n\n"
	}

	return content
}

// renderTests renders the tests tab
func (a *App) renderTests() string {
	content := matrixBaseStyle.Render("Tests\n\n")

	tr := a.state.TestResults
	content += matrixBaseStyle.Render(fmt.Sprintf("Coverage: %s\n", tr.Coverage))
	content += matrixBaseStyle.Render(fmt.Sprintf("Tests: %d/%d passing\n", tr.Passing, tr.Total))
	content += matrixDimStyle.Render(fmt.Sprintf("Last run: %s\n\n", tr.LastRun.Format("2006-01-02 15:04:05")))

	content += matrixBaseStyle.Render("Quality Gates:\n")
	for i, gate := range tr.QualityGates {
		isSelected := (i == a.state.CursorPos)

		var status string
		var statusStyle lipgloss.Style
		if gate.Passed {
			status = "✓"
			statusStyle = statusCompletedMatrixStyle
		} else {
			status = "✗"
			statusStyle = statusBlockedMatrixStyle
		}

		var prefix string
		if isSelected {
			prefix = matrixSelectedStyle.Render("► ")
		} else {
			prefix = "  "
		}

		content += prefix + statusStyle.Render(status) + " " + gate.Name + "\n"
	}

	return content
}

// renderActivity renders the activity tab
func (a *App) renderActivity() string {
	return "Activity\n\nNo recent activity"
}

// renderFooter renders the footer with keyboard hints
func (a *App) renderFooter() string {
	return matrixFooterStyle.Render("[↑/↓] Navigate [Enter/o] Open [r]efresh [q]uit [1-4] Tabs")
}
