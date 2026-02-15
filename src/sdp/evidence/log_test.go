package evidence

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEvidenceLog_New(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")

	log, err := NewEvidenceLog(logPath)
	if err != nil {
		t.Fatalf("NewEvidenceLog failed: %v", err)
	}

	if log == nil {
		t.Fatal("log should not be nil")
	}
}

func TestEvidenceLog_Append(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")

	log, _ := NewEvidenceLog(logPath)

	event := NewBaseEvent(EventTypeDecision, "")
	if err := log.Append(event); err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Fatal("log file should be created")
	}
}

func TestEvidenceLog_AppendChain(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")

	log, _ := NewEvidenceLog(logPath)

	// Append genesis event
	event1 := NewBaseEvent(EventTypePlan, "")
	if err := log.Append(event1); err != nil {
		t.Fatalf("Append event1 failed: %v", err)
	}

	// Get last hash for next event
	lastHash, err := log.GetLastHash()
	if err != nil {
		t.Fatalf("GetLastHash failed: %v", err)
	}

	// Append second event linked to first
	event2 := NewBaseEvent(EventTypeDecision, lastHash)
	if err := log.Append(event2); err != nil {
		t.Fatalf("Append event2 failed: %v", err)
	}

	// Verify chain
	events, err := log.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}

	// Second event should reference first event's hash
	if events[1].PrevHash != events[0].Hash {
		t.Error("chain link broken: event2.prev_hash != event1.hash")
	}
}

func TestEvidenceLog_ReadAll(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")

	log, _ := NewEvidenceLog(logPath)

	// Append multiple events
	for i := 0; i < 5; i++ {
		event := NewBaseEvent(EventTypeDecision, "")
		if err := log.Append(event); err != nil {
			t.Fatalf("Append failed: %v", err)
		}
	}

	events, err := log.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}

	if len(events) != 5 {
		t.Errorf("expected 5 events, got %d", len(events))
	}
}

func TestEvidenceLog_ReadEmpty(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")

	log, _ := NewEvidenceLog(logPath)

	events, err := log.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll failed on empty log: %v", err)
	}

	if len(events) != 0 {
		t.Errorf("expected 0 events from empty log, got %d", len(events))
	}
}

func TestEvidenceLog_GetLastHash(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")

	log, _ := NewEvidenceLog(logPath)

	// Empty log should return empty hash
	hash, err := log.GetLastHash()
	if err != nil {
		t.Fatalf("GetLastHash on empty log failed: %v", err)
	}
	if hash != "" {
		t.Error("empty log should return empty hash")
	}

	// After append, should return event hash
	event := NewBaseEvent(EventTypeDecision, "")
	log.Append(event)

	hash, err = log.GetLastHash()
	if err != nil {
		t.Fatalf("GetLastHash failed: %v", err)
	}
	if hash != event.Hash {
		t.Error("should return last event's hash")
	}
}

func TestEvidenceLog_AppendAutoChain(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")

	log, _ := NewEvidenceLog(logPath)

	// AppendEvent automatically chains to previous
	event1, err := log.AppendEvent(EventTypePlan, nil)
	if err != nil {
		t.Fatalf("AppendEvent failed: %v", err)
	}

	if event1.PrevHash != "" {
		t.Error("first event should have empty prev_hash")
	}

	event2, err := log.AppendEvent(EventTypeDecision, nil)
	if err != nil {
		t.Fatalf("AppendEvent failed: %v", err)
	}

	if event2.PrevHash != event1.Hash {
		t.Error("second event should chain to first")
	}
}

func TestEvidenceLog_ConcurrentAppend(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")

	log, _ := NewEvidenceLog(logPath)

	// Concurrent appends (sequential due to mutex, but tests thread safety)
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func() {
			_, err := log.AppendEvent(EventTypeDecision, nil)
			if err != nil {
				t.Errorf("concurrent append failed: %v", err)
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	events, _ := log.ReadAll()
	if len(events) != 10 {
		t.Errorf("expected 10 events, got %d", len(events))
	}
}

func TestEvidenceLog_Truncate(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")

	log, _ := NewEvidenceLog(logPath)

	// Add events
	for i := 0; i < 5; i++ {
		log.AppendEvent(EventTypeDecision, nil)
	}

	// Truncate to 3 events
	if err := log.Truncate(3); err != nil {
		t.Fatalf("Truncate failed: %v", err)
	}

	events, _ := log.ReadAll()
	if len(events) != 3 {
		t.Errorf("expected 3 events after truncate, got %d", len(events))
	}
}

func TestEvidenceLog_Exists(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")

	log, _ := NewEvidenceLog(logPath)

	if log.Exists() {
		t.Error("log should not exist before any writes")
	}

	log.AppendEvent(EventTypeDecision, nil)

	if !log.Exists() {
		t.Error("log should exist after write")
	}
}

func TestEvidenceLog_Path(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "events.jsonl")

	log, _ := NewEvidenceLog(logPath)

	if log.Path() != logPath {
		t.Errorf("expected path %s, got %s", logPath, log.Path())
	}
}
