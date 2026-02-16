// Package recovery provides recovery playbooks for SDP failures.
// Each playbook defines actionable steps to recover from specific error types.
package recovery

import (
	"fmt"
	"strings"

	"github.com/fall-out-bug/sdp/internal/errors"
)

// Step represents a single recovery step.
type Step struct {
	Order       int      `json:"order"`
	Description string   `json:"description"`
	Command     string   `json:"command,omitempty"`
	Expected    string   `json:"expected,omitempty"`
	Notes       []string `json:"notes,omitempty"`
}

// Playbook defines a complete recovery procedure for an error class or code.
type Playbook struct {
	Code        errors.ErrorCode `json:"code"`
	Title       string           `json:"title"`
	Severity    string           `json:"severity"` // P0, P1, P2
	FastPath    []Step           `json:"fast_path"`
	DeepPath    []Step           `json:"deep_path,omitempty"`
	RelatedDocs []string         `json:"related_docs,omitempty"`
}

// PlaybookRegistry holds all available playbooks.
type PlaybookRegistry struct {
	playbooks map[errors.ErrorCode]*Playbook
	classDefaults map[errors.ErrorClass]*Playbook
}

// NewPlaybookRegistry creates a new playbook registry with built-in playbooks.
func NewPlaybookRegistry() *PlaybookRegistry {
	r := &PlaybookRegistry{
		playbooks: make(map[errors.ErrorCode]*Playbook),
		classDefaults: make(map[errors.ErrorClass]*Playbook),
	}
	r.initBuiltins()
	return r
}

// Get retrieves a playbook for a specific error code.
func (r *PlaybookRegistry) Get(code errors.ErrorCode) *Playbook {
	if pb, ok := r.playbooks[code]; ok {
		return pb
	}
	// Fall back to class default
	if pb, ok := r.classDefaults[code.Class()]; ok {
		return pb
	}
	return nil
}

// GetForError retrieves a playbook for any error type.
func (r *PlaybookRegistry) GetForError(err error) *Playbook {
	return r.Get(errors.GetCode(err))
}

// Register adds a new playbook to the registry.
func (r *PlaybookRegistry) Register(pb *Playbook) {
	r.playbooks[pb.Code] = pb
}

// RegisterClassDefault sets the default playbook for an error class.
func (r *PlaybookRegistry) RegisterClassDefault(class errors.ErrorClass, pb *Playbook) {
	r.classDefaults[class] = pb
}

// List returns all registered playbooks.
func (r *PlaybookRegistry) List() []*Playbook {
	result := make([]*Playbook, 0, len(r.playbooks))
	for _, pb := range r.playbooks {
		result = append(result, pb)
	}
	return result
}

