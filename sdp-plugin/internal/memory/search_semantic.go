package memory

// semanticSearch performs semantic similarity search (AC2)
func (s *Searcher) semanticSearch(query string, opts SearchOptions) ([]ScoredArtifact, error) {
	if s.embeddingFn == nil {
		return s.fullTextSearch(query, opts)
	}

	queryEmbedding, err := s.embeddingFn(query)
	if err != nil {
		return s.fullTextSearch(query, opts)
	}

	artifacts, err := s.store.ListAll()
	if err != nil {
		return nil, err
	}

	var results []ScoredArtifact
	for _, a := range artifacts {
		if opts.FeatureID != "" && a.FeatureID != opts.FeatureID {
			continue
		}
		if len(a.Embedding) == 0 {
			continue
		}
		similarity := cosineSimilarity(queryEmbedding, a.Embedding)
		results = append(results, ScoredArtifact{Artifact: a, Score: similarity})
	}

	return results, nil
}

// graphSearch performs graph traversal for related artifacts (AC3)
func (s *Searcher) graphSearch(query string, opts SearchOptions) ([]ScoredArtifact, error) {
	if opts.FeatureID == "" {
		return s.fullTextSearch(query, opts)
	}

	related := s.graph.FindRelated(opts.FeatureID, 2)

	var results []ScoredArtifact
	for i, a := range related {
		score := 1.0 / float64(i+1)
		results = append(results, ScoredArtifact{Artifact: a, Score: score})
	}

	return results, nil
}

// hybridSearch combines FTS, semantic, and graph search (AC4)
func (s *Searcher) hybridSearch(query string, opts SearchOptions) ([]ScoredArtifact, error) {
	ftsResults, _ := s.fullTextSearch(query, opts)
	semResults, _ := s.semanticSearch(query, opts)
	var graphResults []ScoredArtifact
	if opts.FeatureID != "" {
		graphResults, _ = s.graphSearch(query, opts)
	}

	scoreMap := make(map[string]*ScoredArtifact)

	for _, r := range ftsResults {
		r.Score = r.Score * 0.4
		scoreMap[r.ID] = &r
	}

	for _, r := range semResults {
		if existing, ok := scoreMap[r.ID]; ok {
			existing.Score += r.Score * 0.4
		} else {
			r.Score = r.Score * 0.4
			scoreMap[r.ID] = &r
		}
	}

	for _, r := range graphResults {
		if existing, ok := scoreMap[r.ID]; ok {
			existing.Score += r.Score * 0.2
		} else {
			r.Score = r.Score * 0.2
			scoreMap[r.ID] = &r
		}
	}

	results := make([]ScoredArtifact, 0, len(scoreMap))
	for _, r := range scoreMap {
		results = append(results, *r)
	}

	return results, nil
}
