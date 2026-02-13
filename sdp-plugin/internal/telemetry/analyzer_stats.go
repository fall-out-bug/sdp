package telemetry

import "fmt"

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
