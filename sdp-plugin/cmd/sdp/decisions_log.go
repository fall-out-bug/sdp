package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fall-out-bug/sdp/internal/decision"
	"github.com/fall-out-bug/sdp/internal/evidence"
	"github.com/spf13/cobra"
)

// decisionsLogCmd logs a new decision
func decisionsLogCmd() *cobra.Command {
	var decisionType, featureID, workstreamID, question, decisionStr, rationale, alternatives, outcome, maker, reverses string
	var tags []string

	cmd := &cobra.Command{
		Use:   "log",
		Short: "Log a new decision",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate required fields
			if question == "" || decisionStr == "" {
				return fmt.Errorf("required flags: --question, --decision")
			}

			// Validate field lengths
			if err := validateFieldLength("question", question, maxFieldLength); err != nil {
				return err
			}
			if err := validateFieldLength("decision", decisionStr, maxFieldLength); err != nil {
				return err
			}
			if rationale != "" {
				if err := validateFieldLength("rationale", rationale, maxFieldLength); err != nil {
					return err
				}
			}
			if alternatives != "" {
				if err := validateFieldLength("alternatives", alternatives, maxFieldLength); err != nil {
					return err
				}
			}

			// Strip control characters
			question = stripControlChars(question)
			decisionStr = stripControlChars(decisionStr)
			rationale = stripControlChars(rationale)
			alternatives = stripControlChars(alternatives)

			// Validate decision type
			if decisionType != "" {
				validTypes := []string{
					decision.DecisionTypeVision,
					decision.DecisionTypeTechnical,
					decision.DecisionTypeTradeoff,
					decision.DecisionTypeExplicit,
				}
				valid := false
				for _, t := range validTypes {
					if decisionType == t {
						valid = true
						break
					}
				}
				if !valid {
					return fmt.Errorf("invalid decision type %q, must be one of: vision, technical, tradeoff, explicit", decisionType)
				}
			}

			root, err := findProjectRoot()
			if err != nil {
				return err
			}

			logger, err := decision.NewLogger(root)
			if err != nil {
				return err
			}

			// Parse alternatives (comma-separated)
			var altList []string
			if alternatives != "" {
				altList = strings.Split(alternatives, ",")
				for i := range altList {
					altList[i] = strings.TrimSpace(altList[i])
				}
			}

			// Parse tags (comma-separated)
			var tagList []string
			if len(tags) > 0 {
				tagList = tags
			}

			// Default maker to "user" if not specified
			if maker == "" {
				maker = "user"
			}

			// Create decision (bridge in logger.Log emits to evidence log)
			d := decision.Decision{
				Question:      question,
				Decision:      decisionStr,
				Rationale:     rationale,
				Type:          decisionType,
				FeatureID:     featureID,
				WorkstreamID:  workstreamID,
				Alternatives:  altList,
				Outcome:       outcome,
				DecisionMaker: maker,
				Tags:          tagList,
				Reverses:      reverses,
			}

			if err := logger.Log(d); err != nil {
				return err
			}
			// AC1: Emit decision event to evidence log (bridge in CLI to avoid decision->evidence cycle)
			var rev *string
			if reverses != "" {
				rev = &reverses
			}
			if err := evidence.EmitSync(evidence.DecisionEvent(workstreamID, question, decisionStr, rationale, altList, 0, tagList, rev)); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "warning: evidence emit: %v\n", err)
			}

			fmt.Printf("✓ Logged decision: %s\n", d.Decision)
			return nil
		},
	}

	cmd.Flags().StringVar(&decisionType, "type", decision.DecisionTypeExplicit, "Decision type")
	cmd.Flags().StringVar(&featureID, "feature-id", "", "Feature ID")
	cmd.Flags().StringVar(&workstreamID, "workstream-id", "", "Workstream ID")
	cmd.Flags().StringVar(&question, "question", "", "Question or problem")
	cmd.Flags().StringVar(&decisionStr, "decision", "", "Decision made")
	cmd.Flags().StringVar(&rationale, "rationale", "", "Rationale for decision")
	cmd.Flags().StringVar(&alternatives, "alternatives", "", "Alternatives considered (comma-separated)")
	cmd.Flags().StringVar(&outcome, "outcome", "", "Expected outcome")
	cmd.Flags().StringVar(&maker, "maker", "", "Decision maker (user/claude/system)")
	cmd.Flags().StringSliceVar(&tags, "tags", []string{}, "Tags for categorization")
	cmd.Flags().StringVar(&reverses, "reverses", "", "ID of previous decision being overturned (AC7)")

	return cmd
}
