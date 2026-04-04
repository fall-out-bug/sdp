package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestShouldAskForTelemetryConsent(t *testing.T) {
	tests := []struct {
		name       string
		auto       bool
		headless   bool
		wantPrompt bool
	}{
		{
			name:       "interactive command prompts",
			wantPrompt: true,
		},
		{
			name:       "auto init skips prompt",
			auto:       true,
			wantPrompt: false,
		},
		{
			name:       "headless init skips prompt",
			headless:   true,
			wantPrompt: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{Use: "init"}
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
