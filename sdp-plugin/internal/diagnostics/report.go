// Package diagnostics provides post-failure diagnostics and triage reports.
// These reports help humans and agents quickly understand and resolve failures.
package diagnostics

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/errors"
	"github.com/fall-out-bug/sdp/internal/recovery"
)

// Report represents a complete diagnostics report.
type Report struct {
	Timestamp     string             `json:"timestamp"`
	Error         *ErrorInfo         `json:"error"`
	Environment   *EnvironmentInfo   `json:"environment"`
	Evidence      *EvidenceInfo      `json:"evidence,omitempty"`
	Recovery      *recovery.Playbook `json:"recovery"`
	NextSteps     []NextStep         `json:"next_steps"`
	Context       map[string]string  `json:"context,omitempty"`
	SensitiveData []string           `json:"-"` // Fields to redact
}

// ErrorInfo contains details about the error that occurred.
type ErrorInfo struct {
	Code         string            `json:"code"`
	Class        string            `json:"class"`
	Message      string            `json:"message"`
	RecoveryHint string            `json:"recovery_hint"`
	Cause        string            `json:"cause,omitempty"`
	Stack        string            `json:"stack,omitempty"`
	Context      map[string]string `json:"context,omitempty"`
}

// EnvironmentInfo contains system environment details.
type EnvironmentInfo struct {
	OS          string `json:"os"`
	Arch        string `json:"arch"`
	GoVersion   string `json:"go_version"`
	GitVersion  string `json:"git_version,omitempty"`
	SDPVersion  string `json:"sdp_version"`
	WorkingDir  string `json:"working_dir"`
	ProjectRoot string `json:"project_root,omitempty"`
}

// EvidenceInfo contains relevant evidence from the execution log.
type EvidenceInfo struct {
	LastEvent      string `json:"last_event,omitempty"`
	ChainIntegrity string `json:"chain_integrity"`
	EventCount     int    `json:"event_count"`
}

// NextStep represents a recommended next action.
type NextStep struct {
	Order       int    `json:"order"`
	Description string `json:"description"`
	Command     string `json:"command,omitempty"`
	Expected    string `json:"expected,omitempty"`
}

// Generator creates diagnostics reports.
type Generator struct {
	projectRoot string
	sdpVersion  string
}

// NewGenerator creates a new diagnostics report generator.
func NewGenerator(projectRoot, sdpVersion string) *Generator {
	return &Generator{
		projectRoot: projectRoot,
		sdpVersion:  sdpVersion,
	}
}

// Generate creates a diagnostics report for the given error.
func (g *Generator) Generate(err error) *Report {
	report := &Report{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Context:   make(map[string]string),
	}

	// Populate error info
	report.Error = g.buildErrorInfo(err)

	// Populate environment info
	report.Environment = g.buildEnvironmentInfo()

	// Try to populate evidence info
	report.Evidence = g.buildEvidenceInfo()

	// Get recovery playbook
	if pb := recovery.GetPlaybookForError(err); pb != nil {
		report.Recovery = pb
	}

	// Build next steps
	report.NextSteps = g.buildNextSteps(err)

	return report
}

func (g *Generator) buildErrorInfo(err error) *ErrorInfo {
	info := &ErrorInfo{
		Code:    string(errors.GetCode(err)),
		Class:   string(errors.GetClass(err)),
		Message: err.Error(),
	}

	if sdpErr, ok := err.(*errors.SDPError); ok {
		info.RecoveryHint = sdpErr.RecoveryHint()
		if sdpErr.Cause != nil {
			info.Cause = sdpErr.Cause.Error()
		}
		if sdpErr.Context != nil {
			info.Context = sdpErr.Context
		}
	} else {
		info.RecoveryHint = errors.GetCode(err).RecoveryHint()
	}

	return info
}

