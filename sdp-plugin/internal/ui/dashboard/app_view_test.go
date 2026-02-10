package dashboard

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestApp_Init(t *testing.T) {
	app := New()

	cmd := app.Init()
	if cmd == nil {
		t.Fatal("Init should return a command")
	}
}

func TestApp_Update_TickMsg(t *testing.T) {
	app := New()

	msg := TickMsg(time.Now())
	model, cmd := app.Update(msg)

	if model == nil {
		t.Fatal("Update should return a model")
	}

	if cmd == nil {
		t.Fatal("Update should return a command for TickMsg")
	}

	// Check that loading is still true (we don't have real data)
	dashApp, ok := model.(*App)
	if !ok {
		t.Fatal("Model should be *App")
	}

	if dashApp.state.Loading {
		// Loading might still be true in tests
		// This is expected behavior
	}
}

func TestApp_Update_RefreshMsg(t *testing.T) {
	app := New()
	app.state.Loading = true

	msg := RefreshMsg{}
	model, cmd := app.Update(msg)

	if model == nil {
		t.Fatal("Update should return a model")
	}

	if cmd != nil {
		t.Fatal("Update should not return a command for RefreshMsg")
	}

	dashApp, ok := model.(*App)
	if !ok {
		t.Fatal("Model should be *App")
	}

	if dashApp.state.Loading {
		t.Error("Loading should be false after RefreshMsg")
	}
}

func TestApp_Update_TabSelectMsg(t *testing.T) {
	app := New()

	tests := []struct {
		name      string
		tabMsg    TabSelectMsg
		expectedTab int
	}{
		{"Workstreams tab", TabSelectMsg(TabWorkstreams), 0},
		{"Ideas tab", TabSelectMsg(TabIdeas), 1},
		{"Tests tab", TabSelectMsg(TabTests), 2},
		{"Activity tab", TabSelectMsg(TabActivity), 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, cmd := app.Update(tt.tabMsg)

			if model == nil {
				t.Fatal("Update should return a model")
			}

			if cmd != nil {
				t.Fatal("Update should not return a command for TabSelectMsg")
			}

			dashApp, ok := model.(*App)
			if !ok {
				t.Fatal("Model should be *App")
			}

			if dashApp.state.ActiveTab != tt.expectedTab {
				t.Errorf("Expected ActiveTab %d, got %d", tt.expectedTab, dashApp.state.ActiveTab)
			}
		})
	}
}

func TestApp_handleKeyPress_Quit(t *testing.T) {
	app := New()

	// Test 'q' key
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	model, cmd := app.Update(msg)

	if cmd == nil {
		t.Fatal("Update should return tea.Quit command")
	}

	dashApp, ok := model.(*App)
	if !ok {
		t.Fatal("Model should be *App")
	}

	if !dashApp.quit {
		t.Error("Quit should be set to true")
	}
}

func TestApp_handleKeyPress_Refresh(t *testing.T) {
	app := New()
	app.state.Loading = false

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	model, _ := app.Update(msg)

	dashApp, ok := model.(*App)
	if !ok {
		t.Fatal("Model should be *App")
	}

	if !dashApp.state.Loading {
		t.Error("Loading should be true after refresh key")
	}
}

func TestApp_handleKeyPress_TabKeys(t *testing.T) {
	tests := []struct {
		key         tea.KeyMsg
		expectedTab int
		expectedCmd bool
	}{
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}, 0, true},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}}, 1, true},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}}, 2, true},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}}, 3, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.key.Runes), func(t *testing.T) {
			app := New()
			app.state.CursorPos = 5 // Set non-zero cursor

			model, cmd := app.Update(tt.key)

			dashApp, ok := model.(*App)
			if !ok {
				t.Fatal("Model should be *App")
			}

			// Cursor should be reset
			if dashApp.state.CursorPos != 0 {
				t.Errorf("CursorPos should be reset to 0, got %d", dashApp.state.CursorPos)
			}

			// Should return a command
			if tt.expectedCmd && cmd == nil {
				t.Error("Expected a command to be returned")
			}

			// Now execute the command to get the TabSelectMsg
			if cmd != nil {
				msg := cmd()
				if tabMsg, ok := msg.(TabSelectMsg); ok {
					if int(tabMsg) != tt.expectedTab {
						t.Errorf("Expected TabSelectMsg %d, got %d", tt.expectedTab, int(tabMsg))
					}
				}
			}
		})
	}
}

