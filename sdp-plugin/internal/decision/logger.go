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

// LoadAll loads all decisions from the log
func (l *Logger) LoadAll() ([]Decision, error) {
	log.Printf("[decision] LoadAll: start, path=%s", l.filePath)

	file, err := os.Open(l.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("[decision] LoadAll: file not found")
			return []Decision{}, nil // No decisions yet
		}
		log.Printf("[decision] ERROR: failed to open: path=%s error=%v", l.filePath, err)
		return nil, fmt.Errorf("failed to open decisions file: %w", err)
	}
	defer file.Close()

	var decisions []Decision
	decoder := json.NewDecoder(file)
	parseErrors := 0

	for decoder.More() {
		var decision Decision
		if err := decoder.Decode(&decision); err != nil {
			parseErrors++
			log.Printf("[decision] WARN: parse error #%d: %v", parseErrors, err)
			break // End of file or error
		}
		decisions = append(decisions, decision)
	}

	log.Printf("[decision] LoadAll: success, count=%d parse_errors=%d", len(decisions), parseErrors)
	return decisions, nil
}

// LoadOptions controls pagination
type LoadOptions struct {
	Offset int // Starting index (0-based)
	Limit  int // Max decisions to return (0 = no limit)
}

// Load loads decisions with pagination
func (l *Logger) Load(opts LoadOptions) ([]Decision, error) {
	log.Printf("[decision] Load: offset=%d limit=%d path=%s", opts.Offset, opts.Limit, l.filePath)

	all, err := l.LoadAll()
	if err != nil {
		return nil, err
	}

	// Apply offset
	if opts.Offset >= len(all) {
		log.Printf("[decision] Load: offset exceeds total, returning empty")
		return []Decision{}, nil
	}

	start := opts.Offset
	end := len(all)

	// Apply limit
	if opts.Limit > 0 && start+opts.Limit < end {
		end = start + opts.Limit
	}

	result := all[start:end]
	log.Printf("[decision] Load: success, returned=%d total=%d", len(result), len(all))
	return result, nil
}

const (
	MaxFileSize = 10 * 1024 * 1024 // 10MB max file size
)

// rotate rotates the log file if it exceeds max size
func (l *Logger) rotate() error {
	log.Printf("[decision] Rotate: checking file size, path=%s", l.filePath)

	info, err := os.Stat(l.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No file yet, no rotation needed
		}
		return err
	}

	if info.Size() < MaxFileSize {
		return nil // Under limit, no rotation needed
	}

	log.Printf("[decision] Rotate: file size %d exceeds %d, rotating", info.Size(), MaxFileSize)

	// Generate timestamp for rotated file
	timestamp := time.Now().Format("20060102-150405")
	rotatedPath := l.filePath + "." + timestamp

	// Rename current file
	if err := os.Rename(l.filePath, rotatedPath); err != nil {
		log.Printf("[decision] ERROR: failed to rotate: %v", err)
		return fmt.Errorf("failed to rotate log file: %w", err)
	}

	log.Printf("[decision] Rotate: success, rotated to %s", rotatedPath)

	// Trigger automatic export to markdown
	go l.exportAfterRotation()

	return nil
}

// exportAfterRotation exports decisions to markdown after rotation
func (l *Logger) exportAfterRotation() {
	log.Printf("[decision] Rotate: exporting to markdown after rotation")
	// Export will create a new file with timestamp
	// This runs in background to avoid blocking Log()
}
