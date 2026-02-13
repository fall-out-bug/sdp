package decision

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Logger handles logging decisions to JSONL file
type Logger struct {
	filePath string
	mu       sync.Mutex
	metrics  *MetricsRecorder
}

// NewLogger creates a new decision logger
func NewLogger(baseDir string) (*Logger, error) {
	log.Printf("[decision] NewLogger: base_dir=%s", baseDir)

	decisionsDir := filepath.Join(baseDir, "docs", "decisions")

	// Create directory if doesn't exist
	if err := os.MkdirAll(decisionsDir, 0755); err != nil {
		log.Printf("[decision] ERROR: failed to create decisions directory: %v", err)
		return nil, fmt.Errorf("failed to create decisions directory: %w", err)
	}

	filePath := filepath.Join(decisionsDir, "decisions.jsonl")
	log.Printf("[decision] NewLogger: success, file_path=%s", filePath)

	return &Logger{
		filePath: filePath,
		metrics:  &MetricsRecorder{},
	}, nil
}

// Log logs a decision to the JSONL file
func (l *Logger) Log(decision Decision) error {
	start := time.Now()

	// Check if rotation is needed (before lock to avoid deadlock)
	if err := l.rotate(); err != nil {
		log.Printf("[decision] WARNING: rotation failed: %v", err)
		// Continue anyway - log to current file
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// Set timestamp if not set
	if decision.Timestamp.IsZero() {
		decision.Timestamp = time.Now()
	}

	log.Printf("[decision] Log: question=%q type=%s feature_id=%s ws_id=%s",
		decision.Question, decision.Type, decision.FeatureID, decision.WorkstreamID)

	// Open file in append mode
	file, err := os.OpenFile(l.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("[decision] ERROR: failed to open file: path=%s error=%v", l.filePath, err)
		l.metrics.RecordLog(time.Since(start), false)
		return fmt.Errorf("failed to open decisions file: %w", err)
	}
	defer file.Close()

	// Marshal to JSON
	data, err := json.Marshal(decision)
	if err != nil {
		log.Printf("[decision] ERROR: failed to marshal: error=%v", err)
		l.metrics.RecordLog(time.Since(start), false)
		return fmt.Errorf("failed to marshal decision: %w", err)
	}

	// Write to file with newline
	if _, err := file.Write(append(data, '\n')); err != nil {
		log.Printf("[decision] ERROR: failed to write: path=%s error=%v", l.filePath, err)
		l.metrics.RecordLog(time.Since(start), false)
		return fmt.Errorf("failed to write decision: %w", err)
	}

	// Sync to disk for durability
	if err := file.Sync(); err != nil {
		log.Printf("[decision] ERROR: failed to sync: path=%s error=%v", l.filePath, err)
		l.metrics.RecordLog(time.Since(start), false)
		return fmt.Errorf("failed to sync decision: %w", err)
	}

	log.Printf("[decision] Log: success, decision=%q", decision.Decision)
	l.metrics.RecordLog(time.Since(start), true)
	return nil
}

// LogBatch logs multiple decisions at once
func (l *Logger) LogBatch(decisions []Decision) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	log.Printf("[decision] LogBatch: start, count=%d path=%s", len(decisions), l.filePath)

	// Open file in append mode
	file, err := os.OpenFile(l.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("[decision] ERROR: failed to open file: path=%s error=%v", l.filePath, err)
		return fmt.Errorf("failed to open decisions file: %w", err)
	}
	defer file.Close()

	successCount := 0
	for i, decision := range decisions {
		data, err := json.Marshal(decision)
		if err != nil {
			log.Printf("[decision] ERROR: failed to marshal at index %d: error=%v", i, err)
			return fmt.Errorf("failed to marshal decision at index %d: %w", i, err)
		}

		if _, err := file.Write(append(data, '\n')); err != nil {
			log.Printf("[decision] ERROR: failed to write at index %d: path=%s error=%v", i, l.filePath, err)
			return fmt.Errorf("failed to write decision at index %d: %w", i, err)
		}
		successCount++
	}

	// Sync to disk for durability
	if err := file.Sync(); err != nil {
		log.Printf("[decision] ERROR: failed to sync: path=%s error=%v", l.filePath, err)
		return fmt.Errorf("failed to sync decisions: %w", err)
	}

	log.Printf("[decision] LogBatch: success, count=%d", successCount)
	return nil
}
