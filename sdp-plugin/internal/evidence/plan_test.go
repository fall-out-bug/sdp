package evidence

import "testing"

func TestPlanEvent(t *testing.T) {
	ev := PlanEvent("00-054-01", []string{"schema/index.json"})
	if ev.Type != "plan" || ev.WSID != "00-054-01" {
		t.Errorf("PlanEvent: got %+v", ev)
	}
	if ev.Data == nil {
		t.Fatal("PlanEvent Data is nil")
	}
}

func TestPlanEventWithFeature(t *testing.T) {
	ev := PlanEventWithFeature("00-056-02", "F056", []string{"parse.go", "emitter.go"})
	if ev.Type != "plan" || ev.WSID != "00-056-02" {
		t.Errorf("PlanEventWithFeature: got %+v", ev)
	}
	m, _ := ev.Data.(map[string]interface{})
	if m["feature_id"] != "F056" {
		t.Errorf("feature_id: got %v", m["feature_id"])
	}
}

func TestPlanEventForDesign(t *testing.T) {
	ev := PlanEventForDesign("00-056-00", "F056", 4, []string{"a.md", "b.md"}, map[string]interface{}{"deps": "00-054-05"})
	if ev.Type != "plan" || ev.WSID != "00-056-00" {
		t.Errorf("PlanEventForDesign: got %+v", ev)
	}
	m, _ := ev.Data.(map[string]interface{})
	if m["feature_id"] != "F056" || m["ws_count"].(int) != 4 || m["deps"] != "00-054-05" {
		t.Errorf("PlanEventForDesign data: got %v", m)
	}
}

func TestPlanEventForIdea(t *testing.T) {
	qa := []QAPair{{Q: "Who?", A: "User"}, {Q: "What?", A: "Auth"}}
	ev := PlanEventForIdea("00-056-00", "F056", 2, "Auth requirements", qa)
	if ev.Type != "plan" || ev.WSID != "00-056-00" {
		t.Errorf("PlanEventForIdea: got %+v", ev)
	}
	m, _ := ev.Data.(map[string]interface{})
	if m["feature_id"] != "F056" || m["question_count"].(int) != 2 || m["summary"] != "Auth requirements" {
		t.Errorf("PlanEventForIdea data: got %v", m)
	}
	if m["qa_pairs"] == nil {
		t.Error("qa_pairs missing")
	}
}
