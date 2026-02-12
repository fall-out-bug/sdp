package memory

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/safetylog"
	_ "modernc.org/sqlite"
)

// Store provides SQLite-based artifact storage with FTS5 search
type Store struct {
	db     *sql.DB
	dbPath string
}

// NewStore creates a new artifact store
func NewStore(dbPath string) (*Store, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	store := &Store{db: db, dbPath: dbPath}
	if err := store.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

// initSchema creates the database schema
func (s *Store) initSchema() error {
	start := time.Now()
	safetylog.Debug("memory: initializing schema")

	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS artifacts (
			id TEXT PRIMARY KEY,
			path TEXT NOT NULL UNIQUE,
			type TEXT NOT NULL,
			title TEXT,
			content TEXT,
			feature_id TEXT,
			workstream_id TEXT,
			tags TEXT,
			file_hash TEXT NOT NULL,
			indexed_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create artifacts table: %w", err)
	}

	// Create FTS virtual table with error handling
	if _, err := s.db.Exec(`CREATE VIRTUAL TABLE IF NOT EXISTS artifacts_fts USING fts5(
		id UNINDEXED, title, content, content='artifacts', content_rowid='rowid'
	)`); err != nil {
		safetylog.Warn("memory: failed to create FTS table: %v", err)
		return fmt.Errorf("failed to create FTS table: %w", err)
	}

	// Create indexes with error handling
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_artifacts_type ON artifacts(type)",
		"CREATE INDEX IF NOT EXISTS idx_artifacts_file_hash ON artifacts(file_hash)",
		"CREATE INDEX IF NOT EXISTS idx_artifacts_feature_id ON artifacts(feature_id)",
	}
	for _, idx := range indexes {
		if _, err := s.db.Exec(idx); err != nil {
			safetylog.Warn("memory: failed to create index: %v", err)
		}
	}

	safetylog.Debug("memory: schema initialized (%v)", time.Since(start))
	return nil
}

// Close closes the database connection
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Save stores an artifact
func (s *Store) Save(artifact *Artifact) error {
	return s.SaveContext(context.Background(), artifact)
}

// SaveContext stores an artifact with context support
func (s *Store) SaveContext(ctx context.Context, artifact *Artifact) error {
	start := time.Now()
	tagsStr := strings.Join(artifact.Tags, ",")
	_, err := s.db.ExecContext(ctx, `
		INSERT OR REPLACE INTO artifacts
		(id, path, type, title, content, feature_id, workstream_id, tags, file_hash, indexed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, artifact.ID, artifact.Path, artifact.Type, artifact.Title, artifact.Content,
		artifact.FeatureID, artifact.WorkstreamID, tagsStr, artifact.FileHash, time.Now())
	if err != nil {
		return fmt.Errorf("failed to save artifact: %w", err)
	}

	if _, err := s.db.ExecContext(ctx, `INSERT OR REPLACE INTO artifacts_fts(rowid, id, title, content)
		SELECT rowid, id, title, content FROM artifacts WHERE id = ?`, artifact.ID); err != nil {
		safetylog.Warn("memory: failed to update FTS index for %s: %v", artifact.ID, err)
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
