package dashboard

import (
	"testing"
)

// TestRenderNextStep tests rendering the next step block.
func TestRenderNextStep(t *testing.T) {
	app := New()
	app.state.Loading = false
	app.state.NextStep = NextStepInfo{
		Command:    "sdp apply --ws 00-069-01",
		Reason:     "Ready to execute workstream",
		Confidence: 0.95,
		Category:   "execution",
	}

	output := app.renderNextStep()
	if output == "" {
		t.Error("Expected non-empty next step output")
	}

	// Should contain the command
	if len(output) < 10 {
		t.Error("Output seems too short")
	}
}

// TestRenderNextStepEmpty tests rendering with empty next step.
func TestRenderNextStepEmpty(t *testing.T) {
	app := New()
	app.state.Loading = false
	app.state.NextStep = NextStepInfo{}

	output := app.renderNextStep()
	if output != "" {
		t.Error("Expected empty output for empty next step")
	}
}

// TestRenderNextStepCategories tests different category styling.
func TestRenderNextStepCategories(t *testing.T) {
	categories := []struct {
		category string
	}{
		{"execution"},
		{"recovery"},
		{"planning"},
		{"setup"},
		{"information"},
	}

	for _, tc := range categories {
		t.Run(tc.category, func(t *testing.T) {
			app := New()
			app.state.Loading = false
			app.state.NextStep = NextStepInfo{
				Command:    "sdp test",
				Reason:     "Test reason",
				Confidence: 0.5,
				Category:   tc.category,
			}

			output := app.renderNextStep()
			if output == "" {
				t.Error("Expected non-empty output")
			}
		})
	}
}

// TestRenderNextStepConfidence tests confidence level styling.
func TestRenderNextStepConfidence(t *testing.T) {
	levels := []struct {
		name       string
		confidence float64
	}{
		{"high", 0.95},
		{"medium", 0.7},
		{"low", 0.3},
	}

	for _, tc := range levels {
		t.Run(tc.name, func(t *testing.T) {
			app := New()
			app.state.Loading = false
			app.state.NextStep = NextStepInfo{
				Command:    "sdp test",
				Reason:     "Test reason",
				Confidence: tc.confidence,
				Category:   "execution",
			}

			output := app.renderNextStep()
			if output == "" {
				t.Error("Expected non-empty output")
			}
		})
	}
}
