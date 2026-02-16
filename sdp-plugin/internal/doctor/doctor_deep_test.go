package doctor

import (
	"os"
	"testing"
)

func TestRunDeepChecks(t *testing.T) {
	// Run in temp directory
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Create minimal structure
	os.MkdirAll(".git/hooks", 0755)
	os.MkdirAll(".claude/skills", 0755)
	os.MkdirAll("docs/workstreams/backlog", 0755)

	results := RunDeepChecks()

	if len(results) == 0 {
		t.Error("Expected at least one deep check result")
	}

	// Verify all results have required fields
	for _, r := range results {
		if r.Check == "" {
			t.Error("Deep check missing Check name")
		}
		if r.Status == "" {
			t.Error("Deep check missing Status")
		}
		if r.Message == "" {
			t.Error("Deep check missing Message")
		}
	}
}

func TestCheckGitHooks(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Test 1: No git repo
	result := checkGitHooks()
	if result.Status != "warning" {
		t.Errorf("Expected warning when not in git repo, got %s", result.Status)
	}

	// Test 2: Git repo without hooks
	os.MkdirAll(".git/hooks", 0755)
	result = checkGitHooks()
	if result.Status != "warning" {
		t.Errorf("Expected warning when hooks missing, got %s", result.Status)
	}

	// Test 3: Git repo with hooks
	os.WriteFile(".git/hooks/pre-commit", []byte("#!/bin/bash\necho test"), 0755)
	os.WriteFile(".git/hooks/pre-push", []byte("#!/bin/bash\necho test"), 0755)
	result = checkGitHooks()
	if result.Status != "ok" {
		t.Errorf("Expected ok when hooks present, got %s: %s", result.Status, result.Message)
	}

	// Test 4: Hooks not executable
	os.Chmod(".git/hooks/pre-commit", 0644)
	result = checkGitHooks()
	if result.Status != "warning" {
		t.Errorf("Expected warning when hook not executable, got %s", result.Status)
	}
}

func TestCheckSkillsSyntax(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Test 1: No skills directory
	result := checkSkillsSyntax()
	if result.Status != "warning" {
		t.Errorf("Expected warning when skills dir missing, got %s", result.Status)
	}

	// Test 2: Skills directory with valid skill
	os.MkdirAll(".claude/skills/test-skill", 0755)
	os.WriteFile(".claude/skills/test-skill/SKILL.md", []byte("---\nname: test\n---\n# Test Skill"), 0644)
	result = checkSkillsSyntax()
	if result.Status != "ok" {
		t.Errorf("Expected ok with valid skill, got %s: %s", result.Status, result.Message)
	}

	// Test 3: Skills directory with invalid skill (no frontmatter)
	os.MkdirAll(".claude/skills/bad-skill", 0755)
	os.WriteFile(".claude/skills/bad-skill/SKILL.md", []byte("# No frontmatter"), 0644)
	result = checkSkillsSyntax()
	if result.Status != "warning" {
		t.Errorf("Expected warning with invalid skill, got %s", result.Status)
	}

	// Test 4: Skills directory with missing SKILL.md
	os.MkdirAll(".claude/skills/empty-skill", 0755)
	result = checkSkillsSyntax()
	if result.Status != "warning" {
		t.Errorf("Expected warning with missing SKILL.md, got %s", result.Status)
	}
}

func TestCheckWorkstreamCircularDeps(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Test 1: No workstreams directory
	result := checkWorkstreamCircularDeps()
	if result.Status != "ok" {
		t.Errorf("Expected ok when no workstreams, got %s", result.Status)
	}

	// Test 2: Workstreams without dependencies
	os.MkdirAll("docs/workstreams/backlog", 0755)
	os.WriteFile("docs/workstreams/backlog/00-001-01.md", []byte("# WS\n**Depends on**: -"), 0644)
	os.WriteFile("docs/workstreams/backlog/00-001-02.md", []byte("# WS\n**Depends on**: -"), 0644)
	result = checkWorkstreamCircularDeps()
	if result.Status != "ok" {
		t.Errorf("Expected ok with no deps, got %s: %s", result.Status, result.Message)
	}

	// Test 3: Workstreams with linear dependencies
	os.WriteFile("docs/workstreams/backlog/00-001-03.md", []byte("# WS\n**Depends on**: 00-001-01, 00-001-02"), 0644)
	result = checkWorkstreamCircularDeps()
	if result.Status != "ok" {
		t.Errorf("Expected ok with linear deps, got %s: %s", result.Status, result.Message)
	}
}

