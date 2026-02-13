package memory

import (
	"database/sql"
	"fmt"
	"strings"
)

// scanner is the common interface for sql.Row and sql.Rows
type scanner interface {
	Scan(dest ...any) error
}

// scanArtifactFields scans artifact fields from a row into the artifact struct
func scanArtifactFields(sc scanner, a *Artifact) error {
	var tagsStr string
	var indexedAt sql.NullTime

	err := sc.Scan(
		&a.ID, &a.Path, &a.Type, &a.Title, &a.Content,
		&a.FeatureID, &a.WorkstreamID, &tagsStr, &a.FileHash, &indexedAt,
	)
	if err != nil {
		return err
	}

	// Parse tags
	if tagsStr != "" {
		a.Tags = strings.Split(tagsStr, ",")
	}

	// Parse timestamp
	if indexedAt.Valid {
		a.IndexedAt = indexedAt.Time
	}

	return nil
}

// scanArtifact scans a single artifact from a row
func scanArtifact(row *sql.Row) (*Artifact, error) {
	var a Artifact
	if err := scanArtifactFields(row, &a); err != nil {
		return nil, fmt.Errorf("artifact not found: %w", err)
	}
	return &a, nil
}

// scanArtifacts scans multiple artifacts from rows
func scanArtifacts(rows *sql.Rows) ([]*Artifact, error) {
	var artifacts []*Artifact
	for rows.Next() {
		var a Artifact
		if err := scanArtifactFields(rows, &a); err != nil {
			return nil, fmt.Errorf("failed to scan artifact: %w", err)
		}
		artifacts = append(artifacts, &a)
	}
	return artifacts, rows.Err()
}

const selectArtifactFields = `id, path, type, title, content, feature_id, workstream_id, tags, file_hash, indexed_at`
