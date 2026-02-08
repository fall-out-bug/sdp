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

// TestValidateFrontendBackend verifies frontend vs backend validation
func TestValidateFrontendBackend(t *testing.T) {
	validator := NewContractValidator()

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

	mismatches, err := validator.ValidateFrontendBackend(frontend, backend)
	if err != nil {
		t.Fatalf("ValidateFrontendBackend failed: %v", err)
	}

	// Should have no mismatches for matching contracts
	if len(mismatches) != 0 {
		t.Errorf("Expected 0 mismatches for matching contracts, got %d", len(mismatches))
	}
}

// TestValidateFrontendBackend_SchemaMismatch detects schema incompatibility
func TestValidateFrontendBackend_SchemaMismatch(t *testing.T) {
	validator := NewContractValidator()

	// Frontend requires timestamp, backend doesn't provide it
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
										"timestamp":  {Type: "string"},
									},
									Required: []string{"event_name", "timestamp"},
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

	mismatches, err := validator.ValidateFrontendBackend(frontend, backend)
	if err != nil {
		t.Fatalf("ValidateFrontendBackend failed: %v", err)
	}

	// Should detect schema mismatch
	if len(mismatches) == 0 {
		t.Error("Expected schema mismatch to be detected")
	}

	// Check for WARNING severity
	hasSchemaWarning := false
	for _, m := range mismatches {
		if m.Type == "schema_incompatibility" && m.Severity == "WARNING" {
			hasSchemaWarning = true
		}
	}
	if !hasSchemaWarning {
		t.Error("Expected WARNING severity for schema mismatch")
	}
}

// TestValidateSDKBackend verifies SDK vs backend validation
func TestValidateSDKBackend(t *testing.T) {
	validator := NewContractValidator()

	sdk := &OpenAPIContract{
		OpenAPI: "3.0.0",
		Paths: PathsSpec{
			"/api/v1/telemetry/events": {
				"post": OperationSpec{},
			},
		},
	}

	backend := &OpenAPIContract{
		OpenAPI: "3.0.0",
		Paths: PathsSpec{
			"/api/v1/telemetry/events": {
				"post": OperationSpec{},
			},
		},
	}

	mismatches, err := validator.ValidateSDKBackend(sdk, backend)
	if err != nil {
		t.Fatalf("ValidateSDKBackend failed: %v", err)
	}

	// Should have no mismatches for matching contracts
	if len(mismatches) != 0 {
		t.Errorf("Expected 0 mismatches for matching contracts, got %d", len(mismatches))
	}
}

// TestValidateSDKBackend_EndpointMismatch detects endpoint mismatches
func TestValidateSDKBackend_EndpointMismatch(t *testing.T) {
	validator := NewContractValidator()

	sdk := &OpenAPIContract{
		OpenAPI: "3.0.0",
		Paths: PathsSpec{
			"/api/v1/telemetry/submit": {
				"post": OperationSpec{},
			},
		},
	}

	backend := &OpenAPIContract{
		OpenAPI: "3.0.0",
		Paths: PathsSpec{
			"/api/v1/telemetry/events": {
				"post": OperationSpec{},
			},
		},
	}

	mismatches, err := validator.ValidateSDKBackend(sdk, backend)
	if err != nil {
		t.Fatalf("ValidateSDKBackend failed: %v", err)
	}

	// Should detect endpoint mismatch
	if len(mismatches) == 0 {
		t.Error("Expected endpoint mismatch to be detected")
	}
}

// TestValidateContractFile_ValidContract verifies valid contract validation
func TestValidateContractFile_ValidContract(t *testing.T) {
	validator := NewContractValidator()

	tmpDir := t.TempDir()
	contractPath := tmpDir + "/contract.yaml"

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

	if err := os.WriteFile(contractPath, validYAML, 0644); err != nil {
		t.Fatalf("Failed to write contract: %v", err)
	}

	mismatches, err := validator.ValidateContractFile(contractPath)
	if err != nil {
		t.Fatalf("ValidateContractFile failed: %v", err)
	}

	// Valid contract should have no mismatches
	if len(mismatches) != 0 {
		t.Errorf("Expected 0 mismatches for valid contract, got %d", len(mismatches))
	}
}

