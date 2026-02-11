package monitoring

import (
	"fmt"
	"sync"
	"time"
)

// MetricsCollector collects and tracks metrics for contract validation
type MetricsCollector struct {
	mu sync.RWMutex

	// Validation metrics
	validationsTotal      int64
	validationsSuccess    int64
	validationsFailed     int64
	validationsLatency    []time.Duration
	validationsBySeverity map[string]int64 // ERROR, WARNING, INFO

	// Schema parsing metrics
	schemaParseTotal   int64
	schemaParseSuccess int64
	schemaParseFailed  int64

	// Report generation metrics
	reportGenTotal   int64
	reportGenSuccess int64
	reportGenFailed  int64
	reportGenLatency []time.Duration
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		validationsBySeverity: make(map[string]int64),
		validationsLatency:    make([]time.Duration, 0, 1000),
		reportGenLatency:      make([]time.Duration, 0, 100),
	}
}

// RecordValidation records a contract validation operation
func (m *MetricsCollector) RecordValidation(success bool, duration time.Duration, errorCount, warningCount, infoCount int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.validationsTotal++
	if success {
		m.validationsSuccess++
	} else {
		m.validationsFailed++
	}

	// Keep last 1000 latency samples
	if len(m.validationsLatency) >= 1000 {
		m.validationsLatency = m.validationsLatency[1:]
	}
	m.validationsLatency = append(m.validationsLatency, duration)

	// Track severity distribution
	m.validationsBySeverity["ERROR"] += int64(errorCount)
	m.validationsBySeverity["WARNING"] += int64(warningCount)
	m.validationsBySeverity["INFO"] += int64(infoCount)
}

// RecordSchemaParse records a schema parsing operation
func (m *MetricsCollector) RecordSchemaParse(success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.schemaParseTotal++
	if success {
		m.schemaParseSuccess++
	} else {
		m.schemaParseFailed++
	}
}

// RecordReportGeneration records a report generation operation
func (m *MetricsCollector) RecordReportGeneration(success bool, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.reportGenTotal++
	if success {
		m.reportGenSuccess++
	} else {
		m.reportGenFailed++
	}

	// Keep last 100 latency samples
	if len(m.reportGenLatency) >= 100 {
		m.reportGenLatency = m.reportGenLatency[1:]
	}
	m.reportGenLatency = append(m.reportGenLatency, duration)
}

// GetMetrics returns a snapshot of current metrics
func (m *MetricsCollector) GetMetrics() *MetricsSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	snapshot := &MetricsSnapshot{
		Validations: ValidationMetrics{
			Total:       m.validationsTotal,
			Success:     m.validationsSuccess,
			Failed:      m.validationsFailed,
			SuccessRate: calculateSuccessRate(m.validationsSuccess, m.validationsTotal),
		},
		SchemaParse: SchemaParseMetrics{
			Total:       m.schemaParseTotal,
			Success:     m.schemaParseSuccess,
			Failed:      m.schemaParseFailed,
			SuccessRate: calculateSuccessRate(m.schemaParseSuccess, m.schemaParseTotal),
		},
		ReportGeneration: ReportGenerationMetrics{
			Total:       m.reportGenTotal,
			Success:     m.reportGenSuccess,
			Failed:      m.reportGenFailed,
			SuccessRate: calculateSuccessRate(m.reportGenSuccess, m.reportGenTotal),
		},
		SeverityDistribution: make(map[string]int64),
	}

	// Copy severity distribution
	for k, v := range m.validationsBySeverity {
		snapshot.SeverityDistribution[k] = v
	}

	// Calculate latency percentiles
	if len(m.validationsLatency) > 0 {
		snapshot.Validations.Latency = calculatePercentiles(m.validationsLatency)
	}

	if len(m.reportGenLatency) > 0 {
		snapshot.ReportGeneration.Latency = calculatePercentiles(m.reportGenLatency)
	}

	return snapshot
}

// Reset resets all metrics
func (m *MetricsCollector) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.validationsTotal = 0
	m.validationsSuccess = 0
	m.validationsFailed = 0
	m.validationsLatency = make([]time.Duration, 0, 1000)
	m.validationsBySeverity = make(map[string]int64)

	m.schemaParseTotal = 0
	m.schemaParseSuccess = 0
	m.schemaParseFailed = 0

	m.reportGenTotal = 0
	m.reportGenSuccess = 0
	m.reportGenFailed = 0
	m.reportGenLatency = make([]time.Duration, 0, 100)
}

