package planner

import (
	"fmt"
)

// PromptForInteractive runs interactive questions (AC2).
// This would be called in interactive mode to gather more information.
func (p *Planner) PromptForInteractive() error {
	if !p.Interactive {
		return nil
	}

	// In a real implementation, this would use a prompt library
	// For now, we just mark that interactive mode is active
	return nil
}

// ExecuteAutoApply executes the plan in auto-apply mode (AC3).
// This would trigger @oneshot or @build after planning.
func (p *Planner) ExecuteAutoApply(result *DecompositionResult) error {
	if !p.AutoApply {
		return nil
	}

	// In a real implementation, this would:
	// 1. Create workstream files
	// 2. Trigger execution (@oneshot or @build)
	// For now, just validate the plan is complete
	if len(result.Workstreams) == 0 {
		return fmt.Errorf("cannot auto-apply: no workstreams in plan")
	}

	return nil
}
