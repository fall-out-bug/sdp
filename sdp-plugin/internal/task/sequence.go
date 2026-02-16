package task

import (
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// nextSequence finds the next sequence number for a prefix
func (c *Creator) nextSequence(prefix string) int {
	matches, err := filepath.Glob(filepath.Join(c.config.WorkstreamDir, prefix+"-*.md"))
	if err != nil || len(matches) == 0 {
		return 1
	}

	maxSeq := 0
	for _, match := range matches {
		base := filepath.Base(match)
		ext := strings.TrimSuffix(base, ".md")
		parts := strings.Split(ext, "-")
		if len(parts) >= 3 {
			if seq, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
				if seq > maxSeq {
					maxSeq = seq
				}
			}
		}
	}

	return maxSeq + 1
}

// nextIssueSequence finds the next issue sequence number
func (c *Creator) nextIssueSequence() int {
	matches, err := filepath.Glob(filepath.Join(c.config.IssuesDir, "ISSUE-*.md"))
	if err != nil || len(matches) == 0 {
		return 1
	}

	sequences := make([]int, 0, len(matches))
	for _, match := range matches {
		base := filepath.Base(match)
		numStr := strings.TrimPrefix(base, "ISSUE-")
		numStr = strings.TrimSuffix(numStr, ".md")
		if seq, err := strconv.Atoi(numStr); err == nil {
			sequences = append(sequences, seq)
		}
	}

	if len(sequences) == 0 {
		return 1
	}

	sort.Ints(sequences)
	return sequences[len(sequences)-1] + 1
}
