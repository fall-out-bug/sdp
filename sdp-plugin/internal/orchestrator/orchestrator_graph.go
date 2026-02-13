package orchestrator

import (
	"fmt"
)

// BuildDependencyGraph builds a dependency graph from workstreams
func BuildDependencyGraph(workstreams []WorkstreamNode) (Graph, error) {
	graph := make(Graph)

	// Create nodes
	for _, ws := range workstreams {
		graph[ws.ID] = &DependencyNode{
			Workstream: ws,
			InDegree:   0,
		}
	}

	// Add edges
	for _, ws := range workstreams {
		node := graph[ws.ID]
		for _, depID := range ws.Dependencies {
			// Check for self-dependency
			if depID == ws.ID {
				return nil, fmt.Errorf("%w: workstream %s depends on itself",
					ErrCircularDependency, ws.ID)
			}

			// Check if dependency exists
			depNode, exists := graph[depID]
			if !exists {
				return nil, fmt.Errorf("%w: workstream %s depends on non-existent %s",
					ErrMissingDependency, ws.ID, depID)
			}

			// Add edge
			depNode.Dependencies = append(depNode.Dependencies, node)
			node.InDegree++
		}
	}

	// Check for cycles
	if hasCycle(graph) {
		return nil, ErrCircularDependency
	}

	return graph, nil
}

// TopologicalSort performs topological sort on the dependency graph
func TopologicalSort(graph Graph) ([]string, error) {
	if len(graph) == 0 {
		return []string{}, nil
	}

	// Kahn's algorithm
	inDegree := make(map[string]int)
	queue := []string{}

	// Initialize in-degree copy and find nodes with no dependencies
	for id, node := range graph {
		inDegree[id] = node.InDegree
		if node.InDegree == 0 {
			queue = append(queue, id)
		}
	}

	result := []string{}

	for len(queue) > 0 {
		// Dequeue
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// Reduce in-degree for all neighbors
		node := graph[current]
		for _, neighbor := range node.Dependencies {
			inDegree[neighbor.Workstream.ID]--
			if inDegree[neighbor.Workstream.ID] == 0 {
				queue = append(queue, neighbor.Workstream.ID)
			}
		}
	}

	// Check if all nodes were processed (cycle detection)
	if len(result) != len(graph) {
		return nil, ErrCircularDependency
	}

	return result, nil
}

// hasCycle checks if the graph has a cycle using DFS
func hasCycle(graph Graph) bool {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for id := range graph {
		if !visited[id] {
			if hasCycleDFS(graph, id, visited, recStack) {
				return true
			}
		}
	}

	return false
}

// hasCycleDFS performs DFS to detect cycles
func hasCycleDFS(graph Graph, nodeID string, visited, recStack map[string]bool) bool {
	visited[nodeID] = true
	recStack[nodeID] = true

	node := graph[nodeID]
	for _, neighbor := range node.Dependencies {
		neighborID := neighbor.Workstream.ID
		if !visited[neighborID] {
			if hasCycleDFS(graph, neighborID, visited, recStack) {
				return true
			}
		} else if recStack[neighborID] {
			return true
		}
	}

	recStack[nodeID] = false
	return false
}
