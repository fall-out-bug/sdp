package agents

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGoBackendAnalyzer_ExtractRoutes verifies route extraction from Go code
func TestGoBackendAnalyzer_ExtractRoutes(t *testing.T) {
	// Create temporary Go file with route definitions
	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "routes.go")
	goCode := `package main

import (
	"net/http"
	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/telemetry/events", handleTelemetryEvents).Methods("POST")
	r.HandleFunc("/api/v1/telemetry/events/{id}", handleTelemetryEvent).Methods("GET")
	r.HandleFunc("/api/v1/health", handleHealth).Methods("GET")
}

func handleTelemetryEvents(w http.ResponseWriter, r *http.Request) {
	// Handler implementation
}

func handleTelemetryEvent(w http.ResponseWriter, r *http.Request) {
	// Handler implementation
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	// Handler implementation
}
`
	if err := os.WriteFile(goFile, []byte(goCode), 0644); err != nil {
		t.Fatalf("Failed to create Go file: %v", err)
	}

	analyzer := NewCodeAnalyzer()
	routes, err := analyzer.AnalyzeGoBackend(goFile)
	if err != nil {
		t.Fatalf("AnalyzeGoBackend failed: %v", err)
	}

	// Verify 3 routes extracted
	if len(routes) != 3 {
		t.Errorf("Expected 3 routes, got %d", len(routes))
	}

	// Verify POST route
	postRoute := findRouteByPath(routes, "/api/v1/telemetry/events")
	if postRoute == nil {
		t.Fatal("POST route not found")
	}
	if postRoute.Method != "POST" {
		t.Errorf("Expected method POST, got %s", postRoute.Method)
	}
	if postRoute.File != goFile {
		t.Errorf("Expected file %s, got %s", goFile, postRoute.File)
	}
}

// TestGoBackendAnalyzer_GinFramework verifies gin framework support
func TestGoBackendAnalyzer_GinFramework(t *testing.T) {
	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "gin_routes.go")
	goCode := `package main

import (
	"github.com/gin-gonic/gin"
)

func SetupGinRoutes(r *gin.Engine) {
	r.POST("/api/v1/telemetry/events", handleTelemetryEvents)
	r.GET("/api/v1/telemetry/events/:id", handleTelemetryEvent)
}

func handleTelemetryEvents(c *gin.Context) {
	// Handler implementation
}

func handleTelemetryEvent(c *gin.Context) {
	// Handler implementation
}
`
	if err := os.WriteFile(goFile, []byte(goCode), 0644); err != nil {
		t.Fatalf("Failed to create Go file: %v", err)
	}

	analyzer := NewCodeAnalyzer()
	routes, err := analyzer.AnalyzeGoBackend(goFile)
	if err != nil {
		t.Fatalf("AnalyzeGoBackend failed: %v", err)
	}

	if len(routes) != 2 {
		t.Errorf("Expected 2 routes, got %d", len(routes))
	}
}

// TestTypeScriptFrontendAnalyzer_ExtractFetchCalls verifies fetch call extraction
func TestTypeScriptFrontendAnalyzer_ExtractFetchCalls(t *testing.T) {
	tmpDir := t.TempDir()
	tsFile := filepath.Join(tmpDir, "api.ts")
	tsCode := `// API client functions

export async function submitTelemetry(event: TelemetryEvent): Promise<Response> {
	return fetch("/api/v1/telemetry/events", {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify(event)
	});
}

export async function getTelemetryEvent(id: string): Promise<TelemetryEvent> {
	return fetch("/api/v1/telemetry/events/" + id)
		.then(res => res.json());
}
`
	if err := os.WriteFile(tsFile, []byte(tsCode), 0644); err != nil {
		t.Fatalf("Failed to create TypeScript file: %v", err)
	}

	analyzer := NewCodeAnalyzer()
	calls, err := analyzer.AnalyzeTypeScriptFrontend(tsFile)
	if err != nil {
		t.Fatalf("AnalyzeTypeScriptFrontend failed: %v", err)
	}

	if len(calls) != 2 {
		t.Errorf("Expected 2 API calls, got %d", len(calls))
	}

	// Verify POST call
	postCall := findCallByPath(calls, "/api/v1/telemetry/events")
	if postCall == nil {
		t.Fatal("POST call not found")
	}
	if postCall.Method != "POST" {
		t.Errorf("Expected method POST, got %s", postCall.Method)
	}
}

// TestTypeScriptFrontendAnalyzer_AxiosCalls verifies axios call extraction
func TestTypeScriptFrontendAnalyzer_AxiosCalls(t *testing.T) {
	tmpDir := t.TempDir()
	tsFile := filepath.Join(tmpDir, "axios_api.ts")
	tsCode := `import axios from 'axios';

export async function submitTelemetry(event: TelemetryEvent) {
	return axios.post("/api/v1/telemetry/events", event);
}

export async function getTelemetryEvent(id: string) {
	return axios.get("/api/v1/telemetry/events/" + id);
}
`
	if err := os.WriteFile(tsFile, []byte(tsCode), 0644); err != nil {
		t.Fatalf("Failed to create TypeScript file: %v", err)
	}

	analyzer := NewCodeAnalyzer()
	calls, err := analyzer.AnalyzeTypeScriptFrontend(tsFile)
	if err != nil {
		t.Fatalf("AnalyzeTypeScriptFrontend failed: %v", err)
	}

	if len(calls) != 2 {
		t.Errorf("Expected 2 API calls, got %d", len(calls))
	}
}

