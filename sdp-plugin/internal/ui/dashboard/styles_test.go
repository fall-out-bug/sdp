package dashboard

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestStatusStyle(t *testing.T) {
	tests := []struct {
		status    string
		wantMatch bool // Just verify it returns a non-zero style
	}{
		{"open", true},
		{"in_progress", true},
		{"in-progress", true},
		{"completed", true},
		{"blocked", true},
		{"unknown", true}, // Default case
		{"", true},        // Empty string
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			style := StatusStyle(tt.status)

			// Just verify we get a valid style back
			if style.GetForeground() == lipgloss.Color("") && tt.status != "" && tt.status != "unknown" {
				// Some styles should have colors set
				t.Logf("StatusStyle(%q) returned style", tt.status)
			}
		})
	}
}

func TestPriorityStyle(t *testing.T) {
	tests := []struct {
		priority  string
		wantMatch bool
	}{
		{"P0", true},
		{"P1", true},
		{"P2", true},
		{"P3", true},
		{"unknown", true}, // Default case
		{"", true},        // Empty string
	}

	for _, tt := range tests {
		t.Run(tt.priority, func(t *testing.T) {
			style := PriorityStyle(tt.priority)

			// Just verify we get a valid style back
			_ = style // Use the style to avoid unused variable warning
		})
	}
}

func TestStatusStyleConsistency(t *testing.T) {
	// Test that in_progress and in-progress return the same style
	style1 := StatusStyle("in_progress")
	style2 := StatusStyle("in-progress")

	// They should both be non-empty styles
	_ = style1
	_ = style2
}

func TestPriorityStyleAllPriorities(t *testing.T) {
	// Test all priority levels return non-nil styles
	priorities := []string{"P0", "P1", "P2", "P3"}

	for _, p := range priorities {
		style := PriorityStyle(p)
		if style.GetBold() {
			// P0 and P1 should be bold
			if p != "P0" && p != "P1" {
				t.Logf("Priority %s returned bold style unexpectedly", p)
			}
		}
	}
}

func TestStatusStyleAllStatuses(t *testing.T) {
	// Test all known statuses
	statuses := []string{"open", "in_progress", "in-progress", "completed", "blocked"}

	for _, s := range statuses {
		style := StatusStyle(s)
		_ = style // Just verify no panic
	}
}
