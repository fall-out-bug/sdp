package agents

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall-out-bug/sdp/src/sdp/synthesis"
)

// TestContractSynthesizer_AnalyzeRequirements verifies that the agent
// can parse and analyze requirements files
func TestContractSynthesizer_AnalyzeRequirements(t *testing.T) {
	// Create temporary requirements file
	tmpDir := t.TempDir()
	reqPath := filepath.Join(tmpDir, "telemetry-requirements.md")
	reqContent := `# Telemetry Feature Requirements

## API Endpoints

### Submit Telemetry Event
- POST /api/v1/telemetry/events
- Request: {event_name, timestamp, metadata}
- Response: {event_id, status}

### Get Telemetry Event
- GET /api/v1/telemetry/events/{id}
- Response: {event_id, event_name, timestamp, metadata}
`
	if err := os.WriteFile(reqPath, []byte(reqContent), 0644); err != nil {
		t.Fatalf("Failed to create requirements file: %v", err)
	}

	agent := NewContractSynthesizer()
	requirements, err := agent.AnalyzeRequirements(reqPath)

	if err != nil {
		t.Fatalf("AnalyzeRequirements failed: %v", err)
	}

	if requirements.FeatureName != "telemetry" {
		t.Errorf("Expected feature name 'telemetry', got '%s'", requirements.FeatureName)
	}

	if len(requirements.Endpoints) != 2 {
		t.Errorf("Expected 2 endpoints, got %d", len(requirements.Endpoints))
	}

	// Verify POST endpoint
	postEndpoint := findEndpointByPath(requirements.Endpoints, "/api/v1/telemetry/events")
	if postEndpoint == nil {
		t.Fatal("POST endpoint not found")
	}

	if postEndpoint.Method != "POST" {
		t.Errorf("Expected method POST, got %s", postEndpoint.Method)
	}
}

// TestContractSynthesizer_ProposeContract verifies OpenAPI contract generation
func TestContractSynthesizer_ProposeContract(t *testing.T) {
	agent := NewContractSynthesizer()

	requirements := &ContractRequirements{
		FeatureName: "telemetry",
		Endpoints: []EndpointSpec{
			{
				Path:   "/api/v1/telemetry/events",
				Method: "POST",
				Request: SchemaSpec{
					Fields: []FieldSpec{
						{Name: "event_name", Type: "string", Required: true},
						{Name: "timestamp", Type: "string", Required: true},
					},
				},
				Response: SchemaSpec{
					Fields: []FieldSpec{
						{Name: "event_id", Type: "string", Required: true},
						{Name: "status", Type: "string", Required: true},
					},
				},
			},
		},
	}

	contract, err := agent.ProposeContract(requirements)
	if err != nil {
		t.Fatalf("ProposeContract failed: %v", err)
	}

	// Verify OpenAPI structure
	if contract.OpenAPI != "3.0.0" {
		t.Errorf("Expected OpenAPI version 3.0.0, got %s", contract.OpenAPI)
	}

	if contract.Info.Title != "Telemetry API" {
		t.Errorf("Expected title 'Telemetry API', got %s", contract.Info.Title)
	}

	if len(contract.Paths) == 0 {
		t.Error("Expected at least one path, got none")
	}

	// Verify endpoint exists
	if _, ok := contract.Paths["/api/v1/telemetry/events"]; !ok {
		t.Error("Expected path /api/v1/telemetry/events not found in contract")
	}
}

// TestContractSynthesizer_ApplySynthesisRules verifies conflict resolution
func TestContractSynthesizer_ApplySynthesisRules(t *testing.T) {
	agent := NewContractSynthesizer()

	// Create conflicting proposals
	proposals := []*synthesis.Proposal{
		synthesis.NewProposal(
			"frontend",
			EndpointProposal{Path: "/api/v1/telemetry/submit", Method: "POST"},
			0.9,
			"Frontend prefers this naming",
		),
		synthesis.NewProposal(
			"backend",
			EndpointProposal{Path: "/api/v1/telemetry/events", Method: "POST"},
			0.95,
			"Backend prefers this naming for consistency",
		),
		synthesis.NewProposal(
			"sdk",
			EndpointProposal{Path: "/api/v1/telemetry/events", Method: "POST"},
			0.85,
			"SDK agrees with backend",
		),
	}

	result, err := agent.ApplySynthesisRules(proposals)
	if err != nil {
		t.Fatalf("ApplySynthesisRules failed: %v", err)
	}

	// Backend should win (highest confidence + 2 agents agree)
	if result.Rule != "domain_expertise" {
		t.Errorf("Expected rule 'domain_expertise', got '%s'", result.Rule)
	}

	endpoint := result.Solution.(EndpointProposal)
	if endpoint.Path != "/api/v1/telemetry/events" {
		t.Errorf("Expected path '/api/v1/telemetry/events', got '%s'", endpoint.Path)
	}
}

