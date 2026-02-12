package memory

import (
	"database/sql"
	"fmt"
	"strings"
)

// scanArtifact scans a single artifact from a row
func scanArtifact(row *sql.Row) (*Artifact, error) {
	var a Artifact
	var tagsStr string
	var indexedAt sql.NullTime

	err := row.Scan(
		&a.ID, &a.Path, &a.Type, &a.Title, &a.Content,
		&a.FeatureID, &a.WorkstreamID, &tagsStr, &a.FileHash, &indexedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("artifact not found: %w", err)
	}

	if tagsStr != "" {
		a.Tags = strings.Split(tagsStr, ",")
	}
	if indexedAt.Valid {
		a.IndexedAt = indexedAt.Time
	}

	return &a, nil
}

// scanArtifacts scans multiple artifacts from rows
func scanArtifacts(rows *sql.Rows) ([]*Artifact, error) {
	var artifacts []*Artifact
	for rows.Next() {
		var a Artifact
		var tagsStr string
		var indexedAt sql.NullTime

		err := rows.Scan(
			&a.ID, &a.Path, &a.Type, &a.Title, &a.Content,
			&a.FeatureID, &a.WorkstreamID, &tagsStr, &a.FileHash, &indexedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan artifact: %w", err)
		}

		if tagsStr != "" {
			a.Tags = strings.Split(tagsStr, ",")
		}
		if indexedAt.Valid {
			a.IndexedAt = indexedAt.Time
		}

		artifacts = append(artifacts, &a)
	}

	return artifacts, rows.Err()
}

const selectArtifactFields = `id, path, type, title, content, feature_id, workstream_id, tags, file_hash, indexed_at`
