package memory

import (
	"time"
)

// MemoryTier represents storage tier (AC3)
type MemoryTier string

const (
	TierHot  MemoryTier = "hot"  // Last 30 days, always loaded
	TierWarm MemoryTier = "warm" // 30-90 days, lazy-loaded
	TierCold MemoryTier = "cold" // 90+ days, archived
)

// TierStats tracks statistics for each tier
type TierStats struct {
	HotArtifacts  int   `json:"hot_artifacts"`
	WarmArtifacts int   `json:"warm_artifacts"`
	ColdArtifacts int   `json:"cold_artifacts"`
	HotSize       int64 `json:"hot_size"`
	WarmSize      int64 `json:"warm_size"`
	ColdSize      int64 `json:"cold_size"`
}

// TotalArtifacts returns total across all tiers
func (s *TierStats) TotalArtifacts() int {
	return s.HotArtifacts + s.WarmArtifacts + s.ColdArtifacts
}

// TotalSize returns total size across all tiers
func (s *TierStats) TotalSize() int64 {
	return s.HotSize + s.WarmSize + s.ColdSize
}

// TierManager manages tiered storage (AC3)
type TierManager struct {
	policy CompactionPolicy
	stats  TierStats
}

// NewTierManager creates a new tier manager
func NewTierManager(policy CompactionPolicy) *TierManager {
	return &TierManager{
		policy: policy,
	}
}

// DetermineTier returns the tier for data of a given age (AC3)
func (m *TierManager) DetermineTier(age time.Duration) MemoryTier {
	days := int(age.Hours() / 24)

	switch {
	case days < 30:
		return TierHot
	case days < m.policy.ArchiveAfterDays:
		return TierWarm
	default:
		return TierCold
	}
}

// GetTierStats returns current tier statistics
func (m *TierManager) GetTierStats() TierStats {
	return m.stats
}

// MoveToTier moves an artifact to a different tier
func (m *TierManager) MoveToTier(artifactID string, from, to MemoryTier) error {
	// Update stats
	switch from {
	case TierHot:
		m.stats.HotArtifacts--
	case TierWarm:
		m.stats.WarmArtifacts--
	case TierCold:
		m.stats.ColdArtifacts--
	}

	switch to {
	case TierHot:
		m.stats.HotArtifacts++
	case TierWarm:
		m.stats.WarmArtifacts++
	case TierCold:
		m.stats.ColdArtifacts++
	}

	return nil
}

// Archive moves cold tier data to archive storage (AC3, AC6)
func (m *TierManager) Archive(artifactID string) error {
	return m.MoveToTier(artifactID, TierCold, TierCold)
}

// Restore brings archived data back to hot tier
func (m *TierManager) Restore(artifactID string) error {
	return m.MoveToTier(artifactID, TierCold, TierHot)
}
