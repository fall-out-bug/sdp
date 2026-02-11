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
	FileName       string      `json:"fileName"`
	TypeName       string      `json:"typeName"`
	Fields         []FieldInfo `json:"fields"`
	Features       []string    `json:"features"`
	Recommendation string      `json:"recommendation"`
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

		// Resolve relative path before parsing (bug fix for sdp-zidp)
		resolvedFile := resolveFilePath(file)

		// Parse Go file to extract types
		types, err := extractGoTypes(resolvedFile)
		if err != nil {
			continue // Skip files that can't be parsed
		}

		for _, typeName := range types {
			fields := extractStructFields(resolvedFile, typeName)
			boundaries = append(boundaries, SharedBoundary{
				FileName:       file, // Store original path
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
// Deduplicates feature IDs to avoid same-feature false positives.
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
			// Deduplicate: only add if featureID not already present
			featureIDs := fileToFeatures[file]
			if !stringSliceContains(featureIDs, f.FeatureID) {
				fileToFeatures[file] = append(featureIDs, f.FeatureID)
			}
		}
	}
	return fileToFeatures
}

// stringSliceContains checks if a string slice contains a value.
func stringSliceContains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// resolveFilePath resolves a relative file path to an absolute one.
// Tries common base directories relative to current working directory.
func resolveFilePath(filePath string) string {
	// If already absolute, return as-is
	if filepath.IsAbs(filePath) {
		return filePath
	}

	// Try to find the file relative to current directory or testdata
	for _, base := range []string{".", "testdata", "..", "../.."} {
		candidate := filepath.Join(base, filePath)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	// If not found, return original (will fail during parse)
	return filePath
}

// extractGoTypes parses a Go file and returns type names defined in it.
// Note: filePath should already be resolved via resolveFilePath().
func extractGoTypes(filePath string) ([]string, error) {
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
// Note: filePath should already be resolved via resolveFilePath().
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
