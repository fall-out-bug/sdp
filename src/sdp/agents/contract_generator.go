package agents

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaInferrer infers request/response schemas from code
type SchemaInferrer struct{}

// ContractGenerator enhances contracts with inferred schemas
type ContractGenerator struct {
	inferrer *SchemaInferrer
}

// NewSchemaInferrer creates a new schema inferrer
func NewSchemaInferrer() *SchemaInferrer {
	return &SchemaInferrer{}
}

// NewContractGenerator creates a new contract generator
func NewContractGenerator() *ContractGenerator {
	return &ContractGenerator{
		inferrer: NewSchemaInferrer(),
	}
}

// InferFromStruct extracts schema from Go struct definition
func (si *SchemaInferrer) InferFromStruct(structName, goCode string) (*SchemaSpec, error) {
	schema := &SchemaSpec{Fields: []FieldSpec{}}

	// Simplified struct parsing - look for struct and extract fields
	lines := strings.Split(goCode, "\n")
	inStruct := false
	structFound := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Find struct definition
		if strings.Contains(line, "type "+structName) && strings.Contains(line, "struct") {
			inStruct = true
			structFound = true
			continue
		}

		// End of struct
		if inStruct && line == "}" {
			break
		}

		// Extract fields (simplified - just Name Type)
		if inStruct && !strings.HasPrefix(line, "//") && line != "" && !strings.HasPrefix(line, "func") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				// Filter out non-field lines
				if len(parts[0]) > 0 && parts[0][0] >= 'A' && parts[0][0] <= 'Z' {
					fieldName := strings.ToLower(parts[0])
					goType := parts[1]

					schema.Fields = append(schema.Fields, FieldSpec{
						Name:     fieldName,
						Type:     mapGoTypeToJSON(goType),
						Required: true,
					})
				}
			}
		}
	}

	if !structFound {
		return nil, fmt.Errorf("struct %s not found", structName)
	}

	return schema, nil
}

// InferFromHandler extracts schema from handler function
func (si *SchemaInferrer) InferFromHandler(handlerName, goCode string) (*SchemaSpec, error) {
	schema := &SchemaSpec{
		Fields: []FieldSpec{},
	}

	// Find handler function - simplified implementation
	// Full implementation would parse handler body and find struct decode
	if !strings.Contains(goCode, "func "+handlerName) {
		return nil, fmt.Errorf("handler %s not found", handlerName)
	}

	return schema, nil
}

// InferFromTypeScript extracts schema from TypeScript interface
func (si *SchemaInferrer) InferFromTypeScript(interfaceName, tsCode string) (*SchemaSpec, error) {
	schema := &SchemaSpec{Fields: []FieldSpec{}}

	// Find interface definition
	lines := strings.Split(tsCode, "\n")
	inInterface := false
	interfaceFound := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Find interface definition
		if strings.Contains(line, "interface "+interfaceName) {
			inInterface = true
			interfaceFound = true
			continue
		}

		// End of interface
		if inInterface && line == "}" {
			break
		}

		// Extract fields
		if inInterface && !strings.HasPrefix(line, "//") && line != "" {
			// Format: fieldName: type;
			if strings.Contains(line, ":") && !strings.Contains(line, "function") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					fieldName := strings.TrimSpace(parts[0])
					tsType := strings.TrimSpace(parts[1])
					tsType = strings.TrimSuffix(tsType, ";")
					tsType = strings.TrimSpace(tsType)

					schema.Fields = append(schema.Fields, FieldSpec{
						Name:     fieldName,
						Type:     mapTSTypeToJSON(tsType),
						Required: !strings.HasSuffix(tsType, "?"),
					})
				}
			}
		}
	}

	if !interfaceFound {
		return nil, fmt.Errorf("interface %s not found", interfaceName)
	}

	return schema, nil
}

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

		// Add request body for POST/PUT/PATCH
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
	// Create a copy to avoid mutating the original
	enhanced := &OpenAPIContract{
		OpenAPI: contract.OpenAPI,
		Info:    contract.Info,
		Paths:   make(PathsSpec),
	}

	// Copy paths and enhance schemas
	for path, pathSpec := range contract.Paths {
		enhanced.Paths[path] = make(PathSpec)

		for method, operation := range pathSpec {
			enhancedOp := operation

			// Enhance request schema - recreate request body
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

			// Enhance response schemas - recreate responses
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

// Helper functions

// mapGoTypeToJSON maps Go types to JSON types
func mapGoTypeToJSON(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int32", "int64", "uint32", "uint64":
		return "integer"
	case "float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	default:
		if strings.HasPrefix(goType, "[]") {
			return "array"
		}
		if strings.HasPrefix(goType, "map[") {
			return "object"
		}
		return "object"
	}
}

// mapTSTypeToJSON maps TypeScript types to JSON types
func mapTSTypeToJSON(tsType string) string {
	// Remove optional marker
	tsType = strings.TrimSuffix(tsType, "?")
	tsType = strings.TrimSpace(tsType)

	switch tsType {
	case "string":
		return "string"
	case "number":
		return "number"
	case "boolean":
		return "boolean"
	default:
		if strings.HasPrefix(tsType, "Array<") {
			return "array"
		}
		if strings.HasPrefix(tsType, "Record<") {
			return "object"
		}
		return "object"
	}
}

// FindSchemaInCode finds and infers schema from code snippet
func (si *SchemaInferrer) FindSchemaInCode(codeSnippet, language, typeName string) (*SchemaSpec, error) {
	switch language {
	case "go":
		return si.InferFromStruct(typeName, codeSnippet)
	case "typescript", "javascript":
		return si.InferFromTypeScript(typeName, codeSnippet)
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
}

// ParseSchemaFromComment extracts schema from docstring comments
func (si *SchemaInferrer) ParseSchemaFromComment(comment string) (*SchemaSpec, error) {
	schema := &SchemaSpec{Fields: []FieldSpec{}}

	// Parse comment lines for schema hints
	// Format: @param field_name type description
	lines := strings.Split(comment, "\n")
	paramRe := regexp.MustCompile(`@param\s+(\w+)\s+(\w+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		matches := paramRe.FindStringSubmatch(line)
		if len(matches) >= 3 {
			schema.Fields = append(schema.Fields, FieldSpec{
				Name:     matches[1],
				Type:     matches[2],
				Required: true,
			})
		}
	}

	return schema, nil
}
