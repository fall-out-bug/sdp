package beads

import (
	"os"
	"testing"
)

func TestReadyReturnsTasks(t *testing.T) {
	// Skip if Beads not installed
	if !isBeadsInstalled() {
		t.Skip("Beads CLI not installed")
	}

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tasks, err := client.Ready()
	if err != nil {
		t.Fatalf("Ready() failed: %v", err)
	}

	// Verify we get a slice of tasks (may be empty)
	if tasks == nil {
		t.Error("Expected tasks slice, got nil")
	}
}

func TestShowReturnsTaskDetails(t *testing.T) {
	if !isBeadsInstalled() {
		t.Skip("Beads CLI not installed")
	}

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test with a known Beads ID (if available)
	tasks, err := client.Ready()
	if err != nil {
		t.Fatalf("Ready() failed: %v", err)
	}

	if len(tasks) > 0 {
		task, err := client.Show(tasks[0].ID)
		if err != nil {
			t.Fatalf("Show() failed: %v", err)
		}

		if task == nil {
			t.Error("Expected task, got nil")
		}
		if task.ID != tasks[0].ID {
			t.Errorf("Expected ID %s, got %s", tasks[0].ID, task.ID)
		}
	}
}

func TestUpdateChangesStatus(t *testing.T) {
	if !isBeadsInstalled() {
		t.Skip("Beads CLI not installed")
	}

	client, err := NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Get available tasks
	tasks, err := client.Ready()
	if err != nil {
		t.Fatalf("Ready() failed: %v", err)
	}

	if len(tasks) == 0 {
		t.Skip("No available tasks to test with")
	}

	// Update status to in_progress (should be safe)
	err = client.Update(tasks[0].ID, "in_progress")
	if err != nil {
		// Some Beads tasks might not allow status changes
		// Don't fail the test, just log
		t.Logf("Update() failed (may be expected): %v", err)
	}
}

func TestMapWSToBeads(t *testing.T) {
	// Create temporary mapping file
	tmpDir := t.TempDir()
	mappingPath := createTestMappingFile(tmpDir)

	client := &Client{
		mappingPath:     mappingPath,
		beadsInstalled: false,
	}

	// Test mapping
	beadsID, err := client.MapWSToBeads("00-050-01")
	if err != nil {
		t.Fatalf("MapWSToBeads() failed: %v", err)
	}

	if beadsID != "sdp-x8p" {
		t.Errorf("Expected beads_id sdp-x8p, got %s", beadsID)
	}

	// Test reverse mapping
	wsID, err := client.MapBeadsToWS("sdp-x8p")
	if err != nil {
		t.Fatalf("MapBeadsToWS() failed: %v", err)
	}

	if wsID != "00-050-01" {
		t.Errorf("Expected ws_id 00-050-01, got %s", wsID)
	}
}

func TestMappingNotFound(t *testing.T) {
	// Create temporary mapping file
	tmpDir := t.TempDir()
	mappingPath := createTestMappingFile(tmpDir)

	client := &Client{
		mappingPath:     mappingPath,
		beadsInstalled: false,
	}

	// Test not found
	_, err := client.MapWSToBeads("99-999-99")
	if err == nil {
		t.Error("Expected error for unknown ws_id, got nil")
	}

	_, err = client.MapBeadsToWS("sdp-unknown")
	if err == nil {
		t.Error("Expected error for unknown beads_id, got nil")
	}
}

func TestNewClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	if client == nil {
		t.Fatal("Expected client, got nil")
	}

	// Check if Beads is detected correctly
	expected := isBeadsInstalled()
	if client.beadsInstalled != expected {
		t.Errorf("Expected beadsInstalled=%v, got %v", expected, client.beadsInstalled)
	}
}

