# sdp-xul: Implement ApprovalGateManager (WS-006) - Quality Gates

> **Issue ID:** sdp-xul
> **Type:** Bug (Critical - P0)
> **Priority:** 0

## Goal

Implement ApprovalGateManager to enforce quality checkpoints as required by F024 Unified Workflow.

## Problem

F024 specification claims WS-006 (ApprovalGateManager implementation) was completed, but `internal/orchestrator/approval.go` does not exist. This blocks:
- Quality gate enforcement before workstream execution
- Human-in-the-loop approval for critical changes
- Team coordination with approval workflow

## Acceptance Criteria

### AC1: ApprovalGate Data Structure
- [ ] `ApprovalGate` struct defined with: ID, Name, RequiredApprovers, Approvers []string, Status (pending/approved/rejected)
- [ ] `ApprovalRequest` struct with: GateID, Requester, Reason, Timestamp, ApproverResponses
- [ ] `ApprovalGateManager` struct implements approval logic

### AC2: Gate Management Operations
- [ ] `CreateGate(gate ApprovalGate)` - creates new approval gate
- [ ] `GetGate(gateID)` - retrieves gate by ID
- [ ] `ListGates()` - returns all gates
- [ ] `DeleteGate(gateID)` - removes gate

### AC3: Request/Approval Workflow
- [ ] `RequestApproval(gateID, requester, reason)` - creates approval request
- [ ] `Approve(gateID, approver, response)` - approves a request
- [ ] `Reject(gateID, approver, reason)` - rejects a request
- [ ] `CheckGateApproved(gateID)` - verifies if gate has required approvals
- [ ] `GetPendingApprovals()` - returns gates awaiting approval

### AC4: Gate Enforcement
- [ ] `BlockExecutionUntilApproved(gateID)` - checks if gate approved, returns error if not
- [ ] Integration with orchestrator to block workstream execution
- [ ] Thread-safe operations with `sync.RWMutex`

### AC5: Testing
- [ ] Unit tests for all operations
- [ ] Concurrent access tests (race detector clean)
- [ ] Gate enforcement tests
- [ ] Approval workflow tests (request → approve → check)

## Scope Files

**NEW:**
- `internal/orchestrator/approval.go` (~200 LOC)
- `internal/orchestrator/approval_test.go` (~250 LOC)

**MODIFY:**
- `internal/orchestrator/orchestrator.go` (integrate approval checks)

## Implementation Steps

1. Define data structures:
   ```go
   type ApprovalStatus string
   const (
       StatusPending  ApprovalStatus = "pending"
       StatusApproved ApprovalStatus = "approved"
       StatusRejected ApprovalStatus = "rejected"
   )

   type ApprovalGate struct {
       ID               string
       Name             string
       Description      string
       RequiredApprovers int
       Approvers        map[string]bool // approver -> approved flag
       Status           ApprovalStatus
   }

   type ApprovalRequest struct {
       GateID      string
       Requester   string
       Reason      string
       Timestamp   time.Time
       Responses   map[string]string // approver -> response
   }
   ```

2. Implement ApprovalGateManager:
   ```go
   type ApprovalGateManager struct {
       mu    sync.RWMutex
       gates map[string]*ApprovalGate
   }

   func NewApprovalGateManager() *ApprovalGateManager
   func (am *ApprovalGateManager) CreateGate(gate ApprovalGate) error
   func (am *ApprovalGateManager) GetGate(gateID string) (*ApprovalGate, error)
   func (am *ApprovalGateManager) Approve(gateID, approver, response string) error
   func (am *ApprovalGateManager) Reject(gateID, approver, reason string) error
   func (am *ApprovalGateManager) CheckGateApproved(gateID) error
   ```

3. Add tests:
   - TestCreateGate, TestApprove, TestReject
   - TestCheckGateApproved
   - TestBlockExecutionUntilApproved
   - TestConcurrentAccess
   - TestApprovalWorkflow

## Quality Gates

- Test coverage ≥80%
- No race conditions (`go test -race`)
- No TODOs/FIXMEs
- Files <200 LOC
- go vet clean

## Dependencies

- TeamManager (sdp-p35) - for role coordination

## Estimated Scope

- ~200 LOC implementation
- ~250 LOC tests
- Duration: 2-3 hours
