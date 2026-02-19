package sdpinit

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	// Create temp directory with prompts/
	tmpDir := t.TempDir()
	promptsDir := filepath.Join(tmpDir, "prompts")
	skillsDir := filepath.Join(promptsDir, "skills")
	agentsDir := filepath.Join(promptsDir, "agents")

	// Create test prompts structure
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	// Create test files
	testSkill := filepath.Join(skillsDir, "test.md")
	if err := os.WriteFile(testSkill, []byte("# Test Skill"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	testAgent := filepath.Join(agentsDir, "test.md")
	if err := os.WriteFile(testAgent, []byte("# Test Agent"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Change to temp directory
	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Run init
	cfg := Config{ProjectType: "go"}
	if err := Run(cfg); err != nil {
		t.Fatalf("Run() failed: %v", err)
	}

	// Check .claude/ directory was created
	claudeDir := ".claude"
	if _, err := os.Stat(claudeDir); os.IsNotExist(err) {
		t.Fatal(".claude/ directory was not created")
	}

	// Check subdirectories
	expectedDirs := []string{"skills", "agents", "validators"}
	for _, dir := range expectedDirs {
		dirPath := filepath.Join(claudeDir, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Subdirectory %s was not created", dir)
		}
	}

	// Check prompts were copied
	copiedSkill := filepath.Join(claudeDir, "skills", "test.md")
	if _, err := os.Stat(copiedSkill); os.IsNotExist(err) {
		t.Error("Test skill was not copied")
	}

	copiedAgent := filepath.Join(claudeDir, "agents", "test.md")
	if _, err := os.Stat(copiedAgent); os.IsNotExist(err) {
		t.Error("Test agent was not copied")
	}

	// Check settings.json was created
	settingsPath := filepath.Join(claudeDir, "settings.json")
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		t.Fatal("settings.json was not created")
	}

	// Check settings.json content
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	settingsStr := string(content)
	if !strings.Contains(settingsStr, `"projectType": "go"`) {
		t.Errorf("settings.json missing projectType: %s", settingsStr)
	}

	if !strings.Contains(settingsStr, `"skills":`) {
		t.Errorf("settings.json missing skills: %s", settingsStr)
	}

	// Check file permissions (0600)
	info, err := os.Stat(settingsPath)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}

	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("settings.json has wrong permissions: got %o, want 0600", perm)
	}
}

func TestRun_NoPromptsDir(t *testing.T) {
	// Create temp directory WITHOUT prompts/
	tmpDir := t.TempDir()

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	originalArchiveURL := os.Getenv("SDP_PROMPTS_ARCHIVE_URL")
	t.Cleanup(func() { _ = os.Setenv("SDP_PROMPTS_ARCHIVE_URL", originalArchiveURL) })
	if err := os.Setenv("SDP_PROMPTS_ARCHIVE_URL", "http://127.0.0.1:1/prompts.zip"); err != nil {
		t.Fatalf("setenv SDP_PROMPTS_ARCHIVE_URL: %v", err)
	}

	// Run init - should fail
	cfg := Config{ProjectType: "go"}
	err := Run(cfg)
	if err == nil {
		t.Fatal("Run() should fail when prompts/ doesn't exist")
	}

	if !strings.Contains(err.Error(), "prompts not found locally") {
		t.Errorf("Wrong error: %v", err)
	}
}

