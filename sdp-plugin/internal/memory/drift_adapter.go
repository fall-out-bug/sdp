package memory

import (
	"encoding/json"

	"github.com/fall-out-bug/sdp/internal/drift"
)

// DriftAdapter provides integration between drift detection and memory store
type DriftAdapter struct {
	store *Store
}

// NewDriftAdapter creates a new adapter for drift-memory integration
func NewDriftAdapter(store *Store) *DriftAdapter {
	return &DriftAdapter{store: store}
}

// SaveDriftReport stores a drift report as an artifact in memory
func (a *DriftAdapter) SaveDriftReport(report *drift.DriftReport) error {
	artifact := a.reportToArtifact(report)
	return a.store.Save(artifact)
}

// SaveEnhancedDriftReport stores an enhanced drift report as an artifact
func (a *DriftAdapter) SaveEnhancedDriftReport(report *drift.EnhancedDriftReport) error {
	artifact := a.enhancedReportToArtifact(report)
	return a.store.Save(artifact)
}

// GetDriftReports retrieves all drift reports for a workstream
func (a *DriftAdapter) GetDriftReports(workstreamID string) ([]*Artifact, error) {
	// Get all artifacts and filter by type and workstream
	all, err := a.store.ListAll()
	if err != nil {
		return nil, err
	}

	var reports []*Artifact
	for _, r := range all {
		if r.Type == "drift" && r.WorkstreamID == workstreamID {
			reports = append(reports, r)
		}
	}
	return reports, nil
}

// GetLatestDriftReport gets the most recent drift report for a workstream
func (a *DriftAdapter) GetLatestDriftReport(workstreamID string) (*Artifact, error) {
	reports, err := a.GetDriftReports(workstreamID)
	if err != nil {
		return nil, err
	}
	if len(reports) == 0 {
		return nil, nil
	}

	// Find most recent
	latest := reports[0]
	for _, r := range reports[1:] {
		if r.IndexedAt.After(latest.IndexedAt) {
			latest = r
		}
	}
	return latest, nil
}

// reportToArtifact converts a DriftReport to an Artifact
func (a *DriftAdapter) reportToArtifact(report *drift.DriftReport) *Artifact {
	// Generate unique ID and path (include timestamp to allow multiple reports per WS)
	ts := report.Timestamp.Format("20060102-150405")
	artifactID := "drift-" + report.WorkstreamID + "-" + ts

	// Build searchable content
	content := a.buildDriftContent(report)

	return &Artifact{
		ID:           artifactID,
		Path:         ".sdp/drift/" + report.WorkstreamID + "/" + ts + ".json",
		Type:         "drift",
		Title:        "Drift Report: " + report.WorkstreamID + " (" + report.Verdict + ")",
		Content:      content,
		FeatureID:    extractFeatureFromWSID(report.WorkstreamID),
		WorkstreamID: report.WorkstreamID,
		Tags:         []string{"drift", report.Verdict},
		FileHash:     hashDriftReport(report),
		IndexedAt:    report.Timestamp,
	}
}

// enhancedReportToArtifact converts an EnhancedDriftReport to an Artifact
func (a *DriftAdapter) enhancedReportToArtifact(report *drift.EnhancedDriftReport) *Artifact {
	// Generate unique ID and path (include timestamp to allow multiple reports per WS)
	ts := report.Timestamp.Format("20060102-150405")
	artifactID := "drift-enh-" + report.WorkstreamID + "-" + ts

	// Serialize report to JSON for content (ignore error - content is optional)
	dataBytes, _ := json.Marshal(report) //nolint:errcheck // content is optional

	// Build searchable content
	content := a.buildEnhancedDriftContent(report, dataBytes)

	return &Artifact{
		ID:           artifactID,
		Path:         ".sdp/drift/enhanced/" + report.WorkstreamID + "/" + ts + ".json",
		Type:         "drift",
		Title:        "Enhanced Drift: " + report.WorkstreamID + " (" + report.Verdict + ")",
		Content:      content,
		FeatureID:    extractFeatureFromWSID(report.WorkstreamID),
		WorkstreamID: report.WorkstreamID,
		Tags:         a.buildDriftTags(report),
		FileHash:     hashEnhancedDriftReport(report),
		IndexedAt:    report.Timestamp,
	}
}
