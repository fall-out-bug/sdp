package telemetry

import (
	"archive/tar"
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"
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
	defer file.Close()

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

// UploadResult represents the result of packaging telemetry for upload
type UploadResult struct {
	EventCount int    `json:"event_count"`
	Format     string `json:"format"`
	Size       int64  `json:"size"`
	File       string `json:"file"`
}

// PackForUpload packages telemetry data for upload
// Supported formats: "json" or "archive" (tar.gz)
func PackForUpload(telemetryFile, outputPath, format string) (*UploadResult, error) {
	// Read all events from telemetry file
	events, err := readEventsFromFile(telemetryFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read events: %w", err)
	}

	// Package based on format
	switch format {
	case "json":
		return packAsJSON(events, outputPath)
	case "archive", "tar.gz", "tgz":
		return packAsArchive(events, telemetryFile, outputPath)
	default:
		return nil, fmt.Errorf("unsupported format: %s (use 'json' or 'archive')", format)
	}
}

// packAsJSON packages events as structured JSON
func packAsJSON(events []Event, outputPath string) (*UploadResult, error) {
	// Create upload structure
	uploadData := struct {
		Metadata map[string]interface{} `json:"metadata"`
		Events   []Event                `json:"events"`
	}{
		Metadata: map[string]interface{}{
			"version":      "1.0",
			"generated_at": time.Now().Format(time.RFC3339),
			"event_count":  len(events),
			"format":       "json",
		},
		Events: events,
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(uploadData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal upload data: %w", err)
	}

	// Write to file
	if err := os.WriteFile(outputPath, data, 0600); err != nil {
		return nil, fmt.Errorf("failed to write upload file: %w", err)
	}

	return &UploadResult{
		EventCount: len(events),
		Format:     "json",
		Size:       int64(len(data)),
		File:       outputPath,
	}, nil
}

// packAsArchive packages events as tar.gz archive
func packAsArchive(events []Event, telemetryFile, outputPath string) (*UploadResult, error) {
	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create archive: %w", err)
	}
	defer file.Close()

	// Create gzip writer
	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	// Add telemetry.jsonl to archive
	data, err := os.ReadFile(telemetryFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read telemetry file: %w", err)
	}

	header := &tar.Header{
		Name:    "telemetry.jsonl",
		Mode:    0600,
		Size:    int64(len(data)),
		ModTime: time.Now(),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		return nil, fmt.Errorf("failed to write tar header: %w", err)
	}

	if _, err := tarWriter.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write tar content: %w", err)
	}

	// Get file size
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat archive: %w", err)
	}

	return &UploadResult{
		EventCount: len(events),
		Format:     "archive",
		Size:       info.Size(),
		File:       outputPath,
	}, nil
}

// readEventsFromFile reads all events from a telemetry file
func readEventsFromFile(filePath string) ([]Event, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("telemetry file does not exist: %s", filePath)
	}

	// Read file
	data, err := os.ReadFile(filePath)
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
			// Skip invalid lines
			continue
		}

		events = append(events, event)
	}

	return events, nil
}
