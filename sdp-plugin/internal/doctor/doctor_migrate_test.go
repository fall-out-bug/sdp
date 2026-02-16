package doctor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMigrateConfig_NoConfig(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	_, err := MigrateConfig(false)
	if err == nil {
		t.Error("Expected error when no config exists")
	}
}

func TestMigrateConfig_AlreadyLatest(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Create config with latest version
	os.MkdirAll(".sdp", 0755)
	os.WriteFile(".sdp/config.yml", []byte("version: 1\nfoo: bar"), 0644)

	m, err := MigrateConfig(false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !m.Success {
		t.Error("Expected success for already latest version")
	}
	if m.Message != "Config already at latest version" {
		t.Errorf("Unexpected message: %s", m.Message)
	}
}

func TestMigrateConfig_DryRun(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Create config without version (v0)
	os.MkdirAll(".sdp", 0755)
	os.WriteFile(".sdp/config.yml", []byte("foo: bar"), 0644)

	m, err := MigrateConfig(true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if m.FromVersion != 0 {
		t.Errorf("Expected from version 0, got %d", m.FromVersion)
	}
	if m.BackupPath == "" {
		t.Error("Expected backup path to be set")
	}
}

func TestMigrateConfig_Actual(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Create config without version (v0)
	os.MkdirAll(".sdp", 0755)
	os.WriteFile(".sdp/config.yml", []byte("foo: bar"), 0644)

	m, err := MigrateConfig(false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !m.Success {
		t.Errorf("Expected success, got message: %s", m.Message)
	}

	// Verify config was migrated
	content, _ := os.ReadFile(".sdp/config.yml")
	if string(content)[:10] != "version: 1" {
		t.Errorf("Config not migrated correctly: %s", string(content)[:50])
	}
}

func TestRollbackMigration(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Create config and backup
	os.MkdirAll(".sdp/backups", 0755)
	os.WriteFile(".sdp/config.yml", []byte("version: 1\nnew: content"), 0644)
	os.WriteFile(".sdp/backups/config-old.yml", []byte("old: content"), 0644)

	// Rollback
	err := RollbackMigration(".sdp/backups/config-old.yml")
	if err != nil {
		t.Fatalf("Rollback failed: %v", err)
	}

	// Verify rollback
	content, _ := os.ReadFile(".sdp/config.yml")
	if string(content) != "old: content" {
		t.Errorf("Rollback didn't restore correctly: %s", string(content))
	}
}

func TestRollbackMigration_NoBackup(t *testing.T) {
	err := RollbackMigration("/nonexistent/backup.yml")
	if err == nil {
		t.Error("Expected error for nonexistent backup")
	}
}

func TestDetectConfigVersion(t *testing.T) {
	tests := []struct {
		content  string
		expected int
	}{
		{"version: 1\nfoo: bar", 1},
		{"version: 2\nfoo: bar", 2},
		{"foo: bar", 0},
		{"", 0},
		{"# comment\nversion: 3\nfoo: bar", 0}, // Not at start
	}

	for _, tt := range tests {
		result := detectConfigVersion(tt.content)
		if result != tt.expected {
			t.Errorf("detectConfigVersion(%q) = %d, want %d", tt.content, result, tt.expected)
		}
	}
}

func TestGetLatestVersion(t *testing.T) {
	latest := getLatestVersion()
	if latest < 1 {
		t.Errorf("Expected latest version >= 1, got %d", latest)
	}
}

func TestCreateConfigBackup(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	os.MkdirAll(".sdp", 0755)
	os.WriteFile(".sdp/config.yml", []byte("test: content"), 0644)

	backupPath, err := createConfigBackup(".sdp/config.yml")
	if err != nil {
		t.Fatalf("createConfigBackup failed: %v", err)
	}

	// Verify backup was created
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Error("Backup file not created")
	}

	// Verify backup content
	backupContent, _ := os.ReadFile(backupPath)
	if string(backupContent) != "test: content" {
		t.Error("Backup content mismatch")
	}
}

func TestListBackups(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// No backups directory
	backups, err := ListBackups()
	if err != nil {
		t.Fatalf("ListBackups failed: %v", err)
	}
	if len(backups) != 0 {
		t.Errorf("Expected 0 backups, got %d", len(backups))
	}

	// Create backups
	os.MkdirAll(".sdp/backups", 0755)
	os.WriteFile(".sdp/backups/config-20260101-120000.yml", []byte("old"), 0644)
	os.WriteFile(".sdp/backups/config-20260102-120000.yml", []byte("newer"), 0644)
	os.WriteFile(".sdp/backups/other.txt", []byte("not a backup"), 0644)

	backups, err = ListBackups()
	if err != nil {
		t.Fatalf("ListBackups failed: %v", err)
	}
	if len(backups) != 2 {
		t.Errorf("Expected 2 backups, got %d: %v", len(backups), backups)
	}
}

func TestLogMigration(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	os.MkdirAll(".sdp", 0755)

	m := &Migration{
		FromVersion: 0,
		ToVersion:   1,
		Success:     true,
		BackupPath:  ".sdp/backups/config-test.yml",
	}
	m.Timestamp = m.Timestamp.UTC()

	logMigration(m)

	// Verify log was created
	content, err := os.ReadFile(".sdp/migrations.log")
	if err != nil {
		t.Fatalf("Migration log not created: %v", err)
	}

	logStr := string(content)
	if logStr == "" {
		t.Error("Migration log is empty")
	}
}

func TestMigrateV0ToV1(t *testing.T) {
	tmpDir := t.TempDir()

	configPath := filepath.Join(tmpDir, "config.yml")
	os.WriteFile(configPath, []byte("foo: bar"), 0644)

	err := migrateV0ToV1(configPath)
	if err != nil {
		t.Fatalf("migrateV0ToV1 failed: %v", err)
	}

	content, _ := os.ReadFile(configPath)
	if string(content)[:10] != "version: 1" {
		t.Errorf("Config not migrated: %s", string(content))
	}
}
