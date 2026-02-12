package session

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSessionFormat(t *testing.T) {
	// AC1: Create `.sdp/session.json` format with required fields
	tests := []struct {
		name     string
		session  Session
		expected map[string]bool // fields that should be present
	}{
		{
			name: "all required fields present",
			session: Session{
				Version:        "1.0",
				WorktreePath:   "/Users/user/projects/sdp-F065",
				FeatureID:      "F065",
				ExpectedBranch: "feature/F065",
				ExpectedRemote: "origin/feature/F065",
				CreatedAt:      time.Date(2026, 2, 12, 10, 0, 0, 0, time.UTC),
				CreatedBy:      "sdp worktree create",
			},
			expected: map[string]bool{
				"version":         true,
				"worktree_path":   true,
				"feature_id":      true,
				"expected_branch": true,
				"expected_remote": true,
				"created_at":      true,
				"created_by":      true,
				"hash":            true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.MarshalIndent(tt.session, "", "  ")
			if err != nil {
				t.Fatalf("failed to marshal session: %v", err)
			}

			var result map[string]interface{}
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatalf("failed to unmarshal result: %v", err)
			}

			for field := range tt.expected {
				if _, ok := result[field]; !ok {
					t.Errorf("missing required field: %s", field)
				}
			}
		})
	}
}

func TestSessionInit(t *testing.T) {
	// AC2: Implement `sdp session init` command to create session file
	tmpDir := t.TempDir()
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("failed to create .sdp dir: %v", err)
	}

	tests := []struct {
		name       string
		featureID  string
		worktreePath string
		wantErr    bool
	}{
		{
			name:         "create session for feature",
			featureID:    "F065",
			worktreePath: tmpDir,
			wantErr:      false,
		},
		{
			name:         "empty feature ID fails",
			featureID:    "",
			worktreePath: tmpDir,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := Init(tt.featureID, tt.worktreePath, "sdp session init")
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if session.FeatureID != tt.featureID {
					t.Errorf("FeatureID = %v, want %v", session.FeatureID, tt.featureID)
				}
				if session.Version != "1.0" {
					t.Errorf("Version = %v, want 1.0", session.Version)
				}
				if session.Hash == "" {
					t.Error("Hash should not be empty")
				}
			}
		})
	}
}

func TestSessionSaveAndLoad(t *testing.T) {
	// AC2: Session file should be saved and loaded correctly
	tmpDir := t.TempDir()
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("failed to create .sdp dir: %v", err)
	}

	session, err := Init("F065", tmpDir, "sdp session init")
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Save session
	if err := session.Save(tmpDir); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Load session
	loaded, err := Load(tmpDir)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.FeatureID != session.FeatureID {
		t.Errorf("FeatureID mismatch: got %v, want %v", loaded.FeatureID, session.FeatureID)
	}
	if loaded.Version != session.Version {
		t.Errorf("Version mismatch: got %v, want %v", loaded.Version, session.Version)
	}
	if loaded.Hash != session.Hash {
		t.Errorf("Hash mismatch: got %v, want %v", loaded.Hash, session.Hash)
	}
}

func TestSessionSync(t *testing.T) {
	// AC3: Implement `sdp session sync` to update session from git state
	tmpDir := t.TempDir()
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("failed to create .sdp dir: %v", err)
	}

	// Create initial session
	session, err := Init("F065", tmpDir, "sdp session init")
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Capture original hash before sync
	originalHash := session.Hash

	// Modify and sync
	branch := "feature/F065-updated"
	remote := "origin/feature/F065-updated"
	updated := session.Sync(branch, remote)

	if updated.ExpectedBranch != branch {
		t.Errorf("ExpectedBranch = %v, want %v", updated.ExpectedBranch, branch)
	}
	if updated.ExpectedRemote != remote {
		t.Errorf("ExpectedRemote = %v, want %v", updated.ExpectedRemote, remote)
	}
	// Hash should be recalculated after sync
	if updated.Hash == originalHash {
		t.Errorf("Hash should change after sync: original=%s, updated=%s", originalHash, updated.Hash)
	}
}

func TestSessionRepair(t *testing.T) {
	// AC4: Implement `sdp session repair` to rebuild corrupted sessions
	tmpDir := t.TempDir()
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("failed to create .sdp dir: %v", err)
	}

	// Create initial session
	session, err := Init("F065", tmpDir, "sdp session init")
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := session.Save(tmpDir); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Corrupt the session file by modifying it without updating hash
	sessionPath := filepath.Join(tmpDir, ".sdp", "session.json")
	corruptData, err := os.ReadFile(sessionPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	// Modify the feature ID without updating hash
	corruptData = []byte(strings.Replace(string(corruptData), "F065", "F066", 1))
	if err := os.WriteFile(sessionPath, corruptData, 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	// Verify corruption is detected
	_, err = Load(tmpDir)
	if err == nil {
		t.Error("Load() should fail with corrupted hash")
	}

	// Repair should fix it
	repaired, err := Repair(tmpDir, "F066", "feature/F066", "origin/feature/F066")
	if err != nil {
		t.Fatalf("Repair() error = %v", err)
	}

	// Verify repaired session is valid
	if repaired.FeatureID != "F066" {
		t.Errorf("FeatureID = %v, want F066", repaired.FeatureID)
	}
	if !repaired.IsValid() {
		t.Error("Repaired session should be valid")
	}
}

func TestHashVerification(t *testing.T) {
	// AC5: Add hash verification for tamper detection
	tmpDir := t.TempDir()
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("failed to create .sdp dir: %v", err)
	}

	session, err := Init("F065", tmpDir, "sdp session init")
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Verify hash is calculated correctly
	if session.Hash == "" {
		t.Fatal("Hash should not be empty")
	}

	// Verify hash format (sha256:)
	if !strings.HasPrefix(session.Hash, "sha256:") {
		t.Errorf("Hash should start with 'sha256:', got %v", session.Hash)
	}

	// Verify hash is correct
	expectedHash := calculateExpectedHash(session)
	if session.Hash != expectedHash {
		t.Errorf("Hash = %v, want %v", session.Hash, expectedHash)
	}

	// Verify IsValid returns true for valid session
	if !session.IsValid() {
		t.Error("IsValid() should return true for valid session")
	}

	// Tamper with session
	session.FeatureID = "F066"
	if session.IsValid() {
		t.Error("IsValid() should return false for tampered session")
	}
}

