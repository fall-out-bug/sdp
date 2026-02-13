package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestIndexer_New(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "indexer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create docs directory
	docsDir := filepath.Join(tmpDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatalf("Failed to create docs dir: %v", err)
	}

	dbPath := filepath.Join(tmpDir, "memory.db")
	indexer, err := NewIndexer(docsDir, dbPath)
	if err != nil {
		t.Fatalf("Failed to create indexer: %v", err)
	}
	defer indexer.Close()

	if indexer == nil {
		t.Error("Expected non-nil indexer")
	}
}

func TestIndexer_IndexDirectory(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "indexer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create docs directory with test files
	docsDir := filepath.Join(tmpDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatalf("Failed to create docs dir: %v", err)
	}

	// Create test markdown files
	testFiles := map[string]string{
		"guide.md": `---
title: User Guide
tags: [user, guide]
feature_id: F001
---

# User Guide

This is a comprehensive guide for users.
`,
		"api.md": `---
title: API Reference
tags: [api, reference]
ws_id: 00-001-01
---

# API Reference

REST API documentation.
`,
		"no-frontmatter.md": `# Plain Document

This document has no frontmatter.
`,
	}

	for name, content := range testFiles {
		if err := os.WriteFile(filepath.Join(docsDir, name), []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
	}

	dbPath := filepath.Join(tmpDir, "memory.db")
	indexer, err := NewIndexer(docsDir, dbPath)
	if err != nil {
		t.Fatalf("Failed to create indexer: %v", err)
	}
	defer indexer.Close()

	// Index the directory
	stats, err := indexer.IndexDirectory()
	if err != nil {
		t.Fatalf("Failed to index directory: %v", err)
	}

	// Verify stats
	if stats.TotalFiles != 3 {
		t.Errorf("Expected 3 total files, got %d", stats.TotalFiles)
	}
	if stats.Indexed != 3 {
		t.Errorf("Expected 3 indexed files, got %d", stats.Indexed)
	}
}

func TestIndexer_IncrementalUpdate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "indexer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	docsDir := filepath.Join(tmpDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatalf("Failed to create docs dir: %v", err)
	}

	// Create initial file
	file1 := filepath.Join(docsDir, "file1.md")
	if err := os.WriteFile(file1, []byte("---\ntitle: File 1\n---\n\nContent 1"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	dbPath := filepath.Join(tmpDir, "memory.db")
	indexer, err := NewIndexer(docsDir, dbPath)
	if err != nil {
		t.Fatalf("Failed to create indexer: %v", err)
	}

	// First index
	stats1, err := indexer.IndexDirectory()
	if err != nil {
		t.Fatalf("First index failed: %v", err)
	}
	if stats1.Indexed != 1 {
		t.Errorf("Expected 1 indexed, got %d", stats1.Indexed)
	}

	// Modify file (change content)
	if err := os.WriteFile(file1, []byte("---\ntitle: File 1 Modified\n---\n\nNew Content"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	// Second index (incremental)
	stats2, err := indexer.IndexDirectory()
	if err != nil {
		t.Fatalf("Second index failed: %v", err)
	}

	// Should detect change and reindex
	if stats2.Updated != 1 {
		t.Errorf("Expected 1 updated, got %d", stats2.Updated)
	}

	// Third index (no changes)
	stats3, err := indexer.IndexDirectory()
	if err != nil {
		t.Fatalf("Third index failed: %v", err)
	}

	// Should skip unchanged file
	if stats3.Skipped != 1 {
		t.Errorf("Expected 1 skipped, got %d", stats3.Skipped)
	}

	indexer.Close()
}

func TestIndexer_ExtractMetadata(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		filename     string
		expectTitle  string
		expectFID    string
		expectWSID   string
		expectTags   int
	}{
		{
			name: "full frontmatter",
			content: `---
title: Test Doc
feature_id: F050
ws_id: 00-050-01
tags: [test, example]
---

# Content`,
			filename:    "test.md",
			expectTitle: "Test Doc",
			expectFID:   "F050",
			expectWSID:  "00-050-01",
			expectTags:  2,
		},
		{
			name: "no frontmatter, extract from heading",
			content: `# My Document

Some content here.`,
			filename:    "plain.md",
			expectTitle: "My Document",
			expectFID:   "",
			expectWSID:  "",
			expectTags:  0,
		},
		{
			name: "extract ws_id from filename",
			content: `---
title: Workstream Doc
---

Content`,
			filename:    "00-051-02.md",
			expectTitle: "Workstream Doc",
			expectFID:   "F051",
			expectWSID:  "00-051-02",
			expectTags:  0,
		},
		{
			name: "string tags format",
			content: `---
title: String Tags
tags: tag1, tag2, tag3
---

Content`,
			filename:    "tags.md",
			expectTitle: "String Tags",
			expectTags:  3,
		},
		{
			name: "invalid frontmatter falls back to content",
			content: `---
invalid yaml content [
---

# Fallback Title`,
			filename:    "invalid.md",
			expectTitle: "Fallback Title",
		},
		{
			name: "no title uses filename",
			content: `Just some content without headings.`,
			filename:    "no-title-file.md",
			expectTitle: "no-title-file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "metadata-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			docsDir := filepath.Join(tmpDir, "docs")
			if err := os.MkdirAll(docsDir, 0755); err != nil {
				t.Fatalf("Failed to create docs dir: %v", err)
			}

			dbPath := filepath.Join(tmpDir, "memory.db")
			indexer, err := NewIndexer(docsDir, dbPath)
			if err != nil {
				t.Fatalf("Failed to create indexer: %v", err)
			}
			defer indexer.Close()

			artifact, err := indexer.ParseFile(tt.content, tt.filename)
			if err != nil {
				t.Fatalf("ParseFile failed: %v", err)
			}

			if artifact.Title != tt.expectTitle {
				t.Errorf("Expected title %q, got %q", tt.expectTitle, artifact.Title)
			}
			if artifact.FeatureID != tt.expectFID {
				t.Errorf("Expected feature_id %q, got %q", tt.expectFID, artifact.FeatureID)
			}
			if artifact.WorkstreamID != tt.expectWSID {
				t.Errorf("Expected ws_id %q, got %q", tt.expectWSID, artifact.WorkstreamID)
			}
			if tt.expectTags > 0 && len(artifact.Tags) != tt.expectTags {
				t.Errorf("Expected %d tags, got %d", tt.expectTags, len(artifact.Tags))
			}
		})
	}
}

func TestIndexer_ComputeHash(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "hash-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	docsDir := filepath.Join(tmpDir, "docs")
	dbPath := filepath.Join(tmpDir, "memory.db")

	indexer, err := NewIndexer(docsDir, dbPath)
	if err != nil {
		t.Fatalf("Failed to create indexer: %v", err)
	}
	defer indexer.Close()

	content := "test content for hashing"
	hash1 := indexer.ComputeHash(content)
	hash2 := indexer.ComputeHash(content)

	// Same content should produce same hash
	if hash1 != hash2 {
		t.Error("Same content should produce same hash")
	}

	// Different content should produce different hash
	hash3 := indexer.ComputeHash("different content")
	if hash1 == hash3 {
		t.Error("Different content should produce different hash")
	}

	// Hash should be SHA256 length (64 hex chars)
	if len(hash1) != 64 {
		t.Errorf("Expected hash length 64, got %d", len(hash1))
	}
}

func TestIndexer_GetStats(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "stats-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	docsDir := filepath.Join(tmpDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatalf("Failed to create docs dir: %v", err)
	}

	// Create test files (different content to get different hashes)
	for i := 0; i < 5; i++ {
		content := fmt.Sprintf("---\ntitle: Doc %d\n---\n\nContent %d", i, i)
		filename := fmt.Sprintf("doc%d.md", i)
		if err := os.WriteFile(filepath.Join(docsDir, filename), []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}
	}

	dbPath := filepath.Join(tmpDir, "memory.db")
	indexer, err := NewIndexer(docsDir, dbPath)
	if err != nil {
		t.Fatalf("Failed to create indexer: %v", err)
	}
	defer indexer.Close()

	// Index
	_, err = indexer.IndexDirectory()
	if err != nil {
		t.Fatalf("Index failed: %v", err)
	}

	// Get stats
	stats, err := indexer.GetStats()
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}

	if stats.TotalArtifacts != 5 {
		t.Errorf("Expected 5 artifacts, got %d", stats.TotalArtifacts)
	}
}
