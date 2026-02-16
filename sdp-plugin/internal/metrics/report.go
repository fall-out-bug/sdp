package metrics

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Reporter generates benchmark reports from metrics and taxonomy (AC1).
type Reporter struct {
	metricsPath    string
	taxonomyPath   string
	historicalPath string
}

// NewReporter creates a reporter for given paths.
func NewReporter(metricsPath, taxonomyPath string) *Reporter {
	return &Reporter{
		metricsPath:  metricsPath,
		taxonomyPath: taxonomyPath,
	}
}

// SetHistoricalPath sets path to historical benchmark data (AC3).
func (r *Reporter) SetHistoricalPath(path string) {
	r.historicalPath = path
}

// GetDefaultOutputPath returns the default report path (AC6).
func (r *Reporter) GetDefaultOutputPath() string {
	quarter := r.GetCurrentQuarter()
	return fmt.Sprintf(".sdp/metrics/benchmark-%s.md", quarter)
}

// GetCurrentQuarter returns current quarter identifier (e.g., "2026-Q1").
func (r *Reporter) GetCurrentQuarter() string {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	quarter := (month-1)/3 + 1
	return fmt.Sprintf("%d-Q%d", year, quarter)
}

// Save writes the report to the default location.
func (r *Reporter) Save() error {
	report, err := r.GenerateMarkdown()
	if err != nil {
		return fmt.Errorf("generate report: %w", err)
	}

	outputPath := r.GetDefaultOutputPath()
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	return os.WriteFile(outputPath, []byte(report), 0644)
}
