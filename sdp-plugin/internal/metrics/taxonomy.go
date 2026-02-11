package metrics

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

// FailureClassification represents a classified failure (AC1).
type FailureClassification struct {
	EventID     string `json:"event_id"`
	WSID        string `json:"ws_id"`
	ModelID     string `json:"model_id"`
	Language    string `json:"language"`
	FailureType string `json:"failure_type"`
	Severity    string `json:"severity"`
	Notes       string `json:"notes,omitempty"`
}

// TaxonomyStats provides summary statistics (AC1).
type TaxonomyStats struct {
	TotalClassifications int                `json:"total_classifications"`
	TotalByModel       map[string]int      `json:"total_by_model"`
	TotalByType        map[string]int      `json:"total_by_type"`
	TotalByLanguage    map[string]int      `json:"total_by_language"`
	TotalBySeverity    map[string]int      `json:"total_by_severity"`
}

// FailureType enum (AC2).
const (
	FailureWrongLogic        = "wrong_logic"
	FailureMissingEdgeCase  = "missing_edge_case"
	FailureHallucinatedAPI  = "hallucinated_api"
	FailureTypeError        = "type_error"
	FailureTestPassingWrong = "test_passing_but_wrong"
	FailureCompilationError  = "compilation_error"
	FailureImportError      = "import_error"
	FailureUnknown         = "unknown"
)

// Severity levels.
const (
	SeverityLow      = "LOW"
	SeverityMedium   = "MEDIUM"
	SeverityHigh     = "HIGH"
	SeverityCritical = "CRITICAL"
)

// Taxonomy manages failure classifications (AC3-AC6).
type Taxonomy struct {
	path          string
	classifications map[string]*FailureClassification
	mu            sync.RWMutex
}

// NewTaxonomy creates a taxonomy manager for the given path.
func NewTaxonomy(path string) *Taxonomy {
	return &Taxonomy{
		path:          path,
		classifications: make(map[string]*FailureClassification),
	}
}

// Load loads existing taxonomy from file (AC6).
func (t *Taxonomy) Load() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	data, err := os.ReadFile(t.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No existing taxonomy
		}
		return fmt.Errorf("read taxonomy: %w", err)
	}

	if len(data) == 0 {
		return nil // Empty file
	}

	var list []FailureClassification
	if err := json.Unmarshal(data, &list); err != nil {
		return fmt.Errorf("parse taxonomy: %w", err)
	}

	// Rebuild map
	t.classifications = make(map[string]*FailureClassification)
	for i := range list {
		fc := &list[i]
		t.classifications[fc.EventID] = fc
	}

	return nil
}

// Save writes taxonomy to file (AC6).
func (t *Taxonomy) Save() error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// Convert map to slice for JSON
	list := make([]FailureClassification, 0, len(t.classifications))
	for _, fc := range t.classifications {
		list = append(list, *fc)
	}

	// Ensure directory exists
	dir := filepath.Dir(t.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create taxonomy dir: %w", err)
	}

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal taxonomy: %w", err)
	}

	return os.WriteFile(t.path, data, 0644)
}

// ClassifyFromOutput auto-classifies failure from verification output (AC3, AC4).
func (t *Taxonomy) ClassifyFromOutput(eventID, wsID, modelID, language, output string) FailureClassification {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Determine failure type based on patterns
	failureType := t.classifyByPattern(output)
	severity := t.severityForType(failureType)

	fc := &FailureClassification{
		EventID:     eventID,
		WSID:        wsID,
		ModelID:     modelID,
		Language:    language,
		FailureType: failureType,
		Severity:    severity,
	}

	t.classifications[eventID] = fc
	return *fc
}

