package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTaskCreateCmd_Workstream(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatal(err)
	}

	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	cmd := taskCreateCmd()
	cmd.SetArgs([]string{
		"--type=bug",
		"--title=Test Bug",
		"--feature=F064",
		"--priority=1",
	})

	if err := cmd.Execute(); err != nil {
		t.Errorf("taskCreateCmd.Execute() error = %v", err)
	}
}

func TestTaskCreateCmd_Issue(t *testing.T) {
	tmpDir := t.TempDir()
	issuesDir := filepath.Join(tmpDir, "docs", "issues")
	if err := os.MkdirAll(issuesDir, 0755); err != nil {
		t.Fatal(err)
	}

	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	cmd := taskCreateCmd()
	cmd.SetArgs([]string{
		"--type=bug",
		"--title=Test Issue",
		"--issue",
	})

	if err := cmd.Execute(); err != nil {
		t.Errorf("taskCreateCmd.Execute() error = %v", err)
	}
}

func TestTaskCreateCmd_JSON(t *testing.T) {
	tmpDir := t.TempDir()
	wsDir := filepath.Join(tmpDir, "docs", "workstreams", "backlog")
	if err := os.MkdirAll(wsDir, 0755); err != nil {
		t.Fatal(err)
	}

	oldDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)

	cmd := taskCreateCmd()
	cmd.SetArgs([]string{
		"--type=task",
		"--title=Test",
		"--feature=F064",
		"--json",
	})

	if err := cmd.Execute(); err != nil {
		t.Errorf("taskCreateCmd.Execute() --json error = %v", err)
	}
}

func TestTaskCreateCmd_InvalidType(t *testing.T) {
	cmd := taskCreateCmd()
	cmd.SetArgs([]string{
		"--type=invalid",
		"--title=Test",
		"--feature=F064",
	})

	if err := cmd.Execute(); err == nil {
		t.Error("expected error for invalid type")
	} else if !strings.Contains(err.Error(), "invalid task type") {
		t.Errorf("expected invalid task type error, got: %v", err)
	}
}

func TestTaskCreateCmd_InvalidPriority(t *testing.T) {
	cmd := taskCreateCmd()
	cmd.SetArgs([]string{
		"--type=bug",
		"--title=Test",
		"--feature=F064",
		"--priority=5",
	})

	if err := cmd.Execute(); err == nil {
		t.Error("expected error for invalid priority")
	}
}

func TestParseTaskType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"bug", "bug", false},
		{"task", "task", false},
		{"hotfix", "hotfix", false},
		{"invalid", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := parseTaskType(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("expected error")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if string(result) != tt.expected {
					t.Errorf("parseTaskType(%s) = %s, want %s", tt.input, result, tt.expected)
				}
			}
		})
	}
}
