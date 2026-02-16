package recovery

import (
	"strings"
	"testing"

	"github.com/fall-out-bug/sdp/internal/errors"
)

func TestNewPlaybookRegistry(t *testing.T) {
	r := NewPlaybookRegistry()
	if r == nil {
		t.Fatal("NewPlaybookRegistry returned nil")
	}
	if len(r.playbooks) == 0 {
		t.Error("Registry should have built-in playbooks")
	}
}

func TestPlaybookRegistry_Get(t *testing.T) {
	r := NewPlaybookRegistry()

	tests := []struct {
		code     errors.ErrorCode
		wantNil  bool
		wantName string
	}{
		{errors.ErrGitNotFound, false, "Install Git"},
		{errors.ErrBeadsNotFound, false, "Install Beads CLI"},
		{errors.ErrInvalidWorkstreamID, false, "Fix Workstream ID Format"},
		{errors.ErrBlockedWorkstream, false, "Resolve Blocked Workstream"},
		{errors.ErrCoverageLow, false, "Increase Test Coverage"},
		{errors.ErrFileTooLarge, false, "Split Large File"},
		{errors.ErrCommandFailed, false, "Debug Command Failure"},
		{errors.ErrInternalError, false, "Report Internal Error"},
	}

	for _, tt := range tests {
		t.Run(string(tt.code), func(t *testing.T) {
			pb := r.Get(tt.code)
			if (pb == nil) != tt.wantNil {
				t.Errorf("Get(%s) returned nil=%v, want nil=%v", tt.code, pb == nil, tt.wantNil)
				return
			}
			if pb != nil && pb.Title != tt.wantName {
				t.Errorf("Get(%s).Title = %q, want %q", tt.code, pb.Title, tt.wantName)
			}
		})
	}
}

func TestPlaybookRegistry_Get_ClassDefault(t *testing.T) {
	r := NewPlaybookRegistry()

	// Test that unknown codes fall back to class defaults
	tests := []struct {
		class     errors.ErrorClass
		wantFound bool
	}{
		{errors.ClassEnvironment, true},
		{errors.ClassProtocol, true},
		{errors.ClassDependency, true},
		{errors.ClassValidation, true},
		{errors.ClassRuntime, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.class), func(t *testing.T) {
			pb := r.classDefaults[tt.class]
			if (pb != nil) != tt.wantFound {
				t.Errorf("Class default for %s exists=%v, want %v", tt.class, pb != nil, tt.wantFound)
			}
		})
	}
}

func TestPlaybookRegistry_GetForError(t *testing.T) {
	r := NewPlaybookRegistry()

	t.Run("sdp_error", func(t *testing.T) {
		err := errors.New(errors.ErrGitNotFound, nil)
		pb := r.GetForError(err)
		if pb == nil {
			t.Error("GetForError should return playbook for SDPError")
		}
		if pb != nil && pb.Code != errors.ErrGitNotFound {
			t.Errorf("Playbook code = %s, want %s", pb.Code, errors.ErrGitNotFound)
		}
	})

	t.Run("standard_error", func(t *testing.T) {
		err := &testError{msg: "standard error"}
		pb := r.GetForError(err)
		// Should fall back to internal error playbook
		if pb == nil {
			t.Error("GetForError should return playbook for standard error")
		}
	})
}

func TestPlaybookRegistry_Register(t *testing.T) {
	r := NewPlaybookRegistry()

	newCode := errors.ErrorCode("TEST001")
	pb := &Playbook{
		Code:  newCode,
		Title: "Test Playbook",
	}

	r.Register(pb)

	retrieved := r.Get(newCode)
	if retrieved == nil {
		t.Fatal("Registered playbook not found")
	}
	if retrieved.Title != "Test Playbook" {
		t.Errorf("Retrieved title = %q, want %q", retrieved.Title, "Test Playbook")
	}
}

func TestPlaybookRegistry_List(t *testing.T) {
	r := NewPlaybookRegistry()
	list := r.List()

	if len(list) == 0 {
		t.Error("List should return non-empty slice")
	}
}

