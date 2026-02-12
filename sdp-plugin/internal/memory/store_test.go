package memory

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStore_New(t *testing.T) {
	// Create temp directory for test
	tmpDir, err := os.MkdirTemp("", "memory-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "memory.db")
	store, err := NewStore(dbPath)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Verify database file was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}
}

func TestStore_SaveAndGetArtifact(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "memory-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	artifact := &Artifact{
		ID:           "test-id-1",
		Path:         "docs/test.md",
		Type:         "doc",
		Title:        "Test Document",
		Content:      "This is test content",
		FeatureID:    "F001",
		WorkstreamID: "00-001-01",
		Tags:         []string{"test", "example"},
		FileHash:     "abc123",
	}

	// Save artifact
	if err := store.Save(artifact); err != nil {
		t.Fatalf("Failed to save artifact: %v", err)
	}

	// Retrieve artifact
	retrieved, err := store.GetByID("test-id-1")
	if err != nil {
		t.Fatalf("Failed to get artifact: %v", err)
	}

	if retrieved.Path != artifact.Path {
		t.Errorf("Expected path %s, got %s", artifact.Path, retrieved.Path)
	}
	if retrieved.Title != artifact.Title {
		t.Errorf("Expected title %s, got %s", artifact.Title, retrieved.Title)
	}
	if len(retrieved.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(retrieved.Tags))
	}
}

func TestStore_Search(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "memory-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Add test artifacts
	artifacts := []*Artifact{
		{
			ID:       "1",
			Path:     "docs/api.md",
			Type:     "doc",
			Title:    "API Documentation",
			Content:  "This document describes the REST API endpoints for authentication and user management.",
			FileHash: "h1",
		},
		{
			ID:       "2",
			Path:     "docs/guide.md",
			Type:     "doc",
			Title:    "User Guide",
			Content:  "A comprehensive guide for end users on how to use the application.",
			FileHash: "h2",
		},
		{
			ID:       "3",
			Path:     "docs/arch.md",
			Type:     "doc",
			Title:    "Architecture Overview",
			Content:  "System architecture including API design and database schema.",
			FileHash: "h3",
		},
	}

	for _, a := range artifacts {
		if err := store.Save(a); err != nil {
			t.Fatalf("Failed to save artifact: %v", err)
		}
	}

	// Search for "API"
	results, err := store.Search("API")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	// Should find 2 results (API Documentation and Architecture)
	if len(results) < 2 {
		t.Errorf("Expected at least 2 results for 'API', got %d", len(results))
	}
}

func TestStore_GetByFileHash(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "memory-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	artifact := &Artifact{
		ID:       "hash-test",
		Path:     "docs/hash.md",
		Type:     "doc",
		Title:    "Hash Test",
		Content:  "Content",
		FileHash: "unique-hash-123",
	}

	if err := store.Save(artifact); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	retrieved, err := store.GetByFileHash("unique-hash-123")
	if err != nil {
		t.Fatalf("Failed to get by hash: %v", err)
	}

	if retrieved.ID != "hash-test" {
		t.Errorf("Expected ID hash-test, got %s", retrieved.ID)
	}
}

func TestStore_Delete(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "memory-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	artifact := &Artifact{
		ID:       "delete-test",
		Path:     "docs/delete.md",
		Type:     "doc",
		Title:    "To Delete",
		Content:  "Content",
		FileHash: "del-hash",
	}

	if err := store.Save(artifact); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	if err := store.Delete("delete-test"); err != nil {
		t.Fatalf("Failed to delete: %v", err)
	}

	// Verify deleted
	_, err = store.GetByID("delete-test")
	if err == nil {
		t.Error("Expected error getting deleted artifact")
	}
}

func TestStore_ListByType(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "memory-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Add artifacts of different types
	artifacts := []*Artifact{
		{ID: "d1", Path: "docs/a.md", Type: "doc", Title: "Doc 1", Content: "C1", FileHash: "h1"},
		{ID: "d2", Path: "docs/b.md", Type: "doc", Title: "Doc 2", Content: "C2", FileHash: "h2"},
		{ID: "c1", Path: "src/a.go", Type: "code", Title: "Code 1", Content: "C3", FileHash: "h3"},
	}

	for _, a := range artifacts {
		if err := store.Save(a); err != nil {
			t.Fatalf("Failed to save: %v", err)
		}
	}

	docs, err := store.ListByType("doc")
	if err != nil {
		t.Fatalf("Failed to list by type: %v", err)
	}

	if len(docs) != 2 {
		t.Errorf("Expected 2 docs, got %d", len(docs))
	}
}

