package dashboard

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// renderWorkstreams renders the workstreams tab
func (a *App) renderWorkstreams() string {
	if len(a.state.Workstreams) == 0 {
		return "Workstreams\n\nNo workstreams found"
	}

	var content string
	content += matrixBaseStyle.Render("Workstreams\n\n")

	statusOrder := []string{"open", "in_progress", "completed", "blocked"}
	statusLabels := map[string]string{
		"open":        "Open",
		"in_progress": "In Progress",
		"completed":   "Completed",
		"blocked":     "Blocked",
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

	content += matrixBaseStyle.Render(fmt.Sprintf("Total: %d workstream(s)\n", totalCount))

	return content
}

// renderPriorityMatrix renders priority with matrix colors
func (a *App) renderPriorityMatrix(priority string) string {
	switch priority {
	case "P0":
		return priorityP0MatrixStyle.Render("[" + priority + "]")
	case "P1":
		return priorityP1MatrixStyle.Render("[" + priority + "]")
	case "P2":
		return priorityP2MatrixStyle.Render("[" + priority + "]")
	case "P3":
		return priorityP3MatrixStyle.Render("[" + priority + "]")
	default:
		return matrixBaseStyle.Render("[" + priority + "]")
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
		content += "    " + matrixBaseStyle.Render("Last modified: "+dateStr) + "\n\n"
	}

	return content
}

// renderTests renders the tests tab
func (a *App) renderTests() string {
	content := matrixBaseStyle.Render("Tests\n\n")

	tr := a.state.TestResults
	content += matrixBaseStyle.Render(fmt.Sprintf("Coverage: %s\n", tr.Coverage))
	content += matrixBaseStyle.Render(fmt.Sprintf("Tests: %d/%d passing\n", tr.Passing, tr.Total))
	content += matrixBaseStyle.Render(fmt.Sprintf("Last run: %s\n\n", tr.LastRun.Format("2006-01-02 15:04:05")))

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

// renderNextStep renders the next step recommendation block
func (a *App) renderNextStep() string {
	ns := a.state.NextStep
	if ns.Command == "" {
		return ""
	}

	var content string
	content += matrixBrightStyle.Render("Next Step") + "\n"

	// Confidence indicator
	confidenceStr := fmt.Sprintf("%.0f%%", ns.Confidence*100)
	if ns.Confidence >= 0.8 {
		confidenceStr = statusCompletedMatrixStyle.Render(confidenceStr)
	} else if ns.Confidence >= 0.5 {
		confidenceStr = statusInProgressMatrixStyle.Render(confidenceStr)
	} else {
		confidenceStr = matrixDimStyle.Render(confidenceStr)
	}

	// Category styling
	categoryStr := ns.Category
	switch categoryStr {
	case "execution":
		categoryStr = statusOpenMatrixStyle.Render(categoryStr)
	case "recovery":
		categoryStr = statusBlockedMatrixStyle.Render(categoryStr)
	case "planning":
		categoryStr = statusInProgressMatrixStyle.Render(categoryStr)
	case "setup":
		categoryStr = statusInProgressMatrixStyle.Render(categoryStr)
	default:
		categoryStr = matrixDimStyle.Render(categoryStr)
	}

	content += matrixBaseStyle.Render("  ") + statusCompletedMatrixStyle.Render(ns.Command) + "\n"
	content += matrixBaseStyle.Render("  "+ns.Reason) + "\n"
	content += matrixBaseStyle.Render("  ") + confidenceStr + " " + categoryStr + "\n"

	return content
}
