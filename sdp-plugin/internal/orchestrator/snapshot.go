package orchestrator

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Snapshot represents a point-in-time state (AC1, AC2)
type Snapshot struct {
	ID          string          `json:"id"`
	FeatureID   string          `json:"feature_id"`
	Timestamp   time.Time       `json:"timestamp"`
	CompletedWS []string        `json:"completed_ws"`
	PendingWS   []string        `json:"pending_ws"`
	InProgress  *WorkInProgress `json:"in_progress,omitempty"`
	Metrics     SnapshotMetrics `json:"metrics"`
	ParentID    string          `json:"parent_id,omitempty"`
	Trigger     string          `json:"trigger,omitempty"`
}

// WorkInProgress represents work in progress state (AC2)
type WorkInProgress struct {
	WSID             string    `json:"ws_id"`
	Stage            string    `json:"stage"`
	StartedAt        time.Time `json:"started_at"`
	PartialArtifacts []string  `json:"partial_artifacts,omitempty"`
}

// SnapshotMetrics tracks execution metrics
type SnapshotMetrics struct {
	Duration    time.Duration `json:"duration"`
	Coverage    float64       `json:"coverage"`
	TestsPassed int           `json:"tests_passed"`
	TestsTotal  int           `json:"tests_total"`
}

// SnapshotDiff represents differences between snapshots (AC4)
type SnapshotDiff struct {
	FromID           string   `json:"from_id"`
	ToID             string   `json:"to_id"`
	AddedCompleted   []string `json:"added_completed"`
	RemovedCompleted []string `json:"removed_completed"`
	AddedPending     []string `json:"added_pending"`
	RemovedPending   []string `json:"removed_pending"`
}

// SnapshotManager manages snapshots (AC1, AC3, AC4, AC5)
type SnapshotManager struct {
	snapshotsDir         string
	mu                   sync.RWMutex
	completionsSinceSnap int
	autoSnapshotInterval int
	currentState         map[string]*Snapshot // featureID -> latest state
}

// NewSnapshotManager creates a new snapshot manager
func NewSnapshotManager(snapshotsDir string) *SnapshotManager {
	return &SnapshotManager{
		snapshotsDir:         snapshotsDir,
		autoSnapshotInterval: 5,
		currentState:         make(map[string]*Snapshot),
	}
}

// SetAutoSnapshotInterval sets the auto-snapshot interval
func (m *SnapshotManager) SetAutoSnapshotInterval(interval int) {
	m.autoSnapshotInterval = interval
}

// Create creates a new snapshot (AC1)
func (m *SnapshotManager) Create(featureID string, completed, pending []string) *Snapshot {
	m.mu.Lock()
	defer m.mu.Unlock()

	snap := &Snapshot{
		ID:          generateSnapshotID(featureID),
		FeatureID:   featureID,
		Timestamp:   time.Now(),
		CompletedWS: completed,
		PendingWS:   pending,
	}

	if prev, ok := m.currentState[featureID]; ok {
		snap.ParentID = prev.ID
	}

	m.currentState[featureID] = snap
	return snap
}

// Save persists a snapshot to disk
func (m *SnapshotManager) Save(snap *Snapshot) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err := os.MkdirAll(m.snapshotsDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}

	path := filepath.Join(m.snapshotsDir, snap.ID+".json")
	return os.WriteFile(path, data, 0644)
}

// Load retrieves a snapshot from disk (AC3)
func (m *SnapshotManager) Load(id string) (*Snapshot, error) {
	path := filepath.Join(m.snapshotsDir, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, err
	}

	return &snap, nil
}

// List returns all snapshots for a feature
func (m *SnapshotManager) List(featureID string) ([]*Snapshot, error) {
	var snaps []*Snapshot

	entries, err := os.ReadDir(m.snapshotsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return snaps, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		snap, err := m.Load(entry.Name()[:len(entry.Name())-5])
		if err != nil {
			continue
		}

		if snap.FeatureID == featureID {
			snaps = append(snaps, snap)
		}
	}

	return snaps, nil
}

// Rollback restores state to a previous snapshot (AC3)
func (m *SnapshotManager) Rollback(snapshotID string) (*Snapshot, error) {
	snap, err := m.Load(snapshotID)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	m.currentState[snap.FeatureID] = snap
	m.mu.Unlock()

	return snap, nil
}

// CreatePreRiskSnapshot creates a snapshot before risky operations (AC5)
func (m *SnapshotManager) CreatePreRiskSnapshot(featureID, operation string, completed, pending []string) *Snapshot {
	snap := m.Create(featureID, completed, pending)
	snap.Trigger = "pre_risk:" + operation
	m.Save(snap)
	return snap
}

// generateSnapshotID generates a unique snapshot ID
func generateSnapshotID(featureID string) string {
	ts := time.Now().UnixNano()
	data := featureID + time.Now().String()
	hash := sha256.Sum256([]byte(data))
	return featureID + "-" + hex.EncodeToString(hash[:8]) + "-" + time.Unix(0, ts).Format("20060102-150405")
}
