package ciloop

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// DefFixer describes a deterministic fixer: command + regex to match failure log.
type DefFixer struct {
	Name      string
	Command   string
	AppliesTo string
	Timeout   int // seconds
}

// builtinFixers are the default deterministic fixers (goimports, go mod tidy, go fmt).
var builtinFixers = []DefFixer{
	{
		Name:      "goimports",
		Command:   "goimports -w .",
		AppliesTo: `could not import|imported and not used|undefined:`,
		Timeout:   30,
	},
	{
		Name:      "go-mod-tidy",
		Command:   "go mod tidy",
		AppliesTo: `missing go\.sum entry|go\.mod file not found|cannot find package`,
		Timeout:   30,
	},
	{
		Name:      "go-fmt",
		Command:   "go fmt ./...",
		AppliesTo: `gofmt|formatting`,
		Timeout:   30,
	},
}

// AutofixerRegistry holds built-in and config-loaded fixers.
type AutofixerRegistry struct {
	Fixers []DefFixer
}

// NewAutofixerRegistry returns a registry with built-ins; optionally loads .sdp/auto-fixers.yaml.
func NewAutofixerRegistry(projectRoot string) *AutofixerRegistry {
	r := &AutofixerRegistry{Fixers: append([]DefFixer{}, builtinFixers...)}
	cfgPath := filepath.Join(projectRoot, ".sdp", "auto-fixers.yaml")
	if data, err := os.ReadFile(cfgPath); err == nil {
		extra, err := ParseAutoFixersYAML(data)
		if err == nil {
			r.Fixers = append(r.Fixers, extra...)
		}
	}
	return r
}

type autoFixersYAML struct {
	Fixers []struct {
		Name      string `yaml:"name"`
		Command   string `yaml:"command"`
		AppliesTo string `yaml:"applies_to"`
		Timeout   int    `yaml:"timeout"`
	} `yaml:"fixers"`
}

// ParseAutoFixersYAML parses .sdp/auto-fixers.yaml format. Exported for testing.
func ParseAutoFixersYAML(data []byte) ([]DefFixer, error) {
	var cfg autoFixersYAML
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	out := make([]DefFixer, 0, len(cfg.Fixers))
	for _, f := range cfg.Fixers {
		if f.Name != "" && f.Command != "" && f.AppliesTo != "" {
			t := f.Timeout
			if t <= 0 {
				t = 30
			}
			out = append(out, DefFixer{Name: f.Name, Command: f.Command, AppliesTo: f.AppliesTo, Timeout: t})
		}
	}
	return out, nil
}

// MatchingFixers returns fixers whose AppliesTo regex matches the failure log.
func (r *AutofixerRegistry) MatchingFixers(failureLog string) []DefFixer {
	var out []DefFixer
	for _, f := range r.Fixers {
		re, err := regexp.Compile(f.AppliesTo)
		if err != nil {
			continue
		}
		if re.MatchString(failureLog) {
			out = append(out, f)
		}
	}
	return out
}

// RunDeterministicFixersOpts configures RunDeterministicFixers.
type RunDeterministicFixersOpts struct {
	Ctx             context.Context
	ProjectRoot     string
	FailureLog     string
	Registry       *AutofixerRegistry
	Committer      Committer
	DecisionLogger func(decision, rationale string) error
	RunFileLogger  func(fixerNames []string, duration time.Duration)
}

// RunDeterministicFixers runs matching fixers in order. If any produces changes,
// commits and pushes, returns true. Otherwise returns false (fall through to LLM).
// Uses exec directly for fixer commands (need Dir, Stdout, Stderr).
func RunDeterministicFixers(ctx context.Context, projectRoot string, failureLog string, registry *AutofixerRegistry, committer Committer, decisionLogger func(decision, rationale string) error, runFileLogger func(fixerNames []string, duration time.Duration)) (changed bool, err error) {
	return runDeterministicFixers(RunDeterministicFixersOpts{
		Ctx: ctx, ProjectRoot: projectRoot, FailureLog: failureLog,
		Registry: registry, Committer: committer,
		DecisionLogger: decisionLogger, RunFileLogger: runFileLogger,
	})
}

func runDeterministicFixers(opts RunDeterministicFixersOpts) (changed bool, err error) {
	matching := opts.Registry.MatchingFixers(opts.FailureLog)
	if len(matching) == 0 {
		return false, nil
	}

	start := time.Now()
	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	for _, f := range matching {
		timeout := time.Duration(f.Timeout) * time.Second
		if timeout <= 0 {
			timeout = 30 * time.Second
		}
		runCtx, cancel := context.WithTimeout(ctx, timeout)
		parts := SplitCommand(f.Command)
		if len(parts) == 0 {
			cancel()
			continue
		}
		cmd := exec.CommandContext(runCtx, parts[0], parts[1:]...)
		cmd.Dir = opts.ProjectRoot
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if runErr := cmd.Run(); runErr != nil {
			cancel()
			continue // fixer failed, try next
		}
		cancel()
	}

	// Check if anything changed
	diffCmd := exec.CommandContext(ctx, "git", "diff", "--quiet")
	diffCmd.Dir = opts.ProjectRoot
	if diffErr := diffCmd.Run(); diffErr == nil {
		return false, nil // no changes
	}

	// Changes produced: commit and push
	names := make([]string, len(matching))
	for i, f := range matching {
		names[i] = f.Name
	}
	msg := fmt.Sprintf("fix(ci): auto-fix %s [deterministic]", strings.Join(names, ", "))
	if err := opts.Committer.Commit(ctx, msg); err != nil {
		return false, fmt.Errorf("commit after deterministic fix: %w", err)
	}
	if err := opts.Committer.Push(ctx); err != nil {
		return false, fmt.Errorf("push after deterministic fix: %w", err)
	}
	if opts.DecisionLogger != nil {
		_ = opts.DecisionLogger("AUTO-FIX", fmt.Sprintf("Deterministic fixers applied: %s", strings.Join(names, ", ")))
	}
	if opts.RunFileLogger != nil {
		opts.RunFileLogger(names, time.Since(start))
	}
	return true, nil
}

// SplitCommand splits a command string into executable and args (handles quoted args).
func SplitCommand(s string) []string {
	var parts []string
	var cur strings.Builder
	inQuote := false
	for _, r := range s {
		switch {
		case r == '"' || r == '\'':
			inQuote = !inQuote
		case (r == ' ' || r == '\t') && !inQuote:
			if cur.Len() > 0 {
				parts = append(parts, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteRune(r)
		}
	}
	if cur.Len() > 0 {
		parts = append(parts, cur.String())
	}
	return parts
}
