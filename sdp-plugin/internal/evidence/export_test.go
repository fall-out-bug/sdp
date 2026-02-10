package evidence

import (
	"encoding/csv"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestExport_CSV(t *testing.T) {
	now := time.Now()
	events := []Event{
		{
			ID:        "1",
			Type:      "generation",
			Timestamp: now.Format(time.RFC3339),
			WSID:      "00-054-03",
			Data: map[string]interface{}{
				"model_id": "claude-sonnet-4",
				"action":   "implemented feature",
			},
		},
		{
			ID:        "2",
			Type:      "decision",
			Timestamp: now.Add(time.Minute).Format(time.RFC3339),
			WSID:      "00-054-03",
			Data: map[string]interface{}{
				"choice":    "use Postgres",
				"rationale": "better for complex queries",
			},
		},
	}

	exporter := NewExporter()

	// AC6: Export as CSV
	csvData, err := exporter.ToCSV(events)
	if err != nil {
		t.Fatalf("ToCSV failed: %v", err)
	}

	// Parse CSV to verify structure
	r := csv.NewReader(strings.NewReader(csvData))
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("Failed to parse CSV: %v", err)
	}

	// Should have header + 2 data rows
	if len(records) != 3 {
		t.Errorf("Expected 3 rows (header + 2 data), got %d", len(records))
	}

	// Check header
	header := records[0]
	expectedHeader := []string{"timestamp", "type", "ws_id", "model", "action"}
	if len(header) != len(expectedHeader) {
		t.Errorf("Expected %d columns, got %d", len(expectedHeader), len(header))
	}

	// Check data rows
	if records[1][0] == "" {
		t.Errorf("Timestamp should not be empty")
	}
	if records[1][1] != "generation" {
		t.Errorf("Expected type 'generation', got '%s'", records[1][1])
	}
	if records[1][2] != "00-054-03" {
		t.Errorf("Expected ws_id '00-054-03', got '%s'", records[1][2])
	}
	if records[1][3] != "claude-sonnet-4" {
		t.Errorf("Expected model 'claude-sonnet-4', got '%s'", records[1][3])
	}
}

func TestExport_JSON(t *testing.T) {
	now := time.Now()
	events := []Event{
		{
			ID:        "1",
			Type:      "generation",
			Timestamp: now.Format(time.RFC3339),
			WSID:      "00-054-03",
			Data: map[string]interface{}{
				"model_id": "claude-sonnet-4",
			},
		},
	}

	exporter := NewExporter()

	// AC7: Export as JSON
	jsonStr, err := exporter.ToJSON(events)
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	// Verify it's valid JSON array
	var parsed []Event
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(parsed) != 1 {
		t.Errorf("Expected 1 event, got %d", len(parsed))
	}
	if parsed[0].Type != "generation" {
		t.Errorf("Expected type 'generation', got '%s'", parsed[0].Type)
	}
	if parsed[0].WSID != "00-054-03" {
		t.Errorf("Expected ws_id '00-054-03', got '%s'", parsed[0].WSID)
	}
}

func TestStats_Summary(t *testing.T) {
	now := time.Now()
	events := []Event{
		{Type: "generation", Timestamp: now.Format(time.RFC3339), Data: map[string]interface{}{"model_id": "claude-sonnet-4"}},
		{Type: "generation", Timestamp: now.Add(time.Minute).Format(time.RFC3339), Data: map[string]interface{}{"model_id": "claude-opus-4"}},
		{Type: "verification", Timestamp: now.Add(2 * time.Minute).Format(time.RFC3339)},
		{Type: "decision", Timestamp: now.Add(3 * time.Minute).Format(time.RFC3339)},
		{Type: "generation", Timestamp: now.Add(4 * time.Minute).Format(time.RFC3339), Data: map[string]interface{}{"model_id": "claude-sonnet-4"}},
	}

	exporter := NewExporter()

	// AC8: Stats summary
	stats := exporter.Stats(events)

	// Check event counts by type
	if stats.CountByType["generation"] != 3 {
		t.Errorf("Expected 3 generation events, got %d", stats.CountByType["generation"])
	}
	if stats.CountByType["verification"] != 1 {
		t.Errorf("Expected 1 verification event, got %d", stats.CountByType["verification"])
	}
	if stats.CountByType["decision"] != 1 {
		t.Errorf("Expected 1 decision event, got %d", stats.CountByType["decision"])
	}

	// Check model distribution
	if stats.ModelDistribution["claude-sonnet-4"] != 2 {
		t.Errorf("Expected 2 claude-sonnet-4 events, got %d", stats.ModelDistribution["claude-sonnet-4"])
	}
	if stats.ModelDistribution["claude-opus-4"] != 1 {
		t.Errorf("Expected 1 claude-opus-4 event, got %d", stats.ModelDistribution["claude-opus-4"])
	}

	// Check total
	if stats.Total != 5 {
		t.Errorf("Expected total 5 events, got %d", stats.Total)
	}
}
