package nextstep

import (
	"encoding/json"
	"testing"
)

// TestNextStepRecommendationSchema tests AC2: Define output schema
func TestNextStepRecommendationSchema(t *testing.T) {
	rec := Recommendation{
		Command:     "sdp apply --ws 00-069-01",
		Reason:      "Ready to execute next workstream",
		Confidence:  0.85,
		Category:    CategoryExecution,
		Alternatives: []Alternative{
			{Command: "sdp status", Reason: "View current project state"},
		},
	}

	// Test JSON marshaling (AC5: machine-readable output)
	data, err := json.Marshal(rec)
	if err != nil {
		t.Fatalf("Failed to marshal recommendation: %v", err)
	}

	var unmarshaled Recommendation
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal recommendation: %v", err)
	}

	if unmarshaled.Command != rec.Command {
		t.Errorf("Expected command %s, got %s", rec.Command, unmarshaled.Command)
	}
	if unmarshaled.Confidence != rec.Confidence {
		t.Errorf("Expected confidence %.2f, got %.2f", rec.Confidence, unmarshaled.Confidence)
	}
}

// TestStateInputs tests AC1: Define state inputs
func TestStateInputs(t *testing.T) {
	state := ProjectState{
		Workstreams: []WorkstreamStatus{
			{ID: "00-069-01", Status: StatusReady, Priority: 0, BlockedBy: nil},
			{ID: "00-069-02", Status: StatusBacklog, Priority: 0, BlockedBy: []string{"00-069-01"}},
		},
		LastCommand: "sdp apply --ws 00-068-05",
		Mode:        ModeDrive,
		GitStatus: GitStatusInfo{
			Branch:         "feature/F069-next-step",
			Uncommitted:    false,
			UpstreamDiverg: false,
		},
		Config: ConfigInfo{
			HasSDPConfig: true,
			Version:      "0.10.0",
		},
	}

	if len(state.Workstreams) != 2 {
		t.Errorf("Expected 2 workstreams, got %d", len(state.Workstreams))
	}
	if state.Mode != ModeDrive {
		t.Errorf("Expected ModeDrive, got %v", state.Mode)
	}
	if state.LastCommand == "" {
		t.Error("Expected non-empty last command")
	}
}

// TestTieBreakRules tests AC3: Deterministic tie-break rules
func TestTieBreakRules(t *testing.T) {
	tests := []struct {
		name     string
		ws1      WorkstreamStatus
		ws2      WorkstreamStatus
		expected string // which should be preferred
	}{
		{
			name:     "Lower priority wins",
			ws1:      WorkstreamStatus{ID: "00-069-01", Priority: 0},
			ws2:      WorkstreamStatus{ID: "00-069-02", Priority: 1},
			expected: "00-069-01",
		},
		{
			name:     "Same priority - lower ID wins",
			ws1:      WorkstreamStatus{ID: "00-069-02", Priority: 0},
			ws2:      WorkstreamStatus{ID: "00-069-01", Priority: 0},
			expected: "00-069-01",
		},
		{
			name:     "Ready status over backlog",
			ws1:      WorkstreamStatus{ID: "00-069-01", Status: StatusReady, Priority: 1},
			ws2:      WorkstreamStatus{ID: "00-069-02", Status: StatusBacklog, Priority: 0},
			expected: "00-069-01",
		},
		{
			name:     "Unblocked over blocked",
			ws1:      WorkstreamStatus{ID: "00-069-01", Status: StatusReady, Priority: 1, BlockedBy: nil},
			ws2:      WorkstreamStatus{ID: "00-069-02", Status: StatusReady, Priority: 0, BlockedBy: []string{"00-068-99"}},
			expected: "00-069-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComparePriority(tt.ws1, tt.ws2)
			winnerID := tt.ws1.ID
			if result > 0 {
				winnerID = tt.ws2.ID
			}
			if winnerID != tt.expected {
				t.Errorf("Expected %s to win, got %s (result: %d)", tt.expected, winnerID, result)
			}
		})
	}
}

// TestRecommendationVersioning tests that contract is versioned (DoD)
func TestRecommendationVersioning(t *testing.T) {
	rec := Recommendation{
		Command:    "sdp apply --ws 00-069-01",
		Reason:     "Test",
		Confidence: 0.9,
		Category:   CategoryExecution,
		Version:    ContractVersion,
	}

	if rec.Version == "" {
		t.Error("Recommendation should have version")
	}
	if rec.Version != ContractVersion {
		t.Errorf("Expected version %s, got %s", ContractVersion, rec.Version)
	}
}

