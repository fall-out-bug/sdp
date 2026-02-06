package ui

import (
	"strings"
	"testing"
)

// TestNoColorModeForAllFunctions tests all color functions with NoColor enabled
func TestNoColorModeForAllFunctions(t *testing.T) {
	// Save original value
	originalNoColor := NoColor
	defer func() { NoColor = originalNoColor }()

	// Test with NoColor = true
	NoColor = true

	tests := []struct {
		name        string
		fn          func() string
		mustContain string
	}{
		{
			name:        "Success with NoColor",
			fn:          func() string { return Success("test") },
			mustContain: "test",
		},
		{
			name:        "Error with NoColor",
			fn:          func() string { return Error("test") },
			mustContain: "test",
		},
		{
			name:        "Warning with NoColor",
			fn:          func() string { return Warning("test") },
			mustContain: "test",
		},
		{
			name:        "Info with NoColor",
			fn:          func() string { return Info("test") },
			mustContain: "test",
		},
		{
			name:        "Dim with NoColor",
			fn:          func() string { return Dim("test") },
			mustContain: "test",
		},
		{
			name:        "BoldText with NoColor",
			fn:          func() string { return BoldText("test") },
			mustContain: "test",
		},
		{
			name:        "Checkmark with NoColor",
			fn:          func() string { return Checkmark() },
			mustContain: "[OK]",
		},
		{
			name:        "XMark with NoColor",
			fn:          func() string { return XMark() },
			mustContain: "[FAIL]",
		},
		{
			name:        "WarningSymbol with NoColor",
			fn:          func() string { return WarningSymbol() },
			mustContain: "[WARN]",
		},
		{
			name:        "InfoSymbol with NoColor",
			fn:          func() string { return InfoSymbol() },
			mustContain: "[INFO]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn()

			// Should not contain ANSI codes when NoColor is true
			if strings.Contains(result, "\033[") {
				t.Errorf("%s should not contain ANSI codes when NoColor=true, got: %q", tt.name, result)
			}

			// Should contain expected text
			if !strings.Contains(result, tt.mustContain) {
				t.Errorf("%s should contain %q, got: %q", tt.name, tt.mustContain, result)
			}
		})
	}
}

// TestNoColorModeNoANSICodes tests that no ANSI codes are present when NoColor is true
func TestNoColorModeNoANSICodes(t *testing.T) {
	originalNoColor := NoColor
	defer func() { NoColor = originalNoColor }()

	NoColor = true

	// Test all color functions
	functions := []func(string) string{
		Success,
		Error,
		Warning,
		Info,
		Dim,
		BoldText,
	}

	for _, fn := range functions {
		result := fn("test")
		if strings.Contains(result, "\033[") {
			t.Errorf("Function should not contain ANSI codes when NoColor=true, got: %q", result)
		}
	}
}

// TestBoldTextWithNoColor tests BoldText specifically with NoColor
func TestBoldTextWithNoColor(t *testing.T) {
	originalNoColor := NoColor
	defer func() { NoColor = originalNoColor }()

	NoColor = true
	result := BoldText("bold text")

	if strings.Contains(result, "\033[") {
		t.Error("BoldText should not contain ANSI codes when NoColor=true")
	}

	if !strings.Contains(result, "bold text") {
		t.Errorf("BoldText should contain original text, got: %q", result)
	}

	// Should not have bold ANSI code
	if strings.Contains(result, "\033[1m") {
		t.Error("BoldText should not have bold ANSI code when NoColor=true")
	}
}

// TestCheckmarkWithNoColor tests Checkmark specifically with NoColor
func TestCheckmarkWithNoColor(t *testing.T) {
	originalNoColor := NoColor
	defer func() { NoColor = originalNoColor }()

	NoColor = true
	result := Checkmark()

	// Should return [OK] when NoColor is true
	if result != "[OK]" {
		t.Errorf("Checkmark should return [OK] when NoColor=true, got: %q", result)
	}

	// Should not contain ANSI codes
	if strings.Contains(result, "\033[") {
		t.Error("Checkmark should not contain ANSI codes when NoColor=true")
	}
}

