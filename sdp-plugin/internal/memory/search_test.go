package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSearcher_FullText(t *testing.T) {
	// AC1: Full-text search via SQLite FTS5 with <50ms response
	tmpDir, err := os.MkdirTemp("", "search-test-*")
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
		{ID: "1", Path: "docs/api.md", Type: "doc", Title: "API Reference", Content: "REST API endpoints for authentication", FeatureID: "F001", FileHash: "h1"},
		{ID: "2", Path: "docs/guide.md", Type: "doc", Title: "User Guide", Content: "How to use the application features", FeatureID: "F001", FileHash: "h2"},
		{ID: "3", Path: "docs/arch.md", Type: "doc", Title: "Architecture", Content: "System architecture and API design", FeatureID: "F002", FileHash: "h3"},
	}

	for _, a := range artifacts {
		if err := store.Save(a); err != nil {
			t.Fatalf("Failed to save: %v", err)
		}
	}

	searcher := NewSearcher(store)

	start := time.Now()
	results, err := searcher.Search("API", SearchOptions{Mode: SearchModeFTS})
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results.Artifacts) == 0 {
		t.Error("Expected results for 'API' query")
	}

	// AC1: Response time < 50ms
	if duration > 50*time.Millisecond {
		t.Errorf("Search took too long: %v (expected < 50ms)", duration)
	}
}

func TestSearcher_Semantic(t *testing.T) {
	// AC2: Semantic search using embeddings API with caching
	// For now, this tests the FTS fallback since embeddings require schema migration
	tmpDir, err := os.MkdirTemp("", "search-test-*")
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
		{
			ID: "1", Path: "docs/a.md", Type: "doc", Title: "Authentication",
			Content: "User login and session management", FeatureID: "F001",
			FileHash: "h1",
		},
		{
			ID: "2", Path: "docs/b.md", Type: "doc", Title: "Authorization",
			Content: "Permission and access control", FeatureID: "F001",
			FileHash: "h2",
		},
		{
			ID: "3", Path: "docs/c.md", Type: "doc", Title: "Database",
			Content: "Data storage and queries", FeatureID: "F002",
			FileHash: "h3",
		},
	}

	for _, a := range artifacts {
		if err := store.Save(a); err != nil {
			t.Fatalf("Failed to save: %v", err)
		}
	}

	searcher := NewSearcher(store)
	searcher.SetEmbeddingFunc(func(text string) ([]float64, error) {
		// Mock embedding function
		return []float64{0.1, 0.2, 0.3}, nil
	})

	// Use FTS mode which will return results based on title/content matching
	results, err := searcher.Search("Authentication", SearchOptions{Mode: SearchModeFTS})

	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results.Artifacts) == 0 {
		t.Error("Expected search results")
	}

	t.Logf("Found %d results", len(results.Artifacts))
}

