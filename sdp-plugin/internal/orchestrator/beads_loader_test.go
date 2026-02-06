package orchestrator

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/fall-out-bug/sdp/internal/parser"
)

func TestBeadsLoader_LoadWorkstreams(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Create test workstream files
	ws1Content := `---
ws_id: 00-050-01
parent: sdp-123
feature: F050
status: backlog
size: SMALL
project_id: "00"
---

## WS-00-050-01: First Workstream

### Goal
Test workstream

### Acceptance Criteria
- [ ] AC1: Test
`
	ws2Content := `---
ws_id: 00-050-02
parent: sdp-124
feature: F050
status: backlog
size: SMALL
project_id: "00"
---

## WS-00-050-02: Second Workstream

### Goal
Test workstream 2

### Dependencies
- 00-050-01

### Acceptance Criteria
- [ ] AC1: Test
`
	ws3Content := `---
ws_id: 00-050-03
parent: sdp-125
feature: F051
status: backlog
size: SMALL
project_id: "00"
---

## WS-00-050-03: Different Feature

### Goal
Different feature

### Acceptance Criteria
- [ ] AC1: Test
`

	// Write workstream files
	if err := os.WriteFile(filepath.Join(tmpDir, "00-050-01.md"), []byte(ws1Content), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "00-050-02.md"), []byte(ws2Content), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "00-050-03.md"), []byte(ws3Content), 0644); err != nil {
		t.Fatal(err)
	}

	// Create loader
	loader := NewBeadsLoader(tmpDir, ".beads-sdp-mapping.jsonl")

	// Test loading workstreams for F050
	workstreams, err := loader.LoadWorkstreams("F050")
	if err != nil {
		t.Errorf("LoadWorkstreams failed: %v", err)
	}

	if len(workstreams) != 2 {
		t.Errorf("got %d workstreams, want 2", len(workstreams))
		return
	}

	// Find first and second workstreams (order may vary)
	var ws1, ws2 *WorkstreamNode
	for i := range workstreams {
		if workstreams[i].ID == "00-050-01" {
			ws1 = &workstreams[i]
		}
		if workstreams[i].ID == "00-050-02" {
			ws2 = &workstreams[i]
		}
	}

	if ws1 == nil {
		t.Error("workstream 00-050-01 not found")
	}
	if ws2 == nil {
		t.Error("workstream 00-050-02 not found")
		return
	}

	// Verify second workstream has dependency
	if len(ws2.Dependencies) != 1 {
		t.Errorf("second workstream has %d dependencies, want 1", len(ws2.Dependencies))
		return
	}
	if ws2.Dependencies[0] != "00-050-01" {
		t.Errorf("second workstream dependency = %s, want 00-050-01", ws2.Dependencies[0])
	}
}

func TestBeadsLoader_LoadWorkstreams_NotFound(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	loader := NewBeadsLoader(tmpDir, ".beads-sdp-mapping.jsonl")

	// Test loading workstreams for non-existent feature
	_, err := loader.LoadWorkstreams("F999")
	if err == nil {
		t.Error("expected error for non-existent feature")
	}

	if !errors.Is(err, ErrFeatureNotFound) {
		t.Errorf("error = %v, want ErrFeatureNotFound", err)
	}
}

func TestBeadsLoader_ExtractDependencies(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantDeps []string
	}{
		{
			name: "dependencies section",
			content: `---
ws_id: 00-050-01
feature: F050
---

## Goal
Test

### Dependencies
- 00-050-02
- 00-050-03
`,
			wantDeps: []string{"00-050-02", "00-050-03"},
		},
		{
			name: "bold dependencies header",
			content: `---
ws_id: 00-050-01
feature: F050
---

## Goal
Test

**Dependencies:**
- 00-050-02
`,
			wantDeps: []string{"00-050-02"},
		},
		{
			name: "no dependencies",
			content: `---
ws_id: 00-050-01
feature: F050
---

## Goal
Test
`,
			wantDeps: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			wsPath := filepath.Join(tmpDir, "00-050-01.md")

			if err := os.WriteFile(wsPath, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			loader := NewBeadsLoader(tmpDir, ".beads-sdp-mapping.jsonl")

			// Parse workstream
			ws, err := parser.ParseWorkstream(wsPath)
			if err != nil {
				t.Fatal(err)
			}

			// Extract dependencies
			deps := loader.extractDependencies(ws)

			if len(deps) != len(tt.wantDeps) {
				t.Errorf("got %d dependencies, want %d", len(deps), len(tt.wantDeps))
			}

			for i, dep := range tt.wantDeps {
				if i >= len(deps) || deps[i] != dep {
					t.Errorf("dependency %d = %s, want %s", i, deps[i], dep)
				}
			}
		})
	}
}
