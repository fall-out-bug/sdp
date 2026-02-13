package safetylog

import (
	"bytes"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLevelConstants(t *testing.T) {
	if LevelDebug >= LevelInfo {
		t.Error("LevelDebug should be less than LevelInfo")
	}
	if LevelInfo >= LevelWarn {
		t.Error("LevelInfo should be less than LevelWarn")
	}
	if LevelWarn >= LevelError {
		t.Error("LevelWarn should be less than LevelError")
	}
}

func TestSetLevel(t *testing.T) {
	originalLevel := logLevel
	defer func() { logLevel = originalLevel }()

	SetLevel(LevelDebug)
	if logLevel != LevelDebug {
		t.Errorf("Expected logLevel to be LevelDebug, got %v", logLevel)
	}

	SetLevel(LevelError)
	if logLevel != LevelError {
		t.Errorf("Expected logLevel to be LevelError, got %v", logLevel)
	}
}

func TestInit(t *testing.T) {
	// Reset logger to test init
	logger = nil
	once = sync.Once{}

	Init()
	if logger == nil {
		t.Error("Expected logger to be initialized")
	}

	// Calling Init again should not panic
	Init()
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		input    time.Duration
		expected string
	}{
		{time.Microsecond * 500, "500µs"},
		{time.Millisecond * 5, "5ms"},
		{time.Millisecond * 1500, "1.5s"},
		{time.Second * 2, "2s"},
	}

	for _, tt := range tests {
		result := FormatDuration(tt.input)
		if !strings.Contains(result, strings.TrimSuffix(tt.expected, "s")) {
			t.Errorf("FormatDuration(%v) = %s, want containing %s", tt.input, result, tt.expected)
		}
	}
}

func TestDebugFiltered(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	logger = log.New(&buf, "[sdp] ", 0)
	logLevel = LevelInfo // Debug should be filtered

	Debug("should not appear")

	if strings.Contains(buf.String(), "should not appear") {
		t.Error("Debug message should be filtered at Info level")
	}
}

func TestInfoLogged(t *testing.T) {
	var buf bytes.Buffer
	logger = log.New(&buf, "[sdp] ", 0)
	logLevel = LevelInfo

	Info("test message")

	if !strings.Contains(buf.String(), "test message") {
		t.Error("Info message should be logged")
	}
}

func TestWarnLogged(t *testing.T) {
	var buf bytes.Buffer
	logger = log.New(&buf, "[sdp] ", 0)
	logLevel = LevelInfo

	Warn("warning message")

	if !strings.Contains(buf.String(), "warning message") {
		t.Error("Warn message should be logged")
	}
}

func TestErrorLogged(t *testing.T) {
	var buf bytes.Buffer
	logger = log.New(&buf, "[sdp] ", 0)
	logLevel = LevelInfo

	Error("error message")

	if !strings.Contains(buf.String(), "error message") {
		t.Error("Error message should be logged")
	}
}

func TestOperation(t *testing.T) {
	var buf bytes.Buffer
	logger = log.New(&buf, "[sdp] ", 0)

	Operation("build", "F067", "success", time.Second*5)

	output := buf.String()
	if !strings.Contains(output, "[OP]") {
		t.Error("Operation should contain [OP] prefix")
	}
	if !strings.Contains(output, "F067") {
		t.Error("Operation should contain feature ID")
	}
}

func TestSession(t *testing.T) {
	var buf bytes.Buffer
	logger = log.New(&buf, "[sdp] ", 0)

	Session("start", "F067", "feature/F067")

	output := buf.String()
	if !strings.Contains(output, "[SESSION]") {
		t.Error("Session should contain [SESSION] prefix")
	}
}

func TestContext(t *testing.T) {
	var buf bytes.Buffer
	logger = log.New(&buf, "[sdp] ", 0)

	Context("switch", "worktree-1")

	output := buf.String()
	if !strings.Contains(output, "[CONTEXT]") {
		t.Error("Context should contain [CONTEXT] prefix")
	}
}

func TestGuard(t *testing.T) {
	var buf bytes.Buffer
	logger = log.New(&buf, "[sdp] ", 0)

	tests := []struct {
		allowed  bool
		expected string
	}{
		{true, "allowed"},
		{false, "blocked"},
	}

	for _, tt := range tests {
		buf.Reset()
		Guard("activate", "F067", "main", tt.allowed)

		if !strings.Contains(buf.String(), tt.expected) {
			t.Errorf("Guard with allowed=%v should contain %s", tt.allowed, tt.expected)
		}
	}
}

func TestMain(m *testing.M) {
	// Initialize logger for tests
	Init()
	os.Exit(m.Run())
}
