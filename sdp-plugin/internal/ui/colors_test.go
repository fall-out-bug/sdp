package ui

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestSuccess(t *testing.T) {
	NoColor = false
	result := Success("test message")
	if result == "" {
		t.Error("Success() should return non-empty string")
	}
	// When NoColor is false and IS terminal, should contain ANSI codes
	// But in test environment stdout is not a terminal, so we check text is preserved
	if !strings.Contains(result, "test message") {
		t.Errorf("Success() should preserve original text, got: %q", result)
	}
}

func TestError(t *testing.T) {
	NoColor = false
	result := Error("test message")
	if !strings.Contains(result, "test message") {
		t.Errorf("Error() should preserve original text, got: %q", result)
	}
}

func TestWarning(t *testing.T) {
	NoColor = false
	result := Warning("test message")
	if !strings.Contains(result, "test message") {
		t.Errorf("Warning() should preserve original text, got: %q", result)
	}
}

func TestInfo(t *testing.T) {
	NoColor = false
	result := Info("test message")
	if !strings.Contains(result, "test message") {
		t.Errorf("Info() should preserve original text, got: %q", result)
	}
}

func TestDim(t *testing.T) {
	NoColor = false
	result := Dim("test message")
	if !strings.Contains(result, "test message") {
		t.Errorf("Dim() should preserve original text, got: %q", result)
	}
}

func TestBoldText(t *testing.T) {
	NoColor = false
	result := BoldText("test message")
	if !strings.Contains(result, "test message") {
		t.Errorf("BoldText() should preserve original text, got: %q", result)
	}
}

func TestCheckmark(t *testing.T) {
	NoColor = false
	result := Checkmark()
	// In non-terminal environment, returns [OK]
	if result != "✓" && result != "[OK]" {
		t.Errorf("Checkmark() should return ✓ or [OK], got: %q", result)
	}
}

func TestXMark(t *testing.T) {
	NoColor = false
	result := XMark()
	// In non-terminal environment, returns [FAIL]
	if result != "✗" && result != "[FAIL]" {
		t.Errorf("XMark() should return ✗ or [FAIL], got: %q", result)
	}
}

func TestWarningSymbol(t *testing.T) {
	NoColor = false
	result := WarningSymbol()
	// In non-terminal environment, returns [WARN]
	if result != "⚠" && result != "[WARN]" {
		t.Errorf("WarningSymbol() should return ⚠ or [WARN], got: %q", result)
	}
}

func TestInfoSymbol(t *testing.T) {
	NoColor = false
	result := InfoSymbol()
	// In non-terminal environment, returns [INFO]
	if result != "ℹ" && result != "[INFO]" {
		t.Errorf("InfoSymbol() should return ℹ or [INFO], got: %q", result)
	}
}

func TestNoColorMode(t *testing.T) {
	NoColor = true
	result := Success("test message")
	if strings.Contains(result, "\033[") {
		t.Error("Success() should not contain ANSI codes when NoColor is true")
	}
	// Reset for other tests
	NoColor = false
}

func TestSuccessLine(t *testing.T) {
	var buf bytes.Buffer
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	SuccessLine("Operation completed successfully")

	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)

	output := buf.String()
	if !strings.Contains(output, "✓") && !strings.Contains(output, "[OK]") {
		t.Errorf("SuccessLine() should contain checkmark, got: %q", output)
	}
	if !strings.Contains(output, "Operation completed successfully") {
		t.Errorf("SuccessLine() should contain message, got: %q", output)
	}
}

func TestErrorLine(t *testing.T) {
	var buf bytes.Buffer
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	ErrorLine("Operation failed")

	w.Close()
	os.Stderr = oldStderr
	buf.ReadFrom(r)

	output := buf.String()
	if !strings.Contains(output, "✗") && !strings.Contains(output, "[FAIL]") {
		t.Errorf("ErrorLine() should contain X mark, got: %q", output)
	}
	if !strings.Contains(output, "Operation failed") {
		t.Errorf("ErrorLine() should contain message, got: %q", output)
	}
}

func TestWarningLine(t *testing.T) {
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	WarningLine("Operation may fail")

	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)

	output := buf.String()
	if !strings.Contains(output, "⚠") && !strings.Contains(output, "[WARN]") {
		t.Errorf("WarningLine() should contain warning symbol, got: %q", output)
	}
	if !strings.Contains(output, "Operation may fail") {
		t.Errorf("WarningLine() should contain message, got: %q", output)
	}
}

func TestInfoLine(t *testing.T) {
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	InfoLine("Processing data")

	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)

	output := buf.String()
	if !strings.Contains(output, "ℹ") && !strings.Contains(output, "[INFO]") {
		t.Errorf("InfoLine() should contain info symbol, got: %q", output)
	}
	if !strings.Contains(output, "Processing data") {
		t.Errorf("InfoLine() should contain message, got: %q", output)
	}
}

func TestHeader(t *testing.T) {
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Header("Test Header")

	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)

	output := buf.String()
	if !strings.Contains(output, "Test Header") {
		t.Errorf("Header() should contain title, got: %q", output)
	}
	if !strings.Contains(output, "===========") {
		t.Errorf("Header() should contain underline, got: %q", output)
	}
}

func TestSubheader(t *testing.T) {
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Subheader("Test Subheader")

	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)

	output := buf.String()
	if !strings.Contains(output, "Test Subheader") {
		t.Errorf("Subheader() should contain title, got: %q", output)
	}
	if !strings.Contains(output, "--------------") {
		t.Errorf("Subheader() should contain underline, got: %q", output)
	}
}
