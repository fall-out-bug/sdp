package evidence

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReader_Verify_ValidChain(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events.jsonl")
	w, err := NewWriter(path)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	for i := 0; i < 3; i++ {
		ev := Event{ID: "e", Type: "plan", Timestamp: "2026-02-09T12:00:00Z", WSID: "00-054-04"}
		if err := w.Append(&ev); err != nil {
			t.Fatalf("Append: %v", err)
		}
	}
	r := NewReader(path)
	err = r.Verify()
	if err != nil {
		t.Errorf("Verify: expected nil, got %v", err)
	}
}

func TestReader_Verify_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.jsonl")
	if err := os.WriteFile(path, []byte(""), 0600); err != nil {
		t.Fatalf("write: %v", err)
	}
	r := NewReader(path)
	if err := r.Verify(); err != nil {
		t.Errorf("Verify empty file: want nil, got %v", err)
	}
}

func TestReader_Verify_MissingFile(t *testing.T) {
	r := NewReader(filepath.Join(t.TempDir(), "nonexistent.jsonl"))
	if err := r.Verify(); err != nil {
		t.Errorf("Verify missing: want nil (skip), got %v", err)
	}
}

func TestReader_ReadAll(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events.jsonl")
	w, err := NewWriter(path)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	w.Append(&Event{ID: "e1", Type: "plan", Timestamp: "2026-02-09T12:00:00Z", WSID: "00-054-04"})
	r := NewReader(path)
	events, err := r.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("ReadAll: want 1 event, got %d", len(events))
	}
	if events[0].ID != "e1" {
		t.Errorf("ReadAll: want ID e1, got %s", events[0].ID)
	}
}

func TestReader_ReadAll_MissingFile(t *testing.T) {
	r := NewReader(filepath.Join(t.TempDir(), "nonexistent.jsonl"))
	events, err := r.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll missing: want nil error, got %v", err)
	}
	if events != nil {
		t.Errorf("ReadAll missing: want nil events, got %d", len(events))
	}
}

func TestReader_Verify_BrokenChain(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events.jsonl")
	w, err := NewWriter(path)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	ev := Event{ID: "e1", Type: "plan", Timestamp: "2026-02-09T12:00:00Z", WSID: "00-054-04"}
	if err := w.Append(&ev); err != nil {
		t.Fatalf("Append: %v", err)
	}
	// Corrupt: append a line with wrong prev_hash
	f, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	f.WriteString(`{"id":"e2","type":"plan","timestamp":"x","ws_id":"00-054-04","prev_hash":"wrong"}` + "\n")
	f.Close()

	r := NewReader(path)
	err = r.Verify()
	if err == nil {
		t.Error("Verify: expected error for broken chain")
	}
}
