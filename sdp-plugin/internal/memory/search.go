package memory

import (
	"sort"
	"time"
)

// Search performs search using the specified mode
func (s *Searcher) Search(query string, opts SearchOptions) (*SearchResult, error) {
	start := time.Now()

	if opts.Limit == 0 {
		opts.Limit = 50
	}

	var results []ScoredArtifact
	var err error

	switch opts.Mode {
	case SearchModeFTS:
		results, err = s.fullTextSearch(query, opts)
	case SearchModeSemantic:
		results, err = s.semanticSearch(query, opts)
	case SearchModeGraph:
		results, err = s.graphSearch(query, opts)
	case SearchModeHybrid:
		results, err = s.hybridSearch(query, opts)
	default:
		results, err = s.fullTextSearch(query, opts)
	}

	if err != nil {
		return nil, err
	}

	// Filter by min score
	if opts.MinScore > 0 {
		filtered := make([]ScoredArtifact, 0, len(results))
		for _, r := range results {
			if r.Score >= opts.MinScore {
				filtered = append(filtered, r)
			}
		}
		results = filtered
	}

	// Sort by score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Apply limit
	total := len(results)
	if len(results) > opts.Limit {
		results = results[:opts.Limit]
	}

	return &SearchResult{
		Artifacts: results,
		Total:     total,
		Duration:  time.Since(start),
	}, nil
}

// fullTextSearch performs FTS5 full-text search (AC1)
func (s *Searcher) fullTextSearch(query string, opts SearchOptions) ([]ScoredArtifact, error) {
	artifacts, err := s.store.Search(query)
	if err != nil {
		return nil, err
	}

	var results []ScoredArtifact
	for _, a := range artifacts {
		if opts.FeatureID != "" && a.FeatureID != opts.FeatureID {
			continue
		}
		score := s.calculateFTSScore(query, a)
		results = append(results, ScoredArtifact{Artifact: a, Score: score})
	}

	return results, nil
}

// calculateFTSScore calculates relevance score for FTS results
func (s *Searcher) calculateFTSScore(query string, a *Artifact) float64 {
	score := 0.0
	if containsMatch(a.Title, query) {
		score += 0.5
	}
	if containsMatch(a.Content, query) {
		score += 0.3
	}
	if a.FeatureID == query {
		score += 0.2
	}
	if score > 1.0 {
		score = 1.0
	}
	return score
}
