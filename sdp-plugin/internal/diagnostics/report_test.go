package diagnostics

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall-out-bug/sdp/internal/errors"
)

func TestNewGenerator(t *testing.T) {
	gen := NewGenerator("/tmp/project", "1.0.0")
	if gen == nil {
		t.Fatal("NewGenerator returned nil")
	}
	if gen.projectRoot != "/tmp/project" {
		t.Errorf("projectRoot = %q, want %q", gen.projectRoot, "/tmp/project")
	}
	if gen.sdpVersion != "1.0.0" {
		t.Errorf("sdpVersion = %q, want %q", gen.sdpVersion, "1.0.0")
	}
}

func TestGenerator_Generate(t *testing.T) {
	gen := NewGenerator("", "1.0.0")
	err := errors.New(errors.ErrGitNotFound, nil).WithContext("file", "config.yml")

	report := gen.Generate(err)

	if report == nil {
		t.Fatal("Generate returned nil")
	}
	if report.Timestamp == "" {
		t.Error("Report should have timestamp")
	}
	if report.Error == nil {
		t.Fatal("Report should have error info")
	}
	if report.Error.Code != "ENV001" {
		t.Errorf("Error.Code = %q, want %q", report.Error.Code, "ENV001")
	}
	if report.Environment == nil {
		t.Error("Report should have environment info")
	}
}

func TestGenerator_Generate_WithCause(t *testing.T) {
	gen := NewGenerator("", "1.0.0")
	cause := errors.New(errors.ErrFileNotWritable, nil)
	err := errors.New(errors.ErrGitNotFound, cause)

	report := gen.Generate(err)

	if report.Error.Cause == "" {
		t.Error("Error.Cause should not be empty")
	}
}

func TestGenerator_BuildErrorInfo(t *testing.T) {
	gen := NewGenerator("", "1.0.0")

	t.Run("sdp_error", func(t *testing.T) {
		err := errors.New(errors.ErrGitNotFound, nil).WithContext("file", "test.yml")
		info := gen.buildErrorInfo(err)

		if info.Code != "ENV001" {
			t.Errorf("Code = %q, want %q", info.Code, "ENV001")
		}
		if info.Class != "ENV" {
			t.Errorf("Class = %q, want %q", info.Class, "ENV")
		}
		if info.Context["file"] != "test.yml" {
			t.Errorf("Context[file] = %q, want %q", info.Context["file"], "test.yml")
		}
	})

	t.Run("standard_error", func(t *testing.T) {
		err := &testError{msg: "standard error"}
		info := gen.buildErrorInfo(err)

		if info.Code != "RUNTIME006" {
			t.Errorf("Code for non-SDP error should be RUNTIME006, got %s", info.Code)
		}
	})
}

func TestGenerator_BuildEnvironmentInfo(t *testing.T) {
	gen := NewGenerator("/test/project", "1.0.0")
	info := gen.buildEnvironmentInfo()

	if info.OS == "" {
		t.Error("OS should not be empty")
	}
	if info.Arch == "" {
		t.Error("Arch should not be empty")
	}
	if info.GoVersion == "" {
		t.Error("GoVersion should not be empty")
	}
	if info.SDPVersion != "1.0.0" {
		t.Errorf("SDPVersion = %q, want %q", info.SDPVersion, "1.0.0")
	}
	if info.ProjectRoot != "/test/project" {
		t.Errorf("ProjectRoot = %q, want %q", info.ProjectRoot, "/test/project")
	}
}

func TestGenerator_BuildEvidenceInfo(t *testing.T) {
	t.Run("no_project_root", func(t *testing.T) {
		gen := NewGenerator("", "1.0.0")
		info := gen.buildEvidenceInfo()

		if info.ChainIntegrity != "unknown" {
			t.Errorf("ChainIntegrity should be 'unknown' without project root")
		}
	})

	t.Run("no_log_file", func(t *testing.T) {
		tmpDir := t.TempDir()
		gen := NewGenerator(tmpDir, "1.0.0")
		info := gen.buildEvidenceInfo()

		if info.ChainIntegrity != "no_log" {
			t.Errorf("ChainIntegrity should be 'no_log' when log file missing")
		}
	})

	t.Run("with_log_file", func(t *testing.T) {
		tmpDir := t.TempDir()
		logDir := filepath.Join(tmpDir, ".sdp", "log")
		if err := os.MkdirAll(logDir, 0755); err != nil {
			t.Fatal(err)
		}

		logPath := filepath.Join(logDir, "events.jsonl")
		testData := `{"type":"test1"}
{"type":"test2"}`
		if err := os.WriteFile(logPath, []byte(testData), 0600); err != nil {
			t.Fatal(err)
		}

		gen := NewGenerator(tmpDir, "1.0.0")
		info := gen.buildEvidenceInfo()

		if info.EventCount != 2 {
			t.Errorf("EventCount = %d, want 2", info.EventCount)
		}
	})
}

