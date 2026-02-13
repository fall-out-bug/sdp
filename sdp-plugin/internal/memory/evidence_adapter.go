package memory

import (
	"encoding/json"
	"time"

	"github.com/fall-out-bug/sdp/internal/evidence"
)

// EvidenceAdapter provides integration between evidence.jsonl and memory store
type EvidenceAdapter struct {
	store *Store
}

// NewEvidenceAdapter creates a new adapter for evidence-memory integration
func NewEvidenceAdapter(store *Store) *EvidenceAdapter {
	return &EvidenceAdapter{store: store}
}

// ImportEvents imports evidence events into memory as artifacts
func (a *EvidenceAdapter) ImportEvents(events []evidence.Event) (int, error) {
	imported := 0
	for _, ev := range events {
		artifact := a.eventToArtifact(ev)
		if artifact == nil {
			continue
		}
		if err := a.store.Save(artifact); err != nil {
			continue
		}
		imported++
	}
	return imported, nil
}

// eventToArtifact converts an evidence.Event to a memory.Artifact
func (a *EvidenceAdapter) eventToArtifact(ev evidence.Event) *Artifact {
	if ev.ID == "" || ev.WSID == "" {
		return nil
	}

	// Generate artifact ID from event ID
	artifactID := "ev-" + ev.ID

	// Create content from event
	content := a.buildEventContent(ev)

	return &Artifact{
		ID:           artifactID,
		Path:         ".sdp/log/events.jsonl#" + ev.ID,
		Type:         "evidence",
		Title:        ev.Type + ": " + ev.WSID,
		Content:      content,
		FeatureID:    extractFeatureFromWSID(ev.WSID),
		WorkstreamID: ev.WSID,
		Tags:         []string{"evidence", ev.Type},
		FileHash:     ev.ID, // Use event ID as hash for evidence
		IndexedAt:    time.Now(),
	}
}

// buildEventContent creates searchable content from an event
func (a *EvidenceAdapter) buildEventContent(ev evidence.Event) string {
	content := ev.Type + " event for " + ev.WSID

	// Add data based on event type
	switch ev.Type {
	case "generation":
		if gen, ok := ev.Data.(evidence.GenerationData); ok {
			content += " model:" + gen.ModelID
			if len(gen.FilesChanged) > 0 {
				content += " files:" + joinFiles(gen.FilesChanged)
			}
		}
	case "verification":
		if ver, ok := ev.Data.(evidence.VerificationData); ok {
			if ver.Passed {
				content += " passed"
			} else {
				content += " failed"
			}
			if ver.GateName != "" {
				content += " gate:" + ver.GateName
			}
		}
	case "decision":
		if dec, ok := ev.Data.(evidence.DecisionEventData); ok {
			content += " question:" + dec.Question + " choice:" + dec.Choice
		}
	case "lesson":
		if les, ok := ev.Data.(evidence.LessonEventData); ok {
			content += " insight:" + les.Insight
		}
	}

	// Include full data as JSON
	dataBytes, _ := json.Marshal(ev.Data)
	content += " " + string(dataBytes)

	return content
}

// joinFiles joins file paths for content
func joinFiles(files []string) string {
	if len(files) == 0 {
		return ""
	}
	if len(files) > 3 {
		return files[0] + "," + files[1] + " and " + string(rune(len(files)-2)) + " more"
	}
	result := files[0]
	for i := 1; i < len(files); i++ {
		result += "," + files[i]
	}
	return result
}

// extractFeatureFromWSID extracts feature ID from workstream ID (PP-FFF-SS -> FFF)
func extractFeatureFromWSID(wsid string) string {
	if len(wsid) < 6 {
		return ""
	}
	// Format: PP-FFF-SS (e.g., 00-050-01)
	parts := splitWSID(wsid)
	if len(parts) >= 2 {
		return "F" + parts[1]
	}
	return ""
}

// splitWSID splits a workstream ID into its components
func splitWSID(wsid string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(wsid); i++ {
		if wsid[i] == '-' {
			parts = append(parts, wsid[start:i])
			start = i + 1
		}
	}
	if start < len(wsid) {
		parts = append(parts, wsid[start:])
	}
	return parts
}
