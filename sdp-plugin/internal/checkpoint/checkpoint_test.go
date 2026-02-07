package checkpoint

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCheckpointFormat(t *testing.T) {
	// AC1: Checkpoint format defined
	cp := Checkpoint{
		ID:        "test-feature-001",
		FeatureID: "F001",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    StatusInProgress,
		CompletedWorkstreams: []string{
			"00-001-01",
			"00-001-02",
		},
		CurrentWorkstream: "00-001-03",
		Metadata: map[string]interface{}{
			"agent_id":   "abc123",
			"total_ws":   5,
			"completed":  2,
			"start_time": time.Now().Format(time.RFC3339),
		},
	}

	// Verify it can be serialized to JSON
	data, err := json.Marshal(cp)
	if err != nil {
		t.Fatalf("Failed to marshal checkpoint: %v", err)
	}

	// Verify it can be deserialized from JSON
	var decoded Checkpoint
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal checkpoint: %v", err)
	}

	// Verify key fields
	if decoded.ID != cp.ID {
		t.Errorf("Expected ID %s, got %s", cp.ID, decoded.ID)
	}
	if decoded.FeatureID != cp.FeatureID {
		t.Errorf("Expected FeatureID %s, got %s", cp.FeatureID, decoded.FeatureID)
	}
	if len(decoded.CompletedWorkstreams) != 2 {
		t.Errorf("Expected 2 completed workstreams, got %d", len(decoded.CompletedWorkstreams))
	}
}

func TestSaveCheckpoint(t *testing.T) {
	// AC2: Save checkpoint
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	cp := Checkpoint{
		ID:                   "test-save-001",
		FeatureID:            "F001",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		Status:               StatusInProgress,
		CurrentWorkstream:    "00-001-01",
		CompletedWorkstreams: []string{},
		Metadata:             map[string]interface{}{},
	}

	if err := manager.Save(cp); err != nil {
		t.Fatalf("Failed to save checkpoint: %v", err)
	}

	// Verify file was created
	expectedPath := filepath.Join(tmpDir, "test-save-001.json")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Checkpoint file was not created at %s", expectedPath)
	}

	// Verify file contains valid JSON
	data, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("Failed to read checkpoint file: %v", err)
	}

	var saved Checkpoint
	if err := json.Unmarshal(data, &saved); err != nil {
		t.Fatalf("Checkpoint file contains invalid JSON: %v", err)
	}

	if saved.ID != cp.ID {
		t.Errorf("Expected ID %s, got %s", cp.ID, saved.ID)
	}
}

func TestLoadCheckpoint(t *testing.T) {
	// AC3: Load checkpoint
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	// Create a checkpoint
	cp := Checkpoint{
		ID:                   "test-load-001",
		FeatureID:            "F001",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		Status:               StatusInProgress,
		CurrentWorkstream:    "00-001-02",
		CompletedWorkstreams: []string{"00-001-01"},
		Metadata:             map[string]interface{}{},
	}

	if err := manager.Save(cp); err != nil {
		t.Fatalf("Failed to save checkpoint: %v", err)
	}

	// Load it back
	loaded, err := manager.Load("test-load-001")
	if err != nil {
		t.Fatalf("Failed to load checkpoint: %v", err)
	}

	if loaded.ID != cp.ID {
		t.Errorf("Expected ID %s, got %s", cp.ID, loaded.ID)
	}
	if loaded.FeatureID != cp.FeatureID {
		t.Errorf("Expected FeatureID %s, got %s", cp.FeatureID, loaded.FeatureID)
	}
	if loaded.CurrentWorkstream != cp.CurrentWorkstream {
		t.Errorf("Expected CurrentWorkstream %s, got %s", cp.CurrentWorkstream, loaded.CurrentWorkstream)
	}
	if len(loaded.CompletedWorkstreams) != 1 {
		t.Errorf("Expected 1 completed workstream, got %d", len(loaded.CompletedWorkstreams))
	}
}

func TestLoadCheckpoint_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	_, err := manager.Load("nonexistent")
	if err == nil {
		t.Error("Expected error when loading nonexistent checkpoint, got nil")
	}
	if !os.IsNotExist(err) {
		t.Errorf("Expected ErrNotExist, got %v", err)
	}
}

