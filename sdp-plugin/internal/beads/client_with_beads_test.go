package beads

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestShowWithFakeBeads tests Show with a fake beads binary
func TestShowWithFakeBeads(t *testing.T) {
	// Create a temporary directory with a fake "bd" binary
	tmpDir := t.TempDir()

	// Create a fake bd script that outputs valid task data
	bdScript := `#!/bin/bash
if [ "$1" = "show" ]; then
	echo "ID: $2"
	echo "Title: Test Task"
	echo "Status: open"
	echo "Priority: 2"
elif [ "$1" = "ready" ]; then
	echo "sdp-abc Test task 1"
	echo "sdp-def Test task 2"
elif [ "$1" = "update" ]; then
	# Successful update
	exit 0
fi
`
	bdPath := filepath.Join(tmpDir, "bd")
	if err := os.WriteFile(bdPath, []byte(bdScript), 0755); err != nil {
		t.Fatalf("Failed to create fake bd: %v", err)
	}

	// Add tmpDir to PATH
	oldPath := os.Getenv("PATH")
	t.Cleanup(func() { os.Setenv("PATH", oldPath) })
	os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath)

	// Create client - it should now detect beads as "installed"
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	if !client.beadsInstalled {
		t.Error("Expected beadsInstalled=true with fake bd in PATH")
	}

	// Test Show command
	task, err := client.Show("sdp-abc")
	if err != nil {
		t.Fatalf("Show() failed: %v", err)
	}

	if task == nil {
		t.Fatal("Expected task, got nil")
	}

	if task.ID != "sdp-abc" {
		t.Errorf("Expected ID sdp-abc, got %s", task.ID)
	}

	if task.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got %s", task.Title)
	}

	if task.Status != "open" {
		t.Errorf("Expected status 'open', got %s", task.Status)
	}

	if task.Priority != "2" {
		t.Errorf("Expected priority '2', got %s", task.Priority)
	}
}

// TestUpdateWithFakeBeads tests Update with a fake beads binary
func TestUpdateWithFakeBeads(t *testing.T) {
	// Create a temporary directory with a fake "bd" binary
	tmpDir := t.TempDir()

	// Create a fake bd script
	bdScript := `#!/bin/bash
if [ "$1" = "update" ]; then
	# Successful update
	exit 0
fi
`
	bdPath := filepath.Join(tmpDir, "bd")
	if err := os.WriteFile(bdPath, []byte(bdScript), 0755); err != nil {
		t.Fatalf("Failed to create fake bd: %v", err)
	}

	// Add tmpDir to PATH
	oldPath := os.Getenv("PATH")
	t.Cleanup(func() { os.Setenv("PATH", oldPath) })
	os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath)

	// Create client
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	// Test Update command
	err = client.Update("sdp-abc", "in_progress")
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}
}

// TestShowWithInvalidOutput tests Show with invalid output from beads
func TestShowWithInvalidOutput(t *testing.T) {
	// Create a temporary directory with a fake "bd" binary that returns invalid output
	tmpDir := t.TempDir()

	// Create a fake bd script that returns empty output
	bdScript := `#!/bin/bash
if [ "$1" = "show" ]; then
	# Return empty output
	echo ""
fi
`
	bdPath := filepath.Join(tmpDir, "bd")
	if err := os.WriteFile(bdPath, []byte(bdScript), 0755); err != nil {
		t.Fatalf("Failed to create fake bd: %v", err)
	}

	// Add tmpDir to PATH
	oldPath := os.Getenv("PATH")
	t.Cleanup(func() { os.Setenv("PATH", oldPath) })
	os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath)

	// Create client
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	// Show should still succeed, just with empty fields
	task, err := client.Show("sdp-abc")
	if err != nil {
		t.Fatalf("Show() failed: %v", err)
	}

	if task == nil {
		t.Fatal("Expected task, got nil")
	}

	// Fields should be empty since output was empty
	if task.Title != "" {
		t.Errorf("Expected empty title, got %s", task.Title)
	}
}

