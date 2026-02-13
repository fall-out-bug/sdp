package coordination

import (
	"testing"
)

func TestRole_String(t *testing.T) {
	tests := []struct {
		role     Role
		expected string
	}{
		{RoleImplementer, "implementer"},
		{RoleReviewer, "reviewer"},
		{RoleCoordinator, "coordinator"},
		{RoleSpecialist, "specialist"},
	}

	for _, tt := range tests {
		t.Run(string(tt.role), func(t *testing.T) {
			if string(tt.role) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.role))
			}
		})
	}
}

func TestRoleCapabilities(t *testing.T) {
	cap := GetCapabilities(RoleImplementer)

	if !cap.CanExecute {
		t.Error("Implementer should be able to execute")
	}
	if cap.CanReview {
		t.Error("Implementer should not be able to review")
	}
	if cap.CanDispatch {
		t.Error("Implementer should not be able to dispatch")
	}
}

func TestRoleCapabilities_Reviewer(t *testing.T) {
	cap := GetCapabilities(RoleReviewer)

	if cap.CanExecute {
		t.Error("Reviewer should not be able to execute")
	}
	if !cap.CanReview {
		t.Error("Reviewer should be able to review")
	}
}

func TestRoleCapabilities_Coordinator(t *testing.T) {
	cap := GetCapabilities(RoleCoordinator)

	if !cap.CanDispatch {
		t.Error("Coordinator should be able to dispatch")
	}
	if !cap.CanReview {
		t.Error("Coordinator should be able to review")
	}
}

func TestRoleCapabilities_Specialist(t *testing.T) {
	cap := GetCapabilities(RoleSpecialist)

	if !cap.CanExecute {
		t.Error("Specialist should be able to execute")
	}
	if len(cap.Specialties) == 0 {
		t.Error("Specialist should have specialties")
	}
}
