package memory

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

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

	// Enable WAL mode for better concurrency and crash recovery
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		_ = db.Close() //nolint:errcheck // cleanup on error path
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Set WAL auto-checkpoint threshold (1000 pages ~ 4MB)
	if _, err := db.Exec("PRAGMA wal_autocheckpoint=1000"); err != nil {
		// Log warning but don't fail
	}

	store := &Store{db: db, dbPath: dbPath}
	if err := store.initSchema(); err != nil {
		_ = db.Close() //nolint:errcheck // cleanup on error path
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}
