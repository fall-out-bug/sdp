package graph

import (
	"errors"
	"testing"
)

// TestNewDependencyGraph verifies graph creation
func TestNewDependencyGraph(t *testing.T) {
	g := NewDependencyGraph()
	if g == nil {
		t.Fatal("NewDependencyGraph returned nil")
	}
	if g.nodes == nil {
		t.Error("nodes map not initialized")
	}
	if g.edges == nil {
		t.Error("edges map not initialized")
	}
}

// TestAddNode_SingleNode verifies adding a single node
func TestAddNode_SingleNode(t *testing.T) {
	g := NewDependencyGraph()
	err := g.AddNode("ws-001", nil)
	if err != nil {
		t.Errorf("AddNode failed: %v", err)
	}
	if len(g.nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(g.nodes))
	}
	if g.nodes["ws-001"].Indegree != 0 {
		t.Errorf("expected indegree 0, got %d", g.nodes["ws-001"].Indegree)
	}
}

// TestAddNode_DuplicateNode verifies duplicate node detection
func TestAddNode_DuplicateNode(t *testing.T) {
	g := NewDependencyGraph()
	g.AddNode("ws-001", nil)
	err := g.AddNode("ws-001", nil)
	if err == nil {
		t.Error("expected error for duplicate node")
	}
	if !errors.Is(err, ErrNodeExists) {
		t.Errorf("expected ErrNodeExists, got %v", err)
	}
}

// TestAddNode_WithDependencies verifies adding node with dependencies
func TestAddNode_WithDependencies(t *testing.T) {
	g := NewDependencyGraph()
	g.AddNode("ws-001", nil)
	g.AddNode("ws-002", nil)

	err := g.AddNode("ws-003", []string{"ws-001", "ws-002"})
	if err != nil {
		t.Errorf("AddNode failed: %v", err)
	}
	if g.nodes["ws-003"].Indegree != 2 {
		t.Errorf("expected indegree 2, got %d", g.nodes["ws-003"].Indegree)
	}
	// Check edges
	if len(g.edges["ws-001"]) != 1 || g.edges["ws-001"][0] != "ws-003" {
		t.Error("edge from ws-001 to ws-003 not created")
	}
}

// TestAddNode_MissingDependency verifies missing dependency detection
func TestAddNode_MissingDependency(t *testing.T) {
	g := NewDependencyGraph()
	err := g.AddNode("ws-002", []string{"ws-001"})
	if err == nil {
		t.Error("expected error for missing dependency")
	}
	if !errors.Is(err, ErrMissingDependency) {
		t.Errorf("expected ErrMissingDependency, got %v", err)
	}
}

// TestAddEdge verifies edge addition
func TestAddEdge(t *testing.T) {
	g := NewDependencyGraph()
	g.AddNode("ws-001", nil)
	g.AddNode("ws-002", nil)

	err := g.AddEdge("ws-001", "ws-002")
	if err != nil {
		t.Errorf("AddEdge failed: %v", err)
	}
	if g.nodes["ws-002"].Indegree != 1 {
		t.Errorf("expected indegree 1, got %d", g.nodes["ws-002"].Indegree)
	}
}

// TestAddEdge_CreatesCycle verifies cycle detection
func TestAddEdge_CreatesCycle(t *testing.T) {
	g := NewDependencyGraph()
	g.AddNode("a", nil)
	g.AddNode("b", nil)
	g.AddNode("c", nil)

	// Create chain: a -> b -> c
	g.AddEdge("a", "b")
	g.AddEdge("b", "c")

	// Try to create cycle: c -> a
	err := g.AddEdge("c", "a")
	if err == nil {
		t.Error("expected error for cycle creation")
	}
	if !errors.Is(err, ErrCircularDependency) {
		t.Errorf("expected ErrCircularDependency, got %v", err)
	}
}