func TestStore_ListAll(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "memory-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Add artifacts
	artifacts := []*Artifact{
		{ID: "a1", Path: "docs/a.md", Type: "doc", Title: "A", Content: "C1", FileHash: "h1"},
		{ID: "a2", Path: "docs/b.md", Type: "doc", Title: "B", Content: "C2", FileHash: "h2"},
	}

	for _, a := range artifacts {
		if err := store.Save(a); err != nil {
			t.Fatalf("Failed to save: %v", err)
		}
	}

	all, err := store.ListAll()
	if err != nil {
		t.Fatalf("Failed to list all: %v", err)
	}

	if len(all) != 2 {
		t.Errorf("Expected 2 artifacts, got %d", len(all))
	}
}

func TestStore_GetByPath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "memory-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	artifact := &Artifact{
		ID:       "path-test",
		Path:     "unique/path.md",
		Type:     "doc",
		Title:    "Path Test",
		Content:  "Content",
		FileHash: "path-hash",
	}

	if err := store.Save(artifact); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Verify unique path constraint - update with same path should work
	updated := &Artifact{
		ID:       "path-test-updated",
		Path:     "unique/path.md",
		Type:     "doc",
		Title:    "Updated",
		Content:  "Updated content",
		FileHash: "new-hash",
	}

	// This should replace the old one due to UNIQUE constraint on path
	if err := store.Save(updated); err != nil {
		t.Fatalf("Failed to update: %v", err)
	}

	// Old ID should no longer exist
	_, err = store.GetByID("path-test")
	if err == nil {
		t.Error("Expected old artifact to be replaced")
	}
}

func TestStore_SearchFallback(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "memory-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Add test artifact
	artifact := &Artifact{
		ID:       "search-1",
		Path:     "docs/search.md",
		Type:     "doc",
		Title:    "Searchable Document",
		Content:  "This document is searchable using LIKE fallback",
		FileHash: "search-hash",
	}

	if err := store.Save(artifact); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Search with a simple query (will use FTS or fallback)
	results, err := store.Search("searchable")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results) == 0 {
		t.Error("Expected results from search")
	}
}

func TestStore_SearchLike(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "memory-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Add test artifacts
	artifacts := []*Artifact{
		{ID: "like-1", Path: "docs/a.md", Type: "doc", Title: "Alpha Document", Content: "Content about alpha", FileHash: "h1"},
		{ID: "like-2", Path: "docs/b.md", Type: "doc", Title: "Beta Document", Content: "Content about beta", FileHash: "h2"},
		{ID: "like-3", Path: "docs/c.md", Type: "doc", Title: "Gamma", Content: "Alpha and beta combined", FileHash: "h3"},
	}

	for _, a := range artifacts {
		if err := store.Save(a); err != nil {
			t.Fatalf("Failed to save: %v", err)
		}
	}

	// Direct test of searchLike method
	results, err := store.searchLike("alpha")
	if err != nil {
		t.Fatalf("searchLike failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results for 'alpha', got %d", len(results))
	}

	// Test search in content only
	results, err = store.searchLike("combined")
	if err != nil {
		t.Fatalf("searchLike failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'combined', got %d", len(results))
	}

	// Test no match
	results, err = store.searchLike("nonexistent")
	if err != nil {
		t.Fatalf("searchLike failed: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results for 'nonexistent', got %d", len(results))
	}
}

func TestStore_SchemaVersion(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "memory-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Check schema version
	version, err := store.GetSchemaVersion()
	if err != nil {
		t.Fatalf("Failed to get schema version: %v", err)
	}

	if version != 1 {
		t.Errorf("Expected schema version 1, got %d", version)
	}
}

func TestStore_Checkpoint(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "memory-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Add some data
	artifact := &Artifact{
		ID:       "checkpoint-test",
		Path:     "docs/checkpoint.md",
		Type:     "doc",
		Title:    "Checkpoint Test",
		Content:  "Content",
		FileHash: "cp-hash",
	}
	if err := store.Save(artifact); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Test checkpoint
	if err := store.Checkpoint(); err != nil {
		t.Fatalf("Checkpoint failed: %v", err)
	}

	store.Close()
}
