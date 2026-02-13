package memory

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/evidence"
)

func TestEvidenceAdapter_ImportEvents(t *testing.T) {
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

	adapter := NewEvidenceAdapter(store)

	events := []evidence.Event{
		{
			ID:        "ev-001",
			Type:      "generation",
			Timestamp: "2026-02-13T10:00:00Z",
			WSID:      "00-051-01",
			Data: evidence.GenerationData{
				ModelID:      "claude-sonnet-4",
				ModelVersion: "20250514",
				FilesChanged: []string{"store.go", "queries.go"},
			},
		},
		{
			ID:        "ev-002",
			Type:      "verification",
			Timestamp: "2026-02-13T10:05:00Z",
			WSID:      "00-051-01",
			Data: evidence.VerificationData{
				Passed:   true,
				GateName: "coverage",
				Coverage: 85.5,
			},
		},
		{
			ID:        "ev-003",
			Type:      "decision",
			Timestamp: "2026-02-13T10:10:00Z",
			WSID:      "00-051-02",
			Data: evidence.DecisionEventData{
				Question:  "Use SQLite or PostgreSQL?",
				Choice:    "SQLite",
				Rationale: "Simpler deployment for single-user use case",
			},
		},
	}

	imported, err := adapter.ImportEvents(events)
	if err != nil {
		t.Fatalf("Failed to import events: %v", err)
	}

	if imported != 3 {
		t.Errorf("Expected 3 imported events, got %d", imported)
	}

	// Verify artifacts were created
	all, err := store.ListAll()
	if err != nil {
		t.Fatalf("Failed to list artifacts: %v", err)
	}

	evidenceCount := 0
	for _, a := range all {
		if a.Type == "evidence" {
			evidenceCount++
		}
	}

	if evidenceCount != 3 {
		t.Errorf("Expected 3 evidence artifacts, got %d", evidenceCount)
	}
}

func TestEvidenceAdapter_EventToArtifact(t *testing.T) {
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

	adapter := NewEvidenceAdapter(store)

	// Test generation event
	genEvent := evidence.Event{
		ID:        "ev-gen",
		Type:      "generation",
		Timestamp: "2026-02-13T10:00:00Z",
		WSID:      "00-051-01",
		Data: evidence.GenerationData{
			ModelID:      "claude-sonnet-4",
			ModelVersion: "20250514",
			FilesChanged: []string{"store.go"},
		},
	}

	artifact := adapter.eventToArtifact(genEvent)
	if artifact == nil {
		t.Fatal("Expected artifact, got nil")
	}

	if artifact.Type != "evidence" {
		t.Errorf("Expected type 'evidence', got %s", artifact.Type)
	}

	if artifact.WorkstreamID != "00-051-01" {
		t.Errorf("Expected WSID '00-051-01', got %s", artifact.WorkstreamID)
	}

	if artifact.FeatureID != "F051" {
		t.Errorf("Expected FeatureID 'F051', got %s", artifact.FeatureID)
	}
}

func TestEvidenceAdapter_EmptyEvent(t *testing.T) {
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

	adapter := NewEvidenceAdapter(store)

	// Test empty event (should return nil)
	emptyEvent := evidence.Event{}
	artifact := adapter.eventToArtifact(emptyEvent)
	if artifact != nil {
		t.Error("Expected nil for empty event")
	}

	// Test event with missing WSID
	noWSID := evidence.Event{ID: "test", Type: "generation", Timestamp: "2026-02-13T10:00:00Z"}
	artifact = adapter.eventToArtifact(noWSID)
	if artifact != nil {
		t.Error("Expected nil for event without WSID")
	}
}

func TestExtractFeatureFromWSID(t *testing.T) {
	tests := []struct {
		wsid     string
		expected string
	}{
		{"00-051-01", "F051"},
		{"00-050-02", "F050"},
		{"99-F064-01", "FF064"}, // Fix workstream
		{"short", ""},
		{"", ""},
	}

	for _, tt := range tests {
		result := extractFeatureFromWSID(tt.wsid)
		if result != tt.expected {
			t.Errorf("extractFeatureFromWSID(%s) = %s, want %s", tt.wsid, result, tt.expected)
		}
	}
}