// TestTopologicalSort_EmptyGraph verifies empty graph sort
func TestTopologicalSort_EmptyGraph(t *testing.T) {
	g := NewDependencyGraph()
	result, err := g.TopologicalSort()
	if err != nil {
		t.Errorf("TopologicalSort failed: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

// TestTopologicalSort_SingleNode verifies single node sort
func TestTopologicalSort_SingleNode(t *testing.T) {
	g := NewDependencyGraph()
	g.AddNode("a", nil)

	result, err := g.TopologicalSort()
	if err != nil {
		t.Errorf("TopologicalSort failed: %v", err)
	}
	if len(result) != 1 || result[0] != "a" {
		t.Errorf("expected [a], got %v", result)
	}
}

// TestTopologicalSort_LinearChain verifies linear chain sort
func TestTopologicalSort_LinearChain(t *testing.T) {
	// Graph: a -> b -> c
	g := NewDependencyGraph()
	g.AddNode("a", nil)
	g.AddNode("b", []string{"a"})
	g.AddNode("c", []string{"b"})

	result, err := g.TopologicalSort()
	if err != nil {
		t.Errorf("TopologicalSort failed: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(result))
	}
	// Verify order: a must come before b, b before c
	orderMap := make(map[string]int)
	for i, id := range result {
		orderMap[id] = i
	}
	if orderMap["a"] > orderMap["b"] {
		t.Error("a should come before b")
	}
	if orderMap["b"] > orderMap["c"] {
		t.Error("b should come before c")
	}
}

// TestTopologicalSort_DiamondDependency verifies diamond dependency sort
func TestTopologicalSort_DiamondDependency(t *testing.T) {
	// Graph:
	//       a
	//      / \
	//     b   c
	//      \ /
	//       d
	g := NewDependencyGraph()
	g.AddNode("a", nil)
	g.AddNode("b", []string{"a"})
	g.AddNode("c", []string{"a"})
	g.AddNode("d", []string{"b", "c"})

	result, err := g.TopologicalSort()
	if err != nil {
		t.Errorf("TopologicalSort failed: %v", err)
	}
	if len(result) != 4 {
		t.Errorf("expected 4 nodes, got %d", len(result))
	}
	// Verify order constraints
	orderMap := make(map[string]int)
	for i, id := range result {
		orderMap[id] = i
	}
	if orderMap["a"] > orderMap["b"] || orderMap["a"] > orderMap["c"] {
		t.Error("a should come before b and c")
	}
	if orderMap["b"] > orderMap["d"] || orderMap["c"] > orderMap["d"] {
		t.Error("b and c should come before d")
	}
}

// TestGetReady verifies getting ready nodes
func TestGetReady(t *testing.T) {
	g := NewDependencyGraph()
	g.AddNode("a", nil)
	g.AddNode("b", []string{"a"})
	g.AddNode("c", nil)

	ready := g.GetReady()
	if len(ready) != 2 {
		t.Errorf("expected 2 ready nodes, got %d", len(ready))
	}
	readyMap := make(map[string]bool)
	for _, id := range ready {
		readyMap[id] = true
	}
	if !readyMap["a"] || !readyMap["c"] {
		t.Error("expected a and c to be ready")
	}
}

// TestGetReady_AfterCompletion verifies ready nodes after completion
func TestGetReady_AfterCompletion(t *testing.T) {
	g := NewDependencyGraph()
	g.AddNode("a", nil)
	g.AddNode("b", []string{"a"})

	// Initially only 'a' is ready
	ready := g.GetReady()
	if len(ready) != 1 || ready[0] != "a" {
		t.Errorf("expected [a] ready, got %v", ready)
	}

	// Complete 'a'
	g.MarkComplete("a")

	// Now 'b' should be ready
	ready = g.GetReady()
	if len(ready) != 1 || ready[0] != "b" {
		t.Errorf("expected [b] ready after a completion, got %v", ready)
	}
}

// TestMarkComplete verifies marking nodes complete
func TestMarkComplete(t *testing.T) {
	g := NewDependencyGraph()
	g.AddNode("a", nil)
	g.AddNode("b", []string{"a"})

	g.MarkComplete("a")

	if !g.nodes["a"].Completed {
		t.Error("node a should be marked complete")
	}
	if g.nodes["b"].Indegree != 0 {
		t.Errorf("node b should have indegree 0, got %d", g.nodes["b"].Indegree)
	}
}

// TestMarkComplete_NonExistent verifies marking non-existent node
func TestMarkComplete_NonExistent(t *testing.T) {
	g := NewDependencyGraph()
	// Should not panic
	g.MarkComplete("nonexistent")
}

// TestWorkstreamNode_Fields verifies WorkstreamNode struct
func TestWorkstreamNode_Fields(t *testing.T) {
	node := &WorkstreamNode{
		ID:        "ws-001",
		DependsOn: []string{"ws-000"},
		Indegree:  1,
		Completed: false,
	}

	if node.ID != "ws-001" {
		t.Error("ID not set correctly")
	}
	if len(node.DependsOn) != 1 {
		t.Error("DependsOn not set correctly")
	}
	if node.Indegree != 1 {
		t.Error("Indegree not set correctly")
	}
	if node.Completed {
		t.Error("Completed should be false")
	}
}
