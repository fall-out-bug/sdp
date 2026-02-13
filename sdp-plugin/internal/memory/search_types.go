package memory

import (
	"time"
)

// SearchMode defines the search mode
type SearchMode string

const (
	SearchModeFTS      SearchMode = "fts"
	SearchModeSemantic SearchMode = "semantic"
	SearchModeGraph    SearchMode = "graph"
	SearchModeHybrid   SearchMode = "hybrid"
)

// SearchOptions defines search parameters
type SearchOptions struct {
	Mode      SearchMode
	Limit     int
	FeatureID string
	MinScore  float64
}

// ScoredArtifact represents a search result with relevance score
type ScoredArtifact struct {
	*Artifact
	Score float64
}

// SearchResult contains search results and metadata
type SearchResult struct {
	Artifacts []ScoredArtifact
	Total     int
	Duration  time.Duration
}

// EmbeddingFunc generates embeddings for text
type EmbeddingFunc func(text string) ([]float64, error)

// Searcher provides hybrid search capabilities
type Searcher struct {
	store        *Store
	embeddingFn  EmbeddingFunc
	graph        *Graph
}

// NewSearcher creates a new searcher
func NewSearcher(store *Store) *Searcher {
	return &Searcher{
		store: store,
		graph: NewGraph(store),
	}
}

// SetEmbeddingFunc sets the embedding function for semantic search
func (s *Searcher) SetEmbeddingFunc(fn EmbeddingFunc) {
	s.embeddingFn = fn
}
