package resolver

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectIDType(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected IDType
	}{
		// Workstream IDs: ^\d{2}-\d{3}-\d{2}$
		{"standard workstream", "00-064-01", TypeWorkstream},
		{"workstream with zeros", "00-001-00", TypeWorkstream},
		{"workstream high numbers", "99-999-99", TypeWorkstream},

		// Fix workstream IDs: ^\d{2}-[A-Z]\d{3}-\d{4}$
		{"fix workstream", "99-F064-0001", TypeWorkstream},
		{"fix workstream feature", "01-F001-0010", TypeWorkstream},

		// Beads IDs: ^[a-z]{3}-[a-z0-9]+$
		{"beads ID standard", "sdp-ushh", TypeBeads},
		{"beads ID with numbers", "abc-123", TypeBeads},
		{"beads ID short suffix", "xyz-a", TypeBeads},

		// Issue IDs: ^ISSUE-\d+$
		{"issue ID", "ISSUE-0001", TypeIssue},
		{"issue ID simple", "ISSUE-1", TypeIssue},
		{"issue ID large", "ISSUE-99999", TypeIssue},

		// Invalid/unknown
		{"empty string", "", TypeUnknown},
		{"random string", "random", TypeUnknown},
		{"uuid-like", "a1b2c3d4-e5f6-7890", TypeUnknown},
		{"camelCase", "MyTaskID", TypeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectIDType(tt.id)
			if result != tt.expected {
				t.Errorf("DetectIDType(%q) = %v, want %v", tt.id, result, tt.expected)
			}
		})
	}
}