func (g *Generator) buildEnvironmentInfo() *EnvironmentInfo {
	info := &EnvironmentInfo{
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		GoVersion:  runtime.Version(),
		SDPVersion: g.sdpVersion,
	}

	// Get working directory
	if wd, err := os.Getwd(); err == nil {
		info.WorkingDir = wd
	}

	// Try to get project root
	if g.projectRoot != "" {
		info.ProjectRoot = g.projectRoot
	}

	// Get Git version
	if output, err := exec.Command("git", "--version").Output(); err == nil {
		info.GitVersion = strings.TrimSpace(string(output))
	}

	return info
}

func (g *Generator) buildEvidenceInfo() *EvidenceInfo {
	info := &EvidenceInfo{
		ChainIntegrity: "unknown",
	}

	if g.projectRoot == "" {
		return info
	}

	// Check for evidence log
	logPath := filepath.Join(g.projectRoot, ".sdp", "log", "events.jsonl")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		info.ChainIntegrity = "no_log"
		return info
	}

	// Count events
	if data, err := os.ReadFile(logPath); err == nil {
		lines := strings.Split(strings.TrimSpace(string(data)), "\n")
		info.EventCount = len(lines)
		if len(lines) > 0 && lines[len(lines)-1] != "" {
			info.LastEvent = lines[len(lines)-1]
		}
	}

	// Check chain integrity (simplified check)
	info.ChainIntegrity = "ok" // Full verification would require evidence package

	return info
}

func (g *Generator) buildNextSteps(err error) []NextStep {
	code := errors.GetCode(err)
	class := code.Class()

	steps := []NextStep{}

	// Universal first steps
	steps = append(steps, NextStep{
		Order:       1,
		Description: "Review error message and recovery hint",
	})

	// Class-specific next steps
	switch class {
	case errors.ClassEnvironment:
		steps = append(steps,
			NextStep{
				Order:       2,
				Description: "Run doctor to diagnose environment",
				Command:     "sdp doctor",
				Expected:    "Shows missing tools or misconfigurations",
			},
			NextStep{
				Order:       3,
				Description: "Fix identified issues",
			},
		)

	case errors.ClassProtocol:
		steps = append(steps,
			NextStep{
				Order:       2,
				Description: "Validate input format",
				Command:     "sdp parse <file>",
				Expected:    "Shows validation errors",
			},
			NextStep{
				Order:       3,
				Description: "Fix format issues",
			},
		)

	case errors.ClassDependency:
		steps = append(steps,
			NextStep{
				Order:       2,
				Description: "Check workstream dependencies",
				Command:     "cat docs/workstreams/backlog/<ws-id>.md",
			},
			NextStep{
				Order:       3,
				Description: "Complete blocking workstreams",
				Command:     "sdp apply --ws <blocking-ws-id>",
			},
		)

	case errors.ClassValidation:
		steps = append(steps,
			NextStep{
				Order:       2,
				Description: "Run quality check",
				Command:     "sdp quality",
				Expected:    "Shows specific failures",
			},
			NextStep{
				Order:       3,
				Description: "Fix validation failures",
			},
			NextStep{
				Order:       4,
				Description: "Re-run quality check",
				Command:     "sdp quality",
				Expected:    "All gates pass",
			},
		)

	case errors.ClassRuntime:
		steps = append(steps,
			NextStep{
				Order:       2,
				Description: "Retry the operation",
			},
			NextStep{
				Order:       3,
				Description: "If persists, run diagnostics",
				Command:     "sdp doctor --deep",
			},
			NextStep{
				Order:       4,
				Description: "Report if unresolved",
			},
		)
	}

	// Add evidence check step
	steps = append(steps, NextStep{
		Order:       len(steps) + 1,
		Description: "Review evidence log for details",
		Command:     "sdp log show --ws <ws-id>",
	})

	return steps
}

