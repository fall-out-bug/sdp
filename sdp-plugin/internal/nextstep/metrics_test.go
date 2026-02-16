package nextstep

import (
	"testing"
	"time"
)

// TestMetricsRecord tests recording metric events.
func TestMetricsRecord(t *testing.T) {
	collector := NewMetricsCollector()

	// Record an accepted recommendation
	collector.Record(MetricEvent{
		Type: MetricAccepted,
		Recommendation: Recommendation{
			Command:    "sdp apply --ws 00-069-01",
			Category:   CategoryExecution,
			Confidence: 0.9,
		},
		Timestamp: time.Now(),
	})

	stats := collector.Stats()
	if stats.TotalEvents != 1 {
		t.Errorf("Expected 1 event, got %d", stats.TotalEvents)
	}
	if stats.AcceptedCount != 1 {
		t.Errorf("Expected 1 accepted, got %d", stats.AcceptedCount)
	}
}

// TestMetricsAcceptanceRate tests acceptance rate calculation.
func TestMetricsAcceptanceRate(t *testing.T) {
	collector := NewMetricsCollector()

	// Record 10 events: 8 accepted, 2 rejected
	for i := 0; i < 8; i++ {
		collector.Record(MetricEvent{
			Type: MetricAccepted,
			Recommendation: Recommendation{
				Command:    "sdp apply --ws test",
				Category:   CategoryExecution,
				Confidence: 0.9,
			},
			Timestamp: time.Now(),
		})
	}
	for i := 0; i < 2; i++ {
		collector.Record(MetricEvent{
			Type: MetricRejected,
			Recommendation: Recommendation{
				Command:    "sdp apply --ws test",
				Category:   CategoryExecution,
				Confidence: 0.9,
			},
			Timestamp: time.Now(),
		})
	}

	stats := collector.Stats()
	expectedRate := 0.8 // 8/10
	if stats.AcceptanceRate < expectedRate-0.01 || stats.AcceptanceRate > expectedRate+0.01 {
		t.Errorf("Expected acceptance rate ~%.2f, got %.2f", expectedRate, stats.AcceptanceRate)
	}
}

// TestMetricsCorrectionRate tests correction rate calculation.
func TestMetricsCorrectionRate(t *testing.T) {
	collector := NewMetricsCollector()

	// Record events: 6 accepted, 2 refined (corrected), 2 rejected
	for i := 0; i < 6; i++ {
		collector.Record(MetricEvent{Type: MetricAccepted, Timestamp: time.Now()})
	}
	for i := 0; i < 2; i++ {
		collector.Record(MetricEvent{Type: MetricRefined, Timestamp: time.Now()})
	}
	for i := 0; i < 2; i++ {
		collector.Record(MetricEvent{Type: MetricRejected, Timestamp: time.Now()})
	}

	stats := collector.Stats()
	// Correction rate = refined / (accepted + rejected + refined) = 2/10 = 0.20
	expectedRate := 0.20
	if stats.CorrectionRate < expectedRate-0.01 || stats.CorrectionRate > expectedRate+0.01 {
		t.Errorf("Expected correction rate ~%.2f, got %.2f", expectedRate, stats.CorrectionRate)
	}
}

// TestMetricsByCategory tests metrics grouped by category.
func TestMetricsByCategory(t *testing.T) {
	collector := NewMetricsCollector()

	// Record events by category
	collector.Record(MetricEvent{
		Type:           MetricAccepted,
		Recommendation: Recommendation{Category: CategoryExecution},
		Timestamp:      time.Now(),
	})
	collector.Record(MetricEvent{
		Type:           MetricAccepted,
		Recommendation: Recommendation{Category: CategoryExecution},
		Timestamp:      time.Now(),
	})
	collector.Record(MetricEvent{
		Type:           MetricRejected,
		Recommendation: Recommendation{Category: CategoryRecovery},
		Timestamp:      time.Now(),
	})

	stats := collector.Stats()
	if stats.ByCategory == nil {
		t.Fatal("Expected ByCategory map")
	}
	if stats.ByCategory[CategoryExecution] != 2 {
		t.Errorf("Expected 2 execution events, got %d", stats.ByCategory[CategoryExecution])
	}
	if stats.ByCategory[CategoryRecovery] != 1 {
		t.Errorf("Expected 1 recovery event, got %d", stats.ByCategory[CategoryRecovery])
	}
}

// TestMetricsThresholds tests quality threshold checking.
func TestMetricsThresholds(t *testing.T) {
	thresholds := DefaultQualityThresholds()

	if thresholds.MinAcceptanceRate <= 0 {
		t.Error("Expected positive minimum acceptance rate")
	}
	if thresholds.MaxCorrectionRate <= 0 {
		t.Error("Expected positive maximum correction rate")
	}
}

// TestMetricsQualityCheck tests quality gate evaluation.
func TestMetricsQualityCheck(t *testing.T) {
	collector := NewMetricsCollector()
	thresholds := DefaultQualityThresholds()

	// Record enough events to meet minimum sample with confidence values
	for i := 0; i < thresholds.MinSampleSize; i++ {
		collector.Record(MetricEvent{
			Type: MetricAccepted,
			Recommendation: Recommendation{
				Confidence: 0.9,
			},
			Timestamp: time.Now(),
		})
	}

	stats := collector.Stats()
	result := CheckQuality(stats, thresholds)

	if !result.MeetsThresholds {
		t.Errorf("Expected to meet thresholds with all accepted, failures: %v", result.Failures)
	}
}

// TestMetricsQualityCheckFails tests quality gate failure.
func TestMetricsQualityCheckFails(t *testing.T) {
	collector := NewMetricsCollector()
	thresholds := DefaultQualityThresholds()
	thresholds.MinAcceptanceRate = 0.9 // Require 90% acceptance

	// Record events with only 70% acceptance
	for i := 0; i < 7; i++ {
		collector.Record(MetricEvent{Type: MetricAccepted, Timestamp: time.Now()})
	}
	for i := 0; i < 3; i++ {
		collector.Record(MetricEvent{Type: MetricRejected, Timestamp: time.Now()})
	}

	stats := collector.Stats()
	result := CheckQuality(stats, thresholds)

	if result.MeetsThresholds {
		t.Error("Expected to fail thresholds with low acceptance")
	}
	if len(result.Failures) == 0 {
		t.Error("Expected failure reasons")
	}
}

// TestMetricsExport tests exporting metrics for reporting.
func TestMetricsExport(t *testing.T) {
	collector := NewMetricsCollector()

	collector.Record(MetricEvent{
		Type: MetricAccepted,
		Recommendation: Recommendation{
			Command:    "sdp test",
			Category:   CategoryExecution,
			Confidence: 0.9,
		},
		Timestamp: time.Now(),
	})

	data, err := collector.Export()
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}
	if len(data) == 0 {
		t.Error("Expected non-empty export data")
	}
}
