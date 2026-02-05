package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Analyzer analyzes telemetry data to generate insights
type Analyzer struct {
	filePath string
}

// NewAnalyzer creates a new telemetry analyzer
func NewAnalyzer(filePath string) (*Analyzer, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}

	return &Analyzer{
		filePath: filePath,
	}, nil
}

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

// CalculateSuccessRate calculates success rate by command
func (a *Analyzer) CalculateSuccessRate() (map[string]float64, error) {
	events, err := a.readEvents()
	if err != nil {
		return nil, fmt.Errorf("failed to read events: %w", err)
	}

	// Track command results: command -> [success, total]
	commandResults := make(map[string][2]int)

	for _, event := range events {
		if event.Type != EventTypeCommandComplete {
			continue
		}

		command, ok := event.Data["command"].(string)
		if !ok {
			continue
		}

		success, ok := event.Data["success"].(bool)
		if !ok {
			continue
		}

		stats := commandResults[command]
		if success {
			stats[0]++ // success count
		}
		stats[1]++ // total count
		commandResults[command] = stats
	}

	// Calculate rates
	rates := make(map[string]float64)
	for command, stats := range commandResults {
		total := stats[1]
		if total > 0 {
			rates[command] = float64(stats[0]) / float64(total)
		}
	}

	return rates, nil
}

// CalculateAverageDuration calculates average duration by command
func (a *Analyzer) CalculateAverageDuration() (map[string]int, error) {
	events, err := a.readEvents()
	if err != nil {
		return nil, fmt.Errorf("failed to read events: %w", err)
	}

	// Track command durations: command -> [total_duration, count]
	commandDurations := make(map[string][2]int)

	for _, event := range events {
		if event.Type != EventTypeCommandComplete {
			continue
		}

		command, ok := event.Data["command"].(string)
		if !ok {
			continue
		}

		durationFloat, ok := event.Data["duration"].(float64)
		if !ok {
			continue
		}

		duration := int(durationFloat)

		stats := commandDurations[command]
		stats[0] += duration
		stats[1]++
		commandDurations[command] = stats
	}

	// Calculate averages
	averages := make(map[string]int)
	for command, stats := range commandDurations {
		count := stats[1]
		if count > 0 {
			averages[command] = stats[0] / count
		}
	}

	return averages, nil
}

// TopErrorCategories returns the top N error categories
func (a *Analyzer) TopErrorCategories(n int) ([]ErrorCategory, error) {
	events, err := a.readEvents()
	if err != nil {
		return nil, fmt.Errorf("failed to read events: %w", err)
	}

	// Count errors by message
	errorCounts := make(map[string]int)

	for _, event := range events {
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

	// Convert to slice and sort by count
	errors := make([]ErrorCategory, 0, len(errorCounts))
	for msg, count := range errorCounts {
		errors = append(errors, ErrorCategory{
			Message: msg,
			Count:   count,
		})
	}

	// Sort by count (descending)
	for i := 0; i < len(errors); i++ {
		for j := i + 1; j < len(errors); j++ {
			if errors[j].Count > errors[i].Count {
				errors[i], errors[j] = errors[j], errors[i]
			}
		}
	}

	// Return top N
	if n > len(errors) {
		n = len(errors)
	}
	return errors[:n], nil
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

	report.TopErrors = errors
	if len(errors) > 5 {
		report.TopErrors = errors[:5]
	}

	return report, nil
}

// commandDataInternal tracks internal calculation data
type commandDataInternal struct {
	totalRuns     int
	successCount  int
	totalDuration int
	durationCount int
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
