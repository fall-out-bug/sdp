package orchestrator

import (
	"testing"
	"time"
)

func TestSnapshot_Create(t *testing.T) {
	mgr := NewSnapshotManager("/tmp/test-snapshots")

	snap := mgr.Create("F051", []string{"00-051-01", "00-051-02"}, []string{"00-051-03"})

	if snap.ID == "" {
		t.Error("Expected snapshot ID to be generated")
	}
	if snap.FeatureID != "F051" {
		t.Errorf("Expected F051, got %s", snap.FeatureID)
	}
	if len(snap.CompletedWS) != 2 {
		t.Errorf("Expected 2 completed WS, got %d", len(snap.CompletedWS))
	}
	if len(snap.PendingWS) != 1 {
		t.Errorf("Expected 1 pending WS, got %d", len(snap.PendingWS))
	}
}

func TestSnapshot_WithInProgress(t *testing.T) {
	mgr := NewSnapshotManager("/tmp/test-snapshots")

	wip := &WorkInProgress{
		WSID:      "00-051-03",
		Stage:     "implementing",
		StartedAt: time.Now(),
	}

	snap := mgr.Create("F051", []string{"00-051-01"}, []string{"00-051-03"})
	snap.InProgress = wip

	if snap.InProgress.WSID != "00-051-03" {
		t.Errorf("Expected in-progress WSID, got %s", snap.InProgress.WSID)
	}
}

func TestSnapshotManager_SaveAndLoad(t *testing.T) {
	mgr := NewSnapshotManager(t.TempDir())

	snap := mgr.Create("F051", []string{"00-051-01"}, []string{"00-051-02", "00-051-03"})
	snap.Metrics = SnapshotMetrics{
		Duration:    30 * time.Minute,
		Coverage:    85.5,
		TestsPassed: 10,
		TestsTotal:  10,
	}

	// Save
	if err := mgr.Save(snap); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load
	loaded, err := mgr.Load(snap.ID)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.FeatureID != snap.FeatureID {
		t.Errorf("Expected %s, got %s", snap.FeatureID, loaded.FeatureID)
	}
	if loaded.Metrics.Coverage != snap.Metrics.Coverage {
		t.Errorf("Expected coverage %f, got %f", snap.Metrics.Coverage, loaded.Metrics.Coverage)
	}
}

func TestSnapshotManager_List(t *testing.T) {
	mgr := NewSnapshotManager(t.TempDir())

	// Create multiple snapshots
	snap1 := mgr.Create("F051", []string{"00-051-01"}, []string{"00-051-02"})
	snap2 := mgr.Create("F051", []string{"00-051-01", "00-051-02"}, []string{"00-051-03"})
	snap3 := mgr.Create("F052", []string{"00-052-01"}, []string{})

	mgr.Save(snap1)
	mgr.Save(snap2)
	mgr.Save(snap3)

	// List all for F051
	snaps, err := mgr.List("F051")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(snaps) != 2 {
		t.Errorf("Expected 2 snapshots for F051, got %d", len(snaps))
	}
}

func TestSnapshotManager_Rollback(t *testing.T) {
	mgr := NewSnapshotManager(t.TempDir())

	// Create initial snapshot
	snap1 := mgr.Create("F051", []string{"00-051-01"}, []string{"00-051-02", "00-051-03"})
	mgr.Save(snap1)

	// Create later snapshot
	snap2 := mgr.Create("F051", []string{"00-051-01", "00-051-02"}, []string{"00-051-03"})
	mgr.Save(snap2)

	// Rollback to first snapshot
	state, err := mgr.Rollback(snap1.ID)
	if err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}

	if len(state.CompletedWS) != 1 {
		t.Errorf("Expected 1 completed WS after rollback, got %d", len(state.CompletedWS))
	}
	if state.CompletedWS[0] != "00-051-01" {
		t.Errorf("Expected 00-051-01, got %s", state.CompletedWS[0])
	}
}

func TestSnapshotManager_Diff(t *testing.T) {
	mgr := NewSnapshotManager(t.TempDir())

	snap1 := mgr.Create("F051", []string{"00-051-01"}, []string{"00-051-02", "00-051-03"})
	snap2 := mgr.Create("F051", []string{"00-051-01", "00-051-02"}, []string{"00-051-03"})

	mgr.Save(snap1)
	mgr.Save(snap2)

	diff, err := mgr.Diff(snap1.ID, snap2.ID)
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}

	if len(diff.AddedCompleted) != 1 {
		t.Errorf("Expected 1 added completed, got %d", len(diff.AddedCompleted))
	}
	if diff.AddedCompleted[0] != "00-051-02" {
		t.Errorf("Expected 00-051-02 added, got %v", diff.AddedCompleted)
	}
}

func TestSnapshotManager_AutoSnapshot(t *testing.T) {
	tmpDir := t.TempDir()
	mgr := NewSnapshotManager(tmpDir)
	mgr.SetAutoSnapshotInterval(2) // Auto-snapshot every 2 completed WS

	// First completion - no auto-snapshot
	snap1 := mgr.RecordCompletion("F051", "00-051-01")
	if snap1 != nil {
		t.Error("Expected no auto-snapshot after 1 completion")
	}

	// Verify no snapshots on disk
	snaps, err := mgr.List("F051")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(snaps) != 0 {
		t.Errorf("Expected 0 snapshots, got %d", len(snaps))
	}

	// Second completion - auto-snapshot should trigger
	snap2 := mgr.RecordCompletion("F051", "00-051-02")
	if snap2 == nil {
		t.Fatal("Expected auto-snapshot after 2 completions")
	}
	if snap2.Trigger != "auto" {
		t.Errorf("Expected trigger 'auto', got %s", snap2.Trigger)
	}

	// Verify snapshot was saved to disk
	snaps, _ = mgr.List("F051")
	if len(snaps) != 1 {
		t.Errorf("Expected 1 snapshot on disk, got %d", len(snaps))
	}
}

func TestSnapshotManager_CreatePreRiskSnapshot(t *testing.T) {
	mgr := NewSnapshotManager(t.TempDir())

	snap := mgr.CreatePreRiskSnapshot("F051", "merge", []string{"00-051-01"}, []string{"00-051-02"})

	if snap == nil {
		t.Fatal("Expected pre-risk snapshot to be created")
	}
	if snap.Trigger != "pre_risk:merge" {
		t.Errorf("Expected trigger 'pre_risk:merge', got %s", snap.Trigger)
	}

	// Verify it was saved
	loaded, err := mgr.Load(snap.ID)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.Trigger != snap.Trigger {
		t.Error("Trigger mismatch in loaded snapshot")
	}
}