// TestShowWithError tests Show when beads command fails
func TestShowWithError(t *testing.T) {
	// Create a temporary directory with a fake "bd" binary that fails
	tmpDir := t.TempDir()

	// Create a fake bd script that exits with error
	bdScript := `#!/bin/bash
if [ "$1" = "show" ]; then
	echo "Error: task not found" >&2
	exit 1
fi
`
	bdPath := filepath.Join(tmpDir, "bd")
	if err := os.WriteFile(bdPath, []byte(bdScript), 0755); err != nil {
		t.Fatalf("Failed to create fake bd: %v", err)
	}

	// Add tmpDir to PATH
	oldPath := os.Getenv("PATH")
	t.Cleanup(func() { os.Setenv("PATH", oldPath) })
	os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath)

	// Create client
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	// Show should return error
	_, err = client.Show("sdp-nonexistent")
	if err == nil {
		t.Error("Expected error when beads command fails")
	}

	if !strings.Contains(err.Error(), "bd show failed") {
		t.Errorf("Expected 'bd show failed' error, got: %v", err)
	}
}

// TestUpdateWithError tests Update when beads command fails
func TestUpdateWithError(t *testing.T) {
	// Create a temporary directory with a fake "bd" binary that fails
	tmpDir := t.TempDir()

	// Create a fake bd script that exits with error
	bdScript := `#!/bin/bash
if [ "$1" = "update" ]; then
	echo "Error: task not found" >&2
	exit 1
fi
`
	bdPath := filepath.Join(tmpDir, "bd")
	if err := os.WriteFile(bdPath, []byte(bdScript), 0755); err != nil {
		t.Fatalf("Failed to create fake bd: %v", err)
	}

	// Add tmpDir to PATH
	oldPath := os.Getenv("PATH")
	t.Cleanup(func() { os.Setenv("PATH", oldPath) })
	os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath)

	// Create client
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	// Update should return error
	err = client.Update("sdp-nonexistent", "in_progress")
	if err == nil {
		t.Error("Expected error when beads command fails")
	}

	if !strings.Contains(err.Error(), "bd update failed") {
		t.Errorf("Expected 'bd update failed' error, got: %v", err)
	}
}

// TestReadyWithFakeBeads tests Ready with fake beads binary
func TestReadyWithFakeBeads(t *testing.T) {
	// Create a temporary directory with a fake "bd" binary
	tmpDir := t.TempDir()

	// Create a fake bd script
	bdScript := `#!/bin/bash
if [ "$1" = "ready" ]; then
	echo "sdp-abc Test task 1"
	echo "sdp-def Test task 2"
	echo "sdp-ghi Test task 3"
fi
`
	bdPath := filepath.Join(tmpDir, "bd")
	if err := os.WriteFile(bdPath, []byte(bdScript), 0755); err != nil {
		t.Fatalf("Failed to create fake bd: %v", err)
	}

	// Add tmpDir to PATH
	oldPath := os.Getenv("PATH")
	t.Cleanup(func() { os.Setenv("PATH", oldPath) })
	os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath)

	// Create client
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	// Test Ready command
	tasks, err := client.Ready()
	if err != nil {
		t.Fatalf("Ready() failed: %v", err)
	}

	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}

	// Check first task
	if tasks[0].ID != "sdp-abc" {
		t.Errorf("Expected ID sdp-abc, got %s", tasks[0].ID)
	}

	if tasks[0].Title != "Test task 1" {
		t.Errorf("Expected title 'Test task 1', got %s", tasks[0].Title)
	}
}

