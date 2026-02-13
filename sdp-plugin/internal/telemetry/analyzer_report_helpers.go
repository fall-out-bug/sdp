package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
)

// calculateTopErrors extracts top errors from filtered events
func (a *Analyzer) calculateTopErrors(filteredEvents []Event) []ErrorCategory {
	errorCounts := make(map[string]int)
	for _, event := range filteredEvents {
		if event.Type != EventTypeCommandComplete {
			continue
		}

		success, ok := event.Data["success"].(bool)
		if ok && success {
			continue
		}

		errorMsg, ok := event.Data["error"].(string)
		if !ok || errorMsg == "" {
			errorMsg = "unknown error"
		}

		errorCounts[errorMsg]++
	}

	// Sort errors by count
	errors := make([]ErrorCategory, 0, len(errorCounts))
	for msg, count := range errorCounts {
		errors = append(errors, ErrorCategory{
			Message: msg,
			Count:   count,
		})
	}

	for i := 0; i < len(errors) && i < 5; i++ {
		for j := i + 1; j < len(errors); j++ {
			if errors[j].Count > errors[i].Count {
				errors[i], errors[j] = errors[j], errors[i]
			}
		}
	}

	if len(errors) > 5 {
		return errors[:5]
	}
	return errors
}

// readEvents reads all events from the telemetry file
func (a *Analyzer) readEvents() ([]Event, error) {
	// If file doesn't exist, return empty slice
	if _, err := os.Stat(a.filePath); os.IsNotExist(err) {
		return []Event{}, nil
	}

	// Read file
	data, err := os.ReadFile(a.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read telemetry file: %w", err)
	}

	// If file is empty, return empty slice
	if len(data) == 0 {
		return []Event{}, nil
	}

	// Parse JSONL (array format for test compatibility)
	var events []Event
	if err := json.Unmarshal(data, &events); err != nil {
		// If array parsing fails, try JSONL format
		lines := splitLines(data)
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}

			var event Event
			if err := json.Unmarshal(line, &event); err != nil {
				return nil, fmt.Errorf("failed to unmarshal event: %w", err)
			}

			events = append(events, event)
		}
	}

	return events, nil
}
