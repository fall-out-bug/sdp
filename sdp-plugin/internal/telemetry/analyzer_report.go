package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// CommandStats represents statistics for a single command
type CommandStats struct {
	Command     string  `json:"command"`
	TotalRuns   int     `json:"total_runs"`
	SuccessRate float64 `json:"success_rate"`
	AvgDuration int     `json:"avg_duration_ms"`
}

// ErrorCategory represents an error category with count
type ErrorCategory struct {
	Message string `json:"message"`
	Count   int    `json:"count"`
}

// Report represents a telemetry analysis report
type Report struct {
	TotalEvents  int                     `json:"total_events"`
	DateRange    *DateRange              `json:"date_range,omitempty"`
	CommandStats map[string]CommandStats `json:"command_stats"`
	TopErrors    []ErrorCategory         `json:"top_errors"`
}

// DateRange represents a time range for filtering
type DateRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// commandDataInternal tracks internal calculation data
type commandDataInternal struct {
	totalRuns     int
	successCount  int
	totalDuration int
	durationCount int
}

// GenerateReport generates a comprehensive telemetry report
func (a *Analyzer) GenerateReport(startTime, endTime *time.Time) (*Report, error) {
	events, err := a.readEvents()
	if err != nil {
		return nil, fmt.Errorf("failed to read events: %w", err)
	}

	// Filter by date range if provided
	filteredEvents := events
	if startTime != nil || endTime != nil {
		filteredEvents = make([]Event, 0)
		for _, event := range events {
			// Skip if before start time (inclusive of start time)
			if startTime != nil && event.Timestamp.Before(*startTime) {
				continue
			}
			// Skip if after end time (inclusive of end time)
			if endTime != nil && event.Timestamp.After(*endTime) {
				continue
			}
			filteredEvents = append(filteredEvents, event)
		}
	}

	// Build report
	report := &Report{
		TotalEvents:  len(filteredEvents),
		CommandStats: make(map[string]CommandStats),
		TopErrors:    []ErrorCategory{},
	}

	if startTime != nil && endTime != nil {
		report.DateRange = &DateRange{
			Start: *startTime,
			End:   *endTime,
		}
	}

	// Calculate command stats
	commandData := make(map[string]*commandDataInternal)
	for _, event := range filteredEvents {
		if event.Type != EventTypeCommandComplete {
			continue
		}

		command, ok := event.Data["command"].(string)
		if !ok {
			continue
		}

		if commandData[command] == nil {
			commandData[command] = &commandDataInternal{}
		}

		success, ok := event.Data["success"].(bool)
		if ok {
			commandData[command].totalRuns++
			if success {
				commandData[command].successCount++
			}
		}

		durationFloat, ok := event.Data["duration"].(float64)
		if ok {
			commandData[command].totalDuration += int(durationFloat)
			commandData[command].durationCount++
		}
	}

	// Convert to output format
	for command, data := range commandData {
		successRate := 0.0
		if data.totalRuns > 0 {
			successRate = float64(data.successCount) / float64(data.totalRuns)
		}

		avgDuration := 0
		if data.durationCount > 0 {
			avgDuration = data.totalDuration / data.durationCount
		}

		report.CommandStats[command] = CommandStats{
			Command:     command,
			TotalRuns:   data.totalRuns,
			SuccessRate: successRate,
			AvgDuration: avgDuration,
		}
	}

	// Calculate top errors
	report.TopErrors = a.calculateTopErrors(filteredEvents)

	return report, nil
}

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
