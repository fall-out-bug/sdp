package graph

import (
	"errors"
	"fmt"
)

// ErrCircularDependency is returned when a circular dependency is detected
var ErrCircularDependency = errors.New("circular dependency detected")

// ErrNodeExists is returned when attempting to add a duplicate node
var ErrNodeExists = errors.New("node already exists")

// ErrMissingDependency is returned when a node depends on a non-existent node
var ErrMissingDependency = errors.New("dependency node does not exist")

// WorkstreamNode represents a workstream in the dependency graph
type WorkstreamNode struct {
	ID        string
	DependsOn []string
	Indegree  int
	Completed bool
}

// DependencyGraph represents a directed acyclic graph of workstream dependencies
type DependencyGraph struct {
	nodes map[string]*WorkstreamNode
	edges map[string][]string // node -> dependents
}

// NewDependencyGraph creates a new empty dependency graph
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		nodes: make(map[string]*WorkstreamNode),
		edges: make(map[string][]string),
	}
}

// AddNode adds a new workstream node to the graph
func (g *DependencyGraph) AddNode(id string, dependsOn []string) error {
	// Check for duplicate node
	if _, exists := g.nodes[id]; exists {
		return ErrNodeExists
	}

	// Verify all dependencies exist
	for _, dep := range dependsOn {
		if _, exists := g.nodes[dep]; !exists {
			return fmt.Errorf("%w: %s depends on non-existent %s", ErrMissingDependency, id, dep)
		}
	}

	// Create node
	node := &WorkstreamNode{
		ID:        id,
		DependsOn: dependsOn,
		Indegree:  len(dependsOn),
		Completed: false,
	}
	g.nodes[id] = node

	// Initialize edges slice
	g.edges[id] = []string{}

	// Add edges from dependencies to this node
	for _, dep := range dependsOn {
		g.edges[dep] = append(g.edges[dep], id)
	}

	return nil
}

// AddEdge adds an edge from one node to another
// This is used for testing and manual graph construction
func (g *DependencyGraph) AddEdge(from, to string) error {
	// Check if edge would create a cycle
	if g.wouldCreateCycle(from, to) {
		return ErrCircularDependency
	}

	// Add edge
	g.edges[from] = append(g.edges[from], to)
	g.nodes[to].Indegree++

	return nil
}

// wouldCreateCycle checks if adding an edge would create a cycle using DFS
func (g *DependencyGraph) wouldCreateCycle(from, to string) bool {
	visited := make(map[string]bool)
	return g.hasCycleDFS(to, from, visited)
}

// hasCycleDFS performs DFS to detect cycles
func (g *DependencyGraph) hasCycleDFS(current, target string, visited map[string]bool) bool {
	if current == target {
		return true
	}

	if visited[current] {
		return false
	}

	visited[current] = true

	// Check all neighbors
	for _, neighbor := range g.edges[current] {
		if g.hasCycleDFS(neighbor, target, visited) {
			return true
		}
	}

	return false
}

// TopologicalSort performs Kahn's algorithm to return a topological ordering
func (g *DependencyGraph) TopologicalSort() ([]string, error) {
	if len(g.nodes) == 0 {
		return []string{}, nil
	}

	// Create a copy of indegrees to avoid modifying the original
	indegrees := make(map[string]int)
	for id, node := range g.nodes {
		indegrees[id] = node.Indegree
	}

	// Initialize queue with nodes that have indegree 0
	queue := []string{}
	for id, indegree := range indegrees {
		if indegree == 0 {
			queue = append(queue, id)
		}
	}

	result := []string{}

	// Process nodes
	for len(queue) > 0 {
		// Dequeue a node
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// For each neighbor, reduce indegree
		for _, neighbor := range g.edges[current] {
			indegrees[neighbor]--
			if indegrees[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// Check for cycle (if result doesn't contain all nodes)
	if len(result) != len(g.nodes) {
		return nil, ErrCircularDependency
	}

	return result, nil
}

// GetReady returns all nodes that are ready to execute (indegree = 0 and not completed)
func (g *DependencyGraph) GetReady() []string {
	ready := []string{}

	for id, node := range g.nodes {
		// Node must have indegree 0 AND not be completed
		if node.Indegree == 0 && !node.Completed {
			ready = append(ready, id)
		}
	}

	return ready
}

// MarkComplete marks a node as completed and updates indegrees of dependent nodes
func (g *DependencyGraph) MarkComplete(id string) {
	node, exists := g.nodes[id]
	if !exists {
		return
	}

	node.Completed = true

	// Reduce indegree of all dependent nodes
	for _, dependent := range g.edges[id] {
		if depNode, exists := g.nodes[dependent]; exists {
			depNode.Indegree--
		}
	}
}