func TestResolver_ResolveWorkstream(t *testing.T) {
	// Create temp directory with workstream files
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test workstream file
	wsContent := `---
ws_id: 00-064-01
feature_id: F064
title: "Test Workstream"
status: backlog
---
## Goal
Test goal
`
	wsPath := filepath.Join(wsDir, "00-064-01.md")
	if err := os.WriteFile(wsPath, []byte(wsContent), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewResolver(WithWorkstreamDir(wsDir))

	t.Run("resolve existing workstream", func(t *testing.T) {
		result, err := r.Resolve("00-064-01")
		if err != nil {
			t.Fatalf("Resolve() error = %v", err)
		}
		if result.Type != TypeWorkstream {
			t.Errorf("Resolve() type = %v, want %v", result.Type, TypeWorkstream)
		}
		if result.Path != wsPath {
			t.Errorf("Resolve() path = %v, want %v", result.Path, wsPath)
		}
	})

	t.Run("resolve non-existent workstream", func(t *testing.T) {
		_, err := r.Resolve("00-999-99")
		if err == nil {
			t.Error("Resolve() expected error for non-existent workstream")
		}
	})
}

func TestResolver_ResolveBeads(t *testing.T) {
	// Create temp directory with workstream files containing beads_id
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test workstream file with beads_id
	wsContent := `---
ws_id: 00-064-01
feature_id: F064
beads_id: sdp-ushh
title: "Test Workstream"
---
## Goal
Test goal
`
	wsPath := filepath.Join(wsDir, "00-064-01.md")
	if err := os.WriteFile(wsPath, []byte(wsContent), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewResolver(WithWorkstreamDir(wsDir))

	t.Run("resolve beads to workstream", func(t *testing.T) {
		result, err := r.Resolve("sdp-ushh")
		if err != nil {
			t.Fatalf("Resolve() error = %v", err)
		}
		if result.Type != TypeBeads {
			t.Errorf("Resolve() type = %v, want %v", result.Type, TypeBeads)
		}
		if result.Path != wsPath {
			t.Errorf("Resolve() path = %v, want %v", result.Path, wsPath)
		}
		if result.WSID != "00-064-01" {
			t.Errorf("Resolve() ws_id = %v, want %v", result.WSID, "00-064-01")
		}
	})

	t.Run("resolve non-existent beads", func(t *testing.T) {
		_, err := r.Resolve("xyz-unknown")
		if err == nil {
			t.Error("Resolve() expected error for non-existent beads ID")
		}
	})
}

func TestResolver_ResolveIssue(t *testing.T) {
	// Create temp directory with issues
	tmpDir := t.TempDir()
	issuesDir := filepath.Join(tmpDir, "docs", "issues")
	if err := os.MkdirAll(issuesDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test issue file
	issueContent := `---
issue_id: ISSUE-0001
title: "Test Issue"
status: open
---
## Symptom
Test symptom
`
	issuePath := filepath.Join(issuesDir, "ISSUE-0001.md")
	if err := os.WriteFile(issuePath, []byte(issueContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create index file
	indexDir := filepath.Join(tmpDir, ".sdp")
	if err := os.MkdirAll(indexDir, 0755); err != nil {
		t.Fatal(err)
	}
	indexContent := `{"issue_id":"ISSUE-0001","title":"Test Issue","status":"open","file":"docs/issues/ISSUE-0001.md"}
`
	if err := os.WriteFile(filepath.Join(indexDir, "issues-index.jsonl"), []byte(indexContent), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewResolver(
		WithIssuesDir(issuesDir),
		WithIndexFile(filepath.Join(indexDir, "issues-index.jsonl")),
	)

	t.Run("resolve issue from index", func(t *testing.T) {
		result, err := r.Resolve("ISSUE-0001")
		if err != nil {
			t.Fatalf("Resolve() error = %v", err)
		}
		if result.Type != TypeIssue {
			t.Errorf("Resolve() type = %v, want %v", result.Type, TypeIssue)
		}
		if result.Path != issuePath {
			t.Errorf("Resolve() path = %v, want %v", result.Path, issuePath)
		}
	})

	t.Run("resolve non-existent issue", func(t *testing.T) {
		_, err := r.Resolve("ISSUE-9999")
		if err == nil {
			t.Error("Resolve() expected error for non-existent issue")
		}
	})
}

func TestResolver_ResolveIssueWithoutIndex(t *testing.T) {
	// Create temp directory with issues (no index)
	tmpDir := t.TempDir()
	issuesDir := filepath.Join(tmpDir, "docs", "issues")
	if err := os.MkdirAll(issuesDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create test issue file
	issueContent := `---
issue_id: ISSUE-0042
title: "Test Issue"
---
## Symptom
Test
`
	issuePath := filepath.Join(issuesDir, "ISSUE-0042.md")
	if err := os.WriteFile(issuePath, []byte(issueContent), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewResolver(WithIssuesDir(issuesDir))

	t.Run("resolve issue via filesystem fallback", func(t *testing.T) {
		result, err := r.Resolve("ISSUE-0042")
		if err != nil {
			t.Fatalf("Resolve() error = %v", err)
		}
		if result.Type != TypeIssue {
			t.Errorf("Resolve() type = %v, want %v", result.Type, TypeIssue)
		}
		if result.Path != issuePath {
			t.Errorf("Resolve() path = %v, want %v", result.Path, issuePath)
		}
	})
}

func TestResolver_AmbiguousID(t *testing.T) {
	t.Run("ambiguous format returns valid detection", func(t *testing.T) {
		// This ID doesn't match any known format
		result := DetectIDType("ambiguous-id")
		if result != TypeUnknown {
			t.Errorf("DetectIDType('ambiguous-id') = %v, want %v", result, TypeUnknown)
		}
	})
}

func TestResolver_ErrorHandling(t *testing.T) {
	r := NewResolver()

	t.Run("empty ID", func(t *testing.T) {
		_, err := r.Resolve("")
		if err == nil {
			t.Error("Resolve('') expected error")
		}
	})

	t.Run("whitespace ID", func(t *testing.T) {
		_, err := r.Resolve("   ")
		if err == nil {
			t.Error("Resolve('   ') expected error")
		}
	})
}
