package guard

import (
	"path/filepath"
	"strings"
)

// matchGlob performs glob matching for path patterns.
// Supports:
//   - * matches any sequence within a path segment
//   - ** matches any sequence across path segments
func matchGlob(pattern, path string) bool {
	if pattern == "" {
		return false
	}

	if pattern == path {
		return true
	}

	if strings.Contains(pattern, "**") {
		return matchDoubleStar(pattern, path)
	}

	if strings.Contains(pattern, "*") {
		matched, err := filepath.Match(pattern, path)
		return err == nil && matched
	}

	return false
}

// matchDoubleStar handles ** glob patterns
func matchDoubleStar(pattern, path string) bool {
	patternParts := strings.Split(pattern, "**")
	if len(patternParts) < 2 {
		return false
	}

	prefix := strings.TrimSuffix(patternParts[0], "/")
	suffix := strings.TrimPrefix(patternParts[len(patternParts)-1], "/")

	// Check prefix
	if prefix != "" {
		if !strings.HasPrefix(path, prefix) {
			return false
		}
		if len(path) > len(prefix) && path[len(prefix)] != '/' {
			return false
		}
	}

	// Check suffix
	if suffix != "" {
		if strings.Contains(suffix, "*") {
			pathBase := filepath.Base(path)
			matched, err := filepath.Match(suffix, pathBase)
			if err != nil || !matched {
				return false
			}
		} else if !strings.HasSuffix(path, suffix) {
			return false
		}
	}

	return true
}