// TestReadyWithEmptyOutput tests Ready when beads returns no tasks
func TestReadyWithEmptyOutput(t *testing.T) {
	// Create a temporary directory with a fake "bd" binary
	tmpDir := t.TempDir()

	// Create a fake bd script that returns empty output
	bdScript := `#!/bin/bash
if [ "$1" = "ready" ]; then
	# No tasks available
	echo ""
fi
`
	bdPath := filepath.Join(tmpDir, "bd")
	if err := os.WriteFile(bdPath, []byte(bdScript), 0755); err != nil {
		t.Fatalf("Failed to create fake bd: %v", err)
	}

	// Add tmpDir to PATH
	oldPath := os.Getenv("PATH")
	t.Cleanup(func() { os.Setenv("PATH", oldPath) })
	os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath)

	// Create client
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	// Test Ready command
	tasks, err := client.Ready()
	if err != nil {
		t.Fatalf("Ready() failed: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}
}

// TestNewClientWithFakeBeads tests NewClient detects fake beads
func TestNewClientWithFakeBeads(t *testing.T) {
	// Create a temporary directory with a fake "bd" binary
	tmpDir := t.TempDir()

	// Create a fake bd script
	bdScript := `#!/bin/bash
echo "fake beads"
`
	bdPath := filepath.Join(tmpDir, "bd")
	if err := os.WriteFile(bdPath, []byte(bdScript), 0755); err != nil {
		t.Fatalf("Failed to create fake bd: %v", err)
	}

	// Add tmpDir to PATH
	oldPath := os.Getenv("PATH")
	t.Cleanup(func() { os.Setenv("PATH", oldPath) })
	os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath)

	// Verify exec.LookPath can find it
	_, err := exec.LookPath("bd")
	if err != nil {
		t.Fatalf("exec.LookPath should find fake bd: %v", err)
	}

	// Create client - should detect beads as installed
	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	if !client.beadsInstalled {
		t.Error("Expected beadsInstalled=true with fake bd in PATH")
	}
}

// TestFindMappingFileNotFound tests findMappingFile when file not found
func TestFindMappingFileNotFound(t *testing.T) {
	// Create empty temp directory (no mapping file)
	tmpDir := t.TempDir()

	oldWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(oldWd) })
	os.Chdir(tmpDir)

	_, err := findMappingFile()
	if err == nil {
		t.Error("Expected error when mapping file not found")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' error, got: %v", err)
	}
}

// TestIsBeadsInstalledNoBeads tests isBeadsInstalled when beads not installed
func TestIsBeadsInstalledNoBeads(t *testing.T) {
	// Create temp directory and ensure no "bd" in PATH
	tmpDir := t.TempDir()

	oldPath := os.Getenv("PATH")
	t.Cleanup(func() { os.Setenv("PATH", oldPath) })

	// Set PATH to only include tmpDir (which has no bd)
	os.Setenv("PATH", tmpDir)

	installed := isBeadsInstalled()
	if installed {
		t.Error("Expected beadsInstalled=false with empty PATH")
	}
}

// TestRunBeadsCommandErrorHandling tests runBeadsCommand error handling
func TestRunBeadsCommandErrorHandling(t *testing.T) {
	// Create a temporary directory with a fake "bd" that fails
	tmpDir := t.TempDir()

	bdScript := `#!/bin/bash
echo "Error: command failed" >&2
exit 1
`
	bdPath := filepath.Join(tmpDir, "bd")
	if err := os.WriteFile(bdPath, []byte(bdScript), 0755); err != nil {
		t.Fatalf("Failed to create fake bd: %v", err)
	}

	oldPath := os.Getenv("PATH")
	t.Cleanup(func() { os.Setenv("PATH", oldPath) })
	os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath)

	client := &Client{
		beadsInstalled: true,
	}

	_, err := client.runBeadsCommand("invalid-command")
	if err == nil {
		t.Error("Expected error from failing bd command")
	}

	if !strings.Contains(err.Error(), "command failed") {
		t.Errorf("Expected 'command failed' error, got: %v", err)
	}
}

// TestWriteMappingErrorHandling tests writeMapping error handling
func TestWriteMappingErrorHandling(t *testing.T) {
	// Try to write to an invalid path
	client := &Client{
		mappingPath: "/root/invalid/path/mapping.jsonl", // Permission denied
	}

	entries := []mappingEntry{
		{SdpID: "00-001-01", BeadsID: "sdp-abc"},
	}

	err := client.writeMapping(entries)
	if err == nil {
		t.Error("Expected error when writing to invalid path")
	}
}