func TestResumeFromCheckpoint(t *testing.T) {
	// AC4: Resume from checkpoint
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	// Create a checkpoint with completed WS
	cp := Checkpoint{
		ID:                   "test-resume-001",
		FeatureID:            "F001",
		CreatedAt:            time.Now().Add(-time.Hour),
		UpdatedAt:            time.Now(),
		Status:               StatusInProgress,
		CurrentWorkstream:    "00-001-03",
		CompletedWorkstreams: []string{"00-001-01", "00-001-02"},
		Metadata: map[string]interface{}{
			"agent_id": "abc123",
		},
	}

	if err := manager.Save(cp); err != nil {
		t.Fatalf("Failed to save checkpoint: %v", err)
	}

	// Resume from checkpoint
	resumed, err := manager.Resume("test-resume-001")
	if err != nil {
		t.Fatalf("Failed to resume checkpoint: %v", err)
	}

	// Verify we can continue from where we left off
	if resumed.ID != cp.ID {
		t.Errorf("Expected ID %s, got %s", cp.ID, resumed.ID)
	}
	if resumed.CurrentWorkstream != "00-001-03" {
		t.Errorf("Expected to resume at 00-001-03, got %s", resumed.CurrentWorkstream)
	}
	if len(resumed.CompletedWorkstreams) != 2 {
		t.Errorf("Expected 2 completed workstreams, got %d", len(resumed.CompletedWorkstreams))
	}

	// Verify metadata is preserved
	if resumed.Metadata["agent_id"] != "abc123" {
		t.Errorf("Expected agent_id abc123, got %v", resumed.Metadata["agent_id"])
	}
}

func TestResumeFromCheckpoint_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	_, err := manager.Resume("nonexistent")
	if err == nil {
		t.Error("Expected error when resuming nonexistent checkpoint, got nil")
	}
}

func TestListCheckpoints(t *testing.T) {
	// AC4: List available checkpoints
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	// Create multiple checkpoints
	checkpoints := []Checkpoint{
		{ID: "feature-001", FeatureID: "F001", CreatedAt: time.Now(), UpdatedAt: time.Now(), Status: StatusInProgress, CurrentWorkstream: "00-001-01", CompletedWorkstreams: []string{}, Metadata: map[string]interface{}{}},
		{ID: "feature-002", FeatureID: "F002", CreatedAt: time.Now(), UpdatedAt: time.Now(), Status: StatusCompleted, CurrentWorkstream: "", CompletedWorkstreams: []string{"00-002-01"}, Metadata: map[string]interface{}{}},
		{ID: "feature-003", FeatureID: "F003", CreatedAt: time.Now(), UpdatedAt: time.Now(), Status: StatusInProgress, CurrentWorkstream: "00-003-01", CompletedWorkstreams: []string{}, Metadata: map[string]interface{}{}},
	}

	for _, cp := range checkpoints {
		if err := manager.Save(cp); err != nil {
			t.Fatalf("Failed to save checkpoint: %v", err)
		}
	}

	// List all checkpoints
	list, err := manager.List()
	if err != nil {
		t.Fatalf("Failed to list checkpoints: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("Expected 3 checkpoints, got %d", len(list))
	}
}

func TestListCheckpoints_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	list, err := manager.List()
	if err != nil {
		t.Fatalf("Failed to list checkpoints: %v", err)
	}

	if len(list) != 0 {
		t.Errorf("Expected 0 checkpoints, got %d", len(list))
	}
}

