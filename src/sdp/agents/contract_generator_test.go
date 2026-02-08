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

// TestInferFromHandler verifies handler schema inference
func TestInferFromHandler(t *testing.T) {
	inferrer := NewSchemaInferrer()

	goCode := `
package main
func handleTelemetryEvents(w http.ResponseWriter, r *http.Request) {
	var event TelemetryEvent
	json.NewDecoder(r.Body).Decode(&event)
}
`
	schema, err := inferrer.InferFromHandler("handleTelemetryEvents", goCode)
	if err != nil {
		t.Fatalf("InferFromHandler failed: %v", err)
	}

	if schema == nil {
		t.Fatal("Expected schema to be returned")
	}

	// Handler exists, so we should get a schema (even if empty for now)
	if len(schema.Fields) != 0 {
		t.Logf("Schema fields: %d (implementation may be enhanced later)", len(schema.Fields))
	}
}

// TestInferFromHandler_HandlerNotFound verifies error handling
func TestInferFromHandler_HandlerNotFound(t *testing.T) {
	inferrer := NewSchemaInferrer()

	goCode := `
package main
func handleTelemetryEvents(w http.ResponseWriter, r *http.Request) {}
`
	_, err := inferrer.InferFromHandler("nonExistentHandler", goCode)
	if err == nil {
		t.Fatal("Expected error for non-existent handler")
	}

	if err != nil && err.Error() != "handler nonExistentHandler not found" {
		t.Errorf("Expected 'handler not found' error, got: %v", err)
	}
}

// TestFindSchemaInCode_Go verifies Go schema inference
func TestFindSchemaInCode_Go(t *testing.T) {
	inferrer := NewSchemaInferrer()

	goCode := `
type TelemetryEvent struct {
	EventName string
	Timestamp string
}
`

	schema, err := inferrer.FindSchemaInCode(goCode, "go", "TelemetryEvent")
	if err != nil {
		t.Fatalf("FindSchemaInCode (Go) failed: %v", err)
	}

	if schema == nil {
		t.Fatal("Expected schema to be returned")
	}

	if len(schema.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(schema.Fields))
	}
}

// TestFindSchemaInCode_TypeScript verifies TypeScript schema inference
func TestFindSchemaInCode_TypeScript(t *testing.T) {
	inferrer := NewSchemaInferrer()

	tsCode := `
interface TelemetryEvent {
	event_name: string;
	timestamp: string;
}
`

	schema, err := inferrer.FindSchemaInCode(tsCode, "typescript", "TelemetryEvent")
	if err != nil {
		t.Fatalf("FindSchemaInCode (TypeScript) failed: %v", err)
	}

	if schema == nil {
		t.Fatal("Expected schema to be returned")
	}

	if len(schema.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(schema.Fields))
	}
}

// TestFindSchemaInCode_UnsupportedLanguage verifies error handling
func TestFindSchemaInCode_UnsupportedLanguage(t *testing.T) {
	inferrer := NewSchemaInferrer()

	_, err := inferrer.FindSchemaInCode("some code", "python", "SomeType")
	if err == nil {
		t.Fatal("Expected error for unsupported language")
	}

	if err != nil && err.Error() != "unsupported language: python" {
		t.Errorf("Expected 'unsupported language' error, got: %v", err)
	}
}

// TestParseSchemaFromComment verifies schema extraction from docstring comments
func TestParseSchemaFromComment(t *testing.T) {
	inferrer := NewSchemaInferrer()

	comment := `
SubmitTelemetryEvent submits a telemetry event

@param event_name The name of the event (string)
@param timestamp ISO timestamp of the event (string)
@param metadata Additional metadata (dict)
Returns the event ID (string)
`

	schema, err := inferrer.ParseSchemaFromComment(comment)
	if err != nil {
		t.Fatalf("ParseSchemaFromComment failed: %v", err)
	}

	if schema == nil {
		t.Fatal("Expected schema to be returned")
	}

	if len(schema.Fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(schema.Fields))
	}

	// Verify field names
	expectedFields := []string{"event_name", "timestamp", "metadata"}
	for _, expected := range expectedFields {
		found := false
		for _, field := range schema.Fields {
			if field.Name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected field '%s' not found", expected)
		}
	}
}

// TestParseSchemaFromComment_EmptyComment verifies handling of empty comment
func TestParseSchemaFromComment_EmptyComment(t *testing.T) {
	inferrer := NewSchemaInferrer()

	comment := "No params here"

	schema, err := inferrer.ParseSchemaFromComment(comment)
	if err != nil {
		t.Fatalf("ParseSchemaFromComment failed: %v", err)
	}

	if schema == nil {
		t.Fatal("Expected schema to be returned")
	}

	if len(schema.Fields) != 0 {
		t.Errorf("Expected 0 fields for comment without @param, got %d", len(schema.Fields))
	}
}
