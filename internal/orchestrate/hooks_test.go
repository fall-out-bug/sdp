package orchestrate_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/orchestrate"
)

func TestLoadHookConfig_MissingFile(t *testing.T) {
	dir := t.TempDir()
	cfg, err := orchestrate.LoadHookConfig(dir)
	if err != nil {
		t.Fatalf("LoadHookConfig: %v", err)
	}
	if cfg != nil {
		t.Error("expected nil config when file missing")
	}
}

func TestLoadHookConfig_Valid(t *testing.T) {
	dir := t.TempDir()
	sdp := filepath.Join(dir, ".sdp")
	if err := os.MkdirAll(sdp, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(sdp, "pipeline-hooks.yaml")
	content := `
hooks:
  - phase: build
    when: post
    command: "echo post-build"
    on_fail: halt
  - phase: review
    when: pre
    command: "echo pre-review"
    on_fail: warn
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := orchestrate.LoadHookConfig(dir)
	if err != nil {
		t.Fatalf("LoadHookConfig: %v", err)
	}
	if cfg == nil || len(cfg.Hooks) != 2 {
		t.Fatalf("expected 2 hooks, got %v", cfg)
	}
	if cfg.Hooks[0].Phase != "build" || cfg.Hooks[0].When != "post" || cfg.Hooks[0].OnFail != "halt" {
		t.Errorf("hook 0: %+v", cfg.Hooks[0])
	}
	if cfg.Hooks[1].Phase != "review" || cfg.Hooks[1].When != "pre" || cfg.Hooks[1].OnFail != "warn" {
		t.Errorf("hook 1: %+v", cfg.Hooks[1])
	}
}

func TestRunHooks_PreBuildHalt(t *testing.T) {
	dir := t.TempDir()
	sdp := filepath.Join(dir, ".sdp")
	if err := os.MkdirAll(sdp, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(sdp, "pipeline-hooks.yaml")
	content := `
hooks:
  - phase: build
    when: pre
    command: "exit 1"
    on_fail: halt
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	env := orchestrate.HookEnv{WSID: "00-024-01", FeatureID: "F024", Phase: "build"}
	err := orchestrate.RunHooks(ctx, dir, "build", "pre", env, nil)
	if err == nil {
		t.Error("expected error from halt hook")
	}
}

func TestRunHooks_PostBuildWarn(t *testing.T) {
	dir := t.TempDir()
	sdp := filepath.Join(dir, ".sdp")
	if err := os.MkdirAll(sdp, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(sdp, "pipeline-hooks.yaml")
	content := `
hooks:
  - phase: build
    when: post
    command: "exit 1"
    on_fail: warn
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	env := orchestrate.HookEnv{WSID: "00-024-01", FeatureID: "F024", Phase: "build"}
	err := orchestrate.RunHooks(ctx, dir, "build", "post", env, nil)
	if err != nil {
		t.Errorf("warn should not fail: %v", err)
	}
}

func TestRunHooks_Ignore(t *testing.T) {
	dir := t.TempDir()
	sdp := filepath.Join(dir, ".sdp")
	if err := os.MkdirAll(sdp, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(sdp, "pipeline-hooks.yaml")
	content := `
hooks:
  - phase: ci
    when: post
    command: "exit 42"
    on_fail: ignore
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	err := orchestrate.RunHooks(ctx, dir, "ci", "post", orchestrate.HookEnv{}, nil)
	if err != nil {
		t.Errorf("ignore should not fail: %v", err)
	}
}

func TestRunHooks_MissingConfig(t *testing.T) {
	dir := t.TempDir()
	ctx := context.Background()
	err := orchestrate.RunHooks(ctx, dir, "build", "pre", orchestrate.HookEnv{}, nil)
	if err != nil {
		t.Errorf("missing config should not fail: %v", err)
	}
}
