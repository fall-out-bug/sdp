package coordination

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStore_New(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "coord-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logPath := filepath.Join(tmpDir, "events.jsonl")
	store, err := NewStore(logPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	if store == nil {
		t.Error("Expected non-nil store")
	}
}

func TestStore_AppendAndRead(t *testing.T) {
	// AC3: Event store reader for replaying agent history
	tmpDir, err := os.MkdirTemp("", "coord-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logPath := filepath.Join(tmpDir, "events.jsonl")
	store, err := NewStore(logPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Append events
	events := []*AgentEvent{
		{
			ID:        "evt-1",
			Type:      EventTypeAgentStart,
			AgentID:   "agent-1",
			Role:      "implementer",
			TaskID:    "task-1",
			Timestamp: time.Date(2026, 2, 12, 10, 0, 0, 0, time.UTC),
			Payload:   map[string]interface{}{"ws_id": "00-051-01"},
		},
		{
			ID:        "evt-2",
			Type:      EventTypeAgentAction,
			AgentID:   "agent-1",
			Role:      "implementer",
			TaskID:    "task-1",
			Timestamp: time.Date(2026, 2, 12, 10, 5, 0, 0, time.UTC),
			Payload:   map[string]interface{}{"action": "code_generation"},
		},
		{
			ID:        "evt-3",
			Type:      EventTypeAgentComplete,
			AgentID:   "agent-1",
			Role:      "implementer",
			TaskID:    "task-1",
			Timestamp: time.Date(2026, 2, 12, 10, 30, 0, 0, time.UTC),
			Payload:   map[string]interface{}{"result": "success"},
		},
	}

	for _, e := range events {
		if err := store.Append(e); err != nil {
			t.Fatalf("Failed to append event: %v", err)
		}
	}

	// Read all events
	readEvents, err := store.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read events: %v", err)
	}

	if len(readEvents) != 3 {
		t.Errorf("Expected 3 events, got %d", len(readEvents))
	}
}

func TestStore_FilterByAgent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "coord-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logPath := filepath.Join(tmpDir, "events.jsonl")
	store, err := NewStore(logPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Add events from different agents
	events := []*AgentEvent{
		{ID: "1", Type: EventTypeAgentStart, AgentID: "agent-1", Role: "impl", Timestamp: time.Now()},
		{ID: "2", Type: EventTypeAgentStart, AgentID: "agent-2", Role: "reviewer", Timestamp: time.Now()},
		{ID: "3", Type: EventTypeAgentAction, AgentID: "agent-1", Role: "impl", Timestamp: time.Now()},
	}

	for _, e := range events {
		if err := store.Append(e); err != nil {
			t.Fatalf("Failed to append: %v", err)
		}
	}

	// Filter by agent-1
	filtered, err := store.FilterByAgent("agent-1")
	if err != nil {
		t.Fatalf("FilterByAgent failed: %v", err)
	}

	if len(filtered) != 2 {
		t.Errorf("Expected 2 events for agent-1, got %d", len(filtered))
	}
}

func TestStore_FilterByTask(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "coord-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logPath := filepath.Join(tmpDir, "events.jsonl")
	store, err := NewStore(logPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	events := []*AgentEvent{
		{ID: "1", Type: EventTypeAgentStart, AgentID: "a1", TaskID: "task-1", Timestamp: time.Now()},
		{ID: "2", Type: EventTypeAgentStart, AgentID: "a2", TaskID: "task-2", Timestamp: time.Now()},
		{ID: "3", Type: EventTypeAgentComplete, AgentID: "a1", TaskID: "task-1", Timestamp: time.Now()},
	}

	for _, e := range events {
		if err := store.Append(e); err != nil {
			t.Fatalf("Failed to append: %v", err)
		}
	}

	filtered, err := store.FilterByTask("task-1")
	if err != nil {
		t.Fatalf("FilterByTask failed: %v", err)
	}

	if len(filtered) != 2 {
		t.Errorf("Expected 2 events for task-1, got %d", len(filtered))
	}
}

func TestStore_GetAggregatedStats(t *testing.T) {
	// AC4: Event aggregation for status summaries
	tmpDir, err := os.MkdirTemp("", "coord-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logPath := filepath.Join(tmpDir, "events.jsonl")
	store, err := NewStore(logPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Add events
	events := []*AgentEvent{
		{ID: "1", Type: EventTypeAgentStart, AgentID: "a1", Role: "impl", TaskID: "t1", Timestamp: time.Now()},
		{ID: "2", Type: EventTypeAgentStart, AgentID: "a2", Role: "reviewer", TaskID: "t1", Timestamp: time.Now()},
		{ID: "3", Type: EventTypeAgentComplete, AgentID: "a1", Role: "impl", TaskID: "t1", Timestamp: time.Now()},
		{ID: "4", Type: EventTypeAgentError, AgentID: "a2", Role: "reviewer", TaskID: "t1", Timestamp: time.Now()},
	}

	for _, e := range events {
		if err := store.Append(e); err != nil {
			t.Fatalf("Failed to append: %v", err)
		}
	}

	stats, err := store.GetAggregatedStats()
	if err != nil {
		t.Fatalf("GetAggregatedStats failed: %v", err)
	}

	if stats.TotalEvents != 4 {
		t.Errorf("Expected 4 total events, got %d", stats.TotalEvents)
	}
	if stats.ByType[EventTypeAgentStart] != 2 {
		t.Errorf("Expected 2 start events, got %d", stats.ByType[EventTypeAgentStart])
	}
	if stats.ByType[EventTypeAgentComplete] != 1 {
		t.Errorf("Expected 1 complete event, got %d", stats.ByType[EventTypeAgentComplete])
	}
	if stats.ByType[EventTypeAgentError] != 1 {
		t.Errorf("Expected 1 error event, got %d", stats.ByType[EventTypeAgentError])
	}
}

func TestStore_VerifyHashChain(t *testing.T) {
	// AC5: Hash chain verification for event integrity
	tmpDir, err := os.MkdirTemp("", "coord-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logPath := filepath.Join(tmpDir, "events.jsonl")
	store, err := NewStore(logPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	ts := time.Date(2026, 2, 12, 10, 0, 0, 0, time.UTC)

	// Add events (hash chain should be built automatically)
	event1 := &AgentEvent{
		ID:        "1",
		Type:      EventTypeAgentStart,
		AgentID:   "a1",
		Role:      "impl",
		Timestamp: ts,
	}
	if err := store.Append(event1); err != nil {
		t.Fatalf("Failed to append: %v", err)
	}

	event2 := &AgentEvent{
		ID:        "2",
		Type:      EventTypeAgentComplete,
		AgentID:   "a1",
		Role:      "impl",
		Timestamp: ts.Add(time.Minute),
	}
	if err := store.Append(event2); err != nil {
		t.Fatalf("Failed to append: %v", err)
	}

	// Verify hash chain
	if err := store.VerifyHashChain(); err != nil {
		t.Errorf("Hash chain verification failed: %v", err)
	}
}

func TestStore_Replay(t *testing.T) {
	// AC3: Event replay works correctly
	tmpDir, err := os.MkdirTemp("", "coord-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	logPath := filepath.Join(tmpDir, "events.jsonl")
	store, err := NewStore(logPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	ts := time.Date(2026, 2, 12, 10, 0, 0, 0, time.UTC)

	events := []*AgentEvent{
		{ID: "1", Type: EventTypeAgentStart, AgentID: "a1", Role: "impl", TaskID: "t1", Timestamp: ts},
		{ID: "2", Type: EventTypeAgentAction, AgentID: "a1", Role: "impl", TaskID: "t1", Timestamp: ts.Add(5 * time.Minute)},
		{ID: "3", Type: EventTypeAgentComplete, AgentID: "a1", Role: "impl", TaskID: "t1", Timestamp: ts.Add(30 * time.Minute)},
	}

	for _, e := range events {
		if err := store.Append(e); err != nil {
			t.Fatalf("Failed to append: %v", err)
		}
	}

	// Close and reopen to simulate replay
	store.Close()

	store2, err := NewStore(logPath)
	if err != nil {
		t.Fatalf("Failed to reopen store: %v", err)
	}
	defer store2.Close()

	// Replay events
	replayed, err := store2.ReadAll()
	if err != nil {
		t.Fatalf("Replay failed: %v", err)
	}

	if len(replayed) != 3 {
		t.Errorf("Expected 3 replayed events, got %d", len(replayed))
	}

	// Verify order is preserved
	if replayed[0].ID != "1" || replayed[1].ID != "2" || replayed[2].ID != "3" {
		t.Error("Event order not preserved during replay")
	}
}
