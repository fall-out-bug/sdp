package monitoring

import "fmt"

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
