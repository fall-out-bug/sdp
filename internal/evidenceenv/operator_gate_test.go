package evidenceenv

import (
	"strings"
	"testing"
)

func TestValidateRoleLogOK(t *testing.T) {
	log := `noise line
{"run_id":"run-1","role":"analyst","status":"ok","summary":"done","artifacts":[{"id":"a1"}]}
more noise`
	res := ValidateRoleLog("analyst", "run-1", log)
	if !res.OK {
		t.Fatalf("expected ok, got %+v", res)
	}
}

func TestValidateRoleLogProviderError(t *testing.T) {
	log := `ProviderModelNotFoundError: Model not found: zai/glm-5.`
	res := ValidateRoleLog("coder", "run-1", log)
	if res.OK {
		t.Fatalf("expected failure for provider error")
	}
}

func TestValidateRoleLogConnectivityError(t *testing.T) {
	log := `Error: Unable to connect. Is the computer able to access the url?`
	res := ValidateRoleLog("coder", "run-1", log)
	if res.OK {
		t.Fatalf("expected failure for connectivity error")
	}
}

func TestValidateRoleLogRoleMismatch(t *testing.T) {
	log := `{"run_id":"run-1","role":"analyst","status":"ok","summary":"done","artifacts":[]}`
	res := ValidateRoleLog("reviewer", "run-1", log)
	if res.OK {
		t.Fatalf("expected role mismatch failure")
	}
}

func TestValidateRoleLogNeedsChanges(t *testing.T) {
	log := `{"run_id":"run-1","role":"coder","status":"needs_changes","summary":"fix requested","artifacts":[{"id":"a1"}]}`
	res := ValidateRoleLog("coder", "run-1", log)
	if !res.OK {
		t.Fatalf("needs_changes should pass: %+v", res)
	}
}

func TestValidateRoleLogRunIDMismatch(t *testing.T) {
	log := `{"run_id":"run-2","role":"analyst","status":"ok","summary":"done","artifacts":[]}`
	res := ValidateRoleLog("analyst", "run-1", log)
	if res.OK {
		t.Fatalf("expected run_id mismatch failure")
	}
}

func TestValidateRoleLogInvalidStatus(t *testing.T) {
	log := `{"run_id":"run-1","role":"coder","status":"failed","summary":"err","artifacts":[]}`
	res := ValidateRoleLog("coder", "run-1", log)
	if res.OK {
		t.Fatalf("expected invalid status failure: %+v", res)
	}
	if !strings.Contains(res.Reason, "invalid envelope status") {
		t.Errorf("reason: %s", res.Reason)
	}
}

func TestValidateRoleLogMissingEnvelope(t *testing.T) {
	log := `no json here at all`
	res := ValidateRoleLog("analyst", "run-1", log)
	if res.OK {
		t.Fatalf("expected missing envelope failure")
	}
	if !strings.Contains(res.Reason, "missing valid role envelope") {
		t.Errorf("reason: %s", res.Reason)
	}
}
