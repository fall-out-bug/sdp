package evidence

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// FilterByType returns events whose type matches (AC4: sdp log show --type=decision).
func FilterByType(events []Event, eventType string) []Event {
	if eventType == "" {
		return events
	}
	var out []Event
	for _, e := range events {
		if e.Type == eventType {
			out = append(out, e)
		}
	}
	return out
}

// FilterBySearch returns events whose Data contains query (AC5: full-text search in question/choice/rationale).
func FilterBySearch(events []Event, query string) []Event {
	if query == "" {
		return events
	}
	q := strings.ToLower(query)
	var out []Event
	for _, e := range events {
		if dataContains(e.Data, q) {
			out = append(out, e)
		}
	}
	return out
}

func dataContains(data interface{}, q string) bool {
	if data == nil {
		return false
	}
	m, ok := data.(map[string]interface{})
	if !ok {
		return false
	}
	for _, key := range []string{"question", "choice", "rationale"} {
		if v, ok := m[key]; ok {
			if s, ok := v.(string); ok && strings.Contains(strings.ToLower(s), q) {
				return true
			}
		}
	}
	return false
}

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
	m, ok := e.Data.(map[string]interface{})
	if !ok {
		return ""
	}
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
	if e.Type == "decision" {
		if v, ok := m["choice"]; ok {
			if s, ok := v.(string); ok && len(s) > 0 {
				if len(s) > 40 {
					return s[:37] + "..."
				}
				return s
			}
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
