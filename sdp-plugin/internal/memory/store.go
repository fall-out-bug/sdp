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

// Schema version for migrations
const schemaVersion = 1

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

	// Enable WAL mode for better concurrency and crash recovery
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Set WAL auto-checkpoint threshold (1000 pages ~ 4MB)
	if _, err := db.Exec("PRAGMA wal_autocheckpoint=1000"); err != nil {
		safetylog.Warn("memory: failed to set WAL autocheckpoint: %v", err)
	}

	store := &Store{db: db, dbPath: dbPath}
	if err := store.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

// initSchema creates the database schema with version tracking
func (s *Store) initSchema() error {
	start := time.Now()
	safetylog.Debug("memory: initializing schema")

	// Create schema version table
	if _, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_version (
			version INTEGER PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return fmt.Errorf("failed to create schema_version table: %w", err)
	}

	// Check current version
	var currentVersion int
	row := s.db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_version")
	if err := row.Scan(&currentVersion); err != nil {
		return fmt.Errorf("failed to check schema version: %w", err)
	}

	// Apply migrations if needed
	if currentVersion < 1 {
		if err := s.migrateV1(); err != nil {
			return err
		}
	}

	safetylog.Debug("memory: schema initialized (version %d, %v)", schemaVersion, time.Since(start))
	return nil
}

// migrateV1 applies the initial schema
func (s *Store) migrateV1() error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin migration: %w", err)
	}
	defer tx.Rollback()

	// Create artifacts table
	if _, err := tx.Exec(`
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
	`); err != nil {
		return fmt.Errorf("failed to create artifacts table: %w", err)
	}

	// Create FTS virtual table
	if _, err := tx.Exec(`CREATE VIRTUAL TABLE IF NOT EXISTS artifacts_fts USING fts5(
		id UNINDEXED, title, content, content='artifacts', content_rowid='rowid'
	)`); err != nil {
		return fmt.Errorf("failed to create FTS table: %w", err)
	}

	// Create indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_artifacts_type ON artifacts(type)",
		"CREATE INDEX IF NOT EXISTS idx_artifacts_file_hash ON artifacts(file_hash)",
		"CREATE INDEX IF NOT EXISTS idx_artifacts_feature_id ON artifacts(feature_id)",
	}
	for _, idx := range indexes {
		if _, err := tx.Exec(idx); err != nil {
			safetylog.Warn("memory: failed to create index: %v", err)
		}
	}

	// Record migration
	if _, err := tx.Exec("INSERT INTO schema_version (version) VALUES (1)"); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	return tx.Commit()
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

	// Insert artifact
	_, err = tx.ExecContext(ctx, `
		INSERT OR REPLACE INTO artifacts
		(id, path, type, title, content, feature_id, workstream_id, tags, file_hash, indexed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, artifact.ID, artifact.Path, artifact.Type, artifact.Title, artifact.Content,
		artifact.FeatureID, artifact.WorkstreamID, tagsStr, artifact.FileHash, time.Now())
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

// GetSchemaVersion returns the current schema version
func (s *Store) GetSchemaVersion() (int, error) {
	var version int
	err := s.db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_version").Scan(&version)
	return version, err
}
