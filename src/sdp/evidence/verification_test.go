package evidence

import (
	"encoding/json"
	"testing"
)

func TestVerificationEvent_New(t *testing.T) {
	event := NewVerificationEvent("pytest", "pytest tests/", true, "")

	if event.Tool != "pytest" {
		t.Errorf("expected tool pytest, got %s", event.Tool)
	}
	if event.Command != "pytest tests/" {
		t.Errorf("expected command, got %s", event.Command)
	}
	if !event.Passed {
		t.Error("expected passed=true")
	}
}

func TestVerificationEvent_WithOutput(t *testing.T) {
	event := NewVerificationEvent("go", "go test ./...", true, "")
	event.Output = "ok  github.com/example  0.123s"

	if event.Output == "" {
		t.Error("output should be set")
	}
}

func TestVerificationEvent_WithCoverage(t *testing.T) {
	event := NewVerificationEvent("go", "go test -cover ./...", true, "")
	event.Coverage = 85.5

	if event.Coverage != 85.5 {
		t.Errorf("expected coverage 85.5, got %f", event.Coverage)
	}
}

func TestVerificationEvent_Failed(t *testing.T) {
	event := NewVerificationEvent("pytest", "pytest tests/", false, "")
	event.Output = "FAILED test_foo.py::test_bar"

	if event.Passed {
		t.Error("should be failed")
	}
}

func TestVerificationEvent_JSON(t *testing.T) {
	original := NewVerificationEvent("go", "go test ./...", true, "")
	original.Coverage = 90.0
	original.Duration = 1234
	original.Workstream = "00-054-06"

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var parsed VerificationEvent
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if parsed.Tool != original.Tool {
		t.Error("tool mismatch")
	}
	if parsed.Coverage != original.Coverage {
		t.Error("coverage mismatch")
	}
}

func TestVerificationEvent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		event   VerificationEvent
		wantErr bool
	}{
		{
			name:    "valid passed",
			event:   NewVerificationEvent("pytest", "pytest", true, ""),
			wantErr: false,
		},
		{
			name:    "valid failed",
			event:   NewVerificationEvent("pytest", "pytest", false, ""),
			wantErr: false,
		},
		{
			name: "missing tool",
			event: VerificationEvent{
				BaseEvent: NewBaseEvent(EventTypeVerification, ""),
				Command:   "cmd",
				Passed:    true,
			},
			wantErr: true,
		},
		{
			name: "missing command",
			event: VerificationEvent{
				BaseEvent: NewBaseEvent(EventTypeVerification, ""),
				Tool:      "pytest",
				Passed:    true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
