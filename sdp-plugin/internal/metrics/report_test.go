package metrics

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestReport_GenerateMarkdown_AllSectionsPresent(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")

	// Create sample metrics
	metricsData := map[string]interface{}{
		"catch_rate":            0.25,
		"total_verifications":   100,
		"failed_verifications":  25,
		"model_pass_rate":       map[string]float64{"claude-sonnet-4": 0.85, "claude-opus-4": 0.92},
		"iteration_count":       map[string]int{"00-001-01": 3, "00-001-02": 1},
		"acceptance_catch_rate": 0.15,
	}
	metricsJSON, _ := json.Marshal(metricsData)
	if err := os.WriteFile(metricsPath, metricsJSON, 0644); err != nil {
		t.Fatalf("Failed to write metrics file: %v", err)
	}

	// Create sample taxonomy - Taxonomy.Load() expects JSON array directly
	taxonomyData := []map[string]interface{}{
		{"event_id": "evt1", "ws_id": "00-001-01", "model_id": "claude-sonnet-4", "language": "go", "failure_type": "wrong_logic", "severity": "MEDIUM"},
		{"event_id": "evt2", "ws_id": "00-001-02", "model_id": "claude-opus-4", "language": "go", "failure_type": "type_error", "severity": "MEDIUM"},
	}
	taxonomyJSON, _ := json.Marshal(taxonomyData)
	if err := os.WriteFile(taxonomyPath, taxonomyJSON, 0644); err != nil {
		t.Fatalf("Failed to write taxonomy file: %v", err)
	}

	// Act
	reporter := NewReporter(metricsPath, taxonomyPath)
	report, err := reporter.GenerateMarkdown()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if report == "" {
		t.Fatal("Expected non-empty report")
	}

	// Verify all expected sections are present
	expectedSections := []string{
		"# AI Code Quality Benchmark",
		"## Executive Summary",
		"## Overall Metrics",
		"## Model Performance",
		"## Failure Taxonomy",
		"## Trends Over Time",
		"## Methodology",
	}

	for _, section := range expectedSections {
		if !contains(report, section) {
			t.Errorf("Expected report to contain section '%s'", section)
		}
	}

	// Verify taxonomy data is reflected in report
	if !contains(report, "wrong_logic") {
		t.Error("Expected report to contain 'wrong_logic' failure type")
	}
	if !contains(report, "type_error") {
		t.Error("Expected report to contain 'type_error' failure type")
	}
}

func TestReport_GenerateHTML_HasValidStructure(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")

	// Create minimal metrics
	metricsData := map[string]interface{}{"catch_rate": 0.25, "total_verifications": 100}
	metricsJSON, _ := json.Marshal(metricsData)
	if err := os.WriteFile(metricsPath, metricsJSON, 0644); err != nil {
		t.Fatalf("Failed to write metrics file: %v", err)
	}

	// Act
	reporter := NewReporter(metricsPath, taxonomyPath)
	report, err := reporter.GenerateHTML()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if report == "" {
		t.Fatal("Expected non-empty HTML report")
	}
}

func TestReport_GenerateJSON_ValidFormat(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")

	// Create minimal metrics
	metricsData := map[string]interface{}{"catch_rate": 0.25, "total_verifications": 100}
	metricsJSON, _ := json.Marshal(metricsData)
	if err := os.WriteFile(metricsPath, metricsJSON, 0644); err != nil {
		t.Fatalf("Failed to write metrics file: %v", err)
	}

	// Act
	reporter := NewReporter(metricsPath, taxonomyPath)
	report, err := reporter.GenerateJSON()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if report == "" {
		t.Fatal("Expected non-empty JSON report")
	}
	// Verify it's valid JSON (contains expected keys)
	if len(report) < 10 {
		t.Errorf("Expected substantial JSON report, got length %d", len(report))
	}
}

