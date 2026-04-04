package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestShouldAskForTelemetryConsent(t *testing.T) {
	tests := []struct {
		name       string
		use        string
		auto       bool
		headless   bool
		wantPrompt bool
	}{
		{
			name:       "interactive command prompts",
			use:        "plan",
			wantPrompt: true,
		},
		{
			name:       "auto init skips prompt",
			use:        "init",
			auto:       true,
			wantPrompt: false,
		},
		{
			name:       "headless init skips prompt",
			use:        "init",
			headless:   true,
			wantPrompt: false,
		},
		{
			name:       "doctor skips prompt",
			use:        "doctor",
			wantPrompt: false,
		},
		{
			name:       "status skips prompt",
			use:        "status",
			wantPrompt: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{Use: tt.use}
			cmd.Flags().Bool("auto", false, "")
			cmd.Flags().Bool("headless", false, "")

			if err := cmd.Flags().Set("auto", boolString(tt.auto)); err != nil {
				t.Fatalf("set auto flag: %v", err)
			}
			if err := cmd.Flags().Set("headless", boolString(tt.headless)); err != nil {
				t.Fatalf("set headless flag: %v", err)
			}

			if got := shouldAskForTelemetryConsent(cmd); got != tt.wantPrompt {
				t.Fatalf("shouldAskForTelemetryConsent() = %t, want %t", got, tt.wantPrompt)
			}
		})
	}
}

func TestShouldAskForTelemetryConsentNilCommand(t *testing.T) {
	if shouldAskForTelemetryConsent(nil) {
		t.Fatal("nil command should not request telemetry consent")
	}
}

func boolString(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