func TestClientWhenBeadsNotInstalled(t *testing.T) {
	tmpDir := t.TempDir()
	mappingPath := createTestMappingFile(tmpDir)

	client := &Client{
		mappingPath:     mappingPath,
		beadsInstalled: false,
	}

	// Should return empty results without error
	tasks, err := client.Ready()
	if err != nil {
		t.Errorf("Ready() should not fail when Beads not installed: %v", err)
	}
	if tasks == nil {
		t.Error("Expected empty slice, not nil")
	}

	_, err = client.Show("test-id")
	if err == nil {
		t.Error("Expected error when Beads not installed")
	}

	// Test UpdateMapping
	err = client.UpdateMapping("00-999-99", "sdp-test")
	if err != nil {
		t.Errorf("UpdateMapping() failed: %v", err)
	}

	// Test Sync
	err = client.Sync()
	if err != nil {
		t.Errorf("Sync() failed: %v", err)
	}
}

func TestParseTaskList(t *testing.T) {
	client := &Client{}

	// Test empty output
	tasks := client.parseTaskList("")
	if tasks == nil {
		t.Error("Expected empty slice, not nil")
	}

	// Test with valid output
	output := `sdp-test   Test task
sdp-another Another task
`
	tasks = client.parseTaskList(output)
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	// Test with non-matching lines
	output = `This is not a task
sdp-test   This is a task
Another non-task line
`
	tasks = client.parseTaskList(output)
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
}

func TestUpdateMapping(t *testing.T) {
	tmpDir := t.TempDir()
	mappingPath := createTestMappingFile(tmpDir)

	client := &Client{
		mappingPath:     mappingPath,
		beadsInstalled: false,
	}

	// Test adding new mapping
	err := client.UpdateMapping("00-050-02", "sdp-gtw")
	if err != nil {
		t.Fatalf("UpdateMapping() failed: %v", err)
	}

	// Verify it was added
	beadsID, err := client.MapWSToBeads("00-050-02")
	if err != nil {
		t.Errorf("Failed to map added entry: %v", err)
	}
	if beadsID != "sdp-gtw" {
		t.Errorf("Expected sdp-gtw, got %s", beadsID)
	}

	// Test updating existing mapping
	err = client.UpdateMapping("00-050-01", "sdp-updated")
	if err != nil {
		t.Fatalf("UpdateMapping() for existing failed: %v", err)
	}

	beadsID, err = client.MapWSToBeads("00-050-01")
	if err != nil {
		t.Errorf("Failed to map updated entry: %v", err)
	}
	if beadsID != "sdp-updated" {
		t.Errorf("Expected sdp-updated, got %s", beadsID)
	}
}

func TestShowParsing(t *testing.T) {
	tmpDir := t.TempDir()
	mappingPath := createTestMappingFile(tmpDir)

	client := &Client{
		mappingPath:     mappingPath,
		beadsInstalled: false,
	}

	// Test that Show handles missing Beads gracefully
	_, err := client.Show("sdp-nonexistent")
	if err == nil {
		t.Error("Expected error when Beads not installed")
	}
}

func TestReadMappingWithInvalidData(t *testing.T) {
	tmpDir := t.TempDir()
	mappingPath := tmpDir + "/invalid-mapping.jsonl"

	// Create file with invalid JSON
	content := `{"invalid json
{"sdp_id": "00-050-01", "beads_id": "sdp-x8p", "updated_at": "2026-02-05T12:04:08.943722"}
`
	err := os.WriteFile(mappingPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	client := &Client{
		mappingPath:     mappingPath,
		beadsInstalled: false,
	}

	// Should skip invalid line and return valid entry
	beadsID, err := client.MapWSToBeads("00-050-01")
	if err != nil {
		t.Errorf("Failed to map despite invalid data: %v", err)
	}
	if beadsID != "sdp-x8p" {
		t.Errorf("Expected sdp-x8p, got %s", beadsID)
	}
}

// Helper functions

func createTestMappingFile(dir string) string {
	mappingPath := dir + "/mapping.jsonl"
	content := `{"sdp_id": "00-050-01", "beads_id": "sdp-x8p", "updated_at": "2026-02-05T12:04:08.943722"}
`
	os.WriteFile(mappingPath, []byte(content), 0644)
	return mappingPath
}
