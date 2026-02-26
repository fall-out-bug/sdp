package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_LogPathValidationWhenProjectRootEmpty(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".sdp")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(cfgDir, "config.yml")
	// Malicious log_path that would escape when resolved from cwd
	cfgContent := "version: 1\nevidence:\n  enabled: true\n  log_path: \"../../etc/passwd\"\n"
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0644); err != nil {
		t.Fatal(err)
	}
	origWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	// Load with empty projectRoot — should still validate log_path (using cwd)
	_, err := Load("")
	if err == nil {
		t.Error("expected error for log_path outside root when projectRoot empty, got nil")
	}
}

func TestLoad_FileNotExist_ReturnsDefault(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.Version != 1 {
		t.Errorf("expected version 1, got %d", cfg.Version)
	}
	if cfg.Evidence.LogPath != ".sdp/log/events.jsonl" {
		t.Errorf("expected default log path, got %s", cfg.Evidence.LogPath)
	}
}

func TestLoad_InvalidYAML_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".sdp")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(cfgDir, "config.yml")
	if err := os.WriteFile(cfgPath, []byte("invalid: yaml: content: ["), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(dir)
	if err == nil {
		t.Error("expected parse error for invalid YAML")
	}
}

func TestLoad_ValidYAML_MergesWithDefaults(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".sdp")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(cfgDir, "config.yml")
	content := "version: 1\nquality:\n  coverage_threshold: 90\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.Quality.CoverageThreshold != 90 {
		t.Errorf("expected coverage_threshold 90, got %d", cfg.Quality.CoverageThreshold)
	}
	if cfg.Evidence.LogPath != ".sdp/log/events.jsonl" {
		t.Errorf("expected default log path, got %s", cfg.Evidence.LogPath)
	}
}

func TestLoad_InvalidVersion_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".sdp")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(cfgDir, "config.yml")
	if err := os.WriteFile(cfgPath, []byte("version: 0\n"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(dir)
	if err == nil {
		t.Error("expected validation error for version < 1")
	}
}

func TestLoad_InvalidDuration_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".sdp")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(cfgDir, "config.yml")
	content := "version: 1\ntimeouts:\n  verification_command: \"invalid\"\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(dir)
	if err == nil {
		t.Error("expected validation error for invalid duration")
	}
}

func TestLoad_GuardRulesFileOutsideRoot_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".sdp")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(cfgDir, "config.yml")
	content := "version: 1\nguard:\n  rules_file: \"../../etc/guard-rules.yml\"\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(dir)
	if err == nil {
		t.Error("expected error for guard.rules_file outside root")
	}
}

func TestTimeoutFromEnv_NotSet_ReturnsFallback(t *testing.T) {
	os.Unsetenv("SDP_TEST_TIMEOUT")
	d := TimeoutFromEnv("SDP_TEST_TIMEOUT", 30)
	if d != 30 {
		t.Errorf("expected 30, got %v", d)
	}
}

func TestTimeoutFromEnv_Set_ReturnsParsed(t *testing.T) {
	os.Setenv("SDP_TEST_TIMEOUT", "45s")
	defer os.Unsetenv("SDP_TEST_TIMEOUT")
	d := TimeoutFromEnv("SDP_TEST_TIMEOUT", 30)
	if d != 45*1e9 {
		t.Errorf("expected 45s, got %v", d)
	}
}

func TestTimeoutFromEnv_Invalid_ReturnsFallback(t *testing.T) {
	os.Setenv("SDP_TEST_TIMEOUT", "not-a-duration")
	defer os.Unsetenv("SDP_TEST_TIMEOUT")
	d := TimeoutFromEnv("SDP_TEST_TIMEOUT", 30)
	if d != 30 {
		t.Errorf("expected fallback 30, got %v", d)
	}
}

func TestTimeoutFromConfigOrEnv_ConfigFirst(t *testing.T) {
	d := TimeoutFromConfigOrEnv("90s", "SDP_TEST_TIMEOUT", 30)
	if d != 90*1e9 {
		t.Errorf("expected 90s from config, got %v", d)
	}
}

func TestTimeoutFromConfigOrEnv_ConfigInvalid_FallsToEnv(t *testing.T) {
	os.Setenv("SDP_TEST_TIMEOUT", "60s")
	defer os.Unsetenv("SDP_TEST_TIMEOUT")
	d := TimeoutFromConfigOrEnv("invalid", "SDP_TEST_TIMEOUT", 30)
	if d != 60*1e9 {
		t.Errorf("expected 60s from env, got %v", d)
	}
}

func TestTimeoutFromConfigOrEnv_ConfigEmpty_FallsToEnv(t *testing.T) {
	os.Setenv("SDP_TEST_TIMEOUT", "120s")
	defer os.Unsetenv("SDP_TEST_TIMEOUT")
	d := TimeoutFromConfigOrEnv("", "SDP_TEST_TIMEOUT", 30)
	if d != 120*1e9 {
		t.Errorf("expected 120s from env, got %v", d)
	}
}

func TestFindProjectRoot_InProjectDir(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".sdp")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatal(err)
	}

	origWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	root, err := FindProjectRoot()
	if err != nil {
		t.Fatalf("FindProjectRoot failed: %v", err)
	}
	// Normalize for comparison (macOS /var -> /private/var)
	normDir, _ := filepath.EvalSymlinks(dir)
	normRoot, _ := filepath.EvalSymlinks(root)
	if normRoot != normDir {
		t.Errorf("expected root %s, got %s", normDir, normRoot)
	}
}

func TestFindProjectRoot_InSubdir_FindsParent(t *testing.T) {
	dir := t.TempDir()
	cfgDir := filepath.Join(dir, ".sdp")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatal(err)
	}
	subDir := filepath.Join(dir, "sub", "nested")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	origWd, _ := os.Getwd()
	if err := os.Chdir(subDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	root, err := FindProjectRoot()
	if err != nil {
		t.Fatalf("FindProjectRoot failed: %v", err)
	}
	normDir, _ := filepath.EvalSymlinks(dir)
	normRoot, _ := filepath.EvalSymlinks(root)
	if normRoot != normDir {
		t.Errorf("expected root %s, got %s", normDir, normRoot)
	}
}

func TestFindProjectRoot_WithGit_FindsRoot(t *testing.T) {
	dir := t.TempDir()
	gitDir := filepath.Join(dir, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatal(err)
	}
	subDir := filepath.Join(dir, "src")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	origWd, _ := os.Getwd()
	if err := os.Chdir(subDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origWd)

	root, err := FindProjectRoot()
	if err != nil {
		t.Fatalf("FindProjectRoot failed: %v", err)
	}
	normDir, _ := filepath.EvalSymlinks(dir)
	normRoot, _ := filepath.EvalSymlinks(root)
	if normRoot != normDir {
		t.Errorf("expected root %s, got %s", normDir, normRoot)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}
	if cfg.Version != 1 {
		t.Errorf("expected version 1, got %d", cfg.Version)
	}
	if cfg.Evidence.LogPath != ".sdp/log/events.jsonl" {
		t.Errorf("expected default log path, got %s", cfg.Evidence.LogPath)
	}
	if cfg.Guard.RulesFile != ".sdp/guard-rules.yml" {
		t.Errorf("expected default rules file, got %s", cfg.Guard.RulesFile)
	}
}

func TestConfigValidate_InvalidVersion(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Version = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for version < 1")
	}
}

func TestConfigValidate_InvalidTimeout(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Timeouts.VerificationCommand = "not-a-duration"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid timeout")
	}
}