func TestRun_PromptsInSdpSubdirectory(t *testing.T) {
	tmpDir := t.TempDir()
	promptsDir := filepath.Join(tmpDir, "sdp", "prompts")
	skillsDir := filepath.Join(promptsDir, "skills")

	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	testSkill := filepath.Join(skillsDir, "test.md")
	if err := os.WriteFile(testSkill, []byte("# Test Skill"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	cfg := Config{ProjectType: "go"}
	if err := Run(cfg); err != nil {
		t.Fatalf("Run() failed with sdp/prompts fallback: %v", err)
	}

	copiedSkill := filepath.Join(".claude", "skills", "test.md")
	if _, err := os.Stat(copiedSkill); os.IsNotExist(err) {
		t.Error("Test skill was not copied from sdp/prompts fallback")
	}
}

func TestRun_WithSymlinkedClaudeDirs(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink test requires elevated privileges on windows")
	}

	tmpDir := t.TempDir()
	promptsDir := filepath.Join(tmpDir, "sdp", "prompts")
	skillsDir := filepath.Join(promptsDir, "skills")
	agentsDir := filepath.Join(promptsDir, "agents")

	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		t.Fatalf("mkdir skills: %v", err)
	}
	if err := os.MkdirAll(agentsDir, 0755); err != nil {
		t.Fatalf("mkdir agents: %v", err)
	}

	if err := os.WriteFile(filepath.Join(skillsDir, "test.md"), []byte("# skill"), 0644); err != nil {
		t.Fatalf("write skill: %v", err)
	}
	if err := os.WriteFile(filepath.Join(agentsDir, "test.md"), []byte("# agent"), 0644); err != nil {
		t.Fatalf("write agent: %v", err)
	}

	claudeDir := filepath.Join(tmpDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatalf("mkdir .claude: %v", err)
	}

	if err := os.Symlink(filepath.Join("..", "sdp", "prompts", "skills"), filepath.Join(claudeDir, "skills")); err != nil {
		t.Fatalf("symlink skills: %v", err)
	}
	if err := os.Symlink(filepath.Join("..", "sdp", "prompts", "agents"), filepath.Join(claudeDir, "agents")); err != nil {
		t.Fatalf("symlink agents: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	if err := Run(Config{ProjectType: "go"}); err != nil {
		t.Fatalf("Run() failed with symlinked .claude dirs: %v", err)
	}
}

func TestResolvePromptsDir_FromEnvSourceDir(t *testing.T) {
	tmpDir := t.TempDir()
	sourcePrompts := filepath.Join(tmpDir, "external-prompts")
	if err := os.MkdirAll(filepath.Join(sourcePrompts, "skills"), 0755); err != nil {
		t.Fatalf("mkdir skills: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourcePrompts, "skills", "test.md"), []byte("# env source"), 0644); err != nil {
		t.Fatalf("write test prompt: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	originalEnv := os.Getenv("SDP_PROMPTS_SOURCE_DIR")
	t.Cleanup(func() { _ = os.Setenv("SDP_PROMPTS_SOURCE_DIR", originalEnv) })
	if err := os.Setenv("SDP_PROMPTS_SOURCE_DIR", sourcePrompts); err != nil {
		t.Fatalf("setenv: %v", err)
	}

	resolved, err := resolvePromptsDir()
	if err != nil {
		t.Fatalf("resolvePromptsDir() failed: %v", err)
	}

	if resolved != sourcePrompts {
		t.Fatalf("resolved prompts dir mismatch: got %q want %q", resolved, sourcePrompts)
	}
}

func TestResolvePromptsDir_EnvSourceDirMissingSkills(t *testing.T) {
	tmpDir := t.TempDir()
	badSource := filepath.Join(tmpDir, "bad-prompts")
	if err := os.MkdirAll(badSource, 0755); err != nil {
		t.Fatalf("mkdir bad source: %v", err)
	}

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	originalEnv := os.Getenv("SDP_PROMPTS_SOURCE_DIR")
	t.Cleanup(func() { _ = os.Setenv("SDP_PROMPTS_SOURCE_DIR", originalEnv) })
	if err := os.Setenv("SDP_PROMPTS_SOURCE_DIR", badSource); err != nil {
		t.Fatalf("setenv: %v", err)
	}

	_, err := resolvePromptsDir()
	if err == nil {
		t.Fatalf("resolvePromptsDir() should fail for missing skills in env source dir")
	}
}

func TestRun_CreateDirError(t *testing.T) {
	// This tests error handling when directory creation fails
	// We can't easily mock os.MkdirAll, so we'll test the error path
	// by trying to create a directory in an invalid location
	cfg := Config{ProjectType: "go"}

	// Create temp directory
	tmpDir := t.TempDir()

	originalWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(originalWd) })
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	// Create a file named .claude (not a directory)
	if err := os.WriteFile(".claude", []byte("not a directory"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Run init - should fail
	err := Run(cfg)
	if err == nil {
		t.Fatal("Run() should fail when .claude is a file, not directory")
	}
}

func TestCreateSettings(t *testing.T) {
	tmpDir := t.TempDir()

	cfg := Config{
		ProjectType: "python",
		SkipBeads:   false,
	}

	// Create settings
	claudeDir := tmpDir
	if err := createSettings(claudeDir, cfg); err != nil {
		t.Fatalf("createSettings() failed: %v", err)
	}

	// Check file exists
	settingsPath := filepath.Join(claudeDir, "settings.json")
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		t.Fatal("settings.json was not created")
	}

	// Check content
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	settingsStr := string(content)
	if !strings.Contains(settingsStr, `"projectType": "python"`) {
		t.Errorf("Wrong projectType in settings: %s", settingsStr)
	}

	if !strings.Contains(settingsStr, `"sdpVersion": "1.0.0"`) {
		t.Errorf("Missing sdpVersion in settings: %s", settingsStr)
	}

	// Check skills list
	expectedSkills := []string{"feature", "idea", "design", "build", "review", "deploy", "debug", "bugfix", "hotfix", "oneshot"}
	for _, skill := range expectedSkills {
		if !strings.Contains(settingsStr, `"`+skill+`"`) {
			t.Errorf("Missing skill %s in settings: %s", skill, settingsStr)
		}
	}

	// Check permissions
	info, err := os.Stat(settingsPath)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}

	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("Wrong permissions: got %o, want 0600", perm)
	}
}
