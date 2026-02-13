package memory

import (
	"time"

	"github.com/fall-out-bug/sdp/internal/evidence"
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

// CompactableEvent is an alias for evidence.Event for backward compatibility
// Deprecated: Use evidence.Event directly
type CompactableEvent = evidence.Event

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

	// Parse timestamps and find earliest/latest
	startTime := parseEventTimestamp(events[0].Timestamp)
	endTime := parseEventTimestamp(events[len(events)-1].Timestamp)

	// Group by WSID
	wsidMap := make(map[string]int)
	for _, ev := range events {
		wsidMap[ev.WSID]++
	}

	// Build summary with WSID info
	summary := "Compacted " + intToStr(len(events)) + " events across " + intToStr(len(wsidMap)) + " workstreams"

	return EventSummary{
		ID:         generateSummaryID(),
		WSID:       events[0].WSID, // Primary WSID
		Summary:    summary,
		EventCount: len(events),
		StartTime:  startTime,
		EndTime:    endTime,
	}
}

// parseEventTimestamp parses an evidence event timestamp string
func parseEventTimestamp(s string) time.Time {
	if s == "" {
		return time.Now()
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Now()
	}
	return t
}

// generateSummaryID generates a unique summary ID
func generateSummaryID() string {
	return "sum-" + time.Now().Format("20060102150405")
}