func (r *PlaybookRegistry) initBuiltins() {
	// Environment errors
	r.Register(&Playbook{
		Code:     errors.ErrGitNotFound,
		Title:    "Install Git",
		Severity: "P0",
		FastPath: []Step{
			{Order: 1, Description: "Install Git", Command: "brew install git", Expected: "git --version succeeds"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "Download Git from official source", Command: "https://git-scm.com/downloads"},
			{Order: 2, Description: "Verify installation", Command: "git --version", Expected: "git version X.Y.Z"},
			{Order: 3, Description: "Configure Git identity", Command: "git config --global user.name \"Your Name\""},
		},
		RelatedDocs: []string{"https://git-scm.com/book/en/v2/Getting-Started-Installing-Git"},
	})

	r.Register(&Playbook{
		Code:     errors.ErrBeadsNotFound,
		Title:    "Install Beads CLI",
		Severity: "P1",
		FastPath: []Step{
			{Order: 1, Description: "Install via Homebrew", Command: "brew tap beads-dev/tap && brew install beads", Expected: "bd --version succeeds"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "Alternative: Install via curl", Command: "curl -sSL https://raw.githubusercontent.com/beads-dev/beads/main/install.sh | bash"},
			{Order: 2, Description: "Verify installation", Command: "bd --version", Expected: "beads version X.Y.Z"},
		},
	})

	r.Register(&Playbook{
		Code:     errors.ErrPermissionDenied,
		Title:    "Fix File Permissions",
		Severity: "P1",
		FastPath: []Step{
			{Order: 1, Description: "Check current permissions", Command: "ls -la <file>"},
			{Order: 2, Description: "Fix permissions", Command: "chmod 644 <file>", Expected: "File is readable"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "Check file ownership", Command: "ls -la <file>"},
			{Order: 2, Description: "Take ownership if needed", Command: "chown $(whoami) <file>"},
			{Order: 3, Description: "Set correct permissions", Command: "chmod 644 <file>"},
			{Order: 4, Description: "For directories", Command: "chmod 755 <directory>"},
		},
	})

	// Protocol errors
	r.Register(&Playbook{
		Code:     errors.ErrInvalidWorkstreamID,
		Title:    "Fix Workstream ID Format",
		Severity: "P1",
		FastPath: []Step{
			{Order: 1, Description: "Verify ID format is PP-FFF-SS (e.g., 00-070-01)"},
			{Order: 2, Description: "Check for typos or missing dashes"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "PP = Project (00 for SDP)"},
			{Order: 2, Description: "FFF = Feature number (001-999)"},
			{Order: 3, Description: "SS = Step number (01-99)"},
			{Order: 4, Description: "Example: 00-070-01 = SDP project, feature 70, step 1"},
		},
		RelatedDocs: []string{"docs/PROTOCOL.md#workstream-ids"},
	})

	r.Register(&Playbook{
		Code:     errors.ErrHashChainBroken,
		Title:    "Repair Evidence Chain",
		Severity: "P1",
		FastPath: []Step{
			{Order: 1, Description: "Verify chain integrity", Command: "sdp log trace --verify", Expected: "Shows broken links"},
			{Order: 2, Description: "If corrupted, backup and remove", Command: "cp .sdp/log/events.jsonl .sdp/log/events.jsonl.bak"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "Trace the chain", Command: "sdp log trace"},
			{Order: 2, Description: "Identify the break point"},
			{Order: 3, Description: "Backup existing log", Command: "cp .sdp/log/events.jsonl events-backup.jsonl"},
			{Order: 4, Description: "Remove corrupted entries", Command: "head -n <line> events-backup.jsonl > .sdp/log/events.jsonl"},
			{Order: 5, Description: "Verify repair", Command: "sdp log trace --verify", Expected: "Chain is intact"},
		},
		RelatedDocs: []string{"docs/compliance/COMPLIANCE.md#evidence"},
	})

	r.Register(&Playbook{
		Code:     errors.ErrSessionCorrupted,
		Title:    "Repair Session State",
		Severity: "P1",
		FastPath: []Step{
			{Order: 1, Description: "Delete corrupted session", Command: "rm .sdp/session.json"},
			{Order: 2, Description: "Reinitialize", Command: "sdp init", Expected: "New session created"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "Backup session file", Command: "cp .sdp/session.json session-backup.json"},
			{Order: 2, Description: "Try repair command", Command: "sdp session repair"},
			{Order: 3, Description: "If repair fails, delete and reinit", Command: "rm .sdp/session.json && sdp init"},
			{Order: 4, Description: "Verify state", Command: "sdp status", Expected: "Shows valid state"},
		},
	})

	// Dependency errors
	r.Register(&Playbook{
		Code:     errors.ErrBlockedWorkstream,
		Title:    "Resolve Blocked Workstream",
		Severity: "P1",
		FastPath: []Step{
			{Order: 1, Description: "Check dependencies", Command: "cat docs/workstreams/backlog/<ws-id>.md | grep depends_on"},
			{Order: 2, Description: "Complete blocking workstreams first", Command: "sdp apply --ws <blocking-ws-id>"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "List all pending workstreams", Command: "bd list --status pending"},
			{Order: 2, Description: "Build dependency graph", Command: "sdp plan --graph"},
			{Order: 3, Description: "Execute in dependency order"},
			{Order: 4, Description: "Use --force to skip (not recommended)", Command: "sdp apply --ws <ws-id> --force"},
		},
		RelatedDocs: []string{"docs/PROTOCOL.md#dependencies"},
	})

	r.Register(&Playbook{
		Code:     errors.ErrCircularDependency,
		Title:    "Break Circular Dependency",
		Severity: "P0",
		FastPath: []Step{
			{Order: 1, Description: "Review workstream files for circular refs"},
			{Order: 2, Description: "Break cycle by removing one dependency"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "Analyze dependency graph", Command: "sdp plan --graph"},
			{Order: 2, Description: "Identify cycle nodes"},
			{Order: 3, Description: "Refactor workstreams to remove cycle"},
			{Order: 4, Description: "Split workstream if needed"},
			{Order: 5, Description: "Validate fix", Command: "sdp parse docs/workstreams/backlog/*.md"},
		},
	})

	// Validation errors
	r.Register(&Playbook{
		Code:     errors.ErrCoverageLow,
		Title:    "Increase Test Coverage",
		Severity: "P1",
		FastPath: []Step{
			{Order: 1, Description: "Run coverage report", Command: "go test -coverprofile=coverage.out ./..."},
			{Order: 2, Description: "View uncovered code", Command: "go tool cover -func=coverage.out | grep -v 100.0%"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "Generate detailed report", Command: "go tool cover -html=coverage.out -o coverage.html"},
			{Order: 2, Description: "Identify uncovered functions"},
			{Order: 3, Description: "Write tests for uncovered code"},
			{Order: 4, Description: "Target 80%+ coverage", Command: "go test -cover ./..."},
			{Order: 5, Description: "Verify threshold met", Expected: "Coverage >= 80%"},
		},
	})

	r.Register(&Playbook{
		Code:     errors.ErrFileTooLarge,
		Title:    "Split Large File",
		Severity: "P2",
		FastPath: []Step{
			{Order: 1, Description: "Count lines", Command: "wc -l <file>"},
			{Order: 2, Description: "Identify logical splits (by responsibility)"},
			{Order: 3, Description: "Extract to new files", Expected: "Each file < 200 LOC"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "Analyze file structure"},
			{Order: 2, Description: "Identify single responsibilities"},
			{Order: 3, Description: "Create new files for each responsibility"},
			{Order: 4, Description: "Update imports"},
			{Order: 5, Description: "Run tests to verify", Command: "go test ./..."},
			{Order: 6, Description: "Check file sizes", Command: "find . -name '*.go' -exec wc -l {} \\;"},
		},
		RelatedDocs: []string{"docs/reference/PRINCIPLES.md#file-size"},
	})

	r.Register(&Playbook{
		Code:     errors.ErrTestFailed,
		Title:    "Fix Failing Tests",
		Severity: "P0",
		FastPath: []Step{
			{Order: 1, Description: "Run tests with verbose output", Command: "go test -v ./..."},
			{Order: 2, Description: "Identify failing test"},
			{Order: 3, Description: "Fix the issue"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "Run specific failing test", Command: "go test -v -run TestName ./..."},
			{Order: 2, Description: "Analyze error messages"},
			{Order: 3, Description: "Debug with print statements or debugger"},
			{Order: 4, Description: "Fix underlying issue"},
			{Order: 5, Description: "Run all tests", Command: "go test ./...", Expected: "All tests pass"},
		},
	})

	r.Register(&Playbook{
		Code:     errors.ErrDriftDetected,
		Title:    "Resolve Code-Documentation Drift",
		Severity: "P2",
		FastPath: []Step{
			{Order: 1, Description: "Run drift detection", Command: "sdp drift detect", Expected: "Shows drift details"},
			{Order: 2, Description: "Update code or docs to match"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "Get detailed drift report", Command: "sdp drift report --output=drift.md"},
			{Order: 2, Description: "Review each drift item"},
			{Order: 3, Description: "Decide: update code or update docs"},
			{Order: 4, Description: "Make changes"},
			{Order: 5, Description: "Verify drift resolved", Command: "sdp drift detect", Expected: "No drift detected"},
		},
		RelatedDocs: []string{"docs/reference/commands.md#drift"},
	})

	// Runtime errors
	r.Register(&Playbook{
		Code:     errors.ErrCommandFailed,
		Title:    "Debug Command Failure",
		Severity: "P1",
		FastPath: []Step{
			{Order: 1, Description: "Check command output for errors"},
			{Order: 2, Description: "Verify command exists", Command: "which <command>"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "Run command manually with verbose flags"},
			{Order: 2, Description: "Check environment variables", Command: "env | grep <relevant>"},
			{Order: 3, Description: "Check PATH includes command location"},
			{Order: 4, Description: "Try with absolute path"},
			{Order: 5, Description: "Check for permission issues"},
		},
	})

	r.Register(&Playbook{
		Code:     errors.ErrTimeoutExceeded,
		Title:    "Handle Timeout",
		Severity: "P2",
		FastPath: []Step{
			{Order: 1, Description: "Retry the operation"},
			{Order: 2, Description: "Check system load", Command: "top"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "Identify slow operation"},
			{Order: 2, Description: "Profile if available", Command: "go test -cpuprofile=cpu.out -bench ."},
			{Order: 3, Description: "Optimize or increase timeout"},
			{Order: 4, Description: "Consider async execution"},
		},
	})

	r.Register(&Playbook{
		Code:     errors.ErrInternalError,
		Title:    "Report Internal Error",
		Severity: "P0",
		FastPath: []Step{
			{Order: 1, Description: "Run diagnostics", Command: "sdp doctor"},
			{Order: 2, Description: "Save error context"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "Capture full error output"},
			{Order: 2, Description: "Gather environment info", Command: "sdp doctor --deep"},
			{Order: 3, Description: "Collect relevant logs", Command: "sdp log show > error-log.txt"},
			{Order: 4, Description: "Report issue with context", Command: "GitHub issue with: error message, sdp doctor output, steps to reproduce"},
		},
		RelatedDocs: []string{"https://github.com/fall-out-bug/sdp/issues"},
	})

	// Class defaults
	r.RegisterClassDefault(errors.ClassEnvironment, &Playbook{
		Code:     "ENV000",
		Title:    "Fix Environment Issue",
		Severity: "P1",
		FastPath: []Step{
			{Order: 1, Description: "Run diagnostics", Command: "sdp doctor", Expected: "Identifies missing tools"},
			{Order: 2, Description: "Fix identified issues"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "Full diagnostic", Command: "sdp doctor --deep"},
			{Order: 2, Description: "Check PATH and environment"},
			{Order: 3, Description: "Install missing dependencies"},
			{Order: 4, Description: "Fix permissions"},
			{Order: 5, Description: "Verify fixes", Command: "sdp doctor", Expected: "All checks pass"},
		},
	})

	r.RegisterClassDefault(errors.ClassProtocol, &Playbook{
		Code:     "PROTO000",
		Title:    "Fix Protocol Error",
		Severity: "P1",
		FastPath: []Step{
			{Order: 1, Description: "Check input format"},
			{Order: 2, Description: "Validate YAML syntax", Command: "sdp parse <file>"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "Review PROTOCOL.md for format requirements"},
			{Order: 2, Description: "Validate all required fields"},
			{Order: 3, Description: "Check for typos"},
			{Order: 4, Description: "Use schema validation"},
		},
		RelatedDocs: []string{"docs/PROTOCOL.md"},
	})

	r.RegisterClassDefault(errors.ClassDependency, &Playbook{
		Code:     "DEP000",
		Title:    "Resolve Dependency Issue",
		Severity: "P1",
		FastPath: []Step{
			{Order: 1, Description: "Check workstream dependencies"},
			{Order: 2, Description: "Complete prerequisites first"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "Build dependency graph", Command: "sdp plan --graph"},
			{Order: 2, Description: "Identify blockers"},
			{Order: 3, Description: "Execute in correct order"},
		},
	})

	r.RegisterClassDefault(errors.ClassValidation, &Playbook{
		Code:     "VAL000",
		Title:    "Fix Validation Failure",
		Severity: "P1",
		FastPath: []Step{
			{Order: 1, Description: "Run quality check", Command: "sdp quality", Expected: "Shows specific failures"},
			{Order: 2, Description: "Fix identified issues"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "Get detailed report", Command: "sdp quality --verbose"},
			{Order: 2, Description: "Address each failure type"},
			{Order: 3, Description: "Re-run validation", Expected: "All gates pass"},
		},
		RelatedDocs: []string{"docs/reference/PRINCIPLES.md#quality-gates"},
	})

	r.RegisterClassDefault(errors.ClassRuntime, &Playbook{
		Code:     "RUNTIME000",
		Title:    "Debug Runtime Error",
		Severity: "P1",
		FastPath: []Step{
			{Order: 1, Description: "Retry the operation"},
			{Order: 2, Description: "Check system state", Command: "sdp doctor"},
		},
		DeepPath: []Step{
			{Order: 1, Description: "Gather diagnostics", Command: "sdp doctor --deep"},
			{Order: 2, Description: "Check logs", Command: "sdp log show"},
			{Order: 3, Description: "Isolate the failure"},
			{Order: 4, Description: "Report if persistent"},
		},
	})
}

