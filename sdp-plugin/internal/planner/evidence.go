package planner

import (
	"fmt"
	"path/filepath"

	"github.com/fall-out-bug/sdp/internal/evidence"
)

// EmitPlanEvent emits a plan event to the evidence log.
// AC5: Emits plan event with workstream information.
func (p *Planner) EmitPlanEvent(result *DecompositionResult) error {
	if p.EvidenceWriter == nil {
		return fmt.Errorf("evidence writer not configured")
	}

	// Build scope files list (all workstream files that would be created)
	scopeFiles := []string{}
	for _, ws := range result.Workstreams {
		scopeFiles = append(scopeFiles, filepath.Join(p.BacklogDir, ws.Filename()))
	}

	// Create plan event using existing evidence package
	ev := evidence.PlanEventWithFeature(
		"00-057-00", // Planning workstream ID
		result.FeatureID,
		scopeFiles,
	)

	// Add additional metadata about the decomposition
	if ev.Data == nil {
		ev.Data = make(map[string]interface{})
	}
	if dataMap, ok := ev.Data.(map[string]interface{}); ok {
		dataMap["workstream_count"] = len(result.Workstreams)
		dataMap["dependency_count"] = len(result.Dependencies)
		dataMap["description"] = p.Description
		dataMap["action"] = "decompose"
	}

	return p.EvidenceWriter.Append(ev)
}
