package dashboard

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
