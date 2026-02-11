package monitoring

import (
	"testing"
	"time"
)

// TestEvaluateSLOs_AllPassing verifies SLO evaluation when all targets met
func TestEvaluateSLOs_AllPassing(t *testing.T) {
	collector := NewMetricsCollector()

	// Record good metrics (fast, successful)
	for i := 0; i < 100; i++ {
		collector.RecordValidation(true, 100*time.Millisecond, 0, 0, 0)
		collector.RecordSchemaParse(true)
		collector.RecordReportGeneration(true, 50*time.Millisecond)
	}

	snapshot := collector.GetMetrics()
	report := EvaluateSLOs(snapshot)

	// All should be compliant
	if !report.ValidationLatency.Compliant {
		t.Error("Expected validation latency to be compliant")
	}

	if !report.ValidationAvailability.Compliant {
		t.Error("Expected validation availability to be compliant")
	}

	if !report.ValidationAccuracy.Compliant {
		t.Error("Expected validation accuracy to be compliant")
	}

	if !report.SchemaParseSuccess.Compliant {
		t.Error("Expected schema parse success to be compliant")
	}

	if report.OverallCompliance != 100.0 {
		t.Errorf("Expected 100%% compliance, got %.2f%%", report.OverallCompliance)
	}
}

// TestEvaluateSLOs_SlowLatency verifies SLO evaluation when latency is too high
func TestEvaluateSLOs_SlowLatency(t *testing.T) {
	collector := NewMetricsCollector()

	// Record some slow validations (p95 will exceed 5s target)
	for i := 0; i < 100; i++ {
		latency := 100 * time.Millisecond
		if i >= 95 {
			latency = 6 * time.Second // Slow p95
		}
		collector.RecordValidation(true, latency, 0, 0, 0)
	}

	snapshot := collector.GetMetrics()
	report := EvaluateSLOs(snapshot)

	if report.ValidationLatency.Compliant {
		t.Error("Expected validation latency to be non-compliant (p95 > 5s)")
	}

	if report.ValidationLatency.Actual <= report.ValidationLatency.Target {
		t.Errorf("Expected p95 latency > target (5s), got %.3fs", report.ValidationLatency.Actual)
	}

	// Other SLOs should still pass
	if !report.ValidationAvailability.Compliant {
		t.Error("Expected validation availability to still be compliant")
	}
}

// TestEvaluateSLOs_LowAvailability verifies SLO evaluation when success rate is low
func TestEvaluateSLOs_LowAvailability(t *testing.T) {
	collector := NewMetricsCollector()

	// Record 50% success rate (below 99.9% target)
	for i := 0; i < 100; i++ {
		success := i%2 == 0
		collector.RecordValidation(success, 100*time.Millisecond, 0, 0, 0)
	}

	snapshot := collector.GetMetrics()
	report := EvaluateSLOs(snapshot)

	if report.ValidationAvailability.Compliant {
		t.Error("Expected validation availability to be non-compliant (< 99.9%)")
	}

	if report.ValidationAvailability.Actual >= report.ValidationAvailability.Target {
		t.Errorf("Expected success rate < target (99.9%%), got %.2f%%", report.ValidationAvailability.Actual)
	}
}

// TestEvaluateSLOs_LowSchemaParseSuccess verifies SLO evaluation when schema parse fails often
func TestEvaluateSLOs_LowSchemaParseSuccess(t *testing.T) {
	collector := NewMetricsCollector()

	// Record 90% schema parse success (below 95% target)
	for i := 0; i < 100; i++ {
		success := i < 90
		collector.RecordSchemaParse(success)
	}

	snapshot := collector.GetMetrics()
	report := EvaluateSLOs(snapshot)

	if report.SchemaParseSuccess.Compliant {
		t.Error("Expected schema parse success to be non-compliant (< 95%)")
	}

	if report.SchemaParseSuccess.Actual >= report.SchemaParseSuccess.Target {
		t.Errorf("Expected schema parse success < target (95%%), got %.2f%%", report.SchemaParseSuccess.Actual)
	}
}

