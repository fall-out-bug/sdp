package evidenceenv

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type RoleGateResult struct {
	Role   string `json:"role"`
	OK     bool   `json:"ok"`
	Reason string `json:"reason"`
}

var roleEnvelopeKeys = []string{"run_id", "role", "status", "summary", "artifacts"}

func ValidateRoleLog(role, runID, log string) RoleGateResult {
	if strings.Contains(log, "ProviderModelNotFoundError") || strings.Contains(log, "Model not found") {
		return RoleGateResult{Role: role, OK: false, Reason: "model/provider resolution failure in logs"}
	}
	if strings.Contains(log, "Unable to connect") {
		return RoleGateResult{Role: role, OK: false, Reason: "provider connectivity failure in logs"}
	}

	env, err := extractEnvelope(log)
	if err != nil {
		return RoleGateResult{Role: role, OK: false, Reason: err.Error()}
	}

	if got, _ := env["role"].(string); got != role {
		return RoleGateResult{Role: role, OK: false, Reason: fmt.Sprintf("envelope role mismatch: got %q", got)}
	}
	if got, _ := env["run_id"].(string); got != runID {
		return RoleGateResult{Role: role, OK: false, Reason: fmt.Sprintf("envelope run_id mismatch: got %q", got)}
	}
	status, _ := env["status"].(string)
	if status != "ok" && status != "needs_changes" {
		return RoleGateResult{Role: role, OK: false, Reason: fmt.Sprintf("invalid envelope status: %q", status)}
	}

	return RoleGateResult{Role: role, OK: true, Reason: "ok"}
}

func extractEnvelope(log string) (map[string]any, error) {
	dec := json.NewDecoder(strings.NewReader(log))
	for {
		var v any
		if err := dec.Decode(&v); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			break
		}
		obj, ok := v.(map[string]any)
		if !ok {
			continue
		}
		if hasEnvelopeShape(obj) {
			return obj, nil
		}
	}

	// Fallback scanner for mixed text logs.
	for i := 0; i < len(log); i++ {
		if log[i] != '{' {
			continue
		}
		decoder := json.NewDecoder(strings.NewReader(log[i:]))
		var obj map[string]any
		if err := decoder.Decode(&obj); err != nil {
			continue
		}
		if hasEnvelopeShape(obj) {
			return obj, nil
		}
	}
	return nil, fmt.Errorf("missing valid role envelope in logs")
}

func hasEnvelopeShape(obj map[string]any) bool {
	for _, k := range roleEnvelopeKeys {
		if _, ok := obj[k]; !ok {
			return false
		}
	}
	if _, ok := obj["artifacts"].([]any); !ok {
		return false
	}
	return true
}
