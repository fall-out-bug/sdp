package nextstep

import "time"

// MetricType represents the type of metric event.
type MetricType int

const (
	// MetricAccepted indicates a recommendation was accepted.
	MetricAccepted MetricType = iota
	// MetricRejected indicates a recommendation was rejected.
	MetricRejected
	// MetricRefined indicates a recommendation was refined/corrected.
	MetricRefined
	// MetricAlternative indicates an alternative was selected.
	MetricAlternative
	// MetricDisplayed indicates a recommendation was displayed.
	MetricDisplayed
)

// MetricEvent represents a single metric event.
type MetricEvent struct {
	Type           MetricType     `json:"type"`
	Recommendation Recommendation `json:"recommendation"`
	Timestamp      time.Time      `json:"timestamp"`
	SessionID      string         `json:"session_id,omitempty"`
}

// MetricsStats represents aggregated metrics statistics.
type MetricsStats struct {
	TotalEvents     int                            `json:"total_events"`
	AcceptedCount   int                            `json:"accepted_count"`
	RejectedCount   int                            `json:"rejected_count"`
	RefinedCount    int                            `json:"refined_count"`
	AcceptanceRate  float64                        `json:"acceptance_rate"`
	CorrectionRate  float64                        `json:"correction_rate"`
	AvgConfidence   float64                        `json:"avg_confidence"`
	ByCategory      map[RecommendationCategory]int `json:"by_category"`
	CollectionStart time.Time                      `json:"collection_start"`
	CollectionEnd   time.Time                      `json:"collection_end"`
}

// QualityThresholds defines the minimum quality thresholds for recommendations.
type QualityThresholds struct {
	MinAcceptanceRate float64 `json:"min_acceptance_rate"`
	MaxCorrectionRate float64 `json:"max_correction_rate"`
	MinAvgConfidence  float64 `json:"min_avg_confidence"`
	MinSampleSize     int     `json:"min_sample_size"`
}

// QualityCheckResult represents the result of a quality check.
type QualityCheckResult struct {
	MeetsThresholds bool     `json:"meets_thresholds"`
	Failures        []string `json:"failures,omitempty"`
	Score           float64  `json:"score"`
}

// DefaultQualityThresholds returns the default quality thresholds.
func DefaultQualityThresholds() QualityThresholds {
	return QualityThresholds{
		MinAcceptanceRate: 0.70,
		MaxCorrectionRate: 0.20,
		MinAvgConfidence:  0.75,
		MinSampleSize:     10,
	}
}
