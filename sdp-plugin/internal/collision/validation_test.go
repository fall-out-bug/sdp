package collision

import (
	"os"
	"path/filepath"
	"testing"
)

// TestExtractImplFields_BasicFunctionality tests extractImplFields basic behavior.
// RED test: Ensure refactoring maintains correctness.
func TestExtractImplFields_BasicFunctionality(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	implPath := filepath.Join(tmpDir, "user.go")
	implContent := `package impl

type User struct {
	Email string
	Name  string
	Age   int
}`
	if err := os.WriteFile(implPath, []byte(implContent), 0644); err != nil {
		t.Fatalf("Failed to create impl: %v", err)
	}

	// Act
	fields, err := extractImplFields(implPath, "User")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(fields))
	}
	if fields["Email"] != "string" {
		t.Errorf("Expected Email type 'string', got '%s'", fields["Email"])
	}
	if fields["Name"] != "string" {
		t.Errorf("Expected Name type 'string', got '%s'", fields["Name"])
	}
	if fields["Age"] != "int" {
		t.Errorf("Expected Age type 'int', got '%s'", fields["Age"])
	}
}

// TestExtractImplFields_PointerTypes tests pointer type handling.
func TestExtractImplFields_PointerTypes(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	implPath := filepath.Join(tmpDir, "user.go")
	implContent := `package impl

type User struct {
	Email *string
	Name  string
}`
	if err := os.WriteFile(implPath, []byte(implContent), 0644); err != nil {
		t.Fatalf("Failed to create impl: %v", err)
	}

	// Act
	fields, err := extractImplFields(implPath, "User")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if fields["Email"] != "*string" {
		t.Errorf("Expected Email type '*string', got '%s'", fields["Email"])
	}
}

// TestExtractImplFields_SliceTypes tests slice type handling.
func TestExtractImplFields_SliceTypes(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	implPath := filepath.Join(tmpDir, "user.go")
	implContent := `package impl

type User struct {
	Tags []string
	Name string
}`
	if err := os.WriteFile(implPath, []byte(implContent), 0644); err != nil {
		t.Fatalf("Failed to create impl: %v", err)
	}

	// Act
	fields, err := extractImplFields(implPath, "User")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if fields["Tags"] != "[]string" {
		t.Errorf("Expected Tags type '[]string', got '%s'", fields["Tags"])
	}
}

// TestExtractImplFields_MapTypes tests map type handling.
func TestExtractImplFields_MapTypes(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	implPath := filepath.Join(tmpDir, "user.go")
	implContent := `package impl

type User struct {
	Meta map[string]string
	Name string
}`
	if err := os.WriteFile(implPath, []byte(implContent), 0644); err != nil {
		t.Fatalf("Failed to create impl: %v", err)
	}

	// Act
	fields, err := extractImplFields(implPath, "User")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if fields["Meta"] != "map[string]string" {
		t.Errorf("Expected Meta type 'map[string]string', got '%s'", fields["Meta"])
	}
}

// TestExtractImplFields_TypeNotFound tests behavior when type doesn't exist.
func TestExtractImplFields_TypeNotFound(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	implPath := filepath.Join(tmpDir, "user.go")
	implContent := `package impl

type User struct {
	Email string
}`
	if err := os.WriteFile(implPath, []byte(implContent), 0644); err != nil {
		t.Fatalf("Failed to create impl: %v", err)
	}

	// Act
	fields, err := extractImplFields(implPath, "NonExistent")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(fields) != 0 {
		t.Errorf("Expected 0 fields for non-existent type, got %d", len(fields))
	}
}

// TestExtractImplFields_MultipleTypes tests handling multiple types in same file.
func TestExtractImplFields_MultipleTypes(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	implPath := filepath.Join(tmpDir, "models.go")
	implContent := `package impl

type User struct {
	Email string
}

type Admin struct {
	Role string
}`
	if err := os.WriteFile(implPath, []byte(implContent), 0644); err != nil {
		t.Fatalf("Failed to create impl: %v", err)
	}

	// Act
	userFields, err := extractImplFields(implPath, "User")
	adminFields, err2 := extractImplFields(implPath, "Admin")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error for User, got %v", err)
	}
	if err2 != nil {
		t.Fatalf("Expected no error for Admin, got %v", err2)
	}
	if len(userFields) != 1 {
		t.Errorf("Expected 1 field for User, got %d", len(userFields))
	}
	if len(adminFields) != 1 {
		t.Errorf("Expected 1 field for Admin, got %d", len(adminFields))
	}
}

// TestExtractImplFields_EmptyStruct tests empty struct handling.
func TestExtractImplFields_EmptyStruct(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	implPath := filepath.Join(tmpDir, "user.go")
	implContent := `package impl

type User struct {
}`
	if err := os.WriteFile(implPath, []byte(implContent), 0644); err != nil {
		t.Fatalf("Failed to create impl: %v", err)
	}

	// Act
	fields, err := extractImplFields(implPath, "User")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(fields) != 0 {
		t.Errorf("Expected 0 fields for empty struct, got %d", len(fields))
	}
}

// TestExtractImplFields_NestedTypes tests nested type handling.
func TestExtractImplFields_NestedTypes(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	implPath := filepath.Join(tmpDir, "user.go")
	implContent := `package impl

type User struct {
	Email string
	Profile Profile
}

type Profile struct {
	Bio string
}`
	if err := os.WriteFile(implPath, []byte(implContent), 0644); err != nil {
		t.Fatalf("Failed to create impl: %v", err)
	}

	// Act
	fields, err := extractImplFields(implPath, "User")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if fields["Email"] != "string" {
		t.Errorf("Expected Email type 'string', got '%s'", fields["Email"])
	}
	if fields["Profile"] != "Profile" {
		t.Errorf("Expected Profile type 'Profile', got '%s'", fields["Profile"])
	}
}

// TestExtractImplFields_NotAStruct tests non-struct type handling.
func TestExtractImplFields_NotAStruct(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	implPath := filepath.Join(tmpDir, "user.go")
	implContent := `package impl

type User string`
	if err := os.WriteFile(implPath, []byte(implContent), 0644); err != nil {
		t.Fatalf("Failed to create impl: %v", err)
	}

	// Act
	fields, err := extractImplFields(implPath, "User")

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(fields) != 0 {
		t.Errorf("Expected 0 fields for non-struct type, got %d", len(fields))
	}
}
