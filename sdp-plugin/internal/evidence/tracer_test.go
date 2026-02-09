package evidence

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFilterByCommit(t *testing.T) {
	events := []Event{
		{ID: "e1", CommitSHA: "abc"},
		{ID: "e2", CommitSHA: "def"},
		{ID: "e3", CommitSHA: "abc"},
	}
	got := FilterByCommit(events, "abc")
	if len(got) != 2 {
		t.Errorf("FilterByCommit(abc): want 2, got %d", len(got))
	}
	got = FilterByCommit(events, "")
	if len(got) != 3 {
		t.Errorf("FilterByCommit(''): want all 3, got %d", len(got))
	}
}

func TestFilterByWS(t *testing.T) {
	events := []Event{
		{ID: "e1", WSID: "00-054-03"},
		{ID: "e2", WSID: "00-054-04"},
		{ID: "e3", WSID: "00-054-03"},
	}
	got := FilterByWS(events, "00-054-03")
	if len(got) != 2 {
		t.Errorf("FilterByWS(00-054-03): want 2, got %d", len(got))
	}
}

func TestFilterByType(t *testing.T) {
	events := []Event{
		{ID: "e1", Type: "decision"},
		{ID: "e2", Type: "plan"},
		{ID: "e3", Type: "decision"},
	}
	got := FilterByType(events, "decision")
	if len(got) != 2 {
		t.Errorf("FilterByType(decision): want 2, got %d", len(got))
	}
	got = FilterByType(events, "")
	if len(got) != 3 {
		t.Errorf("FilterByType(''): want all 3, got %d", len(got))
	}
}

func TestFilterBySearch(t *testing.T) {
	events := []Event{
		{ID: "e1", Type: "decision", Data: map[string]interface{}{"question": "Auth?", "choice": "JWT", "rationale": "Stateless"}},
		{ID: "e2", Type: "decision", Data: map[string]interface{}{"question": "DB?", "choice": "Postgres", "rationale": "ACID"}},
	}
	got := FilterBySearch(events, "jwt")
	if len(got) != 1 {
		t.Errorf("FilterBySearch(jwt): want 1, got %d", len(got))
	}
	got = FilterBySearch(events, "database")
	if len(got) != 0 {
		t.Errorf("FilterBySearch(database): want 0, got %d", len(got))
	}
	got = FilterBySearch(events, "postgres")
	if len(got) != 1 {
		t.Errorf("FilterBySearch(postgres): want 1, got %d", len(got))
	}
}

func TestLastN(t *testing.T) {
	events := []Event{{ID: "e1"}, {ID: "e2"}, {ID: "e3"}, {ID: "e4"}, {ID: "e5"}}
	got := LastN(events, 2)
	if len(got) != 2 {
		t.Fatalf("LastN(2): want 2, got %d", len(got))
	}
	if got[0].ID != "e4" || got[1].ID != "e5" {
		t.Errorf("LastN(2): want e4,e5, got %s,%s", got[0].ID, got[1].ID)
	}
	got = LastN(events, 20)
	if len(got) != 5 {
		t.Errorf("LastN(20): want all 5, got %d", len(got))
	}
}

func TestFormatHuman(t *testing.T) {
	events := []Event{
		{Type: "plan", Timestamp: "2026-02-09T12:00:00Z", WSID: "00-054-03"},
	}
	out := FormatHuman(events)
	if !strings.Contains(out, "plan") || !strings.Contains(out, "00-054-03") {
		t.Errorf("FormatHuman: want plan and ws_id, got %q", out)
	}
}

func TestFormatJSON(t *testing.T) {
	events := []Event{{ID: "e1", Type: "plan"}}
	out, err := FormatJSON(events)
	if err != nil {
		t.Fatalf("FormatJSON: %v", err)
	}
	var decoded []Event
	if err := json.Unmarshal([]byte(out), &decoded); err != nil {
		t.Fatalf("FormatJSON output not valid JSON: %v", err)
	}
	if len(decoded) != 1 || decoded[0].ID != "e1" {
		t.Errorf("FormatJSON roundtrip: got %+v", decoded)
	}
}
