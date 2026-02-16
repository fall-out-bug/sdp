package nextstep

// LoopAction represents the action taken in the interactive loop.
type LoopAction int

const (
	// ActionNone indicates no action has been taken yet.
	ActionNone LoopAction = iota
	// ActionAccepted indicates the user accepted the recommendation.
	ActionAccepted
	// ActionRejected indicates the user rejected all options.
	ActionRejected
	// ActionAlternative indicates the user selected an alternative.
	ActionAlternative
	// ActionRefined indicates the user refined the command.
	ActionRefined
	// ActionExited indicates the user exited the loop.
	ActionExited
)

// LoopResult represents the result of a loop interaction.
type LoopResult struct {
	Action  LoopAction
	Command string
	Reason  string
}

// LoopCheckpoint represents a saved state for resumption.
type LoopCheckpoint struct {
	Primary      *Recommendation
	CurrentIndex int
	Context      map[string]any
	History      []LoopResult
}
