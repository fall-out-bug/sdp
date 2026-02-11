package main

import (
	"fmt"
	"strings"

	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/spf13/cobra"
)

func designCmd() *cobra.Command {
	var featureID string
	var wsCount int
	var wsID string
	var metadataPairs []string

	cmd := &cobra.Command{
		Use:   "design",
		Short: "Design-phase evidence (plan events)",
		Long:  `Record @design completion: emit plan event with WS count and decomposition metadata.`,
	}

	recordCmd := &cobra.Command{
		Use:   "record",
		Short: "Record design completion",
		Long:  `Emit plan event for @design completion (ws count, feature_id, metadata).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if featureID == "" {
				return fmt.Errorf("--feature required")
			}
			if wsID == "" {
				wsID = "00-000-00"
			}
			metadata := make(map[string]interface{})
			for _, p := range metadataPairs {
				kv := strings.SplitN(p, "=", 2)
				if len(kv) == 2 {
					metadata[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
				}
			}
			if evidence.Enabled() {
				ev := evidence.PlanEventForDesign(wsID, featureID, wsCount, nil, metadata)
				if err := evidence.EmitSync(ev); err != nil {
					return err
				}
			}
			fmt.Printf("Design recorded: %s, %d workstreams\n", featureID, wsCount)
			return nil
		},
	}
	recordCmd.Flags().StringVar(&featureID, "feature", "", "Feature ID (e.g. F056)")
	recordCmd.Flags().IntVar(&wsCount, "ws-count", 0, "Number of workstreams")
	recordCmd.Flags().StringVar(&wsID, "ws-id", "00-000-00", "WS ID for event (default 00-000-00)")
	recordCmd.Flags().StringArrayVar(&metadataPairs, "metadata", nil, "key=value (repeatable)")
	cmd.AddCommand(recordCmd)
	return cmd
}
