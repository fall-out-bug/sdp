package evidence

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

// SkillEvent builds a thin evidence event for a skill (F056-03). Non-blocking use: Emit(SkillEvent(...)).
func SkillEvent(skillName, eventType, wsID string, data map[string]interface{}) *Event {
	if data == nil {
		data = make(map[string]interface{})
	}
	data["skill"] = skillName
	return &Event{
		Type: eventType,
		WSID: wsID,
		Data: data,
	}
}
