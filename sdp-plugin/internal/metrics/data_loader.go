package metrics

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReportData holds combined metrics and taxonomy data.
type ReportData struct {
	Metrics    MetricsSummary     `json:"metrics"`
	Taxonomy   TaxonomySummary    `json:"taxonomy"`
	Historical []HistoricalEntry  `json:"historical,omitempty"`
}

// MetricsSummary summarizes key metrics for report (AC2).
type MetricsSummary struct {
	CatchRate           float64            `json:"catch_rate"`
	TotalVerifications  int                `json:"total_verifications"`
	FailedVerifications int                `json:"failed_verifications"`
	ModelPassRate       map[string]float64 `json:"model_pass_rate"`
	IterationCount      map[string]int     `json:"iteration_count"`
	AcceptanceCatchRate float64            `json:"acceptance_catch_rate"`
}

// TaxonomySummary summarizes failure classifications (AC2).
type TaxonomySummary struct {
	TotalClassifications int            `json:"total_classifications"`
	ByType              map[string]int `json:"by_type"`
	ByModel             map[string]int `json:"by_model"`
	BySeverity          map[string]int `json:"by_severity"`
}

// HistoricalEntry represents a past benchmark period (AC3).
type HistoricalEntry struct {
	Period    string  `json:"period"`
	CatchRate float64 `json:"catch_rate"`
	TotalWS   int     `json:"total_workstreams"`
}

// loadReportData loads metrics and taxonomy for report generation.
func (r *Reporter) loadReportData() (ReportData, error) {
	data := ReportData{}

	// Load metrics
	metricsFile, err := os.ReadFile(r.metricsPath)
	if err != nil {
		return data, fmt.Errorf("read metrics: %w", err)
	}

	var metrics Metrics
	if err := json.Unmarshal(metricsFile, &metrics); err != nil {
		return data, fmt.Errorf("parse metrics: %w", err)
	}

	data.Metrics = MetricsSummary{
		CatchRate:           metrics.CatchRate,
		TotalVerifications:  metrics.TotalVerifications,
		FailedVerifications: metrics.FailedVerifications,
		ModelPassRate:       metrics.ModelPassRate,
		IterationCount:      metrics.IterationCount,
		AcceptanceCatchRate: metrics.AcceptanceCatchRate,
	}

	// Load taxonomy
	taxonomyFile, err := os.ReadFile(r.taxonomyPath)
	if err != nil && !os.IsNotExist(err) {
		return data, fmt.Errorf("read taxonomy: %w", err)
	}

	if len(taxonomyFile) > 0 {
		taxonomy := NewTaxonomy(r.taxonomyPath)
		if err := taxonomy.Load(); err != nil {
			return data, fmt.Errorf("load taxonomy: %w", err)
		}

		stats := taxonomy.GetStats()
		data.Taxonomy = TaxonomySummary{
			TotalClassifications: stats.TotalClassifications,
			ByType:               stats.TotalByType,
			ByModel:              stats.TotalByModel,
			BySeverity:           stats.TotalBySeverity,
		}
	}

	// Load historical if path set
	if r.historicalPath != "" {
		historicalFile, err := os.ReadFile(r.historicalPath)
		if err != nil && !os.IsNotExist(err) {
			return data, fmt.Errorf("read historical: %w", err)
		}

		var historical []HistoricalEntry
		if err := json.Unmarshal(historicalFile, &historical); err != nil {
			return data, fmt.Errorf("parse historical: %w", err)
		}
		data.Historical = historical
	}

	return data, nil
}
