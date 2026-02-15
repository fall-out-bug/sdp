package monitoring

import (
	"time"
)

// MetricsSnapshot represents a snapshot of current metrics
type MetricsSnapshot struct {
	Validations          ValidationMetrics       `json:"validations"`
	SchemaParse          SchemaParseMetrics      `json:"schema_parse"`
	ReportGeneration     ReportGenerationMetrics `json:"report_generation"`
	SeverityDistribution map[string]int64        `json:"severity_distribution"`
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
