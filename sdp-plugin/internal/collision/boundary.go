package collision

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// FeatureScope represents a feature's scope for boundary detection.
type FeatureScope struct {
	FeatureID  string
	ScopeFiles []string
}

// SharedBoundary represents a shared type/interface boundary between features.
type SharedBoundary struct {
	FileName       string     `json:"fileName"`
	TypeName       string     `json:"typeName"`
	Fields         []FieldInfo `json:"fields"`
	Features       []string   `json:"features"`
	Recommendation string     `json:"recommendation"`
}

// FieldInfo represents a field in a struct or interface.
type FieldInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// DetectBoundaries analyzes scope files to identify shared boundaries (types/interfaces).
func DetectBoundaries(features []FeatureScope) []SharedBoundary {
	fileToFeatures := buildFileToFeatures(features)
	var boundaries []SharedBoundary

	for file, featureIDs := range fileToFeatures {
		if len(featureIDs) < 2 {
			continue // Not a shared boundary
		}

		// Parse Go file to extract types
		types, err := extractGoTypes(file)
		if err != nil {
			continue // Skip files that can't be parsed
		}

		for _, typeName := range types {
			fields := extractStructFields(file, typeName)
			boundaries = append(boundaries, SharedBoundary{
				FileName:       file,
				TypeName:       typeName,
				Fields:         fields,
				Features:       featureIDs,
				Recommendation: "Create shared interface contract",
			})
		}
	}

	return boundaries
}

// buildFileToFeatures maps files to the features that use them.
func buildFileToFeatures(features []FeatureScope) map[string][]string {
	fileToFeatures := make(map[string][]string)
	for _, f := range features {
		for _, file := range f.ScopeFiles {
			file = normalizePath(file)
			if file == "" {
				continue
			}
			// Check if it's a Go file
			if !strings.HasSuffix(file, ".go") {
				continue
			}
			fileToFeatures[file] = append(fileToFeatures[file], f.FeatureID)
		}
	}
	return fileToFeatures
}

// extractGoTypes parses a Go file and returns type names defined in it.
func extractGoTypes(filePath string) ([]string, error) {
	// Resolve relative path
	if !filepath.IsAbs(filePath) {
		// Try to find the file relative to current directory or testdata
		for _, base := range []string{".", "testdata", "..", "../.."} {
			if _, err := os.Stat(filepath.Join(base, filePath)); err == nil {
				filePath = filepath.Join(base, filePath)
				break
			}
		}
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var types []string
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			types = append(types, typeSpec.Name.Name)
		}
	}
	return types, nil
}

// extractStructFields extracts field information from a struct type.
func extractStructFields(filePath, typeName string) []FieldInfo {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil
	}

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != typeName {
				continue
			}
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				return extractFieldsFromStruct(structType)
			}
		}
	}
	return nil
}

// extractFieldsFromStruct extracts field names and types from a struct.
func extractFieldsFromStruct(structType *ast.StructType) []FieldInfo {
	if structType.Fields == nil {
		return nil
	}
	var fields []FieldInfo
	for _, field := range structType.Fields.List {
		fieldType := formatFieldType(field.Type)
		for _, name := range field.Names {
			fields = append(fields, FieldInfo{
				Name: name.Name,
				Type: fieldType,
			})
		}
		// Handle anonymous fields
		if len(field.Names) == 0 {
			fields = append(fields, FieldInfo{
				Name: "",
				Type: fieldType,
			})
		}
	}
	return fields
}

// formatFieldType converts an AST expression to a string representation.
func formatFieldType(expr ast.Expr) string {
	if expr == nil {
		return "any"
	}
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return formatFieldType(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + formatFieldType(t.X)
	case *ast.ArrayType:
		return "[]" + formatFieldType(t.Elt)
	case *ast.MapType:
		return "map[" + formatFieldType(t.Key) + "]" + formatFieldType(t.Value)
	default:
		return "any"
	}
}

// BoundaryToJSON converts a SharedBoundary to JSON string.
func BoundaryToJSON(boundary SharedBoundary) (string, error) {
	data, err := json.MarshalIndent(boundary, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
