package memory

import (
	"time"
)

// CompactionPolicy configures memory compaction behavior (AC1)
type CompactionPolicy struct {
	MaxDBSizeMB        int  `json:"max_db_size_mb"`
	EventRetentionDays int  `json:"event_retention_days"`
	CompactionRatio    int  `json:"compaction_ratio"`
	ArchiveAfterDays   int  `json:"archive_after_days"`
	AutoCompact        bool `json:"auto_compact"`
}

// DefaultCompactionPolicy returns sensible defaults
func DefaultCompactionPolicy() CompactionPolicy {
	return CompactionPolicy{
		MaxDBSizeMB:        100,
		EventRetentionDays: 30,
		CompactionRatio:    10,
		ArchiveAfterDays:   90,
		AutoCompact:        true,
	}
}

// CompactableEvent represents an event that can be compacted
type CompactableEvent struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	Data      string    `json:"data"`
}

// EventSummary represents a summarized group of events (AC2)
type EventSummary struct {
	ID         string    `json:"id"`
	WSID       string    `json:"ws_id"`
	Summary    string    `json:"summary"`
	EventCount int       `json:"event_count"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
}

// Compactor handles memory compaction (AC5)
type Compactor struct {
	policy CompactionPolicy
}

// NewCompactor creates a new compactor
func NewCompactor(policy CompactionPolicy) *Compactor {
	return &Compactor{policy: policy}
}

// NeedsCompaction checks if compaction is needed based on size or age
func (c *Compactor) NeedsCompaction(dbSizeBytes int64, oldEventCount int) bool {
	// Check size threshold
	if c.policy.MaxDBSizeMB > 0 {
		maxBytes := int64(c.policy.MaxDBSizeMB) * 1024 * 1024
		if dbSizeBytes > maxBytes {
			return true
		}
	}

	// Check event age threshold
	if c.policy.EventRetentionDays > 0 && oldEventCount > 0 {
		return true
	}

	return false
}

// CompactEvents groups events into summaries (AC2)
func (c *Compactor) CompactEvents(events []CompactableEvent) []EventSummary {
	if len(events) == 0 {
		return nil
	}

	ratio := c.policy.CompactionRatio
	if ratio <= 0 {
		ratio = 10
	}

	var summaries []EventSummary

	// Group events into batches of 'ratio' size
	for i := 0; i < len(events); i += ratio {
		end := i + ratio
		if end > len(events) {
			end = len(events)
		}

		batch := events[i:end]
		summary := c.Summarize(batch)
		summaries = append(summaries, summary)
	}

	return summaries
}

// Summarize creates a summary from a batch of events (AC2)
func (c *Compactor) Summarize(events []CompactableEvent) EventSummary {
	if len(events) == 0 {
		return EventSummary{}
	}

	// Find earliest and latest timestamps
	startTime := events[0].Timestamp
	endTime := events[len(events)-1].Timestamp

	// Build a simple summary (in production, this would use LLM)
	summary := "Compacted " + string(rune(len(events))) + " events"

	return EventSummary{
		ID:         generateSummaryID(),
		Summary:    summary,
		EventCount: len(events),
		StartTime:  startTime,
		EndTime:    endTime,
	}
}

// generateSummaryID generates a unique summary ID
func generateSummaryID() string {
	return "sum-" + time.Now().Format("20060102150405")
}
