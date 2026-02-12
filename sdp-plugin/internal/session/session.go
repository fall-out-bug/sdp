// Package session provides per-worktree session state management for git safety.
// It implements the session file format and operations for tracking worktree identity.
package session

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// SessionVersion is the current version of the session file format
	SessionVersion = "1.0"
	// SessionFileName is the name of the session file within .sdp/
	SessionFileName = "session.json"
	// SDPDir is the SDP configuration directory name
	SDPDir = ".sdp"
)

// Session represents the per-worktree session state for git safety.
// It tracks the expected identity of the worktree to prevent branch confusion.
type Session struct {
	// Version is the session file format version
	Version string `json:"version"`
	// WorktreePath is the absolute path to this worktree
	WorktreePath string `json:"worktree_path"`
	// FeatureID is the feature this worktree serves
	FeatureID string `json:"feature_id"`
	// ExpectedBranch is the branch name for this worktree
	ExpectedBranch string `json:"expected_branch"`
	// ExpectedRemote is the remote tracking branch
	ExpectedRemote string `json:"expected_remote"`
	// CreatedAt is when the session was created
	CreatedAt time.Time `json:"created_at"`
	// CreatedBy indicates what created this session
	CreatedBy string `json:"created_by"`
	// Hash is SHA256 of file content (excluding hash field) for tamper detection
	Hash string `json:"hash"`
}

// Init creates a new session for the given feature and worktree path.
// It initializes all required fields and calculates the integrity hash.
func Init(featureID, worktreePath, createdBy string) (*Session, error) {
	if featureID == "" {
		return nil, fmt.Errorf("feature ID cannot be empty")
	}
	if worktreePath == "" {
		return nil, fmt.Errorf("worktree path cannot be empty")
	}

	now := time.Now().UTC()
	branch := fmt.Sprintf("feature/%s", featureID)
	remote := fmt.Sprintf("origin/feature/%s", featureID)

	session := &Session{
		Version:        SessionVersion,
		WorktreePath:   worktreePath,
		FeatureID:      featureID,
		ExpectedBranch: branch,
		ExpectedRemote: remote,
		CreatedAt:      now,
		CreatedBy:      createdBy,
	}

	// Calculate hash after all fields are set
	session.Hash = session.calculateHash()

	return session, nil
}

// Load reads the session file from the given project root.
// Returns an error if the file doesn't exist or if hash validation fails.
func Load(projectRoot string) (*Session, error) {
	sessionPath := filepath.Join(projectRoot, SDPDir, SessionFileName)
	data, err := os.ReadFile(sessionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("session file not found: %w", err)
		}
		return nil, fmt.Errorf("read session file: %w", err)
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("parse session file: %w", err)
	}

	// Verify hash for tamper detection
	if !session.isValidHash() {
		return nil, fmt.Errorf("session file corrupted or tampered (hash mismatch)")
	}

	return &session, nil
}

// Save writes the session file to the given project root.
// It ensures the .sdp directory exists and recalculates the hash before saving.
func (s *Session) Save(projectRoot string) error {
	sdpDir := filepath.Join(projectRoot, SDPDir)

	// Ensure .sdp directory exists
	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		return fmt.Errorf("create .sdp directory: %w", err)
	}

	// Recalculate hash before saving
	s.Hash = s.calculateHash()

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal session: %w", err)
	}

	sessionPath := filepath.Join(sdpDir, SessionFileName)
	if err := os.WriteFile(sessionPath, data, 0644); err != nil {
		return fmt.Errorf("write session file: %w", err)
	}

	return nil
}

// Sync updates the session with new branch and remote information.
// It recalculates the hash after updating.
func (s *Session) Sync(branch, remote string) *Session {
	s.ExpectedBranch = branch
	s.ExpectedRemote = remote
	s.Hash = s.calculateHash()
	return s
}

// Repair creates a new session from scratch, replacing any corrupted session.
// Use this when the session file is tampered or corrupted.
func Repair(projectRoot, featureID, branch, remote string) (*Session, error) {
	session, err := Init(featureID, projectRoot, "sdp session repair")
	if err != nil {
		return nil, err
	}

	// Override with provided values
	session.ExpectedBranch = branch
	session.ExpectedRemote = remote
	session.Hash = session.calculateHash()

	// Save the repaired session
	if err := session.Save(projectRoot); err != nil {
		return nil, fmt.Errorf("save repaired session: %w", err)
	}

	return session, nil
}

// IsValid checks if the session's hash is valid (not tampered).
func (s *Session) IsValid() bool {
	return s.isValidHash()
}

// isValidHash verifies the integrity hash matches the content.
func (s *Session) isValidHash() bool {
	expectedHash := s.calculateHash()
	return s.Hash == expectedHash
}

// calculateHash generates a SHA256 hash of the session content (excluding the hash field).
func (s *Session) calculateHash() string {
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

// Validate checks that all required fields are present and valid.
func (s *Session) Validate() error {
	if s.Version == "" {
		return fmt.Errorf("version is required")
	}
	if s.WorktreePath == "" {
		return fmt.Errorf("worktree_path is required")
	}
	if s.FeatureID == "" {
		return fmt.Errorf("feature_id is required")
	}
	if s.ExpectedBranch == "" {
		return fmt.Errorf("expected_branch is required")
	}
	if s.ExpectedRemote == "" {
		return fmt.Errorf("expected_remote is required")
	}
	if s.CreatedBy == "" {
		return fmt.Errorf("created_by is required")
	}
	if !s.isValidHash() {
		return fmt.Errorf("hash validation failed")
	}
	return nil
}

// Exists checks if a session file exists at the given project root.
func Exists(projectRoot string) bool {
	sessionPath := filepath.Join(projectRoot, SDPDir, SessionFileName)
	_, err := os.Stat(sessionPath)
	return err == nil
}

// Delete removes the session file from the given project root.
func Delete(projectRoot string) error {
	sessionPath := filepath.Join(projectRoot, SDPDir, SessionFileName)
	if err := os.Remove(sessionPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete session file: %w", err)
	}
	return nil
}
