package dashboard

import (
	"github.com/charmbracelet/lipgloss"
)

// Matrix-themed color palette
var (
	// Matrix style - black background with green text
	matrixBackground = lipgloss.Color("0")   // Black
	matrixForeground = lipgloss.Color("82")  // Bright green
	matrixAccent     = lipgloss.Color("46")  // Bright green accent
	matrixHighlight  = lipgloss.Color("226") // Yellow for selection
	matrixError      = lipgloss.Color("196") // Red for errors/blockers

	// Base style with matrix theme
	matrixBaseStyle = lipgloss.NewStyle().
			Foreground(matrixForeground).
			Background(matrixBackground)

	// Header style
	matrixHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(matrixAccent).
				Background(matrixBackground).
				Underline(true).
				Width(80)

	// Active tab style
	matrixActiveTabStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(matrixHighlight).
				Background(matrixBackground).
				Underline(true)

	// Inactive tab style - BRIGHT green (not dim)
	matrixInactiveTabStyle = lipgloss.NewStyle().
				Foreground(matrixForeground).
				Background(matrixBackground)

	// Selected item style (for arrow key navigation)
	matrixSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(matrixHighlight).
				Background(matrixBackground)

	// NO dim style - everything is bright green
	matrixBrightStyle = lipgloss.NewStyle().
				Foreground(matrixAccent).
				Background(matrixBackground)

	// Status colors (matrix theme) - ALL BRIGHT
	statusOpenMatrixStyle = lipgloss.NewStyle().
				Foreground(matrixForeground).
				Background(matrixBackground)

	statusInProgressMatrixStyle = lipgloss.NewStyle().
					Foreground(matrixAccent).
					Background(matrixBackground) // Bright green

	statusCompletedMatrixStyle = lipgloss.NewStyle().
					Foreground(matrixForeground).
					Background(matrixBackground) // Bright green (not dim!)

	statusBlockedMatrixStyle = lipgloss.NewStyle().
					Foreground(matrixError).
					Background(matrixBackground).
					Bold(true)

	// Priority colors (matrix theme) - ALL BRIGHT
	priorityP0MatrixStyle = lipgloss.NewStyle().
				Foreground(matrixError).
				Background(matrixBackground).
				Bold(true)

	priorityP1MatrixStyle = lipgloss.NewStyle().
				Foreground(matrixHighlight).
				Background(matrixBackground).
				Bold(true)

	priorityP2MatrixStyle = lipgloss.NewStyle().
				Foreground(matrixForeground).
				Background(matrixBackground)

	priorityP3MatrixStyle = lipgloss.NewStyle().
				Foreground(matrixForeground).
				Background(matrixBackground) // Bright green (not dim!)

	// Footer style - BRIGHT GREEN (not dim)
	matrixFooterStyle = lipgloss.NewStyle().
				Foreground(matrixForeground).
				Background(matrixBackground)

	// Dim style for secondary information
	matrixDimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("242")).
			Background(matrixBackground)
)
