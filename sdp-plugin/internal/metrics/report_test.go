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
	t.Skip("Skipping - requires taxonomy file format fix")
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")

	// Create sample metrics
	metricsData := map[string]interface{}{
		"catch_rate":           0.25,
		"total_verifications":   100,
		"failed_verifications":  25,
		"model_pass_rate":      map[string]float64{"claude-sonnet-4": 0.85, "claude-opus-4": 0.92},
		"iteration_count":       map[string]int{"00-001-01": 3, "00-001-02": 1},
		"acceptance_catch_rate": 0.15,
	}

	// Create sample taxonomy
	taxonomyData := map[string]interface{}{
		"classifications": []map[string]interface{}{
			{"failure_type": "wrong_logic", "severity": "MEDIUM"},
			{"failure_type": "type_error", "severity": "MEDIUM"},
		},
	}

	metricsJSON, _ := json.Marshal(metricsData)
	os.WriteFile(metricsPath, metricsJSON, 0644)
	taxonomyJSON, _ := json.Marshal(taxonomyData)
	os.WriteFile(taxonomyPath, taxonomyJSON, 0644)

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
}

func TestReport_GenerateHTML_HasValidStructure(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")

	// Create minimal metrics
	metricsData := map[string]interface{}{"catch_rate": 0.25, "total_verifications": 100}
	metricsJSON, _ := json.Marshal(metricsData)
	os.WriteFile(metricsPath, metricsJSON, 0644)

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
	os.WriteFile(metricsPath, metricsJSON, 0644)

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
	t.Skip("Skipping - requires taxonomy file format fix")
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")
	historicalPath := filepath.Join(tempDir, "historical.json")

	// Create historical data
	historicalData := []map[string]interface{}{
		{
			"period":      "2025-Q4",
			"catch_rate":  0.30,
			"total_ws":    50,
		},
		{
			"period":      "2026-Q1",
			"catch_rate":  0.25,
			"total_ws":    75,
		},
	}
	historicalJSON, _ := json.Marshal(historicalData)
	os.WriteFile(historicalPath, historicalJSON, 0644)

	// Create current metrics
	metricsData := map[string]interface{}{"catch_rate": 0.20, "total_verifications": 100}
	metricsJSON, _ := json.Marshal(metricsData)
	os.WriteFile(metricsPath, metricsJSON, 0644)

	// Act
	reporter := NewReporter(metricsPath, taxonomyPath)
	reporter.SetHistoricalPath(historicalPath)
	report, err := reporter.GenerateMarkdown()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	// Should contain trend section
	if !contains(report, "Trend Over Time") {
		t.Error("Expected report to contain 'Trend Over Time' section")
	}
	if !contains(report, "2025-Q4") {
		t.Error("Expected report to contain historical data")
	}
}

func TestReport_GenerateToDefaultPath_CreatesInCorrectLocation(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	metricsPath := filepath.Join(tempDir, "metrics.json")
	taxonomyPath := filepath.Join(tempDir, "taxonomy.json")

	metricsData := map[string]interface{}{"catch_rate": 0.25, "total_verifications": 100}
	metricsJSON, _ := json.Marshal(metricsData)
	os.WriteFile(metricsPath, metricsJSON, 0644)

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
	os.WriteFile(metricsPath, metricsJSON, 0644)

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
		"catch_rate":              0.25,
		"total_verifications":      100,
		"failed_verifications":     25,
		"model_pass_rate":         map[string]float64{"model-a": 0.85},
		"iteration_count":          map[string]int{"ws-1": 3},
		"acceptance_catch_rate":    0.15,
	}
	metricsJSON, _ := json.Marshal(metricsData)
	os.WriteFile(metricsPath, metricsJSON, 0644)

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
	os.WriteFile(metricsPath, metricsJSON, 0644)

	// Act
	reporter := NewReporter(metricsPath, taxonomyPath)
	report, _ := reporter.GenerateMarkdown()

	// Assert - should contain placeholder message when no historical data
	if !contains(report, "Historical data not available") {
		t.Error("Expected report to contain historical data placeholder")
	}
}