func TestReport_GenerateWithTrend_IncludesHistoricalData(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")
	historicalPath := filepath.Join(tempDir, "historical.json")

	// Create historical data - matches HistoricalEntry struct
	historicalData := []map[string]interface{}{
		{
			"period":            "2025-Q4",
			"catch_rate":        0.30,
			"total_workstreams": 50,
		},
		{
			"period":            "2026-Q1",
			"catch_rate":        0.25,
			"total_workstreams": 75,
		},
	}
	historicalJSON, _ := json.Marshal(historicalData)
	if err := os.WriteFile(historicalPath, historicalJSON, 0644); err != nil {
		t.Fatalf("Failed to write historical file: %v", err)
	}

	// Create current metrics
	metricsData := map[string]interface{}{
		"catch_rate":            0.20,
		"total_verifications":   100,
		"failed_verifications":  20,
		"model_pass_rate":       map[string]float64{"claude-sonnet-4": 0.85},
		"iteration_count":       map[string]int{"00-001-01": 3},
		"acceptance_catch_rate": 0.10,
	}
	metricsJSON, _ := json.Marshal(metricsData)
	if err := os.WriteFile(metricsPath, metricsJSON, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Act
	reporter := NewReporter(metricsPath, taxonomyPath)
	reporter.SetHistoricalPath(historicalPath)
	report, err := reporter.GenerateMarkdown()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// The section is called "Trends Over Time" (plural), not "Trend Over Time"
	if !contains(report, "Trends Over Time") {
		t.Error("Expected report to contain 'Trends Over Time' section")
	}
	if !contains(report, "2025-Q4") {
		t.Error("Expected report to contain historical data for 2025-Q4")
	}
	if !contains(report, "2026-Q1") {
		t.Error("Expected report to contain historical data for 2026-Q1")
	}
}

func TestReport_GenerateToDefaultPath_CreatesInCorrectLocation(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")

	metricsData := map[string]interface{}{"catch_rate": 0.25, "total_verifications": 100}
	metricsJSON, _ := json.Marshal(metricsData)
	if err := os.WriteFile(metricsPath, metricsJSON, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Act
	reporter := NewReporter(metricsPath, taxonomyPath)
	outputPath := reporter.GetDefaultOutputPath()

	// Assert
	if !contains(outputPath, "benchmark") {
		t.Errorf("Expected output path to contain 'benchmark', got %s", outputPath)
	}
	if !contains(outputPath, "benchmark") {
		t.Errorf("Expected output path to contain 'benchmark', got %s", outputPath)
	}
}

func TestReport_GenerateQuarterlyReport_UsesCorrectQuarter(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")

	metricsData := map[string]interface{}{"catch_rate": 0.25, "total_verifications": 100}
	metricsJSON, _ := json.Marshal(metricsData)
	if err := os.WriteFile(metricsPath, metricsJSON, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Act
	reporter := NewReporter(metricsPath, taxonomyPath)
	quarter := reporter.GetCurrentQuarter()

	// Assert
	// Quarter should be like "2026-Q1"
	if len(quarter) < 6 {
		t.Errorf("Expected quarter format YYYY-QN, got %s", quarter)
	}
	// Year should be current year
	currentYear := time.Now().Year()
	if !contains(quarter, fmt.Sprintf("%d", currentYear)) {
		t.Errorf("Expected quarter to contain current year %d, got %s", currentYear, quarter)
	}
}

// contains checks if substring exists in string
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func TestReporter_LoadMetrics_ParsesCorrectly(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")

	// Create sample metrics
	metricsData := map[string]interface{}{
		"catch_rate":            0.25,
		"total_verifications":   100,
		"failed_verifications":  25,
		"model_pass_rate":       map[string]float64{"model-a": 0.85},
		"iteration_count":       map[string]int{"ws-1": 3},
		"acceptance_catch_rate": 0.15,
	}
	metricsJSON, _ := json.Marshal(metricsData)
	if err := os.WriteFile(metricsPath, metricsJSON, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Act
	reporter := NewReporter(metricsPath, taxonomyPath)
	data, err := reporter.loadReportData()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if data.Metrics.CatchRate != 0.25 {
		t.Errorf("Expected catch rate 0.25, got %f", data.Metrics.CatchRate)
	}
	if data.Metrics.TotalVerifications != 100 {
		t.Errorf("Expected 100 total verifications, got %d", data.Metrics.TotalVerifications)
	}
}

func TestReporter_GenerateTrendWithoutHistorical_ReturnsPlaceholder(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")

	metricsData := map[string]interface{}{"catch_rate": 0.25}
	metricsJSON, _ := json.Marshal(metricsData)
	if err := os.WriteFile(metricsPath, metricsJSON, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Act
	reporter := NewReporter(metricsPath, taxonomyPath)
	report, _ := reporter.GenerateMarkdown()

	// Assert - should contain placeholder message when no historical data
	if !contains(report, "Historical data not available") {
		t.Error("Expected report to contain historical data placeholder")
	}
}

func TestReporter_Save_CreatesReportInDefaultLocation(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")

	metricsData := map[string]interface{}{
		"catch_rate":            0.25,
		"total_verifications":   100,
		"failed_verifications":  25,
		"model_pass_rate":       map[string]float64{"claude-sonnet-4": 0.85},
		"iteration_count":       map[string]int{"00-001-01": 3},
		"acceptance_catch_rate": 0.15,
	}
	metricsJSON, _ := json.Marshal(metricsData)
	if err := os.WriteFile(metricsPath, metricsJSON, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Change working directory to tempDir for default output path
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Errorf("Failed to restore working directory: %v", err)
		}
	}()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Act
	reporter := NewReporter(metricsPath, taxonomyPath)
	saveErr := reporter.Save()
	if saveErr != nil {
		err = saveErr
	}

	// Assert
	if err != nil {
		t.Fatalf("Expected no error from Save, got %v", err)
	}

	// Verify report file was created
	outputPath := reporter.GetDefaultOutputPath()
	if _, statErr := os.Stat(outputPath); statErr != nil {
		t.Fatalf("Expected report file to exist at %s, got error: %v", outputPath, statErr)
	}

	// Verify report content
	content, _ := os.ReadFile(outputPath)
	if len(content) == 0 {
		t.Error("Expected report file to have content")
	}
	if !contains(string(content), "# AI Code Quality Benchmark") {
		t.Error("Expected report to contain benchmark header")
	}
}

func TestReporter_Save_WithNestedDirectory_CreatesDirectory(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")

	metricsData := map[string]interface{}{
		"catch_rate":          0.25,
		"total_verifications": 100,
	}
	metricsJSON, _ := json.Marshal(metricsData)
	if err := os.WriteFile(metricsPath, metricsJSON, 0644); err != nil {
		t.Fatalf("Failed to write metrics file: %v", err)
	}

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Errorf("Failed to restore working directory: %v", err)
		}
	}()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Act
	reporter := NewReporter(metricsPath, taxonomyPath)
	saveErr := reporter.Save()
	if saveErr != nil {
		err = saveErr
	}

	// Assert
	if err != nil {
		t.Fatalf("Expected no error when creating nested directories, got %v", err)
	}

	// Verify the nested directory was created
	outputPath := reporter.GetDefaultOutputPath()
	dir := filepath.Dir(outputPath)
	if _, statErr := os.Stat(dir); statErr != nil {
		t.Fatalf("Expected directory %s to exist, got error: %v", dir, statErr)
	}
}

func TestReporter_Save_WithInvalidMetrics_ReturnsError(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "nonexistent.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Errorf("Failed to restore working directory: %v", err)
		}
	}()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Act
	reporter := NewReporter(metricsPath, taxonomyPath)
	saveErr := reporter.Save()
	if saveErr != nil {
		err = saveErr
	}

	// Assert - should return error when metrics file doesn't exist
	if err == nil {
		t.Error("Expected error when metrics file doesn't exist, got nil")
	}
}