// TestEvaluateSLOs_PartialCompliance verifies partial compliance calculation
func TestEvaluateSLOs_PartialCompliance(t *testing.T) {
	collector := NewMetricsCollector()

	// Record metrics that will fail 2 out of 4 SLOs
	// Fast latency for most, but slow p95
	for i := 0; i < 100; i++ {
		latency := 100 * time.Millisecond
		if i >= 95 {
			latency = 6 * time.Second // Slow p95 - FAILS
		}
		success := i < 90 // Fail availability (90% < 99.9%) - FAILS
		collector.RecordValidation(success, latency, 0, 0, 0)
	}

	// Record schema parses (98% success - passes 95% target)
	for i := 0; i < 100; i++ {
		parseSuccess := i < 98
		collector.RecordSchemaParse(parseSuccess)
	}

	snapshot := collector.GetMetrics()
	report := EvaluateSLOs(snapshot)

	// Should be 25% compliant (1 out of 4 SLOs passing)
	// - Latency: FAIL (p95 > 5s)
	// - Availability: FAIL (90% < 99.9%)
	// - Accuracy: FAIL (90% < 99%)
	// - Schema parse: PASS (98% >= 95%)
	expectedCompliance := 25.0
	if report.OverallCompliance < expectedCompliance-1 || report.OverallCompliance > expectedCompliance+1 {
		t.Errorf("Expected ~%.1f%% compliance, got %.2f%%", expectedCompliance, report.OverallCompliance)
	}
}

// TestSLOReportString verifies report string format
func TestSLOReportString(t *testing.T) {
	collector := NewMetricsCollector()
	collector.RecordValidation(true, 100*time.Millisecond, 0, 0, 0)

	snapshot := collector.GetMetrics()
	report := EvaluateSLOs(snapshot)

	output := report.String()

	// Check for required sections
	requiredStrings := []string{
		"SLO Compliance Report",
		"Overall Status",
		"Validation Latency",
		"Validation Availability",
		"Validation Accuracy",
		"Schema Parse Success",
	}

	for _, required := range requiredStrings {
		if !contains(output, required) {
			t.Errorf("Expected report to contain %q", required)
		}
	}
}

// TestHealthCheck_Passing verifies health check passes
func TestHealthCheck_Passing(t *testing.T) {
	collector := NewMetricsCollector()

	// Record successful validation
	collector.RecordValidation(true, 100*time.Millisecond, 0, 0, 0)

	snapshot := collector.GetMetrics()
	err := HealthCheck(snapshot)

	if err != nil {
		t.Errorf("Expected health check to pass, got error: %v", err)
	}
}

// TestHealthCheck_NoMetrics verifies health check fails when no metrics
func TestHealthCheck_NoMetrics(t *testing.T) {
	collector := NewMetricsCollector()
	snapshot := collector.GetMetrics()

	err := HealthCheck(snapshot)

	if err == nil {
		t.Error("Expected health check to fail when no metrics available")
	}

	if !contains(err.Error(), "no validation metrics") {
		t.Errorf("Expected 'no validation metrics' error, got: %v", err)
	}
}

// TestHealthCheck_LowSuccessRate verifies health check fails when success rate is low
func TestHealthCheck_LowSuccessRate(t *testing.T) {
	collector := NewMetricsCollector()

	// Record 80% success rate (below 90% minimum)
	for i := 0; i < 100; i++ {
		success := i < 80
		collector.RecordValidation(success, 100*time.Millisecond, 0, 0, 0)
	}

	snapshot := collector.GetMetrics()
	err := HealthCheck(snapshot)

	if err == nil {
		t.Error("Expected health check to fail with low success rate")
	}

	if !contains(err.Error(), "success rate below minimum") {
		t.Errorf("Expected 'success rate below minimum' error, got: %v", err)
	}
}

// TestSLOConstants verifies SLO target constants
func TestSLOConstants(t *testing.T) {
	if ValidationLatencyTarget != 5*time.Second {
		t.Errorf("Expected ValidationLatencyTarget=5s, got %v", ValidationLatencyTarget)
	}

	if ValidationAvailabilityTarget != 99.9 {
		t.Errorf("Expected ValidationAvailabilityTarget=99.9, got %f", ValidationAvailabilityTarget)
	}

	if ValidationAccuracyTarget != 99.0 {
		t.Errorf("Expected ValidationAccuracyTarget=99.0, got %f", ValidationAccuracyTarget)
	}

	if SchemaParseSuccessTarget != 95.0 {
		t.Errorf("Expected SchemaParseSuccessTarget=95.0, got %f", SchemaParseSuccessTarget)
	}
}
