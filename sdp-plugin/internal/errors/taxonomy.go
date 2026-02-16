// Package errors provides a structured error taxonomy for SDP.
// All critical execution paths use these error types for consistent
// classification and recovery.
package errors

import (
	"fmt"
	"strings"
)

// ErrorClass represents the top-level category of an error.
type ErrorClass string

const (
	// ClassEnvironment indicates issues with the runtime environment
	// (missing tools, permissions, filesystem state).
	ClassEnvironment ErrorClass = "ENV"
	// ClassProtocol indicates violations of SDP protocol rules
	// (invalid workstream IDs, malformed files, missing required fields).
	ClassProtocol ErrorClass = "PROTO"
	// ClassDependency indicates issues with workstream dependencies
	// (blocked workstreams, circular dependencies, missing prerequisites).
	ClassDependency ErrorClass = "DEP"
	// ClassValidation indicates input or state validation failures
	// (coverage below threshold, file too large, type errors).
	ClassValidation ErrorClass = "VAL"
	// ClassRuntime indicates unexpected runtime failures
	// (network errors, external command failures, resource exhaustion).
	ClassRuntime ErrorClass = "RUNTIME"
)

// IsValid returns true if the error class is valid.
func (c ErrorClass) IsValid() bool {
	switch c {
	case ClassEnvironment, ClassProtocol, ClassDependency, ClassValidation, ClassRuntime:
		return true
	default:
		return false
	}
}

// String returns the string representation of the error class.
func (c ErrorClass) String() string {
	return string(c)
}

// Description returns a human-readable description of the error class.
func (c ErrorClass) Description() string {
	switch c {
	case ClassEnvironment:
		return "Environment issue"
	case ClassProtocol:
		return "Protocol violation"
	case ClassDependency:
		return "Dependency problem"
	case ClassValidation:
		return "Validation failure"
	case ClassRuntime:
		return "Runtime error"
	default:
		return "Unknown error class"
	}
}

// ErrorCode represents a structured error code (e.g., ENV001, PROTO002).
type ErrorCode string

const (
	// Environment errors (ENV001-ENV099)
	ErrGitNotFound       ErrorCode = "ENV001"
	ErrGoNotFound        ErrorCode = "ENV002"
	ErrClaudeNotFound    ErrorCode = "ENV003"
	ErrBeadsNotFound     ErrorCode = "ENV004"
	ErrPermissionDenied  ErrorCode = "ENV005"
	ErrWorktreeNotFound  ErrorCode = "ENV006"
	ErrConfigNotFound    ErrorCode = "ENV007"
	ErrDirectoryNotFound ErrorCode = "ENV008"
	ErrFileNotWritable   ErrorCode = "ENV009"

	// Protocol errors (PROTO001-PROTO099)
	ErrInvalidWorkstreamID ErrorCode = "PROTO001"
	ErrInvalidFeatureID    ErrorCode = "PROTO002"
	ErrMalformedYAML       ErrorCode = "PROTO003"
	ErrMissingRequired     ErrorCode = "PROTO004"
	ErrInvalidStatus       ErrorCode = "PROTO005"
	ErrHashChainBroken     ErrorCode = "PROTO006"
	ErrSessionCorrupted    ErrorCode = "PROTO007"
	ErrInvalidEventType    ErrorCode = "PROTO008"
	ErrSchemaViolation     ErrorCode = "PROTO009"

	// Dependency errors (DEP001-DEP099)
	ErrBlockedWorkstream   ErrorCode = "DEP001"
	ErrCircularDependency  ErrorCode = "DEP002"
	ErrMissingPrerequisite ErrorCode = "DEP003"
	ErrFeatureNotFound     ErrorCode = "DEP004"
	ErrWorkstreamNotFound  ErrorCode = "DEP005"
	ErrCollisionDetected   ErrorCode = "DEP006"

	// Validation errors (VAL001-VAL099)
	ErrCoverageLow       ErrorCode = "VAL001"
	ErrFileTooLarge      ErrorCode = "VAL002"
	ErrTestFailed        ErrorCode = "VAL003"
	ErrLintFailed        ErrorCode = "VAL004"
	ErrTypeMismatch      ErrorCode = "VAL005"
	ErrQualityGateFailed ErrorCode = "VAL006"
	ErrDriftDetected     ErrorCode = "VAL007"
	ErrScopeViolation    ErrorCode = "VAL008"

	// Runtime errors (RUNTIME001-RUNTIME099)
	ErrCommandFailed     ErrorCode = "RUNTIME001"
	ErrNetworkError      ErrorCode = "RUNTIME002"
	ErrTimeoutExceeded   ErrorCode = "RUNTIME003"
	ErrResourceExhausted ErrorCode = "RUNTIME004"
	ErrUnexpectedState   ErrorCode = "RUNTIME005"
	ErrInternalError     ErrorCode = "RUNTIME006"
)

