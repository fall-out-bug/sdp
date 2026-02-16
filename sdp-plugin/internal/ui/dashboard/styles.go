package dashboard

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles for dashboard elements
var (
	// Header style
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")). // Light blue
			Width(80)

	// Active tab style
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("226")). // Yellow
			Underline(true)

	// Inactive tab style
	inactiveTabStyle = lipgloss.NewStyle().
				Faint(true)

	// Status colors
	statusOpenStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("76"))  // Green
	statusInProgressStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("226")) // Yellow
	statusCompletedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))  // Blue
	statusBlockedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // Red

	// Priority colors
	priorityP0Style = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true) // Red
	priorityP1Style = lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true) // Orange
	priorityP2Style = lipgloss.NewStyle().Foreground(lipgloss.Color("226"))            // Yellow
	priorityP3Style = lipgloss.NewStyle().Faint(true)                                  // Gray

	// Footer style
	footerStyle = lipgloss.NewStyle().
			Faint(true).
			Width(80)
)

// StatusStyle returns the style for a given status
func StatusStyle(status string) lipgloss.Style {
	switch status {
	case "open":
		return statusOpenStyle
	case "in_progress", "in-progress":
		return statusInProgressStyle
	case "completed":
		return statusCompletedStyle
	case "blocked":
		return statusBlockedStyle
	default:
		return lipgloss.NewStyle()
	}
}

// PriorityStyle returns the style for a given priority
func PriorityStyle(priority string) lipgloss.Style {
	switch priority {
	case "P0":
		return priorityP0Style
	case "P1":
		return priorityP1Style
	case "P2":
		return priorityP2Style
	case "P3":
		return priorityP3Style
	default:
		return lipgloss.NewStyle()
	}
}