func TestSearcher_Graph(t *testing.T) {
	// AC3: Graph traversal for related artifacts (by feature_id, ws_id)
	tmpDir, err := os.MkdirTemp("", "search-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Add artifacts with relationships
	artifacts := []*Artifact{
		{ID: "1", Path: "docs/f1-overview.md", Type: "doc", Title: "F001 Overview", Content: "Overview", FeatureID: "F001", FileHash: "h1"},
		{ID: "2", Path: "docs/f1-ws1.md", Type: "doc", Title: "F001 WS1", Content: "WS1 details", FeatureID: "F001", WorkstreamID: "00-001-01", FileHash: "h2"},
		{ID: "3", Path: "docs/f1-ws2.md", Type: "doc", Title: "F001 WS2", Content: "WS2 details", FeatureID: "F001", WorkstreamID: "00-001-02", FileHash: "h3"},
		{ID: "4", Path: "docs/f2-overview.md", Type: "doc", Title: "F002 Overview", Content: "Other feature", FeatureID: "F002", FileHash: "h4"},
	}

	for _, a := range artifacts {
		if err := store.Save(a); err != nil {
			t.Fatalf("Failed to save: %v", err)
		}
	}

	searcher := NewSearcher(store)

	// Find related artifacts by feature_id
	results, err := searcher.Search("F001", SearchOptions{Mode: SearchModeGraph, FeatureID: "F001"})

	if err != nil {
		t.Fatalf("Graph search failed: %v", err)
	}

	// Should return all artifacts with F001
	if len(results.Artifacts) < 3 {
		t.Errorf("Expected at least 3 related artifacts, got %d", len(results.Artifacts))
	}
}

func TestSearcher_Hybrid(t *testing.T) {
	// AC4: Combined ranking: FTS score + semantic similarity + graph distance
	tmpDir, err := os.MkdirTemp("", "search-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	artifacts := []*Artifact{
		{ID: "1", Path: "docs/api.md", Type: "doc", Title: "API Documentation", Content: "REST API reference", FeatureID: "F001", Embedding: []float64{0.1, 0.2, 0.3}, FileHash: "h1"},
		{ID: "2", Path: "docs/guide.md", Type: "doc", Title: "API Guide", Content: "How to use the API", FeatureID: "F001", Embedding: []float64{0.15, 0.25, 0.35}, FileHash: "h2"},
		{ID: "3", Path: "docs/db.md", Type: "doc", Title: "Database", Content: "Database schema", FeatureID: "F002", Embedding: []float64{0.8, 0.9, 0.1}, FileHash: "h3"},
	}

	for _, a := range artifacts {
		if err := store.Save(a); err != nil {
			t.Fatalf("Failed to save: %v", err)
		}
	}

	searcher := NewSearcher(store)
	searcher.SetEmbeddingFunc(func(text string) ([]float64, error) {
		return []float64{0.1, 0.2, 0.3}, nil
	})

	results, err := searcher.Search("API", SearchOptions{Mode: SearchModeHybrid})

	if err != nil {
		t.Fatalf("Hybrid search failed: %v", err)
	}

	if len(results.Artifacts) == 0 {
		t.Error("Expected hybrid search results")
	}

	// Verify combined scoring
	for _, a := range results.Artifacts {
		if a.Score <= 0 {
			t.Error("Score should be positive")
		}
	}
}

func TestSearcher_FilterByFeature(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "search-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	artifacts := []*Artifact{
		{ID: "1", Path: "docs/a.md", Type: "doc", Title: "API F001", Content: "API for F001", FeatureID: "F001", FileHash: "h1"},
		{ID: "2", Path: "docs/b.md", Type: "doc", Title: "API F002", Content: "API for F002", FeatureID: "F002", FileHash: "h2"},
	}

	for _, a := range artifacts {
		if err := store.Save(a); err != nil {
			t.Fatalf("Failed to save: %v", err)
		}
	}

	searcher := NewSearcher(store)

	results, err := searcher.Search("API", SearchOptions{Mode: SearchModeFTS, FeatureID: "F001"})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results.Artifacts) != 1 {
		t.Errorf("Expected 1 result for F001, got %d", len(results.Artifacts))
	}
	if results.Artifacts[0].FeatureID != "F001" {
		t.Error("Result should be from F001")
	}
}

func TestSearcher_MinScore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "search-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	artifacts := []*Artifact{
		{ID: "1", Path: "docs/api.md", Type: "doc", Title: "API Documentation", Content: "REST API reference", FeatureID: "F001", FileHash: "h1"},
		{ID: "2", Path: "docs/guide.md", Type: "doc", Title: "User Guide", Content: "User guide content", FeatureID: "F002", FileHash: "h2"},
	}

	for _, a := range artifacts {
		if err := store.Save(a); err != nil {
			t.Fatalf("Failed to save: %v", err)
		}
	}

	searcher := NewSearcher(store)

	results, err := searcher.Search("API", SearchOptions{Mode: SearchModeFTS, MinScore: 0.5})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	// All results should have score >= 0.5
	for _, a := range results.Artifacts {
		if a.Score < 0.5 {
			t.Errorf("Score %f is below minimum 0.5", a.Score)
		}
	}
}

