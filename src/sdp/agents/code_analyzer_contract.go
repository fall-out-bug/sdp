package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// GenerateOpenAPIContract generates an OpenAPI contract from extracted data
func (ca *CodeAnalyzer) GenerateOpenAPIContract(
	componentName string,
	backendRoutes []ExtractedRoute,
	frontendCalls []ExtractedCall,
) (*OpenAPIContract, error) {
	contract := &OpenAPIContract{
		OpenAPI: "3.0.0",
		Info: InfoSpec{
			Title:   fmt.Sprintf("%s API", strings.Title(componentName)),
			Version: "1.0.0",
		},
		Paths: make(PathsSpec),
	}

	for _, route := range backendRoutes {
		if _, exists := contract.Paths[route.Path]; !exists {
			contract.Paths[route.Path] = make(PathSpec)
		}

		operation := OperationSpec{
			Summary: fmt.Sprintf("%s %s", route.Method, route.Path),
			Responses: ResponsesSpec{
				"200": ResponseSpec{
					Description: "Success",
					Content: map[string]MediaSpec{
						"application/json": {
							Schema: SchemaRefSpec{Type: "object"},
						},
					},
				},
			},
		}

		if route.Method == "POST" || route.Method == "PUT" || route.Method == "PATCH" {
			operation.RequestBody = &RequestSpec{
				Required: true,
				Content: map[string]MediaSpec{
					"application/json": {
						Schema: SchemaRefSpec{Type: "object"},
					},
				},
			}
		}

		contract.Paths[route.Path][strings.ToLower(route.Method)] = operation
	}

	return contract, nil
}

// WriteExtractedContract writes the extracted contract to a YAML file
func (ca *CodeAnalyzer) WriteExtractedContract(contract *OpenAPIContract, outputPath string) error {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := yaml.Marshal(contract)
	if err != nil {
		return fmt.Errorf("failed to marshal contract: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write contract: %w", err)
	}

	return nil
}

// ExtractComponentContract performs end-to-end contract extraction
func (ca *CodeAnalyzer) ExtractComponentContract(
	componentName string,
	componentType string,
	filePaths []string,
	outputPath string,
) error {
	var contract *OpenAPIContract
	var err error

	switch componentType {
	case "backend":
		var allRoutes []ExtractedRoute
		for _, filePath := range filePaths {
			routes, err := ca.AnalyzeGoBackend(filePath)
			if err != nil {
				return fmt.Errorf("failed to analyze %s: %w", filePath, err)
			}
			allRoutes = append(allRoutes, routes...)
		}
		contract, err = ca.GenerateOpenAPIContract(componentName, allRoutes, nil)
		if err != nil {
			return fmt.Errorf("failed to generate contract: %w", err)
		}

	case "frontend":
		var allCalls []ExtractedCall
		for _, filePath := range filePaths {
			calls, err := ca.AnalyzeTypeScriptFrontend(filePath)
			if err != nil {
				return fmt.Errorf("failed to analyze %s: %w", filePath, err)
			}
			allCalls = append(allCalls, calls...)
		}
		var routes []ExtractedRoute
		for _, call := range allCalls {
			routes = append(routes, ExtractedRoute{
				Path:   call.Path,
				Method: call.Method,
				File:   call.File,
				Line:   call.Line,
			})
		}
		contract, err = ca.GenerateOpenAPIContract(componentName, routes, nil)
		if err != nil {
			return fmt.Errorf("failed to generate contract: %w", err)
		}

	case "sdk":
		contract = &OpenAPIContract{
			OpenAPI: "3.0.0",
			Info: InfoSpec{
				Title:   fmt.Sprintf("%s SDK", strings.Title(componentName)),
				Version: "1.0.0",
			},
			Paths: make(PathsSpec),
		}

	default:
		return fmt.Errorf("unknown component type: %s", componentType)
	}

	if err := ca.WriteExtractedContract(contract, outputPath); err != nil {
		return fmt.Errorf("failed to write contract: %w", err)
	}

	return nil
}