func TestFormatPlaybook(t *testing.T) {
	r := NewPlaybookRegistry()
	pb := r.Get(errors.ErrGitNotFound)

	output := FormatPlaybook(pb)

	// Check for expected content
	if !strings.Contains(output, "Install Git") {
		t.Error("Output should contain playbook title")
	}
	if !strings.Contains(output, "Quick Fix:") {
		t.Error("Output should contain 'Quick Fix:' section")
	}
	if !strings.Contains(output, "brew install git") {
		t.Error("Output should contain command")
	}
}

func TestFormatPlaybook_WithDeepPath(t *testing.T) {
	r := NewPlaybookRegistry()
	pb := r.Get(errors.ErrGitNotFound)

	output := FormatPlaybook(pb)

	if !strings.Contains(output, "Full Recovery:") {
		t.Error("Output should contain 'Full Recovery:' section for playbooks with deep path")
	}
}

func TestFormatPlaybook_WithRelatedDocs(t *testing.T) {
	r := NewPlaybookRegistry()
	pb := r.Get(errors.ErrGitNotFound)

	output := FormatPlaybook(pb)

	if !strings.Contains(output, "Related Docs:") {
		t.Error("Output should contain 'Related Docs:' section")
	}
}

func TestGlobalRegistry(t *testing.T) {
	pb := GetPlaybook(errors.ErrGitNotFound)
	if pb == nil {
		t.Error("GetPlaybook should return playbook from global registry")
	}
}

func TestGetPlaybookForError(t *testing.T) {
	err := errors.New(errors.ErrCoverageLow, nil)
	pb := GetPlaybookForError(err)

	if pb == nil {
		t.Fatal("GetPlaybookForError returned nil")
	}
	if pb.Code != errors.ErrCoverageLow {
		t.Errorf("Playbook code = %s, want %s", pb.Code, errors.ErrCoverageLow)
	}
}

func TestPlaybook_HasRequiredFields(t *testing.T) {
	r := NewPlaybookRegistry()

	// All built-in playbooks should have required fields
	for code, pb := range r.playbooks {
		t.Run(string(code), func(t *testing.T) {
			if pb.Code == "" {
				t.Error("Playbook should have Code")
			}
			if pb.Title == "" {
				t.Error("Playbook should have Title")
			}
			if pb.Severity == "" {
				t.Error("Playbook should have Severity")
			}
			if len(pb.FastPath) == 0 {
				t.Error("Playbook should have FastPath steps")
			}
		})
	}
}

func TestPlaybook_StepFields(t *testing.T) {
	r := NewPlaybookRegistry()
	pb := r.Get(errors.ErrGitNotFound)

	if pb == nil {
		t.Fatal("Playbook not found")
	}

	for _, step := range pb.FastPath {
		if step.Description == "" {
			t.Error("Step should have Description")
		}
		if step.Order == 0 {
			t.Error("Step should have non-zero Order")
		}
	}
}

func TestAllErrorCodes_HavePlaybooks(t *testing.T) {
	r := NewPlaybookRegistry()

	// Test that common error codes have playbooks
	codes := []errors.ErrorCode{
		errors.ErrGitNotFound,
		errors.ErrGoNotFound,
		errors.ErrClaudeNotFound,
		errors.ErrBeadsNotFound,
		errors.ErrPermissionDenied,
		errors.ErrWorktreeNotFound,
		errors.ErrConfigNotFound,
		errors.ErrInvalidWorkstreamID,
		errors.ErrInvalidFeatureID,
		errors.ErrMalformedYAML,
		errors.ErrHashChainBroken,
		errors.ErrSessionCorrupted,
		errors.ErrBlockedWorkstream,
		errors.ErrCircularDependency,
		errors.ErrCoverageLow,
		errors.ErrFileTooLarge,
		errors.ErrTestFailed,
		errors.ErrDriftDetected,
		errors.ErrCommandFailed,
		errors.ErrTimeoutExceeded,
		errors.ErrInternalError,
	}

	for _, code := range codes {
		t.Run(string(code), func(t *testing.T) {
			pb := r.Get(code)
			// Either specific playbook or class default should exist
			classDefault := r.classDefaults[code.Class()]
			if pb == nil && classDefault == nil {
				t.Errorf("No playbook or class default for %s", code)
			}
		})
	}
}

// testError is a simple error implementation for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
