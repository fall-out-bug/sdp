package decision

import (
	"fmt"
	"log"
	"os"
	"time"
)

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