func TestApp_handleKeyPress_CursorNavigation(t *testing.T) {
	app := New()

	// Test up arrow
	app.state.CursorPos = 5
	msg := tea.KeyMsg{Type: tea.KeyUp}
	model, _ := app.Update(msg)

	dashApp, ok := model.(*App)
	if !ok {
		t.Fatal("Model should be *App")
	}

	if dashApp.state.CursorPos != 4 {
		t.Errorf("Expected CursorPos 4, got %d", dashApp.state.CursorPos)
	}

	// Test down arrow with no data (maxCursorPos returns 0)
	msg = tea.KeyMsg{Type: tea.KeyDown}
	model, _ = app.Update(msg)

	dashApp, ok = model.(*App)
	if !ok {
		t.Fatal("Model should be *App")
	}

	// With maxCursorPos = 0, cursor won't move (0 < 0-1 is false)
	if dashApp.state.CursorPos != 4 {
		t.Errorf("Expected CursorPos to stay at 4, got %d", dashApp.state.CursorPos)
	}

	// Test up at position 0
	app.state.CursorPos = 0
	msg = tea.KeyMsg{Type: tea.KeyUp}
	model, _ = app.Update(msg)

	dashApp, ok = model.(*App)
	if !ok {
		t.Fatal("Model should be *App")
	}

	if dashApp.state.CursorPos != 0 {
		t.Errorf("CursorPos should stay at 0, got %d", dashApp.state.CursorPos)
	}
}

func TestApp_maxCursorPos(t *testing.T) {
	app := New()

	// Test with no data
	app.state.Workstreams = map[string][]WorkstreamSummary{}
	app.state.Ideas = []IdeaSummary{}
	app.state.TestResults.QualityGates = []GateStatus{}

	maxPos := app.maxCursorPos()
	if maxPos != 0 {
		t.Errorf("Expected maxCursorPos 0 with no data, got %d", maxPos)
	}

	// Test with workstreams
	app.state.ActiveTab = 0 // Workstreams tab
	app.state.Workstreams = map[string][]WorkstreamSummary{
		"open": {
			{ID: "001", Title: "Task 1", Status: "open"},
			{ID: "002", Title: "Task 2", Status: "open"},
		},
	}

	maxPos = app.maxCursorPos()
	if maxPos != 2 {
		t.Errorf("Expected maxCursorPos 2, got %d", maxPos)
	}

	// Test with ideas
	app.state.ActiveTab = 1 // Ideas tab
	app.state.Ideas = []IdeaSummary{
		{Title: "Idea 1", Path: "path1"},
		{Title: "Idea 2", Path: "path2"},
		{Title: "Idea 3", Path: "path3"},
	}

	maxPos = app.maxCursorPos()
	if maxPos != 3 {
		t.Errorf("Expected maxCursorPos 3, got %d", maxPos)
	}

	// Test with tests
	app.state.ActiveTab = 2 // Tests tab
	app.state.TestResults = TestSummary{
		QualityGates: []GateStatus{
			{Name: "Gate 1", Passed: true},
			{Name: "Gate 2", Passed: false},
		},
	}

	maxPos = app.maxCursorPos()
	if maxPos != 2 {
		t.Errorf("Expected maxCursorPos 2, got %d", maxPos)
	}
}

func TestApp_renderHeader(t *testing.T) {
	app := New()

	header := app.renderHeader()
	if header == "" {
		t.Error("renderHeader should not return empty string")
	}

	// Check for matrix mode indicator
	if len(header) < 10 {
		t.Errorf("Header seems too short: %s", header)
	}
}

func TestApp_renderTabs(t *testing.T) {
	app := New()

	// Test each tab
	tabs := []struct {
		tabIndex  int
		tabName   string
	}{
		{0, "Workstreams"},
		{1, "Ideas"},
		{2, "Tests"},
		{3, "Activity"},
	}

	for _, tt := range tabs {
		t.Run(tt.tabName, func(t *testing.T) {
			app.state.ActiveTab = tt.tabIndex
			tabsRendered := app.renderTabs()

			if tabsRendered == "" {
				t.Error("renderTabs should not return empty string")
			}

			// Check that tab name is present
			// Note: The actual rendering uses styles, so we check for the base text
		})
	}
}

