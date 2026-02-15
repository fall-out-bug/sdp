package evidence

import (
	"encoding/json"
	"testing"
)

func TestAcceptanceEvent_New(t *testing.T) {
	event := NewAcceptanceEvent("npm run start", true, "Server started on :3000", 30, "")

	if event.Command != "npm run start" {
		t.Error("command not set")
	}
	if !event.Passed {
		t.Error("should be passed")
	}
	if event.TimeoutSeconds != 30 {
		t.Error("timeout not set")
	}
}

func TestAcceptanceEvent_Failed(t *testing.T) {
	event := NewAcceptanceEvent("npm run start", false, "Connection refused", 30, "")

	if event.Passed {
		t.Error("should be failed")
	}
}

func TestAcceptanceEvent_JSON(t *testing.T) {
	original := NewAcceptanceEvent("./server", true, "OK", 30, "")
	original.Workstream = "00-054-07"

	data, _ := json.Marshal(original)
	var parsed AcceptanceEvent
	json.Unmarshal(data, &parsed)

	if parsed.Command != original.Command {
		t.Error("command mismatch")
	}
}

func TestAcceptanceEvent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		event   AcceptanceEvent
		wantErr bool
	}{
		{"valid", NewAcceptanceEvent("cmd", true, "", 30, ""), false},
		{"missing command", AcceptanceEvent{BaseEvent: NewBaseEvent(EventTypeAcceptance, "")}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApprovalEvent_New(t *testing.T) {
	event := NewApprovalEvent("user@example.com", "human", "")

	if event.Approver != "user@example.com" {
		t.Error("approver not set")
	}
	if event.Mode != "human" {
		t.Error("mode not set")
	}
}

func TestApprovalEvent_AutoMode(t *testing.T) {
	event := NewApprovalEvent("sdp-auto", "auto", "")
	event.Reasoning = "All quality gates passed"

	if event.Mode != "auto" {
		t.Error("mode should be auto")
	}
}

func TestApprovalEvent_JSON(t *testing.T) {
	original := NewApprovalEvent("user@test.com", "human", "")
	original.Reasoning = "LGTM"
	original.Workstream = "00-054-07"

	data, _ := json.Marshal(original)
	var parsed ApprovalEvent
	json.Unmarshal(data, &parsed)

	if parsed.Approver != original.Approver {
		t.Error("approver mismatch")
	}
}

func TestApprovalEvent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		event   ApprovalEvent
		wantErr bool
	}{
		{"valid human", NewApprovalEvent("user@test.com", "human", ""), false},
		{"valid auto", NewApprovalEvent("sdp", "auto", ""), false},
		{"missing approver", ApprovalEvent{BaseEvent: NewBaseEvent(EventTypeApproval, ""), Mode: "human"}, true},
		{"missing mode", ApprovalEvent{BaseEvent: NewBaseEvent(EventTypeApproval, ""), Approver: "user"}, true},
		{"invalid mode", ApprovalEvent{BaseEvent: NewBaseEvent(EventTypeApproval, ""), Approver: "user", Mode: "invalid"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
