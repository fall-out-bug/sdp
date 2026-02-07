package orchestrator

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// RoleStatus represents the activation state of a team role
type RoleStatus string

const (
	// RoleStatusActive indicates the role is currently active
	RoleStatusActive RoleStatus = "active"
	// RoleStatusDormant indicates the role is currently dormant
	RoleStatusDormant RoleStatus = "dormant"
)

// TeamRole represents a team role with activation status (extends Role loader's Role)
type TeamRole struct {
	ID          string
	Name        string
	Description string
	Permissions []string
	Status      RoleStatus
}

// RoleRegistry defines the interface for role management operations
type RoleRegistry interface {
	AddRole(role TeamRole) error
	RemoveRole(roleID string) error
	GetRole(roleID string) (TeamRole, error)
	ListRoles() []TeamRole
	ListActiveRoles() []TeamRole
	ActivateRole(roleID string) error
	DeactivateRole(roleID string) error
}

// TeamManager manages team roles with thread-safe operations
type TeamManager struct {
	mu       sync.RWMutex
	roles    map[string]TeamRole
	lastSync time.Time
}

// NewTeamManager creates a new TeamManager instance
func NewTeamManager() *TeamManager {
	return &TeamManager{
		roles: make(map[string]TeamRole),
	}
}

// AddRole adds a new role to the registry
func (tm *TeamManager) AddRole(role TeamRole) error {
	// Validate role
	if role.ID == "" {
		return errors.New("role ID cannot be empty")
	}
	if role.Name == "" {
		return errors.New("role name cannot be empty")
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Check for duplicate
	if _, exists := tm.roles[role.ID]; exists {
		return fmt.Errorf("role with ID %s already exists", role.ID)
	}

	tm.roles[role.ID] = role
	return nil
}

// RemoveRole removes a role from the registry
func (tm *TeamManager) RemoveRole(roleID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if _, exists := tm.roles[roleID]; !exists {
		return fmt.Errorf("role with ID %s not found", roleID)
	}

	delete(tm.roles, roleID)
	return nil
}

// GetRole retrieves a role by ID
func (tm *TeamManager) GetRole(roleID string) (TeamRole, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	role, exists := tm.roles[roleID]
	if !exists {
		return TeamRole{}, fmt.Errorf("role with ID %s not found", roleID)
	}

	return role, nil
}

// ListRoles returns all roles in the registry
func (tm *TeamManager) ListRoles() []TeamRole {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	roles := make([]TeamRole, 0, len(tm.roles))
	for _, role := range tm.roles {
		roles = append(roles, role)
	}

	return roles
}

// ListActiveRoles returns only active roles
func (tm *TeamManager) ListActiveRoles() []TeamRole {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	active := make([]TeamRole, 0)
	for _, role := range tm.roles {
		if role.Status == RoleStatusActive {
			active = append(active, role)
		}
	}

	return active
}

// ActivateRole changes a role's status to active
func (tm *TeamManager) ActivateRole(roleID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	role, exists := tm.roles[roleID]
	if !exists {
		return fmt.Errorf("role with ID %s not found", roleID)
	}

	role.Status = RoleStatusActive
	tm.roles[roleID] = role
	return nil
}

// DeactivateRole changes a role's status to dormant
func (tm *TeamManager) DeactivateRole(roleID string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	role, exists := tm.roles[roleID]
	if !exists {
		return fmt.Errorf("role with ID %s not found", roleID)
	}

	role.Status = RoleStatusDormant
	tm.roles[roleID] = role
	return nil
}

// LastSync returns the last synchronization timestamp
func (tm *TeamManager) LastSync() time.Time {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.lastSync
}

// UpdateSync updates the last synchronization timestamp
func (tm *TeamManager) UpdateSync() {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.lastSync = time.Now()
}