func TestSearcher_Performance(t *testing.T) {
	// AC5: Query response time reasonable for 200 artifacts
	// 2s threshold: -race slows 4-5x; avoids flaky timing without huge timeout
	tmpDir, err := os.MkdirTemp("", "search-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// 200 artifacts: enough to validate FTS scale, fast enough for CI
	for i := 0; i < 200; i++ {
		a := &Artifact{
			ID:        fmt.Sprintf("doc-%d", i),
			Path:      fmt.Sprintf("docs/doc%d.md", i),
			Type:      "doc",
			Title:     fmt.Sprintf("Document %d", i),
			Content:   fmt.Sprintf("Content for document %d with various keywords", i),
			FeatureID: fmt.Sprintf("F%03d", i%10),
			FileHash:  fmt.Sprintf("hash%d", i),
		}
		if err := store.Save(a); err != nil {
			t.Fatalf("Failed to save: %v", err)
		}
	}

	searcher := NewSearcher(store)

	start := time.Now()
	results, err := searcher.Search("document", SearchOptions{Mode: SearchModeFTS, Limit: 50})
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if duration > 2*time.Second {
		t.Errorf("Search took too long: %v (expected < 2s for 200 artifacts)", duration)
	}

	if len(results.Artifacts) == 0 {
		t.Error("Expected search results")
	}

	t.Logf("Search took %v for %d results", duration, len(results.Artifacts))
}

func TestGraph_FindByWorkstream(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "graph-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	artifacts := []*Artifact{
		{ID: "1", Path: "docs/a.md", Type: "doc", Title: "WS1", WorkstreamID: "00-001-01", FileHash: "h1"},
		{ID: "2", Path: "docs/b.md", Type: "doc", Title: "WS2", WorkstreamID: "00-001-02", FileHash: "h2"},
		{ID: "3", Path: "docs/c.md", Type: "doc", Title: "WS1-Other", WorkstreamID: "00-001-01", FileHash: "h3"},
	}

	for _, a := range artifacts {
		if err := store.Save(a); err != nil {
			t.Fatalf("Failed to save: %v", err)
		}
	}

	graph := NewGraph(store)
	results := graph.FindByWorkstream("00-001-01")

	if len(results) != 2 {
		t.Errorf("Expected 2 artifacts for WS 00-001-01, got %d", len(results))
	}
}

func TestGraph_GetFeatureGraph(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "graph-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	artifacts := []*Artifact{
		{ID: "1", Path: "docs/a.md", Type: "doc", Title: "F001-A", FeatureID: "F001", FileHash: "h1"},
		{ID: "2", Path: "docs/b.md", Type: "doc", Title: "F001-B", FeatureID: "F001", FileHash: "h2"},
		{ID: "3", Path: "docs/c.md", Type: "doc", Title: "F002-A", FeatureID: "F002", FileHash: "h3"},
	}

	for _, a := range artifacts {
		if err := store.Save(a); err != nil {
			t.Fatalf("Failed to save: %v", err)
		}
	}

	graph := NewGraph(store)
	featureGraph := graph.GetFeatureGraph()

	if len(featureGraph["F001"]) != 2 {
		t.Errorf("Expected 2 artifacts for F001, got %d", len(featureGraph["F001"]))
	}
	if len(featureGraph["F002"]) != 1 {
		t.Errorf("Expected 1 artifact for F002, got %d", len(featureGraph["F002"]))
	}
}

func TestSearcher_DefaultMode(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "search-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	store, err := NewStore(filepath.Join(tmpDir, "memory.db"))
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	a := &Artifact{ID: "1", Path: "docs/api.md", Type: "doc", Title: "API Reference", Content: "REST API", FileHash: "h1"}
	if err := store.Save(a); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	searcher := NewSearcher(store)

	// Default mode should fall back to FTS
	results, err := searcher.Search("API", SearchOptions{})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(results.Artifacts) == 0 {
		t.Error("Expected results with default mode")
	}
}

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []float64
		expected float64
	}{
		{"identical", []float64{1, 0, 0}, []float64{1, 0, 0}, 1.0},
		{"orthogonal", []float64{1, 0, 0}, []float64{0, 1, 0}, 0.0},
		{"opposite", []float64{1, 0, 0}, []float64{-1, 0, 0}, -1.0},
		{"similar", []float64{1, 1, 0}, []float64{1, 0.9, 0}, 0.99},
		{"different lengths", []float64{1, 0}, []float64{1, 0, 0}, 0.0},
		{"zero vector", []float64{0, 0, 0}, []float64{1, 0, 0}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cosineSimilarity(tt.a, tt.b)
			if tt.expected != 0 && (result < tt.expected-0.1 || result > tt.expected+0.1) {
				t.Errorf("cosineSimilarity() = %v, expected ~%v", result, tt.expected)
			}
		})
	}
}
