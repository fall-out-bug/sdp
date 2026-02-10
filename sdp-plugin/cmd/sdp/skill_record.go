package main

import (
	"fmt"
	"strings"

	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/spf13/cobra"
)

func skillRecord() *cobra.Command {
	var skillName, eventType, wsID string
	var dataPairs []string

	cmd := &cobra.Command{
		Use:   "record",
		Short: "Record skill execution in evidence log (F056)",
		Long:  `Emit a thin evidence event for a skill. Non-blocking; used by @vision, @reality, etc.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if skillName == "" || eventType == "" {
				return fmt.Errorf("--skill and --type required")
			}
			if wsID == "" {
				wsID = "00-000-00"
			}
			data := make(map[string]interface{})
			for _, p := range dataPairs {
				kv := strings.SplitN(p, "=", 2)
				if len(kv) == 2 {
					data[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
				}
			}
			_ = evidence.EmitSync(evidence.SkillEvent(skillName, eventType, wsID, data))
			fmt.Printf("Recorded: %s %s\n", skillName, eventType)
			return nil
		},
	}
	cmd.Flags().StringVar(&skillName, "skill", "", "Skill name (vision, reality, oneshot, prototype, hotfix, bugfix, issue, debug)")
	cmd.Flags().StringVar(&eventType, "type", "", "Event type (plan, verification, generation, approval)")
	cmd.Flags().StringVar(&wsID, "ws-id", "00-000-00", "WS ID for event")
	cmd.Flags().StringArrayVar(&dataPairs, "data", nil, "key=value (repeatable)")
	return cmd
}