// Class returns the error class for this error code.
func (c ErrorCode) Class() ErrorClass {
	switch {
	case strings.HasPrefix(string(c), "ENV"):
		return ClassEnvironment
	case strings.HasPrefix(string(c), "PROTO"):
		return ClassProtocol
	case strings.HasPrefix(string(c), "DEP"):
		return ClassDependency
	case strings.HasPrefix(string(c), "VAL"):
		return ClassValidation
	case strings.HasPrefix(string(c), "RUNTIME"):
		return ClassRuntime
	default:
		return ClassRuntime
	}
}

// IsValid returns true if the error code is recognized.
func (c ErrorCode) IsValid() bool {
	switch c {
	case ErrGitNotFound, ErrGoNotFound, ErrClaudeNotFound, ErrBeadsNotFound,
		ErrPermissionDenied, ErrWorktreeNotFound, ErrConfigNotFound,
		ErrDirectoryNotFound, ErrFileNotWritable,
		ErrInvalidWorkstreamID, ErrInvalidFeatureID, ErrMalformedYAML,
		ErrMissingRequired, ErrInvalidStatus, ErrHashChainBroken,
		ErrSessionCorrupted, ErrInvalidEventType, ErrSchemaViolation,
		ErrBlockedWorkstream, ErrCircularDependency, ErrMissingPrerequisite,
		ErrFeatureNotFound, ErrWorkstreamNotFound, ErrCollisionDetected,
		ErrCoverageLow, ErrFileTooLarge, ErrTestFailed, ErrLintFailed,
		ErrTypeMismatch, ErrQualityGateFailed, ErrDriftDetected, ErrScopeViolation,
		ErrCommandFailed, ErrNetworkError, ErrTimeoutExceeded,
		ErrResourceExhausted, ErrUnexpectedState, ErrInternalError:
		return true
	default:
		return false
	}
}

// errorMessages maps error codes to user-facing messages.
var errorMessages = map[ErrorCode]string{
	ErrGitNotFound:       "Git is not installed or not found in PATH",
	ErrGoNotFound:        "Go is not installed or not found in PATH",
	ErrClaudeNotFound:    "Claude Code CLI is not installed",
	ErrBeadsNotFound:     "Beads CLI is not installed (required for task tracking)",
	ErrPermissionDenied:  "Permission denied",
	ErrWorktreeNotFound:  "Git worktree not found",
	ErrConfigNotFound:    "SDP configuration file not found",
	ErrDirectoryNotFound: "Required directory not found",
	ErrFileNotWritable:   "File is not writable",
	ErrInvalidWorkstreamID: "Invalid workstream ID format (expected PP-FFF-SS)",
	ErrInvalidFeatureID:    "Invalid feature ID format (expected FNNN)",
	ErrMalformedYAML:       "YAML parsing error",
	ErrMissingRequired:     "Required field is missing",
	ErrInvalidStatus:       "Invalid status value",
	ErrHashChainBroken:     "Evidence hash chain is broken",
	ErrSessionCorrupted:    "Session file is corrupted or tampered",
	ErrInvalidEventType:    "Invalid event type",
	ErrSchemaViolation:     "Schema validation failed",
	ErrBlockedWorkstream:   "Workstream is blocked by unresolved dependencies",
	ErrCircularDependency:  "Circular dependency detected",
	ErrMissingPrerequisite: "Required prerequisite is not satisfied",
	ErrFeatureNotFound:     "Feature not found",
	ErrWorkstreamNotFound:  "Workstream not found",
	ErrCollisionDetected:   "File scope collision detected between workstreams",
	ErrCoverageLow:         "Test coverage is below required threshold",
	ErrFileTooLarge:        "File exceeds maximum allowed size",
	ErrTestFailed:          "Tests failed",
	ErrLintFailed:          "Linting failed",
	ErrTypeMismatch:        "Type mismatch",
	ErrQualityGateFailed:   "Quality gate failed",
	ErrDriftDetected:       "Drift detected between code and documentation",
	ErrScopeViolation:      "Edit scope violation",
	ErrCommandFailed:       "External command failed",
	ErrNetworkError:        "Network error",
	ErrTimeoutExceeded:     "Operation timed out",
	ErrResourceExhausted:   "Resource exhausted",
	ErrUnexpectedState:     "Unexpected state encountered",
	ErrInternalError:       "Internal error",
}

// Message returns the default user-facing message for this error code.
func (c ErrorCode) Message() string {
	if msg, ok := errorMessages[c]; ok {
		return msg
	}
	return "Unknown error"
}

