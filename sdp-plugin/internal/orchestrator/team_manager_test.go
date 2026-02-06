package orchestrator

import (
	"sync"
	"testing"
)

// TestRoleStatus tests role status constants
func TestRoleStatus(t *testing.T) {
	tests := []struct {
		name  string
		status RoleStatus
		want  string
	}{
		{"Active status", RoleStatusActive, "active"},
		{"Dormant status", RoleStatusDormant, "dormant"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(tt.status); got != tt.want {
				t.Errorf("RoleStatus = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNewTeamManager tests TeamManager creation
func TestNewTeamManager(t *testing.T) {
	tm := NewTeamManager()

	if tm == nil {
		t.Fatal("NewTeamManager() returned nil")
	}

	if tm.roles == nil {
		t.Error("roles map not initialized")
	}

	// Mutex is always initialized as zero value, no need to check
}

// TestAddRole tests adding a role to the registry
func TestAddRole(t *testing.T) {
	tm := NewTeamManager()

	role := TeamRole{
		ID:          "test-role",
		Name:        "Test Role",
		Description: "A test role",
		Permissions: []string{"read", "write"},
		Status:      RoleStatusActive,
	}

	err := tm.AddRole(role)
	if err != nil {
		t.Fatalf("AddRole() error = %v", err)
	}

	got, err := tm.GetRole("test-role")
	if err != nil {
		t.Fatalf("GetRole() error = %v", err)
	}

	if got.ID != role.ID {
		t.Errorf("GetRole() ID = %v, want %v", got.ID, role.ID)
	}

	if got.Name != role.Name {
		t.Errorf("GetRole() Name = %v, want %v", got.Name, role.Name)
	}
}

// TestAddRoleDuplicate tests adding duplicate role IDs
func TestAddRoleDuplicate(t *testing.T) {
	tm := NewTeamManager()

	role := TeamRole{
		ID:     "dup-role",
		Name:   "Duplicate Role",
		Status: RoleStatusActive,
	}

	err := tm.AddRole(role)
	if err != nil {
		t.Fatalf("First AddRole() error = %v", err)
	}

	err = tm.AddRole(role)
	if err == nil {
		t.Error("AddRole() duplicate should return error")
	}
}

// TestRemoveRole tests removing a role
func TestRemoveRole(t *testing.T) {
	tm := NewTeamManager()

	role := TeamRole{
		ID:     "remove-role",
		Name:   "Remove Me",
		Status: RoleStatusActive,
	}

	tm.AddRole(role)

	err := tm.RemoveRole("remove-role")
	if err != nil {
		t.Fatalf("RemoveRole() error = %v", err)
	}

	_, err = tm.GetRole("remove-role")
	if err == nil {
		t.Error("GetRole() after remove should return error")
	}
}

// TestRemoveRoleNotFound tests removing non-existent role
func TestRemoveRoleNotFound(t *testing.T) {
	tm := NewTeamManager()

	err := tm.RemoveRole("non-existent")
	if err == nil {
		t.Error("RemoveRole() non-existent should return error")
	}
}

// TestActivateRole tests activating a dormant role
func TestActivateRole(t *testing.T) {
	tm := NewTeamManager()

	role := TeamRole{
		ID:     "dormant-role",
		Name:   "Dormant Role",
		Status: RoleStatusDormant,
	}

	tm.AddRole(role)

	err := tm.ActivateRole("dormant-role")
	if err != nil {
		t.Fatalf("ActivateRole() error = %v", err)
	}

	got, _ := tm.GetRole("dormant-role")
	if got.Status != RoleStatusActive {
		t.Errorf("ActivateRole() status = %v, want %v", got.Status, RoleStatusActive)
	}
}

// TestDeactivateRole tests deactivating an active role
func TestDeactivateRole(t *testing.T) {
	tm := NewTeamManager()

	role := TeamRole{
		ID:     "active-role",
		Name:   "Active Role",
		Status: RoleStatusActive,
	}

	tm.AddRole(role)

	err := tm.DeactivateRole("active-role")
	if err != nil {
		t.Fatalf("DeactivateRole() error = %v", err)
	}

	got, _ := tm.GetRole("active-role")
	if got.Status != RoleStatusDormant {
		t.Errorf("DeactivateRole() status = %v, want %v", got.Status, RoleStatusDormant)
	}
}

// TestListRoles tests listing all roles
func TestListRoles(t *testing.T) {
	tm := NewTeamManager()

	role1 := TeamRole{ID: "role-1", Name: "Role 1", Status: RoleStatusActive}
	role2 := TeamRole{ID: "role-2", Name: "Role 2", Status: RoleStatusDormant}

	tm.AddRole(role1)
	tm.AddRole(role2)

	roles := tm.ListRoles()
	if len(roles) != 2 {
		t.Fatalf("ListRoles() count = %v, want 2", len(roles))
	}
}

// TestListActiveRoles tests listing only active roles
func TestListActiveRoles(t *testing.T) {
	tm := NewTeamManager()

	role1 := TeamRole{ID: "active-1", Name: "Active 1", Status: RoleStatusActive}
	role2 := TeamRole{ID: "dormant-1", Name: "Dormant 1", Status: RoleStatusDormant}
	role3 := TeamRole{ID: "active-2", Name: "Active 2", Status: RoleStatusActive}

	tm.AddRole(role1)
	tm.AddRole(role2)
	tm.AddRole(role3)

	active := tm.ListActiveRoles()
	if len(active) != 2 {
		t.Fatalf("ListActiveRoles() count = %v, want 2", len(active))
	}

	for _, role := range active {
		if role.Status != RoleStatusActive {
			t.Errorf("ListActiveRoles() returned dormant role: %v", role.ID)
		}
	}
}

// TestConcurrentAccess tests concurrent role operations
func TestConcurrentAccess(t *testing.T) {
	tm := NewTeamManager()
	var wg sync.WaitGroup

	// Launch 10 goroutines adding roles
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			role := TeamRole{
				ID:     string(rune('a' + idx)),
				Name:   "Concurrent Role",
				Status: RoleStatusActive,
			}
			tm.AddRole(role)
		}(i)
	}

	wg.Wait()

	roles := tm.ListRoles()
	if len(roles) != 10 {
		t.Errorf("Concurrent access: expected 10 roles, got %d", len(roles))
	}
}

// TestRoleValidation tests role field validation
func TestRoleValidation(t *testing.T) {
	tm := NewTeamManager()

	tests := []struct {
		name    string
		role    TeamRole
		wantErr bool
	}{
		{
			name: "Valid role",
			role: TeamRole{
				ID:     "valid",
				Name:   "Valid Role",
				Status: RoleStatusActive,
			},
			wantErr: false,
		},
		{
			name: "Empty ID",
			role: TeamRole{
				ID:     "",
				Name:   "No ID",
				Status: RoleStatusActive,
			},
			wantErr: true,
		},
		{
			name: "Empty name",
			role: TeamRole{
				ID:     "no-name",
				Name:   "",
				Status: RoleStatusActive,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tm.AddRole(tt.role)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// BenchmarkAddRole benchmarks role addition performance
func BenchmarkAddRole(b *testing.B) {
	tm := NewTeamManager()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		role := TeamRole{
			ID:     string(rune(i)),
			Name:   "Benchmark Role",
			Status: RoleStatusActive,
		}
		tm.AddRole(role)
	}
}