// TestConflictResolution_UnanimousAgreement verifies unanimous rule
func TestConflictResolution_UnanimousAgreement(t *testing.T) {
	agent := NewContractSynthesizer()

	// All agents agree
	proposals := []*synthesis.Proposal{
		synthesis.NewProposal(
			"frontend",
			EndpointProposal{Path: "/api/v1/telemetry/events", Method: "POST"},
			0.9,
			"Agrees with backend",
		),
		synthesis.NewProposal(
			"backend",
			EndpointProposal{Path: "/api/v1/telemetry/events", Method: "POST"},
			0.95,
			"Proposes this endpoint",
		),
		synthesis.NewProposal(
			"sdk",
			EndpointProposal{Path: "/api/v1/telemetry/events", Method: "POST"},
			0.85,
			"Agrees with backend",
		),
	}

	result, err := agent.ApplySynthesisRules(proposals)
	if err != nil {
		t.Fatalf("ApplySynthesisRules failed: %v", err)
	}

	if result.Rule != "unanimous" {
		t.Errorf("Expected rule 'unanimous', got '%s'", result.Rule)
	}
}

// TestConflictResolution_DomainExpertiseVeto verifies veto power
func TestConflictResolution_DomainExpertiseVeto(t *testing.T) {
	agent := NewContractSynthesizer()

	// Backend has highest confidence (domain expert)
	proposals := []*synthesis.Proposal{
		synthesis.NewProposal(
			"frontend",
			EndpointProposal{Path: "/api/submit", Method: "POST"},
			0.7,
			"Simpler path",
		),
		synthesis.NewProposal(
			"backend",
			EndpointProposal{Path: "/api/v1/telemetry/events", Method: "POST"},
			0.98,
			"RESTful convention, versioned path",
		),
		synthesis.NewProposal(
			"sdk",
			EndpointProposal{Path: "/api/submit", Method: "POST"},
			0.75,
			"Simpler for SDK users",
		),
	}

	result, err := agent.ApplySynthesisRules(proposals)
	if err != nil {
		t.Fatalf("ApplySynthesisRules failed: %v", err)
	}

	if result.Rule != "domain_expertise" {
		t.Errorf("Expected rule 'domain_expertise', got '%s'", result.Rule)
	}

	if result.WinningAgent != "backend" {
		t.Errorf("Expected winning agent 'backend', got '%s'", result.WinningAgent)
	}
}

// TestContractSynthesizer_WriteContract verifies contract file writing
func TestContractSynthesizer_WriteContract(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test-contract.yaml")

	agent := NewContractSynthesizer()

	contract := &OpenAPIContract{
		OpenAPI: "3.0.0",
		Info: InfoSpec{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: PathsSpec{
			"/api/test": {
				"post": OperationSpec{
					Summary: "Test endpoint",
					Responses: ResponsesSpec{
						"200": ResponseSpec{
							Description: "Success",
						},
					},
				},
			},
		},
	}

	err := agent.WriteContract(contract, outputPath)
	if err != nil {
		t.Fatalf("WriteContract failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Contract file was not created")
	}

	// Verify file contains expected content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read contract file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "openapi: 3.0.0") {
		t.Error("Contract file missing openapi version")
	}

	if !strings.Contains(contentStr, "Test API") {
		t.Error("Contract file missing title")
	}
}

// TestContractSynthesizer_SynthesizeContract verifies end-to-end synthesis
func TestContractSynthesizer_SynthesizeContract(t *testing.T) {
	tmpDir := t.TempDir()
	reqPath := filepath.Join(tmpDir, "telemetry-requirements.md")
	outputPath := filepath.Join(tmpDir, "telemetry.yaml")

	reqContent := `# Telemetry Feature Requirements

## API Endpoints

### Submit Telemetry Event
- POST /api/v1/telemetry/events
- Request: {event_name, timestamp}
- Response: {event_id}
`
	if err := os.WriteFile(reqPath, []byte(reqContent), 0644); err != nil {
		t.Fatalf("Failed to create requirements file: %v", err)
	}

	agent := NewContractSynthesizer()
	result, err := agent.SynthesizeContract("telemetry", reqPath, outputPath)

	if err != nil {
		t.Fatalf("SynthesizeContract failed: %v", err)
	}

	if result.Rule == "" {
		t.Error("Expected synthesis rule to be set")
	}

	// Verify contract file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Contract file was not created")
	}
}

// Helper functions

func findEndpointByPath(endpoints []EndpointSpec, path string) *EndpointSpec {
	for i := range endpoints {
		if endpoints[i].Path == path {
			return &endpoints[i]
		}
	}
	return nil
}
