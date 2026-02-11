package metrics

import (
	"path/filepath"
	"testing"
)

func TestTaxonomy_ClassifyFailure_WrongLogic(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")
	taxonomy := NewTaxonomy(taxonomyPath)

	// Act - verification output with "assert" error (wrong logic)
	verificationOutput := `expected 5 but got 3
    Test failed: assertion violated`
	classification := taxonomy.ClassifyFromOutput("evt-1", "00-001-01", "claude-sonnet-4", "go", verificationOutput)

	// Assert
	if classification.FailureType != "wrong_logic" {
		t.Errorf("Expected failure type wrong_logic, got %s", classification.FailureType)
	}
	if classification.Severity != "MEDIUM" {
		t.Errorf("Expected severity MEDIUM, got %s", classification.Severity)
	}
}

func TestTaxonomy_ClassifyFailure_MissingEdgeCase(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")
	taxonomy := NewTaxonomy(taxonomyPath)

	// Act - verification output with nil pointer, out of bounds (edge case)
	verificationOutput := `panic: runtime error
    index out of range
    nil pointer dereference`
	classification := taxonomy.ClassifyFromOutput("evt-2", "00-001-02", "claude-opus-4", "python", verificationOutput)

	// Assert
	if classification.FailureType != "missing_edge_case" {
		t.Errorf("Expected failure type missing_edge_case, got %s", classification.FailureType)
	}
}

func TestTaxonomy_ClassifyFailure_HallucinatedAPI(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")
	taxonomy := NewTaxonomy(taxonomyPath)

	// Act - verification output with "undefined function" or "no such method"
	verificationOutput := `undefined: function.GetName not found
    compilation error: api.DoSomething() undefined`
	classification := taxonomy.ClassifyFromOutput("evt-3", "00-001-03", "claude-sonnet-4", "javascript", verificationOutput)

	// Assert
	if classification.FailureType != "hallucinated_api" {
		t.Errorf("Expected failure type hallucinated_api, got %s", classification.FailureType)
	}
}

func TestTaxonomy_ClassifyFailure_TypeError(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")
	taxonomy := NewTaxonomy(taxonomyPath)

	// Act - verification output with type error, type mismatch
	verificationOutput := `cannot use str (type string) as type int
    type error: int expected
    static type checking failed`
	classification := taxonomy.ClassifyFromOutput("evt-4", "00-001-04", "claude-opus-4", "go", verificationOutput)

	// Assert
	if classification.FailureType != "type_error" {
		t.Errorf("Expected failure type type_error, got %s", classification.FailureType)
	}
}

func TestTaxonomy_ClassifyFailure_CompilationError(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")
	taxonomy := NewTaxonomy(taxonomyPath)

	// Act - verification output with syntax error
	verificationOutput := `syntax error near unexpected token
    compilation failed: invalid syntax
    parse error`
	classification := taxonomy.ClassifyFromOutput("evt-5", "00-001-05", "claude-sonnet-4", "python", verificationOutput)

	// Assert
	if classification.FailureType != "compilation_error" {
		t.Errorf("Expected failure type compilation_error, got %s", classification.FailureType)
	}
}

func TestTaxonomy_ClassifyFailure_ImportError(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")
	taxonomy := NewTaxonomy(taxonomyPath)

	// Act - verification output with import error, module not found
	verificationOutput := `module 'requests' not found
    import error: no such package
    cannot resolve import`
	classification := taxonomy.ClassifyFromOutput("evt-6", "00-001-06", "claude-opus-4", "python", verificationOutput)

	// Assert
	if classification.FailureType != "import_error" {
		t.Errorf("Expected failure type import_error, got %s", classification.FailureType)
	}
}

func TestTaxonomy_ClassifyFailure_UnknownError_ReturnsUnknown(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")
	taxonomy := NewTaxonomy(taxonomyPath)

	// Act - verification output with no recognized pattern
	verificationOutput := `some unknown error pattern`
	classification := taxonomy.ClassifyFromOutput("evt-7", "00-001-07", "claude-sonnet-4", "go", verificationOutput)

	// Assert
	if classification.FailureType != "unknown" {
		t.Errorf("Expected failure type unknown, got %s", classification.FailureType)
	}
}