func TestCleanOldCheckpoints(t *testing.T) {
	// AC5: Clean old checkpoints
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	now := time.Now()

	// Create checkpoints with different ages
	oldCheckpoint := Checkpoint{
		ID:                   "old-feature",
		FeatureID:            "F001",
		CreatedAt:            now.Add(-48 * time.Hour), // 2 days old
		UpdatedAt:            now.Add(-48 * time.Hour),
		Status:               StatusCompleted,
		CurrentWorkstream:    "",
		CompletedWorkstreams: []string{"00-001-01"},
		Metadata:             map[string]interface{}{},
	}

	recentCheckpoint := Checkpoint{
		ID:                   "recent-feature",
		FeatureID:            "F002",
		CreatedAt:            now.Add(-2 * time.Hour), // 2 hours old
		UpdatedAt:            now.Add(-2 * time.Hour),
		Status:               StatusInProgress,
		CurrentWorkstream:    "00-002-01",
		CompletedWorkstreams: []string{},
		Metadata:             map[string]interface{}{},
	}

	if err := manager.Save(oldCheckpoint); err != nil {
		t.Fatalf("Failed to save old checkpoint: %v", err)
	}
	if err := manager.Save(recentCheckpoint); err != nil {
		t.Fatalf("Failed to save recent checkpoint: %v", err)
	}

	// Clean checkpoints older than 24 hours
	deleted, err := manager.Clean(24 * time.Hour)
	if err != nil {
		t.Fatalf("Failed to clean checkpoints: %v", err)
	}

	if deleted != 1 {
		t.Errorf("Expected to delete 1 checkpoint, got %d", deleted)
	}

	// Verify old checkpoint was deleted
	if _, err := manager.Load("old-feature"); !os.IsNotExist(err) {
		t.Error("Expected old checkpoint to be deleted")
	}

	// Verify recent checkpoint still exists
	if _, err := manager.Load("recent-feature"); err != nil {
		t.Errorf("Expected recent checkpoint to exist, got error: %v", err)
	}
}

func TestUpdateCheckpoint(t *testing.T) {
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	// Create initial checkpoint
	cp := Checkpoint{
		ID:                   "test-update-001",
		FeatureID:            "F001",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		Status:               StatusInProgress,
		CurrentWorkstream:    "00-001-01",
		CompletedWorkstreams: []string{},
		Metadata:             map[string]interface{}{},
	}

	if err := manager.Save(cp); err != nil {
		t.Fatalf("Failed to save checkpoint: %v", err)
	}

	// Update checkpoint with completed WS
	cp.CompletedWorkstreams = []string{"00-001-01"}
	cp.CurrentWorkstream = "00-001-02"
	cp.UpdatedAt = time.Now()

	if err := manager.Save(cp); err != nil {
		t.Fatalf("Failed to update checkpoint: %v", err)
	}

	// Verify update persisted
	loaded, err := manager.Load("test-update-001")
	if err != nil {
		t.Fatalf("Failed to load checkpoint: %v", err)
	}

	if len(loaded.CompletedWorkstreams) != 1 {
		t.Errorf("Expected 1 completed workstream, got %d", len(loaded.CompletedWorkstreams))
	}
	if loaded.CurrentWorkstream != "00-001-02" {
		t.Errorf("Expected CurrentWorkstream 00-001-02, got %s", loaded.CurrentWorkstream)
	}
}

func TestCheckpointStatusValidation(t *testing.T) {
	validStatuses := []Status{
		StatusPending,
		StatusInProgress,
		StatusCompleted,
		StatusFailed,
	}

	for _, status := range validStatuses {
		if !status.IsValid() {
			t.Errorf("Expected status %v to be valid", status)
		}
	}

	invalidStatus := Status("invalid")
	if invalidStatus.IsValid() {
		t.Error("Expected invalid status to be invalid")
	}
}

func TestCheckpointDirectory(t *testing.T) {
	// Verify checkpoints are stored in .sdp/checkpoints/
	// This test verifies the default directory structure
	tmpDir := t.TempDir()

	// Create .sdp/checkpoints subdirectory
	checkpointDir := filepath.Join(tmpDir, ".sdp", "checkpoints")
	if err := os.MkdirAll(checkpointDir, 0755); err != nil {
		t.Fatalf("Failed to create checkpoint directory: %v", err)
	}

	manager := NewManager(checkpointDir)

	cp := Checkpoint{
		ID:                   "test-dir-001",
		FeatureID:            "F001",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		Status:               StatusInProgress,
		CurrentWorkstream:    "00-001-01",
		CompletedWorkstreams: []string{},
		Metadata:             map[string]interface{}{},
	}

	if err := manager.Save(cp); err != nil {
		t.Fatalf("Failed to save checkpoint: %v", err)
	}

	// Verify file exists in expected location
	expectedPath := filepath.Join(checkpointDir, "test-dir-001.json")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Checkpoint file not found at %s", expectedPath)
	}
}

