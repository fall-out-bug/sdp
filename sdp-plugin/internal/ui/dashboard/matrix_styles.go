package dashboard

import (
	"github.com/charmbracelet/lipgloss"
)

// Matrix-themed color palette
var (
	// Matrix style - black background with green text
	matrixBackground = lipgloss.Color("0")     // Black
	matrixForeground = lipgloss.Color("34")    // Green (bright)
	matrixDim       = lipgloss.Color("28")    // Dark green
	matrixAccent    = lipgloss.Color("46")    // Bright green accent
	matrixHighlight = lipgloss.Color("226")   // Yellow for selection
	matrixError     = lipgloss.Color("196")   // Red for errors/blockers

	// Base style with matrix theme
	matrixBaseStyle = lipgloss.NewStyle().
			Foreground(matrixForeground).
			Background(matrixBackground)

	// Header style
	matrixHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(matrixAccent).
			Background(matrixBackground).
			Underline(true)

	// Active tab style
	matrixActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(matrixHighlight).
			Background(matrixBackground).
			Underline(true)

	// Inactive tab style
	matrixInactiveTabStyle = lipgloss.NewStyle().
			Foreground(matrixDim).
			Background(matrixBackground).
			Faint(true)

	// Selected item style (for arrow key navigation)
	matrixSelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(matrixHighlight).
			Background(matrixBackground)

	// Dim style for secondary text
	matrixDimStyle = lipgloss.NewStyle().
			Foreground(matrixDim).
			Background(matrixBackground)

	// Status colors (matrix theme)
	statusOpenMatrixStyle = lipgloss.NewStyle().
			Foreground(matrixForeground).
			Background(matrixBackground)

	statusInProgressMatrixStyle = lipgloss.NewStyle().
			Foreground(matrixAccent).
			Background(matrixBackground) // Bright green

	statusCompletedMatrixStyle = lipgloss.NewStyle().
			Foreground(matrixDim).
			Background(matrixBackground) // Dark green

	statusBlockedMatrixStyle = lipgloss.NewStyle().
			Foreground(matrixError).
			Background(matrixBackground).
			Bold(true)

	// Priority colors (matrix theme)
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
			Foreground(matrixDim).
			Background(matrixBackground).
			Faint(true)

	// Footer style
	matrixFooterStyle = lipgloss.NewStyle().
			Foreground(matrixDim).
			Background(matrixBackground).
			Faint(true)
)
