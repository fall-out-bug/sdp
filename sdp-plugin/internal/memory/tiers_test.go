package memory

import (
	"testing"
	"time"
)

func TestTierManager_DetermineTier(t *testing.T) {
	mgr := NewTierManager(DefaultCompactionPolicy())

	// Hot: recent
	age := time.Hour * 24 * 15 // 15 days
	if mgr.DetermineTier(age) != TierHot {
		t.Error("Expected hot tier for 15 days")
	}

	// Warm: middle
	age = time.Hour * 24 * 60 // 60 days
	if mgr.DetermineTier(age) != TierWarm {
		t.Error("Expected warm tier for 60 days")
	}

	// Cold: old
	age = time.Hour * 24 * 120 // 120 days
	if mgr.DetermineTier(age) != TierCold {
		t.Error("Expected cold tier for 120 days")
	}
}

func TestTierManager_GetTierStats(t *testing.T) {
	mgr := NewTierManager(DefaultCompactionPolicy())

	stats := mgr.GetTierStats()

	// Should return a valid stats structure
	if stats.HotArtifacts < 0 || stats.WarmArtifacts < 0 || stats.ColdArtifacts < 0 {
		t.Error("Stats should not be negative")
	}
}

func TestMemoryTier_String(t *testing.T) {
	tests := []struct {
		tier     MemoryTier
		expected string
	}{
		{TierHot, "hot"},
		{TierWarm, "warm"},
		{TierCold, "cold"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if string(tt.tier) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.tier))
			}
		})
	}
}

func TestTierStats_Total(t *testing.T) {
	stats := TierStats{
		HotArtifacts:  100,
		WarmArtifacts: 50,
		ColdArtifacts: 25,
		HotSize:       1024,
		WarmSize:      2048,
		ColdSize:      512,
	}

	total := stats.TotalArtifacts()
	if total != 175 {
		t.Errorf("Expected 175 total, got %d", total)
	}

	totalSize := stats.TotalSize()
	if totalSize != 3584 {
		t.Errorf("Expected 3584 bytes total, got %d", totalSize)
	}
}

func TestTierManager_MoveToTier(t *testing.T) {
	mgr := NewTierManager(DefaultCompactionPolicy())
	mgr.stats.HotArtifacts = 10

	err := mgr.MoveToTier("artifact-1", TierHot, TierWarm)
	if err != nil {
		t.Fatalf("MoveToTier failed: %v", err)
	}

	if mgr.stats.HotArtifacts != 9 {
		t.Errorf("Expected 9 hot, got %d", mgr.stats.HotArtifacts)
	}
	if mgr.stats.WarmArtifacts != 1 {
		t.Errorf("Expected 1 warm, got %d", mgr.stats.WarmArtifacts)
	}
}

func TestTierManager_Archive(t *testing.T) {
	mgr := NewTierManager(DefaultCompactionPolicy())

	err := mgr.Archive("artifact-1")
	if err != nil {
		t.Fatalf("Archive failed: %v", err)
	}
}

func TestTierManager_Restore(t *testing.T) {
	mgr := NewTierManager(DefaultCompactionPolicy())
	mgr.stats.ColdArtifacts = 5

	err := mgr.Restore("artifact-1")
	if err != nil {
		t.Fatalf("Restore failed: %v", err)
	}

	if mgr.stats.HotArtifacts != 1 {
		t.Errorf("Expected 1 hot after restore, got %d", mgr.stats.HotArtifacts)
	}
}
