package security

import (
	"fmt"
	"path/filepath"
	"strings"
)

// SanitizePath sanitizes a user-provided path to prevent path traversal attacks
// Returns cleaned path or error if path contains traversal patterns
func SanitizePath(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	// Clean the path (remove ., resolve .. where safe, clean slashes)
	cleaned := filepath.Clean(input)

	// Check for path traversal patterns
	if strings.Contains(cleaned, "../") {
		return "", fmt.Errorf("path contains traversal pattern '../': %s", input)
	}

	if strings.Contains(cleaned, "..\\") {
		return "", fmt.Errorf("path contains traversal pattern '..\\': %s", input)
	}

	// Block absolute paths (security: user shouldn't access arbitrary system files)
	if filepath.IsAbs(cleaned) {
		return "", fmt.Errorf("absolute paths are not allowed: %s", input)
	}

	// Additional safety: ensure path doesn't start with ..
	if strings.HasPrefix(cleaned, "..") {
		return "", fmt.Errorf("path starts with parent directory reference: %s", input)
	}

	return cleaned, nil
}

// ValidatePathInDirectory ensures that targetPath is within baseDir
// Used to prevent directory traversal attacks
func ValidatePathInDirectory(baseDir, targetPath string) error {
	if baseDir == "" {
		return fmt.Errorf("base directory cannot be empty")
	}

	if targetPath == "" {
		return fmt.Errorf("target path cannot be empty")
	}

	// Clean both paths
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return fmt.Errorf("failed to resolve base directory: %w", err)
	}

	absTarget, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("failed to resolve target path: %w", err)
	}

	// Ensure target starts with base (prevents traversal outside base)
	relPath, err := filepath.Rel(absBase, absTarget)
	if err != nil {
		return fmt.Errorf("failed to compute relative path: %w", err)
	}

	// Check if relative path starts with .. (means it's outside base)
	if strings.HasPrefix(relPath, "..") {
		return fmt.Errorf("target path '%s' is outside base directory '%s'", targetPath, baseDir)
	}

	return nil
}

// SafeJoinPath safely joins base path with user input
// Returns error if the result would escape the base directory
func SafeJoinPath(base, userPath string) (string, error) {
	if base == "" {
		return "", fmt.Errorf("base path cannot be empty")
	}

	if userPath == "" {
		return "", fmt.Errorf("user path cannot be empty")
	}

	// First sanitize the user input
	sanitized, err := SanitizePath(userPath)
	if err != nil {
		return "", err
	}

	// Join paths
	joined := filepath.Join(base, sanitized)

	// Validate the result is within base
	if err := ValidatePathInDirectory(base, joined); err != nil {
		return "", err
	}

	// Clean again to ensure no double slashes or .
	return filepath.Clean(joined), nil
}
