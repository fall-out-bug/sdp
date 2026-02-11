package evidence

import (
	"time"
)

// Browser provides filtering and pagination for event logs.
type Browser struct {
	events []Event
}

// NewBrowser creates a new browser for the given events.
func NewBrowser(events []Event) *Browser {
	return &Browser{events: events}
}

// Page returns the nth page (1-indexed) with pageSize items per page,
// along with the total number of events.
func (b *Browser) Page(n, pageSize int) ([]Event, int) {
	total := len(b.events)
	if total == 0 {
		return []Event{}, total
	}

	start := (n - 1) * pageSize
	if start >= total {
		return []Event{}, total
	}

	end := start + pageSize
	if end > total {
		end = total
	}

	return b.events[start:end], total
}

// FilterByType returns events matching the given type.
func (b *Browser) FilterByType(eventType string) []Event {
	if eventType == "" {
		return b.events
	}
	var out []Event
	for _, e := range b.events {
		if e.Type == eventType {
			out = append(out, e)
		}
	}
	return out
}

// FilterByModel returns generation events matching the given model ID.
func (b *Browser) FilterByModel(modelID string) []Event {
	if modelID == "" {
		return b.events
	}
	var out []Event
	for _, e := range b.events {
		if e.Data == nil {
			continue
		}
		m, ok := e.Data.(map[string]interface{})
		if !ok {
			continue
		}
		if v, ok := m["model_id"]; ok {
			if s, ok := v.(string); ok && s == modelID {
				out = append(out, e)
			}
		}
	}
	return out
}

// FilterBySince returns events on or after the given timestamp (RFC3339).
func (b *Browser) FilterBySince(since string) []Event {
	if since == "" {
		return b.events
	}
	sinceTime, err := time.Parse(time.RFC3339, since)
	if err != nil {
		// If parsing fails, return all events
		return b.events
	}
	var out []Event
	for _, e := range b.events {
		ts, err := time.Parse(time.RFC3339, e.Timestamp)
		if err != nil {
			continue
		}
		if !ts.Before(sinceTime) {
			out = append(out, e)
		}
	}
	return out
}

// FilterByWS returns events matching the given workstream ID.
func (b *Browser) FilterByWS(wsID string) []Event {
	if wsID == "" {
		return b.events
	}
	var out []Event
	for _, e := range b.events {
		if e.WSID == wsID {
			out = append(out, e)
		}
	}
	return out
}
