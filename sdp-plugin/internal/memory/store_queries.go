package memory

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/safetylog"
)

// Search performs full-text search on artifacts
func (s *Store) Search(query string) ([]*Artifact, error) {
	return s.SearchContext(context.Background(), query)
}

// SearchContext performs full-text search with context support
func (s *Store) SearchContext(ctx context.Context, query string) ([]*Artifact, error) {
	start := time.Now()
	rows, err := s.db.QueryContext(ctx, `
		SELECT a.`+strings.ReplaceAll(selectArtifactFields, ", ", ", a.")+`
		FROM artifacts a
		JOIN artifacts_fts fts ON a.id = fts.id
		WHERE artifacts_fts MATCH ?
		ORDER BY bm25(artifacts_fts) LIMIT 50
	`, query)
	if err != nil {
		safetylog.Debug("memory: FTS search failed, using LIKE fallback")
		return s.searchLike(query)
	}
	defer rows.Close()
	results, err := scanArtifacts(rows)
	if err == nil {
		safetylog.Debug("memory: search '%s' found %d results (%v)", query, len(results), time.Since(start))
	}
	return results, err
}

// searchLike performs a LIKE-based search as fallback
func (s *Store) searchLike(query string) ([]*Artifact, error) {
	likeQuery := "%" + query + "%"
	rows, err := s.db.Query(
		`SELECT `+selectArtifactFields+` FROM artifacts WHERE title LIKE ? OR content LIKE ? LIMIT 50`,
		likeQuery, likeQuery)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	defer rows.Close()
	return scanArtifacts(rows)
}

// Delete removes an artifact by ID
func (s *Store) Delete(id string) error {
	return s.DeleteContext(context.Background(), id)
}

// DeleteContext removes an artifact by ID with context support
func (s *Store) DeleteContext(ctx context.Context, id string) error {
	start := time.Now()
	_, err := s.db.ExecContext(ctx, `DELETE FROM artifacts WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete artifact: %w", err)
	}
	if _, err := s.db.ExecContext(ctx, `DELETE FROM artifacts_fts WHERE id = ?`, id); err != nil {
		safetylog.Warn("memory: failed to delete from FTS index %s: %v", id, err)
	}
	safetylog.Debug("memory: deleted artifact %s (%v)", id, time.Since(start))
	return nil
}

// ListByType lists all artifacts of a given type
func (s *Store) ListByType(artifactType string) ([]*Artifact, error) {
	rows, err := s.db.Query(
		`SELECT `+selectArtifactFields+` FROM artifacts WHERE type = ? ORDER BY path`, artifactType)
	if err != nil {
		return nil, fmt.Errorf("failed to list artifacts: %w", err)
	}
	defer rows.Close()
	return scanArtifacts(rows)
}

// ListAll lists all artifacts
func (s *Store) ListAll() ([]*Artifact, error) {
	rows, err := s.db.Query(`SELECT ` + selectArtifactFields + ` FROM artifacts ORDER BY path`)
	if err != nil {
		return nil, fmt.Errorf("failed to list artifacts: %w", err)
	}
	defer rows.Close()
	return scanArtifacts(rows)
}
