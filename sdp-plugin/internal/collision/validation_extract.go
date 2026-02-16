package collision

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

// extractImplFields extracts struct fields from a Go implementation file.
func extractImplFields(path, typeName string) (map[string]string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse implementation: %w", err)
	}

	fields := make(map[string]string)
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
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			extractStructTypeFields(structType, fields)
			break
		}
	}
	return fields, nil
}

// extractStructTypeFields extracts fields from a struct type into the provided map.
// Named to avoid conflict with extractStructFields in boundary.go.
func extractStructTypeFields(structType *ast.StructType, fields map[string]string) {
	for _, field := range structType.Fields.List {
		fieldType := formatFieldType(field.Type)
		addFieldNames(field.Names, fieldType, fields)
	}
}

// addFieldNames adds field names with their type to the fields map.
func addFieldNames(names []*ast.Ident, fieldType string, fields map[string]string) {
	for _, name := range names {
		fields[name.Name] = fieldType
	}
}

// ValidateContractsInDir validates all contracts in a directory.
func ValidateContractsInDir(contractsDir, implDir string) ([]Violation, error) {
	var allViolations []Violation

	entries, err := os.ReadDir(contractsDir)
	if err != nil {
		return nil, fmt.Errorf("read contracts dir: %w", err)
	}

	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".yaml" {
			continue
		}

		contractPath := filepath.Join(contractsDir, e.Name())
		typeName := filepath.Base(e.Name())
		typeName = typeName[:len(typeName)-5] // Remove .yaml

		// Find implementation file
		implPath := filepath.Join(implDir, typeName+".go")
		implFound := true
		if _, err := os.Stat(implPath); os.IsNotExist(err) {
			// Try lowercase
			implPath = filepath.Join(implDir, toLowerFirst(typeName)+".go")
			if _, err := os.Stat(implPath); os.IsNotExist(err) {
				// Implementation file missing - add violation (bug fix for sdp-1lqm)
				implFound = false
			}
		}

		if !implFound {
			// Add violation for missing implementation file
			allViolations = append(allViolations, Violation{
				Type:     "missing_implementation",
				Field:    typeName,
				Expected: typeName + ".go or " + toLowerFirst(typeName) + ".go",
				Actual:   "not found",
				Severity: "error",
				Message:  fmt.Sprintf("Implementation file not found for contract %s: expected %s or %s in %s", typeName, typeName+".go", toLowerFirst(typeName)+".go", implDir),
			})
			continue
		}

		violations, err := ValidateContractAgainstImpl(contractPath, implPath)
		if err != nil {
			return nil, fmt.Errorf("validate %s: %w", typeName, err)
		}

		allViolations = append(allViolations, violations...)
	}

	return allViolations, nil
}

func toLowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	// Check if already lowercase (a-z)
	if s[0] >= 'a' && s[0] <= 'z' {
		return s
	}
	// Convert uppercase (A-Z) to lowercase
	if s[0] >= 'A' && s[0] <= 'Z' {
		return string(s[0]+32) + s[1:]
	}
	return s
}
