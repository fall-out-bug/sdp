package main

import (
	"fmt"
	"strings"

	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/spf13/cobra"
)

func ideaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "idea",
		Short: "Idea-phase evidence (plan events)",
		Long:  `Record @idea completion: emit plan event with question count, summary, optional Q&A pairs.`,
	}

	var featureID, summary string
	var questionCount int
	var wsID string
	var qaPairs []string

	recordCmd := &cobra.Command{
		Use:   "record",
		Short: "Record idea completion",
		Long:  `Emit plan event for @idea completion (questions, summary, qa_pairs). Use --qa "question|answer" for each pair.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if featureID == "" {
				return fmt.Errorf("--feature required")
			}
			if wsID == "" {
				wsID = "00-000-00"
			}
			var pairs []evidence.QAPair
			for _, s := range qaPairs {
				parts := strings.SplitN(s, "|", 2)
				q, a := "", ""
				if len(parts) >= 1 {
					q = strings.TrimSpace(parts[0])
				}
				if len(parts) >= 2 {
					a = strings.TrimSpace(parts[1])
				}
				pairs = append(pairs, evidence.QAPair{Q: q, A: a})
			}
			if evidence.Enabled() {
				ev := evidence.PlanEventForIdea(wsID, featureID, questionCount, summary, pairs)
				if err := evidence.EmitSync(ev); err != nil {
					return err
				}
			}
			fmt.Printf("Idea recorded: %s, %d questions\n", featureID, questionCount)
			return nil
		},
	}
	recordCmd.Flags().StringVar(&featureID, "feature", "", "Feature ID (e.g. F056)")
	recordCmd.Flags().IntVar(&questionCount, "questions", 0, "Number of questions")
	recordCmd.Flags().StringVar(&summary, "summary", "", "Requirements summary")
	recordCmd.Flags().StringVar(&wsID, "ws-id", "00-000-00", "WS ID for event")
	recordCmd.Flags().StringArrayVar(&qaPairs, "qa", nil, "question|answer (repeatable)")
	cmd.AddCommand(recordCmd)
	return cmd
}
