package monitoring

import (
	"sync"
	"time"
)

// MetricsCollector collects and tracks metrics for contract validation
type MetricsCollector struct {
	mu sync.RWMutex

	validationsTotal      int64
	validationsSuccess    int64
	validationsFailed     int64
	validationsLatency    []time.Duration
	validationsBySeverity map[string]int64

	schemaParseTotal   int64
	schemaParseSuccess int64
	schemaParseFailed  int64

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

	if len(m.validationsLatency) >= 1000 {
		m.validationsLatency = m.validationsLatency[1:]
	}
	m.validationsLatency = append(m.validationsLatency, duration)

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

	for k, v := range m.validationsBySeverity {
		snapshot.SeverityDistribution[k] = v
	}

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

	sorted := make([]time.Duration, len(durations))
	copy(sorted, durations)

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
