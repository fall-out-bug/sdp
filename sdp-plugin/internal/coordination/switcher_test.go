package coordination

import (
	"testing"
)

func TestRoleSwitcher_AssignRole_ReviewTask(t *testing.T) {
	switcher := NewRoleSwitcher()
	task := &Task{Type: "review", ID: "task-001", AssignedAgent: "agent-1"}
	history := []AgentEvent{}

	role := switcher.AssignRole(task, history)

	if role != RoleReviewer {
		t.Errorf("Expected reviewer role for review task, got %s", role)
	}
}

func TestRoleSwitcher_AssignRole_ImplementTask(t *testing.T) {
	switcher := NewRoleSwitcher()
	task := &Task{Type: "implement", ID: "task-001", AssignedAgent: "agent-1"}
	history := []AgentEvent{}

	role := switcher.AssignRole(task, history)

	if role != RoleImplementer {
		t.Errorf("Expected implementer role for implement task, got %s", role)
	}
}

func TestRoleSwitcher_AssignRole_SelfReview(t *testing.T) {
	switcher := NewRoleSwitcher()
	task := &Task{Type: "review", ID: "task-001", AssignedAgent: "agent-1"}

	// History shows agent-1 implemented this task
	history := []AgentEvent{
		{Type: EventTypeAgentComplete, AgentID: "agent-1", TaskID: "task-001", Role: string(RoleImplementer)},
	}

	role := switcher.AssignRole(task, history)

	// Should still be reviewer for self-review (but flagged differently in payload)
	if role != RoleReviewer {
		t.Errorf("Expected reviewer role for self-review, got %s", role)
	}
}

func TestRoleSwitcher_AssignRole_SpecialistTask(t *testing.T) {
	switcher := NewRoleSwitcher()
	task := &Task{Type: "implement", ID: "task-001", AssignedAgent: "agent-1", Metadata: map[string]interface{}{"specialty": "security"}}
	history := []AgentEvent{}

	role := switcher.AssignRole(task, history)

	if role != RoleSpecialist {
		t.Errorf("Expected specialist role for specialty task, got %s", role)
	}
}

func TestRoleSwitcher_SwitchRole(t *testing.T) {
	switcher := NewRoleSwitcher()
	switcher.CurrentRole = RoleImplementer

	event := switcher.SwitchRole(RoleReviewer, "task complete, switching to review")

	if event == nil {
		t.Fatal("Expected role switch event")
	}
	if event.Type != EventTypeRoleSwitch {
		t.Errorf("Expected role_switch event type, got %s", event.Type)
	}
	if switcher.CurrentRole != RoleReviewer {
		t.Errorf("Expected current role to be reviewer, got %s", switcher.CurrentRole)
	}
}

func TestRoleSwitcher_SwitchRole_NoChange(t *testing.T) {
	switcher := NewRoleSwitcher()
	switcher.CurrentRole = RoleImplementer

	event := switcher.SwitchRole(RoleImplementer, "no change needed")

	if event != nil {
		t.Error("Expected no event when role doesn't change")
	}
}

func TestRoleSwitcher_IsSelfReview(t *testing.T) {
	switcher := NewRoleSwitcher()

	history := []AgentEvent{
		{Type: EventTypeAgentComplete, AgentID: "agent-1", TaskID: "task-001"},
	}

	selfReview := switcher.IsSelfReview("task-001", "agent-1", history)

	if !selfReview {
		t.Error("Expected to detect self-review")
	}
}

func TestRoleSwitcher_IsSelfReview_DifferentAgent(t *testing.T) {
	switcher := NewRoleSwitcher()

	history := []AgentEvent{
		{Type: EventTypeAgentComplete, AgentID: "agent-1", TaskID: "task-001"},
	}

	selfReview := switcher.IsSelfReview("task-001", "agent-2", history)

	if selfReview {
		t.Error("Should not be self-review when different agent")
	}
}

func TestRoleSwitcher_MultiRole(t *testing.T) {
	switcher := NewRoleSwitcher()
	switcher.CurrentRole = RoleSpecialist

	cap := GetCapabilities(switcher.CurrentRole)

	// Specialist can both execute and review
	if !cap.CanExecute {
		t.Error("Specialist should be able to execute")
	}
	if !cap.CanReview {
		t.Error("Specialist should be able to review")
	}
}

func TestRoleSwitcher_GetPromptContext(t *testing.T) {
	switcher := NewRoleSwitcher()
	switcher.CurrentRole = RoleImplementer

	ctx := switcher.GetPromptContext()

	if ctx["role"] != "implementer" {
		t.Errorf("Expected role implementer, got %v", ctx["role"])
	}
	if ctx["can_execute"] != true {
		t.Error("Implementer should have can_execute true")
	}
	if ctx["can_review"] != false {
		t.Error("Implementer should have can_review false")
	}
}
