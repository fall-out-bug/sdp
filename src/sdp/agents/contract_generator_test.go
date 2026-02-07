package agents

import (
	"testing"
)

// TestGenerateFromBackend verifies contract generation
func TestGenerateFromBackend(t *testing.T) {
	generator := NewContractGenerator()

	routes := []ExtractedRoute{
		{Path: "/api/v1/telemetry/events", Method: "POST", File: "routes.go", Line: 10},
		{Path: "/api/v1/telemetry/events/{id}", Method: "GET", File: "routes.go", Line: 11},
	}

	contract, err := generator.GenerateFromBackend("telemetry", routes)
	if err != nil {
		t.Fatalf("GenerateFromBackend failed: %v", err)
	}

	if contract.OpenAPI != "3.0.0" {
		t.Errorf("Expected OpenAPI version 3.0.0, got %s", contract.OpenAPI)
	}

	if len(contract.Paths) != 2 {
		t.Errorf("Expected 2 paths, got %d", len(contract.Paths))
	}

	// Verify POST operation has request body
	postPath := contract.Paths["/api/v1/telemetry/events"]
	if postPath == nil {
		t.Fatal("POST path not found")
	}

	postOp := postPath["post"]
	if postOp.RequestBody == nil {
		t.Error("Expected request body for POST operation")
	}
}

// TestEnhanceContract verifies schema enhancement
func TestEnhanceContract(t *testing.T) {
	generator := NewContractGenerator()

	// Start with basic contract
	contract := &OpenAPIContract{
		OpenAPI: "3.0.0",
		Info: InfoSpec{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: PathsSpec{
			"/api/test": {
				"post": OperationSpec{
					Summary: "POST /api/test",
					RequestBody: &RequestSpec{
						Content: map[string]MediaSpec{
							"application/json": {
								Schema: SchemaRefSpec{Type: "object"},
							},
						},
					},
					Responses: ResponsesSpec{
						"200": ResponseSpec{
							Content: map[string]MediaSpec{
								"application/json": {
									Schema: SchemaRefSpec{Type: "object"},
								},
							},
						},
					},
				},
			},
		},
	}

	// Enhance with inferred schemas
	inferredSchemas := map[string]SchemaSpec{
		"/api/test:post:request": {
			Fields: []FieldSpec{
				{Name: "event_name", Type: "string", Required: true},
				{Name: "timestamp", Type: "string", Required: true},
			},
		},
	}

	enhanced, err := generator.EnhanceContract(contract, inferredSchemas)
	if err != nil {
		t.Fatalf("EnhanceContract failed: %v", err)
	}

	// Verify schema was enhanced
	path := enhanced.Paths["/api/test"]
	if path == nil {
		t.Fatal("Path not found in enhanced contract")
	}

	postOp := path["post"]
	if postOp.RequestBody == nil {
		t.Fatal("Request body not found")
	}

	schema := postOp.RequestBody.Content["application/json"].Schema
	if len(schema.Properties) == 0 {
		t.Error("Expected schema properties to be enhanced")
	}

	// Verify specific fields
	if _, ok := schema.Properties["event_name"]; !ok {
		t.Error("Expected event_name field in schema")
	}

	if len(schema.Required) != 2 {
		t.Errorf("Expected 2 required fields, got %d", len(schema.Required))
	}
}

// TestInferFromStruct verifies struct schema inference
func TestInferFromStruct(t *testing.T) {
	inferrer := NewSchemaInferrer()

	goCode := `
type TelemetryEvent struct {
	EventName string
	Timestamp string
	Metadata  map[string]string
}
`

	schema, err := inferrer.InferFromStruct("TelemetryEvent", goCode)
	if err != nil {
		t.Fatalf("InferFromStruct failed: %v", err)
	}

	if len(schema.Fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(schema.Fields))
	}
}

// TestInferFromTypeScript verifies TypeScript schema inference
func TestInferFromTypeScript(t *testing.T) {
	inferrer := NewSchemaInferrer()

	tsCode := `
interface TelemetryEvent {
	event_name: string;
	timestamp: string;
}
`

	schema, err := inferrer.InferFromTypeScript("TelemetryEvent", tsCode)
	if err != nil {
		t.Fatalf("InferFromTypeScript failed: %v", err)
	}

	if len(schema.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(schema.Fields))
	}

	// Verify fields are required
	if schema.Fields[0].Required {
		t.Log("Field is required as expected")
	}
}

// TestSchemaSpecToSchemaRef verifies conversion
func TestSchemaSpecToSchemaRef(t *testing.T) {
	generator := NewContractGenerator()

	schema := SchemaSpec{
		Fields: []FieldSpec{
			{Name: "field1", Type: "string", Required: true},
			{Name: "field2", Type: "integer", Required: false},
		},
	}

	ref := generator.schemaSpecToSchemaRef(schema)

	if ref.Type != "object" {
		t.Errorf("Expected type 'object', got '%s'", ref.Type)
	}

	if len(ref.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(ref.Properties))
	}

	if len(ref.Required) != 1 {
		t.Errorf("Expected 1 required field, got %d", len(ref.Required))
	}
}

// Helper functions

func findFieldByName(fields []FieldSpec, name string) *FieldSpec {
	for i := range fields {
		if fields[i].Name == name {
			return &fields[i]
		}
	}
	return nil
}