func TestSaveCheckpoint_CreateDirectory(t *testing.T) {
	// Test that Save creates directory if it doesn't exist
	tmpDir := t.TempDir()
	checkpointDir := filepath.Join(tmpDir, "newdir", "checkpoints")
	manager := NewManager(checkpointDir)

	cp := Checkpoint{
		ID:                   "test-mkdir-001",
		FeatureID:            "F001",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		Status:               StatusInProgress,
		CurrentWorkstream:    "00-001-01",
		CompletedWorkstreams: []string{},
		Metadata:             map[string]interface{}{},
	}

	// Directory doesn't exist yet
	if _, err := os.Stat(checkpointDir); os.IsNotExist(err) {
		// Expected - directory doesn't exist
	} else {
		t.Fatal("Expected directory to not exist")
	}

	// Save should create directory
	if err := manager.Save(cp); err != nil {
		t.Fatalf("Failed to save checkpoint: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(checkpointDir); os.IsNotExist(err) {
		t.Error("Expected directory to be created")
	}
}

func TestSaveCheckpoint_InvalidJSON(t *testing.T) {
	// Test loading an invalid JSON file
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	// Create invalid JSON file
	filename := filepath.Join(tmpDir, "invalid.json")
	if err := os.WriteFile(filename, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("Failed to create invalid file: %v", err)
	}

	// Load should fail
	_, err := manager.Load("invalid")
	if err == nil {
		t.Error("Expected error loading invalid JSON, got nil")
	}
}

func TestListCheckpoints_SortsByNewest(t *testing.T) {
	// Test that List returns checkpoints sorted by created_at (newest first)
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	now := time.Now()

	// Create checkpoints with different timestamps
	old := Checkpoint{
		ID:                   "old",
		FeatureID:            "F001",
		CreatedAt:            now.Add(-2 * time.Hour),
		UpdatedAt:            now.Add(-2 * time.Hour),
		Status:               StatusCompleted,
		CurrentWorkstream:    "",
		CompletedWorkstreams: []string{"00-001-01"},
		Metadata:             map[string]interface{}{},
	}

	new := Checkpoint{
		ID:                   "new",
		FeatureID:            "F002",
		CreatedAt:            now.Add(-1 * time.Hour),
		UpdatedAt:            now.Add(-1 * time.Hour),
		Status:               StatusInProgress,
		CurrentWorkstream:    "00-002-01",
		CompletedWorkstreams: []string{},
		Metadata:             map[string]interface{}{},
	}

	newest := Checkpoint{
		ID:                   "newest",
		FeatureID:            "F003",
		CreatedAt:            now,
		UpdatedAt:            now,
		Status:               StatusPending,
		CurrentWorkstream:    "",
		CompletedWorkstreams: []string{},
		Metadata:             map[string]interface{}{},
	}

	// Save in random order
	if err := manager.Save(newest); err != nil {
		t.Fatalf("Failed to save newest: %v", err)
	}
	if err := manager.Save(old); err != nil {
		t.Fatalf("Failed to save old: %v", err)
	}
	if err := manager.Save(new); err != nil {
		t.Fatalf("Failed to save new: %v", err)
	}

	// List should return in order: newest, new, old
	list, err := manager.List()
	if err != nil {
		t.Fatalf("Failed to list checkpoints: %v", err)
	}

	if len(list) != 3 {
		t.Fatalf("Expected 3 checkpoints, got %d", len(list))
	}

	if list[0].ID != "newest" {
		t.Errorf("Expected first checkpoint to be 'newest', got '%s'", list[0].ID)
	}
	if list[1].ID != "new" {
		t.Errorf("Expected second checkpoint to be 'new', got '%s'", list[1].ID)
	}
	if list[2].ID != "old" {
		t.Errorf("Expected third checkpoint to be 'old', got '%s'", list[2].ID)
	}
}

func TestListCheckpoints_SkipsInvalidFiles(t *testing.T) {
	// Test that List skips non-JSON files and invalid checkpoints
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	// Create valid checkpoint
	valid := Checkpoint{
		ID:                   "valid",
		FeatureID:            "F001",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		Status:               StatusInProgress,
		CurrentWorkstream:    "00-001-01",
		CompletedWorkstreams: []string{},
		Metadata:             map[string]interface{}{},
	}
	if err := manager.Save(valid); err != nil {
		t.Fatalf("Failed to save valid checkpoint: %v", err)
	}

	// Create non-JSON file
	if err := os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("text"), 0644); err != nil {
		t.Fatalf("Failed to create readme: %v", err)
	}

	// Create invalid JSON file
	if err := os.WriteFile(filepath.Join(tmpDir, "invalid.json"), []byte("{bad json}"), 0644); err != nil {
		t.Fatalf("Failed to create invalid json: %v", err)
	}

	// List should only return the valid checkpoint
	list, err := manager.List()
	if err != nil {
		t.Fatalf("Failed to list checkpoints: %v", err)
	}

	if len(list) != 1 {
		t.Errorf("Expected 1 checkpoint, got %d", len(list))
	}
	if list[0].ID != "valid" {
		t.Errorf("Expected checkpoint 'valid', got '%s'", list[0].ID)
	}
}

