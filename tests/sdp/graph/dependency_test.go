package graph

import (
	"testing"

	"github.com/fall-out-bug/sdp/src/sdp/graph"
)

// TestEmptyGraph tests creating an empty dependency graph
func TestEmptyGraph(t *testing.T) {
	g := graph.NewDependencyGraph()

	if g == nil {
		t.Fatal("Expected non-nil graph")
	}

	// Can't access private fields directly, but we can test behavior
	ready := g.GetReady()
	if len(ready) != 0 {
		t.Errorf("Expected empty ready list, got %d nodes", len(ready))
	}
}

// TestAddNode tests adding a single node to the graph
func TestAddNode(t *testing.T) {
	g := graph.NewDependencyGraph()

	err := g.AddNode("00-001-01", []string{})
	if err != nil {
		t.Fatalf("Failed to add node: %v", err)
	}

	ready := g.GetReady()
	if len(ready) != 1 {
		t.Errorf("Expected 1 ready node, got %d", len(ready))
	}

	if ready[0] != "00-001-01" {
		t.Errorf("Expected 00-001-01 to be ready, got %s", ready[0])
	}
}

// TestAddNodeWithDependencies tests adding a node with dependencies
func TestAddNodeWithDependencies(t *testing.T) {
	g := graph.NewDependencyGraph()

	// Add dependency first
	g.AddNode("00-001-01", []string{})

	// Add dependent node
	err := g.AddNode("00-001-02", []string{"00-001-01"})
	if err != nil {
		t.Fatalf("Failed to add node with dependencies: %v", err)
	}

	ready := g.GetReady()
	if len(ready) != 1 {
		t.Errorf("Expected 1 ready node, got %d", len(ready))
	}

	if ready[0] != "00-001-01" {
		t.Errorf("Expected 00-001-01 to be ready, got %s", ready[0])
	}
}

// TestCircularDependency tests detection of circular dependencies
func TestCircularDependency(t *testing.T) {
	g := graph.NewDependencyGraph()

	// Create a simple cycle: A -> B -> C -> A
	g.AddNode("00-001-01", []string{})
	g.AddNode("00-001-02", []string{})
	g.AddNode("00-001-03", []string{})

	// Add edges to create cycle
	g.AddEdge("00-001-02", "00-001-01")
	g.AddEdge("00-001-03", "00-001-02")
	err := g.AddEdge("00-001-01", "00-001-03")

	if err == nil {
		t.Error("Expected error for circular dependency, got nil")
	}

	if err != graph.ErrCircularDependency {
		t.Errorf("Expected ErrCircularDependency, got %v", err)
	}
}

// TestTopologicalSortSimple tests topological sort with a simple DAG
func TestTopologicalSortSimple(t *testing.T) {
	g := graph.NewDependencyGraph()

	// A -> B -> C
	g.AddNode("00-001-01", []string{})
	g.AddNode("00-001-02", []string{"00-001-01"})
	g.AddNode("00-001-03", []string{"00-001-02"})

	order, err := g.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}

	if len(order) != 3 {
		t.Fatalf("Expected 3 nodes, got %d", len(order))
	}

	// Verify order: A must come before B, B must come before C
	orderMap := make(map[string]int)
	for i, id := range order {
		orderMap[id] = i
	}

	if orderMap["00-001-01"] >= orderMap["00-001-02"] {
		t.Error("00-001-01 should come before 00-001-02")
	}

	if orderMap["00-001-02"] >= orderMap["00-001-03"] {
		t.Error("00-001-02 should come before 00-001-03")
	}
}

// TestTopologicalSortComplex tests topological sort with a complex DAG
func TestTopologicalSortComplex(t *testing.T) {
	g := graph.NewDependencyGraph()

	//     A
	//    / \
	//   B   C
	//   |   |
	//   D   E
	//    \ /
	//     F
	g.AddNode("00-001-01", []string{}) // A
	g.AddNode("00-001-02", []string{"00-001-01"}) // B depends on A
	g.AddNode("00-001-03", []string{"00-001-01"}) // C depends on A
	g.AddNode("00-001-04", []string{"00-001-02"}) // D depends on B
	g.AddNode("00-001-05", []string{"00-001-03"}) // E depends on C
	g.AddNode("00-001-06", []string{"00-001-04", "00-001-05"}) // F depends on D and E

	order, err := g.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}

	if len(order) != 6 {
		t.Fatalf("Expected 6 nodes, got %d", len(order))
	}

	// Verify ordering constraints
	orderMap := make(map[string]int)
	for i, id := range order {
		orderMap[id] = i
	}

	// A before B, C
	if orderMap["00-001-01"] >= orderMap["00-001-02"] ||
		orderMap["00-001-01"] >= orderMap["00-001-03"] {
		t.Error("A should come before B and C")
	}

	// B before D
	if orderMap["00-001-02"] >= orderMap["00-001-04"] {
		t.Error("B should come before D")
	}

	// C before E
	if orderMap["00-001-03"] >= orderMap["00-001-05"] {
		t.Error("C should come before E")
	}

	// D and E before F
	if orderMap["00-001-04"] >= orderMap["00-001-06"] ||
		orderMap["00-001-05"] >= orderMap["00-001-06"] {
		t.Error("D and E should come before F")
	}
}

// TestGetReady tests identification of ready nodes
func TestGetReady(t *testing.T) {
	g := graph.NewDependencyGraph()

	//     A
	//    / \
	//   B   C
	g.AddNode("00-001-01", []string{})
	g.AddNode("00-001-02", []string{"00-001-01"})
	g.AddNode("00-001-03", []string{"00-001-01"})

	// Initially only A should be ready
	ready := g.GetReady()
	if len(ready) != 1 {
		t.Errorf("Expected 1 ready node, got %d", len(ready))
	}

	if ready[0] != "00-001-01" {
		t.Errorf("Expected 00-001-01 to be ready, got %s", ready[0])
	}
}