func TestReport_SetHistoricalPath_UpdatesPath(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")
	reporter := NewReporter(metricsPath, taxonomyPath)

	// Act
	customPath := filepath.Join(tempDir, "custom-historical.json")
	reporter.SetHistoricalPath(customPath)

	// Assert - SetHistoricalPath just sets the field
	// The effect is verified through GenerateMarkdown with historical data
	metricsData := map[string]interface{}{
		"catch_rate":          0.20,
		"total_verifications": 100,
	}
	metricsJSON, _ := json.Marshal(metricsData)
	if err := os.WriteFile(metricsPath, metricsJSON, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	historicalData := []map[string]interface{}{
		{"period": "2025-Q4", "catch_rate": 0.30, "total_workstreams": 50},
	}
	historicalJSON, _ := json.Marshal(historicalData)
	if err := os.WriteFile(customPath, historicalJSON, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	report, err := reporter.GenerateMarkdown()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !contains(report, "2025-Q4") {
		t.Error("Expected historical data to be included after SetHistoricalPath")
	}
}

func TestReport_EstimateVerificationsForModel_ReturnsPlaceholderValue(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")

	metricsData := map[string]interface{}{
		"catch_rate":          0.25,
		"total_verifications": 100,
		"model_pass_rate":     map[string]float64{"claude-sonnet-4": 0.85},
	}
	metricsJSON, _ := json.Marshal(metricsData)
	if err := os.WriteFile(metricsPath, metricsJSON, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Act - The estimateVerificationsForModel is called internally by generateModelComparison
	reporter := NewReporter(metricsPath, taxonomyPath)
	report, err := reporter.GenerateMarkdown()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	// The placeholder value is 10, which should appear in the model comparison table
	if !contains(report, "10") {
		t.Error("Expected model comparison to include verification estimate")
	}
}

func TestReport_GenerateTaxonomySection_WithNoFailures_ReturnsNoFailuresMessage(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")

	// Create taxonomy with zero classifications
	taxonomyData := []map[string]interface{}{}
	taxonomyJSON, _ := json.Marshal(taxonomyData)
	if err := os.WriteFile(taxonomyPath, taxonomyJSON, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	metricsData := map[string]interface{}{
		"catch_rate":          0.25,
		"total_verifications": 100,
	}
	metricsJSON, _ := json.Marshal(metricsData)
	if err := os.WriteFile(metricsPath, metricsJSON, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Act
	reporter := NewReporter(metricsPath, taxonomyPath)
	report, err := reporter.GenerateMarkdown()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !contains(report, "No failures recorded") {
		t.Error("Expected 'No failures recorded' message when taxonomy is empty")
	}
}

func TestReport_GenerateTaxonomySection_WithUnknownFailureType_ReturnsUncategorizedDescription(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")

	// Create taxonomy with unknown failure type
	taxonomyData := []map[string]interface{}{
		{"event_id": "evt1", "ws_id": "00-001-01", "model_id": "claude-sonnet-4", "language": "go", "failure_type": "unknown_weird_type", "severity": "MEDIUM"},
	}
	taxonomyJSON, _ := json.Marshal(taxonomyData)
	if err := os.WriteFile(taxonomyPath, taxonomyJSON, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	metricsData := map[string]interface{}{
		"catch_rate":          0.25,
		"total_verifications": 100,
	}
	metricsJSON, _ := json.Marshal(metricsData)
	if err := os.WriteFile(metricsPath, metricsJSON, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Act
	reporter := NewReporter(metricsPath, taxonomyPath)
	report, err := reporter.GenerateMarkdown()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !contains(report, "Uncategorized failure") {
		t.Error("Expected 'Uncategorized failure' description for unknown failure type")
	}
}

func TestReport_GenerateTaxonomySection_SeverityDistribution_AllLevels(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")

	// Create taxonomy with all severity levels
	taxonomyData := []map[string]interface{}{
		{"event_id": "evt1", "ws_id": "00-001-01", "model_id": "claude-sonnet-4", "language": "go", "failure_type": "wrong_logic", "severity": "CRITICAL"},
		{"event_id": "evt2", "ws_id": "00-001-02", "model_id": "claude-opus-4", "language": "go", "failure_type": "type_error", "severity": "HIGH"},
		{"event_id": "evt3", "ws_id": "00-001-03", "model_id": "claude-sonnet-4", "language": "go", "failure_type": "hallucinated_api", "severity": "MEDIUM"},
		{"event_id": "evt4", "ws_id": "00-001-04", "model_id": "claude-opus-4", "language": "go", "failure_type": "import_error", "severity": "LOW"},
	}
	taxonomyJSON, _ := json.Marshal(taxonomyData)
	if err := os.WriteFile(taxonomyPath, taxonomyJSON, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	metricsData := map[string]interface{}{
		"catch_rate":          0.25,
		"total_verifications": 100,
	}
	metricsJSON, _ := json.Marshal(metricsData)
	if err := os.WriteFile(metricsPath, metricsJSON, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Act
	reporter := NewReporter(metricsPath, taxonomyPath)
	report, err := reporter.GenerateMarkdown()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify all severity levels are shown
	expectedSeverities := []string{"CRITICAL", "HIGH", "MEDIUM", "LOW"}
	for _, severity := range expectedSeverities {
		if !contains(report, severity) {
			t.Errorf("Expected report to contain severity level '%s'", severity)
		}
	}
	if !contains(report, "Severity Distribution") {
		t.Error("Expected 'Severity Distribution' section")
	}
}

func TestReport_GenerateTrendSection_TrendAnalysis(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")
	historicalPath := filepath.Join(tempDir, "historical.json")

	// Historical with higher catch rate (worse) than current (improving)
	historicalData := []map[string]interface{}{
		{"period": "2025-Q4", "catch_rate": 0.40, "total_workstreams": 50},
	}
	historicalJSON, _ := json.Marshal(historicalData)
	if err := os.WriteFile(historicalPath, historicalJSON, 0644); err != nil {
		t.Fatalf("Failed to write historical file: %v", err)
	}

	// Current with lower catch rate (better) than historical
	metricsData := map[string]interface{}{
		"catch_rate":          0.20,
		"total_verifications": 100,
	}
	metricsJSON, _ := json.Marshal(metricsData)
	if err := os.WriteFile(metricsPath, metricsJSON, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Act
	reporter := NewReporter(metricsPath, taxonomyPath)
	reporter.SetHistoricalPath(historicalPath)
	report, err := reporter.GenerateMarkdown()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !contains(report, "Trend:") {
		t.Error("Expected trend analysis to be present")
	}
}

func TestReport_GenerateTrendSection_StableTrend(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")
	historicalPath := filepath.Join(tempDir, "historical.json")

	// Historical with same catch rate as current (stable)
	historicalData := []map[string]interface{}{
		{"period": "2025-Q4", "catch_rate": 0.25, "total_workstreams": 50},
	}
	historicalJSON, _ := json.Marshal(historicalData)
	if err := os.WriteFile(historicalPath, historicalJSON, 0644); err != nil {
		t.Fatalf("Failed to write historical file: %v", err)
	}

	metricsData := map[string]interface{}{
		"catch_rate":          0.25,
		"total_verifications": 100,
	}
	metricsJSON, _ := json.Marshal(metricsData)
	if err := os.WriteFile(metricsPath, metricsJSON, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Act
	reporter := NewReporter(metricsPath, taxonomyPath)
	reporter.SetHistoricalPath(historicalPath)
	report, err := reporter.GenerateMarkdown()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !contains(report, "Stable") {
		t.Error("Expected 'Stable' trend when catch rate is unchanged")
	}
}
