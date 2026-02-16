package evidence

import "testing"

func TestLessonEvent(t *testing.T) {
	l := Lesson{WSID: "00-054-11", Outcome: "passed", WhatWorked: []string{"A: ok"}, Category: "verification"}
	ev := LessonEvent(l)
	if ev.Type != "lesson" || ev.WSID != "00-054-11" {
		t.Errorf("LessonEvent: got %+v", ev)
	}
	if ev.Data == nil {
		t.Fatal("LessonEvent Data is nil")
	}
	l2 := Lesson{WSID: "00-054-13", Outcome: "failed", WhatFailed: []string{"B: fail"}}
	ev2 := LessonEvent(l2)
	if ev2.Type != "lesson" {
		t.Errorf("LessonEvent failed: got %+v", ev2)
	}
}

func TestDecisionEvent(t *testing.T) {
	ev := DecisionEvent("00-054-10", "How to store decisions?", "Evidence log", "Single source of truth", []string{"Separate file"}, 0.9, []string{"architecture"}, nil)
	if ev.Type != "decision" || ev.WSID != "00-054-10" {
		t.Errorf("DecisionEvent: got %+v", ev)
	}
	if ev.Data == nil {
		t.Fatal("DecisionEvent Data is nil")
	}
	rev := "evt-123"
	ev2 := DecisionEvent("00-054-10", "Q", "Revert", "Rationale", nil, 0, nil, &rev)
	if ev2.Data == nil {
		t.Fatal("DecisionEvent with reverses: Data is nil")
	}
}
