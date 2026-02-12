// Package session provides per-worktree session state management for git safety.
package session

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// IsValid checks if the session's hash is valid (not tampered).
func (s *Session) IsValid() bool {
	return s.isValidHash()
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

// isValidHash verifies the integrity hash matches the content.
func (s *Session) isValidHash() bool {
	expectedHash := s.calculateHash()
	return s.Hash == expectedHash
}

// calculateHash generates a SHA256 hash of the session content (excluding the hash field).
func (s *Session) calculateHash() string {
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
