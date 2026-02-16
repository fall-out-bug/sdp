package errors

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestErrorClass_IsValid(t *testing.T) {
	tests := []struct {
		class    ErrorClass
		expected bool
	}{
		{ClassEnvironment, true},
		{ClassProtocol, true},
		{ClassDependency, true},
		{ClassValidation, true},
		{ClassRuntime, true},
		{ErrorClass("INVALID"), false},
		{ErrorClass(""), false},
		{ErrorClass("env"), false}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(string(tt.class), func(t *testing.T) {
			if got := tt.class.IsValid(); got != tt.expected {
				t.Errorf("ErrorClass(%q).IsValid() = %v, want %v", tt.class, got, tt.expected)
			}
		})
	}
}

func TestErrorClass_String(t *testing.T) {
	tests := []struct {
		class    ErrorClass
		expected string
	}{
		{ClassEnvironment, "ENV"},
		{ClassProtocol, "PROTO"},
		{ClassDependency, "DEP"},
		{ClassValidation, "VAL"},
		{ClassRuntime, "RUNTIME"},
	}

	for _, tt := range tests {
		t.Run(string(tt.class), func(t *testing.T) {
			if got := tt.class.String(); got != tt.expected {
				t.Errorf("ErrorClass(%q).String() = %q, want %q",
					tt.class, got, tt.expected)
			}
		})
	}
}

func TestErrorClass_Description(t *testing.T) {
	tests := []struct {
		class       ErrorClass
		containsStr string
	}{
		{ClassEnvironment, "Environment"},
		{ClassProtocol, "Protocol"},
		{ClassDependency, "Dependency"},
		{ClassValidation, "Validation"},
		{ClassRuntime, "Runtime"},
		{ErrorClass("UNKNOWN"), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(string(tt.class), func(t *testing.T) {
			got := tt.class.Description()
			if !strings.Contains(got, tt.containsStr) {
				t.Errorf("ErrorClass(%q).Description() = %q, should contain %q",
					tt.class, got, tt.containsStr)
			}
		})
	}
}

func TestErrorCode_Class(t *testing.T) {
	tests := []struct {
		code          ErrorCode
		expectedClass ErrorClass
	}{
		{ErrGitNotFound, ClassEnvironment},
		{ErrGoNotFound, ClassEnvironment},
		{ErrInvalidWorkstreamID, ClassProtocol},
		{ErrMalformedYAML, ClassProtocol},
		{ErrBlockedWorkstream, ClassDependency},
		{ErrCircularDependency, ClassDependency},
		{ErrCoverageLow, ClassValidation},
		{ErrTestFailed, ClassValidation},
		{ErrCommandFailed, ClassRuntime},
		{ErrTimeoutExceeded, ClassRuntime},
		{ErrorCode("UNKNOWN001"), ClassRuntime}, // default
	}

	for _, tt := range tests {
		t.Run(string(tt.code), func(t *testing.T) {
			if got := tt.code.Class(); got != tt.expectedClass {
				t.Errorf("ErrorCode(%q).Class() = %v, want %v", tt.code, got, tt.expectedClass)
			}
		})
	}
}

func TestErrorCode_IsValid(t *testing.T) {
	// Test all defined error codes are valid
	validCodes := []ErrorCode{
		ErrGitNotFound, ErrGoNotFound, ErrClaudeNotFound, ErrBeadsNotFound,
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
		ErrResourceExhausted, ErrUnexpectedState, ErrInternalError,
	}

	for _, code := range validCodes {
		t.Run(string(code), func(t *testing.T) {
			if !code.IsValid() {
				t.Errorf("ErrorCode(%q).IsValid() = false, should be true", code)
			}
		})
	}

	// Test invalid codes
	invalidCodes := []ErrorCode{
		"ENV100",   // undefined
		"PROTO100", // undefined
		"INVALID",  // wrong format
		"",         // empty
	}

	for _, code := range invalidCodes {
		t.Run("invalid_"+string(code), func(t *testing.T) {
			if code.IsValid() {
				t.Errorf("ErrorCode(%q).IsValid() = true, should be false", code)
			}
		})
	}
}

func TestErrorCode_Message(t *testing.T) {
	// Verify all valid codes have non-empty messages
	validCodes := []ErrorCode{
		ErrGitNotFound, ErrGoNotFound, ErrClaudeNotFound, ErrBeadsNotFound,
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
		ErrResourceExhausted, ErrUnexpectedState, ErrInternalError,
	}

	for _, code := range validCodes {
		t.Run(string(code), func(t *testing.T) {
			msg := code.Message()
			if msg == "" || msg == "Unknown error" {
				t.Errorf("ErrorCode(%q).Message() = %q, should have specific message", code, msg)
			}
		})
	}

	// Unknown code should return generic message
	t.Run("unknown_code", func(t *testing.T) {
		msg := ErrorCode("UNKNOWN").Message()
		if msg != "Unknown error" {
			t.Errorf("Unknown code Message() = %q, want 'Unknown error'", msg)
		}
	})
}

func TestErrorCode_RecoveryHint(t *testing.T) {
	// Verify all valid codes have non-empty recovery hints
	validCodes := []ErrorCode{
		ErrGitNotFound, ErrGoNotFound, ErrClaudeNotFound, ErrBeadsNotFound,
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
		ErrResourceExhausted, ErrUnexpectedState, ErrInternalError,
	}

	for _, code := range validCodes {
		t.Run(string(code), func(t *testing.T) {
			hint := code.RecoveryHint()
			if hint == "" || hint == "No recovery hint available" {
				t.Errorf("ErrorCode(%q).RecoveryHint() = %q, should have specific hint", code, hint)
			}
		})
	}
}

func TestSDPError_Error(t *testing.T) {
	t.Run("with_cause", func(t *testing.T) {
		cause := errors.New("underlying error")
		err := New(ErrGitNotFound, cause)
		got := err.Error()

		if !strings.Contains(got, "ENV001") {
			t.Errorf("Error() should contain code ENV001, got %q", got)
		}
		if !strings.Contains(got, "underlying error") {
			t.Errorf("Error() should contain cause, got %q", got)
		}
	})

	t.Run("without_cause", func(t *testing.T) {
		err := New(ErrGitNotFound, nil)
		got := err.Error()

		if !strings.Contains(got, "ENV001") {
			t.Errorf("Error() should contain code ENV001, got %q", got)
		}
		if strings.Contains(got, ":") {
			// Should not have colon for cause when there's no cause
			parts := strings.Split(got, ":")
			if len(parts) > 1 && strings.Contains(parts[1], "underlying") {
				t.Errorf("Error() should not contain cause separator without cause, got %q", got)
			}
		}
	})
}

func TestSDPError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := New(ErrGitNotFound, cause)

	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, cause)
	}

	// Test with nil cause
	errNoCause := New(ErrGitNotFound, nil)
	if errNoCause.Unwrap() != nil {
		t.Errorf("Unwrap() with nil cause should return nil")
	}
}