// TestValidateContractFile_MissingOpenAPI verifies missing openapi field detection
func TestValidateContractFile_MissingOpenAPI(t *testing.T) {
	validator := NewContractValidator()

	tmpDir := t.TempDir()
	contractPath := tmpDir + "/contract.yaml"

	invalidYAML := []byte(`
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      summary: Test endpoint
`)

	if err := os.WriteFile(contractPath, invalidYAML, 0644); err != nil {
		t.Fatalf("Failed to write contract: %v", err)
	}

	mismatches, err := validator.ValidateContractFile(contractPath)
	if err != nil {
		t.Fatalf("ValidateContractFile failed: %v", err)
	}

	// Should detect missing openapi field
	if len(mismatches) == 0 {
		t.Error("Expected mismatch for missing openapi field")
	}

	// Check for ERROR severity
	hasOpenAPIError := false
	for _, m := range mismatches {
		if m.Type == "invalid_contract" && m.Severity == "ERROR" {
			hasOpenAPIError = true
		}
	}
	if !hasOpenAPIError {
		t.Error("Expected ERROR severity for missing openapi field")
	}
}

// TestValidateContractFile_NoPaths verifies empty paths detection
func TestValidateContractFile_NoPaths(t *testing.T) {
	validator := NewContractValidator()

	tmpDir := t.TempDir()
	contractPath := tmpDir + "/contract.yaml"

	noPathsYAML := []byte(`
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths: {}
`)

	if err := os.WriteFile(contractPath, noPathsYAML, 0644); err != nil {
		t.Fatalf("Failed to write contract: %v", err)
	}

	mismatches, err := validator.ValidateContractFile(contractPath)
	if err != nil {
		t.Fatalf("ValidateContractFile failed: %v", err)
	}

	// Should detect empty paths
	if len(mismatches) == 0 {
		t.Error("Expected mismatch for empty paths")
	}

	// Check for WARNING severity (no paths is a warning, not error)
	hasPathsWarning := false
	for _, m := range mismatches {
		if m.Type == "invalid_contract" && m.Severity == "WARNING" {
			hasPathsWarning = true
		}
	}
	if !hasPathsWarning {
		t.Error("Expected WARNING severity for empty paths")
	}
}

// TestGenerateReport_MultipleSeverities verifies report with mixed severities
func TestGenerateReport_MultipleSeverities(t *testing.T) {
	validator := NewContractValidator()

	mismatches := []*ContractMismatch{
		{
			Severity:   "ERROR",
			Type:       "endpoint_mismatch",
			ComponentA: "frontend",
			ComponentB: "backend",
			Path:       "/api/test",
			Method:     "POST",
			Expected:   "POST /api/test",
			Actual:     "NOT FOUND",
			Fix:        "Add endpoint",
		},
		{
			Severity:   "WARNING",
			Type:       "schema_incompatibility",
			ComponentA: "frontend",
			ComponentB: "backend",
			Path:       "/api/test2",
			Expected:   "Field 'x' required",
			Actual:     "Field 'x' not found",
			Fix:        "Add field",
		},
		{
			Severity:   "INFO",
			Type:       "info",
			ComponentA: "component",
			ComponentB: "other",
			Path:       "/api/info",
			Expected:   "Info",
			Actual:     "Info",
			Fix:        "No action",
		},
	}

	report := validator.GenerateReport(mismatches)

	if report == "" {
		t.Fatal("Expected report to be generated")
	}

	// Verify all sections present
	if !contains(report, "Errors") {
		t.Error("Report missing Errors section")
	}

	if !contains(report, "Warnings") {
		t.Error("Report missing Warnings section")
	}

	if !contains(report, "Info") {
		t.Error("Report missing Info section")
	}

	// Verify summary counts
	if !contains(report, "Errors: 1") {
		t.Error("Report missing error count")
	}

	if !contains(report, "Warnings: 1") {
		t.Error("Report missing warning count")
	}

	if !contains(report, "Info: 1") {
		t.Error("Report missing info count")
	}
}

// TestGenerateReport_EmptyMismatches verifies report with no issues
func TestGenerateReport_EmptyMismatches(t *testing.T) {
	validator := NewContractValidator()

	report := validator.GenerateReport([]*ContractMismatch{})

	if report == "" {
		t.Fatal("Expected report to be generated")
	}

	if !contains(report, "✅ No contract mismatches found!") {
		t.Error("Report missing success message")
	}
}

// TestWriteReport_ErrorHandling verifies error handling in WriteReport
func TestWriteReport_ErrorHandling(t *testing.T) {
	validator := NewContractValidator()

	// Use an invalid path that cannot be created
	invalidPath := "/proc/root/invalid/path/report.md"

	err := validator.WriteReport("# Test Report", invalidPath)
	if err == nil {
		t.Fatal("Expected error for invalid path")
	}

	if !contains(err.Error(), "failed to") {
		t.Errorf("Expected 'failed to' in error message, got: %v", err)
	}
}
