package guard

import (
	"bytes"
	"os/exec"
	"strings"
)

// getStagedFiles returns a list of staged files using git diff --cached
// AC6: CI diff-range auto-detection - uses CI_BASE_SHA and CI_HEAD_SHA when available
func getStagedFiles(opts CheckOptions) ([]string, error) {
	var cmd *exec.Cmd

	// Use CI diff range if provided (AC6)
	if opts.Base != "" && opts.Head != "" {
		cmd = exec.Command("git", "diff", "--name-only", "--diff-filter=ACM", opts.Base, opts.Head)
	} else {
		// Use staged files (local)
		cmd = exec.Command("git", "diff", "--cached", "--name-only", "--diff-filter=ACM")
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// If git command fails (not in a git repo), return empty list
		return []string{}, nil
	}

	// Parse output
	output := stdout.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	var files []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, line)
		}
	}

	return files, nil
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