// Helper functions

func calculateSuccessRate(success, total int64) float64 {
	if total == 0 {
		return 0
	}
	return float64(success) / float64(total) * 100
}

func calculatePercentiles(durations []time.Duration) LatencyMetrics {
	if len(durations) == 0 {
		return LatencyMetrics{}
	}

	// Sort durations
	sorted := make([]time.Duration, len(durations))
	copy(sorted, durations)

	// Simple bubble sort (good enough for small datasets)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	p50 := sorted[len(sorted)*50/100]
	p95 := sorted[len(sorted)*95/100]
	p99 := sorted[len(sorted)*99/100]

	return LatencyMetrics{
		P50: p50,
		P95: p95,
		P99: p99,
	}
}

// MetricsSnapshot represents a snapshot of current metrics
type MetricsSnapshot struct {
	Validations        ValidationMetrics     `json:"validations"`
	SchemaParse        SchemaParseMetrics    `json:"schema_parse"`
	ReportGeneration   ReportGenerationMetrics `json:"report_generation"`
	SeverityDistribution map[string]int64    `json:"severity_distribution"`
}

// ValidationMetrics tracks contract validation metrics
type ValidationMetrics struct {
	Total       int64         `json:"total"`
	Success     int64         `json:"success"`
	Failed      int64         `json:"failed"`
	SuccessRate float64       `json:"success_rate_percent"`
	Latency     LatencyMetrics `json:"latency,omitempty"`
}

// SchemaParseMetrics tracks schema parsing metrics
type SchemaParseMetrics struct {
	Total       int64   `json:"total"`
	Success     int64   `json:"success"`
	Failed      int64   `json:"failed"`
	SuccessRate float64 `json:"success_rate_percent"`
}

// ReportGenerationMetrics tracks report generation metrics
type ReportGenerationMetrics struct {
	Total       int64         `json:"total"`
	Success     int64         `json:"success"`
	Failed      int64         `json:"failed"`
	SuccessRate float64       `json:"success_rate_percent"`
	Latency     LatencyMetrics `json:"latency,omitempty"`
}

// LatencyMetrics tracks latency percentiles
type LatencyMetrics struct {
	P50 time.Duration `json:"p50"`
	P95 time.Duration `json:"p95"`
	P99 time.Duration `json:"p99"`
}

// String returns a string representation of metrics (for Prometheus export)
func (m *MetricsSnapshot) String() string {
	return fmt.Sprintf(
		`# HELP contract_validations_total Total number of contract validations
# TYPE contract_validations_total counter
contract_validations_total %d

# HELP contract_validations_success Total number of successful validations
# TYPE contract_validations_success counter
contract_validations_success %d

# HELP contract_validations_failed Total number of failed validations
# TYPE contract_validations_failed counter
contract_validations_failed %d

# HELP contract_validations_success_rate Success rate of validations
# TYPE contract_validations_success_rate gauge
contract_validations_success_rate %.2f

# HELP contract_validation_latency_seconds Validation latency percentiles
# TYPE contract_validation_latency_seconds gauge
contract_validation_latency_p50_seconds %.3f
contract_validation_latency_p95_seconds %.3f
contract_validation_latency_p99_seconds %.3f

# HELP schema_parse_total Total number of schema parses
# TYPE schema_parse_total counter
schema_parse_total %d

# HELP schema_parse_success Total number of successful schema parses
# TYPE schema_parse_success counter
schema_parse_success %d

# HELP schema_parse_failed Total number of failed schema parses
# TYPE schema_parse_failed counter
schema_parse_failed %d

# HELP report_generation_total Total number of report generations
# TYPE report_generation_total counter
report_generation_total %d

# HELP report_generation_success Total number of successful report generations
# TYPE report_generation_success counter
report_generation_success %d

# HELP report_generation_failed Total number of failed report generations
# TYPE report_generation_failed counter
report_generation_failed %d
`,
		m.Validations.Total,
		m.Validations.Success,
		m.Validations.Failed,
		m.Validations.SuccessRate,
		float64(m.Validations.Latency.P50.Seconds()),
		float64(m.Validations.Latency.P95.Seconds()),
		float64(m.Validations.Latency.P99.Seconds()),
		m.SchemaParse.Total,
		m.SchemaParse.Success,
		m.SchemaParse.Failed,
		m.ReportGeneration.Total,
		m.ReportGeneration.Success,
		m.ReportGeneration.Failed,
	)
}
