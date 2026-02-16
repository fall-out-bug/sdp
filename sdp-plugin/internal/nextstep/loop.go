package nextstep

import (
	"fmt"
	"strings"
)

// InteractiveLoop manages the accept/refine/reject loop for recommendations.
type InteractiveLoop struct {
	primary      *Recommendation
	currentIndex int
	context      map[string]any
	history      []LoopResult
}

// NewInteractiveLoop creates a new interactive loop with a recommendation.
func NewInteractiveLoop(rec *Recommendation) *InteractiveLoop {
	return &InteractiveLoop{
		primary:      rec,
		currentIndex: 0,
		context:      copyMetadata(rec.Metadata),
		history:      []LoopResult{},
	}
}

// CurrentRecommendation returns the currently selected recommendation.
func (l *InteractiveLoop) CurrentRecommendation() *Recommendation {
	if l.currentIndex == 0 {
		return l.primary
	}
	if l.currentIndex <= len(l.primary.Alternatives) {
		alt := l.primary.Alternatives[l.currentIndex-1]
		return &Recommendation{
			Command:    alt.Command,
			Reason:     alt.Reason,
			Confidence: l.primary.Confidence * 0.8,
			Category:   l.primary.Category,
			Version:    l.primary.Version,
		}
	}
	return nil
}

// CurrentIndex returns the current selection index.
func (l *InteractiveLoop) CurrentIndex() int {
	return l.currentIndex
}

// Accept accepts the current recommendation.
func (l *InteractiveLoop) Accept() LoopResult {
	rec := l.CurrentRecommendation()
	result := LoopResult{
		Action:  ActionAccepted,
		Command: rec.Command,
		Reason:  rec.Reason,
	}
	l.history = append(l.history, result)
	return result
}

// Reject rejects the current recommendation and moves to the next alternative.
func (l *InteractiveLoop) Reject() LoopResult {
	totalOptions := 1 + len(l.primary.Alternatives)
	l.currentIndex++

	if l.currentIndex >= totalOptions {
		result := LoopResult{
			Action:  ActionRejected,
			Command: "",
			Reason:  "All options rejected",
		}
		l.history = append(l.history, result)
		return result
	}

	rec := l.CurrentRecommendation()
	result := LoopResult{
		Action:  ActionAlternative,
		Command: rec.Command,
		Reason:  rec.Reason,
	}
	l.history = append(l.history, result)
	return result
}

// Refine refines the current command with user input.
func (l *InteractiveLoop) Refine(userInput string) LoopResult {
	rec := l.CurrentRecommendation()
	refinedCmd := refineCommand(rec.Command, userInput)

	result := LoopResult{
		Action:  ActionRefined,
		Command: refinedCmd,
		Reason:  fmt.Sprintf("Refined: %s", userInput),
	}
	l.history = append(l.history, result)
	return result
}

// refineCommand applies user refinements to a command.
func refineCommand(baseCmd, userInput string) string {
	input := strings.TrimSpace(userInput)
	if input == "" {
		return baseCmd
	}

	if strings.HasPrefix(input, "--") || strings.HasPrefix(input, "-") {
		return baseCmd + " " + input
	}

	return baseCmd + " " + input
}

// Why returns an explanation of why this recommendation was made.
func (l *InteractiveLoop) Why() string {
	rec := l.CurrentRecommendation()
	var sb strings.Builder

	sb.WriteString(rec.Reason)

	if wsID, ok := l.context["workstream_id"]; ok {
		sb.WriteString(fmt.Sprintf(" (workstream: %s)", wsID))
	}
	if priority, ok := l.context["priority"]; ok {
		sb.WriteString(fmt.Sprintf(" [priority: %v]", priority))
	}
	if confidence := rec.Confidence; confidence >= 0.9 {
		sb.WriteString(" - High confidence recommendation")
	} else if confidence >= 0.7 {
		sb.WriteString(" - Good confidence recommendation")
	}

	return sb.String()
}

// Context returns the preserved context from the recommendation.
func (l *InteractiveLoop) Context() map[string]any {
	return l.context
}

// History returns the interaction history.
func (l *InteractiveLoop) History() []LoopResult {
	return l.history
}

// CreateCheckpoint creates a checkpoint for later resumption.
func (l *InteractiveLoop) CreateCheckpoint() *LoopCheckpoint {
	return &LoopCheckpoint{
		Primary:      l.primary,
		CurrentIndex: l.currentIndex,
		Context:      copyMetadata(l.context),
		History:      append([]LoopResult{}, l.history...),
	}
}

// ResumeFromCheckpoint creates a new loop from a checkpoint.
func ResumeFromCheckpoint(cp *LoopCheckpoint) *InteractiveLoop {
	return &InteractiveLoop{
		primary:      cp.Primary,
		currentIndex: cp.CurrentIndex,
		context:      copyMetadata(cp.Context),
		history:      append([]LoopResult{}, cp.History...),
	}
}

// copyMetadata creates a copy of the metadata map.
func copyMetadata(m map[string]any) map[string]any {
	if m == nil {
		return make(map[string]any)
	}
	cp := make(map[string]any, len(m))
	for k, v := range m {
		cp[k] = v
	}
	return cp
}
