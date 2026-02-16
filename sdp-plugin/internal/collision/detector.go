package collision

import (
	"path/filepath"
	"strings"
)

// WorkstreamScope is a minimal workstream view for overlap detection.
type WorkstreamScope struct {
	ID         string
	Status     string
	ScopeFiles []string
}

// Overlap represents a file or dir touched by multiple workstreams.
type Overlap struct {
	File        string
	Workstreams []string
	Severity    string // "high" = same file, "low" = same dir
}

// DetectOverlaps returns overlaps across in-progress workstreams (AC4, AC5).
func DetectOverlaps(workstreams []WorkstreamScope) []Overlap {
	fileToWS := buildFileToWS(workstreams)
	overlaps := highSeverityOverlaps(fileToWS)
	dirToWS := buildDirToWS(workstreams)
	overlaps = append(overlaps, lowSeverityOverlaps(dirToWS, overlaps)...)
	return overlaps
}

func buildFileToWS(workstreams []WorkstreamScope) map[string][]string {
	fileToWS := make(map[string][]string)
	for _, ws := range workstreams {
		if ws.Status != "in_progress" {
			continue
		}
		for _, f := range ws.ScopeFiles {
			f = normalizePath(f)
			if f != "" {
				fileToWS[f] = append(fileToWS[f], ws.ID)
			}
		}
	}
	return fileToWS
}

func highSeverityOverlaps(fileToWS map[string][]string) []Overlap {
	var out []Overlap
	for file, ids := range fileToWS {
		if len(ids) > 1 {
			out = append(out, Overlap{File: file, Workstreams: ids, Severity: "high"})
		}
	}
	return out
}

func buildDirToWS(workstreams []WorkstreamScope) map[string]map[string]bool {
	dirToWS := make(map[string]map[string]bool)
	for _, ws := range workstreams {
		if ws.Status != "in_progress" {
			continue
		}
		seen := make(map[string]bool)
		for _, f := range ws.ScopeFiles {
			f = normalizePath(f)
			if f == "" {
				continue
			}
			dir := filepath.Dir(f)
			if dir != "." && !seen[dir] {
				seen[dir] = true
				if dirToWS[dir] == nil {
					dirToWS[dir] = make(map[string]bool)
				}
				dirToWS[dir][ws.ID] = true
			}
		}
	}
	return dirToWS
}

func lowSeverityOverlaps(dirToWS map[string]map[string]bool, existing []Overlap) []Overlap {
	var out []Overlap
	for dir, idSet := range dirToWS {
		if len(idSet) <= 1 {
			continue
		}
		ids := make([]string, 0, len(idSet))
		for id := range idSet {
			ids = append(ids, id)
		}
		already := false
		for _, o := range existing {
			if filepath.Dir(o.File) == dir {
				already = true
				break
			}
		}
		if !already {
			out = append(out, Overlap{File: dir + "/", Workstreams: ids, Severity: "low"})
		}
	}
	return out
}

func normalizePath(p string) string {
	p = strings.TrimSpace(p)
	p = strings.Trim(p, "`")
	return filepath.Clean(p)
}
