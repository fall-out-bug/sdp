package collision

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

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
