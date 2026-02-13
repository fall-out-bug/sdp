package dashboard

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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

// View renders the UI with full black background
func (a *App) View() string {
	if a.quit {
		return ""
	}

	// Build the full view with black background
	content := a.renderHeader()
	content += "\n"
	content += a.renderTabs()
	content += "\n\n"
	content += a.renderContent()
	content += "\n\n"
	content += a.renderFooter()

	// Wrap everything in black background
	return matrixBaseStyle.Render(content)
}

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
		return matrixBrightStyle.Render("Loading...")
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

// renderFooter renders the footer with keyboard hints
func (a *App) renderFooter() string {
	return matrixFooterStyle.Render("[↑/↓] Navigate [Enter/o] Open [r]efresh [q]uit [1-4] Tabs")
}
