package orchestrator

import (
	"errors"
	"fmt"
	"time"
)

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
	if approved, exists := gate.Approvers[approver]; exists && approved {
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