func TestCalculateHash(t *testing.T) {
	// AC5: Hash calculation should be consistent
	session := Session{
		Version:        "1.0",
		WorktreePath:   "/path/to/worktree",
		FeatureID:      "F065",
		ExpectedBranch: "feature/F065",
		ExpectedRemote: "origin/feature/F065",
		CreatedAt:      time.Date(2026, 2, 12, 10, 0, 0, 0, time.UTC),
		CreatedBy:      "sdp worktree create",
	}

	hash1 := session.calculateHash()
	hash2 := session.calculateHash()

	if hash1 != hash2 {
		t.Errorf("Hash should be deterministic: %v != %v", hash1, hash2)
	}

	// Hash should change if content changes
	session.FeatureID = "F066"
	hash3 := session.calculateHash()
	if hash1 == hash3 {
		t.Error("Hash should change when content changes")
	}
}

func TestSessionNotFound(t *testing.T) {
	// Test loading session when file doesn't exist
	tmpDir := t.TempDir()

	_, err := Load(tmpDir)
	if err == nil {
		t.Error("Load() should return error when session file doesn't exist")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		session *Session
		wantErr bool
	}{
		{
			name: "valid session",
			session: &Session{
				Version:        "1.0",
				WorktreePath:   "/path/to/worktree",
				FeatureID:      "F065",
				ExpectedBranch: "feature/F065",
				ExpectedRemote: "origin/feature/F065",
				CreatedAt:      time.Now(),
				CreatedBy:      "sdp session init",
				Hash:           "", // Will be set below
			},
			wantErr: false,
		},
		{
			name: "missing version",
			session: &Session{
				Version:        "",
				WorktreePath:   "/path",
				FeatureID:      "F065",
				ExpectedBranch: "feature/F065",
				ExpectedRemote: "origin/feature/F065",
				CreatedBy:      "sdp session init",
			},
			wantErr: true,
		},
		{
			name: "missing feature_id",
			session: &Session{
				Version:        "1.0",
				WorktreePath:   "/path",
				FeatureID:      "",
				ExpectedBranch: "feature/F065",
				ExpectedRemote: "origin/feature/F065",
				CreatedBy:      "sdp session init",
			},
			wantErr: true,
		},
		{
			name: "missing worktree_path",
			session: &Session{
				Version:        "1.0",
				WorktreePath:   "",
				FeatureID:      "F065",
				ExpectedBranch: "feature/F065",
				ExpectedRemote: "origin/feature/F065",
				CreatedBy:      "sdp session init",
			},
			wantErr: true,
		},
		{
			name: "missing created_by",
			session: &Session{
				Version:        "1.0",
				WorktreePath:   "/path",
				FeatureID:      "F065",
				ExpectedBranch: "feature/F065",
				ExpectedRemote: "origin/feature/F065",
				CreatedBy:      "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set hash for valid session
			if !tt.wantErr {
				tt.session.Hash = tt.session.calculateHash()
			}
			err := tt.session.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExists(t *testing.T) {
	tmpDir := t.TempDir()
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("failed to create .sdp dir: %v", err)
	}

	// Session doesn't exist initially
	if Exists(tmpDir) {
		t.Error("Exists() should return false when session doesn't exist")
	}

	// Create session
	session, err := Init("F065", tmpDir, "sdp session init")
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := session.Save(tmpDir); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Session should exist now
	if !Exists(tmpDir) {
		t.Error("Exists() should return true when session exists")
	}
}

func TestDelete(t *testing.T) {
	tmpDir := t.TempDir()
	sdpDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		t.Fatalf("failed to create .sdp dir: %v", err)
	}

	// Create session
	session, err := Init("F065", tmpDir, "sdp session init")
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := session.Save(tmpDir); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify session exists
	if !Exists(tmpDir) {
		t.Fatal("Session should exist before delete")
	}

	// Delete session
	if err := Delete(tmpDir); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify session is gone
	if Exists(tmpDir) {
		t.Error("Exists() should return false after delete")
	}

	// Delete should be idempotent (no error if file doesn't exist)
	if err := Delete(tmpDir); err != nil {
		t.Errorf("Delete() should not error on non-existent file: %v", err)
	}
}

// Helper function to calculate expected hash
func calculateExpectedHash(s *Session) string {
	// Create a copy without hash for hashing
	copy := struct {
		Version        string    `json:"version"`
		WorktreePath   string    `json:"worktree_path"`
		FeatureID      string    `json:"feature_id"`
		ExpectedBranch string    `json:"expected_branch"`
		ExpectedRemote string    `json:"expected_remote"`
		CreatedAt      time.Time `json:"created_at"`
		CreatedBy      string    `json:"created_by"`
	}{
		Version:        s.Version,
		WorktreePath:   s.WorktreePath,
		FeatureID:      s.FeatureID,
		ExpectedBranch: s.ExpectedBranch,
		ExpectedRemote: s.ExpectedRemote,
		CreatedAt:      s.CreatedAt,
		CreatedBy:      s.CreatedBy,
	}

	data, err := json.Marshal(copy)
	if err != nil {
		return ""
	}

	hash := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(hash[:])
}
