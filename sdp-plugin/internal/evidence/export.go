package evidence

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Exporter provides export and statistics for event logs.
type Exporter struct{}

// NewExporter creates a new exporter.
func NewExporter() *Exporter {
	return &Exporter{}
}

// ToCSV exports events as CSV with columns: timestamp, type, ws_id, model, action.
func (e *Exporter) ToCSV(events []Event) (string, error) {
	var b strings.Builder
	w := csv.NewWriter(&b)

	// Write header
	header := []string{"timestamp", "type", "ws_id", "model", "action"}
	if err := w.Write(header); err != nil {
		return "", fmt.Errorf("write header: %w", err)
	}

	// Write data rows
	for _, ev := range events {
		model, action := e.extractFields(ev)

		// Format timestamp for readability
		ts := ev.Timestamp
		if t, err := time.Parse(time.RFC3339, ev.Timestamp); err == nil {
			ts = t.Format("2006-01-02 15:04:05")
		}

		row := []string{ts, ev.Type, ev.WSID, model, action}
		if err := w.Write(row); err != nil {
			return "", fmt.Errorf("write row: %w", err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return "", fmt.Errorf("flush csv: %w", err)
	}

	return b.String(), nil
}

// extractFields extracts model and action from event data.
func (e *Exporter) extractFields(ev Event) (model, action string) {
	if ev.Data == nil {
		return "", ""
	}
	m, ok := ev.Data.(map[string]interface{})
	if !ok {
		return "", ""
	}

	// Extract model_id
	if v, ok := m["model_id"]; ok {
		if s, ok := v.(string); ok {
			model = s
		}
	}

	// Extract action (for generation events) or choice (for decisions)
	if v, ok := m["action"]; ok {
		if s, ok := v.(string); ok {
			action = s
		}
	} else if v, ok := m["choice"]; ok {
		if s, ok := v.(string); ok {
			action = s
		}
	}

	return model, action
}

// ToJSON exports events as a JSON array.
func (e *Exporter) ToJSON(events []Event) (string, error) {
	b, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal json: %w", err)
	}
	return string(b), nil
}

// Stats represents event statistics.
type Stats struct {
	Total             int
	CountByType       map[string]int
	ModelDistribution map[string]int
	DateDistribution  map[string]int // YYYY-MM-DD -> count
}

// Stats computes summary statistics for events.
func (e *Exporter) Stats(events []Event) Stats {
	stats := Stats{
		Total:             len(events),
		CountByType:       make(map[string]int),
		ModelDistribution: make(map[string]int),
		DateDistribution:  make(map[string]int),
	}

	for _, ev := range events {
		// Count by type
		stats.CountByType[ev.Type]++

		// Count by model
		if ev.Data != nil {
			m, ok := ev.Data.(map[string]interface{})
			if ok {
				if v, ok := m["model_id"]; ok {
					if s, ok := v.(string); ok {
						stats.ModelDistribution[s]++
					}
				}
			}
		}

		// Count by date
		if t, err := time.Parse(time.RFC3339, ev.Timestamp); err == nil {
			date := t.Format("2006-01-02")
			stats.DateDistribution[date]++
		}
	}

	return stats
}
