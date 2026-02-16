package decision

import (
	"expvar"
	"sync"
	"time"
)

// Metrics for decision logging
var (
	// Counters
	logCounter        = expvar.NewInt("decision_log_total")
	logSuccessCounter = expvar.NewInt("decision_log_success")
	logErrorCounter   = expvar.NewInt("decision_log_errors")

	batchLogCounter        = expvar.NewInt("decision_batch_log_total")
	batchLogSuccessCounter = expvar.NewInt("decision_batch_log_success")
	batchLogErrorCounter   = expvar.NewInt("decision_batch_log_errors")

	loadCounter        = expvar.NewInt("decision_load_total")
	loadSuccessCounter = expvar.NewInt("decision_load_success")
	loadErrorCounter   = expvar.NewInt("decision_load_errors")

	// Gauges
	currentFileSize = expvar.NewInt("decision_file_size_bytes")
	parseErrorCount = expvar.NewInt("decision_parse_errors")

	// Histograms (manual)
	logDurationMs      = expvar.NewInt("decision_log_duration_ms")
	batchLogDurationMs = expvar.NewInt("decision_batch_log_duration_ms")
	loadDurationMs     = expvar.NewInt("decision_load_duration_ms")
)

// MetricsRecorder tracks operation metrics
type MetricsRecorder struct {
	mu sync.Mutex
}

// RecordLog records a log operation
func (m *MetricsRecorder) RecordLog(duration time.Duration, success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	logCounter.Add(1)
	logDurationMs.Add(duration.Milliseconds())

	if success {
		logSuccessCounter.Add(1)
	} else {
		logErrorCounter.Add(1)
	}
}

// RecordBatchLog records a batch log operation
func (m *MetricsRecorder) RecordBatchLog(duration time.Duration, success bool, count int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	batchLogCounter.Add(1)
	batchLogDurationMs.Add(duration.Milliseconds())

	if success {
		batchLogSuccessCounter.Add(1)
	} else {
		batchLogErrorCounter.Add(1)
	}
}

// RecordLoad records a load operation
func (m *MetricsRecorder) RecordLoad(duration time.Duration, success bool, parseErrors int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	loadCounter.Add(1)
	loadDurationMs.Add(duration.Milliseconds())
	parseErrorCount.Add(int64(parseErrors))

	if success {
		loadSuccessCounter.Add(1)
	} else {
		loadErrorCounter.Add(1)
	}
}
