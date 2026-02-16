package evidence

import "testing"

func TestVerificationEvent(t *testing.T) {
	ev := VerificationEvent("00-054-01", true, "coverage", 85.0)
	if ev.Type != "verification" || ev.WSID != "00-054-01" {
		t.Errorf("VerificationEvent: got %+v", ev)
	}
}

func TestVerificationEventWithFindings(t *testing.T) {
	ev := VerificationEventWithFindings("00-056-01", false, "QA", 82.0, "Coverage below threshold")
	if ev.Type != "verification" || ev.WSID != "00-056-01" {
		t.Errorf("VerificationEventWithFindings: got %+v", ev)
	}
	if ev.Data == nil {
		t.Fatal("Data is nil")
	}
	m, _ := ev.Data.(map[string]interface{})
	if m["findings"] != "Coverage below threshold" {
		t.Errorf("findings: got %v", m["findings"])
	}
}

func TestApprovalEvent(t *testing.T) {
	ev := ApprovalEvent("00-000-00", "main", "abc123def", "CI")
	if ev.Type != "approval" || ev.WSID != "00-000-00" {
		t.Errorf("ApprovalEvent: got %+v", ev)
	}
	if ev.Data == nil {
		t.Fatal("Data is nil")
	}
	m, _ := ev.Data.(map[string]interface{})
	if m["target_branch"] != "main" || m["commit_sha"] != "abc123def" || m["approved_by"] != "CI" {
		t.Errorf("ApprovalEvent data: got %v", m)
	}
}

func TestSkillEvent(t *testing.T) {
	ev := SkillEvent("vision", "plan", "00-000-00", map[string]interface{}{"strategic": true})
	if ev.Type != "plan" || ev.WSID != "00-000-00" {
		t.Errorf("SkillEvent: got %+v", ev)
	}
	m, _ := ev.Data.(map[string]interface{})
	if m["skill"] != "vision" || m["strategic"] != true {
		t.Errorf("SkillEvent data: got %v", m)
	}
	ev2 := SkillEvent("reality", "verification", "00-056-03", nil)
	if ev2.Data == nil {
		t.Fatal("SkillEvent with nil data: Data is nil")
	}
	m2, _ := ev2.Data.(map[string]interface{})
	if m2["skill"] != "reality" {
		t.Errorf("SkillEvent nil data: got %v", m2)
	}
}

func TestGenerationEvent(t *testing.T) {
	ev := GenerationEvent("00-054-03", []string{"internal/evidence/types.go"})
	if ev.Type != "generation" || ev.WSID != "00-054-03" {
		t.Errorf("GenerationEvent: got %+v", ev)
	}
}
