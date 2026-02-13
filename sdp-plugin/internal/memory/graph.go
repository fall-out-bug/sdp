package memory

import "github.com/fall-out-bug/sdp/internal/safetylog"

// Graph provides graph traversal for related artifacts (AC3)
type Graph struct {
	store *Store
}

// NewGraph creates a new graph for artifact relationships
func NewGraph(store *Store) *Graph {
	return &Graph{store: store}
}

// FindRelated finds artifacts related to the given feature_id
func (g *Graph) FindRelated(featureID string, maxDepth int) []*Artifact {
	if maxDepth <= 0 {
		maxDepth = 1
	}

	artifacts, err := g.store.ListAll()
	if err != nil {
		safetylog.Warn("graph: FindRelated failed to list artifacts: %v", err)
		return nil
	}

	var related []*Artifact
	visited := make(map[string]bool)

	// Find artifacts by feature_id
	for _, a := range artifacts {
		if a.FeatureID == featureID {
			if !visited[a.ID] {
				visited[a.ID] = true
				related = append(related, a)
			}
		}
	}

	// If depth > 1, find related by workstream relationships
	if maxDepth > 1 {
		for _, a := range related {
			if a.WorkstreamID != "" {
				// Find other workstreams in same feature
				for _, other := range artifacts {
					if other.FeatureID == featureID && !visited[other.ID] {
						visited[other.ID] = true
						related = append(related, other)
					}
				}
			}
		}
	}

	return related
}

// FindByWorkstream finds artifacts by workstream_id
func (g *Graph) FindByWorkstream(wsID string) []*Artifact {
	artifacts, err := g.store.ListAll()
	if err != nil {
		safetylog.Warn("graph: FindByWorkstream failed to list artifacts: %v", err)
		return nil
	}

	var results []*Artifact
	for _, a := range artifacts {
		if a.WorkstreamID == wsID {
			results = append(results, a)
		}
	}
	return results
}

// GetFeatureGraph returns all artifacts organized by feature
func (g *Graph) GetFeatureGraph() map[string][]*Artifact {
	artifacts, err := g.store.ListAll()
	if err != nil {
		safetylog.Warn("graph: GetFeatureGraph failed to list artifacts: %v", err)
		return nil
	}

	graph := make(map[string][]*Artifact)
	for _, a := range artifacts {
		if a.FeatureID != "" {
			graph[a.FeatureID] = append(graph[a.FeatureID], a)
		}
	}
	return graph
}