func TestClean_OnlyCleansCompletedCheckpoints(t *testing.T) {
	// Test that Clean only removes completed checkpoints, not in-progress ones
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	now := time.Now()

	// Create old in-progress checkpoint
	inProgress := Checkpoint{
		ID:                   "in-progress",
		FeatureID:            "F001",
		CreatedAt:            now.Add(-48 * time.Hour),
		UpdatedAt:            now.Add(-48 * time.Hour),
		Status:               StatusInProgress,
		CurrentWorkstream:    "00-001-01",
		CompletedWorkstreams: []string{},
		Metadata:             map[string]interface{}{},
	}

	// Create old completed checkpoint
	completed := Checkpoint{
		ID:                   "completed",
		FeatureID:            "F002",
		CreatedAt:            now.Add(-48 * time.Hour),
		UpdatedAt:            now.Add(-48 * time.Hour),
		Status:               StatusCompleted,
		CurrentWorkstream:    "",
		CompletedWorkstreams: []string{"00-002-01"},
		Metadata:             map[string]interface{}{},
	}

	if err := manager.Save(inProgress); err != nil {
		t.Fatalf("Failed to save in-progress checkpoint: %v", err)
	}
	if err := manager.Save(completed); err != nil {
		t.Fatalf("Failed to save completed checkpoint: %v", err)
	}

	// Clean old checkpoints
	deleted, err := manager.Clean(24 * time.Hour)
	if err != nil {
		t.Fatalf("Failed to clean checkpoints: %v", err)
	}

	if deleted != 1 {
		t.Errorf("Expected to delete 1 checkpoint, got %d", deleted)
	}

	// Verify in-progress checkpoint still exists
	if _, err := manager.Load("in-progress"); err != nil {
		t.Errorf("Expected in-progress checkpoint to exist: %v", err)
	}

	// Verify completed checkpoint was deleted
	if _, err := manager.Load("completed"); !os.IsNotExist(err) {
		t.Error("Expected completed checkpoint to be deleted")
	}
}

func TestClean_RespectsAge(t *testing.T) {
	// Test that Clean only removes checkpoints older than specified age
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	now := time.Now()

	// Create old completed checkpoint
	old := Checkpoint{
		ID:                   "old",
		FeatureID:            "F001",
		CreatedAt:            now.Add(-48 * time.Hour),
		UpdatedAt:            now.Add(-48 * time.Hour),
		Status:               StatusCompleted,
		CurrentWorkstream:    "",
		CompletedWorkstreams: []string{"00-001-01"},
		Metadata:             map[string]interface{}{},
	}

	// Create recent completed checkpoint
	recent := Checkpoint{
		ID:                   "recent",
		FeatureID:            "F002",
		CreatedAt:            now.Add(-2 * time.Hour),
		UpdatedAt:            now.Add(-2 * time.Hour),
		Status:               StatusCompleted,
		CurrentWorkstream:    "",
		CompletedWorkstreams: []string{"00-002-01"},
		Metadata:             map[string]interface{}{},
	}

	if err := manager.Save(old); err != nil {
		t.Fatalf("Failed to save old checkpoint: %v", err)
	}
	if err := manager.Save(recent); err != nil {
		t.Fatalf("Failed to save recent checkpoint: %v", err)
	}

	// Clean checkpoints older than 24 hours
	deleted, err := manager.Clean(24 * time.Hour)
	if err != nil {
		t.Fatalf("Failed to clean checkpoints: %v", err)
	}

	if deleted != 1 {
		t.Errorf("Expected to delete 1 checkpoint, got %d", deleted)
	}

	// Verify old checkpoint was deleted
	if _, err := manager.Load("old"); !os.IsNotExist(err) {
		t.Error("Expected old checkpoint to be deleted")
	}

	// Verify recent checkpoint still exists
	if _, err := manager.Load("recent"); err != nil {
		t.Errorf("Expected recent checkpoint to exist: %v", err)
	}
}

