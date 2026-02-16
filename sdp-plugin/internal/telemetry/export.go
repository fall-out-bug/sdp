package telemetry

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

// ExportJSON exports telemetry data to JSON format
func (c *Collector) ExportJSON(exportPath string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Read all events from file
	events, err := c.readEvents()
	if err != nil {
		return fmt.Errorf("failed to read events: %w", err)
	}

	// Marshal to JSON array
	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal events: %w", err)
	}

	// Write to export file (restricted permissions for telemetry data)
	if err := os.WriteFile(exportPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}

	return nil
}

// ExportCSV exports telemetry data to CSV format
func (c *Collector) ExportCSV(exportPath string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Read all events
	events, err := c.readEvents()
	if err != nil {
		return fmt.Errorf("failed to read events: %w", err)
	}

	// Create CSV file (restricted permissions for telemetry data)
	file, err := os.OpenFile(exportPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to create export file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			// Log but don't fail if close fails after successful write
			fmt.Fprintf(os.Stderr, "warning: failed to close export file: %v\n", cerr)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Type", "Timestamp", "Data"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write rows
	for _, event := range events {
		// Marshal data to JSON string for CSV
		dataJSON, err := json.Marshal(event.Data)
		if err != nil {
			return fmt.Errorf("failed to marshal event data: %w", err)
		}

		row := []string{
			string(event.Type),
			event.Timestamp.Format("2006-01-02T15:04:05-07:00"),
			string(dataJSON),
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return nil
}
