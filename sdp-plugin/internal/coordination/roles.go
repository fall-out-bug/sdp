package coordination

// Role represents an agent's current role (AC1)
type Role string

const (
	RoleImplementer Role = "implementer"
	RoleReviewer    Role = "reviewer"
	RoleCoordinator Role = "coordinator"
	RoleSpecialist  Role = "specialist"
)

// RoleCapabilities defines what a role can do
type RoleCapabilities struct {
	Role        Role
	CanExecute  bool
	CanReview   bool
	CanDispatch bool
	Specialties []string
}

// roleCapabilities maps roles to their capabilities
var roleCapabilities = map[Role]RoleCapabilities{
	RoleImplementer: {
		Role:        RoleImplementer,
		CanExecute:  true,
		CanReview:   false,
		CanDispatch: false,
		Specialties: []string{},
	},
	RoleReviewer: {
		Role:        RoleReviewer,
		CanExecute:  false,
		CanReview:   true,
		CanDispatch: false,
		Specialties: []string{"code_quality", "testing"},
	},
	RoleCoordinator: {
		Role:        RoleCoordinator,
		CanExecute:  false,
		CanReview:   true,
		CanDispatch: true,
		Specialties: []string{"orchestration", "planning"},
	},
	RoleSpecialist: {
		Role:        RoleSpecialist,
		CanExecute:  true,
		CanReview:   true,
		CanDispatch: false,
		Specialties: []string{"security", "performance", "testing", "infrastructure"},
	},
}

// GetCapabilities returns capabilities for a role
func GetCapabilities(role Role) RoleCapabilities {
	if cap, ok := roleCapabilities[role]; ok {
		return cap
	}
	return RoleCapabilities{Role: role}
}

// IsValidRole checks if a role is valid
func IsValidRole(role string) bool {
	_, ok := roleCapabilities[Role(role)]
	return ok
}
