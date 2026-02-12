package resolver

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBeadsDetector_IsAvailable(t *testing.T) {
	detector := NewBeadsDetector()

	// This test will pass whether or not beads is installed
	// It just tests the detection logic runs without error
	available := detector.IsAvailable()
	t.Logf("Beads available: %v", available)
}

func TestBeadsDetector_HasBeadsDirectory(t *testing.T) {
	t.Run("directory exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		beadsDir := filepath.Join(tmpDir, ".beads")
		if err := os.MkdirAll(beadsDir, 0755); err != nil {
			t.Fatal(err)
		}

		detector := NewBeadsDetectorWithDir(tmpDir)
		if !detector.HasBeadsDirectory() {
			t.Error("expected .beads directory to be detected")
		}
	})

	t.Run("directory missing", func(t *testing.T) {
		tmpDir := t.TempDir()
		detector := NewBeadsDetectorWithDir(tmpDir)
		if detector.HasBeadsDirectory() {
			t.Error("expected .beads directory NOT to be detected")
		}
	})
}

func TestResolver_ResolveBeadsID(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create workstream with beads_id
	wsContent := `---
ws_id: 00-064-01
feature_id: F064
beads_id: sdp-test1
title: "Test"
---
## Goal
Test
`
	if err := os.WriteFile(filepath.Join(wsDir, "00-064-01.md"), []byte(wsContent), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewResolver(WithWorkstreamDir(wsDir))

	t.Run("find by beads ID", func(t *testing.T) {
		result, err := r.ResolveBeadsID("sdp-test1")
		if err != nil {
			t.Fatalf("ResolveBeadsID() error = %v", err)
		}
		if result.WSID != "00-064-01" {
			t.Errorf("expected WSID '00-064-01', got %s", result.WSID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := r.ResolveBeadsID("sdp-unknown")
		if err == nil {
			t.Error("expected error for unknown beads ID")
		}
	})
}

func TestBeadsDetector_IsEnabled(t *testing.T) {
	t.Run("with beads dir and CLI", func(t *testing.T) {
		tmpDir := t.TempDir()
		beadsDir := filepath.Join(tmpDir, ".beads")
		if err := os.MkdirAll(beadsDir, 0755); err != nil {
			t.Fatal(err)
		}

		detector := NewBeadsDetectorWithDir(tmpDir)
		// IsEnabled also checks CLI, which may not be installed
		// So we just test it runs without panic
		_ = detector.IsEnabled()
	})
}

func TestCreateBeadsIssue_NotEnabled(t *testing.T) {
	// When beads is not enabled, should return empty without error
	// This tests the error path
	_, err := CreateBeadsIssue("Test", "bug")
	// Will fail because beads CLI is not available, which is expected
	if err == nil {
		t.Log("beads CLI available (unexpected in test env)")
	}
}

func TestUpdateBeadsNotes_NotEnabled(t *testing.T) {
	// When beads is not enabled, should return nil
	// Create temp detector that's disabled
	detector := NewBeadsDetector()
	if !detector.IsEnabled() {
		err := UpdateBeadsNotes("sdp-test", "path/to/ws.md")
		if err != nil {
			t.Errorf("expected nil when beads disabled, got: %v", err)
		}
	}
}

func TestParseBeadsIDFromOutput(t *testing.T) {
	tests := []struct {
		output   string
		expected string
	}{
		{"Created issue: sdp-abc123", "sdp-abc123"},
		{"sdp-xyz789", "sdp-xyz789"},
		{"Some output\nsdp-test\nMore", "sdp-test"},
		{"No ID here", ""},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.output, func(t *testing.T) {
			result := parseBeadsIDFromOutput(tt.output)
			if result != tt.expected {
				t.Errorf("parseBeadsIDFromOutput(%q) = %q, want %q", tt.output, result, tt.expected)
			}
		})
	}
}

func TestResolver_FindWorkstreamByBeadsID(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create multiple workstreams
	testCases := []struct {
		filename string
		content  string
	}{
		{
			"00-064-01.md",
			"---\nws_id: 00-064-01\nbeads_id: sdp-abc\n---\n## Goal\nTest1",
		},
		{
			"00-064-02.md",
			"---\nws_id: 00-064-02\nbeads_id: sdp-xyz\n---\n## Goal\nTest2",
		},
		{
			"00-064-03.md",
			"---\nws_id: 00-064-03\n---\n## Goal\nNo beads",
		},
	}

	for _, tc := range testCases {
		if err := os.WriteFile(filepath.Join(wsDir, tc.filename), []byte(tc.content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	r := NewResolver(WithWorkstreamDir(wsDir))

	t.Run("find existing", func(t *testing.T) {
		wsID, path, err := r.FindWorkstreamByBeadsID("sdp-abc")
		if err != nil {
			t.Fatalf("FindWorkstreamByBeadsID() error = %v", err)
		}
		if wsID != "00-064-01" {
			t.Errorf("expected wsID '00-064-01', got %s", wsID)
		}
		if path == "" {
			t.Error("expected non-empty path")
		}
	})

	t.Run("find another", func(t *testing.T) {
		wsID, _, err := r.FindWorkstreamByBeadsID("sdp-xyz")
		if err != nil {
			t.Fatal(err)
		}
		if wsID != "00-064-02" {
			t.Errorf("expected wsID '00-064-02', got %s", wsID)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, _, err := r.FindWorkstreamByBeadsID("sdp-missing")
		if err == nil {
			t.Error("expected error for missing beads ID")
		}
	})
}
