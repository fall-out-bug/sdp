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
	// Load contract
	contractData, err := os.ReadFile(contractPath)
	if err != nil {
		return nil, fmt.Errorf("read contract: %w", err)
	}

	var contract Contract
	if err := yaml.Unmarshal(contractData, &contract); err != nil {
		return nil, fmt.Errorf("parse contract: %w", err)
	}

	// Parse implementation file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, implPath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse implementation: %w", err)
	}

	// Extract struct fields from implementation
	implFields := make(map[string]string)
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != contract.TypeName {
				continue
			}
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				for _, field := range structType.Fields.List {
					fieldType := formatFieldType(field.Type)
					for _, name := range field.Names {
						implFields[name.Name] = fieldType
					}
				}
			}
		}
	}

	// Build contract field map
	contractFields := make(map[string]string)
	for _, f := range contract.Fields {
		contractFields[f.Name] = f.Type
	}

	var violations []Violation

	// Check for missing required fields
	for fieldName, fieldType := range contractFields {
		if implType, exists := implFields[fieldName]; !exists {
			violations = append(violations, Violation{
				Type:     "missing_field",
				Field:    fieldName,
				Expected: fieldType,
				Actual:   "missing",
				Severity: "error",
				Message:  fmt.Sprintf("Missing required field: %s (%s)", fieldName, fieldType),
			})
		} else if implType != fieldType {
			violations = append(violations, Violation{
				Type:     "type_mismatch",
				Field:    fieldName,
				Expected: fieldType,
				Actual:   implType,
				Severity: "warning",
				Message:  fmt.Sprintf("Type mismatch for %s: expected %s, got %s", fieldName, fieldType, implType),
			})
		}
	}

	// Check for extra fields (warning only in P1)
	for fieldName, fieldType := range implFields {
		if _, exists := contractFields[fieldName]; !exists {
			violations = append(violations, Violation{
				Type:     "extra_field",
				Field:    fieldName,
				Expected: "not in contract",
				Actual:   fieldType,
				Severity: "warning",
				Message:  fmt.Sprintf("Extra field not in contract: %s (%s)", fieldName, fieldType),
			})
		}
	}

	return violations, nil
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
		if _, err := os.Stat(implPath); os.IsNotExist(err) {
			// Try lowercase
			implPath = filepath.Join(implDir, toLowerFirst(typeName)+".go")
			if _, err := os.Stat(implPath); os.IsNotExist(err) {
				continue
			}
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