// TestMarkComplete tests marking nodes as complete
func TestMarkComplete(t *testing.T) {
	g := graph.NewDependencyGraph()

	// A -> B -> C
	g.AddNode("00-001-01", []string{})
	g.AddNode("00-001-02", []string{"00-001-01"})
	g.AddNode("00-001-03", []string{"00-001-02"})

	// Initially only A is ready
	ready := g.GetReady()
	if len(ready) != 1 || ready[0] != "00-001-01" {
		t.Fatalf("Expected only 00-001-01 to be ready initially")
	}

	// Mark A as complete
	g.MarkComplete("00-001-01")

	// Now B should be ready
	ready = g.GetReady()
	if len(ready) != 1 {
		t.Errorf("Expected 1 ready node after marking A complete, got %d", len(ready))
	}

	if len(ready) > 0 && ready[0] != "00-001-02" {
		t.Errorf("Expected 00-001-02 to be ready, got %s", ready[0])
	}
}

// TestMultipleReadyNodes tests multiple nodes ready simultaneously
func TestMultipleReadyNodes(t *testing.T) {
	g := graph.NewDependencyGraph()

	//     A
	//    /|\
	//   B C D
	g.AddNode("00-001-01", []string{})
	g.AddNode("00-001-02", []string{"00-001-01"})
	g.AddNode("00-001-03", []string{"00-001-01"})
	g.AddNode("00-001-04", []string{"00-001-01"})

	// Initially only A is ready
	ready := g.GetReady()
	if len(ready) != 1 {
		t.Fatalf("Expected 1 ready node initially, got %d", len(ready))
	}

	// Mark A as complete
	g.MarkComplete("00-001-01")

	// Now B, C, D should all be ready
	ready = g.GetReady()
	if len(ready) != 3 {
		t.Errorf("Expected 3 ready nodes, got %d", len(ready))
	}

	// Verify all expected nodes are present
	readyMap := make(map[string]bool)
	for _, id := range ready {
		readyMap[id] = true
	}

	for _, expected := range []string{"00-001-02", "00-001-03", "00-001-04"} {
		if !readyMap[expected] {
			t.Errorf("Expected %s to be ready", expected)
		}
	}
}

// TestIndependentNodes tests completely independent nodes
func TestIndependentNodes(t *testing.T) {
	g := graph.NewDependencyGraph()

	g.AddNode("00-001-01", []string{})
	g.AddNode("00-001-02", []string{})
	g.AddNode("00-001-03", []string{})

	ready := g.GetReady()
	if len(ready) != 3 {
		t.Errorf("Expected 3 ready nodes, got %d", len(ready))
	}

	// All should be in ready list
	readyMap := make(map[string]bool)
	for _, id := range ready {
		readyMap[id] = true
	}

	for _, expected := range []string{"00-001-01", "00-001-02", "00-001-03"} {
		if !readyMap[expected] {
			t.Errorf("Expected %s to be ready", expected)
		}
	}
}

// TestGetReadyAfterCompletion tests ready nodes after partial completion
func TestGetReadyAfterCompletion(t *testing.T) {
	g := graph.NewDependencyGraph()

	// A -> B -> C
	//  \-> D
	g.AddNode("00-001-01", []string{})
	g.AddNode("00-001-02", []string{"00-001-01"})
	g.AddNode("00-001-03", []string{"00-001-02"})
	g.AddNode("00-001-04", []string{"00-001-01"})

	// Initial: only A ready
	ready := g.GetReady()
	if len(ready) != 1 {
		t.Fatalf("Expected 1 ready node initially, got %d", len(ready))
	}

	// Complete A
	g.MarkComplete("00-001-01")

	// Now B and D should be ready
	ready = g.GetReady()
	if len(ready) != 2 {
		t.Errorf("Expected 2 ready nodes after A, got %d", len(ready))
	}

	readyMap := make(map[string]bool)
	for _, id := range ready {
		readyMap[id] = true
	}

	if !readyMap["00-001-02"] || !readyMap["00-001-04"] {
		t.Error("Expected 00-001-02 and 00-001-04 to be ready")
	}

	// Complete B
	g.MarkComplete("00-001-02")

	// Now C should be ready (D is still ready since we didn't mark it complete)
	ready = g.GetReady()
	if len(ready) != 2 {
		t.Errorf("Expected 2 ready nodes after B, got %d", len(ready))
	}

	readyMap = make(map[string]bool)
	for _, id := range ready {
		readyMap[id] = true
	}

	// Both C and D should be ready
	if !readyMap["00-001-03"] || !readyMap["00-001-04"] {
		t.Error("Expected both 00-001-03 and 00-001-04 to be ready")
	}
}

// TestDuplicateNodeAdd tests adding the same node twice
func TestDuplicateNodeAdd(t *testing.T) {
	g := graph.NewDependencyGraph()

	err := g.AddNode("00-001-01", []string{})
	if err != nil {
		t.Fatalf("Failed to add node: %v", err)
	}

	err = g.AddNode("00-001-01", []string{})
	if err == nil {
		t.Error("Expected error when adding duplicate node, got nil")
	}
}

// TestMissingDependency tests adding a node with non-existent dependency
func TestMissingDependency(t *testing.T) {
	g := graph.NewDependencyGraph()

	// Try to add node with dependency that doesn't exist
	err := g.AddNode("00-001-02", []string{"00-001-01"})
	if err == nil {
		t.Error("Expected error when adding node with missing dependency, got nil")
	}
}