func TestGenerator_BuildNextSteps(t *testing.T) {
	gen := NewGenerator("", "1.0.0")

	tests := []struct {
		class         errors.ErrorClass
		wantMinSteps  int
		wantContains  string
	}{
		{errors.ClassEnvironment, 3, "sdp doctor"},
		{errors.ClassProtocol, 3, "sdp parse"},
		{errors.ClassDependency, 3, "dependencies"},
		{errors.ClassValidation, 4, "sdp quality"},
		{errors.ClassRuntime, 4, "Retry"},
	}

	for _, tt := range tests {
		t.Run(string(tt.class), func(t *testing.T) {
			err := errors.New(errors.ErrGitNotFound, nil)
			// Override class by using specific error code
			var code errors.ErrorCode
			switch tt.class {
			case errors.ClassEnvironment:
				code = errors.ErrGitNotFound
			case errors.ClassProtocol:
				code = errors.ErrInvalidWorkstreamID
			case errors.ClassDependency:
				code = errors.ErrBlockedWorkstream
			case errors.ClassValidation:
				code = errors.ErrCoverageLow
			case errors.ClassRuntime:
				code = errors.ErrCommandFailed
			}
			err = errors.New(code, nil)

			steps := gen.buildNextSteps(err)

			if len(steps) < tt.wantMinSteps {
				t.Errorf("Got %d steps, want at least %d", len(steps), tt.wantMinSteps)
			}

			found := false
			for _, step := range steps {
				if strings.Contains(step.Description, tt.wantContains) ||
					strings.Contains(step.Command, tt.wantContains) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Steps should contain %q", tt.wantContains)
			}
		})
	}
}

func TestReport_FormatText(t *testing.T) {
	gen := NewGenerator("", "1.0.0")
	err := errors.New(errors.ErrGitNotFound, nil).WithContext("file", "config.yml")
	report := gen.Generate(err)

	output := report.FormatText()

	// Check for expected sections
	if !strings.Contains(output, "=== SDP Diagnostics Report ===") {
		t.Error("Output should contain report header")
	}
	if !strings.Contains(output, "--- Error ---") {
		t.Error("Output should contain Error section")
	}
	if !strings.Contains(output, "--- Environment ---") {
		t.Error("Output should contain Environment section")
	}
	if !strings.Contains(output, "--- Next Steps ---") {
		t.Error("Output should contain Next Steps section")
	}
	if !strings.Contains(output, "ENV001") {
		t.Error("Output should contain error code")
	}
}

func TestReport_FormatJSON(t *testing.T) {
	gen := NewGenerator("", "1.0.0")
	sdpErr := errors.New(errors.ErrGitNotFound, nil)
	report := gen.Generate(sdpErr)

	output, fmtErr := report.FormatJSON()
	if fmtErr != nil {
		t.Fatalf("FormatJSON failed: %v", fmtErr)
	}

	// Verify it's valid JSON
	var parsed Report
	if jsonErr := json.Unmarshal([]byte(output), &parsed); jsonErr != nil {
		t.Fatalf("Output is not valid JSON: %v", jsonErr)
	}

	// Check required fields
	if !strings.Contains(output, `"code"`) {
		t.Error("JSON should contain code field")
	}
	if !strings.Contains(output, `"timestamp"`) {
		t.Error("JSON should contain timestamp field")
	}
}

func TestReport_Redact(t *testing.T) {
	gen := NewGenerator("", "1.0.0")
	err := errors.New(errors.ErrGitNotFound, nil).
		WithContext("file", "config.yml").
		WithContext("api_key", "secret123")
	report := gen.Generate(err)

	report.Redact([]string{"api_key"})

	if report.Error.Context["api_key"] != "[REDACTED]" {
		t.Errorf("api_key should be redacted, got %q", report.Error.Context["api_key"])
	}
	if report.Error.Context["file"] != "config.yml" {
		t.Error("file should not be redacted")
	}
}

func TestReport_Save(t *testing.T) {
	gen := NewGenerator("", "1.0.0")
	sdpErr := errors.New(errors.ErrGitNotFound, nil)
	report := gen.Generate(sdpErr)

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "diagnostics", "report.json")

	if saveErr := report.Save(path); saveErr != nil {
		t.Fatalf("Save failed: %v", saveErr)
	}

	// Verify file was created
	if _, statErr := os.Stat(path); os.IsNotExist(statErr) {
		t.Error("Report file was not created")
	}

	// Verify content
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("ReadFile failed: %v", readErr)
	}

	if !strings.Contains(string(data), "ENV001") {
		t.Error("Saved report should contain error code")
	}
}

func TestGenerateForError(t *testing.T) {
	err := errors.New(errors.ErrCoverageLow, nil)
	report := GenerateForError(err, "/test", "1.0.0")

	if report == nil {
		t.Fatal("GenerateForError returned nil")
	}
	if report.Error.Code != "VAL001" {
		t.Errorf("Error.Code = %q, want %q", report.Error.Code, "VAL001")
	}
	if report.Environment.ProjectRoot != "/test" {
		t.Errorf("ProjectRoot = %q, want %q", report.Environment.ProjectRoot, "/test")
	}
}

func TestReport_WithRecovery(t *testing.T) {
	gen := NewGenerator("", "1.0.0")
	err := errors.New(errors.ErrGitNotFound, nil)
	report := gen.Generate(err)

	if report.Recovery == nil {
		t.Error("Report should have recovery playbook")
	}
	if report.Recovery.Title != "Install Git" {
		t.Errorf("Recovery.Title = %q, want %q", report.Recovery.Title, "Install Git")
	}
}

func TestReport_EvidenceSection(t *testing.T) {
	gen := NewGenerator("", "1.0.0")
	err := errors.New(errors.ErrGitNotFound, nil)
	report := gen.Generate(err)

	if report.Evidence == nil {
		t.Error("Report should have evidence info")
	}
}

// testError is a simple error implementation for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
