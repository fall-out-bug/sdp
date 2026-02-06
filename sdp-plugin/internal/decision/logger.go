package decision

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Logger handles logging decisions to JSONL file
type Logger struct {
	filePath string
}

// NewLogger creates a new decision logger
func NewLogger(baseDir string) (*Logger, error) {
	decisionsDir := filepath.Join(baseDir, "docs", "decisions")

	// Create directory if doesn't exist
	if err := os.MkdirAll(decisionsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create decisions directory: %w", err)
	}

	filePath := filepath.Join(decisionsDir, "decisions.jsonl")

	return &Logger{
		filePath: filePath,
	}, nil
}

// Log logs a decision to the JSONL file
func (l *Logger) Log(decision Decision) error {
	// Set timestamp if not set
	if decision.Timestamp.IsZero() {
		decision.Timestamp = time.Now()
	}

	// Open file in append mode
	file, err := os.OpenFile(l.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open decisions file: %w", err)
	}
	defer file.Close()

	// Marshal to JSON
	data, err := json.Marshal(decision)
	if err != nil {
		return fmt.Errorf("failed to marshal decision: %w", err)
	}

	// Write to file with newline
	if _, err := file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write decision: %w", err)
	}

	return nil
}

// LogBatch logs multiple decisions at once
func (l *Logger) LogBatch(decisions []Decision) error {
	// Open file in append mode
	file, err := os.OpenFile(l.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open decisions file: %w", err)
	}
	defer file.Close()

	for _, decision := range decisions {
		data, err := json.Marshal(decision)
		if err != nil {
			return fmt.Errorf("failed to marshal decision: %w", err)
		}

		if _, err := file.Write(append(data, '\n')); err != nil {
			return fmt.Errorf("failed to write decision: %w", err)
		}
	}

	return nil
}

// LoadAll loads all decisions from the log
func (l *Logger) LoadAll() ([]Decision, error) {
	file, err := os.Open(l.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Decision{}, nil // No decisions yet
		}
		return nil, fmt.Errorf("failed to open decisions file: %w", err)
	}
	defer file.Close()

	var decisions []Decision
	decoder := json.NewDecoder(file)

	for decoder.More() {
		var decision Decision
		if err := decoder.Decode(&decision); err != nil {
			break // End of file or error
		}
		decisions = append(decisions, decision)
	}

	return decisions, nil
}
