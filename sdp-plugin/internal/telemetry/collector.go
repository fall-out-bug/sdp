package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Collector manages telemetry collection
type Collector struct {
	filePath   string
	enabled    bool
	mu         sync.Mutex
	eventCount int
}

// NewCollector creates a new telemetry collector
func NewCollector(filePath string, enabled bool) (*Collector, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create telemetry directory: %w", err)
	}

	return &Collector{
		filePath:   filePath,
		enabled:    enabled,
		eventCount: 0,
	}, nil
}

// Record records a telemetry event
func (c *Collector) Record(event Event) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If telemetry is disabled, do nothing
	if !c.enabled {
		return nil
	}

	// Validate event
	if !event.Type.IsValid() {
		return fmt.Errorf("invalid event type: %s", event.Type)
	}

	// Marshal event to JSON
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Append to file
	file, err := os.OpenFile(c.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open telemetry file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write event: %w", err)
	}

	c.eventCount++
	return nil
}

// Status returns the current status
func (c *Collector) Status() Status {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Count actual events in file for accurate status
	actualCount := c.eventCount
	if c.enabled && c.eventCount == 0 {
		if events, err := c.readEvents(); err == nil {
			actualCount = len(events)
		}
	}

	return Status{
		Enabled:    c.enabled,
		EventCount: actualCount,
		FilePath:   c.filePath,
	}
}

// Disable disables telemetry collection
func (c *Collector) Disable() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.enabled = false
}

// Enable enables telemetry collection
func (c *Collector) Enable() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.enabled = true
}

// readEvents reads all events from the telemetry file
func (c *Collector) readEvents() ([]Event, error) {
	// If file doesn't exist, return empty slice
	if _, err := os.Stat(c.filePath); os.IsNotExist(err) {
		return []Event{}, nil
	}

	// Read file
	data, err := os.ReadFile(c.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read telemetry file: %w", err)
	}

	// Parse JSONL
	var events []Event
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

	return events, nil
}

// splitLines splits byte array into lines
func splitLines(data []byte) [][]byte {
	return splitLinesHelper(data, '\n')
}

// splitLinesHelper splits byte array by delimiter
func splitLinesHelper(data []byte, delimiter byte) [][]byte {
	var lines [][]byte
	start := 0

	for i := 0; i < len(data); i++ {
		if data[i] == delimiter {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}

	// Add last line if it doesn't end with delimiter
	if start < len(data) {
		lines = append(lines, data[start:])
	}

	return lines
}
