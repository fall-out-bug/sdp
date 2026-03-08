package nextstep

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
)

type StatusView struct {
	Version       string          `json:"version"`
	HasGit        bool            `json:"has_git"`
	HasClaude     bool            `json:"has_claude"`
	HasSDP        bool            `json:"has_sdp"`
	HasBeads      bool            `json:"has_beads"`
	Environment   EnvironmentView `json:"environment"`
	Workstreams   WorkstreamView  `json:"workstreams"`
	ActiveSession *SessionState   `json:"active_session,omitempty"`
	NextAction    string          `json:"next_action"`
	NextStep      *Recommendation `json:"next_step,omitempty"`
}

type EnvironmentView struct {
	ProjectRoot    string `json:"project_root,omitempty"`
	GitBranch      string `json:"git_branch,omitempty"`
	MainBranch     string `json:"main_branch,omitempty"`
	Uncommitted    bool   `json:"uncommitted"`
	HasSDPConfig   bool   `json:"has_sdp_config"`
	HasSDPDir      bool   `json:"has_sdp_dir"`
	HasClaudeDir   bool   `json:"has_claude_dir"`
	HasBeadsIssues bool   `json:"has_beads_issues"`
	HasEvidenceLog bool   `json:"has_evidence_log"`
}

type WorkstreamView struct {
	Total      int                `json:"total"`
	Backlog    int                `json:"backlog"`
	Completed  int                `json:"completed"`
	Ready      []WorkstreamStatus `json:"ready,omitempty"`
	InProgress []WorkstreamStatus `json:"in_progress,omitempty"`
	Blocked    []WorkstreamStatus `json:"blocked,omitempty"`
	Failed     []WorkstreamStatus `json:"failed,omitempty"`
}

func BuildStatusView(projectRoot string, state ProjectState, nextStep *Recommendation) *StatusView {
	env := EnvironmentView{
		ProjectRoot:    projectRoot,
		GitBranch:      state.GitStatus.Branch,
		MainBranch:     state.GitStatus.MainBranch,
		Uncommitted:    state.GitStatus.Uncommitted,
		HasSDPConfig:   state.Config.HasSDPConfig,
		HasSDPDir:      dirExists(filepath.Join(projectRoot, ".sdp")),
		HasClaudeDir:   dirExists(filepath.Join(projectRoot, ".claude")),
		HasBeadsIssues: fileExists(filepath.Join(projectRoot, ".beads", "issues.jsonl")),
		HasEvidenceLog: state.Config.EvidenceEnabled,
	}

	workstreams := WorkstreamView{}
	for _, ws := range state.Workstreams {
		workstreams.Total++
		switch ws.Status {
		case StatusReady:
			if len(ws.BlockedBy) > 0 {
				workstreams.Blocked = append(workstreams.Blocked, ws)
				continue
			}
			workstreams.Ready = append(workstreams.Ready, ws)
		case StatusInProgress:
			workstreams.InProgress = append(workstreams.InProgress, ws)
		case StatusBlocked:
			workstreams.Blocked = append(workstreams.Blocked, ws)
		case StatusFailed:
			workstreams.Failed = append(workstreams.Failed, ws)
		case StatusCompleted:
			workstreams.Completed++
		default:
			workstreams.Backlog++
		}
	}
	sortWorkstreams(workstreams.Ready)
	sortWorkstreams(workstreams.InProgress)
	sortWorkstreams(workstreams.Blocked)
	sortWorkstreams(workstreams.Failed)

	view := &StatusView{
		Version:       ContractVersion,
		HasGit:        state.GitStatus.IsRepo,
		HasClaude:     env.HasClaudeDir,
		HasSDP:        env.HasSDPDir,
		HasBeads:      env.HasBeadsIssues,
		Environment:   env,
		Workstreams:   workstreams,
		ActiveSession: state.Session,
		NextStep:      nextStep,
	}
	if nextStep != nil {
		view.NextAction = nextStep.Command
	}
	return view
}

func (s *StatusView) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func sortWorkstreams(items []WorkstreamStatus) {
	slices.SortFunc(items, func(a, b WorkstreamStatus) int {
		return ComparePriority(a, b)
	})
}
