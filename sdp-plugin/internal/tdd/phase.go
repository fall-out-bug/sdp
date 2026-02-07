package tdd

// Phase represents a TDD cycle phase
type Phase int

const (
	// Red phase: Write failing test
	Red Phase = iota
	// Green phase: Make test pass
	Green
	// Refactor phase: Improve code without changing behavior
	Refactor
)

// String returns the string representation of the phase
func (p Phase) String() string {
	switch p {
	case Red:
		return "Red"
	case Green:
		return "Green"
	case Refactor:
		return "Refactor"
	default:
		return "Unknown"
	}
}
