package monitoring

import (
	"testing"
	"time"
)

// TestNewMetricsCollector verifies metrics collector creation
func TestNewMetricsCollector(t *testing.T) {
	collector := NewMetricsCollector()
	if collector == nil {
		t.Fatal("Expected non-nil metrics collector")
	}

	if collector.validationsBySeverity == nil {
		t.Error("Expected severity map to be initialized")
	}

	if collector.validationsLatency == nil {
		t.Error("Expected latency slice to be initialized")
	}

	if collector.reportGenLatency == nil {
		t.Error("Expected report latency slice to be initialized")
	}
}

// TestRecordValidation verifies validation recording
func TestRecordValidation(t *testing.T) {
	collector := NewMetricsCollector()

	// Record successful validation
	collector.RecordValidation(true, 100*time.Millisecond, 1, 2, 3)

	snapshot := collector.GetMetrics()
	if snapshot.Validations.Total != 1 {
		t.Errorf("Expected 1 validation, got %d", snapshot.Validations.Total)
	}

	if snapshot.Validations.Success != 1 {
		t.Errorf("Expected 1 success, got %d", snapshot.Validations.Success)
	}

	if snapshot.Validations.Failed != 0 {
		t.Errorf("Expected 0 failures, got %d", snapshot.Validations.Failed)
	}

	if snapshot.SeverityDistribution["ERROR"] != 1 {
		t.Errorf("Expected 1 error, got %d", snapshot.SeverityDistribution["ERROR"])
	}

	if snapshot.SeverityDistribution["WARNING"] != 2 {
		t.Errorf("Expected 2 warnings, got %d", snapshot.SeverityDistribution["WARNING"])
	}

	if snapshot.SeverityDistribution["INFO"] != 3 {
		t.Errorf("Expected 3 info, got %d", snapshot.SeverityDistribution["INFO"])
	}

	// Record failed validation
	collector.RecordValidation(false, 200*time.Millisecond, 0, 0, 0)

	snapshot = collector.GetMetrics()
	if snapshot.Validations.Total != 2 {
		t.Errorf("Expected 2 validations, got %d", snapshot.Validations.Total)
	}

	if snapshot.Validations.Success != 1 {
		t.Errorf("Expected 1 success, got %d", snapshot.Validations.Success)
	}

	if snapshot.Validations.Failed != 1 {
		t.Errorf("Expected 1 failure, got %d", snapshot.Validations.Failed)
	}
}

// TestRecordSchemaParse verifies schema parse recording
func TestRecordSchemaParse(t *testing.T) {
	collector := NewMetricsCollector()

	// Record successful parse
	collector.RecordSchemaParse(true)

	snapshot := collector.GetMetrics()
	if snapshot.SchemaParse.Total != 1 {
		t.Errorf("Expected 1 parse, got %d", snapshot.SchemaParse.Total)
	}

	if snapshot.SchemaParse.Success != 1 {
		t.Errorf("Expected 1 success, got %d", snapshot.SchemaParse.Success)
	}

	// Record failed parse
	collector.RecordSchemaParse(false)

	snapshot = collector.GetMetrics()
	if snapshot.SchemaParse.Total != 2 {
		t.Errorf("Expected 2 parses, got %d", snapshot.SchemaParse.Total)
	}

	if snapshot.SchemaParse.Failed != 1 {
		t.Errorf("Expected 1 failure, got %d", snapshot.SchemaParse.Failed)
	}
}

// TestRecordReportGeneration verifies report generation recording
func TestRecordReportGeneration(t *testing.T) {
	collector := NewMetricsCollector()

	// Record successful generation
	collector.RecordReportGeneration(true, 50*time.Millisecond)

	snapshot := collector.GetMetrics()
	if snapshot.ReportGeneration.Total != 1 {
		t.Errorf("Expected 1 report, got %d", snapshot.ReportGeneration.Total)
	}

	if snapshot.ReportGeneration.Success != 1 {
		t.Errorf("Expected 1 success, got %d", snapshot.ReportGeneration.Success)
	}

	// Record failed generation
	collector.RecordReportGeneration(false, 100*time.Millisecond)

	snapshot = collector.GetMetrics()
	if snapshot.ReportGeneration.Total != 2 {
		t.Errorf("Expected 2 reports, got %d", snapshot.ReportGeneration.Total)
	}

	if snapshot.ReportGeneration.Failed != 1 {
		t.Errorf("Expected 1 failure, got %d", snapshot.ReportGeneration.Failed)
	}
}