// TestColorFunctionPreservesText tests that all color functions preserve input text
func TestColorFunctionPreservesText(t *testing.T) {
	originalNoColor := NoColor
	defer func() { NoColor = originalNoColor }()

	NoColor = false

	testCases := []struct {
		name string
		fn   func(string) string
	}{
		{"Success", Success},
		{"Error", Error},
		{"Warning", Warning},
		{"Info", Info},
		{"Dim", Dim},
		{"BoldText", BoldText},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := "test message"
			result := tc.fn(input)

			// Should always contain the input text
			if !strings.Contains(result, input) {
				t.Errorf("%s(%q) = %q should contain input text", tc.name, input, result)
			}
		})
	}
}

// TestSymbolsReturnExpectedValues tests symbol functions return expected values
func TestSymbolsReturnExpectedValues(t *testing.T) {
	originalNoColor := NoColor
	defer func() { NoColor = originalNoColor }()

	NoColor = false

	// These functions should return either symbols or fallback text
	checkmark := Checkmark()
	if checkmark != "✓" && checkmark != "[OK]" {
		t.Errorf("Checkmark() returned unexpected value: %q", checkmark)
	}

	xmark := XMark()
	if xmark != "✗" && xmark != "[FAIL]" {
		t.Errorf("XMark() returned unexpected value: %q", xmark)
	}

	warning := WarningSymbol()
	if warning != "⚠" && warning != "[WARN]" {
		t.Errorf("WarningSymbol() returned unexpected value: %q", warning)
	}

	info := InfoSymbol()
	if info != "ℹ" && info != "[INFO]" {
		t.Errorf("InfoSymbol() returned unexpected value: %q", info)
	}
}

// TestColorWithEmptyString tests color functions with empty string
func TestColorWithEmptyString(t *testing.T) {
	originalNoColor := NoColor
	defer func() { NoColor = originalNoColor }()

	NoColor = false

	// All functions should handle empty string gracefully
	functions := []func(string) string{
		Success,
		Error,
		Warning,
		Info,
		Dim,
		BoldText,
	}

	for _, fn := range functions {
		result := fn("")
		// Should return something (maybe just ANSI codes or empty string)
		// The important thing is it shouldn't crash
		_ = result
	}
}

// TestCheckmarkNonEmpty tests Checkmark returns non-empty string
func TestCheckmarkNonEmpty(t *testing.T) {
	originalNoColor := NoColor
	defer func() { NoColor = originalNoColor }()

	NoColor = false
	result := Checkmark()

	if result == "" {
		t.Error("Checkmark() should return non-empty string")
	}
}

// TestXMarkNonEmpty tests XMark returns non-empty string
func TestXMarkNonEmpty(t *testing.T) {
	originalNoColor := NoColor
	defer func() { NoColor = originalNoColor }()

	NoColor = false
	result := XMark()

	if result == "" {
		t.Error("XMark() should return non-empty string")
	}
}

// TestWarningSymbolNonEmpty tests WarningSymbol returns non-empty string
func TestWarningSymbolNonEmpty(t *testing.T) {
	originalNoColor := NoColor
	defer func() { NoColor = originalNoColor }()

	NoColor = false
	result := WarningSymbol()

	if result == "" {
		t.Error("WarningSymbol() should return non-empty string")
	}
}

// TestInfoSymbolNonEmpty tests InfoSymbol returns non-empty string
func TestInfoSymbolNonEmpty(t *testing.T) {
	originalNoColor := NoColor
	defer func() { NoColor = originalNoColor }()

	NoColor = false
	result := InfoSymbol()

	if result == "" {
		t.Error("InfoSymbol() should return non-empty string")
	}
}
