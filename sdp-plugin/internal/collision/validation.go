package collision

import (
	"fmt"
	"os"

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
