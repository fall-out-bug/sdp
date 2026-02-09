package evidence

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriter_Append_FirstEvent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".sdp", "log", "events.jsonl")
	w, err := NewWriter(path)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	ev := Event{ID: "e1", Type: "plan", Timestamp: "2026-02-09T12:00:00Z", WSID: "00-054-04"}
	if err := w.Append(&ev); err != nil {
		t.Fatalf("Append: %v", err)
	}
	if ev.PrevHash != genesisHash {
		t.Errorf("first event PrevHash: want %q, got %q", genesisHash, ev.PrevHash)
	}
	// Second event should get prev_hash = hash of first line
	ev2 := Event{ID: "e2", Type: "verification", Timestamp: "2026-02-09T12:01:00Z", WSID: "00-054-04"}
	if err := w.Append(&ev2); err != nil {
		t.Fatalf("Append 2: %v", err)
	}
	if ev2.PrevHash == "" || ev2.PrevHash == genesisHash {
		t.Errorf("second event PrevHash should be hash of first, got %q", ev2.PrevHash)
	}
}

func TestWriter_Append_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".sdp", "log", "events.jsonl")
	_, err := NewWriter(path)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		t.Errorf("dir not created: %v", err)
	}
}

func TestWriter_Append_ResumesExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events.jsonl")
	w1, err := NewWriter(path)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	ev1 := Event{ID: "e1", Type: "plan", Timestamp: "2026-02-09T12:00:00Z", WSID: "00-054-04"}
	if err := w1.Append(&ev1); err != nil {
		t.Fatalf("Append: %v", err)
	}
	// New writer on same file should resume lastHash
	w2, err := NewWriter(path)
	if err != nil {
		t.Fatalf("NewWriter 2: %v", err)
	}
	ev2 := Event{ID: "e2", Type: "verification", Timestamp: "2026-02-09T12:01:00Z", WSID: "00-054-04"}
	if err := w2.Append(&ev2); err != nil {
		t.Fatalf("Append 2: %v", err)
	}
	if ev2.PrevHash == genesisHash {
		t.Error("second writer should have resumed lastHash from file")
	}
}

func TestWriter_Append_Concurrent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events.jsonl")
	w, err := NewWriter(path)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	done := make(chan struct{})
	for i := 0; i < 5; i++ {
		go func(j int) {
			ev := Event{ID: "e", Type: "plan", Timestamp: "2026-02-09T12:00:00Z", WSID: "00-054-04"}
			_ = w.Append(&ev)
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 5; i++ {
		<-done
	}
}
