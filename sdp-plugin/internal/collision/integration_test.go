package collision

import (
	"os"
	"path/filepath"
	"testing"
)

// TestContractValidate_ImplementationVsContract tests contract validation.
func TestContractValidate_ImplementationVsContract(t *testing.T) {
	// Arrange - create a temporary contract directory
	tmpDir := t.TempDir()
	contractPath := filepath.Join(tmpDir, "User.yaml")

	// Create a sample contract
	contractContent := `typeName: User
fields:
  - name: Email
    type: string
  - name: Name
    type: string
requiredBy:
  - F054
  - F055
status: locked
`
	if err := os.WriteFile(contractPath, []byte(contractContent), 0644); err != nil {
		t.Fatalf("Failed to create contract: %v", err)
	}

	// Create a sample implementation file that matches
	implPath := filepath.Join(tmpDir, "user_impl.go")
	implContent := `package impl

type User struct {
	Email string
	Name  string
	Role  string // Extra field not in contract
}
`
	if err := os.WriteFile(implPath, []byte(implContent), 0644); err != nil {
		t.Fatalf("Failed to create impl: %v", err)
	}

	// Act
	violations, err := ValidateContractAgainstImpl(contractPath, implPath)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	// Should have warning about extra field (but not error in P1)
	if len(violations) == 0 {
		t.Error("Expected at least one violation (extra field)")
	}
}

// TestContractValidate_MissingContractField tests missing required field.
func TestContractValidate_MissingContractField(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	contractPath := filepath.Join(tmpDir, "User.yaml")

	// Contract requires Email and Name
	contractContent := `typeName: User
fields:
  - name: Email
    type: string
  - name: Name
    type: string
requiredBy:
  - F054
status: locked
`
	if err := os.WriteFile(contractPath, []byte(contractContent), 0644); err != nil {
		t.Fatalf("Failed to create contract: %v", err)
	}

	// Implementation only has Email (missing Name)
	implPath := filepath.Join(tmpDir, "user_impl.go")
	implContent := `package impl

type User struct {
	Email string
}
`
	if err := os.WriteFile(implPath, []byte(implContent), 0644); err != nil {
		t.Fatalf("Failed to create impl: %v", err)
	}

	// Act
	violations, err := ValidateContractAgainstImpl(contractPath, implPath)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	// Should have error about missing Name field
	hasMissingField := false
	for _, v := range violations {
		if v.Type == "missing_field" && v.Field == "Name" {
			hasMissingField = true
			break
		}
	}
	if !hasMissingField {
		t.Error("Expected violation for missing Name field")
	}
}

// TestContractValidate_TypeMismatch tests type mismatch detection.
func TestContractValidate_TypeMismatch(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	contractPath := filepath.Join(tmpDir, "User.yaml")

	// Contract specifies Email as string
	contractContent := `typeName: User
fields:
  - name: Email
    type: string
requiredBy:
  - F054
status: locked
`
	if err := os.WriteFile(contractPath, []byte(contractContent), 0644); err != nil {
		t.Fatalf("Failed to create contract: %v", err)
	}

	// Implementation has Email as int
	implPath := filepath.Join(tmpDir, "user_impl.go")
	implContent := `package impl

type User struct {
	Email int
}
`
	if err := os.WriteFile(implPath, []byte(implContent), 0644); err != nil {
		t.Fatalf("Failed to create impl: %v", err)
	}

	// Act
	violations, err := ValidateContractAgainstImpl(contractPath, implPath)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	// Should have warning about type mismatch
	hasTypeMismatch := false
	for _, v := range violations {
		if v.Type == "type_mismatch" && v.Field == "Email" {
			hasTypeMismatch = true
			break
		}
	}
	if !hasTypeMismatch {
		t.Error("Expected violation for Email type mismatch")
	}
}