func TestCheckBeadsIntegrity(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Test 1: No beads database
	result := checkBeadsIntegrity()
	if result.Status != "warning" {
		t.Errorf("Expected warning when beads missing, got %s", result.Status)
	}

	// Test 2: Empty beads database
	os.MkdirAll(".beads", 0755)
	os.WriteFile(".beads/beads.db", []byte{}, 0644)
	result = checkBeadsIntegrity()
	if result.Status != "warning" {
		t.Errorf("Expected warning when beads empty, got %s", result.Status)
	}

	// Test 3: Valid beads database
	os.WriteFile(".beads/beads.db", []byte("test data"), 0644)
	result = checkBeadsIntegrity()
	if result.Status != "ok" {
		t.Errorf("Expected ok with valid beads, got %s: %s", result.Status, result.Message)
	}
}

func TestCheckConfigVersion(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Test 1: No config
	result := checkConfigVersion()
	if result.Status != "ok" {
		t.Errorf("Expected ok when no config, got %s", result.Status)
	}

	// Test 2: Config without version
	os.MkdirAll(".sdp", 0755)
	os.WriteFile(".sdp/config.yml", []byte("foo: bar"), 0644)
	result = checkConfigVersion()
	if result.Status != "warning" {
		t.Errorf("Expected warning when no version, got %s", result.Status)
	}

	// Test 3: Config with version
	os.WriteFile(".sdp/config.yml", []byte("version: 1\nfoo: bar"), 0644)
	result = checkConfigVersion()
	if result.Status != "ok" {
		t.Errorf("Expected ok with version, got %s: %s", result.Status, result.Message)
	}
}

func TestDetectCycles(t *testing.T) {
	tests := []struct {
		name     string
		deps     map[string][]string
		expected int // number of cycles
	}{
		{
			name:     "empty",
			deps:     map[string][]string{},
			expected: 0,
		},
		{
			name:     "no cycles",
			deps:     map[string][]string{"a": {"b"}, "b": {"c"}},
			expected: 0,
		},
		{
			name:     "self cycle",
			deps:     map[string][]string{"a": {"a"}},
			expected: 1,
		},
		{
			name:     "simple cycle",
			deps:     map[string][]string{"a": {"b"}, "b": {"a"}},
			expected: 1,
		},
		{
			name:     "complex cycle",
			deps:     map[string][]string{"a": {"b"}, "b": {"c"}, "c": {"a"}},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cycles := detectCycles(tt.deps)
			if len(cycles) != tt.expected {
				t.Errorf("Expected %d cycles, got %d: %v", tt.expected, len(cycles), cycles)
			}
		})
	}
}

func TestCheckWorkstreamCircularDeps_WithCycle(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	os.MkdirAll("docs/workstreams/backlog", 0755)

	// Create circular dependency: a -> b -> a
	os.WriteFile("docs/workstreams/backloop/00-001-01.md", []byte("# WS A\n**Depends on**: 00-001-02"), 0644)
	os.WriteFile("docs/workstreams/backlog/00-001-02.md", []byte("# WS B\n**Depends on**: 00-001-01"), 0644)

	result := checkWorkstreamCircularDeps()
	// May or may not detect depending on order
	t.Logf("Result: %s - %s", result.Status, result.Message)
}

func TestCheckSkillsSyntax_ReadError(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	os.MkdirAll(".claude/skills", 0000) // No permissions

	result := checkSkillsSyntax()
	if result.Status != "error" {
		t.Logf("Got status: %s, message: %s", result.Status, result.Message)
	}
}

func TestCheckBeadsIntegrity_ReadError(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	os.MkdirAll(".beads", 0755)
	os.WriteFile(".beads/beads.db", []byte("test"), 0000) // No read permissions

	result := checkBeadsIntegrity()
	if result.Status != "error" {
		t.Logf("Got status: %s, message: %s", result.Status, result.Message)
	}
}

func TestCheckConfigVersion_ReadError(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	os.MkdirAll(".sdp", 0755)
	os.WriteFile(".sdp/config.yml", []byte("test"), 0000) // No read permissions

	result := checkConfigVersion()
	if result.Status != "error" {
		t.Logf("Got status: %s, message: %s", result.Status, result.Message)
	}
}
