package memory

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Indexer indexes project artifacts into the memory store
type Indexer struct {
	store   *Store
	docsDir string
}

// NewIndexer creates a new indexer
func NewIndexer(docsDir, dbPath string) (*Indexer, error) {
	store, err := NewStore(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}
	return &Indexer{store: store, docsDir: docsDir}, nil
}

// Close closes the indexer and its store
func (i *Indexer) Close() error {
	return i.store.Close()
}

// IndexDirectory indexes all markdown files in the docs directory
func (i *Indexer) IndexDirectory() (*IndexStats, error) {
	stats := &IndexStats{}

	err := filepath.Walk(i.docsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}

		stats.TotalFiles++
		content, err := os.ReadFile(path)
		if err != nil {
			stats.Errors++
			return nil
		}

		hash := i.ComputeHash(string(content))
		relPath, err := filepath.Rel(filepath.Dir(i.docsDir), path)
		if err != nil {
			relPath = path
		}

		existing, err := i.store.GetByFileHash(hash)
		if err == nil && existing != nil {
			stats.Skipped++
			return nil
		}

		artifactID := i.generateID(relPath)
		_, err = i.store.GetByID(artifactID)
		isUpdate := err == nil

		artifact, err := i.ParseFile(string(content), filepath.Base(path))
		if err != nil {
			stats.Errors++
			return nil
		}

		artifact.ID = artifactID
		artifact.Path = relPath
		artifact.FileHash = hash

		if err := i.store.Save(artifact); err != nil {
			stats.Errors++
			return nil
		}

		if isUpdate {
			stats.Updated++
		} else {
			stats.Indexed++
		}
		return nil
	})

	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}
	return stats, nil
}

// ComputeHash computes SHA256 hash of content
func (i *Indexer) ComputeHash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}

// generateID generates a unique ID for an artifact
func (i *Indexer) generateID(path string) string {
	hash := sha256.Sum256([]byte(path))
	return hex.EncodeToString(hash[:16])
}

// GetStats returns statistics about the store
func (i *Indexer) GetStats() (*StoreStats, error) {
	artifacts, err := i.store.ListAll()
	if err != nil {
		return nil, err
	}

	stats := &StoreStats{
		TotalArtifacts: len(artifacts),
		ByType:         make(map[string]int),
	}

	for _, a := range artifacts {
		stats.ByType[a.Type]++
		if a.IndexedAt.After(stats.LastIndexed) {
			stats.LastIndexed = a.IndexedAt
		}
	}
	return stats, nil
}
