package collision

import (
	"testing"
)

func TestDetectOverlaps_None(t *testing.T) {
	ws := []WorkstreamScope{
		{ID: "00-054-01", Status: "in_progress", ScopeFiles: []string{"schema/index.json"}},
		{ID: "00-054-02", Status: "in_progress", ScopeFiles: []string{"internal/config/project.go"}},
	}
	got := DetectOverlaps(ws)
	if len(got) != 0 {
		t.Errorf("expected no overlaps, got %d", len(got))
	}
}

func TestDetectOverlaps_SameFile(t *testing.T) {
	ws := []WorkstreamScope{
		{ID: "00-054-04", Status: "in_progress", ScopeFiles: []string{"internal/evidence/writer.go"}},
		{ID: "00-054-05", Status: "in_progress", ScopeFiles: []string{"internal/evidence/writer.go"}},
	}
	got := DetectOverlaps(ws)
	if len(got) != 1 {
		t.Fatalf("expected 1 overlap, got %d", len(got))
	}
	if got[0].File != "internal/evidence/writer.go" {
		t.Errorf("File: want internal/evidence/writer.go, got %s", got[0].File)
	}
	if len(got[0].Workstreams) != 2 {
		t.Errorf("Workstreams: want 2, got %d", len(got[0].Workstreams))
	}
	if got[0].Severity != "high" {
		t.Errorf("Severity: want high, got %s", got[0].Severity)
	}
}

func TestDetectOverlaps_SkipNonInProgress(t *testing.T) {
	ws := []WorkstreamScope{
		{ID: "00-054-04", Status: "in_progress", ScopeFiles: []string{"internal/evidence/writer.go"}},
		{ID: "00-054-05", Status: "backlog", ScopeFiles: []string{"internal/evidence/writer.go"}},
	}
	got := DetectOverlaps(ws)
	if len(got) != 0 {
		t.Errorf("expected no overlaps (one backlog), got %d", len(got))
	}
}

func TestDetectOverlaps_SameDir(t *testing.T) {
	ws := []WorkstreamScope{
		{ID: "00-054-04", Status: "in_progress", ScopeFiles: []string{"internal/evidence/writer.go"}},
		{ID: "00-054-05", Status: "in_progress", ScopeFiles: []string{"internal/evidence/reader.go"}},
	}
	got := DetectOverlaps(ws)
	if len(got) != 1 {
		t.Fatalf("expected 1 (same-dir) overlap, got %d", len(got))
	}
	if got[0].Severity != "low" {
		t.Errorf("Severity: want low, got %s", got[0].Severity)
	}
	if got[0].File != "internal/evidence/" {
		t.Errorf("File: want internal/evidence/, got %s", got[0].File)
	}
}
