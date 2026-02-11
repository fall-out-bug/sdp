package collision

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Violation represents a contract validation issue.
type Violation struct {
	Type     string `json:"type"`     // missing_field, type_mismatch, extra_field
	Field    string `json:"field"`    // Field name
	Expected string `json:"expected"` // Expected type/value
	Actual   string `json:"actual"`   // Actual type/value
	Severity string `json:"severity"` // error, warning
	Message  string `json:"message"`  // Human-readable message
}

// ValidateContractAgainstImpl validates implementation against a contract.
func ValidateContractAgainstImpl(contractPath, implPath string) ([]Violation, error) {
	contract, err := loadContract(contractPath)
	if err != nil {
		return nil, err
	}

	implFields, err := extractImplFields(implPath, contract.TypeName)
	if err != nil {
		return nil, err
	}

	contractFields := buildContractFieldMap(contract.Fields)
	return compareFields(contractFields, implFields), nil
}

// loadContract loads and parses a contract YAML file.
func loadContract(path string) (Contract, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Contract{}, fmt.Errorf("read contract: %w", err)
	}
	var contract Contract
	if err := yaml.Unmarshal(data, &contract); err != nil {
		return Contract{}, fmt.Errorf("parse contract: %w", err)
	}
	return contract, nil
}

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

// buildContractFieldMap converts field slice to map.
func buildContractFieldMap(fields []FieldInfo) map[string]string {
	m := make(map[string]string)
	for _, f := range fields {
		m[f.Name] = f.Type
	}
	return m
}

// compareFields compares contract fields with implementation fields.
func compareFields(contractFields, implFields map[string]string) []Violation {
	var violations []Violation

	// Check for missing required fields and type mismatches
	for name, cType := range contractFields {
		if iType, exists := implFields[name]; !exists {
			violations = append(violations, Violation{
				Type:     "missing_field",
				Field:    name,
				Expected: cType,
				Actual:   "missing",
				Severity: "error",
				Message:  fmt.Sprintf("Missing required field: %s (%s)", name, cType),
			})
		} else if iType != cType {
			violations = append(violations, Violation{
				Type:     "type_mismatch",
				Field:    name,
				Expected: cType,
				Actual:   iType,
				Severity: "warning",
				Message:  fmt.Sprintf("Type mismatch for %s: expected %s, got %s", name, cType, iType),
			})
		}
	}

	// Check for extra fields (warning only in P1)
	for name, fType := range implFields {
		if _, exists := contractFields[name]; !exists {
			violations = append(violations, Violation{
				Type:     "extra_field",
				Field:    name,
				Expected: "not in contract",
				Actual:   fType,
				Severity: "warning",
				Message:  fmt.Sprintf("Extra field not in contract: %s (%s)", name, fType),
			})
		}
	}

	return violations
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
