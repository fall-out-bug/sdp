package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// findProjectRoot finds the project root by looking for .git directory
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // Reached root
		}
		dir = parent
	}

	return "", fmt.Errorf("not in a git repository")
}

// validateFieldLength validates field length
func validateFieldLength(fieldName, value string, maxLen int) error {
	if len(value) > maxLen {
		return fmt.Errorf("%s exceeds maximum length of %d bytes (got %d)", fieldName, maxLen, len(value))
	}
	return nil
}

// stripControlChars removes control characters except newline/tab
func stripControlChars(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			return -1
		}
		return r
	}, s)
}
