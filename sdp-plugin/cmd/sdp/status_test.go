package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestStatusCmd_TextMode(t *testing.T) {
	// Save original stdout
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()

	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Create minimal project structure
	os.MkdirAll(".git", 0755)
	os.MkdirAll(".claude", 0755)
	os.MkdirAll(".sdp", 0755)

	cmd := statusCmd()
	if err := cmd.Flags().Set("text", "true"); err != nil {
		t.Fatalf("Failed to set text flag: %v", err)
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := cmd.RunE(cmd, []string{})

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("statusCmd() error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check output contains expected sections
	if !strings.Contains(output, "SDP Project Status") {
		t.Error("Output should contain 'SDP Project Status'")
	}
	if !strings.Contains(output, "Environment:") {
		t.Error("Output should contain 'Environment:'")
	}
	if !strings.Contains(output, "Workstreams:") {
		t.Error("Output should contain 'Workstreams:'")
	}
	if !strings.Contains(output, "Next Action:") {
		t.Error("Output should contain 'Next Action:'")
	}
}

func TestStatusCmd_JSONMode(t *testing.T) {
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()

	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Create minimal project structure
	os.MkdirAll(".git", 0755)
	os.MkdirAll(".claude", 0755)
	os.MkdirAll(".sdp", 0755)

	cmd := statusCmd()
	if err := cmd.Flags().Set("json", "true"); err != nil {
		t.Fatalf("Failed to set json flag: %v", err)
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := cmd.RunE(cmd, []string{})

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("statusCmd() error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check JSON output
	if !strings.Contains(output, `"has_git":true`) {
		t.Error("JSON output should contain has_git:true")
	}
	if !strings.Contains(output, `"has_claude":true`) {
		t.Error("JSON output should contain has_claude:true")
	}
	if !strings.Contains(output, `"workstreams"`) {
		t.Error("JSON output should contain workstreams object")
	}
	if !strings.Contains(output, `"next_action"`) {
		t.Error("JSON output should contain next_action")
	}
}

func TestGatherProjectStatus(t *testing.T) {
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()

	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Test with no structure
	status := gatherProjectStatus()
	if status.HasGit {
		t.Error("HasGit should be false without .git")
	}
	if status.HasClaude {
		t.Error("HasClaude should be false without .claude")
	}
	if status.HasSDP {
		t.Error("HasSDP should be false without .sdp")
	}

	// Create structure
	os.MkdirAll(".git", 0755)
	os.MkdirAll(".claude", 0755)
	os.MkdirAll(".sdp", 0755)
	os.MkdirAll(".beads", 0755)
	os.WriteFile(".beads/issues.jsonl", []byte("{}"), 0644)

	status = gatherProjectStatus()
	if !status.HasGit {
		t.Error("HasGit should be true with .git")
	}
	if !status.HasClaude {
		t.Error("HasClaude should be true with .claude")
	}
	if !status.HasSDP {
		t.Error("HasSDP should be true with .sdp")
	}
	if !status.HasBeads {
		t.Error("HasBeads should be true with .beads/issues.jsonl")
	}
}

func TestDetermineNextAction(t *testing.T) {
	tests := []struct {
		name       string
		status     *ProjectStatus
		contains   string
		notContain string
	}{
		{
			name: "no git",
			status: &ProjectStatus{
				HasGit:    false,
				HasClaude: false,
				HasSDP:    false,
			},
			contains: "git init",
		},
		{
			name: "no claude",
			status: &ProjectStatus{
				HasGit:    true,
				HasClaude: false,
				HasSDP:    false,
			},
			contains: "sdp init --guided",
		},
		{
			name: "no sdp",
			status: &ProjectStatus{
				HasGit:    true,
				HasClaude: true,
				HasSDP:    false,
			},
			contains: "sdp init",
		},
		{
			name: "active session",
			status: &ProjectStatus{
				HasGit:    true,
				HasClaude: true,
				HasSDP:    true,
				ActiveSession: &ActiveSession{
					WorkstreamID: "00-001-01",
				},
			},
			contains: "sdp apply --ws 00-001-01",
		},
		{
			name: "ready for feature",
			status: &ProjectStatus{
				HasGit:    true,
				HasClaude: true,
				HasSDP:    true,
				Workstreams: WorkstreamSummary{
					Open: 0,
				},
			},
			contains: "sdp plan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := determineNextAction(tt.status)
			if tt.contains != "" && !strings.Contains(action, tt.contains) {
				t.Errorf("Next action = %q, should contain %q", action, tt.contains)
			}
			if tt.notContain != "" && strings.Contains(action, tt.notContain) {
				t.Errorf("Next action = %q, should NOT contain %q", action, tt.notContain)
			}
		})
	}
}

func TestCountWorkstreams(t *testing.T) {
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()

	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// No directories
	summary := countWorkstreams()
	if summary.Open != 0 {
		t.Errorf("Open should be 0, got %d", summary.Open)
	}

	// Create backlog with workstream files
	os.MkdirAll("docs/workstreams/backlog", 0755)
	os.WriteFile("docs/workstreams/backlog/00-001-01.md", []byte(""), 0644)
	os.WriteFile("docs/workstreams/backlog/00-001-02.md", []byte(""), 0644)
	os.WriteFile("docs/workstreams/backlog/README.md", []byte(""), 0644) // Not a workstream

	summary = countWorkstreams()
	if summary.Open != 2 {
		t.Errorf("Open should be 2, got %d", summary.Open)
	}

	// Create completed directory
	os.MkdirAll("docs/workstreams/completed", 0755)
	os.WriteFile("docs/workstreams/completed/00-000-01.md", []byte(""), 0644)

	summary = countWorkstreams()
	if summary.Completed != 1 {
		t.Errorf("Completed should be 1, got %d", summary.Completed)
	}
}

func TestBoolIcon(t *testing.T) {
	if boolIcon(true) != "[OK]" {
		t.Error("boolIcon(true) should be '[OK]'")
	}
	if boolIcon(false) != "[MISSING]" {
		t.Error("boolIcon(false) should be '[MISSING]'")
	}
}

func TestDirExists(t *testing.T) {
	tmpDir := t.TempDir()

	if dirExists(tmpDir) != true {
		t.Error("dirExists should return true for existing directory")
	}
	if dirExists(tmpDir+"/nonexistent") != false {
		t.Error("dirExists should return false for non-existent directory")
	}
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := tmpDir + "/test.txt"
	os.WriteFile(testFile, []byte("test"), 0644)

	if fileExists(testFile) != true {
		t.Error("fileExists should return true for existing file")
	}
	if fileExists(tmpDir+"/nonexistent.txt") != false {
		t.Error("fileExists should return false for non-existent file")
	}
	if fileExists(tmpDir) != false {
		t.Error("fileExists should return false for directory")
	}
}

func TestFindJSONField(t *testing.T) {
	content := `{"workstream_id": "00-001-01", "feature_id": "F001"}`
	idx := findJSONField(content, "workstream_id")
	if idx < 0 {
		t.Error("findJSONField should find workstream_id")
	}

	idx = findJSONField(content, "nonexistent")
	if idx >= 0 {
		t.Error("findJSONField should return -1 for nonexistent field")
	}
}

func TestExtractJSONString(t *testing.T) {
	content := `{"workstream_id": "00-001-01", "feature_id": "F001"}`
	idx := findJSONField(content, "workstream_id")
	value := extractJSONString(content, idx)
	if value != "00-001-01" {
		t.Errorf("extractJSONString = %q, want '00-001-01'", value)
	}

	// Test with different spacing
	content2 := `{ "name" : "test value" }`
	idx2 := findJSONField(content2, "name")
	value2 := extractJSONString(content2, idx2)
	if value2 != "test value" {
		t.Errorf("extractJSONString = %q, want 'test value'", value2)
	}
}

func TestGetActiveSession(t *testing.T) {
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()

	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// No session file
	session := getActiveSession()
	if session != nil {
		t.Error("getActiveSession should return nil without session file")
	}

	// Create session file
	os.MkdirAll(".sdp", 0755)
	os.WriteFile(".sdp/session.json", []byte(`{"workstream_id": "00-001-01", "feature_id": "F001"}`), 0644)

	session = getActiveSession()
	if session == nil {
		t.Fatal("getActiveSession should return session with session file")
	}
	if session.WorkstreamID != "00-001-01" {
		t.Errorf("WorkstreamID = %q, want '00-001-01'", session.WorkstreamID)
	}
	if session.FeatureID != "F001" {
		t.Errorf("FeatureID = %q, want 'F001'", session.FeatureID)
	}
}

func TestProjectStatus_PrintText(t *testing.T) {
	status := &ProjectStatus{
		HasGit:    true,
		HasClaude: true,
		HasSDP:    true,
		HasBeads:  false,
		Workstreams: WorkstreamSummary{
			Open:       5,
			InProgress: 2,
			Completed:  10,
			Blocked:    1,
		},
		NextAction: "Test action",
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := status.printText()

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("printText error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "[OK]") {
		t.Error("Output should show [OK] for true values")
	}
	if !strings.Contains(output, "[MISSING]") {
		t.Error("Output should show [MISSING] for false values")
	}
}

func TestProjectStatus_PrintJSON(t *testing.T) {
	status := &ProjectStatus{
		HasGit:    true,
		HasClaude: false,
		HasSDP:    true,
		HasBeads:  false,
		Workstreams: WorkstreamSummary{
			Open:       5,
			InProgress: 2,
			Completed:  10,
			Blocked:    1,
		},
		NextAction: "Test action",
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := status.printJSON()

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("printJSON error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, `"has_git":true`) {
		t.Error("JSON should contain has_git:true")
	}
	if !strings.Contains(output, `"has_claude":false`) {
		t.Error("JSON should contain has_claude:false")
	}
}
