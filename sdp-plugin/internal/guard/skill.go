package guard

// Skill implements the guard skill logic
type Skill struct {
	stateManager *StateManager
	activeWS     string
}

// NewSkill creates a new guard skill
func NewSkill(configDir string) *Skill {
	return &Skill{
		stateManager: NewStateManager(configDir),
		activeWS:     "",
	}
}