func TestSDPError_Class(t *testing.T) {
	err := New(ErrGitNotFound, nil)
	if got := err.Class(); got != ClassEnvironment {
		t.Errorf("Class() = %v, want %v", got, ClassEnvironment)
	}

	err2 := New(ErrTestFailed, nil)
	if got := err2.Class(); got != ClassValidation {
		t.Errorf("Class() = %v, want %v", got, ClassValidation)
	}
}

func TestSDPError_RecoveryHint(t *testing.T) {
	err := New(ErrGitNotFound, nil)
	hint := err.RecoveryHint()

	if !strings.Contains(hint, "Install") {
		t.Errorf("RecoveryHint() should contain 'Install', got %q", hint)
	}
}

func TestSDPError_WithContext(t *testing.T) {
	err := New(ErrGitNotFound, nil)
	err.WithContext("file", "/path/to/config")
	err.WithContext("operation", "read")

	if err.Context == nil {
		t.Fatal("Context should not be nil after WithContext")
	}
	if err.Context["file"] != "/path/to/config" {
		t.Errorf("Context[file] = %q, want %q", err.Context["file"], "/path/to/config")
	}
	if err.Context["operation"] != "read" {
		t.Errorf("Context[operation] = %q, want %q", err.Context["operation"], "read")
	}
}

func TestNew(t *testing.T) {
	cause := errors.New("cause")
	err := New(ErrGitNotFound, cause)

	if err.Code != ErrGitNotFound {
		t.Errorf("Code = %v, want %v", err.Code, ErrGitNotFound)
	}
	if err.Cause != cause {
		t.Errorf("Cause = %v, want %v", err.Cause, cause)
	}
	if err.Message != ErrGitNotFound.Message() {
		t.Errorf("Message = %q, want %q", err.Message, ErrGitNotFound.Message())
	}
}

