package main

import (
	"fmt"
	"os"
)

// ProjectStatus holds the project status information
type ProjectStatus struct {
	HasGit        bool              `json:"has_git"`
	HasClaude     bool              `json:"has_claude"`
	HasSDP        bool              `json:"has_sdp"`
	HasBeads      bool              `json:"has_beads"`
	Workstreams   WorkstreamSummary `json:"workstreams"`
	ActiveSession *ActiveSession    `json:"active_session,omitempty"`
	NextAction    string            `json:"next_action"`
}

// WorkstreamSummary holds workstream counts
type WorkstreamSummary struct {
	Open       int `json:"open"`
	InProgress int `json:"in_progress"`
	Completed  int `json:"completed"`
	Blocked    int `json:"blocked"`
}

// ActiveSession holds current session info
type ActiveSession struct {
	WorkstreamID string `json:"workstream_id,omitempty"`
	FeatureID    string `json:"feature_id,omitempty"`
	Branch       string `json:"branch,omitempty"`
}

// runTextStatus generates text or JSON status output
func runTextStatus(jsonOutput bool) error {
	status := gatherProjectStatus()

	if jsonOutput {
		return status.printJSON()
	}
	return status.printText()
}

// gatherProjectStatus collects project status information
func gatherProjectStatus() *ProjectStatus {
	status := &ProjectStatus{
		HasGit:    dirExists(".git"),
		HasClaude: dirExists(".claude"),
		HasSDP:    dirExists(".sdp"),
		HasBeads:  fileExists(".beads/issues.jsonl"),
	}

	status.Workstreams = countWorkstreams()
	status.ActiveSession = getActiveSession()
	status.NextAction = determineNextAction(status)

	return status
}

// countWorkstreams counts workstreams by status
func countWorkstreams() WorkstreamSummary {
	summary := WorkstreamSummary{}

	backlogDir := "docs/workstreams/backlog"
	if entries, err := os.ReadDir(backlogDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if len(name) >= 12 && name[2] == '-' && name[6] == '-' {
				summary.Open++
			}
		}
	}

	completedDir := "docs/workstreams/completed"
	if entries, err := os.ReadDir(completedDir); err == nil {
		summary.Completed = len(entries)
	}

	return summary
}

// getActiveSession checks for an active session
func getActiveSession() *ActiveSession {
	sessionFile := ".sdp/session.json"
	data, err := os.ReadFile(sessionFile)
	if err != nil {
		return nil
	}

	session := &ActiveSession{}
	content := string(data)

	if idx := findJSONField(content, "workstream_id"); idx >= 0 {
		session.WorkstreamID = extractJSONString(content, idx)
	}
	if idx := findJSONField(content, "feature_id"); idx >= 0 {
		session.FeatureID = extractJSONString(content, idx)
	}
	if idx := findJSONField(content, "branch"); idx >= 0 {
		session.Branch = extractJSONString(content, idx)
	}

	if session.WorkstreamID != "" || session.FeatureID != "" {
		return session
	}
	return nil
}

// determineNextAction suggests the next action based on status
func determineNextAction(status *ProjectStatus) string {
	if !status.HasGit {
		return "Run 'git init' to initialize a repository"
	}
	if !status.HasClaude {
		return "Run 'sdp init --guided' to set up SDP"
	}
	if !status.HasSDP {
		return "Run 'sdp init' to create .sdp directory"
	}
	if status.ActiveSession != nil && status.ActiveSession.WorkstreamID != "" {
		return fmt.Sprintf("Run 'sdp apply --ws %s' to continue work", status.ActiveSession.WorkstreamID)
	}
	if status.Workstreams.Open > 0 {
		return "Run 'sdp status' to see open workstreams"
	}
	return "Run 'sdp plan \"your feature idea\"' to plan a feature"
}

// printText outputs status as human-readable text
func (s *ProjectStatus) printText() error {
	fmt.Println("SDP Project Status")
	fmt.Println("=================")
	fmt.Println()

	fmt.Println("Environment:")
	fmt.Printf("  Git:     %s\n", boolIcon(s.HasGit))
	fmt.Printf("  Claude:  %s\n", boolIcon(s.HasClaude))
	fmt.Printf("  SDP:     %s\n", boolIcon(s.HasSDP))
	fmt.Printf("  Beads:   %s\n", boolIcon(s.HasBeads))
	fmt.Println()

	fmt.Println("Workstreams:")
	fmt.Printf("  Open:        %d\n", s.Workstreams.Open)
	fmt.Printf("  In Progress: %d\n", s.Workstreams.InProgress)
	fmt.Printf("  Completed:   %d\n", s.Workstreams.Completed)
	fmt.Printf("  Blocked:     %d\n", s.Workstreams.Blocked)
	fmt.Println()

	if s.ActiveSession != nil {
		fmt.Println("Active Session:")
		if s.ActiveSession.WorkstreamID != "" {
			fmt.Printf("  Workstream: %s\n", s.ActiveSession.WorkstreamID)
		}
		if s.ActiveSession.FeatureID != "" {
			fmt.Printf("  Feature:    %s\n", s.ActiveSession.FeatureID)
		}
		if s.ActiveSession.Branch != "" {
			fmt.Printf("  Branch:     %s\n", s.ActiveSession.Branch)
		}
		fmt.Println()
	}

	fmt.Println("Next Action:")
	fmt.Printf("  %s\n", s.NextAction)

	return nil
}

// printJSON outputs status as JSON
func (s *ProjectStatus) printJSON() error {
	fmt.Printf(`{"has_git":%t,"has_claude":%t,"has_sdp":%t,"has_beads":%t,`,
		s.HasGit, s.HasClaude, s.HasSDP, s.HasBeads)
	fmt.Printf(`"workstreams":{"open":%d,"in_progress":%d,"completed":%d,"blocked":%d},`,
		s.Workstreams.Open, s.Workstreams.InProgress, s.Workstreams.Completed, s.Workstreams.Blocked)
	fmt.Printf(`"next_action":%q}`, s.NextAction)
	fmt.Println()
	return nil
}
