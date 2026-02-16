package nextstep

import (
	"fmt"
	"sort"
)

// RecommendationRule defines a rule for generating recommendations.
type RecommendationRule interface {
	Evaluate(state ProjectState) *Recommendation
}

// defaultRules returns the ordered list of default recommendation rules.
func defaultRules() []RecommendationRule {
	return []RecommendationRule{
		&noGitRepoRule{},
		&failedWorkstreamRule{},
		&uncommittedChangesRule{},
		&inProgressWorkstreamRule{},
		&readyWorkstreamRule{},
		&blockedWorkstreamRule{},
		&allCompleteRule{},
		&freshProjectRule{},
	}
}

// readyWorkstreamRule recommends executing ready workstreams.
type readyWorkstreamRule struct{}

func (r *readyWorkstreamRule) Evaluate(state ProjectState) *Recommendation {
	ready := []WorkstreamStatus{}
	for _, ws := range state.Workstreams {
		if ws.Status == StatusReady && len(ws.BlockedBy) == 0 {
			ready = append(ready, ws)
		}
	}

	if len(ready) == 0 {
		return nil
	}

	sort.Slice(ready, func(i, j int) bool {
		return ComparePriority(ready[i], ready[j]) < 0
	})

	ws := ready[0]
	return &Recommendation{
		Command:    fmt.Sprintf("sdp apply --ws %s", ws.ID),
		Reason:     fmt.Sprintf("Ready to execute workstream %s (%s)", ws.ID, ws.Feature),
		Confidence: 0.95,
		Category:   CategoryExecution,
		Alternatives: []Alternative{
			{Command: "sdp status", Reason: "View all ready workstreams"},
		},
		Metadata: map[string]any{
			"workstream_id": ws.ID,
			"feature_id":    ws.Feature,
			"priority":      ws.Priority,
		},
	}
}

// blockedWorkstreamRule recommends resolving blockers.
type blockedWorkstreamRule struct{}

func (r *blockedWorkstreamRule) Evaluate(state ProjectState) *Recommendation {
	for _, ws := range state.Workstreams {
		if ws.Status == StatusReady && len(ws.BlockedBy) > 0 {
			blockerID := ws.BlockedBy[0]
			return &Recommendation{
				Command:    fmt.Sprintf("sdp apply --ws %s", blockerID),
				Reason:     fmt.Sprintf("Complete blocker %s to unblock %s", blockerID, ws.ID),
				Confidence: 0.9,
				Category:   CategoryExecution,
				Alternatives: []Alternative{
					{Command: "sdp status", Reason: "View dependency status"},
				},
				Metadata: map[string]any{
					"blocked_workstream": ws.ID,
					"blocker":            blockerID,
				},
			}
		}
	}
	return nil
}

// allCompleteRule recommends review/deploy when all workstreams are done.
type allCompleteRule struct{}

func (r *allCompleteRule) Evaluate(state ProjectState) *Recommendation {
	if len(state.Workstreams) == 0 {
		return nil
	}

	allComplete := true
	featureID := ""
	for _, ws := range state.Workstreams {
		if ws.Status != StatusCompleted {
			allComplete = false
			break
		}
		if featureID == "" {
			featureID = ws.Feature
		}
	}

	if allComplete && featureID != "" {
		return &Recommendation{
			Command:    fmt.Sprintf("sdp review %s", featureID),
			Reason:     fmt.Sprintf("All workstreams complete for %s - review before deployment", featureID),
			Confidence: 0.95,
			Category:   CategoryPlanning,
			Alternatives: []Alternative{
				{Command: fmt.Sprintf("sdp deploy %s", featureID), Reason: "Deploy if review passed"},
				{Command: "sdp status", Reason: "Check overall project state"},
			},
			Metadata: map[string]any{
				"feature_id": featureID,
			},
		}
	}

	return nil
}

// freshProjectRule recommends initial setup for fresh projects.
type freshProjectRule struct{}

func (r *freshProjectRule) Evaluate(state ProjectState) *Recommendation {
	if len(state.Workstreams) == 0 && state.ActiveWorkstream == "" {
		if state.Config.HasSDPConfig {
			return &Recommendation{
				Command:    "sdp doctor",
				Reason:     "Start by checking your environment setup",
				Confidence: 0.9,
				Category:   CategorySetup,
				Alternatives: []Alternative{
					{Command: "sdp status", Reason: "View current project state"},
					{Command: "@vision \"your project idea\"", Reason: "Start strategic planning"},
				},
			}
		}
		return &Recommendation{
			Command:    "sdp init",
			Reason:     "Initialize SDP in this project",
			Confidence: 0.95,
			Category:   CategorySetup,
		}
	}
	return nil
}
