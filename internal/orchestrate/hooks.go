package orchestrate

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const defaultHookTimeout = 60 * time.Second

// HookConfig is the schema for .sdp/pipeline-hooks.yaml.
type HookConfig struct {
	Hooks []HookEntry `yaml:"hooks"`
}

// HookEntry defines a single hook.
type HookEntry struct {
	Phase   string `yaml:"phase"`   // build, review, ci
	When    string `yaml:"when"`    // pre, post
	Command string `yaml:"command"`
	OnFail  string `yaml:"on_fail"` // halt, warn, ignore
	Timeout int    `yaml:"timeout"` // seconds; 0 = default 60
}

// LoadHookConfig reads .sdp/pipeline-hooks.yaml. Returns nil if file is missing (graceful degradation).
func LoadHookConfig(projectRoot string) (*HookConfig, error) {
	path := filepath.Join(projectRoot, ".sdp", "pipeline-hooks.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read pipeline-hooks: %w", err)
	}
	var cfg HookConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse pipeline-hooks: %w", err)
	}
	return &cfg, nil
}

// HookEnv holds environment variables for hook execution.
type HookEnv struct {
	WSID           string
	FeatureID      string
	Phase          string
	CheckpointPath string
}

// RunHooks executes hooks matching phase+when. On halt failure, returns error.
// Stdout/stderr are captured and can be logged by the caller.
func RunHooks(ctx context.Context, projectRoot string, phase, when string, env HookEnv, log func(msg string)) error {
	cfg, err := LoadHookConfig(projectRoot)
	if err != nil {
		return err
	}
	if cfg == nil {
		return nil
	}
	for _, h := range cfg.Hooks {
		if h.Phase != phase || h.When != when {
			continue
		}
		if err := runHook(ctx, projectRoot, h, env, log); err != nil {
			return err
		}
	}
	return nil
}

func runHook(ctx context.Context, projectRoot string, h HookEntry, env HookEnv, log func(string)) error {
	timeout := defaultHookTimeout
	if h.Timeout > 0 {
		timeout = time.Duration(h.Timeout) * time.Second
	}
	hookCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(hookCtx, "sh", "-c", h.Command)
	cmd.Dir = projectRoot
	cmd.Env = append(os.Environ(),
		"WS_ID="+env.WSID,
		"FEATURE_ID="+env.FeatureID,
		"PHASE="+env.Phase,
		"CHECKPOINT_PATH="+env.CheckpointPath,
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	out := strings.TrimSpace(stdout.String() + "\n" + stderr.String())
	if out != "" && log != nil {
		log(fmt.Sprintf("hook %s-%s: %s", h.Phase, h.When, out))
	}
	if err == nil {
		return nil
	}
	switch strings.ToLower(h.OnFail) {
	case "ignore":
		return nil
	case "warn":
		if log != nil {
			log(fmt.Sprintf("hook %s-%s failed (warn): %v", h.Phase, h.When, err))
		}
		return nil
	case "halt", "":
		return fmt.Errorf("hook %s-%s failed: %w", h.Phase, h.When, err)
	default:
		return fmt.Errorf("hook %s-%s failed: %w", h.Phase, h.When, err)
	}
}
