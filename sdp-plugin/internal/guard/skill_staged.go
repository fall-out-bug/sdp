package guard

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fall-out-bug/sdp/internal/config"
)

// ParseCheckOptions parses check options from environment variables (AC6: CI diff-range auto-detection)
// Validates SHA format for CI_BASE_SHA and CI_HEAD_SHA (sdp-67l6)
func ParseCheckOptions() CheckOptions {
	base := os.Getenv("CI_BASE_SHA")
	head := os.Getenv("CI_HEAD_SHA")

	// Validate SHA format if provided
	if base != "" && !isValidSHA(base) {
		// Invalid SHA format - use empty to avoid git errors
		base = ""
	}
	if head != "" && !isValidSHA(head) {
		// Invalid SHA format - use empty to avoid git errors
		head = ""
	}

	options := CheckOptions{
		Base: base,
		Head: head,
	}

	return options
}

// isValidSHA validates SHA format (sdp-67l6)
// Empty string is valid (no CI diff range)
// Valid SHA: exactly 40 hexadecimal characters (0-9, a-f, A-F)
func isValidSHA(sha string) bool {
	// Empty is valid (means no CI diff range)
	if sha == "" {
		return true
	}

	// Must be exactly 40 characters
	if len(sha) != 40 {
		return false
	}

	// Must be all hexadecimal
	for _, c := range sha {
		if !isHexChar(c) {
			return false
		}
	}

	return true
}

// isHexChar checks if a rune is a valid hexadecimal digit
func isHexChar(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}
