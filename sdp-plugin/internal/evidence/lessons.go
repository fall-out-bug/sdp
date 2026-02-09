package evidence

import (
	"strings"

	"github.com/fall-out-bug/sdp/internal/verify"
)

// Lesson represents an auto-extracted lesson from a completed workstream (AC2).
type Lesson struct {
	WSID            string   `json:"ws_id"`
	WhatWorked      []string `json:"what_worked,omitempty"`
	WhatFailed      []string `json:"what_failed,omitempty"`
	Category        string   `json:"category,omitempty"`
	RelatedDecisions []string `json:"related_decisions,omitempty"`
	Outcome         string   `json:"outcome"` // "passed", "failed", "mixed"
}

// ExtractLesson builds a lesson from verification result (AC1, AC2).
func ExtractLesson(wsID string, result *verify.VerificationResult) Lesson {
	l := Lesson{WSID: wsID}
	if result == nil {
		l.Outcome = "unknown"
		return l
	}
	if result.Passed {
		l.Outcome = "passed"
		for _, c := range result.Checks {
			if c.Passed {
				l.WhatWorked = append(l.WhatWorked, c.Name+": "+c.Message)
			}
		}
	} else {
		l.Outcome = "failed"
		for _, c := range result.Checks {
			if c.Passed {
				l.WhatWorked = append(l.WhatWorked, c.Name+": "+c.Message)
			} else {
				l.WhatFailed = append(l.WhatFailed, c.Name+": "+c.Message)
			}
		}
		for _, cmd := range result.FailedCommands {
			l.WhatFailed = append(l.WhatFailed, "command: "+cmd)
		}
		for _, f := range result.MissingFiles {
			l.WhatFailed = append(l.WhatFailed, "missing: "+f)
		}
	}
	if len(l.WhatFailed) > 0 && len(l.WhatWorked) > 0 {
		l.Outcome = "mixed"
	}
	if len(l.WhatFailed) > 0 {
		l.Category = "verification"
	}
	return l
}

// MatchesOutcome returns true if lesson outcome matches filter (AC6).
func (l Lesson) MatchesOutcome(filter string) bool {
	if filter == "" {
		return true
	}
	return strings.EqualFold(l.Outcome, filter)
}

// EmitLesson emits a lesson event to the evidence log when enabled (AC1).
func EmitLesson(lesson Lesson) {
	if !Enabled() {
		return
	}
	_ = EmitSync(LessonEvent(lesson))
}
