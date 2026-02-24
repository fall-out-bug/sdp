package guard_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/guard"
)

func TestParseScopeFiles(t *testing.T) {
	content := "---\nws_id: 00-023-01\n---\n\n# WS\n\n## Scope Files\n\n- `internal/guard/scope_check.go` — new\n- `internal/guard/allowlist.go` — new\n- `internal/guard/scope_check_test.go` — test\n\n## Other Section\n\n- `ignored.go`\n"
	paths, err := guard.ParseScopeFilesFromContent(content)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"internal/guard/scope_check.go", "internal/guard/allowlist.go", "internal/guard/scope_check_test.go"}
	if len(paths) != len(want) {
		t.Fatalf("got %d paths, want %d: %v", len(paths), len(want), paths)
	}
	for i, p := range paths {
		if p != want[i] {
			t.Errorf("paths[%d] = %q, want %q", i, p, want[i])
		}
	}
}

func TestIsAllowlisted(t *testing.T) {
	allowlist := []string{"go.sum", "go.mod", "package-lock.json"}
	tests := []struct {
		file string
		want bool
	}{
		{"go.sum", true},
		{"go.mod", true},
		{"internal/foo.go", false},
		{"pkg/bar/go.mod", true},
	}
	for _, tt := range tests {
		got := guard.IsAllowlisted(tt.file, allowlist)
		if got != tt.want {
			t.Errorf("IsAllowlisted(%q) = %v, want %v", tt.file, got, tt.want)
		}
	}
}

func TestCheckScope_InScopeOnly(t *testing.T) {
	dir := t.TempDir()
	setupProject(t, dir)

	wsContent := "---\nws_id: 00-023-01\n---\n\n## Scope Files\n\n- `internal/guard/scope_check.go`\n- `internal/guard/allowlist.go`\n"
	wsPath := filepath.Join(dir, "docs", "workstreams", "backlog", "00-023-01.md")
	if err := os.MkdirAll(filepath.Dir(wsPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(wsPath, []byte(wsContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create in-scope file and commit
	guardDir := filepath.Join(dir, "internal", "guard")
	if err := os.MkdirAll(guardDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(guardDir, "scope_check.go"), []byte("package guard\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "add", "internal/guard/scope_check.go")
	runGit(t, dir, "commit", "-m", "add scope_check")

	verdict, err := guard.CheckScope(dir, "00-023-01", false)
	if err != nil {
		t.Fatal(err)
	}
	if !verdict.Pass {
		t.Errorf("expected pass, got violations: %v", verdict.Violations)
	}
}

func TestCheckScope_OutOfScopeViolation(t *testing.T) {
	dir := t.TempDir()
	setupProject(t, dir)

	wsContent := "---\nws_id: 00-023-01\n---\n\n## Scope Files\n\n- `internal/guard/scope_check.go`\n"
	wsPath := filepath.Join(dir, "docs", "workstreams", "backlog", "00-023-01.md")
	if err := os.MkdirAll(filepath.Dir(wsPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(wsPath, []byte(wsContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create out-of-scope file
	if err := os.MkdirAll(filepath.Join(dir, "cmd", "other"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "cmd", "other", "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "add", "cmd/other/main.go")
	runGit(t, dir, "commit", "-m", "add out of scope")

	verdict, err := guard.CheckScope(dir, "00-023-01", false)
	if err != nil {
		t.Fatal(err)
	}
	if verdict.Pass {
		t.Error("expected fail for out-of-scope change")
	}
	if len(verdict.Violations) != 1 || verdict.Violations[0] != "cmd/other/main.go" {
		t.Errorf("got violations %v", verdict.Violations)
	}
}

func TestCheckScope_Allowlisted(t *testing.T) {
	dir := t.TempDir()
	setupProject(t, dir)

	wsContent := "---\nws_id: 00-023-01\n---\n\n## Scope Files\n\n- `internal/guard/scope_check.go`\n"
	wsPath := filepath.Join(dir, "docs", "workstreams", "backlog", "00-023-01.md")
	if err := os.MkdirAll(filepath.Dir(wsPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(wsPath, []byte(wsContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Change go.mod (allowlisted)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, dir, "add", "go.mod")
	runGit(t, dir, "commit", "-m", "bump deps")

	verdict, err := guard.CheckScope(dir, "00-023-01", false)
	if err != nil {
		t.Fatal(err)
	}
	if !verdict.Pass {
		t.Errorf("expected pass for allowlisted go.mod, got violations: %v", verdict.Violations)
	}
	if len(verdict.Warnings) != 1 || verdict.Warnings[0] != "go.mod" {
		t.Errorf("got warnings %v", verdict.Warnings)
	}
}

func setupProject(t *testing.T, dir string) {
	t.Helper()
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@test")
	runGit(t, dir, "config", "user.name", "Test")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "init", "--allow-empty")
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
