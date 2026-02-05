package ui

import (
	"fmt"
	"os"
	"strings"
)

// Color codes for terminal output
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[90m"
	Bold   = "\033[1m"
)

// NoColor disables color output (set by --no-color flag)
var NoColor = false

// Check if output is a terminal
func isTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// color wraps text in ANSI color codes if colors are enabled
func color(colorCode, text string) string {
	if NoColor || !isTerminal() {
		return text
	}
	return colorCode + text + Reset
}

// Success returns green text for success messages
func Success(text string) string {
	return color(Green, text)
}

// Error returns red text for error messages
func Error(text string) string {
	return color(Red, text)
}

// Warning returns yellow text for warnings
func Warning(text string) string {
	return color(Yellow, text)
}

// Info returns blue text for informational messages
func Info(text string) string {
	return color(Blue, text)
}

// Dim returns gray text for less important text
func Dim(text string) string {
	return color(Gray, text)
}

// BoldText returns bold text
func BoldText(text string) string {
	if NoColor || !isTerminal() {
		return text
	}
	return Bold + text + Reset
}

// Checkmark returns a green checkmark symbol
func Checkmark() string {
	if NoColor || !isTerminal() {
		return "[OK]"
	}
	return Success("✓")
}

// XMark returns a red X mark symbol
func XMark() string {
	if NoColor || !isTerminal() {
		return "[FAIL]"
	}
	return Error("✗")
}

// WarningSymbol returns a yellow warning symbol
func WarningSymbol() string {
	if NoColor || !isTerminal() {
		return "[WARN]"
	}
	return Warning("⚠")
}

// InfoSymbol returns a blue info symbol
func InfoSymbol() string {
	if NoColor || !isTerminal() {
		return "[INFO]"
	}
	return Info("ℹ")
}

// SuccessLine prints a success message with checkmark
func SuccessLine(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", Checkmark(), fmt.Sprintf(format, args...))
}

// ErrorLine prints an error message with X mark
func ErrorLine(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s %s\n", XMark(), fmt.Sprintf(format, args...))
}

// WarningLine prints a warning message with warning symbol
func WarningLine(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", WarningSymbol(), fmt.Sprintf(format, args...))
}

// InfoLine prints an info message with info symbol
func InfoLine(format string, args ...interface{}) {
	fmt.Printf("%s %s\n", InfoSymbol(), fmt.Sprintf(format, args...))
}

// Header prints a bold header
func Header(text string) {
	fmt.Println()
	fmt.Println(BoldText(text))
	fmt.Println(strings.Repeat("=", len(text)))
}

// Subheader prints a subheader
func Subheader(text string) {
	fmt.Println()
	fmt.Println(BoldText(text))
	fmt.Println(strings.Repeat("-", len(text)))
}
