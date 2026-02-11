package metrics

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
	TotalClassifications int            `json:"total_classifications"`
	TotalByModel         map[string]int `json:"total_by_model"`
	TotalByType          map[string]int `json:"total_by_type"`
	TotalByLanguage      map[string]int `json:"total_by_language"`
	TotalBySeverity      map[string]int `json:"total_by_severity"`
}

// FailureType enum (AC2).
const (
	FailureWrongLogic       = "wrong_logic"
	FailureMissingEdgeCase  = "missing_edge_case"
	FailureHallucinatedAPI  = "hallucinated_api"
	FailureTypeError        = "type_error"
	FailureTestPassingWrong = "test_passing_but_wrong"
	FailureCompilationError = "compilation_error"
	FailureImportError      = "import_error"
	FailureUnknown          = "unknown"
)

// Severity levels.
const (
	SeverityLow      = "LOW"
	SeverityMedium   = "MEDIUM"
	SeverityHigh     = "HIGH"
	SeverityCritical = "CRITICAL"
)