// classifyByPattern determines failure type from output patterns (AC4).
func (t *Taxonomy) classifyByPattern(output string) string {
	outputLower := strings.ToLower(output)

	// Check patterns in priority order
	if t.matchesAny(outputLower, []string{
		"assertion failed",
		"assertion error",
		"expected.*but got",
		"assertion violated",
		"value.*does not match",
	}) {
		return FailureWrongLogic
	}

	// Edge case patterns
	if t.matchesAny(outputLower, []string{
		"nil pointer",
		"null pointer",
		"index out of",
		"out of range",
		"out of bounds",
		"panic:",
		"runtime error",
		"segmentation fault",
		"access violation",
	}) {
		return FailureMissingEdgeCase
	}

	// Hallucinated API patterns
	if t.matchesAny(outputLower, []string{
		"undefined.*function",
		"undefined.*method",
		"no such.*function",
		"not a function",
		"no such method",
		"api.*not found",
		"undefined symbol",
		"package.*is not in",
	}) {
		return FailureHallucinatedAPI
	}

	// Type error patterns
	if t.matchesAny(outputLower, []string{
		"type error",
		"cannot use.*as type",
		"type mismatch",
		"cannot convert",
		"incompatible type",
		"static type checking",
	}) {
		return FailureTypeError
	}

	// Import error patterns
	if t.matchesAny(outputLower, []string{
		"import error",
		"module.*not found",
		"cannot resolve import",
		"no such package",
		"no such module",
		"undefined: package",
	}) {
		return FailureImportError
	}

	// Compilation error patterns
	if t.matchesAny(outputLower, []string{
		"syntax error",
		"parse error",
		"invalid syntax",
		"unexpected token",
		"compilation failed",
		"build error",
		"compiler error",
	}) {
		return FailureCompilationError
	}

	// Default to unknown
	return FailureUnknown
}

// matchesAny checks if output matches any of the patterns.
func (t *Taxonomy) matchesAny(output string, patterns []string) bool {
	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, output); matched {
			return true
		}
	}
	return false
}

// severityForType returns severity level for failure type.
func (t *Taxonomy) severityForType(failureType string) string {
	switch failureType {
	case FailureMissingEdgeCase:
		return SeverityHigh
	case FailureTypeError, FailureCompilationError, FailureImportError:
		return SeverityMedium
	case FailureWrongLogic, FailureHallucinatedAPI:
		return SeverityMedium
	case FailureTestPassingWrong:
		return SeverityCritical // Acceptance test catches what unit tests missed
	default:
		return SeverityLow
	}
}

// SetClassification manually sets classification (AC5).
func (t *Taxonomy) SetClassification(eventID, failureType, notes string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if fc, exists := t.classifications[eventID]; exists {
		fc.FailureType = failureType
		fc.Notes = notes
	} else {
		t.classifications[eventID] = &FailureClassification{
			EventID:     eventID,
			FailureType: failureType,
			Severity:    t.severityForType(failureType),
			Notes:       notes,
		}
	}
}

// GetClassification retrieves classification by event ID.
func (t *Taxonomy) GetClassification(eventID string) (FailureClassification, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	fc, exists := t.classifications[eventID]
	if !exists {
		return FailureClassification{}, false
	}
	return *fc, true
}

// GetByModel returns all classifications for a model.
func (t *Taxonomy) GetByModel(modelID string) []FailureClassification {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var result []FailureClassification
	for _, fc := range t.classifications {
		if fc.ModelID == modelID {
			result = append(result, *fc)
		}
	}
	return result
}

// GetByType returns all classifications of a failure type.
func (t *Taxonomy) GetByType(failureType string) []FailureClassification {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var result []FailureClassification
	for _, fc := range t.classifications {
		if fc.FailureType == failureType {
			result = append(result, *fc)
		}
	}
	return result
}

// GetStats returns summary statistics.
func (t *Taxonomy) GetStats() TaxonomyStats {
	t.mu.RLock()
	defer t.mu.RUnlock()

	stats := TaxonomyStats{
		TotalClassifications: len(t.classifications),
		TotalByModel:        make(map[string]int),
		TotalByType:         make(map[string]int),
		TotalByLanguage:     make(map[string]int),
		TotalBySeverity:     make(map[string]int),
	}

	for _, fc := range t.classifications {
		stats.TotalByModel[fc.ModelID]++
		stats.TotalByType[fc.FailureType]++
		stats.TotalByLanguage[fc.Language]++
		stats.TotalBySeverity[fc.Severity]++
	}

	return stats
}

// GetAll returns all classifications.
func (t *Taxonomy) GetAll() []FailureClassification {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make([]FailureClassification, 0, len(t.classifications))
	for _, fc := range t.classifications {
		result = append(result, *fc)
	}
	return result
}
