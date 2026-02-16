package memory

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/safetylog"
)

// Save stores an artifact
func (s *Store) Save(artifact *Artifact) error {
	return s.SaveContext(context.Background(), artifact)
}

// SaveContext stores an artifact with context support and transaction
func (s *Store) SaveContext(ctx context.Context, artifact *Artifact) error {
	start := time.Now()

	// Use transaction for atomicity
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	tagsStr := strings.Join(artifact.Tags, ",")

	// Use artifact's IndexedAt if set, otherwise use current time
	indexedAt := artifact.IndexedAt
	if indexedAt.IsZero() {
		indexedAt = time.Now()
	}

	// Insert artifact
	_, err = tx.ExecContext(ctx, `
		INSERT OR REPLACE INTO artifacts
		(id, path, type, title, content, feature_id, workstream_id, tags, file_hash, indexed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, artifact.ID, artifact.Path, artifact.Type, artifact.Title, artifact.Content,
		artifact.FeatureID, artifact.WorkstreamID, tagsStr, artifact.FileHash, indexedAt)
	if err != nil {
		return fmt.Errorf("failed to save artifact: %w", err)
	}

	// Update FTS index
	if _, err := tx.ExecContext(ctx, `INSERT OR REPLACE INTO artifacts_fts(rowid, id, title, content)
		SELECT rowid, id, title, content FROM artifacts WHERE id = ?`, artifact.ID); err != nil {
		return fmt.Errorf("failed to update FTS index: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	safetylog.Debug("memory: saved artifact %s (%v)", artifact.ID, time.Since(start))
	return nil
}

// GetByID retrieves an artifact by ID
func (s *Store) GetByID(id string) (*Artifact, error) {
	return s.GetByIDContext(context.Background(), id)
}

// GetByIDContext retrieves an artifact by ID with context support
func (s *Store) GetByIDContext(ctx context.Context, id string) (*Artifact, error) {
	return scanArtifact(s.db.QueryRowContext(ctx,
		`SELECT `+selectArtifactFields+` FROM artifacts WHERE id = ?`, id))
}

// GetByFileHash retrieves an artifact by file hash
func (s *Store) GetByFileHash(hash string) (*Artifact, error) {
	return scanArtifact(s.db.QueryRow(
		`SELECT `+selectArtifactFields+` FROM artifacts WHERE file_hash = ?`, hash))
}

// Close closes the database connection with WAL checkpoint
func (s *Store) Close() error {
	if s.db == nil {
		return nil
	}

	// Perform WAL checkpoint before closing for crash recovery
	if _, err := s.db.Exec("PRAGMA wal_checkpoint(TRUNCATE)"); err != nil {
		safetylog.Warn("memory: WAL checkpoint failed: %v", err)
	}

	return s.db.Close()
}

// Checkpoint forces a WAL checkpoint for crash recovery
func (s *Store) Checkpoint() error {
	_, err := s.db.Exec("PRAGMA wal_checkpoint(TRUNCATE)")
	if err != nil {
		return fmt.Errorf("checkpoint failed: %w", err)
	}
	safetylog.Debug("memory: WAL checkpoint completed")
	return nil
}