// TestPythonSDKAnalyzer_ExtractPublicMethods verifies public method extraction
func TestPythonSDKAnalyzer_ExtractPublicMethods(t *testing.T) {
	tmpDir := t.TempDir()
	pyFile := filepath.Join(tmpDir, "client.py")
	pyCode := `"""Telemetry SDK Client"""

class TelemetryClient:
    """Client for telemetry API"""

    def __init__(self, api_key: str):
        self.api_key = api_key

    def submit_event(self, event_name: str, timestamp: str, metadata: dict) -> dict:
        """Submit a telemetry event"""
        return {"event_id": "123", "status": "submitted"}

    def get_event(self, event_id: str) -> dict:
        """Get a telemetry event by ID"""
        return {"event_id": event_id, "event_name": "test"}

    def _private_method(self):
        """Private method - should not be extracted"""
        pass
`
	if err := os.WriteFile(pyFile, []byte(pyCode), 0644); err != nil {
		t.Fatalf("Failed to create Python file: %v", err)
	}

	analyzer := NewCodeAnalyzer()
	methods, err := analyzer.AnalyzePythonSDK(pyFile)
	if err != nil {
		t.Fatalf("AnalyzePythonSDK failed: %v", err)
	}

	// Should extract 2 public methods (not _private_method)
	if len(methods) != 2 {
		t.Errorf("Expected 2 public methods, got %d", len(methods))
	}

	// Verify submit_event method
	submitMethod := findMethodByName(methods, "submit_event")
	if submitMethod == nil {
		t.Fatal("submit_event method not found")
	}
	if submitMethod.Name != "submit_event" {
		t.Errorf("Expected method name 'submit_event', got '%s'", submitMethod.Name)
	}
}

// TestGenerateOpenAPIContract verifies OpenAPI contract generation
func TestGenerateOpenAPIContract(t *testing.T) {
	analyzer := NewCodeAnalyzer()

	// Create sample extracted data
	backendRoutes := []ExtractedRoute{
		{Path: "/api/v1/telemetry/events", Method: "POST", File: "routes.go"},
		{Path: "/api/v1/telemetry/events/{id}", Method: "GET", File: "routes.go"},
	}

	frontendCalls := []ExtractedCall{
		{Path: "/api/v1/telemetry/events", Method: "POST", File: "api.ts"},
		{Path: "/api/v1/telemetry/events/{id}", Method: "GET", File: "api.ts"},
	}

	contract, err := analyzer.GenerateOpenAPIContract("telemetry", backendRoutes, frontendCalls)
	if err != nil {
		t.Fatalf("GenerateOpenAPIContract failed: %v", err)
	}

	// Verify OpenAPI structure
	if contract.OpenAPI != "3.0.0" {
		t.Errorf("Expected OpenAPI version 3.0.0, got %s", contract.OpenAPI)
	}

	if contract.Info.Title != "Telemetry API" {
		t.Errorf("Expected title 'Telemetry API', got %s", contract.Info.Title)
	}

	if len(contract.Paths) != 2 {
		t.Errorf("Expected 2 paths, got %d", len(contract.Paths))
	}

	// Verify path exists
	if _, ok := contract.Paths["/api/v1/telemetry/events"]; !ok {
		t.Error("Expected path /api/v1/telemetry/events not found")
	}
}

// TestWriteExtractedContract verifies contract file writing
func TestWriteExtractedContract(t *testing.T) {
	analyzer := NewCodeAnalyzer()

	contract := &OpenAPIContract{
		OpenAPI: "3.0.0",
		Info: InfoSpec{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: make(PathsSpec),
	}

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test-contract.yaml")

	err := analyzer.WriteExtractedContract(contract, outputPath)
	if err != nil {
		t.Fatalf("WriteExtractedContract failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatal("Contract file was not created")
	}

	// Verify file content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read contract file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Contract file is empty")
	}

	// Verify it's valid YAML
	if !stringContains(content, "openapi: 3.0.0") {
		t.Error("Contract file missing openapi version")
	}
}

// Helper functions

func findRouteByPath(routes []ExtractedRoute, path string) *ExtractedRoute {
	for i := range routes {
		if routes[i].Path == path {
			return &routes[i]
		}
	}
	return nil
}

func findCallByPath(calls []ExtractedCall, path string) *ExtractedCall {
	for i := range calls {
		if calls[i].Path == path {
			return &calls[i]
		}
	}
	return nil
}

func findMethodByName(methods []ExtractedMethod, name string) *ExtractedMethod {
	for i := range methods {
		if methods[i].Name == name {
			return &methods[i]
		}
	}
	return nil
}

func stringContains(content []byte, substr string) bool {
	return len(content) > 0 && stringContainsStr(string(content), substr)
}

func stringContainsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
