package memory

import (
	"testing"
	"time"
)

func TestCompactor_NeedsCompaction(t *testing.T) {
	compactor := NewCompactor(CompactionPolicy{
		MaxDBSizeMB: 10,
	})

	// Under threshold
	if compactor.NeedsCompaction(5*1024*1024, 0) {
		t.Error("Should not need compaction under threshold")
	}

	// Over threshold
	if !compactor.NeedsCompaction(15*1024*1024, 0) {
		t.Error("Should need compaction over threshold")
	}
}

func TestCompactor_NeedsCompaction_ByAge(t *testing.T) {
	compactor := NewCompactor(CompactionPolicy{
		EventRetentionDays: 30,
	})

	// No old events
	if compactor.NeedsCompaction(0, 0) {
		t.Error("Should not need compaction with no old events")
	}

	// Has old events
	if !compactor.NeedsCompaction(0, 100) {
		t.Error("Should need compaction with old events")
	}
}

func TestCompactor_CompactEvents(t *testing.T) {
	compactor := NewCompactor(CompactionPolicy{
		CompactionRatio: 10,
	})

	events := []CompactableEvent{}
	for i := 0; i < 25; i++ {
		events = append(events, CompactableEvent{
			ID:        string(rune(i)),
			Timestamp: time.Now().Add(-time.Duration(i) * time.Hour),
			Type:      "test",
			Data:      "test data",
		})
	}

	summaries := compactor.CompactEvents(events)

	// Should have 3 summaries (25 events / 10 ratio = 2.5, rounded up = 3)
	if len(summaries) != 3 {
		t.Errorf("Expected 3 summaries, got %d", len(summaries))
	}

	// Each summary should cover ~10 events
	for i, s := range summaries {
		if s.EventCount <= 0 {
			t.Errorf("Summary %d should have events", i)
		}
	}
}

func TestCompactor_Summarize(t *testing.T) {
	compactor := NewCompactor(CompactionPolicy{})

	now := time.Now()
	events := []CompactableEvent{
		{ID: "1", Type: "agent_start", Data: "started workstream", Timestamp: now.Add(-2 * time.Hour)},
		{ID: "2", Type: "agent_action", Data: "wrote tests", Timestamp: now.Add(-1 * time.Hour)},
		{ID: "3", Type: "agent_complete", Data: "finished workstream", Timestamp: now},
	}

	summary := compactor.Summarize(events)

	if summary.EventCount != 3 {
		t.Errorf("Expected 3 events, got %d", summary.EventCount)
	}
	if summary.StartTime.IsZero() || summary.EndTime.IsZero() {
		t.Error("Start and end times should be set")
	}
}

func TestCompactionPolicy_Default(t *testing.T) {
	policy := DefaultCompactionPolicy()

	if policy.MaxDBSizeMB != 100 {
		t.Errorf("Expected max size 100MB, got %d", policy.MaxDBSizeMB)
	}
	if policy.EventRetentionDays != 30 {
		t.Errorf("Expected retention 30 days, got %d", policy.EventRetentionDays)
	}
	if policy.CompactionRatio != 10 {
		t.Errorf("Expected ratio 10, got %d", policy.CompactionRatio)
	}
}
