package evidence

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// FilterByCommit returns events whose commit_sha matches (AC1).
func FilterByCommit(events []Event, commitSHA string) []Event {
	if commitSHA == "" {
		return events
	}
	var out []Event
	for _, e := range events {
		if e.CommitSHA == commitSHA {
			out = append(out, e)
		}
	}
	return out
}

// FilterByWS returns events whose ws_id matches (AC2).
func FilterByWS(events []Event, wsID string) []Event {
	if wsID == "" {
		return events
	}
	var out []Event
	for _, e := range events {
		if e.WSID == wsID {
			out = append(out, e)
		}
	}
	return out
}

// LastN returns the last n events (AC7: recent 20).
func LastN(events []Event, n int) []Event {
	if n <= 0 || len(events) <= n {
		return events
	}
	return events[len(events)-n:]
}

// FormatHuman returns human-readable timeline (AC3).
func FormatHuman(events []Event) string {
	var b strings.Builder
	for _, e := range events {
		ts := e.Timestamp
		if t, err := time.Parse(time.RFC3339, e.Timestamp); err == nil {
			ts = t.Format("2006-01-02 15:04:05")
		}
		key := keyData(e)
		fmt.Fprintf(&b, "  %s  %-12s  WS %s  %s\n", ts, e.Type, e.WSID, key)
	}
	return b.String()
}

func keyData(e Event) string {
	if e.Data == nil {
		return ""
	}
	// Minimal key data from Data map
	if m, ok := e.Data.(map[string]interface{}); ok {
		if v, ok := m["passed"]; ok {
			if b, ok := v.(bool); ok && b {
				return "PASS"
			}
			return "FAIL"
		}
		if v, ok := m["gate_name"]; ok {
			return fmt.Sprintf("gate: %v", v)
		}
		if v, ok := m["model_id"]; ok {
			return fmt.Sprintf("model: %v", v)
		}
	}
	return ""
}

// FormatJSON returns JSON array for machine consumption (AC4).
func FormatJSON(events []Event) (string, error) {
	b, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
