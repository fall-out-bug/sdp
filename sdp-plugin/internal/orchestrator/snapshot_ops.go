package orchestrator

import (
	"time"
)

// Diff compares two snapshots (AC4)
func (m *SnapshotManager) Diff(fromID, toID string) (*SnapshotDiff, error) {
	from, err := m.Load(fromID)
	if err != nil {
		return nil, err
	}

	to, err := m.Load(toID)
	if err != nil {
		return nil, err
	}

	diff := &SnapshotDiff{
		FromID: fromID,
		ToID:   toID,
	}

	// Find added/removed completed
	fromSet := make(map[string]bool)
	for _, ws := range from.CompletedWS {
		fromSet[ws] = true
	}
	for _, ws := range to.CompletedWS {
		if !fromSet[ws] {
			diff.AddedCompleted = append(diff.AddedCompleted, ws)
		}
	}

	toSet := make(map[string]bool)
	for _, ws := range to.CompletedWS {
		toSet[ws] = true
	}
	for _, ws := range from.CompletedWS {
		if !toSet[ws] {
			diff.RemovedCompleted = append(diff.RemovedCompleted, ws)
		}
	}

	return diff, nil
}

// RecordCompletion records a workstream completion and triggers auto-snapshot (AC5)
func (m *SnapshotManager) RecordCompletion(featureID, wsID string) *Snapshot {
	m.mu.Lock()

	m.completionsSinceSnap++

	// Get or create current state
	state, ok := m.currentState[featureID]
	if !ok {
		state = &Snapshot{
			FeatureID:   featureID,
			CompletedWS: []string{},
			PendingWS:   []string{},
		}
		m.currentState[featureID] = state
	}

	// Add the completed workstream
	state.CompletedWS = append(state.CompletedWS, wsID)

	var snap *Snapshot
	if m.completionsSinceSnap >= m.autoSnapshotInterval {
		m.completionsSinceSnap = 0

		snap = &Snapshot{
			ID:          generateSnapshotID(featureID),
			FeatureID:   featureID,
			Timestamp:   time.Now(),
			CompletedWS: append([]string{}, state.CompletedWS...),
			PendingWS:   append([]string{}, state.PendingWS...),
			Trigger:     "auto",
		}
		if state.ID != "" {
			snap.ParentID = state.ID
		}
		m.currentState[featureID] = snap
	}

	m.mu.Unlock()

	if snap != nil {
		m.Save(snap)
	}

	return snap
}
