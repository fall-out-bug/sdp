package agents

import (
	"fmt"
	"strings"
)

// GenerateFromBackend generates a complete contract from backend routes
func (cg *ContractGenerator) GenerateFromBackend(
	componentName string,
	routes []ExtractedRoute,
) (*OpenAPIContract, error) {
	contract := &OpenAPIContract{
		OpenAPI: "3.0.0",
		Info: InfoSpec{
			Title:   fmt.Sprintf("%s API", strings.Title(componentName)),
			Version: "1.0.0",
		},
		Paths: make(PathsSpec),
	}

	for _, route := range routes {
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
							Schema: SchemaRefSpec{
								Type:       "object",
								Properties: map[string]PropertySpec{},
							},
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
						Schema: SchemaRefSpec{
							Type:       "object",
							Properties: map[string]PropertySpec{},
						},
					},
				},
			}
		}

		contract.Paths[route.Path][strings.ToLower(route.Method)] = operation
	}

	return contract, nil
}

// EnhanceContract enhances a contract with inferred schemas
func (cg *ContractGenerator) EnhanceContract(
	contract *OpenAPIContract,
	inferredSchemas map[string]SchemaSpec,
) (*OpenAPIContract, error) {
	enhanced := &OpenAPIContract{
		OpenAPI: contract.OpenAPI,
		Info:    contract.Info,
		Paths:   make(PathsSpec),
	}

	for path, pathSpec := range contract.Paths {
		enhanced.Paths[path] = make(PathSpec)

		for method, operation := range pathSpec {
			enhancedOp := operation

			if operation.RequestBody != nil {
				schemaKey := fmt.Sprintf("%s:%s:request", path, method)
				if schema, ok := inferredSchemas[schemaKey]; ok {
					newContent := make(map[string]MediaSpec)
					for mediaType := range operation.RequestBody.Content {
						newContent[mediaType] = MediaSpec{
							Schema: cg.schemaSpecToSchemaRef(schema),
						}
					}
					enhancedOp.RequestBody = &RequestSpec{
						Required: operation.RequestBody.Required,
						Content:  newContent,
					}
				}
			}

			newResponses := make(ResponsesSpec)
			for statusCode, response := range operation.Responses {
				schemaKey := fmt.Sprintf("%s:%s:response:%s", path, method, statusCode)
				if schema, ok := inferredSchemas[schemaKey]; ok {
					newContent := make(map[string]MediaSpec)
					for mediaType := range response.Content {
						newContent[mediaType] = MediaSpec{
							Schema: cg.schemaSpecToSchemaRef(schema),
						}
					}
					newResponses[statusCode] = ResponseSpec{
						Description: response.Description,
						Content:     newContent,
					}
				} else {
					newResponses[statusCode] = response
				}
			}
			enhancedOp.Responses = newResponses

			enhanced.Paths[path][method] = enhancedOp
		}
	}

	return enhanced, nil
}

// schemaSpecToSchemaRef converts a SchemaSpec to SchemaRefSpec
func (cg *ContractGenerator) schemaSpecToSchemaRef(schema SchemaSpec) SchemaRefSpec {
	properties := make(map[string]PropertySpec)
	required := []string{}

	for _, field := range schema.Fields {
		properties[field.Name] = PropertySpec{Type: field.Type}
		if field.Required {
			required = append(required, field.Name)
		}
	}

	return SchemaRefSpec{
		Type:       "object",
		Properties: properties,
		Required:   required,
	}
}
