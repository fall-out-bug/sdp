package evidence

import (
	"fmt"
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
		go func() {
			ev := Event{ID: "e", Type: "plan", Timestamp: "2026-02-09T12:00:00Z", WSID: "00-054-04"}
			_ = w.Append(&ev)
			done <- struct{}{}
		}()
	}
	for i := 0; i < 5; i++ {
		<-done
	}
}

// TestWriter_Append_Concurrent_ValidChain verifies hash chain integrity after concurrent Appends (AC: concurrent writers produce valid chain).
func TestWriter_Append_Concurrent_ValidChain(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events.jsonl")
	w, err := NewWriter(path)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	const n = 20
	done := make(chan error, n)
	for i := 0; i < n; i++ {
		go func(j int) {
			ev := Event{
				ID:        fmt.Sprintf("evt-%d", j),
				Type:      "verification",
				Timestamp: "2026-02-09T12:00:00Z",
				WSID:      "00-054-04",
			}
			done <- w.Append(&ev)
		}(i)
	}
	for i := 0; i < n; i++ {
		if err := <-done; err != nil {
			t.Fatalf("Append %d: %v", i, err)
		}
	}
	r := NewReader(path)
	if err := r.Verify(); err != nil {
		t.Fatalf("chain invalid after concurrent writes: %v", err)
	}
}

func TestNewWriter_EmptyPath(t *testing.T) {
	_, err := NewWriter("")
	if err == nil {
		t.Fatal("NewWriter(\"\"): want error, got nil")
	}
}

func TestWriter_Append_PathIsDir(t *testing.T) {
	dir := t.TempDir()
	pathAsDir := filepath.Join(dir, "events.jsonl")
	if err := os.MkdirAll(pathAsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	w, err := NewWriter(pathAsDir)
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	ev := Event{ID: "e1", Type: "plan", Timestamp: "2026-02-09T12:00:00Z", WSID: "00-054-04"}
	err = w.Append(&ev)
	if err == nil {
		t.Error("Append to path that is a directory: want error, got nil")
	}
}
