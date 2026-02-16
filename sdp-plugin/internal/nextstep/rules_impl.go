package nextstep

import "fmt"

// noGitRepoRule recommends git initialization when not in a repo.
type noGitRepoRule struct{}

func (r *noGitRepoRule) Evaluate(state ProjectState) *Recommendation {
	if !state.GitStatus.IsRepo {
		return &Recommendation{
			Command:    "git init",
			Reason:     "Initialize a git repository to enable SDP workflows",
			Confidence: 0.95,
			Category:   CategorySetup,
			Alternatives: []Alternative{
				{Command: "sdp init", Reason: "Initialize SDP configuration"},
			},
		}
	}
	return nil
}

// failedWorkstreamRule recommends recovery after a failure.
type failedWorkstreamRule struct{}

func (r *failedWorkstreamRule) Evaluate(state ProjectState) *Recommendation {
	for _, ws := range state.Workstreams {
		if ws.Status == StatusFailed {
			return &Recommendation{
				Command:    fmt.Sprintf("sdp debug --ws %s", ws.ID),
				Reason:     fmt.Sprintf("Workstream %s failed: %s", ws.ID, truncateError(ws.LastError)),
				Confidence: 0.9,
				Category:   CategoryRecovery,
				Alternatives: []Alternative{
					{Command: fmt.Sprintf("sdp apply --retry --ws %s", ws.ID), Reason: "Retry the workstream"},
					{Command: "sdp status", Reason: "Check current state before debugging"},
				},
				Metadata: map[string]any{
					"failed_workstream": ws.ID,
					"error":             ws.LastError,
				},
			}
		}
	}

	if state.LastCommandError != "" {
		return &Recommendation{
			Command:    "sdp doctor",
			Reason:     fmt.Sprintf("Last command failed: %s", truncateError(state.LastCommandError)),
			Confidence: 0.85,
			Category:   CategoryRecovery,
			Alternatives: []Alternative{
				{Command: "sdp status", Reason: "Check current project state"},
			},
		}
	}

	return nil
}

// uncommittedChangesRule recommends committing when there are uncommitted changes.
type uncommittedChangesRule struct{}

func (r *uncommittedChangesRule) Evaluate(state ProjectState) *Recommendation {
	if state.GitStatus.Uncommitted {
		return &Recommendation{
			Command:    "git status",
			Reason:     "You have uncommitted changes - review before proceeding",
			Confidence: 0.8,
			Category:   CategoryInformation,
			Alternatives: []Alternative{
				{Command: "git diff", Reason: "Review the changes"},
				{Command: "git add . && git commit", Reason: "Commit the changes"},
			},
		}
	}
	return nil
}

// inProgressWorkstreamRule recommends continuing or checking in-progress work.
type inProgressWorkstreamRule struct{}

func (r *inProgressWorkstreamRule) Evaluate(state ProjectState) *Recommendation {
	for _, ws := range state.Workstreams {
		if ws.Status == StatusInProgress {
			return &Recommendation{
				Command:    "sdp status",
				Reason:     fmt.Sprintf("Workstream %s is in progress - check current state", ws.ID),
				Confidence: 0.85,
				Category:   CategoryInformation,
				Alternatives: []Alternative{
					{Command: fmt.Sprintf("sdp checkpoint resume %s", ws.ID), Reason: "Resume from checkpoint"},
				},
				Metadata: map[string]any{
					"active_workstream": ws.ID,
					"feature_id":        ws.Feature,
				},
			}
		}
	}
	return nil
}

// truncateError truncates an error message for display.
func truncateError(err string) string {
	if len(err) > 50 {
		return err[:47] + "..."
	}
	if err == "" {
		return "unknown error"
	}
	return err
}
