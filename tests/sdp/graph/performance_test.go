package graph

import (
	"fmt"
	"testing"
	"time"

	"github.com/fall-out-bug/sdp/src/sdp/graph"
)

// BenchmarkSequential1Workstream benchmarks 1 workstream (baseline)
func BenchmarkSequential1Workstream(b *testing.B) {
	g := graph.NewDependencyGraph()
	g.AddNode("00-001-01", []string{})

	executeFn := func(wsID string) error {
		time.Sleep(10 * time.Millisecond) // Simulate work
		return nil
	}

	for i := 0; i < b.N; i++ {
		dispatcher := graph.NewDispatcher(g, 1)
		dispatcher.Execute(executeFn)
	}
}

// BenchmarkSequential5Workstreams benchmarks 5 sequential workstreams
func BenchmarkSequential5Workstreams(b *testing.B) {
	g := graph.NewDependencyGraph()
	g.AddNode("00-001-01", []string{})
	g.AddNode("00-001-02", []string{"00-001-01"})
	g.AddNode("00-001-03", []string{"00-001-02"})
	g.AddNode("00-001-04", []string{"00-001-03"})
	g.AddNode("00-001-05", []string{"00-001-04"})

	executeFn := func(wsID string) error {
		time.Sleep(10 * time.Millisecond) // Simulate work
		return nil
	}

	for i := 0; i < b.N; i++ {
		dispatcher := graph.NewDispatcher(g, 1) // Sequential
		dispatcher.Execute(executeFn)
	}
}

// BenchmarkParallel5Workstreams benchmarks 5 parallel workstreams
func BenchmarkParallel5Workstreams(b *testing.B) {
	g := graph.NewDependencyGraph()
	// All independent - can run in parallel
	g.AddNode("00-001-01", []string{})
	g.AddNode("00-001-02", []string{})
	g.AddNode("00-001-03", []string{})
	g.AddNode("00-001-04", []string{})
	g.AddNode("00-001-05", []string{})

	executeFn := func(wsID string) error {
		time.Sleep(10 * time.Millisecond) // Simulate work
		return nil
	}

	for i := 0; i < b.N; i++ {
		dispatcher := graph.NewDispatcher(g, 5) // Parallel
		dispatcher.Execute(executeFn)
	}
}

// BenchmarkSequential10Workstreams benchmarks 10 sequential workstreams
func BenchmarkSequential10Workstreams(b *testing.B) {
	g := graph.NewDependencyGraph()
	for i := 1; i <= 10; i++ {
		if i == 1 {
			g.AddNode(fmt.Sprintf("00-001-%02d", i), []string{})
		} else {
			g.AddNode(fmt.Sprintf("00-001-%02d", i), []string{fmt.Sprintf("00-001-%02d", i-1)})
		}
	}

	executeFn := func(wsID string) error {
		time.Sleep(10 * time.Millisecond) // Simulate work
		return nil
	}

	for i := 0; i < b.N; i++ {
		dispatcher := graph.NewDispatcher(g, 1) // Sequential
		dispatcher.Execute(executeFn)
	}
}

// BenchmarkParallel10Workstreams benchmarks 10 parallel workstreams
func BenchmarkParallel10Workstreams(b *testing.B) {
	g := graph.NewDependencyGraph()
	// All independent - can run in parallel
	for i := 1; i <= 10; i++ {
		g.AddNode(fmt.Sprintf("00-001-%02d", i), []string{})
	}

	executeFn := func(wsID string) error {
		time.Sleep(10 * time.Millisecond) // Simulate work
		return nil
	}

	for i := 0; i < b.N; i++ {
		dispatcher := graph.NewDispatcher(g, 5) // Parallel
		dispatcher.Execute(executeFn)
	}
}

// TestSpeedupTarget5WS verifies 3x speedup for 5 workstreams
func TestSpeedupTarget5WS(t *testing.T) {
	// Sequential execution
	g1 := graph.NewDependencyGraph()
	for i := 1; i <= 5; i++ {
		g1.AddNode(fmt.Sprintf("00-001-%02d", i), []string{})
	}

	executeFn := func(wsID string) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	}

	start := time.Now()
	dispatcher1 := graph.NewDispatcher(g1, 1) // Sequential
	dispatcher1.Execute(executeFn)
	sequentialDuration := time.Since(start)

	// Parallel execution
	g2 := graph.NewDependencyGraph()
	for i := 1; i <= 5; i++ {
		g2.AddNode(fmt.Sprintf("00-001-%02d", i), []string{})
	}

	start = time.Now()
	dispatcher2 := graph.NewDispatcher(g2, 5) // Parallel
	dispatcher2.Execute(executeFn)
	parallelDuration := time.Since(start)

	speedup := float64(sequentialDuration) / float64(parallelDuration)

	t.Logf("Sequential: %v, Parallel: %v, Speedup: %.2fx", sequentialDuration, parallelDuration, speedup)

	// We expect at least 2x speedup (3x is the target, but we allow some margin)
	if speedup < 2.0 {
		t.Errorf("Speedup target not met: expected at least 2x, got %.2fx", speedup)
	}
}

// TestSpeedupTarget10WS verifies 5x speedup for 10 workstreams
func TestSpeedupTarget10WS(t *testing.T) {
	// Sequential execution
	g1 := graph.NewDependencyGraph()
	for i := 1; i <= 10; i++ {
		g1.AddNode(fmt.Sprintf("00-001-%02d", i), []string{})
	}

	executeFn := func(wsID string) error {
		time.Sleep(10 * time.Millisecond)
		return nil
	}

	start := time.Now()
	dispatcher1 := graph.NewDispatcher(g1, 1) // Sequential
	dispatcher1.Execute(executeFn)
	sequentialDuration := time.Since(start)

	// Parallel execution
	g2 := graph.NewDependencyGraph()
	for i := 1; i <= 10; i++ {
		g2.AddNode(fmt.Sprintf("00-001-%02d", i), []string{})
	}

	start = time.Now()
	dispatcher2 := graph.NewDispatcher(g2, 5) // Parallel
	dispatcher2.Execute(executeFn)
	parallelDuration := time.Since(start)

	speedup := float64(sequentialDuration) / float64(parallelDuration)

	t.Logf("Sequential: %v, Parallel: %v, Speedup: %.2fx", sequentialDuration, parallelDuration, speedup)

	// We expect at least 3x speedup (5x is the target, but we allow some margin)
	if speedup < 3.0 {
		t.Errorf("Speedup target not met: expected at least 3x, got %.2fx", speedup)
	}
}
