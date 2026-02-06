package orchestrator

import (
	"sync"
	"testing"
)

// TestApprovalStatus tests approval status constants
func TestApprovalStatus(t *testing.T) {
	tests := []struct {
		name  string
		status ApprovalStatus
		want  string
	}{
		{"Pending status", StatusPending, "pending"},
		{"Approved status", StatusApproved, "approved"},
		{"Rejected status", StatusRejected, "rejected"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(tt.status); got != tt.want {
				t.Errorf("ApprovalStatus = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNewApprovalGateManager tests manager creation
func TestNewApprovalGateManager(t *testing.T) {
	am := NewApprovalGateManager()

	if am == nil {
		t.Fatal("NewApprovalGateManager() returned nil")
	}

	if am.gates == nil {
		t.Error("gates map not initialized")
	}
}

// TestCreateGate tests creating an approval gate
func TestCreateGate(t *testing.T) {
	am := NewApprovalGateManager()

	gate := ApprovalGate{
		ID:                "test-gate",
		Name:              "Test Gate",
		Description:       "A test approval gate",
		RequiredApprovers: 2,
		Approvers:         make(map[string]bool),
		Status:            StatusPending,
	}

	err := am.CreateGate(gate)
	if err != nil {
		t.Fatalf("CreateGate() error = %v", err)
	}

	got, err := am.GetGate("test-gate")
	if err != nil {
		t.Fatalf("GetGate() error = %v", err)
	}

	if got.ID != gate.ID {
		t.Errorf("GetGate() ID = %v, want %v", got.ID, gate.ID)
	}

	if got.Name != gate.Name {
		t.Errorf("GetGate() Name = %v, want %v", got.Name, gate.Name)
	}
}

// TestCreateGateDuplicate tests creating duplicate gate IDs
func TestCreateGateDuplicate(t *testing.T) {
	am := NewApprovalGateManager()

	gate := ApprovalGate{
		ID:                "dup-gate",
		Name:              "Duplicate Gate",
		RequiredApprovers: 1,
		Approvers:         make(map[string]bool),
		Status:            StatusPending,
	}

	err := am.CreateGate(gate)
	if err != nil {
		t.Fatalf("First CreateGate() error = %v", err)
	}

	err = am.CreateGate(gate)
	if err == nil {
		t.Error("CreateGate() duplicate should return error")
	}
}

// TestCreateGateValidation tests gate validation
func TestCreateGateValidation(t *testing.T) {
	am := NewApprovalGateManager()

	tests := []struct {
		name    string
		gate    ApprovalGate
		wantErr bool
	}{
		{
			name: "Empty ID",
			gate: ApprovalGate{
				ID:                "",
				Name:              "No ID",
				RequiredApprovers: 1,
				Approvers:         make(map[string]bool),
				Status:            StatusPending,
			},
			wantErr: true,
		},
		{
			name: "Empty name",
			gate: ApprovalGate{
				ID:                "no-name",
				Name:              "",
				RequiredApprovers: 1,
				Approvers:         make(map[string]bool),
				Status:            StatusPending,
			},
			wantErr: true,
		},
		{
			name: "Zero required approvers",
			gate: ApprovalGate{
				ID:                "zero-approvers",
				Name:              "Zero Approvers",
				RequiredApprovers: 0,
				Approvers:         make(map[string]bool),
				Status:            StatusPending,
			},
			wantErr: true,
		},
		{
			name: "Nil approvers map (should be auto-initialized)",
			gate: ApprovalGate{
				ID:                "nil-approvers",
				Name:              "Nil Approvers",
				RequiredApprovers: 1,
				Approvers:         nil,
				Status:            StatusPending,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := am.CreateGate(tt.gate)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateGate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGetGateNotFound tests getting non-existent gate
func TestGetGateNotFound(t *testing.T) {
	am := NewApprovalGateManager()

	_, err := am.GetGate("non-existent")
	if err == nil {
		t.Error("GetGate() non-existent should return error")
	}
}

// TestApprove tests approving a gate
func TestApprove(t *testing.T) {
	am := NewApprovalGateManager()

	gate := ApprovalGate{
		ID:                "approve-gate",
		Name:              "Approve Gate",
		RequiredApprovers: 2,
		Approvers:         make(map[string]bool),
		Status:            StatusPending,
	}

	am.CreateGate(gate)

	err := am.Approve("approve-gate", "alice", "LGTM")
	if err != nil {
		t.Fatalf("Approve() error = %v", err)
	}

	got, _ := am.GetGate("approve-gate")
	if !got.Approvers["alice"] {
		t.Error("Approve() did not mark alice as approved")
	}

	if got.Status != StatusPending {
		t.Errorf("Approve() status = %v, want %v (still need 2nd approval)", got.Status, StatusPending)
	}
}

// TestApproveAlreadyApproved tests duplicate approval
func TestApproveAlreadyApproved(t *testing.T) {
	am := NewApprovalGateManager()

	gate := ApprovalGate{
		ID:                "already-approved",
		Name:              "Already Approved",
		RequiredApprovers: 1,
		Approvers:         make(map[string]bool),
		Status:            StatusPending,
	}

	am.CreateGate(gate)
	am.Approve("already-approved", "alice", "LGTM")

	err := am.Approve("already-approved", "alice", "Again")
	if err == nil {
		t.Error("Approve() duplicate should return error")
	}
}

// TestApproveComplete tests gate completion
func TestApproveComplete(t *testing.T) {
	am := NewApprovalGateManager()

	gate := ApprovalGate{
		ID:                "complete-gate",
		Name:              "Complete Gate",
		RequiredApprovers: 2,
		Approvers:         make(map[string]bool),
		Status:            StatusPending,
	}

	am.CreateGate(gate)
	am.Approve("complete-gate", "alice", "LGTM")
	am.Approve("complete-gate", "bob", "Good to go")

	got, _ := am.GetGate("complete-gate")
	if got.Status != StatusApproved {
		t.Errorf("Approve() complete status = %v, want %v", got.Status, StatusApproved)
	}
}

// TestReject tests rejecting a gate
func TestReject(t *testing.T) {
	am := NewApprovalGateManager()

	gate := ApprovalGate{
		ID:                "reject-gate",
		Name:              "Reject Gate",
		RequiredApprovers: 2,
		Approvers:         make(map[string]bool),
		Status:            StatusPending,
	}

	am.CreateGate(gate)

	err := am.Reject("reject-gate", "charlie", "Not ready")
	if err != nil {
		t.Fatalf("Reject() error = %v", err)
	}

	got, _ := am.GetGate("reject-gate")
	if got.Status != StatusRejected {
		t.Errorf("Reject() status = %v, want %v", got.Status, StatusRejected)
	}
}

// TestCheckGateApproved tests checking gate approval
func TestCheckGateApproved(t *testing.T) {
	am := NewApprovalGateManager()

	gate := ApprovalGate{
		ID:                "check-gate",
		Name:              "Check Gate",
		RequiredApprovers: 2,
		Approvers:         make(map[string]bool),
		Status:            StatusPending,
	}

	am.CreateGate(gate)

	// Not enough approvals
	err := am.CheckGateApproved("check-gate")
	if err == nil {
		t.Error("CheckGateApproved() should return error when not enough approvers")
	}

	// Add required approvals
	am.Approve("check-gate", "alice", "LGTM")
	am.Approve("check-gate", "bob", "Good")

	err = am.CheckGateApproved("check-gate")
	if err != nil {
		t.Errorf("CheckGateApproved() should not return error when approved: %v", err)
	}
}

// TestBlockExecutionUntilApproved tests blocking execution
func TestBlockExecutionUntilApproved(t *testing.T) {
	am := NewApprovalGateManager()

	gate := ApprovalGate{
		ID:                "block-gate",
		Name:              "Block Gate",
		RequiredApprovers: 1,
		Approvers:         make(map[string]bool),
		Status:            StatusPending,
	}

	am.CreateGate(gate)

	// Should block
	err := am.BlockExecutionUntilApproved("block-gate")
	if err == nil {
		t.Error("BlockExecutionUntilApproved() should block when pending")
	}

	// Approve and retry
	am.Approve("block-gate", "alice", "LGTM")
	err = am.BlockExecutionUntilApproved("block-gate")
	if err != nil {
		t.Errorf("BlockExecutionUntilApproved() should not block when approved: %v", err)
	}
}

// TestListGates tests listing all gates
func TestListGates(t *testing.T) {
	am := NewApprovalGateManager()

	gate1 := ApprovalGate{
		ID:                "gate-1",
		Name:              "Gate 1",
		RequiredApprovers: 1,
		Approvers:         make(map[string]bool),
		Status:            StatusPending,
	}
	gate2 := ApprovalGate{
		ID:                "gate-2",
		Name:              "Gate 2",
		RequiredApprovers: 2,
		Approvers:         make(map[string]bool),
		Status:            StatusPending,
	}

	am.CreateGate(gate1)
	am.CreateGate(gate2)

	gates := am.ListGates()
	if len(gates) != 2 {
		t.Fatalf("ListGates() count = %v, want 2", len(gates))
	}
}

// TestDeleteGate tests deleting a gate
func TestDeleteGate(t *testing.T) {
	am := NewApprovalGateManager()

	gate := ApprovalGate{
		ID:                "delete-gate",
		Name:              "Delete Gate",
		RequiredApprovers: 1,
		Approvers:         make(map[string]bool),
		Status:            StatusPending,
	}

	am.CreateGate(gate)

	err := am.DeleteGate("delete-gate")
	if err != nil {
		t.Fatalf("DeleteGate() error = %v", err)
	}

	_, err = am.GetGate("delete-gate")
	if err == nil {
		t.Error("GetGate() after delete should return error")
	}
}

// TestGetPendingApprovals tests getting pending approvals
func TestGetPendingApprovals(t *testing.T) {
	am := NewApprovalGateManager()

	gate1 := ApprovalGate{
		ID:                "pending-1",
		Name:              "Pending 1",
		RequiredApprovers: 1,
		Approvers:         make(map[string]bool),
		Status:            StatusPending,
	}
	gate2 := ApprovalGate{
		ID:                "approved-1",
		Name:              "Approved 1",
		RequiredApprovers: 1,
		Approvers:         make(map[string]bool),
		Status:            StatusPending,
	}
	gate3 := ApprovalGate{
		ID:                "rejected-1",
		Name:              "Rejected 1",
		RequiredApprovers: 1,
		Approvers:         make(map[string]bool),
		Status:            StatusPending,
	}

	am.CreateGate(gate1)
	am.CreateGate(gate2)
	am.CreateGate(gate3)

	am.Approve("approved-1", "alice", "LGTM")
	am.Reject("rejected-1", "bob", "Not good")

	pending := am.GetPendingApprovals()
	if len(pending) != 1 {
		t.Fatalf("GetPendingApprovals() count = %v, want 1", len(pending))
	}

	if pending[0].ID != "pending-1" {
		t.Errorf("GetPendingApprovals() returned wrong gate: %v", pending[0].ID)
	}
}

// TestGetPendingApprovalsEmpty tests getting pending when none pending
func TestGetPendingApprovalsEmpty(t *testing.T) {
	am := NewApprovalGateManager()

	gate := ApprovalGate{
		ID:                "all-approved",
		Name:              "All Approved",
		RequiredApprovers: 1,
		Approvers:         make(map[string]bool),
		Status:            StatusPending,
	}

	am.CreateGate(gate)
	am.Approve("all-approved", "alice", "LGTM")

	pending := am.GetPendingApprovals()
	if len(pending) != 0 {
		t.Errorf("GetPendingApprovals() count = %v, want 0", len(pending))
	}
}

// TestGetApproverCount tests getting approver count
func TestGetApproverCount(t *testing.T) {
	am := NewApprovalGateManager()

	gate := ApprovalGate{
		ID:                "count-gate",
		Name:              "Count Gate",
		RequiredApprovers: 3,
		Approvers:         make(map[string]bool),
		Status:            StatusPending,
	}

	am.CreateGate(gate)

	// Initial count
	count, err := am.GetApproverCount("count-gate")
	if err != nil {
		t.Fatalf("GetApproverCount() error = %v", err)
	}
	if count != 0 {
		t.Errorf("GetApproverCount() = %v, want 0", count)
	}

	// Add approvers
	am.Approve("count-gate", "alice", "LGTM")
	am.Approve("count-gate", "bob", "Good")

	count, err = am.GetApproverCount("count-gate")
	if err != nil {
		t.Fatalf("GetApproverCount() error = %v", err)
	}
	if count != 2 {
		t.Errorf("GetApproverCount() = %v, want 2", count)
	}
}

// TestConcurrentAccessApproval tests concurrent gate operations
func TestConcurrentAccessApproval(t *testing.T) {
	am := NewApprovalGateManager()
	var wg sync.WaitGroup

	// Create gate first
	gate := ApprovalGate{
		ID:                "concurrent-gate",
		Name:              "Concurrent Gate",
		RequiredApprovers: 10,
		Approvers:         make(map[string]bool),
		Status:            StatusPending,
	}
	am.CreateGate(gate)

	// Launch 10 goroutines approving
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			approver := string(rune('a' + idx))
			am.Approve("concurrent-gate", approver, "Approved")
		}(i)
	}

	wg.Wait()

	got, _ := am.GetGate("concurrent-gate")
	if len(got.Approvers) != 10 {
		t.Errorf("Concurrent access: expected 10 approvers, got %d", len(got.Approvers))
	}

	if got.Status != StatusApproved {
		t.Errorf("Concurrent access: status = %v, want %v", got.Status, StatusApproved)
	}
}

// BenchmarkCreateGate benchmarks gate creation performance
func BenchmarkCreateGate(b *testing.B) {
	am := NewApprovalGateManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gate := ApprovalGate{
			ID:                string(rune(i)),
			Name:              "Benchmark Gate",
			RequiredApprovers: 1,
			Approvers:         make(map[string]bool),
			Status:            StatusPending,
		}
		am.CreateGate(gate)
	}
}

// TestApproveNonExistentGate tests approving non-existent gate
func TestApproveNonExistentGate(t *testing.T) {
	am := NewApprovalGateManager()

	err := am.Approve("non-existent", "alice", "LGTM")
	if err == nil {
		t.Error("Approve() non-existent gate should return error")
	}
}

// TestRejectNonExistentGate tests rejecting non-existent gate
func TestRejectNonExistentGate(t *testing.T) {
	am := NewApprovalGateManager()

	err := am.Reject("non-existent", "bob", "No")
	if err == nil {
		t.Error("Reject() non-existent gate should return error")
	}
}

// TestRejectAlreadyApproved tests rejecting after approving
func TestRejectAlreadyApproved(t *testing.T) {
	am := NewApprovalGateManager()

	gate := ApprovalGate{
		ID:                "reject-after-approve",
		Name:              "Reject After Approve",
		RequiredApprovers: 1,
		Approvers:         make(map[string]bool),
		Status:            StatusPending,
	}

	am.CreateGate(gate)
	am.Approve("reject-after-approve", "alice", "LGTM")

	err := am.Reject("reject-after-approve", "alice", "Changed mind")
	if err == nil {
		t.Error("Reject() after approve should return error")
	}
}

// TestDeleteNonExistentGate tests deleting non-existent gate
func TestDeleteNonExistentGate(t *testing.T) {
	am := NewApprovalGateManager()

	err := am.DeleteGate("non-existent")
	if err == nil {
		t.Error("DeleteGate() non-existent should return error")
	}
}

// TestGetApproverCountNonExistent tests getting count for non-existent gate
func TestGetApproverCountNonExistent(t *testing.T) {
	am := NewApprovalGateManager()

	_, err := am.GetApproverCount("non-existent")
	if err == nil {
		t.Error("GetApproverCount() non-existent should return error")
	}
}
