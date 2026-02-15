package agents

import (
	"fmt"
	"regexp"
	"strings"
)

// InferFromStruct extracts schema from Go struct definition
func (si *SchemaInferrer) InferFromStruct(structName, goCode string) (*SchemaSpec, error) {
	schema := &SchemaSpec{Fields: []FieldSpec{}}

	lines := strings.Split(goCode, "\n")
	inStruct := false
	structFound := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "type "+structName) && strings.Contains(line, "struct") {
			inStruct = true
			structFound = true
			continue
		}

		if inStruct && line == "}" {
			break
		}

		if inStruct && !strings.HasPrefix(line, "//") && line != "" && !strings.HasPrefix(line, "func") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
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
	schema := &SchemaSpec{Fields: []FieldSpec{}}

	if !strings.Contains(goCode, "func "+handlerName) {
		return nil, fmt.Errorf("handler %s not found", handlerName)
	}

	return schema, nil
}

// InferFromTypeScript extracts schema from TypeScript interface
func (si *SchemaInferrer) InferFromTypeScript(interfaceName, tsCode string) (*SchemaSpec, error) {
	schema := &SchemaSpec{Fields: []FieldSpec{}}

	lines := strings.Split(tsCode, "\n")
	inInterface := false
	interfaceFound := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "interface "+interfaceName) {
			inInterface = true
			interfaceFound = true
			continue
		}

		if inInterface && line == "}" {
			break
		}

		if inInterface && !strings.HasPrefix(line, "//") && line != "" {
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
