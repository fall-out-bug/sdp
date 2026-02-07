package graph

import "time"

// Checkpoint represents a serialized state of dispatcher execution
type Checkpoint struct {
	Version        string                   `json:"version"`
	FeatureID      string                   `json:"feature_id"`
	Timestamp      time.Time                `json:"timestamp"`
	Completed      []string                 `json:"completed"`
	Failed         []string                 `json:"failed"`
	Graph          *GraphSnapshot           `json:"graph"`
	CircuitBreaker *CircuitBreakerSnapshot `json:"circuit_breaker"`
}

// GraphSnapshot represents the state of the dependency graph
type GraphSnapshot struct {
	Nodes []NodeSnapshot      `json:"nodes"`
	Edges map[string][]string `json:"edges"`
}

// NodeSnapshot represents the state of a single workstream node
type NodeSnapshot struct {
	ID        string   `json:"id"`
	DependsOn []string `json:"depends_on"`
	Indegree  int      `json:"indegree"`
	Completed bool     `json:"completed"`
}

// CircuitBreakerSnapshot represents the state of the circuit breaker
type CircuitBreakerSnapshot struct {
	State            int       `json:"state"`
	FailureCount     int       `json:"failure_count"`
	SuccessCount     int       `json:"success_count"`
	ConsecutiveOpens int       `json:"consecutive_opens"`
	LastFailureTime  time.Time `json:"last_failure_time"`
}