func TestGetDefaultDir(t *testing.T) {
	// Test GetDefaultDir returns correct path
	dir, err := GetDefaultDir()
	if err != nil {
		t.Fatalf("Failed to get default dir: %v", err)
	}

	// Should end with .sdp/checkpoints
	if filepath.Base(dir) != "checkpoints" {
		t.Errorf("Expected dir to end with 'checkpoints', got %s", filepath.Base(dir))
	}

	parent := filepath.Dir(dir)
	if filepath.Base(parent) != ".sdp" {
		t.Errorf("Expected parent dir to be '.sdp', got %s", filepath.Base(parent))
	}
}

func TestCheckpoint_AllStatusesValid(t *testing.T) {
	// Verify all status constants are valid
	statuses := []Status{
		StatusPending,
		StatusInProgress,
		StatusCompleted,
		StatusFailed,
	}

	for _, status := range statuses {
		if !status.IsValid() {
			t.Errorf("Expected status %s to be valid", status)
		}
	}
}

func TestCheckpoint_MetadataPreserved(t *testing.T) {
	// Test that metadata is preserved through save/load
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	cp := Checkpoint{
		ID:                   "test-meta-001",
		FeatureID:            "F001",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
		Status:               StatusInProgress,
		CurrentWorkstream:    "00-001-01",
		CompletedWorkstreams: []string{},
		Metadata: map[string]interface{}{
			"agent_id":     "abc123",
			"total_ws":     5,
			"completed":    2,
			"start_time":   time.Now().Format(time.RFC3339),
			"custom_field": "custom_value",
			"nested": map[string]interface{}{
				"key": "value",
			},
		},
	}

	if err := manager.Save(cp); err != nil {
		t.Fatalf("Failed to save checkpoint: %v", err)
	}

	loaded, err := manager.Load("test-meta-001")
	if err != nil {
		t.Fatalf("Failed to load checkpoint: %v", err)
	}

	// Verify all metadata fields
	if loaded.Metadata["agent_id"] != "abc123" {
		t.Errorf("Expected agent_id 'abc123', got %v", loaded.Metadata["agent_id"])
	}
	if loaded.Metadata["custom_field"] != "custom_value" {
		t.Errorf("Expected custom_field 'custom_value', got %v", loaded.Metadata["custom_field"])
	}
	if loaded.Metadata["total_ws"] != float64(5) {
		// JSON numbers are unmarshaled as float64
		t.Errorf("Expected total_ws 5, got %v", loaded.Metadata["total_ws"])
	}
}

func TestCheckpoint_CompletedWorkstreamsPreserved(t *testing.T) {
	// Test that completed workstreams list is preserved
	tmpDir := t.TempDir()
	manager := NewManager(tmpDir)

	cp := Checkpoint{
		ID:                "test-ws-001",
		FeatureID:         "F001",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Status:            StatusInProgress,
		CurrentWorkstream: "00-001-04",
		CompletedWorkstreams: []string{
			"00-001-01",
			"00-001-02",
			"00-001-03",
		},
		Metadata: map[string]interface{}{},
	}

	if err := manager.Save(cp); err != nil {
		t.Fatalf("Failed to save checkpoint: %v", err)
	}

	loaded, err := manager.Load("test-ws-001")
	if err != nil {
		t.Fatalf("Failed to load checkpoint: %v", err)
	}

	if len(loaded.CompletedWorkstreams) != 3 {
		t.Fatalf("Expected 3 completed workstreams, got %d", len(loaded.CompletedWorkstreams))
	}

	expected := []string{"00-001-01", "00-001-02", "00-001-03"}
	for i, ws := range expected {
		if loaded.CompletedWorkstreams[i] != ws {
			t.Errorf("Expected %s at index %d, got %s", ws, i, loaded.CompletedWorkstreams[i])
		}
	}
}
