package collision

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGenerateContractFromBoundaries tests contract generation from shared boundaries.
func TestGenerateContractFromBoundaries(t *testing.T) {
	// Arrange
	boundaries := []SharedBoundary{
		{
			FileName:  "internal/model/user.go",
			TypeName:  "User",
			Fields:    []FieldInfo{{Name: "Email", Type: "string"}, {Name: "Name", Type: "string"}},
			Features:  []string{"F054", "F055"},
		},
	}
	outputDir := t.TempDir()

	// Act
	contracts, err := GenerateContracts(boundaries, outputDir)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(contracts) != 1 {
		t.Fatalf("Expected 1 contract, got %d", len(contracts))
	}
	// Verify file was created
	contractPath := filepath.Join(outputDir, "User.yaml")
	if _, err := os.Stat(contractPath); os.IsNotExist(err) {
		t.Errorf("Expected contract file at %s, but it doesn't exist", contractPath)
	}
}

// TestContractYAMLFormat tests that generated contract has correct YAML format.
func TestContractYAMLFormat(t *testing.T) {
	// Arrange
	boundaries := []SharedBoundary{
		{
			FileName:  "internal/model/user.go",
			TypeName:  "User",
			Fields:    []FieldInfo{{Name: "Email", Type: "string"}, {Name: "Name", Type: "string"}},
			Features:  []string{"F054", "F055"},
		},
	}
	outputDir := t.TempDir()

	// Act
	contracts, err := GenerateContracts(boundaries, outputDir)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(contracts) == 0 {
		t.Fatal("Expected at least 1 contract")
	}

	// Read the contract file
	contractPath := filepath.Join(outputDir, "User.yaml")
	data, err := os.ReadFile(contractPath)
	if err != nil {
		t.Fatalf("Failed to read contract file: %v", err)
	}

	content := string(data)
	// Verify YAML contains expected fields
	if !contains(content, "typeName") || !contains(content, "User") {
		t.Error("Expected contract to contain typeName: User")
	}
	if !contains(content, "requiredBy") {
		t.Error("Expected contract to contain requiredBy field")
	}
	if !contains(content, "status") {
		t.Error("Expected contract to contain status field")
	}
}

// TestGenerateContract_FieldAggregation tests field aggregation across features.
func TestGenerateContract_FieldAggregation(t *testing.T) {
	// Arrange - two features with overlapping field needs
	boundaries := []SharedBoundary{
		{
			FileName:  "internal/model/user.go",
			TypeName:  "User",
			Fields:    []FieldInfo{{Name: "Email", Type: "string"}, {Name: "Name", Type: "string"}},
			Features:  []string{"F054", "F055"},
		},
	}
	outputDir := t.TempDir()

	// Act
	contracts, err := GenerateContracts(boundaries, outputDir)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(contracts) != 1 {
		t.Fatalf("Expected 1 contract, got %d", len(contracts))
	}

	c := contracts[0]
	if c.TypeName != "User" {
		t.Errorf("Expected type name 'User', got '%s'", c.TypeName)
	}
	if len(c.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(c.Fields))
	}
	if len(c.RequiredBy) != 2 {
		t.Errorf("Expected 2 features in requiredBy, got %d", len(c.RequiredBy))
	}
}

// TestGenerateContract_MultipleBoundaries tests handling multiple boundaries.
func TestGenerateContract_MultipleBoundaries(t *testing.T) {
	// Arrange
	boundaries := []SharedBoundary{
		{
			FileName:  "internal/model/user.go",
			TypeName:  "User",
			Fields:    []FieldInfo{{Name: "Email", Type: "string"}},
			Features:  []string{"F054", "F055"},
		},
		{
			FileName:  "internal/model/user.go",
			TypeName:  "Profile",
			Fields:    []FieldInfo{{Name: "Bio", Type: "string"}},
			Features:  []string{"F054"},
		},
	}
	outputDir := t.TempDir()

	// Act
	contracts, err := GenerateContracts(boundaries, outputDir)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(contracts) != 2 {
		t.Fatalf("Expected 2 contracts, got %d", len(contracts))
	}
}
