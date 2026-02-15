package evidence

import (
	"encoding/json"
	"testing"
)

func TestPlanEvent_New(t *testing.T) {
	workstreams := []string{"00-054-01", "00-054-02", "00-054-03"}
	event := NewPlanEvent("F054", "Evidence Layer", workstreams, "")

	if event.Feature != "F054" {
		t.Errorf("expected feature F054, got %s", event.Feature)
	}
	if event.Description != "Evidence Layer" {
		t.Errorf("expected description, got %s", event.Description)
	}
	if len(event.Workstreams) != 3 {
		t.Errorf("expected 3 workstreams, got %d", len(event.Workstreams))
	}
	if event.BaseEvent.Type != EventTypePlan {
		t.Errorf("expected type plan, got %s", event.BaseEvent.Type)
	}
}

func TestPlanEvent_WithDependencies(t *testing.T) {
	event := NewPlanEvent("F055", "Feature", []string{"WS1"}, "")
	event.Dependencies = []string{"F054"}

	if len(event.Dependencies) != 1 {
		t.Error("dependencies not set")
	}
}

func TestPlanEvent_WithCostEstimate(t *testing.T) {
	event := NewPlanEvent("F056", "Feature", []string{"WS1"}, "")
	event.CostEstimate = "2h"

	if event.CostEstimate != "2h" {
		t.Error("cost estimate not set")
	}
}

func TestPlanEvent_JSON(t *testing.T) {
	original := NewPlanEvent("F054", "Evidence", []string{"WS1", "WS2"}, "")
	original.Dependencies = []string{"F053"}
	original.CostEstimate = "5h"

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var parsed PlanEvent
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if parsed.Feature != original.Feature {
		t.Error("feature mismatch")
	}
	if len(parsed.Workstreams) != 2 {
		t.Error("workstreams mismatch")
	}
}

func TestPlanEvent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		event   PlanEvent
		wantErr bool
	}{
		{
			name:    "valid",
			event:   NewPlanEvent("F054", "Evidence", []string{"WS1"}, ""),
			wantErr: false,
		},
		{
			name: "missing feature",
			event: PlanEvent{
				BaseEvent:   NewBaseEvent(EventTypePlan, ""),
				Description: "Desc",
				Workstreams: []string{"WS1"},
			},
			wantErr: true,
		},
		{
			name: "no workstreams",
			event: PlanEvent{
				BaseEvent: NewBaseEvent(EventTypePlan, ""),
				Feature:   "F054",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
