package metrics

import (
	"regexp"
	"strings"
)

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
