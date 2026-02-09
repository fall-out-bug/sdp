package parser

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestParseValidWorkstream(t *testing.T) {
	// Create a temporary valid workstream file
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "00-050-01.md")
	content := `---
ws_id: 00-050-01
parent: sdp-79u
feature: F050
status: backlog
size: MEDIUM
project_id: 00
---

## WS-00-050-01: Workstream Parser

### Goal

Parse workstream markdown files

### Acceptance Criteria

- [ ] AC1: Parse valid WS markdown
- [ ] AC2: Validate WS ID format

### Scope Files

- internal/parser/workstream.go
`
	err := os.WriteFile(wsPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test parsing
	ws, err := ParseWorkstream(wsPath)
	if err != nil {
		t.Fatalf("ParseWorkstream failed: %v", err)
	}

	// Validate parsed data
	if ws.ID != "00-050-01" {
		t.Errorf("Expected ID 00-050-01, got %s", ws.ID)
	}
	if ws.Feature != "F050" {
		t.Errorf("Expected Feature F050, got %s", ws.Feature)
	}
	if ws.Status != "backlog" {
		t.Errorf("Expected Status backlog, got %s", ws.Status)
	}
	if ws.Goal == "" {
		t.Error("Goal should not be empty")
	}
	if len(ws.Acceptance) == 0 {
		t.Error("Acceptance criteria should not be empty")
	}
}

func TestValidateInvalidWSID(t *testing.T) {
	tests := []struct {
		name  string
		wsID  string
		valid bool
	}{
		{"Valid format", "00-001-01", true},
		{"Valid format", "99-999-99", true},
		{"Missing leading zeros", "0-1-1", false},
		{"Wrong separator", "00-001_01", false},
		{"Too many digits", "000-001-01", false},
		{"Letters", "00-abc-01", false},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws := &Workstream{ID: tt.wsID}
			err := ws.Validate()
			if tt.valid && err != nil {
				t.Errorf("Expected valid WS ID %s, got error: %v", tt.wsID, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("Expected invalid WS ID %s, but got no error", tt.wsID)
			}
		})
	}
}

func TestParseMissingRequiredFields(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "00-050-01.md")

	// Missing ws_id field
	content := `---
feature: F050
status: backlog
---
# Test
`
	err := os.WriteFile(wsPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = ParseWorkstream(wsPath)
	if err == nil {
		t.Error("Expected error for missing required fields, got nil")
	}
}

func TestParseEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "00-050-01.md")

	err := os.WriteFile(wsPath, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = ParseWorkstream(wsPath)
	if err == nil {
		t.Error("Expected error for empty file, got nil")
	}
}

func TestParseMalformedYAML(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "00-050-01.md")

	// Malformed YAML (unclosed bracket)
	content := `---
ws_id: [00-050-01
feature: F050
---
# Test
`
	err := os.WriteFile(wsPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = ParseWorkstream(wsPath)
	if err == nil {
		t.Error("Expected error for malformed YAML, got nil")
	}
}

func TestParseMissingFrontmatter(t *testing.T) {
	tmpDir := t.TempDir()
	wsPath := filepath.Join(tmpDir, "00-050-01.md")

	// No frontmatter delimiters
	content := `ws_id: 00-050-01
feature: F050
# Test
`
	err := os.WriteFile(wsPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = ParseWorkstream(wsPath)
	if err == nil {
		t.Error("Expected error for missing frontmatter, got nil")
	}
}

func TestValidateFile(t *testing.T) {
	tmpDir := t.TempDir()
	validPath := filepath.Join(tmpDir, "00-050-01.md")
	content := `---
ws_id: 00-050-01
parent: sdp-79u
feature: F050
status: backlog
size: MEDIUM
project_id: 00
---

## Test Workstream

### Goal

Test workstream for validation

### Acceptance Criteria

- [ ] AC1: Valid workstream
`
	os.WriteFile(validPath, []byte(content), 0644)

	// Test valid file
	issues, err := ValidateFile(validPath)
	if err != nil {
		t.Fatalf("ValidateFile failed: %v", err)
	}
	if len(issues) > 0 {
		for _, issue := range issues {
			t.Logf("Validation issue: %s - %s (%s)", issue.Field, issue.Message, issue.Severity)
		}
		// Only fail if there are ERROR severity issues
		hasErrors := false
		for _, issue := range issues {
			if issue.Severity == "ERROR" {
				hasErrors = true
				break
			}
		}
		if hasErrors {
			t.Errorf("Expected no ERROR validation issues for valid file, got %d", len(issues))
		}
	}

	// Test invalid WS ID
	invalidPath := filepath.Join(tmpDir, "00-050-02.md")
	invalidContent := `---
ws_id: invalid-id
feature: F050
status: backlog
---
# Test
`
	os.WriteFile(invalidPath, []byte(invalidContent), 0644)

	issues, err = ValidateFile(invalidPath)
	if err != nil {
		t.Fatalf("ValidateFile failed: %v", err)
	}
	if len(issues) == 0 {
		t.Error("Expected validation issues for invalid WS ID, got none")
	}
}

// schemaIndexEntry represents an entry in schema/index.json
type schemaIndexEntry struct {
	ID    string `json:"id"`
	Path  string `json:"path"`
	Title string `json:"title"`
}

// schemaIndex represents schema/index.json structure
type schemaIndex struct {
	Version int                `json:"version"`
	Schemas []schemaIndexEntry `json:"schemas"`
}

// TestSchemaRegistryLoads verifies that the canonical schema registry and
// schema files exist and are valid JSON (AC5: any code loading schemas still works).
func TestSchemaRegistryLoads(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	// Walk up to find schema/ (repo root when run from sdp-plugin is ..)
	var schemaDir string
	for base := wd; base != ""; base = filepath.Dir(base) {
		idx := filepath.Join(base, "schema", "index.json")
		if _, err := os.Stat(idx); err == nil {
			schemaDir = filepath.Join(base, "schema")
			break
		}
		if base == filepath.Dir(base) {
			break
		}
	}
	if schemaDir == "" {
		t.Skip("schema/index.json not found (run from repo containing schema/)")
	}
	indexPath := filepath.Join(schemaDir, "index.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("read index: %v", err)
	}
	var idx schemaIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		t.Fatalf("index.json invalid JSON: %v", err)
	}
	if idx.Version < 1 || len(idx.Schemas) == 0 {
		t.Error("index.json must have version >= 1 and non-empty schemas")
	}
	for _, s := range idx.Schemas {
		p := filepath.Join(schemaDir, s.Path)
		body, err := os.ReadFile(p)
		if err != nil {
			t.Errorf("schema file %s: %v", s.Path, err)
			continue
		}
		var js map[string]interface{}
		if err := json.Unmarshal(body, &js); err != nil {
			t.Errorf("schema %s invalid JSON: %v", s.Path, err)
		}
	}
}