// recoveryHints maps error codes to recovery hints.
var recoveryHints = map[ErrorCode]string{
	ErrGitNotFound:       "Install Git from https://git-scm.com",
	ErrGoNotFound:        "Install Go from https://go.dev/dl/",
	ErrClaudeNotFound:    "Install Claude Code CLI from Anthropic",
	ErrBeadsNotFound:     "Install Beads: brew tap beads-dev/tap && brew install beads",
	ErrPermissionDenied:  "Check file permissions or run with appropriate privileges",
	ErrWorktreeNotFound:  "Verify you are in a valid git worktree",
	ErrConfigNotFound:    "Run 'sdp init' to create configuration",
	ErrDirectoryNotFound: "Ensure required directories exist",
	ErrFileNotWritable:   "Check file permissions",
	ErrInvalidWorkstreamID: "Verify the ID format matches PP-FFF-SS or FNNN",
	ErrInvalidFeatureID:    "Verify the ID format matches PP-FFF-SS or FNNN",
	ErrMalformedYAML:       "Check YAML syntax and structure",
	ErrMissingRequired:     "Provide all required fields",
	ErrInvalidStatus:       "Use valid status: pending, in_progress, completed, failed",
	ErrHashChainBroken:     "Run 'sdp log trace --verify' to diagnose",
	ErrSessionCorrupted:    "Run 'sdp session repair' or delete .sdp/session.json",
	ErrInvalidEventType:    "Use valid event types: plan, generation, verification",
	ErrSchemaViolation:     "Verify file matches expected schema",
	ErrBlockedWorkstream:   "Complete blocking workstreams first or use --force",
	ErrCircularDependency:  "Review workstream dependencies for cycles",
	ErrMissingPrerequisite: "Ensure all prerequisites are satisfied",
	ErrFeatureNotFound:     "Verify the ID exists in docs/workstreams/",
	ErrWorkstreamNotFound:  "Verify the ID exists in docs/workstreams/",
	ErrCollisionDetected:   "Review workstream scope files for overlaps",
	ErrCoverageLow:         "Add tests to increase coverage to >= 80%",
	ErrFileTooLarge:        "Split file into smaller modules (< 200 LOC)",
	ErrTestFailed:          "Run tests with verbose output to diagnose failures",
	ErrLintFailed:          "Fix linting errors reported by linter",
	ErrTypeMismatch:        "Verify types match expected signatures",
	ErrQualityGateFailed:   "Review quality gate output for specific failures",
	ErrDriftDetected:       "Run 'sdp drift detect' for details and sync",
	ErrScopeViolation:      "Stay within workstream scope or use 'sdp guard deactivate'",
	ErrCommandFailed:       "Check command output for details",
	ErrNetworkError:        "Check network connectivity and retry",
	ErrTimeoutExceeded:     "Increase timeout or optimize operation",
	ErrResourceExhausted:   "Free up resources and retry",
	ErrUnexpectedState:     "Run 'sdp doctor' to diagnose environment",
	ErrInternalError:       "Report this issue with full error context",
}

// RecoveryHint returns a brief hint for recovering from this error.
func (c ErrorCode) RecoveryHint() string {
	if hint, ok := recoveryHints[c]; ok {
		return hint
	}
	return "No recovery hint available"
}

// SDPError is the primary error type for SDP operations.
type SDPError struct {
	Code    ErrorCode
	Message string
	Cause   error
	Context map[string]string
}

// Error implements the error interface.
func (e *SDPError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause.
func (e *SDPError) Unwrap() error {
	return e.Cause
}

// Class returns the error class.
func (e *SDPError) Class() ErrorClass {
	return e.Code.Class()
}

// RecoveryHint returns the recovery hint for this error.
func (e *SDPError) RecoveryHint() string {
	return e.Code.RecoveryHint()
}

// WithContext adds context information to the error.
func (e *SDPError) WithContext(key, value string) *SDPError {
	if e.Context == nil {
		e.Context = make(map[string]string)
	}
	e.Context[key] = value
	return e
}

// New creates a new SDPError with the given code and optional cause.
func New(code ErrorCode, cause error) *SDPError {
	return &SDPError{
		Code:    code,
		Message: code.Message(),
		Cause:   cause,
	}
}

// Newf creates a new SDPError with a custom message.
func Newf(code ErrorCode, format string, args ...interface{}) *SDPError {
	return &SDPError{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// Wrap wraps an existing error with an SDP error code.
func Wrap(code ErrorCode, cause error, message string) *SDPError {
	return &SDPError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// IsSDPError checks if an error is an SDPError.
func IsSDPError(err error) bool {
	_, ok := err.(*SDPError)
	return ok
}

// GetCode extracts the error code from an error.
func GetCode(err error) ErrorCode {
	if sdpErr, ok := err.(*SDPError); ok {
		return sdpErr.Code
	}
	return ErrInternalError
}

// GetClass extracts the error class from an error.
func GetClass(err error) ErrorClass {
	return GetCode(err).Class()
}
