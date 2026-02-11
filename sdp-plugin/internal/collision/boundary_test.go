package collision

import (
	"testing"
)

// TestDetectBoundaries_NoOverlaps_ReturnsEmpty tests the case where no shared boundaries exist.
func TestDetectBoundaries_NoOverlaps_ReturnsEmpty(t *testing.T) {
	// Arrange
	workstreams := []FeatureScope{
		{
			FeatureID: "F054",
			ScopeFiles: []string{
				"internal/auth/user.go",
				"internal/auth/session.go",
			},
		},
		{
			FeatureID: "F055",
			ScopeFiles: []string{
				"internal/billing/payment.go",
				"internal/billing/invoice.go",
			},
		},
	}

	// Act
	boundaries := DetectBoundaries(workstreams)

	// Assert
	if len(boundaries) != 0 {
		t.Errorf("Expected 0 boundaries, got %d", len(boundaries))
	}
}

// TestDetectBoundaries_SharedType_ReturnsBoundary tests detection of shared types.
func TestDetectBoundaries_SharedType_ReturnsBoundary(t *testing.T) {
	// Arrange - use test data file for actual parsing
	workstreams := []FeatureScope{
		{
			FeatureID: "F054",
			ScopeFiles: []string{
				"testdata/user_model.go", // Contains User struct
			},
		},
		{
			FeatureID: "F055",
			ScopeFiles: []string{
				"testdata/user_model.go", // Same file - shared type
			},
		},
	}

	// Act
	boundaries := DetectBoundaries(workstreams)

	// Assert - should detect User and Profile types
	if len(boundaries) < 1 {
		t.Fatalf("Expected at least 1 boundary, got %d", len(boundaries))
	}
	// Check first boundary has both features
	if len(boundaries[0].Features) != 2 {
		t.Errorf("Expected 2 features in first boundary, got %d", len(boundaries[0].Features))
	}
}

// TestDetectBoundaries_ParsesGoStructs tests parsing Go structs from files.
func TestDetectBoundaries_ParsesGoStructs(t *testing.T) {
	// Arrange - This test will need test data files
	workstreams := []FeatureScope{
		{
			FeatureID: "F054",
			ScopeFiles: []string{
				"testdata/user_model.go", // Test data file with User struct
			},
		},
		{
			FeatureID: "F055",
			ScopeFiles: []string{
				"testdata/user_model.go",
			},
		},
	}

	// Act
	boundaries := DetectBoundaries(workstreams)

	// Assert - should detect User type
	if len(boundaries) < 1 {
		t.Fatalf("Expected at least 1 boundary, got %d", len(boundaries))
	}
	// Find User boundary
	var userBoundary *SharedBoundary
	for i := range boundaries {
		if boundaries[i].TypeName == "User" {
			userBoundary = &boundaries[i]
			break
		}
	}
	if userBoundary == nil {
		t.Fatal("Expected to find User boundary")
	}
	if len(userBoundary.Fields) == 0 {
		t.Error("Expected fields to be parsed, got empty slice")
	}
}

// TestDetectBoundaries_FieldOverlap tests that field-level overlap is detected.
func TestDetectBoundaries_FieldOverlap(t *testing.T) {
	// Arrange
	workstreams := []FeatureScope{
		{
			FeatureID: "F054",
			ScopeFiles: []string{
				"testdata/user_model.go", // User with Email, Name fields
			},
		},
		{
			FeatureID: "F055",
			ScopeFiles: []string{
				"testdata/user_model.go", // Same User struct
			},
		},
	}

	// Act
	boundaries := DetectBoundaries(workstreams)

	// Assert - find User boundary and check for Email field
	var userBoundary *SharedBoundary
	for i := range boundaries {
		if boundaries[i].TypeName == "User" {
			userBoundary = &boundaries[i]
			break
		}
	}
	if userBoundary == nil {
		t.Fatal("Expected to find User boundary")
	}
	// Check that Email field is detected
	hasEmail := false
	for _, f := range userBoundary.Fields {
		if f.Name == "Email" {
			hasEmail = true
			break
		}
	}
	if !hasEmail {
		t.Error("Expected 'Email' field to be detected in boundary")
	}
}

// TestBoundaryToJSON tests JSON output generation.
func TestBoundaryToJSON(t *testing.T) {
	// Arrange
	boundary := SharedBoundary{
		FileName:       "internal/model/user.go",
		TypeName:       "User",
		Fields:         []FieldInfo{{Name: "Email", Type: "string"}},
		Features:       []string{"F054", "F055"},
		Recommendation: "Create shared interface contract",
	}

	// Act
	jsonOutput, err := BoundaryToJSON(boundary)

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if jsonOutput == "" {
		t.Error("Expected JSON output, got empty string")
	}
	// Verify it contains expected fields (case sensitive)
	if !contains(jsonOutput, "typeName") && !contains(jsonOutput, "TypeName") {
		t.Error("Expected JSON to contain typeName")
	}
	if !contains(jsonOutput, "User") {
		t.Error("Expected JSON to contain User")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr) >= 0
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
