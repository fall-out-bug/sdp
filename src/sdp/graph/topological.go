package graph

// TopologicalSort performs Kahn's algorithm for topological sorting
// This file is a wrapper for the algorithm implemented in dependency.go
//
// Kahn's Algorithm:
// 1. Identify all nodes with indegree 0 (no dependencies)
// 2. Add them to a queue
// 3. While queue is not empty:
//    a. Remove a node from the queue
//    b. Add it to the result
//    c. For each neighbor (dependent), reduce indegree by 1
//    d. If neighbor's indegree becomes 0, add it to the queue
// 4. If result contains all nodes, return success
// 5. Otherwise, a cycle exists (circular dependency)
//
// Time Complexity: O(V + E) where V = vertices, E = edges
// Space Complexity: O(V)
