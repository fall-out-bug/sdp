package coordination

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Store manages agent event storage (AC3, AC5)
type Store struct {
	path     string
	mu       sync.Mutex
	lastHash string
}

// AggregatedStats holds aggregated event statistics (AC4)
type AggregatedStats struct {
	TotalEvents int
	ByType      map[string]int
	ByAgent     map[string]int
	ByRole      map[string]int
}

// NewStore creates a new event store
func NewStore(path string) (*Store, error) {
	s := &Store{path: path}
	if err := s.loadLastHash(); err != nil {
		return nil, fmt.Errorf("failed to load last hash: %w", err)
	}
	return s, nil
}

// loadLastHash loads the last hash from the existing log file
func (s *Store) loadLastHash() error {
	return scanEvents(s.path, func(e *AgentEvent) error {
		s.lastHash = e.Hash
		return nil
	})
}

// Close closes the store
func (s *Store) Close() error {
	return nil
}

// Append appends an event to the log with hash chain
func (s *Store) Append(event *AgentEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event.PrevHash = s.lastHash
	event.Hash = event.ComputeHash()

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	f, err := os.OpenFile(s.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log: %w", err)
	}
	defer f.Close()

	if _, err := fmt.Fprintln(f, string(data)); err != nil {
		return fmt.Errorf("failed to write event: %w", err)
	}

	s.lastHash = event.Hash
	return nil
}

// ReadAll reads all events from the store (AC3)
func (s *Store) ReadAll() ([]*AgentEvent, error) {
	var events []*AgentEvent
	err := scanEvents(s.path, func(e *AgentEvent) error {
		events = append(events, e)
		return nil
	})
	return events, err
}

// FilterByAgent returns events for a specific agent
func (s *Store) FilterByAgent(agentID string) ([]*AgentEvent, error) {
	var filtered []*AgentEvent
	err := scanEvents(s.path, func(e *AgentEvent) error {
		if e.AgentID == agentID {
			filtered = append(filtered, e)
		}
		return nil
	})
	return filtered, err
}

// FilterByTask returns events for a specific task
func (s *Store) FilterByTask(taskID string) ([]*AgentEvent, error) {
	var filtered []*AgentEvent
	err := scanEvents(s.path, func(e *AgentEvent) error {
		if e.TaskID == taskID {
			filtered = append(filtered, e)
		}
		return nil
	})
	return filtered, err
}

// GetAggregatedStats returns aggregated event statistics (AC4)
func (s *Store) GetAggregatedStats() (*AggregatedStats, error) {
	stats := &AggregatedStats{
		ByType:  make(map[string]int),
		ByAgent: make(map[string]int),
		ByRole:  make(map[string]int),
	}

	err := scanEvents(s.path, func(e *AgentEvent) error {
		stats.TotalEvents++
		stats.ByType[e.Type]++
		stats.ByAgent[e.AgentID]++
		stats.ByRole[e.Role]++
		return nil
	})
	return stats, err
}

// VerifyHashChain verifies the hash chain integrity (AC5)
func (s *Store) VerifyHashChain() error {
	return verifyHashChain(s.path)
}
