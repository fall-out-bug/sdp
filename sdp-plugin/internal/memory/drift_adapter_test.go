package memory

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fall-out-bug/sdp/internal/drift"
)

func TestDriftAdapter_SaveAndGetDriftReport(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "memory-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	adapter := NewDriftAdapter(store)

	report := &drift.DriftReport{
		WorkstreamID: "00-051-01",
		Timestamp:    time.Now(),
		Verdict:      "PASS",
		Issues: []drift.DriftIssue{
			{
				File:           "store.go",
				Status:         drift.StatusOK,
				Expected:       "File exists",
				Actual:         "File exists",
				Recommendation: "",
			},
		},
	}

	if err := adapter.SaveDriftReport(report); err != nil {
		t.Fatalf("Failed to save drift report: %v", err)
	}

	// Retrieve reports
	reports, err := adapter.GetDriftReports("00-051-01")
	if err != nil {
		t.Fatalf("Failed to get drift reports: %v", err)
	}

	if len(reports) != 1 {
		t.Fatalf("Expected 1 report, got %d", len(reports))
	}

	if reports[0].Type != "drift" {
		t.Errorf("Expected type 'drift', got %s", reports[0].Type)
	}
}

func TestDriftAdapter_SaveEnhancedDriftReport(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "memory-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	adapter := NewDriftAdapter(store)

	report := &drift.EnhancedDriftReport{
		WorkstreamID: "00-051-02",
		Timestamp:    time.Now(),
		Verdict:      "FAIL",
		DriftTypes: []drift.DriftTypeReport{
			{
				Type:     drift.DriftTypeCodeDocs,
				Severity: drift.SeverityError,
				Issues: []drift.EnhancedDriftIssue{
					{
						File:     "queries.go",
						Line:     42,
						Message:  "Function not documented",
						Severity: drift.SeverityError,
					},
				},
				Suggestions: []string{"Add doc comment to queries.go:42"},
			},
		},
	}

	if err := adapter.SaveEnhancedDriftReport(report); err != nil {
		t.Fatalf("Failed to save enhanced drift report: %v", err)
	}

	// Retrieve reports
	reports, err := adapter.GetDriftReports("00-051-02")
	if err != nil {
		t.Fatalf("Failed to get drift reports: %v", err)
	}

	if len(reports) != 1 {
		t.Errorf("Expected 1 report, got %d", len(reports))
	}
}

func TestDriftAdapter_GetLatestDriftReport(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "memory-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	adapter := NewDriftAdapter(store)

	// Save multiple reports with distinct timestamps
	baseTime := time.Date(2026, 2, 13, 10, 0, 0, 0, time.UTC)
	for i := 0; i < 3; i++ {
		report := &drift.DriftReport{
			WorkstreamID: "00-051-03",
			Timestamp:    baseTime.Add(time.Duration(i) * time.Hour),
			Verdict:      "PASS",
			Issues:       []drift.DriftIssue{},
		}
		if err := adapter.SaveDriftReport(report); err != nil {
			t.Fatalf("Failed to save report %d: %v", i, err)
		}
		t.Logf("Saved report %d at %v", i, report.Timestamp)
	}

	// Get all reports to debug
	all, err := store.ListAll()
	if err != nil {
		t.Fatalf("Failed to list all: %v", err)
	}
	t.Logf("Total artifacts: %d", len(all))
	for i, a := range all {
		t.Logf("  [%d] ID=%s Type=%s WSID=%s IndexedAt=%v", i, a.ID, a.Type, a.WorkstreamID, a.IndexedAt)
	}

	// Get latest
	latest, err := adapter.GetLatestDriftReport("00-051-03")
	if err != nil {
		t.Fatalf("Failed to get latest report: %v", err)
	}

	if latest == nil {
		t.Fatal("Expected latest report, got nil")
	}

	t.Logf("Latest report: ID=%s IndexedAt=%v", latest.ID, latest.IndexedAt)

	// Verify it's the most recent (12:00 UTC = baseTime + 2 hours)
	expectedTime := baseTime.Add(2 * time.Hour)
	// Use Equal for exact time comparison since we control the timestamps
	if !latest.IndexedAt.Equal(expectedTime) {
		t.Errorf("Expected latest report at %v, got %v", expectedTime, latest.IndexedAt)
	}
}

func TestDriftAdapter_EmptyWorkstream(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "memory-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	adapter := NewDriftAdapter(store)

	// Get reports for non-existent workstream
	reports, err := adapter.GetDriftReports("00-999-99")
	if err != nil {
		t.Fatalf("Failed to get drift reports: %v", err)
	}

	if len(reports) != 0 {
		t.Errorf("Expected 0 reports for non-existent workstream, got %d", len(reports))
	}

	// Get latest for non-existent workstream
	latest, err := adapter.GetLatestDriftReport("00-999-99")
	if err != nil {
		t.Fatalf("Failed to get latest report: %v", err)
	}

	if latest != nil {
		t.Error("Expected nil for non-existent workstream")
	}
}

func TestSimpleHash(t *testing.T) {
	// Test that simpleHash produces consistent results
	hash1 := simpleHash("test string")
	hash2 := simpleHash("test string")

	if hash1 != hash2 {
		t.Error("simpleHash should produce consistent results")
	}

	// Test different strings produce different hashes
	hash3 := simpleHash("different string")
	if hash1 == hash3 {
		t.Error("simpleHash should produce different results for different inputs")
	}
}

func TestIntToStr(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{123, "123"},
		{-456, "-456"},
	}

	for _, tt := range tests {
		result := intToStr(tt.input)
		if result != tt.expected {
			t.Errorf("intToStr(%d) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}
