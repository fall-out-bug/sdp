package graph

import "time"

// CreateCheckpoint creates a checkpoint from the current dispatcher state
func (cm *CheckpointManager) CreateCheckpoint(graph *DependencyGraph, featureID string, completed []string, failed []string) *Checkpoint {
	// Snapshot graph nodes
	nodes := make([]NodeSnapshot, 0, len(graph.nodes))
	for _, node := range graph.nodes {
		nodes = append(nodes, NodeSnapshot{
			ID:        node.ID,
			DependsOn: copyStringSlice(node.DependsOn),
			Indegree:  node.Indegree,
			Completed: node.Completed,
		})
	}

	// Snapshot graph edges (copy to avoid mutation)
	edges := make(map[string][]string)
	for from, toList := range graph.edges {
		edges[from] = copyStringSlice(toList)
	}

	return &Checkpoint{
		Version:   "1.0",
		FeatureID: featureID,
		Timestamp: time.Now().UTC(),
		Completed: copyStringSlice(completed),
		Failed:    copyStringSlice(failed),
		Graph: &GraphSnapshot{
			Nodes: nodes,
			Edges: edges,
		},
		CircuitBreaker: nil, // Will be populated by dispatcher
	}
}

// RestoreGraph recreates a dependency graph from a checkpoint
func (cm *CheckpointManager) RestoreGraph(checkpoint *Checkpoint) *DependencyGraph {
	graph := NewDependencyGraph()

	// Restore nodes
	for _, nodeSnapshot := range checkpoint.Graph.Nodes {
		node := &WorkstreamNode{
			ID:        nodeSnapshot.ID,
			DependsOn: copyStringSlice(nodeSnapshot.DependsOn),
			Indegree:  nodeSnapshot.Indegree,
			Completed: nodeSnapshot.Completed,
		}
		graph.nodes[nodeSnapshot.ID] = node
	}

	// Restore edges
	for from, toList := range checkpoint.Graph.Edges {
		graph.edges[from] = copyStringSlice(toList)
	}

	return graph
}

// copyStringSlice creates a deep copy of a string slice
func copyStringSlice(slice []string) []string {
	if slice == nil {
		return nil
	}
	copied := make([]string, len(slice))
	copy(copied, slice)
	return copied
}
