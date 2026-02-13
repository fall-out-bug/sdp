package coordination

import (
	"time"
)

// Event type for role switching (AC3)
const EventTypeRoleSwitch = "role_switch"

// Task represents a task to be assigned
type Task struct {
	ID            string
	Type          string
	AssignedAgent string
	Metadata      map[string]interface{}
}

// RoleSwitcher handles dynamic role switching (AC2, AC3, AC4, AC5)
type RoleSwitcher struct {
	CurrentRole Role
	AgentID     string
}

// NewRoleSwitcher creates a new role switcher
func NewRoleSwitcher() *RoleSwitcher {
	return &RoleSwitcher{
		CurrentRole: RoleImplementer,
	}
}

// AssignRole assigns a role based on task type and history (AC2)
func (s *RoleSwitcher) AssignRole(task *Task, history []AgentEvent) Role {
	// Check for specialty metadata (AC4)
	if specialty, ok := task.Metadata["specialty"]; ok {
		if _, isString := specialty.(string); isString {
			return RoleSpecialist
		}
	}

	// Check task type
	switch task.Type {
	case "review":
		return RoleReviewer
	case "coordinate", "orchestrate":
		return RoleCoordinator
	case "specialist", "security", "performance":
		return RoleSpecialist
	default:
		return RoleImplementer
	}
}

// SwitchRole performs a role transition and returns an event (AC3)
func (s *RoleSwitcher) SwitchRole(newRole Role, reason string) *AgentEvent {
	if s.CurrentRole == newRole {
		return nil
	}

	oldRole := s.CurrentRole
	s.CurrentRole = newRole

	return &AgentEvent{
		ID:        generateEventID(),
		Type:      EventTypeRoleSwitch,
		AgentID:   s.AgentID,
		Role:      string(newRole),
		Timestamp: time.Now(),
		Payload: map[string]interface{}{
			"from_role": string(oldRole),
			"to_role":   string(newRole),
			"reason":    reason,
		},
	}
}

// IsSelfReview checks if an agent is reviewing their own work (AC4)
func (s *RoleSwitcher) IsSelfReview(taskID, agentID string, history []AgentEvent) bool {
	for _, event := range history {
		if event.TaskID == taskID && event.AgentID == agentID && event.Type == EventTypeAgentComplete {
			return true
		}
	}
	return false
}

// GetPromptContext returns context to inject based on role (AC4)
func (s *RoleSwitcher) GetPromptContext() map[string]interface{} {
	cap := GetCapabilities(s.CurrentRole)
	return map[string]interface{}{
		"role":         string(s.CurrentRole),
		"can_execute":  cap.CanExecute,
		"can_review":   cap.CanReview,
		"can_dispatch": cap.CanDispatch,
		"specialties":  cap.Specialties,
	}
}
