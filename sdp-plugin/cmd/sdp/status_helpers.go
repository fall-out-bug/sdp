package main

import (
	"os"
)

// dirExists checks if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// fileExists checks if a file exists (not a directory)
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// boolIcon returns a string representation of a boolean
func boolIcon(b bool) string {
	if b {
		return "[OK]"
	}
	return "[MISSING]"
}

// findJSONField finds the position after a JSON field name
func findJSONField(content string, field string) int {
	search := `"` + field + `"`
	for i := 0; i < len(content)-len(search); i++ {
		if content[i:i+len(search)] == search {
			return i + len(search)
		}
	}
	return -1
}

// extractJSONString extracts a string value after a field position
func extractJSONString(content string, startIdx int) string {
	for i := startIdx; i < len(content); i++ {
		if content[i] == ':' {
			for j := i + 1; j < len(content); j++ {
				if content[j] == '"' {
					for k := j + 1; k < len(content); k++ {
						if content[k] == '"' {
							return content[j+1 : k]
						}
					}
				}
			}
		}
	}
	return ""
}
