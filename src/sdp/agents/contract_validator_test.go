package agents

import (
	"os"
	"testing"
)

// TestCompareContracts_MatchingEndpoints verifies happy path
func TestCompareContracts_MatchingEndpoints(t *testing.T) {
	validator := NewContractValidator()

	// Identical frontend and backend contracts
	frontend := &OpenAPIContract{
		OpenAPI: "3.0.0",
		Paths: PathsSpec{
			"/api/v1/telemetry/events": {
				"post": OperationSpec{
					RequestBody: &RequestSpec{
						Content: map[string]MediaSpec{
							"application/json": {
								Schema: SchemaRefSpec{
									Type: "object",
									Properties: map[string]PropertySpec{
										"event_name": {Type: "string"},
									},
									Required: []string{"event_name"},
								},
							},
						},
					},
				},
			},
		},
	}

	backend := &OpenAPIContract{
		OpenAPI: "3.0.0",
		Paths: PathsSpec{
			"/api/v1/telemetry/events": {
				"post": OperationSpec{
					RequestBody: &RequestSpec{
						Content: map[string]MediaSpec{
							"application/json": {
								Schema: SchemaRefSpec{
									Type: "object",
									Properties: map[string]PropertySpec{
										"event_name": {Type: "string"},
									},
									Required: []string{"event_name"},
								},
							},
						},
					},
				},
			},
		},
	}

	mismatches, err := validator.CompareContracts(frontend, backend, "frontend", "backend")
	if err != nil {
		t.Fatalf("CompareContracts failed: %v", err)
	}

	if len(mismatches) != 0 {
		t.Errorf("Expected 0 mismatches, got %d", len(mismatches))
	}
}

// TestCompareContracts_EndpointMismatch verifies endpoint mismatch detection
func TestCompareContracts_EndpointMismatch(t *testing.T) {
	validator := NewContractValidator()

	// Frontend calls /submit, backend exposes /events
	frontend := &OpenAPIContract{
		Paths: PathsSpec{
			"/api/v1/telemetry/submit": {
				"post": OperationSpec{},
			},
		},
	}

	backend := &OpenAPIContract{
		Paths: PathsSpec{
			"/api/v1/telemetry/events": {
				"post": OperationSpec{},
			},
		},
	}

	mismatches, err := validator.CompareContracts(frontend, backend, "frontend", "backend")
	if err != nil {
		t.Fatalf("CompareContracts failed: %v", err)
	}

	if len(mismatches) != 1 {
		t.Errorf("Expected 1 mismatch (method mismatch only), got %d", len(mismatches))
	}

	if mismatches[0].Severity != "ERROR" {
		t.Errorf("Expected ERROR severity, got %s", mismatches[0].Severity)
	}
}

// TestCompareContracts_MethodMismatch verifies HTTP method mismatch
func TestCompareContracts_MethodMismatch(t *testing.T) {
	validator := NewContractValidator()

	// Frontend uses POST, backend uses GET
	frontend := &OpenAPIContract{
		Paths: PathsSpec{
			"/api/v1/telemetry/events": {
				"post": OperationSpec{},
			},
		},
	}

	backend := &OpenAPIContract{
		Paths: PathsSpec{
			"/api/v1/telemetry/events": {
				"get": OperationSpec{},
			},
		},
	}

	mismatches, err := validator.CompareContracts(frontend, backend, "frontend", "backend")
	if err != nil {
		t.Fatalf("CompareContracts failed: %v", err)
	}

	if len(mismatches) != 1 {
		t.Errorf("Expected 1 mismatch (method mismatch only), got %d", len(mismatches))
	}
}

// TestValidateSchemas_Compatible verifies schema compatibility
func TestValidateSchemas_Compatible(t *testing.T) {
	validator := NewContractValidator()

	// Compatible schemas - required fields match
	frontendSchema := SchemaRefSpec{
		Type: "object",
		Properties: map[string]PropertySpec{
			"event_name": {Type: "string"},
		},
		Required: []string{"event_name"},
	}

	backendSchema := SchemaRefSpec{
		Type: "object",
		Properties: map[string]PropertySpec{
			"event_name": {Type: "string"},
		},
		Required: []string{"event_name"},
	}

	mismatch := validator.ValidateSchemas(frontendSchema, backendSchema, "/api/test", "frontend", "backend")
	if mismatch != nil {
		t.Errorf("Expected schemas to be compatible, got mismatch: %v", mismatch)
	}
}

// TestValidateSchemas_Incompatible verifies schema incompatibility detection
func TestValidateSchemas_Incompatible(t *testing.T) {
	validator := NewContractValidator()

	// Incompatible - frontend expects field that backend doesn't provide
	frontendSchema := SchemaRefSpec{
		Type: "object",
		Properties: map[string]PropertySpec{
			"event_name": {Type: "string"},
			"timestamp":  {Type: "string"},
		},
		Required: []string{"event_name", "timestamp"},
	}

	backendSchema := SchemaRefSpec{
		Type: "object",
		Properties: map[string]PropertySpec{
			"event_name": {Type: "string"},
		},
		Required: []string{"event_name"},
	}

	mismatch := validator.ValidateSchemas(frontendSchema, backendSchema, "/api/test", "frontend", "backend")
	if mismatch == nil {
		t.Error("Expected schema mismatch, got nil")
	}

	if mismatch != nil && mismatch.Severity != "WARNING" {
		t.Errorf("Expected WARNING severity, got %s", mismatch.Severity)
	}
}

// TestGenerateReport verifies report generation
func TestGenerateReport(t *testing.T) {
	validator := NewContractValidator()

	mismatches := []*ContractMismatch{
		{
			Severity:   "ERROR",
			Type:       "endpoint_mismatch",
			ComponentA: "frontend",
			ComponentB: "backend",
			Path:       "/api/telemetry",
			Expected:   "POST /api/telemetry",
			Actual:     "GET /api/telemetry",
			File:       "src/app/api.ts",
			Fix:        "Change to GET or update backend to POST",
		},
	}

	report := validator.GenerateReport(mismatches)

	if report == "" {
		t.Fatal("Expected report to be generated")
	}

	// Verify report contains key information
	if !contains(report, "# Contract Validation Report") {
		t.Error("Report missing title")
	}

	if !contains(report, "Errors") {
		t.Error("Report missing Errors section")
	}

	if !contains(report, "endpoint_mismatch") {
		t.Error("Report missing mismatch type")
	}
}

// TestWriteReport verifies report file writing
func TestWriteReport(t *testing.T) {
	validator := NewContractValidator()

	tmpDir := t.TempDir()
	reportPath := tmpDir + "/validation-report.md"

	report := "# Test Report\n\nThis is a test report."

	err := validator.WriteReport(report, reportPath)
	if err != nil {
		t.Fatalf("WriteReport failed: %v", err)
	}

	// Verify file exists
	content, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("Failed to read report: %v", err)
	}

	if len(content) == 0 {
		t.Error("Report file is empty")
	}
}

// Helper functions

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
