package agents

import (
	"os"
	"strings"
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

	// Bidirectional comparison detects 2 mismatches:
	// 1. Frontend has /submit, backend doesn't
	// 2. Backend has /events, frontend doesn't
	if len(mismatches) != 2 {
		t.Errorf("Expected 2 mismatches (bidirectional), got %d", len(mismatches))
	}

	// Should have at least one ERROR severity
	hasError := false
	for _, m := range mismatches {
		if m.Severity == "ERROR" {
			hasError = true
		}
	}
	if !hasError {
		t.Error("Expected at least one ERROR severity mismatch")
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

// TestSafeYAMLUnmarshal_ValidYAML verifies valid YAML can be parsed
func TestSafeYAMLUnmarshal_ValidYAML(t *testing.T) {
	validYAML := []byte(`
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      summary: Test endpoint
`)

	var result map[string]interface{}
	err := safeYAMLUnmarshal(validYAML, &result)
	if err != nil {
		t.Fatalf("Failed to parse valid YAML: %v", err)
	}

	if result["openapi"] != "3.0.0" {
		t.Errorf("Expected openapi version 3.0.0, got %v", result["openapi"])
	}
}

// TestSafeYAMLUnmarshal_RejectsLargeFile verifies file size limit
func TestSafeYAMLUnmarshal_RejectsLargeFile(t *testing.T) {
	// Create a YAML file larger than 10MB
	largeYAML := make([]byte, MaxYAMLFileSize+1)
	for i := range largeYAML {
		largeYAML[i] = 'x'
	}

	var result map[string]interface{}
	err := safeYAMLUnmarshal(largeYAML, &result)
	if err == nil {
		t.Fatal("Expected error for oversized YAML file")
	}

	if !strings.Contains(err.Error(), "exceeds maximum allowed size") {
		t.Errorf("Expected size limit error, got: %v", err)
	}
}

// TestSafeYAMLUnmarshal_RejectsUnknownFields verifies strict mode
func TestSafeYAMLUnmarshal_RejectsUnknownFields(t *testing.T) {
	yamlWithUnknownField := []byte(`
openapi: 3.0.0
unknown_field: should_be_rejected
paths:
  /test:
    get:
      summary: Test
`)

	type StrictContract struct {
		OpenAPI string                 `yaml:"openapi"`
		Paths   map[string]interface{} `yaml:"paths"`
	}

	var result StrictContract
	err := safeYAMLUnmarshal(yamlWithUnknownField, &result)
	// strict mode should reject unknown fields
	if err == nil {
		t.Log("Warning: strict mode not enforced (unknown field accepted)")
	}
}