// TestGetMetrics verifies metrics snapshot
func TestGetMetrics(t *testing.T) {
	collector := NewMetricsCollector()

	// Record some data
	collector.RecordValidation(true, 100*time.Millisecond, 1, 0, 0)
	collector.RecordValidation(true, 200*time.Millisecond, 0, 1, 0)
	collector.RecordValidation(false, 300*time.Millisecond, 0, 0, 1)
	collector.RecordSchemaParse(true)
	collector.RecordSchemaParse(false)
	collector.RecordReportGeneration(true, 50*time.Millisecond)

	snapshot := collector.GetMetrics()

	// Check validation metrics
	if snapshot.Validations.Total != 3 {
		t.Errorf("Expected 3 validations, got %d", snapshot.Validations.Total)
	}

	if snapshot.Validations.Success != 2 {
		t.Errorf("Expected 2 successes, got %d", snapshot.Validations.Success)
	}

	expectedSuccessRate := 66.66666666666666
	if snapshot.Validations.SuccessRate < expectedSuccessRate-0.01 ||
		snapshot.Validations.SuccessRate > expectedSuccessRate+0.01 {
		t.Errorf("Expected success rate ~%.2f, got %.2f", expectedSuccessRate, snapshot.Validations.SuccessRate)
	}

	// Check schema parse metrics
	if snapshot.SchemaParse.Total != 2 {
		t.Errorf("Expected 2 parses, got %d", snapshot.SchemaParse.Total)
	}

	// Check report generation metrics
	if snapshot.ReportGeneration.Total != 1 {
		t.Errorf("Expected 1 report, got %d", snapshot.ReportGeneration.Total)
	}

	// Check latency
	if snapshot.Validations.Latency.P50 == 0 {
		t.Error("Expected non-zero p50 latency")
	}

	if snapshot.Validations.Latency.P95 == 0 {
		t.Error("Expected non-zero p95 latency")
	}

	if snapshot.Validations.Latency.P99 == 0 {
		t.Error("Expected non-zero p99 latency")
	}
}

// TestReset verifies metrics reset
func TestReset(t *testing.T) {
	collector := NewMetricsCollector()

	// Record some data
	collector.RecordValidation(true, 100*time.Millisecond, 1, 0, 0)
	collector.RecordSchemaParse(true)
	collector.RecordReportGeneration(true, 50*time.Millisecond)

	// Reset
	collector.Reset()

	snapshot := collector.GetMetrics()

	if snapshot.Validations.Total != 0 {
		t.Errorf("Expected 0 validations after reset, got %d", snapshot.Validations.Total)
	}

	if snapshot.SchemaParse.Total != 0 {
		t.Errorf("Expected 0 parses after reset, got %d", snapshot.SchemaParse.Total)
	}

	if snapshot.ReportGeneration.Total != 0 {
		t.Errorf("Expected 0 reports after reset, got %d", snapshot.ReportGeneration.Total)
	}

	if len(snapshot.SeverityDistribution) != 0 {
		t.Errorf("Expected empty severity distribution after reset, got %d items", len(snapshot.SeverityDistribution))
	}
}

// TestLatencyPercentiles verifies percentile calculations
func TestLatencyPercentiles(t *testing.T) {
	collector := NewMetricsCollector()

	// Record latencies: 100ms, 200ms, 300ms, 400ms, 500ms
	latencies := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		300 * time.Millisecond,
		400 * time.Millisecond,
		500 * time.Millisecond,
	}

	for _, latency := range latencies {
		collector.RecordValidation(true, latency, 0, 0, 0)
	}

	snapshot := collector.GetMetrics()

	// For sorted [100, 200, 300, 400, 500]:
	// p50 should be 300 (index 2)
	// p95 should be 500 (index 4)
	// p99 should be 500 (index 4)

	if snapshot.Validations.Latency.P50 != 300*time.Millisecond {
		t.Errorf("Expected p50=300ms, got %v", snapshot.Validations.Latency.P50)
	}

	if snapshot.Validations.Latency.P95 != 500*time.Millisecond {
		t.Errorf("Expected p95=500ms, got %v", snapshot.Validations.Latency.P95)
	}

	if snapshot.Validations.Latency.P99 != 500*time.Millisecond {
		t.Errorf("Expected p99=500ms, got %v", snapshot.Validations.Latency.P99)
	}
}

// TestMetricsSnapshotString verifies Prometheus export format
func TestMetricsSnapshotString(t *testing.T) {
	collector := NewMetricsCollector()
	collector.RecordValidation(true, 100*time.Millisecond, 1, 0, 0)
	collector.RecordSchemaParse(true)
	collector.RecordReportGeneration(true, 50*time.Millisecond)

	snapshot := collector.GetMetrics()
	output := snapshot.String()

	// Check for Prometheus format keywords
	requiredStrings := []string{
		"HELP contract_validations_total",
		"TYPE contract_validations_total counter",
		"contract_validations_total 1",
		"contract_validations_success 1",
		"contract_validations_success_rate",
		"contract_validation_latency_p50_seconds",
		"schema_parse_total",
		"report_generation_total",
	}

	for _, required := range requiredStrings {
		if !contains(output, required) {
			t.Errorf("Expected output to contain %q", required)
		}
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
