package evidence

import "strings"

// DecisionEvent builds a decision event (AC6, AC7). reverses links to a previous decision being overturned.
func DecisionEvent(wsID, question, choice, rationale string, alternatives []string, confidence float64, tags []string, reverses *string) *Event {
	data := map[string]interface{}{
		"question":   question,
		"choice":     choice,
		"rationale":  rationale,
		"confidence": confidence,
	}
	if len(alternatives) > 0 {
		data["alternatives"] = alternatives
	}
	if len(tags) > 0 {
		data["tags"] = tags
	}
	if reverses != nil && *reverses != "" {
		data["reverses"] = *reverses
	}
	return &Event{
		Type: "decision",
		WSID: wsID,
		Data: data,
	}
}

// LessonEvent builds a lesson event (AC10). Outcome: passed→worked, failed→failed, mixed→mixed.
func LessonEvent(lesson Lesson) *Event {
	outcome := lesson.Outcome
	if outcome == "passed" {
		outcome = "worked"
	}
	data := map[string]interface{}{
		"category":     lesson.Category,
		"insight":     lessonInsight(lesson),
		"source_ws_id": lesson.WSID,
		"outcome":     outcome,
	}
	if len(lesson.RelatedDecisions) > 0 {
		data["related_decisions"] = lesson.RelatedDecisions
	}
	return &Event{
		Type: "lesson",
		WSID: lesson.WSID,
		Data: data,
	}
}

func lessonInsight(l Lesson) string {
	if len(l.WhatFailed) > 0 && len(l.WhatWorked) > 0 {
		return "mixed: some checks passed, some failed"
	}
	if len(l.WhatFailed) > 0 {
		return "failed: " + strings.Join(l.WhatFailed, "; ")
	}
	if len(l.WhatWorked) > 0 {
		return "worked: " + strings.Join(l.WhatWorked, "; ")
	}
	return l.Outcome
}
