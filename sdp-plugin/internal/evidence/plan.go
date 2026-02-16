package evidence

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
		"feature_id":  featureID,
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
