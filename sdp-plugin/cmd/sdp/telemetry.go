package main

import (
	"github.com/spf13/cobra"
)

var telemetryCmd = &cobra.Command{
	Use:   "telemetry",
	Short: "Manage telemetry collection",
	Long: `Manage telemetry collection for SDP.

Telemetry tracks anonymized usage metrics to help improve SDP:
  - Command invocations
  - Execution duration
  - Success/failure rates

🔒 Privacy Policy:
  - No PII (names, emails, usernames) collected
  - No data transmitted remotely (stored locally)
  - Opt-out available: sdp telemetry disable
  - Auto-cleanup after 90 days
  - See docs/PRIVACY.md for details

All data is stored locally in ~/.sdp/telemetry.jsonl`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default to status if no subcommand
		return telemetryStatusCmd.RunE(cmd, args)
	},
}

func init() {
	telemetryCmd.AddCommand(telemetryStatusCmd)
	telemetryCmd.AddCommand(telemetryExportCmd)
	telemetryCmd.AddCommand(telemetryUploadCmd)
	telemetryCmd.AddCommand(telemetryDisableCmd)
	telemetryCmd.AddCommand(telemetryEnableCmd)
	telemetryCmd.AddCommand(telemetryAnalyzeCmd)
	telemetryCmd.AddCommand(telemetryConsentCmd)
}
