package guard

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestStateManagerSaveAndLoad(t *testing.T) {
	configDir := t.TempDir()
	sm := NewStateManager(configDir)

	state := GuardState{
		ActiveWS:    "00-001-01",
		ActivatedAt: time.Now().Format(time.RFC3339),
		ScopeFiles:  []string{"/file1.go", "/file2.go"},
		Timestamp:   "",
	}

	// Save
	err := sm.Save(state)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists and permissions
	statePath := filepath.Join(configDir, GuardStateFile)
	info, err := os.Stat(statePath)
	if err != nil {
		t.Fatalf("State file not created: %v", err)
	}

	perms := info.Mode().Perm()
	if perms != 0600 {
		t.Errorf("File permissions = %04o, want 0600", perms)
	}

	// Load
	loaded, err := sm.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.ActiveWS != state.ActiveWS {
		t.Errorf("ActiveWS = %s, want %s", loaded.ActiveWS, state.ActiveWS)
	}

	if len(loaded.ScopeFiles) != len(state.ScopeFiles) {
		t.Errorf("ScopeFiles count = %d, want %d", len(loaded.ScopeFiles), len(state.ScopeFiles))
	}

	if loaded.Timestamp == "" {
		t.Error("Timestamp should be set by Save")
	}
}

func TestStateManagerLoadNotExists(t *testing.T) {
	configDir := t.TempDir()
	sm := NewStateManager(configDir)

	// Load when file doesn't exist should return empty state
	loaded, err := sm.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.ActiveWS != "" {
		t.Errorf("ActiveWS should be empty, got %s", loaded.ActiveWS)
	}
}

func TestStateManagerClear(t *testing.T) {
	configDir := t.TempDir()
	sm := NewStateManager(configDir)

	// Create state file
	state := GuardState{
		ActiveWS:    "00-001-01",
		ActivatedAt: time.Now().Format(time.RFC3339),
	}
	sm.Save(state)

	// Verify exists
	statePath := filepath.Join(configDir, GuardStateFile)
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		t.Fatal("State file should exist")
	}

	// Clear
	err := sm.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// Verify removed
	if _, err := os.Stat(statePath); !os.IsNotExist(err) {
		t.Error("State file should be removed")
	}

	// Clear when already gone should not error
	err = sm.Clear()
	if err != nil {
		t.Errorf("Clear on non-existent file failed: %v", err)
	}
}

func TestGuardStateIsExpired(t *testing.T) {
	tests := []struct {
		name        string
		activatedAt string
		wantExpired bool
	}{
		{
			name:        "No active WS",
			activatedAt: "",
			wantExpired: true,
		},
		{
			name:        "Invalid timestamp",
			activatedAt: "invalid",
			wantExpired: true,
		},
		{
			name:        "Recent activation",
			activatedAt: time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			wantExpired: false,
		},
		{
			name:        "Expired (>24 hours)",
			activatedAt: time.Now().Add(-25 * time.Hour).Format(time.RFC3339),
			wantExpired: true,
		},
		{
			name:        "Just under 24 hours",
			activatedAt: time.Now().Add(-23*time.Hour - 59*time.Minute).Format(time.RFC3339),
			wantExpired: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &GuardState{
				ActiveWS:    "00-001-01",
				ActivatedAt: tt.activatedAt,
			}

			got := state.IsExpired()
			if got != tt.wantExpired {
				t.Errorf("IsExpired() = %v, want %v", got, tt.wantExpired)
			}
		})
	}
}
