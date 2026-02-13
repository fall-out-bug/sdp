package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// getValidTypes returns list of valid failure types.
func getValidTypes() []string {
	return []string{
		"wrong_logic",
		"missing_edge_case",
		"hallucinated_api",
		"type_error",
		"test_passing_but_wrong",
		"compilation_error",
		"import_error",
	}
}

// evidenceEvent represents a simplified evidence event for classification.
type evidenceEvent struct {
	ID   string                 `json:"id"`
	Type string                 `json:"type"`
	WSID string                 `json:"ws_id"`
	Data map[string]interface{} `json:"data"`
}

// readEventsJSONL reads events from JSONL file.
func readEventsJSONL(path string) ([]evidenceEvent, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []evidenceEvent{}, nil
		}
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	events := make([]evidenceEvent, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var ev evidenceEvent
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			continue // Skip invalid lines
		}
		if ev.Data == nil {
			ev.Data = make(map[string]interface{})
		}
		events = append(events, ev)
	}

	return events, nil
}

// initMetricsDir ensures metrics directory exists
func initMetricsDir() {
	metricsDir := filepath.Join(".sdp", "metrics")
	if err := os.MkdirAll(metricsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "warning: failed to create metrics directory: %v\n", err)
	}
}