// TestConsumerSurfaceParsing tests AC4: All consumer surfaces can parse output
func TestConsumerSurfaceParsing(t *testing.T) {
	rec := Recommendation{
		Command:    "sdp apply --ws 00-069-01",
		Reason:     "Ready for execution",
		Confidence: 0.95,
		Category:   CategoryExecution,
		Version:    ContractVersion,
		Alternatives: []Alternative{
			{Command: "sdp status", Reason: "Check project state"},
		},
		Metadata: map[string]any{
			"workstream_id": "00-069-01",
			"feature_id":    "F069",
		},
	}

	// Test JSON output for scripting/UI
	jsonData, err := json.Marshal(rec)
	if err != nil {
		t.Fatalf("Failed to marshal for JSON consumer: %v", err)
	}

	// Verify JSON is valid
	var parsed map[string]any
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Verify required fields exist
	requiredFields := []string{"command", "reason", "confidence", "category", "version"}
	for _, field := range requiredFields {
		if _, ok := parsed[field]; !ok {
			t.Errorf("Missing required field in JSON output: %s", field)
		}
	}
}

// TestSchemaValidation tests AC5: Schema validation tests
func TestSchemaValidation(t *testing.T) {
	tests := []struct {
		name    string
		rec     Recommendation
		wantErr bool
	}{
		{
			name: "Valid recommendation",
			rec: Recommendation{
				Command:    "sdp apply --ws 00-069-01",
				Reason:     "Valid reason",
				Confidence: 0.8,
				Category:   CategoryExecution,
				Version:    ContractVersion,
			},
			wantErr: false,
		},
		{
			name: "Empty command invalid",
			rec: Recommendation{
				Command:    "",
				Reason:     "Invalid",
				Confidence: 0.5,
				Category:   CategoryExecution,
			},
			wantErr: true,
		},
		{
			name: "Confidence out of range",
			rec: Recommendation{
				Command:    "sdp status",
				Reason:     "Invalid confidence",
				Confidence: 1.5,
				Category:   CategoryExecution,
			},
			wantErr: true,
		},
		{
			name: "Empty category invalid",
			rec: Recommendation{
				Command:    "sdp status",
				Reason:     "Invalid category",
				Confidence: 0.5,
				Category:   "",
			},
			wantErr: true,
		},
		{
			name: "Negative confidence invalid",
			rec: Recommendation{
				Command:    "sdp status",
				Reason:     "Invalid confidence",
				Confidence: -0.1,
				Category:   CategoryExecution,
			},
			wantErr: true,
		},
		{
			name: "Zero confidence valid",
			rec: Recommendation{
				Command:    "sdp status",
				Reason:     "Zero confidence is valid",
				Confidence: 0.0,
				Category:   CategoryExecution,
			},
			wantErr: false,
		},
		{
			name: "Confidence 1.0 valid",
			rec: Recommendation{
				Command:    "sdp status",
				Reason:     "Full confidence is valid",
				Confidence: 1.0,
				Category:   CategoryExecution,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rec.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestRecommendationString tests the String() method.
func TestRecommendationString(t *testing.T) {
	rec := Recommendation{
		Command:    "sdp apply --ws 00-069-01",
		Reason:     "Ready to execute",
		Confidence: 0.85,
		Category:   CategoryExecution,
	}

	expected := "[execution] sdp apply --ws 00-069-01 (85% confidence): Ready to execute"
	if got := rec.String(); got != expected {
		t.Errorf("String() = %q, want %q", got, expected)
	}
}

// TestComparePriorityEdgeCases tests edge cases for priority comparison.
func TestComparePriorityEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		ws1      WorkstreamStatus
		ws2      WorkstreamStatus
		expected int
	}{
		{
			name:     "Identical workstreams",
			ws1:      WorkstreamStatus{ID: "00-069-01", Status: StatusReady, Priority: 0, BlockedBy: nil},
			ws2:      WorkstreamStatus{ID: "00-069-01", Status: StatusReady, Priority: 0, BlockedBy: nil},
			expected: 0,
		},
		{
			name:     "Both ready and blocked - priority wins",
			ws1:      WorkstreamStatus{ID: "00-069-01", Status: StatusReady, Priority: 1, BlockedBy: []string{"00-068-99"}},
			ws2:      WorkstreamStatus{ID: "00-069-02", Status: StatusReady, Priority: 0, BlockedBy: []string{"00-068-98"}},
			expected: 1, // ws2 has lower priority value
		},
		{
			name:     "Neither ready - compare blocked status",
			ws1:      WorkstreamStatus{ID: "00-069-01", Status: StatusBacklog, Priority: 0, BlockedBy: nil},
			ws2:      WorkstreamStatus{ID: "00-069-02", Status: StatusBacklog, Priority: 0, BlockedBy: []string{"00-068-99"}},
			expected: -1, // ws1 is unblocked
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComparePriority(tt.ws1, tt.ws2)
			if result != tt.expected {
				t.Errorf("ComparePriority() = %d, want %d", result, tt.expected)
			}
		})
	}
}