// TestReadMappingErrorHandling tests readMapping error handling
func TestReadMappingErrorHandling(t *testing.T) {
	client := &Client{
		mappingPath: "/nonexistent/file.jsonl",
	}

	_, err := client.readMapping()
	if err == nil {
		t.Error("Expected error when reading nonexistent file")
	}
}

// TestShowParsingWithPartialOutput tests Show parsing with partial output
func TestShowParsingWithPartialOutput(t *testing.T) {
	// Create a temporary directory with a fake "bd" binary
	tmpDir := t.TempDir()

	// Create a fake bd script that returns partial output (only Title)
	bdScript := `#!/bin/bash
if [ "$1" = "show" ]; then
	echo "Title: Only title provided"
fi
`
	bdPath := filepath.Join(tmpDir, "bd")
	if err := os.WriteFile(bdPath, []byte(bdScript), 0755); err != nil {
		t.Fatalf("Failed to create fake bd: %v", err)
	}

	oldPath := os.Getenv("PATH")
	t.Cleanup(func() { os.Setenv("PATH", oldPath) })
	os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	task, err := client.Show("sdp-abc")
	if err != nil {
		t.Fatalf("Show() failed: %v", err)
	}

	// Should have title but other fields empty
	if task.Title != "Only title provided" {
		t.Errorf("Expected title 'Only title provided', got %s", task.Title)
	}

	if task.Status != "" {
		t.Errorf("Expected empty status, got %s", task.Status)
	}
}

// TestMultipleUpdates tests multiple sequential Update calls
func TestMultipleUpdates(t *testing.T) {
	// Create a temporary directory with a fake "bd" binary
	tmpDir := t.TempDir()

	bdScript := `#!/bin/bash
if [ "$1" = "update" ]; then
	exit 0
fi
`
	bdPath := filepath.Join(tmpDir, "bd")
	if err := os.WriteFile(bdPath, []byte(bdScript), 0755); err != nil {
		t.Fatalf("Failed to create fake bd: %v", err)
	}

	oldPath := os.Getenv("PATH")
	t.Cleanup(func() { os.Setenv("PATH", oldPath) })
	os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	// Multiple updates should all succeed
	for i := 0; i < 3; i++ {
		err = client.Update(fmt.Sprintf("sdp-%d", i), "in_progress")
		if err != nil {
			t.Errorf("Update %d failed: %v", i, err)
		}
	}
}

// TestSyncWithFakeBeads tests Sync with fake beads
func TestSyncWithFakeBeads(t *testing.T) {
	// Create a temporary directory with a fake "bd" binary
	tmpDir := t.TempDir()

	bdScript := `#!/bin/bash
if [ "$1" = "sync" ]; then
	exit 0
fi
`
	bdPath := filepath.Join(tmpDir, "bd")
	if err := os.WriteFile(bdPath, []byte(bdScript), 0755); err != nil {
		t.Fatalf("Failed to create fake bd: %v", err)
	}

	oldPath := os.Getenv("PATH")
	t.Cleanup(func() { os.Setenv("PATH", oldPath) })
	os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	err = client.Sync()
	if err != nil {
		t.Errorf("Sync() failed: %v", err)
	}
}

// TestSyncWithError tests Sync when beads command fails
func TestSyncWithError(t *testing.T) {
	// Create a temporary directory with a fake "bd" binary
	tmpDir := t.TempDir()

	bdScript := `#!/bin/bash
if [ "$1" = "sync" ]; then
	echo "Sync failed" >&2
	exit 1
fi
`
	bdPath := filepath.Join(tmpDir, "bd")
	if err := os.WriteFile(bdPath, []byte(bdScript), 0755); err != nil {
		t.Fatalf("Failed to create fake bd: %v", err)
	}

	oldPath := os.Getenv("PATH")
	t.Cleanup(func() { os.Setenv("PATH", oldPath) })
	os.Setenv("PATH", tmpDir+string(os.PathListSeparator)+oldPath)

	client, err := NewClient()
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	err = client.Sync()
	if err == nil {
		t.Error("Expected error when sync fails")
	}
}
