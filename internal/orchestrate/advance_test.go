package orchestrate_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/orchestrate"
)

func TestCurrentBuildWS(t *testing.T) {
	tests := []struct {
		cp   *orchestrate.Checkpoint
		want string
	}{
		{
			cp:   &orchestrate.Checkpoint{Phase: orchestrate.PhaseBuild, Workstreams: []orchestrate.WSStatus{{ID: "00-023-01", Status: "pending"}}},
			want: "00-023-01",
		},
		{
			cp:   &orchestrate.Checkpoint{Phase: orchestrate.PhaseBuild, Workstreams: []orchestrate.WSStatus{{ID: "00-023-01", Status: "done"}, {ID: "00-023-02", Status: "pending"}}},
			want: "00-023-02",
		},
		{
			cp:   &orchestrate.Checkpoint{Phase: orchestrate.PhaseBuild, Workstreams: []orchestrate.WSStatus{{ID: "00-023-01", Status: "done"}, {ID: "00-023-02", Status: "done"}}},
			want: "",
		},
		{
			cp:   &orchestrate.Checkpoint{Phase: orchestrate.PhaseReview},
			want: "",
		},
	}
	for _, tt := range tests {
		got := orchestrate.CurrentBuildWS(tt.cp)
		if got != tt.want {
			t.Errorf("CurrentBuildWS() = %q, want %q", got, tt.want)
		}
	}
}

func TestRunGuardCheck_AdvanceWithCleanScope(t *testing.T) {
	dir := t.TempDir()
	setupGuardTestProject(t, dir)

	// Commit in-scope change
	guardDir := filepath.Join(dir, "internal", "guard")
	if err := os.MkdirAll(guardDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(guardDir, "scope_check.go"), []byte("package guard\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "add", "internal/guard/scope_check.go")
	runGit(t, dir, "commit", "-m", "add scope_check")

	err := orchestrate.RunGuardCheck(dir, "00-023-01")
	if err != nil {
		t.Errorf("expected pass, got: %v", err)
	}
}

func TestRunGuardCheck_AdvanceWithViolationBlocked(t *testing.T) {
	dir := t.TempDir()
	setupGuardTestProject(t, dir)

	// Commit out-of-scope change
	if err := os.MkdirAll(filepath.Join(dir, "cmd", "other"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "cmd", "other", "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "add", "cmd/other/main.go")
	runGit(t, dir, "commit", "-m", "add out of scope")

	err := orchestrate.RunGuardCheck(dir, "00-023-01")
	if err == nil {
		t.Fatal("expected scope violation error")
	}
	var scopeErr *orchestrate.ScopeViolationError
	if !errors.As(err, &scopeErr) {
		t.Errorf("expected ScopeViolationError, got %T", err)
	}
	if scopeErr.WSID != "00-023-01" || len(scopeErr.Violations) == 0 {
		t.Errorf("got WSID=%q violations=%v", scopeErr.WSID, scopeErr.Violations)
	}
}

func setupGuardTestProject(t *testing.T, dir string) {
	t.Helper()
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@test")
	runGit(t, dir, "config", "user.name", "Test")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "init", "--allow-empty")

	wsContent := "---\nws_id: 00-023-01\n---\n\n## Scope Files\n\n- `internal/guard/scope_check.go`\n"
	wsPath := filepath.Join(dir, "docs", "workstreams", "backlog", "00-023-01.md")
	if err := os.MkdirAll(filepath.Dir(wsPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(wsPath, []byte(wsContent), 0o644); err != nil {
		t.Fatal(err)
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_AUTHOR_DATE=2020-01-01T00:00:00Z", "GIT_COMMITTER_DATE=2020-01-01T00:00:00Z")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}
