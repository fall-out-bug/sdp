package metrics

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Taxonomy manages failure classifications (AC3-AC6).
type Taxonomy struct {
	path            string
	classifications map[string]*FailureClassification
	mu              sync.RWMutex
}

// NewTaxonomy creates a taxonomy manager for the given path.
func NewTaxonomy(path string) *Taxonomy {
	return &Taxonomy{
		path:            path,
		classifications: make(map[string]*FailureClassification),
	}
}

// Load loads existing taxonomy from file (AC6).
func (t *Taxonomy) Load() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	data, err := os.ReadFile(t.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No existing taxonomy
		}
		return fmt.Errorf("read taxonomy: %w", err)
	}

	if len(data) == 0 {
		return nil // Empty file
	}

	var list []FailureClassification
	if err := json.Unmarshal(data, &list); err != nil {
		return fmt.Errorf("parse taxonomy: %w", err)
	}

	// Rebuild map
	t.classifications = make(map[string]*FailureClassification)
	for i := range list {
		fc := &list[i]
		t.classifications[fc.EventID] = fc
	}

	return nil
}

// Save writes taxonomy to file (AC6).
func (t *Taxonomy) Save() error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// Convert map to slice for JSON
	list := make([]FailureClassification, 0, len(t.classifications))
	for _, fc := range t.classifications {
		list = append(list, *fc)
	}

	// Ensure directory exists
	dir := filepath.Dir(t.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create taxonomy dir: %w", err)
	}

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal taxonomy: %w", err)
	}

	return os.WriteFile(t.path, data, 0644)
}
