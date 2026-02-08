package agents

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAnalyzeRequirements_ValidMarkdown verifies requirements parsing
func TestAnalyzeRequirements_ValidMarkdown(t *testing.T) {
	cs := NewContractSynthesizer()
	tmpDir := t.TempDir()

	reqPath := filepath.Join(tmpDir, "sdp-telemetry-requirements.md")
	reqContent := `# Telemetry Feature Requirements

## Endpoints

### POST /api/v1/telemetry/events
Request: {event_name, timestamp}
Response: {event_id}

### GET /api/v1/telemetry/events/{id}
Response: {event}
`

	if err := os.WriteFile(reqPath, []byte(reqContent), 0644); err != nil {
		t.Fatalf("Failed to write requirements: %v", err)
	}

	req, err := cs.AnalyzeRequirements(reqPath)
	if err != nil {
		t.Fatalf("AnalyzeRequirements failed: %v", err)
	}

	if req == nil {
		t.Fatal("Expected requirements to be parsed")
	}

	if req.FeatureName != "telemetry" {
		t.Errorf("Expected feature name 'telemetry', got '%s'", req.FeatureName)
	}

	if len(req.Endpoints) != 2 {
		t.Errorf("Expected 2 endpoints, got %d", len(req.Endpoints))
	}
}

// TestAnalyzeRequirements_ErrorHandling verifies error handling
func TestAnalyzeRequirements_ErrorHandling(t *testing.T) {
	cs := NewContractSynthesizer()

	_, err := cs.AnalyzeRequirements("/nonexistent/file.md")
	if err == nil {
		t.Fatal("Expected error for non-existent file")
	}

	if !strings.Contains(err.Error(), "failed to read") {
		t.Errorf("Expected 'failed to read' error, got: %v", err)
	}
}

// TestWriteContract_Success verifies contract writing
func TestWriteContract_Success(t *testing.T) {
	cs := NewContractSynthesizer()
	tmpDir := t.TempDir()

	outputPath := filepath.Join(tmpDir, "contract.yaml")
	contract := &OpenAPIContract{
		OpenAPI: "3.0.0",
		Info:    InfoSpec{Title: "Test", Version: "1.0"},
		Paths:   make(PathsSpec),
	}

	err := cs.WriteContract(contract, outputPath)
	if err != nil {
		t.Fatalf("WriteContract failed: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Contract file was not created")
	}

	content, _ := os.ReadFile(outputPath)
	if len(content) == 0 {
		t.Error("Contract file is empty")
	}
}

// TestWriteContract_ErrorHandling verifies error handling
func TestWriteContract_ErrorHandling(t *testing.T) {
	cs := NewContractSynthesizer()

	contract := &OpenAPIContract{
		OpenAPI: "3.0.0",
		Info:    InfoSpec{Title: "Test", Version: "1.0"},
		Paths:   make(PathsSpec),
	}

	invalidPath := "/proc/root/invalid/path/contract.yaml"

	err := cs.WriteContract(contract, invalidPath)
	if err == nil {
		t.Fatal("Expected error for invalid path")
	}
}

// TestSynthesizeContract_EndToEnd verifies full synthesis workflow
func TestSynthesizeContract_EndToEnd(t *testing.T) {
	cs := NewContractSynthesizer()
	tmpDir := t.TempDir()

	reqPath := filepath.Join(tmpDir, "sdp-test-requirements.md")
	reqContent := `# Test

### GET /api/test
Response: {data}
`

	os.WriteFile(reqPath, []byte(reqContent), 0644)
	outputPath := filepath.Join(tmpDir, "contract.yaml")

	result, err := cs.SynthesizeContract("test", reqPath, outputPath)
	if err != nil {
		t.Fatalf("SynthesizeContract failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected synthesis result")
	}
}