func TestNewf(t *testing.T) {
	err := Newf(ErrGitNotFound, "custom message: %s", "detail")

	if err.Code != ErrGitNotFound {
		t.Errorf("Code = %v, want %v", err.Code, ErrGitNotFound)
	}
	if err.Message != "custom message: detail" {
		t.Errorf("Message = %q, want %q", err.Message, "custom message: detail")
	}
}

func TestWrap(t *testing.T) {
	cause := errors.New("underlying")
	err := Wrap(ErrGitNotFound, cause, "wrapped message")

	if err.Code != ErrGitNotFound {
		t.Errorf("Code = %v, want %v", err.Code, ErrGitNotFound)
	}
	if err.Cause != cause {
		t.Errorf("Cause = %v, want %v", err.Cause, cause)
	}
	if err.Message != "wrapped message" {
		t.Errorf("Message = %q, want %q", err.Message, "wrapped message")
	}
}

func TestIsSDPError(t *testing.T) {
	sdpErr := New(ErrGitNotFound, nil)
	stdErr := errors.New("standard error")

	if !IsSDPError(sdpErr) {
		t.Error("IsSDPError(sdpErr) should be true")
	}
	if IsSDPError(stdErr) {
		t.Error("IsSDPError(stdErr) should be false")
	}
}

func TestGetCode(t *testing.T) {
	sdpErr := New(ErrGitNotFound, nil)
	stdErr := errors.New("standard error")

	if got := GetCode(sdpErr); got != ErrGitNotFound {
		t.Errorf("GetCode(sdpErr) = %v, want %v", got, ErrGitNotFound)
	}
	if got := GetCode(stdErr); got != ErrInternalError {
		t.Errorf("GetCode(stdErr) = %v, want %v", got, ErrInternalError)
	}
}

func TestGetClass(t *testing.T) {
	sdpErr := New(ErrGitNotFound, nil)
	stdErr := errors.New("standard error")

	if got := GetClass(sdpErr); got != ClassEnvironment {
		t.Errorf("GetClass(sdpErr) = %v, want %v", got, ClassEnvironment)
	}
	if got := GetClass(stdErr); got != ClassRuntime {
		t.Errorf("GetClass(stdErr) = %v, want %v", got, ClassRuntime)
	}
}

func TestErrorInterface(t *testing.T) {
	// Ensure SDPError implements error interface
	var err error = New(ErrGitNotFound, nil)
	if err.Error() == "" {
		t.Error("SDPError should implement error interface")
	}
}

func TestErrorChain(t *testing.T) {
	// Test error chain with multiple wraps
	cause := errors.New("root cause")
	err1 := New(ErrGitNotFound, cause)
	err2 := fmt.Errorf("wrapped: %w", err1)

	// Should be able to unwrap to find SDPError
	var sdpErr *SDPError
	if !errors.As(err2, &sdpErr) {
		t.Error("Should be able to unwrap to SDPError")
	}
	if sdpErr.Code != ErrGitNotFound {
		t.Errorf("Unwrapped code = %v, want %v", sdpErr.Code, ErrGitNotFound)
	}
}

func TestAllErrorCodesCovered(t *testing.T) {
	// Verify all error code ranges are consistent
	envCount := 0
	protoCount := 0
	depCount := 0
	valCount := 0
	runtimeCount := 0

	codes := []ErrorCode{
		ErrGitNotFound, ErrGoNotFound, ErrClaudeNotFound, ErrBeadsNotFound,
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
		ErrResourceExhausted, ErrUnexpectedState, ErrInternalError,
	}

	for _, code := range codes {
		switch code.Class() {
		case ClassEnvironment:
			envCount++
		case ClassProtocol:
			protoCount++
		case ClassDependency:
			depCount++
		case ClassValidation:
			valCount++
		case ClassRuntime:
			runtimeCount++
		}
	}

	// Verify we have codes in each class
	if envCount == 0 {
		t.Error("No environment error codes defined")
	}
	if protoCount == 0 {
		t.Error("No protocol error codes defined")
	}
	if depCount == 0 {
		t.Error("No dependency error codes defined")
	}
	if valCount == 0 {
		t.Error("No validation error codes defined")
	}
	if runtimeCount == 0 {
		t.Error("No runtime error codes defined")
	}

	t.Logf("Error code counts: ENV=%d, PROTO=%d, DEP=%d, VAL=%d, RUNTIME=%d",
		envCount, protoCount, depCount, valCount, runtimeCount)
}
