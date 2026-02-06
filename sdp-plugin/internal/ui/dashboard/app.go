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
			Workstreams: make(map[string][]WorkstreamSummary),
			Ideas:       []IdeaSummary{},
			Loading:     true,
		},
		quit: false,
	}
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	// Start ticker for auto-refresh (every 2 seconds)
	return tea.Batch(
		a.tickCmd(),
		a.refreshCmd(),
	)
}

// tickCmd returns a command that ticks every 2 seconds
func (a *App) tickCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
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
		return a, func() tea.Msg {
			return TabSelectMsg(TabWorkstreams)
		}

	case "2":
		return a, func() tea.Msg {
			return TabSelectMsg(TabIdeas)
		}

	case "3":
		return a, func() tea.Msg {
			return TabSelectMsg(TabTests)
		}

	case "4":
		return a, func() tea.Msg {
			return TabSelectMsg(TabActivity)
		}
	}

	return a, nil
}

// View renders the UI
func (a *App) View() string {
	if a.quit {
		return ""
	}

	// Build the view
	var view string

	view += a.renderHeader()
	view += "\n"
	view += a.renderTabs()
	view += "\n"
	view += a.renderContent()
	view += "\n"
	view += a.renderFooter()

	return view
}

// renderHeader renders the dashboard header
func (a *App) renderHeader() string {
	return headerStyle.Render("🚀 SDP Dashboard")
}

// renderTabs renders the tab bar
func (a *App) renderTabs() string {
	tabs := []TabType{TabWorkstreams, TabIdeas, TabTests, TabActivity}
	var rendered string

	for i, tab := range tabs {
		tabName := fmt.Sprintf("%d. %s", i+1, tab.String())
		if i == a.state.ActiveTab {
			rendered += activeTabStyle.Render(tabName) + " "
		} else {
			rendered += inactiveTabStyle.Render(tabName) + " "
		}
	}

	return rendered
}

// renderContent renders the active tab content
func (a *App) renderContent() string {
	if a.state.Loading {
		return lipgloss.NewStyle().Faint(true).Render("Loading...")
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
	content += "Workstreams\n\n"

	statusOrder := []string{"open", "in_progress", "completed", "blocked"}
	statusLabels := map[string]string{
		"open":         "Open",
		"in_progress":  "In Progress",
		"completed":    "Completed",
		"blocked":      "Blocked",
	}

	totalCount := 0
	for _, status := range statusOrder {
		wss, ok := a.state.Workstreams[status]
		if !ok || len(wss) == 0 {
			continue
		}

		label := statusLabels[status]
		content += StatusStyle(status).Render(fmt.Sprintf("%s (%d)", label, len(wss))) + "\n"

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

			priorityStyled := PriorityStyle(priority).Render("["+priority+"]")
			content += fmt.Sprintf("  %s: %s%s%s %s\n", ws.ID, ws.Title, assignee, size, priorityStyled)
		}

		content += "\n"
		totalCount += len(wss)
	}

	content += fmt.Sprintf("Total: %d workstream(s)\n", totalCount)

	return content
}

// renderIdeas renders the ideas tab
func (a *App) renderIdeas() string {
	if len(a.state.Ideas) == 0 {
		return "Ideas\n\nNo ideas found"
	}

	var content string
	content += fmt.Sprintf("Ideas (%d)\n\n", len(a.state.Ideas))

	for _, idea := range a.state.Ideas {
		// Format date
		dateStr := idea.Date.Format("2006-01-02")
		content += fmt.Sprintf("• %s\n  %s\n  Last modified: %s\n\n", idea.Title, idea.Path, dateStr)
	}

	return content
}

// renderTests renders the tests tab
func (a *App) renderTests() string {
	content := "Tests\n\n"

	tr := a.state.TestResults
	content += fmt.Sprintf("Coverage: %s\n", tr.Coverage)
	content += fmt.Sprintf("Tests: %d/%d passing\n", tr.Passing, tr.Total)
	content += fmt.Sprintf("Last run: %s\n\n", tr.LastRun.Format("2006-01-02 15:04:05"))

	content += "Quality Gates:\n"
	for _, gate := range tr.QualityGates {
		status := "✗"
		if gate.Passed {
			status = "✓"
		}
		content += fmt.Sprintf("  %s %s\n", status, gate.Name)
	}

	return content
}

// renderActivity renders the activity tab
func (a *App) renderActivity() string {
	return "Activity\n\nNo recent activity"
}

// renderFooter renders the footer with keyboard hints
func (a *App) renderFooter() string {
	return "[r]efresh [q]uit [1-4] Tabs"
}
