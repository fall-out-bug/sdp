package verify

import "time"

// CheckResult represents a single verification check
type CheckResult struct {
	Name     string `json:"name"`
	Passed   bool   `json:"passed"`
	Message  string `json:"message,omitempty"`
	Evidence string `json:"evidence,omitempty"`
}

// VerificationResult represents the complete verification result
type VerificationResult struct {
	WSID           string        `json:"ws_id"`
	Passed         bool          `json:"passed"`
	Checks         []CheckResult `json:"checks"`
	CoverageActual float64       `json:"coverage_actual,omitempty"`
	MissingFiles   []string      `json:"missing_files,omitempty"`
	FailedCommands []string      `json:"failed_commands,omitempty"`
	Duration       time.Duration `json:"duration"`
}

// WorkstreamData represents parsed workstream frontmatter
type WorkstreamData struct {
	WSID                 string   `json:"ws_id" yaml:"ws_id"`
	Title                string   `json:"title" yaml:"title"`
	Status               string   `json:"status" yaml:"status"`
	ScopeFiles           []string `json:"scope_files" yaml:"scope_files"`
	VerificationCommands []string `json:"verification_commands" yaml:"verification_commands"`
	CoverageThreshold    float64  `json:"coverage_threshold" yaml:"coverage_threshold"`
}