func TestTaxonomy_ManualOverride_SetsType(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")
	taxonomy := NewTaxonomy(taxonomyPath)

	// Act - manually override classification
	taxonomy.SetClassification("evt-1", "wrong_logic", "Manual correction by user")

	// Assert
	classification, exists := taxonomy.GetClassification("evt-1")
	if !exists {
		t.Fatal("Expected classification to exist after manual override")
	}
	if classification.FailureType != "wrong_logic" {
		t.Errorf("Expected failure type wrong_logic, got %s", classification.FailureType)
	}
	if classification.Notes != "Manual correction by user" {
		t.Errorf("Expected notes 'Manual correction by user', got %s", classification.Notes)
	}
}

func TestTaxonomy_SaveAndLoad_PreservesData(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")
	taxonomy := NewTaxonomy(taxonomyPath)

	// Act - add classification and save
	taxonomy.SetClassification("evt-1", "wrong_logic", "Test note")
	if err := taxonomy.Save(); err != nil {
		t.Fatalf("Failed to save taxonomy: %v", err)
	}

	// Load new taxonomy instance
	taxonomy2 := NewTaxonomy(taxonomyPath)
	if err := taxonomy2.Load(); err != nil {
		t.Fatalf("Failed to load taxonomy: %v", err)
	}

	// Assert
	classification, exists := taxonomy2.GetClassification("evt-1")
	if !exists {
		t.Fatal("Expected classification to exist after save/load")
	}
	if classification.FailureType != "wrong_logic" {
		t.Errorf("Expected failure type wrong_logic, got %s", classification.FailureType)
	}
}

func TestTaxonomy_GetByModel_ReturnsModelSpecificFailures(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")
	taxonomy := NewTaxonomy(taxonomyPath)

	// Classify directly with full info
	_ = taxonomy.ClassifyFromOutput("evt-1", "00-001-01", "claude-sonnet-4", "go", "assertion failed")
	_ = taxonomy.ClassifyFromOutput("evt-2", "00-001-02", "claude-opus-4", "python", "type error")
	_ = taxonomy.ClassifyFromOutput("evt-3", "00-001-03", "claude-sonnet-4", "go", "expected 5 but got 3")

	// Act - get classifications by model
	sonnetFailures := taxonomy.GetByModel("claude-sonnet-4")

	// Assert
	if len(sonnetFailures) != 2 {
		t.Errorf("Expected 2 failures for claude-sonnet-4, got %d", len(sonnetFailures))
	}
}

func TestTaxonomy_GetByType_ReturnsTypeSpecificFailures(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")
	taxonomy := NewTaxonomy(taxonomyPath)

	taxonomy.SetClassification("evt-1", "wrong_logic", "Note 1")
	taxonomy.SetClassification("evt-2", "type_error", "Note 2")
	taxonomy.SetClassification("evt-3", "wrong_logic", "Note 3")

	// Act - get classifications by type
	logicFailures := taxonomy.GetByType("wrong_logic")

	// Assert
	if len(logicFailures) != 2 {
		t.Errorf("Expected 2 wrong_logic failures, got %d", len(logicFailures))
	}
}

func TestTaxonomy_GetStats_ReturnsSummary(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")
	taxonomy := NewTaxonomy(taxonomyPath)

	// Add classifications via ClassifyFromOutput
	_ = taxonomy.ClassifyFromOutput("evt-1", "00-001-01", "claude-sonnet-4", "go", "assertion failed")
	_ = taxonomy.ClassifyFromOutput("evt-2", "00-001-02", "claude-opus-4", "python", "type error")
	_ = taxonomy.ClassifyFromOutput("evt-3", "00-001-03", "claude-sonnet-4", "go", "expected but got")

	// Act
	stats := taxonomy.GetStats()

	// Assert
	if stats.TotalClassifications != 3 {
		t.Errorf("Expected 3 total classifications, got %d", stats.TotalClassifications)
	}
	if stats.TotalByModel["claude-sonnet-4"] != 2 {
		t.Errorf("Expected 2 classifications for claude-sonnet-4, got %d", stats.TotalByModel["claude-sonnet-4"])
	}
	if stats.TotalByType["wrong_logic"] != 2 {
		t.Errorf("Expected 2 wrong_logic classifications, got %d", stats.TotalByType["wrong_logic"])
	}
}