// TestE2E_ContractWorkflow tests end-to-end contract workflow.
func TestE2E_ContractWorkflow(t *testing.T) {
	// Arrange - simulate two features with shared boundary
	tmpDir := t.TempDir()

	boundaries := []SharedBoundary{
		{
			FileName:  "internal/model/user.go",
			TypeName:  "User",
			Fields:    []FieldInfo{{Name: "Email", Type: "string"}, {Name: "Name", Type: "string"}},
			Features:  []string{"F054", "F055"},
		},
	}

	// Act - generate contracts
	contracts, err := GenerateContracts(boundaries, tmpDir)

	// Assert
	if err != nil {
		t.Fatalf("Failed to generate contracts: %v", err)
	}
	if len(contracts) != 1 {
		t.Fatalf("Expected 1 contract, got %d", len(contracts))
	}

	// Verify contract file exists
	contractPath := filepath.Join(tmpDir, "User.yaml")
	if _, err := os.Stat(contractPath); os.IsNotExist(err) {
		t.Errorf("Contract file not created at %s", contractPath)
	}

	// Verify contract can be loaded
	data, err := os.ReadFile(contractPath)
	if err != nil {
		t.Fatalf("Failed to read contract: %v", err)
	}
	if len(data) == 0 {
		t.Error("Contract file is empty")
	}
}

// TestValidateContractsInDir tests directory-level validation.
func TestValidateContractsInDir(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	contractsDir := filepath.Join(tmpDir, "contracts")
	implDir := filepath.Join(tmpDir, "impl")

	if err := os.MkdirAll(contractsDir, 0755); err != nil {
		t.Fatalf("Failed to create contracts dir: %v", err)
	}
	if err := os.MkdirAll(implDir, 0755); err != nil {
		t.Fatalf("Failed to create impl dir: %v", err)
	}

	// Create contract
	contractContent := `typeName: User
fields:
  - name: Email
    type: string
requiredBy:
  - F054
status: locked
`
	if err := os.WriteFile(filepath.Join(contractsDir, "User.yaml"), []byte(contractContent), 0644); err != nil {
		t.Fatalf("Failed to create contract: %v", err)
	}

	// Create implementation
	implContent := `package impl

type User struct {
	Email string
}
`
	if err := os.WriteFile(filepath.Join(implDir, "User.go"), []byte(implContent), 0644); err != nil {
		t.Fatalf("Failed to create impl: %v", err)
	}

	// Act
	violations, err := ValidateContractsInDir(contractsDir, implDir)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	// Should have no violations (Email matches)
	if len(violations) != 0 {
		t.Errorf("Expected 0 violations, got %d", len(violations))
	}
}

// TestValidateContractsInDir_NoMatchingImpl tests when no impl exists.
func TestValidateContractsInDir_NoMatchingImpl(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	contractsDir := filepath.Join(tmpDir, "contracts")
	implDir := filepath.Join(tmpDir, "impl")

	if err := os.MkdirAll(contractsDir, 0755); err != nil {
		t.Fatalf("Failed to create contracts dir: %v", err)
	}
	if err := os.MkdirAll(implDir, 0755); err != nil {
		t.Fatalf("Failed to create impl dir: %v", err)
	}

	// Create contract without matching impl
	contractContent := `typeName: User
fields:
  - name: Email
    type: string
requiredBy:
  - F054
status: locked
`
	if err := os.WriteFile(filepath.Join(contractsDir, "User.yaml"), []byte(contractContent), 0644); err != nil {
		t.Fatalf("Failed to create contract: %v", err)
	}

	// Act - should not error, just skip
	violations, err := ValidateContractsInDir(contractsDir, implDir)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(violations) != 0 {
		t.Errorf("Expected 0 violations (no impl to validate), got %d", len(violations))
	}
}

// TestToLowerFirst tests the toLowerFirst helper.
func TestToLowerFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"User", "user"},
		{"Name", "name"},
		{"Email", "email"},
		{"A", "a"},
		{"", ""},
		{"a", "a"}, // Already lowercase
	}

	for _, tt := range tests {
		result := toLowerFirst(tt.input)
		if result != tt.expected {
			t.Errorf("toLowerFirst(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