// TestProposeContract verifies contract proposal generation
func TestProposeContract(t *testing.T) {
	cs := NewContractSynthesizer()

	requirements := &ContractRequirements{
		FeatureName: "telemetry",
		Endpoints: []EndpointSpec{
			{
				Path:     "/api/v1/telemetry/events",
				Method:   "POST",
				Request:  SchemaSpec{Fields: []FieldSpec{{Name: "event_name", Type: "string", Required: true}}},
				Response: SchemaSpec{Fields: []FieldSpec{{Name: "event_id", Type: "string", Required: true}}},
			},
		},
	}

	contract, err := cs.ProposeContract(requirements)
	if err != nil {
		t.Fatalf("ProposeContract failed: %v", err)
	}

	if contract == nil {
		t.Fatal("Expected contract to be proposed")
	}

	if contract.OpenAPI != "3.0.0" {
		t.Errorf("Expected OpenAPI 3.0.0, got %s", contract.OpenAPI)
	}
}

// TestValidateEndpointPath verifies path validation
func TestValidateEndpointPath(t *testing.T) {
	tests := []struct {
		path      string
		shouldErr bool
	}{
		{"/api/test", false},
		{"/api/{id}", false},
		{"api/test", true},
		{"/../api/test", true},
	}

	for _, tt := range tests {
		err := validateEndpointPath(tt.path)
		if tt.shouldErr && err == nil {
			t.Errorf("Expected error for path %q", tt.path)
		}
		if !tt.shouldErr && err != nil {
			t.Errorf("Unexpected error for path %q: %v", tt.path, err)
		}
	}
}

// TestValidateHTTPMethod verifies HTTP method validation
func TestValidateHTTPMethod(t *testing.T) {
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	
	for _, method := range validMethods {
		err := validateHTTPMethod(method)
		if err != nil {
			t.Errorf("Unexpected error for method %s: %v", method, err)
		}
	}
}

// TestSanitizeFieldName verifies field name sanitization
func TestSanitizeFieldName(t *testing.T) {
	tests := []struct {
		name      string
		shouldErr bool
	}{
		{"field_name", false},
		{"FieldName", false},
		{"_private", false},
		{"123field", true},
		{"field-name", true},
		{"field;name", true},
	}

	for _, tt := range tests {
		_, err := sanitizeFieldName(tt.name)
		if tt.shouldErr && err == nil {
			t.Errorf("Expected error for field %q", tt.name)
		}
		if !tt.shouldErr && err != nil {
			t.Errorf("Unexpected error for field %q: %v", tt.name, err)
		}
	}
}

// TestAnalyzeRequirements_InvalidPath validates path validation
func TestAnalyzeRequirements_InvalidPath(t *testing.T) {
	cs := NewContractSynthesizer()
	tmpDir := t.TempDir()

	reqPath := filepath.Join(tmpDir, "sdp-test-requirements.md")
	reqContent := `# Test

### GET /api/../test
Response: {data}
`

	os.WriteFile(reqPath, []byte(reqContent), 0644)

	_, err := cs.AnalyzeRequirements(reqPath)
	if err == nil {
		t.Error("Expected error for invalid path with ..")
	}
}

// TestSynthesizeContract_InvalidMethod validates method validation
func TestSynthesizeContract_InvalidMethod(t *testing.T) {
	cs := NewContractSynthesizer()
	tmpDir := t.TempDir()

	reqPath := filepath.Join(tmpDir, "sdp-test-requirements.md")
	reqContent := `# Test

### INVALID /api/test
Response: {data}
`

	os.WriteFile(reqPath, []byte(reqContent), 0644)
	outputPath := filepath.Join(tmpDir, "contract.yaml")

	_, err := cs.SynthesizeContract("test", reqPath, outputPath)
	if err == nil {
		t.Error("Expected error for invalid HTTP method")
	}
}

// TestSynthesizeContract_AnalyzeError validates error propagation
func TestSynthesizeContract_AnalyzeError(t *testing.T) {
	cs := NewContractSynthesizer()
	tmpDir := t.TempDir()

	reqPath := filepath.Join(tmpDir, "sdp-test-requirements.md")
	outputPath := filepath.Join(tmpDir, "contract.yaml")

	// Don't create the file - should fail at analyze step
	_, err := cs.SynthesizeContract("test", reqPath, outputPath)
	if err == nil {
		t.Error("Expected error for missing requirements file")
	}
}
