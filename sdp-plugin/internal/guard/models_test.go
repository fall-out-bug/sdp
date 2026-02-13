package guard

import (
	"testing"
	"time"
)

func TestGuardState_AddFinding(t *testing.T) {
	gs := &GuardState{}

	// Add first finding
	f1 := ReviewFinding{ID: "F001", Priority: 0, Status: "open"}
	gs.AddFinding(f1)

	if len(gs.ReviewFindings) != 1 {
		t.Errorf("Expected 1 finding, got %d", len(gs.ReviewFindings))
	}

	// Add second finding
	f2 := ReviewFinding{ID: "F002", Priority: 2, Status: "open"}
	gs.AddFinding(f2)

	if len(gs.ReviewFindings) != 2 {
		t.Errorf("Expected 2 findings, got %d", len(gs.ReviewFindings))
	}
}

func TestGuardState_ResolveFinding(t *testing.T) {
	gs := &GuardState{
		ReviewFindings: []ReviewFinding{
			{ID: "F001", Priority: 0, Status: "open"},
			{ID: "F002", Priority: 1, Status: "open"},
		},
	}

	// Resolve existing finding
	resolved := gs.ResolveFinding("F001", "test-user")
	if !resolved {
		t.Error("Expected finding to be resolved")
	}

	// Check status changed
	if gs.ReviewFindings[0].Status != "resolved" {
		t.Errorf("Expected status 'resolved', got '%s'", gs.ReviewFindings[0].Status)
	}
	if gs.ReviewFindings[0].ResolvedBy != "test-user" {
		t.Errorf("Expected ResolvedBy 'test-user', got '%s'", gs.ReviewFindings[0].ResolvedBy)
	}

	// Try to resolve non-existent finding
	resolved = gs.ResolveFinding("F999", "test-user")
	if resolved {
		t.Error("Expected non-existent finding to not be resolved")
	}
}

func TestGuardState_GetOpenFindings(t *testing.T) {
	gs := &GuardState{
		ReviewFindings: []ReviewFinding{
			{ID: "F001", Priority: 0, Status: "open"},
			{ID: "F002", Priority: 1, Status: "resolved"},
			{ID: "F003", Priority: 2, Status: "open"},
		},
	}

	open := gs.GetOpenFindings()
	if len(open) != 2 {
		t.Errorf("Expected 2 open findings, got %d", len(open))
	}
}

func TestGuardState_GetBlockingFindings(t *testing.T) {
	gs := &GuardState{
		ReviewFindings: []ReviewFinding{
			{ID: "F001", Priority: 0, Status: "open"},       // Blocking
			{ID: "F002", Priority: 1, Status: "open"},       // Blocking
			{ID: "F003", Priority: 2, Status: "open"},       // Not blocking
			{ID: "F004", Priority: 0, Status: "resolved"},   // Not blocking (resolved)
		},
	}

	blocking := gs.GetBlockingFindings()
	if len(blocking) != 2 {
		t.Errorf("Expected 2 blocking findings, got %d", len(blocking))
	}
}

func TestGuardState_HasBlockingFindings(t *testing.T) {
	// No findings
	gs := &GuardState{}
	if gs.HasBlockingFindings() {
		t.Error("Expected no blocking findings when empty")
	}

	// Only non-blocking findings
	gs.ReviewFindings = []ReviewFinding{
		{ID: "F001", Priority: 2, Status: "open"},
	}
	if gs.HasBlockingFindings() {
		t.Error("Expected no blocking findings for P2")
	}

	// Has blocking finding
	gs.ReviewFindings = append(gs.ReviewFindings, ReviewFinding{
		ID: "F002", Priority: 0, Status: "open",
	})
	if !gs.HasBlockingFindings() {
		t.Error("Expected blocking findings for P0")
	}

	// All resolved
	gs.ReviewFindings[1].Status = "resolved"
	if gs.HasBlockingFindings() {
		t.Error("Expected no blocking findings when resolved")
	}
}

func TestGuardState_FindingCount(t *testing.T) {
	gs := &GuardState{
		ReviewFindings: []ReviewFinding{
			{ID: "F001", Priority: 0, Status: "open"},
			{ID: "F002", Priority: 1, Status: "open"},
			{ID: "F003", Priority: 2, Status: "resolved"},
			{ID: "F004", Priority: 0, Status: "resolved"},
		},
	}

	open, resolved, blocking := gs.FindingCount()

	if open != 2 {
		t.Errorf("Expected 2 open, got %d", open)
	}
	if resolved != 2 {
		t.Errorf("Expected 2 resolved, got %d", resolved)
	}
	if blocking != 2 {
		t.Errorf("Expected 2 blocking (P0/P1), got %d", blocking)
	}
}

func TestGuardState_IsExpired(t *testing.T) {
	// No ActiveWS means expired
	gs := &GuardState{}
	if !gs.IsExpired() {
		t.Error("Expected expired with no ActiveWS")
	}

	// Recently activated with ActiveWS - not expired
	gs.ActiveWS = "00-067-15"
	gs.ActivatedAt = time.Now().Format(time.RFC3339)
	if gs.IsExpired() {
		t.Error("Expected not expired for recent activation")
	}

	// Activated 25 hours ago - expired
	gs.ActivatedAt = time.Now().Add(-25 * time.Hour).Format(time.RFC3339)
	if !gs.IsExpired() {
		t.Error("Expected expired for old activation")
	}

	// Invalid time format - expired
	gs.ActivatedAt = "invalid-time"
	if !gs.IsExpired() {
		t.Error("Expected expired for invalid time format")
	}
}