// FormatPlaybook renders a playbook as human-readable text.
func FormatPlaybook(pb *Playbook) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("== %s [%s] ==\n", pb.Title, pb.Severity))
	sb.WriteString("\n")

	if len(pb.FastPath) > 0 {
		sb.WriteString("Quick Fix:\n")
		for _, step := range pb.FastPath {
			sb.WriteString(fmt.Sprintf("  %d. %s\n", step.Order, step.Description))
			if step.Command != "" {
				sb.WriteString(fmt.Sprintf("     $ %s\n", step.Command))
			}
			if step.Expected != "" {
				sb.WriteString(fmt.Sprintf("     Expected: %s\n", step.Expected))
			}
		}
		sb.WriteString("\n")
	}

	if len(pb.DeepPath) > 0 {
		sb.WriteString("Full Recovery:\n")
		for _, step := range pb.DeepPath {
			sb.WriteString(fmt.Sprintf("  %d. %s\n", step.Order, step.Description))
			if step.Command != "" {
				sb.WriteString(fmt.Sprintf("     $ %s\n", step.Command))
			}
			if step.Expected != "" {
				sb.WriteString(fmt.Sprintf("     Expected: %s\n", step.Expected))
			}
		}
		sb.WriteString("\n")
	}

	if len(pb.RelatedDocs) > 0 {
		sb.WriteString("Related Docs:\n")
		for _, doc := range pb.RelatedDocs {
			sb.WriteString(fmt.Sprintf("  - %s\n", doc))
		}
	}

	return sb.String()
}

// Global registry instance
var globalRegistry = NewPlaybookRegistry()

// GetPlaybook retrieves a playbook from the global registry.
func GetPlaybook(code errors.ErrorCode) *Playbook {
	return globalRegistry.Get(code)
}

// GetPlaybookForError retrieves a playbook for any error.
func GetPlaybookForError(err error) *Playbook {
	return globalRegistry.GetForError(err)
}
