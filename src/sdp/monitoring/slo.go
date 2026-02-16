package monitoring

import (
	"fmt"
	"time"
)

// SLOs (Service Level Objectives) for contract validation
const (
	// ValidationLatencyTarget: 95% of validations should complete within 5 seconds
	ValidationLatencyTarget = 5 * time.Second

	// ValidationAvailabilityTarget: 99.9% availability (0.1% downtime allowed)
	ValidationAvailabilityTarget = 99.9

	// ValidationAccuracyTarget: 99% accuracy (no false positives/negatives)
	ValidationAccuracyTarget = 99.0

	// SchemaParseSuccessTarget: 95% of schemas should parse successfully
	SchemaParseSuccessTarget = 95.0
)

// SLOReport represents an SLO compliance report
type SLOReport struct {
	ValidationLatency    SLOStatus `json:"validation_latency"`
	ValidationAvailability SLOStatus `json:"validation_availability"`
	ValidationAccuracy   SLOStatus `json:"validation_accuracy"`
	SchemaParseSuccess   SLOStatus `json:"schema_parse_success"`
	OverallCompliance    float64   `json:"overall_compliance"`
}

// SLOStatus represents SLO compliance status
type SLOStatus struct {
	Target      float64       `json:"target"`
	Actual      float64       `json:"actual"`
	Compliant   bool          `json:"compliant"`
	Description string        `json:"description"`
}

// EvaluateSLOs evaluates SLO compliance against metrics
func EvaluateSLOs(snapshot *MetricsSnapshot) *SLOReport {
	report := &SLOReport{}

	// Evaluate validation latency (p95)
	report.ValidationLatency = SLOStatus{
		Target:      float64(ValidationLatencyTarget.Seconds()),
		Actual:      float64(snapshot.Validations.Latency.P95.Seconds()),
		Description: "95th percentile validation latency (seconds)",
	}
	report.ValidationLatency.Compliant = report.ValidationLatency.Actual <= report.ValidationLatency.Target

	// Evaluate validation availability (success rate)
	report.ValidationAvailability = SLOStatus{
		Target:      ValidationAvailabilityTarget,
		Actual:      snapshot.Validations.SuccessRate,
		Description: "Validation success rate (percent)",
	}
	report.ValidationAvailability.Compliant = report.ValidationAvailability.Actual >= report.ValidationAvailability.Target

	// Evaluate validation accuracy (assume 100% if no failures)
	accuracy := 100.0
	if snapshot.Validations.Total > 0 {
		accuracy = snapshot.Validations.SuccessRate
	}
	report.ValidationAccuracy = SLOStatus{
		Target:      ValidationAccuracyTarget,
		Actual:      accuracy,
		Description: "Validation accuracy (percent)",
	}
	report.ValidationAccuracy.Compliant = report.ValidationAccuracy.Actual >= report.ValidationAccuracy.Target

	// Evaluate schema parse success rate
	report.SchemaParseSuccess = SLOStatus{
		Target:      SchemaParseSuccessTarget,
		Actual:      snapshot.SchemaParse.SuccessRate,
		Description: "Schema parse success rate (percent)",
	}
	report.SchemaParseSuccess.Compliant = report.SchemaParseSuccess.Actual >= report.SchemaParseSuccess.Target

	// Calculate overall compliance
	compliantCount := 0
	if report.ValidationLatency.Compliant {
		compliantCount++
	}
	if report.ValidationAvailability.Compliant {
		compliantCount++
	}
	if report.ValidationAccuracy.Compliant {
		compliantCount++
	}
	if report.SchemaParseSuccess.Compliant {
		compliantCount++
	}

	report.OverallCompliance = float64(compliantCount) / 4.0 * 100.0

	return report
}

// String returns a string representation of SLO report
func (r *SLOReport) String() string {
	status := "✅ PASS"
	if r.OverallCompliance < 100 {
		status = "⚠️ PARTIAL"
	}
	if r.OverallCompliance < 75 {
		status = "❌ FAIL"
	}

	return fmt.Sprintf(
		`# SLO Compliance Report
**Overall Status: %s** (%.1f%% compliant)

## Validation Latency (p95)
- Target: %.3fs
- Actual: %.3fs
- Status: %s

## Validation Availability
- Target: %.2f%%
- Actual: %.2f%%
- Status: %s

## Validation Accuracy
- Target: %.2f%%
- Actual: %.2f%%
- Status: %s

## Schema Parse Success
- Target: %.2f%%
- Actual: %.2f%%
- Status: %s
`,
		status,
		r.OverallCompliance,
		r.ValidationLatency.Target,
		r.ValidationLatency.Actual,
		getStatusIcon(r.ValidationLatency.Compliant),
		r.ValidationAvailability.Target,
		r.ValidationAvailability.Actual,
		getStatusIcon(r.ValidationAvailability.Compliant),
		r.ValidationAccuracy.Target,
		r.ValidationAccuracy.Actual,
		getStatusIcon(r.ValidationAccuracy.Compliant),
		r.SchemaParseSuccess.Target,
		r.SchemaParseSuccess.Actual,
		getStatusIcon(r.SchemaParseSuccess.Compliant),
	)
}

func getStatusIcon(compliant bool) string {
	if compliant {
		return "✅ PASS"
	}
	return "❌ FAIL"
}

// HealthCheck checks if the system is healthy
func HealthCheck(snapshot *MetricsSnapshot) error {
	// Basic health check: at least one validation in the last 5 minutes
	if snapshot.Validations.Total == 0 {
		return fmt.Errorf("no validation metrics available - system may be down")
	}

	// Check if success rate is above minimum threshold (90%)
	if snapshot.Validations.SuccessRate < 90.0 {
		return fmt.Errorf("validation success rate below minimum: %.2f%% < 90%%", snapshot.Validations.SuccessRate)
	}

	return nil
}
