package nextstep

import (
	"encoding/json"
	"sync"
	"time"
)

// MetricsCollector collects and aggregates recommendation metrics.
type MetricsCollector struct {
	mu     sync.RWMutex
	events []MetricEvent
}

// NewMetricsCollector creates a new metrics collector.
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		events: []MetricEvent{},
	}
}

// Record records a metric event.
func (c *MetricsCollector) Record(event MetricEvent) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	c.events = append(c.events, event)
}

// Stats returns aggregated statistics.
func (c *MetricsCollector) Stats() MetricsStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stats := MetricsStats{
		ByCategory: make(map[RecommendationCategory]int),
	}

	if len(c.events) == 0 {
		return stats
	}

	stats.TotalEvents = len(c.events)
	stats.CollectionStart = c.events[0].Timestamp
	stats.CollectionEnd = c.events[len(c.events)-1].Timestamp

	var totalConfidence float64
	var confidenceCount int

	for _, event := range c.events {
		switch event.Type {
		case MetricAccepted:
			stats.AcceptedCount++
		case MetricRejected:
			stats.RejectedCount++
		case MetricRefined:
			stats.RefinedCount++
		}

		if event.Recommendation.Category != "" {
			stats.ByCategory[event.Recommendation.Category]++
		}

		if event.Recommendation.Confidence > 0 {
			totalConfidence += event.Recommendation.Confidence
			confidenceCount++
		}
	}

	decisions := stats.AcceptedCount + stats.RejectedCount + stats.RefinedCount
	if decisions > 0 {
		stats.AcceptanceRate = float64(stats.AcceptedCount) / float64(decisions)
		stats.CorrectionRate = float64(stats.RefinedCount) / float64(decisions)
	}

	if confidenceCount > 0 {
		stats.AvgConfidence = totalConfidence / float64(confidenceCount)
	}

	return stats
}

// Export exports metrics as JSON.
func (c *MetricsCollector) Export() ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return json.Marshal(c.events)
}

// Reset clears all collected metrics.
func (c *MetricsCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.events = []MetricEvent{}
}

// CheckQuality evaluates metrics against thresholds.
func CheckQuality(stats MetricsStats, thresholds QualityThresholds) QualityCheckResult {
	result := QualityCheckResult{
		MeetsThresholds: true,
		Failures:        []string{},
	}

	if stats.TotalEvents < thresholds.MinSampleSize {
		result.MeetsThresholds = false
		result.Failures = append(result.Failures, "insufficient sample size")
		return result
	}

	if stats.AcceptanceRate < thresholds.MinAcceptanceRate {
		result.MeetsThresholds = false
		result.Failures = append(result.Failures, "acceptance rate below threshold")
	}

	if stats.CorrectionRate > thresholds.MaxCorrectionRate {
		result.MeetsThresholds = false
		result.Failures = append(result.Failures, "correction rate above threshold")
	}

	if stats.AvgConfidence < thresholds.MinAvgConfidence {
		result.MeetsThresholds = false
		result.Failures = append(result.Failures, "average confidence below threshold")
	}

	result.Score = calculateScore(stats)
	return result
}

// calculateScore calculates an overall quality score.
func calculateScore(stats MetricsStats) float64 {
	if stats.TotalEvents == 0 {
		return 0
	}

	acceptanceWeight := 0.4
	correctionWeight := 0.3
	confidenceWeight := 0.3

	acceptanceScore := stats.AcceptanceRate * 100
	if acceptanceScore > 100 {
		acceptanceScore = 100
	}

	correctionScore := (1 - stats.CorrectionRate) * 100
	if correctionScore < 0 {
		correctionScore = 0
	}

	confidenceScore := stats.AvgConfidence * 100

	return (acceptanceScore * acceptanceWeight) +
		(correctionScore * correctionWeight) +
		(confidenceScore * confidenceWeight)
}
