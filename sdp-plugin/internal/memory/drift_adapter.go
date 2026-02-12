package memory

import (
	"encoding/json"
	"time"

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

	// Serialize report to JSON for content
	dataBytes, _ := json.Marshal(report)

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

// buildDriftContent creates searchable content from a DriftReport
func (a *DriftAdapter) buildDriftContent(report *drift.DriftReport) string {
	content := "Drift report for " + report.WorkstreamID + "\n"
	content += "Verdict: " + report.Verdict + "\n"
	content += "Timestamp: " + report.Timestamp.Format(time.RFC3339) + "\n"

	for _, issue := range report.Issues {
		content += "Issue: " + issue.File + " " + string(issue.Status) + "\n"
		content += "Expected: " + issue.Expected + "\n"
		if issue.Actual != "" {
			content += "Actual: " + issue.Actual + "\n"
		}
		if issue.Recommendation != "" {
			content += "Recommendation: " + issue.Recommendation + "\n"
		}
	}

	return content
}

// buildEnhancedDriftContent creates searchable content from an EnhancedDriftReport
func (a *DriftAdapter) buildEnhancedDriftContent(report *drift.EnhancedDriftReport, dataBytes []byte) string {
	content := "Enhanced drift report for " + report.WorkstreamID + "\n"
	content += "Verdict: " + report.Verdict + "\n"
	content += "Timestamp: " + report.Timestamp.Format(time.RFC3339) + "\n"

	for _, dt := range report.DriftTypes {
		content += "Type: " + string(dt.Type) + " (" + string(dt.Severity) + ")\n"
		for _, issue := range dt.Issues {
			content += "Issue: " + issue.File
			if issue.Line > 0 {
				content += ":" + string(rune(issue.Line))
			}
			content += " " + issue.Message + "\n"
		}
		for _, s := range dt.Suggestions {
			content += "Suggestion: " + s + "\n"
		}
	}

	// Include full JSON for detailed search
	content += " " + string(dataBytes)

	return content
}

// buildDriftTags builds tags from drift report
func (a *DriftAdapter) buildDriftTags(report *drift.EnhancedDriftReport) []string {
	tags := []string{"drift", report.Verdict}

	for _, dt := range report.DriftTypes {
		tags = append(tags, string(dt.Type))
		if dt.Severity == drift.SeverityError {
			tags = append(tags, "error")
		} else if dt.Severity == drift.SeverityWarning {
			tags = append(tags, "warning")
		}
	}

	return tags
}

// hashDriftReport generates a hash for a drift report
func hashDriftReport(report *drift.DriftReport) string {
	// Simple hash based on content
	data, _ := json.Marshal(report)
	return simpleHash(string(data))
}

// hashEnhancedDriftReport generates a hash for an enhanced drift report
func hashEnhancedDriftReport(report *drift.EnhancedDriftReport) string {
	data, _ := json.Marshal(report)
	return simpleHash(string(data))
}

// simpleHash generates a simple hash from a string
func simpleHash(s string) string {
	// Simple FNV-1a style hash
	h := uint32(2166136261)
	for _, c := range s {
		h ^= uint32(c)
		h *= 16777619
	}
	return "h" + intToStr(int(h))
}

// intToStr converts int to string without importing strconv
func intToStr(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}
