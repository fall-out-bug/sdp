// Package safetylog provides structured logging for git safety operations.
package safetylog

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// Level represents log severity.
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var (
	logger   *log.Logger
	once     sync.Once
	logLevel = LevelInfo
)

// Init initializes the logger. Can be called multiple times safely.
func Init() {
	once.Do(func() {
		logger = log.New(os.Stderr, "[sdp] ", log.LstdFlags)
	})
}

// SetLevel sets the minimum log level.
func SetLevel(level Level) {
	logLevel = level
}

// Debug logs a debug message.
func Debug(format string, args ...interface{}) {
	if logLevel <= LevelDebug {
		initLogger()
		logger.Printf("[DEBUG] "+format, args...)
	}
}

// Info logs an info message.
func Info(format string, args ...interface{}) {
	if logLevel <= LevelInfo {
		initLogger()
		logger.Printf("[INFO] "+format, args...)
	}
}

// Warn logs a warning message.
func Warn(format string, args ...interface{}) {
	if logLevel <= LevelWarn {
		initLogger()
		logger.Printf("[WARN] "+format, args...)
	}
}

// Error logs an error message.
func Error(format string, args ...interface{}) {
	if logLevel <= LevelError {
		initLogger()
		logger.Printf("[ERROR] "+format, args...)
	}
}

// Operation logs a structured operation event.
func Operation(op, featureID, status string, dur time.Duration) {
	initLogger()
	logger.Printf("[OP] %s feature=%s status=%s duration=%v",
		op, featureID, status, dur)
}

// Session logs session-related events.
func Session(event, featureID, branch string) {
	initLogger()
	logger.Printf("[SESSION] %s feature=%s branch=%s", event, featureID, branch)
}

// Context logs context-related events.
func Context(event, worktree string) {
	initLogger()
	logger.Printf("[CONTEXT] %s worktree=%s", event, worktree)
}

// Guard logs guard-related events.
func Guard(event, featureID, branch string, allowed bool) {
	initLogger()
	status := "allowed"
	if !allowed {
		status = "blocked"
	}
	logger.Printf("[GUARD] %s feature=%s branch=%s status=%s",
		event, featureID, branch, status)
}

func initLogger() {
	if logger == nil {
		Init()
	}
}

// FormatDuration formats a duration for logging.
func FormatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return d.Round(time.Millisecond).String()
}
