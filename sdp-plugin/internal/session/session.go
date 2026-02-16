// Package session provides per-worktree session state management for git safety.
package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fall-out-bug/sdp/internal/safetylog"
)

const (
	SessionVersion  = "1.0"
	SessionFileName = "session.json"
	SDPDir          = ".sdp"
)

// Session represents the per-worktree session state for git safety.
type Session struct {
	Version        string    `json:"version"`
	WorktreePath   string    `json:"worktree_path"`
	FeatureID      string    `json:"feature_id"`
	ExpectedBranch string    `json:"expected_branch"`
	ExpectedRemote string    `json:"expected_remote"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      string    `json:"created_by"`
	Hash           string    `json:"hash"`
}

// Init creates a new session for the given feature and worktree path.
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

	s := &Session{
		Version:        SessionVersion,
		WorktreePath:   worktreePath,
		FeatureID:      featureID,
		ExpectedBranch: branch,
		ExpectedRemote: remote,
		CreatedAt:      now,
		CreatedBy:      createdBy,
	}
	s.Hash = s.calculateHash()

	safetylog.Session("init", featureID, branch)
	return s, nil
}

// Load reads the session file from the given project root.
func Load(projectRoot string) (*Session, error) {
	sessionPath := filepath.Join(projectRoot, SDPDir, SessionFileName)
	data, err := os.ReadFile(sessionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("session file not found: %w", err)
		}
		return nil, fmt.Errorf("read session file: %w", err)
	}

	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse session file: %w", err)
	}

	if !s.isValidHash() {
		return nil, fmt.Errorf("session file corrupted or tampered (hash mismatch)")
	}

	return &s, nil
}

// Save writes the session file to the given project root.
func (s *Session) Save(projectRoot string) error {
	sdpDir := filepath.Join(projectRoot, SDPDir)

	if err := os.MkdirAll(sdpDir, 0755); err != nil {
		return fmt.Errorf("create .sdp directory: %w", err)
	}

	s.Hash = s.calculateHash()

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal session: %w", err)
	}

	sessionPath := filepath.Join(sdpDir, SessionFileName)
	if err := os.WriteFile(sessionPath, data, 0644); err != nil {
		return fmt.Errorf("write session file: %w", err)
	}

	safetylog.Debug("session saved: %s", sessionPath)
	return nil
}

// Sync updates the session with new branch and remote information.
func (s *Session) Sync(branch, remote string) *Session {
	s.ExpectedBranch = branch
	s.ExpectedRemote = remote
	s.Hash = s.calculateHash()
	return s
}

// Repair creates a new session from scratch, replacing any corrupted session.
func Repair(projectRoot, featureID, branch, remote string) (*Session, error) {
	s, err := Init(featureID, projectRoot, "sdp session repair")
	if err != nil {
		return nil, err
	}

	s.ExpectedBranch = branch
	s.ExpectedRemote = remote
	s.Hash = s.calculateHash()

	if err := s.Save(projectRoot); err != nil {
		return nil, fmt.Errorf("save repaired session: %w", err)
	}

	return s, nil
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
