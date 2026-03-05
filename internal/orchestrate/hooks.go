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

const hookDisallowedChars = ";|&`$<>()\\\n\r"

var allowedHookCommands = map[string]bool{
	"bd":                 true,
	"echo":               true,
	"false":              true,
	"git":                true,
	"go":                 true,
	"make":               true,
	"notify":             true,
	"sdp":                true,
	"sdp-doc-sync":       true,
	"sdp-evidence":       true,
	"sdp-protocol-check": true,
	"slack-notify":       true,
	"trivy":              true,
	"true":               true,
}

// HookConfig is the schema for .sdp/pipeline-hooks.yaml.
type HookConfig struct {
	Hooks []HookEntry `yaml:"hooks"`
}

// HookEntry defines a single hook.
type HookEntry struct {
	Phase   string `yaml:"phase"` // build, review, ci
	When    string `yaml:"when"`  // pre, post
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
	parts, err := parseAndValidateHookCommand(projectRoot, h.Command)
	if err != nil {
		return fmt.Errorf("hook %s-%s: %w", h.Phase, h.When, err)
	}

	timeout := defaultHookTimeout
	if h.Timeout > 0 {
		timeout = time.Duration(h.Timeout) * time.Second
	}
	hookCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(hookCtx, parts[0], parts[1:]...)
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

	err = cmd.Run()
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

func parseAndValidateHookCommand(projectRoot, cmd string) ([]string, error) {
	trimmed := strings.TrimSpace(cmd)
	if trimmed == "" {
		return nil, fmt.Errorf("empty command")
	}
	if strings.ContainsAny(trimmed, hookDisallowedChars) {
		return nil, fmt.Errorf("command contains disallowed shell metacharacters")
	}

	parts, err := parseCommandParts(trimmed)
	if err != nil {
		return nil, err
	}

	if !isAllowedHookExecutable(projectRoot, parts[0]) {
		return nil, fmt.Errorf("command %q is not in allowlist", parts[0])
	}

	return parts, nil
}

func parseCommandParts(cmd string) ([]string, error) {
	parts := make([]string, 0, 4)
	var current strings.Builder
	var quote rune
	escaped := false

	for _, r := range cmd {
		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}

		if r == '\\' {
			escaped = true
			continue
		}

		if quote != 0 {
			if r == quote {
				quote = 0
			} else {
				current.WriteRune(r)
			}
			continue
		}

		switch r {
		case '\'', '"':
			quote = r
		case ' ', '\t':
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if escaped {
		return nil, fmt.Errorf("invalid trailing escape")
	}
	if quote != 0 {
		return nil, fmt.Errorf("unterminated quote")
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid command")
	}

	return parts, nil
}

func isAllowedHookExecutable(projectRoot, command string) bool {
	if allowedHookCommands[command] {
		return true
	}

	if strings.HasPrefix(command, "./") || strings.Contains(command, "/") {
		candidate := command
		if !filepath.IsAbs(candidate) {
			candidate = filepath.Join(projectRoot, candidate)
		}
		clean := filepath.Clean(candidate)
		rel, err := filepath.Rel(projectRoot, clean)
		if err != nil {
			return false
		}
		if strings.HasPrefix(rel, "..") {
			return false
		}
		return true
	}

	return false
}
