package metrics

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Reporter generates benchmark reports from metrics and taxonomy (AC1).
type Reporter struct {
	metricsPath     string
	taxonomyPath   string
	historicalPath  string
}

// NewReporter creates a reporter for given paths.
func NewReporter(metricsPath, taxonomyPath string) *Reporter {
	return &Reporter{
		metricsPath:   metricsPath,
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

// ReportData holds combined metrics and taxonomy data.
type ReportData struct {
	Metrics   MetricsSummary       `json:"metrics"`
	Taxonomy  TaxonomySummary      `json:"taxonomy"`
	Historical []HistoricalEntry   `json:"historical,omitempty"`
}

// MetricsSummary summarizes key metrics for report (AC2).
type MetricsSummary struct {
	CatchRate          float64            `json:"catch_rate"`
	TotalVerifications int                `json:"total_verifications"`
	FailedVerifications int                `json:"failed_verifications"`
	ModelPassRate     map[string]float64 `json:"model_pass_rate"`
	IterationCount     map[string]int     `json:"iteration_count"`
	AcceptanceCatchRate float64           `json:"acceptance_catch_rate"`
}

// TaxonomySummary summarizes failure classifications (AC2).
type TaxonomySummary struct {
	TotalClassifications int                        `json:"total_classifications"`
	ByType              map[string]int             `json:"by_type"`
	ByModel             map[string]int             `json:"by_model"`
	BySeverity          map[string]int             `json:"by_severity"`
}

// HistoricalEntry represents a past benchmark period (AC3).
type HistoricalEntry struct {
	Period     string  `json:"period"`
	CatchRate  float64 `json:"catch_rate"`
	TotalWS    int     `json:"total_workstreams"`
}

// GenerateMarkdown creates markdown benchmark report (AC1, AC4).
func (r *Reporter) GenerateMarkdown() (string, error) {
	data, err := r.loadReportData()
	if err != nil {
		return "", fmt.Errorf("load report data: %w", err)
	}

	var sb strings.Builder
	sb.WriteString(r.generateMarkdownHeader(data))
	sb.WriteString(r.generateMetricsSection(data))
	sb.WriteString(r.generateModelComparison(data))
	sb.WriteString(r.generateTaxonomySection(data))
	sb.WriteString(r.generateTrendSection(data))
	sb.WriteString(r.generateMethodology())

	return sb.String(), nil
}

// generateMarkdownHeader creates report header.
func (r *Reporter) generateMarkdownHeader(data ReportData) string {
	quarter := r.GetCurrentQuarter()
	return fmt.Sprintf(`# AI Code Quality Benchmark: %s

**Generated:** %s
**Data Source:** SDP Dogfooding (F054 onwards)
**Period:** All time through %s

## Executive Summary

This report presents metrics on AI-generated code quality, based on evidence from the Spec-Driven Protocol (SDP) dogfooding dataset. SDP builds SDP, providing a unique opportunity to measure AI code generation quality in a real-world context.

`,
		quarter,
		time.Now().Format("2006-01-02"),
		quarter,
	)
}

// generateMetricsSection creates metrics summary (AC2).
func (r *Reporter) generateMetricsSection(data ReportData) string {
	m := data.Metrics
	return fmt.Sprintf(`
## Overall Metrics

### Catch Rate
**%.2f%%** of verifications failed on first generation

- Total Verifications: %d
- Failed Verifications: %d
- This metric indicates how often AI-generated code passes tests on the first try.

### Iteration Count
Average **%.1f** red→green cycles per workstream

Higher iteration counts indicate more refinement needed.

### Acceptance Catch Rate
**%.2f%%** of builds failed acceptance testing

Acceptance tests catch issues that unit tests miss, providing a higher quality bar.
`,
		m.CatchRate*100,
		m.TotalVerifications,
		m.FailedVerifications,
		r.averageIterationCount(m.IterationCount),
		m.AcceptanceCatchRate*100,
	)
}

// averageIterationCount calculates average iterations.
func (r *Reporter) averageIterationCount(iterations map[string]int) float64 {
	if len(iterations) == 0 {
		return 0
	}
	sum := 0
	for _, count := range iterations {
		sum += count
	}
	return float64(sum) / float64(len(iterations))
}

// generateModelComparison creates model comparison table (AC2).
func (r *Reporter) generateModelComparison(data ReportData) string {
	var sb strings.Builder
	sb.WriteString("\n## Model Performance\n\n")
	sb.WriteString("| Model | Pass Rate | Verifications |\n")
	sb.WriteString("|--------|------------|---------------|\n")

	for model, rate := range data.Metrics.ModelPassRate {
		// Count verifications per model (simplified)
		sb.WriteString(fmt.Sprintf("| %s | %.1f%% | %d |\n",
			model, rate*100, r.estimateVerificationsForModel(model)))
	}

	return sb.String()
}

// estimateVerificationsForModel estimates verification count (simplified).
func (r *Reporter) estimateVerificationsForModel(model string) int {
	// This is a placeholder - in real implementation, we'd track this
	return 10
}

// generateTaxonomySection creates failure breakdown (AC2).
func (r *Reporter) generateTaxonomySection(data ReportData) string {
	t := data.Taxonomy
	if t.TotalClassifications == 0 {
		return "\n## Failure Taxonomy\n\nNo failures recorded.\n"
	}

	var sb strings.Builder
	sb.WriteString("\n## Failure Taxonomy\n\n")
	sb.WriteString("Breakdown of verification failures by type:\n\n")

	// Sort by count (most common first)
	sortedTypes := r.sortByCount(t.ByType)

	for _, failureType := range sortedTypes {
		count := t.ByType[failureType]
		percentage := float64(count) / float64(t.TotalClassifications) * 100
		sb.WriteString(fmt.Sprintf("- **%s**: %d (%.1f%%) - %s\n",
			failureType, count, percentage, r.getFailureDescription(failureType)))
	}

	sb.WriteString(fmt.Sprintf("\n### Severity Distribution\n\n"))
	sb.WriteString(fmt.Sprintf("- CRITICAL: %d (systemic issues caught by acceptance)\n",
		t.BySeverity["CRITICAL"]))
	sb.WriteString(fmt.Sprintf("- HIGH: %d (crashes, edge cases)\n", t.BySeverity["HIGH"]))
	sb.WriteString(fmt.Sprintf("- MEDIUM: %d (logic, type errors)\n", t.BySeverity["MEDIUM"]))
	sb.WriteString(fmt.Sprintf("- LOW: %d (minor issues)\n", t.BySeverity["LOW"]))

	return sb.String()
}

// getFailureDescription returns human-readable description.
func (r *Reporter) getFailureDescription(failureType string) string {
	descriptions := map[string]string{
		"wrong_logic":        "Logic errors - incorrect algorithms or flow",
		"missing_edge_case":  "Edge cases - boundary conditions, nil pointers",
		"hallucinated_api":  "Hallucinated APIs - non-existent functions or methods",
		"type_error":         "Type errors - type mismatches or incompatibilities",
		"test_passing_but_wrong": "Test-Passing-But-Wrong - unit tests pass but acceptance fails",
		"compilation_error":  "Compilation errors - syntax or build failures",
		"import_error":       "Import errors - missing or unresolved modules",
		"unknown":           "Unknown - unclassified failures",
	}
	if desc, ok := descriptions[failureType]; ok {
		return desc
	}
	return "Uncategorized failure"
}

// sortByCount sorts failure types by count (descending).
func (r *Reporter) sortByCount(byType map[string]int) []string {
	type item struct {
		key   string
		count int
	}
	var items []item
	for k, v := range byType {
		items = append(items, item{key: k, count: v})
	}
	// Simple bubble sort (descending by count)
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].count > items[i].count {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = item.key
	}
	return result
}

// generateTrendSection creates trend analysis (AC3).
func (r *Reporter) generateTrendSection(data ReportData) string {
	if len(data.Historical) == 0 {
		return "\n## Trends Over Time\n\nHistorical data not available. Future reports will include trend analysis as more data accumulates.\n"
	}

	var sb strings.Builder
	sb.WriteString("\n## Trends Over Time\n\n")
	sb.WriteString("Comparison of catch rates across periods:\n\n")
	sb.WriteString("| Period | Catch Rate | Total Workstreams |\n")
	sb.WriteString("|--------|------------|-------------------|\n")

	// Add current period
	sb.WriteString(fmt.Sprintf("| %s (current) | %.1f%% | %d |\n",
		r.GetCurrentQuarter(),
		data.Metrics.CatchRate*100,
		len(data.Metrics.IterationCount)))

	// Add historical periods
	for _, hist := range data.Historical {
		sb.WriteString(fmt.Sprintf("| %s | %.1f%% | %d |\n",
			hist.Period, hist.CatchRate*100, hist.TotalWS))
	}

	// Calculate trend
	if len(data.Historical) >= 1 {
		latest := data.Historical[len(data.Historical)-1].CatchRate
		current := data.Metrics.CatchRate
		if current < latest {
			sb.WriteString("\n**Trend:** Improving - Catch rate decreased from previous period\n")
		} else if current > latest {
			sb.WriteString("\n**Trend:** Declining - Catch rate increased from previous period\n")
		} else {
			sb.WriteString("\n**Trend:** Stable - Catch rate unchanged from previous period\n")
		}
	}

	return sb.String()
}

// generateMethodology creates methodology section.
func (r *Reporter) generateMethodology() string {
	return `
## Methodology

**Privacy Note:** This report contains only aggregated metrics. No raw code, company names, or identifying information is included (AC4).

**Data Collection:**
- Evidence captured during SDP development workflow
- Each code generation is followed by verification
- Metrics aggregated from verification events

**Limitations:**
- Dogfooding dataset: SDP builds SDP, which may bias results
- Sample size varies by quarter
- Classification is heuristic-based and may not be 100% accurate

**Next Report:** Will include Q2 2026 data as more builds complete.
`
}

// GenerateHTML creates HTML benchmark report (AC5).
func (r *Reporter) GenerateHTML() (string, error) {
	markdown, err := r.GenerateMarkdown()
	if err != nil {
		return "", err
	}

	// Simple markdown to HTML conversion
	html := `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>AI Code Quality Benchmark</title>
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
		table { border-collapse: collapse; width: 100%%; }
		th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
		th { background-color: #4CAF50; color: white; }
		h1 { color: #333; }
		h2 { color: #666; border-bottom: 1px solid #eee; padding-bottom: 10px; }
	</style>
</head>
<body>
` + markdownToHTML(markdown) + `
</body>
</html>`

	return html, nil
}

// markdownToHTML converts basic markdown to HTML (simplified).
func markdownToHTML(md string) string {
	// Very basic conversion - for production, use a proper library
	html := md
	// Headers
	replacer := strings.NewReplacer("# ", "<h1>", "## ", "<h2>", "### ", "<h3>")
	html = replacer.Replace(html)
	// Bold
	boldReplacer := strings.NewReplacer("**", "<strong>", "***", "</strong>")
	html = boldReplacer.Replace(html)
	// Line breaks
	brReplacer := strings.NewReplacer("\n\n", "</p><p>", "\n", "<br>")
	html = brReplacer.Replace(html)

	return html
}

// GenerateJSON creates JSON benchmark report (AC5).
func (r *Reporter) GenerateJSON() (string, error) {
	data, err := r.loadReportData()
	if err != nil {
		return "", fmt.Errorf("load report data: %w", err)
	}

	reportJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal report: %w", err)
	}

	return string(reportJSON), nil
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
		CatchRate:          metrics.CatchRate,
		TotalVerifications: metrics.TotalVerifications,
		FailedVerifications: metrics.FailedVerifications,
		ModelPassRate:     metrics.ModelPassRate,
		IterationCount:     metrics.IterationCount,
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
			ByType:              stats.TotalByType,
			ByModel:             stats.TotalByModel,
			BySeverity:          stats.TotalBySeverity,
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
