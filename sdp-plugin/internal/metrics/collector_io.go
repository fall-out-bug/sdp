package metrics

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func (c *Collector) readEvents() ([]evidenceEvent, error) {
	f, err := os.Open(c.logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []evidenceEvent{}, nil
		}
		return nil, fmt.Errorf("open log: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			// Log close error but don't override read error
		}
	}()

	var events []evidenceEvent
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var ev evidenceEvent
		if err := json.Unmarshal(line, &ev); err != nil {
			continue // Skip invalid lines
		}
		// Convert data to map[string]interface{} for easier access
		if ev.Data == nil {
			ev.Data = make(map[string]interface{})
		}
		events = append(events, ev)
	}

	return events, sc.Err()
}

func (c *Collector) writeOutput(metrics *Metrics) error {
	// Ensure output directory exists
	dir := filepath.Dir(c.outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	data, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal metrics: %w", err)
	}

	return os.WriteFile(c.outputPath, data, 0644)
}

func (c *Collector) loadWatermark() map[string]bool {
	data, err := os.ReadFile(c.watermarkPath)
	if err != nil {
		return make(map[string]bool)
	}

	var ids []string
	if err := json.Unmarshal(data, &ids); err != nil {
		return make(map[string]bool)
	}

	result := make(map[string]bool, len(ids))
	for _, id := range ids {
		result[id] = true
	}
	return result
}

func (c *Collector) saveWatermark(processedIDs map[string]bool) error {
	ids := make([]string, 0, len(processedIDs))
	for id := range processedIDs {
		ids = append(ids, id)
	}

	data, err := json.Marshal(ids)
	if err != nil {
		return fmt.Errorf("marshal watermark: %w", err)
	}

	// Ensure watermark directory exists
	dir := filepath.Dir(c.watermarkPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create watermark dir: %w", err)
	}

	return os.WriteFile(c.watermarkPath, data, 0644)
}
