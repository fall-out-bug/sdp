package model

import (
	"testing"
)

func TestRouter_SelectModel(t *testing.T) {
	r := NewRouter()

	tests := []struct {
		taskType   string
		complexity int
		wantSpeed  string
	}{
		{"planning", 8, "quality"},
		{"planning", 3, "balanced"},
		{"code", 5, "balanced"},
		{"review", 3, "quality"},
		{"debug", 5, "fast"},
		{"quick", 5, "fast"},
		{"unknown", 5, "balanced"},
	}

	for _, tt := range tests {
		p := r.SelectModel(tt.taskType, tt.complexity)
		if p.Speed != tt.wantSpeed {
			t.Errorf("SelectModel(%s, %d) = %s, want %s", tt.taskType, tt.complexity, p.Speed, tt.wantSpeed)
		}
	}
}

func TestRouter_SelectBySpeed(t *testing.T) {
	r := NewRouter()

	tests := []struct {
		speed     string
		wantName  string
	}{
		{"fast", "fast"},
		{"balanced", "balanced"},
		{"quality", "quality"},
	}

	for _, tt := range tests {
		p := r.SelectBySpeed(tt.speed)
		if p.Name != tt.wantName {
			t.Errorf("SelectBySpeed(%s) = %s, want %s", tt.speed, p.Name, tt.wantName)
		}
	}
}

func TestRouter_SelectByCost(t *testing.T) {
	r := NewRouter()

	p := r.SelectByCost(0.01)
	if p.CostPer1K > 0.01 {
		t.Errorf("SelectByCost(0.01) returned model with cost %f", p.CostPer1K)
	}

	// Should return cheapest if all are affordable
	p = r.SelectByCost(1.0)
	if p.Name != "fast" && p.Name != "balanced" {
		t.Errorf("SelectByCost(1.0) should return cheapest, got %s", p.Name)
	}
}

func TestRouter_AddProfile(t *testing.T) {
	r := NewRouter()

	custom := Profile{
		Name: "custom", ModelID: "custom-model",
		Speed: "custom", MaxTokens: 100000, CostPer1K: 0.001,
	}
	r.AddProfile(custom)

	p, ok := r.GetProfile("custom")
	if !ok {
		t.Error("Custom profile not found")
	}
	if p.ModelID != "custom-model" {
		t.Errorf("Expected custom-model, got %s", p.ModelID)
	}
}

func TestRouter_AddRule(t *testing.T) {
	r := NewRouter()

	r.AddRule(RoutingRule{
		TaskType: "custom-task", MinComplex: 0, MaxComplex: 10,
		Profile: "fast", Priority: 100,
	})

	p := r.SelectModel("custom-task", 5)
	if p.Name != "fast" {
		t.Errorf("Custom rule not applied, got %s", p.Name)
	}
}

func TestRouter_SetDefault(t *testing.T) {
	r := NewRouter()
	r.SetDefault("quality")

	// Unknown task type should use quality
	p := r.SelectModel("nonexistent", 5)
	if p.Name != "quality" {
		t.Errorf("Default should be quality, got %s", p.Name)
	}
}

func TestRouter_ListProfiles(t *testing.T) {
	r := NewRouter()

	profiles := r.ListProfiles()
	if len(profiles) < 3 {
		t.Errorf("Expected at least 3 profiles, got %d", len(profiles))
	}
}

func TestEstimateCost(t *testing.T) {
	p := Profile{CostPer1K: 0.01}
	cost := EstimateCost(5000, p)

	expected := 0.05 // 5000 / 1000 * 0.01
	if cost != expected {
		t.Errorf("EstimateCost(5000) = %f, want %f", cost, expected)
	}
}

func TestRouter_ComplexityBounds(t *testing.T) {
	r := NewRouter()

	// Low complexity planning
	p := r.SelectModel("planning", 0)
	if p.Speed != "balanced" {
		t.Errorf("Low complexity planning should be balanced, got %s", p.Speed)
	}

	// High complexity planning
	p = r.SelectModel("planning", 10)
	if p.Speed != "quality" {
		t.Errorf("High complexity planning should be quality, got %s", p.Speed)
	}
}
