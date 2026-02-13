package telemetry

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

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
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close archive: %v\n", cerr)
		}
	}()

	// Create gzip writer
	gzWriter := gzip.NewWriter(file)
	defer func() {
		if cerr := gzWriter.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close gzip writer: %v\n", cerr)
		}
	}()

	// Create tar writer
	tarWriter := tar.NewWriter(gzWriter)
	defer func() {
		if cerr := tarWriter.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to close tar writer: %v\n", cerr)
		}
	}()

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
	lines := splitLines(data)
	events := make([]Event, 0, len(lines))
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
