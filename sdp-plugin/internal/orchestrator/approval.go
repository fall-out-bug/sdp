package orchestrator

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ApprovalStatus represents the state of an approval gate
type ApprovalStatus string

const (
	// StatusPending indicates the gate is awaiting approval
	StatusPending ApprovalStatus = "pending"
	// StatusApproved indicates the gate has been approved
	StatusApproved ApprovalStatus = "approved"
	// StatusRejected indicates the gate has been rejected
	StatusRejected ApprovalStatus = "rejected"
)

// ApprovalGate represents a quality checkpoint requiring approval
type ApprovalGate struct {
	ID                string
	Name              string
	Description       string
	RequiredApprovers int
	Approvers         map[string]bool // approver -> approved flag
	Status            ApprovalStatus
	RejectionReason   string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// ApprovalGateManager manages approval gates with thread-safe operations
type ApprovalGateManager struct {
	mu    sync.RWMutex
	gates map[string]*ApprovalGate
}

// NewApprovalGateManager creates a new ApprovalGateManager instance
func NewApprovalGateManager() *ApprovalGateManager {
	return &ApprovalGateManager{
		gates: make(map[string]*ApprovalGate),
	}
}

// CreateGate adds a new approval gate
func (am *ApprovalGateManager) CreateGate(gate ApprovalGate) error {
	// Validate gate
	if gate.ID == "" {
		return errors.New("gate ID cannot be empty")
	}
	if gate.Name == "" {
		return errors.New("gate name cannot be empty")
	}
	if gate.RequiredApprovers < 1 {
		return errors.New("required approvers must be at least 1")
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	// Check for duplicate
	if _, exists := am.gates[gate.ID]; exists {
		return fmt.Errorf("gate with ID %s already exists", gate.ID)
	}

	// Initialize approvers map if nil
	if gate.Approvers == nil {
		gate.Approvers = make(map[string]bool)
	}

	gate.CreatedAt = time.Now()
	gate.UpdatedAt = time.Now()

	am.gates[gate.ID] = &gate
	return nil
}

// GetGate retrieves a gate by ID
func (am *ApprovalGateManager) GetGate(gateID string) (*ApprovalGate, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	gate, exists := am.gates[gateID]
	if !exists {
		return nil, fmt.Errorf("gate with ID %s not found", gateID)
	}

	return gate, nil
}

// ListGates returns all gates
func (am *ApprovalGateManager) ListGates() []*ApprovalGate {
	am.mu.RLock()
	defer am.mu.RUnlock()

	gates := make([]*ApprovalGate, 0, len(am.gates))
	for _, gate := range am.gates {
		gates = append(gates, gate)
	}

	return gates
}

// DeleteGate removes a gate
func (am *ApprovalGateManager) DeleteGate(gateID string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if _, exists := am.gates[gateID]; !exists {
		return fmt.Errorf("gate with ID %s not found", gateID)
	}

	delete(am.gates, gateID)
	return nil
}