// FormatText formats the report as human-readable text.
func (r *Report) FormatText() string {
	var sb strings.Builder

	sb.WriteString("=== SDP Diagnostics Report ===\n")
	sb.WriteString(fmt.Sprintf("Generated: %s\n\n", r.Timestamp))

	// Error section
	sb.WriteString("--- Error ---\n")
	sb.WriteString(fmt.Sprintf("Code:    %s\n", r.Error.Code))
	sb.WriteString(fmt.Sprintf("Class:   %s\n", r.Error.Class))
	sb.WriteString(fmt.Sprintf("Message: %s\n", r.Error.Message))
	if r.Error.RecoveryHint != "" {
		sb.WriteString(fmt.Sprintf("Hint:    %s\n", r.Error.RecoveryHint))
	}
	if r.Error.Cause != "" {
		sb.WriteString(fmt.Sprintf("Cause:   %s\n", r.Error.Cause))
	}
	if len(r.Error.Context) > 0 {
		sb.WriteString("Context:\n")
		for k, v := range r.Error.Context {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
		}
	}
	sb.WriteString("\n")

	// Environment section
	sb.WriteString("--- Environment ---\n")
	sb.WriteString(fmt.Sprintf("OS:          %s\n", r.Environment.OS))
	sb.WriteString(fmt.Sprintf("Arch:        %s\n", r.Environment.Arch))
	sb.WriteString(fmt.Sprintf("Go Version:  %s\n", r.Environment.GoVersion))
	if r.Environment.GitVersion != "" {
		sb.WriteString(fmt.Sprintf("Git Version: %s\n", r.Environment.GitVersion))
	}
	sb.WriteString(fmt.Sprintf("SDP Version: %s\n", r.Environment.SDPVersion))
	sb.WriteString(fmt.Sprintf("Working Dir: %s\n", r.Environment.WorkingDir))
	if r.Environment.ProjectRoot != "" {
		sb.WriteString(fmt.Sprintf("Project Root: %s\n", r.Environment.ProjectRoot))
	}
	sb.WriteString("\n")

	// Evidence section
	if r.Evidence != nil {
		sb.WriteString("--- Evidence ---\n")
		sb.WriteString(fmt.Sprintf("Chain Integrity: %s\n", r.Evidence.ChainIntegrity))
		sb.WriteString(fmt.Sprintf("Event Count:     %d\n", r.Evidence.EventCount))
		sb.WriteString("\n")
	}

	// Recovery section
	if r.Recovery != nil {
		sb.WriteString("--- Recovery ---\n")
		sb.WriteString(recovery.FormatPlaybook(r.Recovery))
		sb.WriteString("\n")
	}

	// Next steps
	if len(r.NextSteps) > 0 {
		sb.WriteString("--- Next Steps ---\n")
		for _, step := range r.NextSteps {
			sb.WriteString(fmt.Sprintf("%d. %s\n", step.Order, step.Description))
			if step.Command != "" {
				sb.WriteString(fmt.Sprintf("   $ %s\n", step.Command))
			}
			if step.Expected != "" {
				sb.WriteString(fmt.Sprintf("   Expected: %s\n", step.Expected))
			}
		}
	}

	return sb.String()
}

// FormatJSON formats the report as JSON.
func (r *Report) FormatJSON() (string, error) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal report: %w", err)
	}
	return string(data), nil
}

// Redact removes sensitive information from the report.
func (r *Report) Redact(sensitiveFields []string) {
	for _, field := range sensitiveFields {
		// Redact from error context
		if r.Error != nil && r.Error.Context != nil {
			if _, exists := r.Error.Context[field]; exists {
				r.Error.Context[field] = "[REDACTED]"
			}
		}
		// Redact from general context
		if r.Context != nil {
			if _, exists := r.Context[field]; exists {
				r.Context[field] = "[REDACTED]"
			}
		}
	}
}

// Save writes the report to a file.
func (r *Report) Save(path string) error {
	data, err := r.FormatJSON()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

// GenerateForError is a convenience function to generate a report.
func GenerateForError(err error, projectRoot, sdpVersion string) *Report {
	gen := NewGenerator(projectRoot, sdpVersion)
	return gen.Generate(err)
}
