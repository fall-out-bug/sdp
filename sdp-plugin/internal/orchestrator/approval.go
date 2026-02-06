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

// Approve adds an approval to a gate
func (am *ApprovalGateManager) Approve(gateID, approver, response string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	gate, exists := am.gates[gateID]
	if !exists {
		return fmt.Errorf("gate with ID %s not found", gateID)
	}

	// Check if already approved
	if _, approved := gate.Approvers[approver]; approved {
		return fmt.Errorf("approver %s has already approved this gate", approver)
	}

	// Cannot approve a rejected gate
	if gate.Status == StatusRejected {
		return errors.New("cannot approve a rejected gate")
	}

	gate.Approvers[approver] = true
	gate.UpdatedAt = time.Now()

	// Check if gate is now fully approved
	if len(gate.Approvers) >= gate.RequiredApprovers {
		gate.Status = StatusApproved
	}

	return nil
}

// Reject rejects a gate with a reason
func (am *ApprovalGateManager) Reject(gateID, approver, reason string) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	gate, exists := am.gates[gateID]
	if !exists {
		return fmt.Errorf("gate with ID %s not found", gateID)
	}

	// Check if approver has already approved
	if approved, _ := gate.Approvers[approver]; approved {
		return errors.New("cannot reject gate after approving")
	}

	gate.Status = StatusRejected
	gate.RejectionReason = reason
	gate.UpdatedAt = time.Now()

	return nil
}

// CheckGateApproved verifies if a gate has sufficient approvals
func (am *ApprovalGateManager) CheckGateApproved(gateID string) error {
	am.mu.RLock()
	defer am.mu.RUnlock()

	gate, exists := am.gates[gateID]
	if !exists {
		return fmt.Errorf("gate with ID %s not found", gateID)
	}

	if gate.Status == StatusRejected {
		return fmt.Errorf("gate %s has been rejected: %s", gateID, gate.RejectionReason)
	}

	if gate.Status != StatusApproved {
		return fmt.Errorf("gate %s is not approved (status: %s, approvers: %d/%d)",
			gateID, gate.Status, len(gate.Approvers), gate.RequiredApprovers)
	}

	return nil
}

// BlockExecutionUntilApproved blocks until gate is approved
func (am *ApprovalGateManager) BlockExecutionUntilApproved(gateID string) error {
	return am.CheckGateApproved(gateID)
}

// GetPendingApprovals returns all gates that are pending approval
func (am *ApprovalGateManager) GetPendingApprovals() []*ApprovalGate {
	am.mu.RLock()
	defer am.mu.RUnlock()

	pending := make([]*ApprovalGate, 0)
	for _, gate := range am.gates {
		if gate.Status == StatusPending {
			pending = append(pending, gate)
		}
	}

	return pending
}

// GetApproverCount returns the number of approvers for a gate
func (am *ApprovalGateManager) GetApproverCount(gateID string) (int, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	gate, exists := am.gates[gateID]
	if !exists {
		return 0, fmt.Errorf("gate with ID %s not found", gateID)
	}

	return len(gate.Approvers), nil
}