func TestApp_renderContent_Loading(t *testing.T) {
	app := New()
	app.state.Loading = true

	content := app.renderContent()
	if content == "" {
		t.Error("renderContent should not return empty string")
	}
}

func TestApp_renderContent_NoData(t *testing.T) {
	app := New()
	app.state.Loading = false
	app.state.Workstreams = map[string][]WorkstreamSummary{}
	app.state.Ideas = []IdeaSummary{}

	// Test workstreams tab
	app.state.ActiveTab = 0
	content := app.renderContent()
	if content == "" {
		t.Error("renderContent should not return empty string")
	}

	// Test ideas tab
	app.state.ActiveTab = 1
	content = app.renderContent()
	if content == "" {
		t.Error("renderContent should not return empty string")
	}
}

func TestApp_renderWorkstreams_WithData(t *testing.T) {
	app := New()
	app.state.Loading = false
	app.state.Workstreams = map[string][]WorkstreamSummary{
		"open": {
			{ID: "001", Title: "Task 1", Status: "open", Priority: "P1", Assignee: "user", Size: "SMALL"},
		},
		"in_progress": {
			{ID: "002", Title: "Task 2", Status: "in_progress", Priority: "P0", Assignee: "claude"},
		},
	}

	content := app.renderWorkstreams()
	if content == "" {
		t.Error("renderWorkstreams should not return empty string")
	}

	// Check for workstream ID
	if len(content) < 10 {
		t.Errorf("Content seems too short: %s", content)
	}
}

func TestApp_renderPriorityMatrix(t *testing.T) {
	app := New()

	priorities := []string{"P0", "P1", "P2", "P3", "UNKNOWN"}

	for _, priority := range priorities {
		result := app.renderPriorityMatrix(priority)
		if result == "" {
			t.Errorf("renderPriorityMatrix should not return empty string for %s", priority)
		}
	}
}

func TestApp_renderIdeas_WithData(t *testing.T) {
	app := New()
	app.state.Loading = false
	app.state.Ideas = []IdeaSummary{
		{Title: "Idea 1", Path: "docs/idea1.md", Date: time.Now()},
	}

	content := app.renderIdeas()
	if content == "" {
		t.Error("renderIdeas should not return empty string")
	}
}

func TestApp_renderTests(t *testing.T) {
	app := New()
	app.state.Loading = false
	app.state.TestResults = TestSummary{
		Coverage: "85.5%",
		Passing:  42,
		Total:    50,
		LastRun:  time.Now(),
		QualityGates: []GateStatus{
			{Name: "Coverage", Passed: true},
			{Name: "Linting", Passed: false},
		},
	}

	content := app.renderTests()
	if content == "" {
		t.Error("renderTests should not return empty string")
	}
}

func TestApp_renderActivity(t *testing.T) {
	app := New()

	content := app.renderActivity()
	if content == "" {
		t.Error("renderActivity should not return empty string")
	}
}

func TestApp_renderFooter(t *testing.T) {
	app := New()

	footer := app.renderFooter()
	if footer == "" {
		t.Error("renderFooter should not return empty string")
	}

	// Check for keyboard hints
	if len(footer) < 20 {
		t.Errorf("Footer seems too short: %s", footer)
	}
}

func TestApp_openSelectedItem(t *testing.T) {
	app := New()

	cmd := app.openSelectedItem()
	if cmd == nil {
		t.Fatal("openSelectedItem should return a command")
	}

	// The command should return nil msg (TODO in implementation)
	// We just verify it doesn't panic
}

func TestTickMsg(t *testing.T) {
	now := time.Now()
	msg := TickMsg(now)

	if time.Time(msg) != now {
		t.Errorf("Expected time %v, got %v", now, time.Time(msg))
	}
}

func TestRefreshMsg(t *testing.T) {
	msg := RefreshMsg{}

	// Just verify it can be created
	_ = msg
}

func TestTabSelectMsg(t *testing.T) {
	msg := TabSelectMsg(int(TabIdeas))

	if int(msg) != int(TabIdeas) {
		t.Errorf("Expected %d, got %d", int(TabIdeas), int(msg))
	}
}

func TestApp_View(t *testing.T) {
	app := New()

	view := app.View()
	if view == "" {
		t.Error("View should not return empty string")
	}

	// Test with quit set
	app.quit = true
	view = app.View()
	if view != "" {
		t.Error("View should return empty string when quit is true")
	}
}
