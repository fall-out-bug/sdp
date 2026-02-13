package memory

import (
	"fmt"
	"time"

	"github.com/fall-out-bug/sdp/internal/safetylog"
)

// Schema version for migrations
const schemaVersion = 1

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

// GetSchemaVersion returns the current schema version
func (s *Store) GetSchemaVersion() (int, error) {
	var version int
	err := s.db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_version").Scan(&version)
	return version, err
}
