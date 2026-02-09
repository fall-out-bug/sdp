package evidence

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fall-out-bug/sdp/internal/config"
)

// Emit appends an event to the evidence log (AC6, AC7). Non-blocking; errors are ignored.
func Emit(ev *Event) {
	if ev == nil {
		return
	}
	ev2 := *ev
	if ev2.ID == "" {
		ev2.ID = "evt-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	if ev2.Timestamp == "" {
		ev2.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	go func() {
		_ = emitSync(&ev2)
	}()
}

// EmitSync writes the event immediately (use from CLI so process exit doesn't drop it).
func EmitSync(ev *Event) error {
	if ev == nil {
		return nil
	}
	ev2 := *ev
	if ev2.ID == "" {
		ev2.ID = "evt-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	if ev2.Timestamp == "" {
		ev2.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	return emitSync(&ev2)
}

// emitSync writes event to log; returns error (caller may ignore).
func emitSync(ev *Event) error {
	root, err := config.FindProjectRoot()
	if err != nil {
		return err
	}
	cfg, err := config.Load(root)
	if err != nil {
		return err
	}
	if cfg == nil || !cfg.Evidence.Enabled {
		return nil
	}
	logPath := cfg.Evidence.LogPath
	if logPath == "" {
		logPath = ".sdp/log/events.jsonl"
	}
	path := filepath.Join(root, logPath)
	w, err := NewWriter(path)
	if err != nil {
		return err
	}
	return w.Append(ev)
}

// Enabled returns whether evidence emission is enabled (AC8).
func Enabled() bool {
	root, err := config.FindProjectRoot()
	if err != nil {
		return false
	}
	cfg, err := config.Load(root)
	if err != nil || cfg == nil {
		return true
	}
	return cfg.Evidence.Enabled
}

// ModelID returns SDP_MODEL_ID or ANTHROPIC_MODEL or "unknown" (AC5).
func ModelID() string {
	if s := os.Getenv("SDP_MODEL_ID"); s != "" {
		return s
	}
	if s := os.Getenv("ANTHROPIC_MODEL"); s != "" {
		return s
	}
	return "unknown"
}

// PlanEvent builds a plan event (AC1).
func PlanEvent(wsID string, scopeFiles []string) *Event {
	return PlanEventWithFeature(wsID, "", scopeFiles)
}

// PlanEventWithFeature builds a plan event with feature_id (F056).
func PlanEventWithFeature(wsID, featureID string, scopeFiles []string) *Event {
	data := map[string]interface{}{
		"scope_files": scopeFiles,
		"action":      "activate",
	}
	if featureID != "" {
		data["feature_id"] = featureID
	}
	return &Event{
		Type: "plan",
		WSID: wsID,
		Data: data,
	}
}

// PlanEventForDesign builds a plan event for @design completion (F056).
func PlanEventForDesign(wsID, featureID string, wsCount int, scopeFiles []string, metadata map[string]interface{}) *Event {
	data := map[string]interface{}{
		"feature_id":   featureID,
		"ws_count":    wsCount,
		"scope_files": scopeFiles,
		"action":      "design_complete",
	}
	for k, v := range metadata {
		data[k] = v
	}
	return &Event{
		Type: "plan",
		WSID: wsID,
		Data: data,
	}
}

// QAPair is a question-answer pair for idea evidence (F056).
type QAPair struct {
	Q string
	A string
}

// PlanEventForIdea builds a plan event for @idea completion (F056).
func PlanEventForIdea(wsID, featureID string, questionCount int, summary string, qaPairs []QAPair) *Event {
	data := map[string]interface{}{
		"feature_id":     featureID,
		"question_count": questionCount,
		"summary":        summary,
		"action":         "idea_complete",
	}
	if len(qaPairs) > 0 {
		qa := make([]map[string]string, len(qaPairs))
		for i, p := range qaPairs {
			qa[i] = map[string]string{"question": p.Q, "answer": p.A}
		}
		data["qa_pairs"] = qa
	}
	return &Event{
		Type: "plan",
		WSID: wsID,
		Data: data,
	}
}

// GenerationEvent builds a generation event (AC2).
func GenerationEvent(wsID string, filesChanged []string) *Event {
	return &Event{
		Type: "generation",
		WSID: wsID,
		Data: map[string]interface{}{
			"model_id":      ModelID(),
			"model_version": "",
			"prompt_hash":   "",
			"files_changed": filesChanged,
		},
	}
}

// VerificationEvent builds a verification event (AC3, AC4).
func VerificationEvent(wsID string, passed bool, gateName string, coverage float64) *Event {
	return VerificationEventWithFindings(wsID, passed, gateName, coverage, "")
}

// VerificationEventWithFindings builds a verification event with reviewer output (F056).
func VerificationEventWithFindings(wsID string, passed bool, gateName string, coverage float64, findings string) *Event {
	data := map[string]interface{}{
		"passed":    passed,
		"gate_name": gateName,
		"coverage":  coverage,
	}
	if findings != "" {
		data["findings"] = findings
	}
	return &Event{
		Type: "verification",
		WSID: wsID,
		Data: data,
	}
}

// ApprovalEvent builds an approval event (F056: deploy).
func ApprovalEvent(wsID, targetBranch, commitSHA, approvedBy string) *Event {
	data := map[string]interface{}{
		"target_branch": targetBranch,
		"commit_sha":    commitSHA,
		"approved_by":   approvedBy,
	}
	return &Event{
		Type: "approval",
		WSID: wsID,
		Data: data,
	}
}

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
		"insight":      lessonInsight(lesson),
		"source_ws_id": lesson.WSID,
		"outcome":      outcome,
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
