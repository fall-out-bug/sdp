package nextstep

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	gitwrap "github.com/fall-out-bug/sdp/internal/git"
	"github.com/fall-out-bug/sdp/internal/parser"
	"github.com/fall-out-bug/sdp/internal/session"
)

// StateCollector gathers project state from various sources.
type StateCollector struct {
	projectRoot string
}

// NewStateCollector creates a new state collector.
func NewStateCollector(projectRoot string) *StateCollector {
	return &StateCollector{projectRoot: projectRoot}
}

// Collect gathers the current project state for recommendation.
func (c *StateCollector) Collect() (ProjectState, error) {
	state := ProjectState{
		Workstreams: c.collectWorkstreams(),
		GitStatus:   c.collectGitStatus(),
		Config:      c.collectConfig(),
		Mode:        ModeDrive, // Default mode
	}
	state.Session = c.collectSession()
	if state.Session != nil {
		state.ActiveWorkstream = state.Session.WorkstreamID
	}

	return state, nil
}

// collectWorkstreams gathers workstream status from docs/workstreams/.
func (c *StateCollector) collectWorkstreams() []WorkstreamStatus {
	parsed := []*parser.Workstream{}

	// Search for workstream files
	patterns := []string{
		filepath.Join(c.projectRoot, "docs", "workstreams", "*", "backlog", "*.md"),
		filepath.Join(c.projectRoot, "docs", "workstreams", "backlog", "*.md"),
	}

	seen := make(map[string]bool)

	for _, pattern := range patterns {
		files, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}

		for _, file := range files {
			ws, err := parser.ParseWorkstream(file)
			if err != nil {
				continue
			}

			// Avoid duplicates
			if seen[ws.ID] {
				continue
			}
			seen[ws.ID] = true
			parsed = append(parsed, ws)
		}
	}

	statuses := make(map[string]WorkstreamState, len(parsed))
	for _, ws := range parsed {
		statuses[ws.ID] = mapStatus(ws.Status)
	}
	workstreams := make([]WorkstreamStatus, 0, len(parsed))

	for _, ws := range parsed {
		status := WorkstreamStatus{
			ID:        ws.ID,
			Feature:   ws.Feature,
			Status:    mapStatus(ws.Status),
			Priority:  ws.Priority,
			Size:      ws.Size,
			BlockedBy: extractUnmetBlockers(ws, statuses),
		}

		workstreams = append(workstreams, status)
	}

	return workstreams
}

// mapStatus converts string status to WorkstreamState.
func mapStatus(status string) WorkstreamState {
	switch strings.ToLower(status) {
	case "backlog":
		return StatusBacklog
	case "ready", "open":
		return StatusReady
	case "in_progress", "in-progress", "started":
		return StatusInProgress
	case "blocked":
		return StatusBlocked
	case "completed", "done":
		return StatusCompleted
	case "failed", "error":
		return StatusFailed
	default:
		return StatusBacklog
	}
}

func extractUnmetBlockers(ws *parser.Workstream, statuses map[string]WorkstreamState) []string {
	if len(ws.DependsOn) == 0 {
		return nil
	}

	blockers := make([]string, 0, len(ws.DependsOn))
	for _, dep := range ws.DependsOn {
		status, ok := statuses[dep]
		if !ok || status == StatusCompleted {
			continue
		}
		blockers = append(blockers, dep)
	}
	if len(blockers) == 0 {
		return nil
	}
	return blockers
}

// collectGitStatus gathers git repository status.
func (c *StateCollector) collectGitStatus() GitStatusInfo {
	info := GitStatusInfo{
		IsRepo:     false,
		MainBranch: "main",
	}

	gitDir := filepath.Join(c.projectRoot, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		info.IsRepo = true
	}

	if info.IsRepo {
		if branch, err := gitwrap.GetCurrentBranch(c.projectRoot); err == nil {
			info.Branch = branch
		}

		if out, err := c.runGit("status", "--porcelain"); err == nil {
			info.Uncommitted = strings.TrimSpace(string(out)) != ""
		}

		if _, err := c.runGit("rev-parse", "--verify", "master"); err == nil {
			info.MainBranch = "master"
		}

		if out, err := c.runGit("rev-list", "--left-right", "--count", "HEAD...@{u}"); err == nil {
			parts := strings.Fields(strings.TrimSpace(string(out)))
			if len(parts) == 2 {
				info.UpstreamDiverg = parts[0] != "0" && parts[1] != "0"
			}
		}
	}

	return info
}

func (c *StateCollector) collectSession() *SessionState {
	if session.Exists(c.projectRoot) {
		if current, err := session.Load(c.projectRoot); err == nil {
			return &SessionState{
				FeatureID:      current.FeatureID,
				WorktreePath:   current.WorktreePath,
				ExpectedBranch: current.ExpectedBranch,
				ExpectedRemote: current.ExpectedRemote,
				Source:         "session",
			}
		}
	}

	return c.collectLegacySession()
}

func (c *StateCollector) collectLegacySession() *SessionState {
	path := filepath.Join(c.projectRoot, ".sdp", "session.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil
	}

	state := &SessionState{
		WorkstreamID:   stringFromMap(payload, "workstream_id"),
		FeatureID:      stringFromMap(payload, "feature_id"),
		WorktreePath:   stringFromMap(payload, "worktree_path"),
		ExpectedBranch: stringFromMap(payload, "branch"),
		ExpectedRemote: stringFromMap(payload, "expected_remote"),
		Source:         "legacy",
	}

	if state.WorkstreamID == "" && state.FeatureID == "" && state.WorktreePath == "" && state.ExpectedBranch == "" {
		return nil
	}

	return state
}

func stringFromMap(payload map[string]any, key string) string {
	value, ok := payload[key]
	if !ok {
		return ""
	}
	text, ok := value.(string)
	if !ok {
		return ""
	}
	return text
}

func (c *StateCollector) runGit(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = c.projectRoot
	return cmd.Output()
}

// collectConfig gathers SDP configuration.
func (c *StateCollector) collectConfig() ConfigInfo {
	info := ConfigInfo{
		HasSDPConfig: false,
	}

	configPath := filepath.Join(c.projectRoot, ".sdp", "config.yml")
	if _, err := os.Stat(configPath); err == nil {
		info.HasSDPConfig = true
	}

	configJSONPath := filepath.Join(c.projectRoot, ".sdp", "config.json")
	if _, err := os.Stat(configJSONPath); err == nil {
		info.HasSDPConfig = true
	}

	// Check for evidence log
	evidencePath := filepath.Join(c.projectRoot, ".sdp", "log", "events.jsonl")
	if _, err := os.Stat(evidencePath); err == nil {
		info.EvidenceEnabled = true
	}

	info.ProjectRoot = c.projectRoot

	return info
}
